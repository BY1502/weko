package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
)

var listKnowledgeChunksTool = BaseTool{
	name: ToolListKnowledgeChunks,
	description: `Retrieve full chunk content for a document by knowledge_id.

## Use After grep_chunks or knowledge_search:
1. grep_chunks(["keyword", "variant"]) → get knowledge_id  
2. list_knowledge_chunks(knowledge_id) → read full content

## When to Use:
- Need full content of chunks from a known document
- Want to see context around specific chunks
- Check how many chunks a document has

## Parameters:
- knowledge_id (required): Document ID
- limit (optional): Chunks per page (default 20, max 100)
- offset (optional): Start position (default 0)

## Output:
Full chunk content with chunk_id, chunk_index, and content text.`,
	schema: json.RawMessage(`{
  "type": "object",
  "properties": {
    "knowledge_id": {
      "type": "string",
      "description": "Document ID to retrieve chunks from"
    },
    "limit": {
      "type": "integer",
      "description": "Chunks per page (default 20, max 100)",
      "default": 20,
      "minimum": 1,
      "maximum": 100
    },
    "offset": {
      "type": "integer",
      "description": "Start position (default 0)",
      "default": 0,
      "minimum": 0
    }
  },
  "required": ["knowledge_id", "limit", "offset"]
}`),
}

// ListKnowledgeChunksInput defines the input parameters for list knowledge chunks tool
type ListKnowledgeChunksInput struct {
	KnowledgeID string `json:"knowledge_id"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
}

// ListKnowledgeChunksTool retrieves chunk snapshots for a specific knowledge document.
type ListKnowledgeChunksTool struct {
	BaseTool
	chunkService     interfaces.ChunkService
	knowledgeService interfaces.KnowledgeService
}

// NewListKnowledgeChunksTool creates a new tool instance.
func NewListKnowledgeChunksTool(
	knowledgeService interfaces.KnowledgeService,
	chunkService interfaces.ChunkService,
) *ListKnowledgeChunksTool {
	return &ListKnowledgeChunksTool{
		BaseTool:         listKnowledgeChunksTool,
		chunkService:     chunkService,
		knowledgeService: knowledgeService,
	}
}

// Execute performs the chunk fetch against the chunk service.
func (t *ListKnowledgeChunksTool) Execute(ctx context.Context, args json.RawMessage) (*types.ToolResult, error) {
	tenantID := uint64(0)
	if tid, ok := ctx.Value(types.TenantIDContextKey).(uint64); ok {
		tenantID = tid
	}

	// Parse args from json.RawMessage
	var input ListKnowledgeChunksInput
	if err := json.Unmarshal(args, &input); err != nil {
		return &types.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse args: %v", err),
		}, err
	}

	knowledgeID := input.KnowledgeID
	ok := knowledgeID != ""
	if !ok || strings.TrimSpace(knowledgeID) == "" {
		return &types.ToolResult{
			Success: false,
			Error:   "knowledge_id is required",
		}, fmt.Errorf("knowledge_id is required")
	}
	knowledgeID = strings.TrimSpace(knowledgeID)

	chunkLimit := 20
	if input.Limit > 0 {
		chunkLimit = input.Limit
	}
	offset := 0
	if input.Offset > 0 {
		offset = input.Offset
	}
	if offset < 0 {
		offset = 0
	}

	pagination := &types.Pagination{
		Page:     offset/chunkLimit + 1,
		PageSize: chunkLimit,
	}

	chunks, total, err := t.chunkService.GetRepository().ListPagedChunksByKnowledgeID(ctx,
		tenantID, knowledgeID, pagination, []types.ChunkType{types.ChunkTypeText, types.ChunkTypeFAQ}, "", "", "", "", "")
	if err != nil {
		return &types.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("failed to list chunks: %v", err),
		}, err
	}
	if chunks == nil {
		return &types.ToolResult{
			Success: false,
			Error:   "chunk query returned no data",
		}, fmt.Errorf("chunk query returned no data")
	}

	totalChunks := total
	fetched := len(chunks)

	knowledgeTitle := t.lookupKnowledgeTitle(ctx, knowledgeID)

	output := t.buildOutput(knowledgeID, knowledgeTitle, totalChunks, fetched, chunks)

	formattedChunks := make([]map[string]interface{}, 0, len(chunks))
	for idx, c := range chunks {
		chunkData := map[string]interface{}{
			"seq":             idx + 1,
			"chunk_id":        c.ID,
			"chunk_index":     c.ChunkIndex,
			"content":         c.Content,
			"chunk_type":      c.ChunkType,
			"knowledge_id":    c.KnowledgeID,
			"knowledge_base":  c.KnowledgeBaseID,
			"start_at":        c.StartAt,
			"end_at":          c.EndAt,
			"parent_chunk_id": c.ParentChunkID,
		}

		// 이미지 정보 추가
		if c.ImageInfo != "" {
			var imageInfos []types.ImageInfo
			if err := json.Unmarshal([]byte(c.ImageInfo), &imageInfos); err == nil && len(imageInfos) > 0 {
				imageList := make([]map[string]string, 0, len(imageInfos))
				for _, img := range imageInfos {
					imgData := make(map[string]string)
					if img.URL != "" {
						imgData["url"] = img.URL
					}
					if img.Caption != "" {
						imgData["caption"] = img.Caption
					}
					if img.OCRText != "" {
						imgData["ocr_text"] = img.OCRText
					}
					if len(imgData) > 0 {
						imageList = append(imageList, imgData)
					}
				}
				if len(imageList) > 0 {
					chunkData["images"] = imageList
				}
			}
		}

		formattedChunks = append(formattedChunks, chunkData)
	}

	return &types.ToolResult{
		Success: true,
		Output:  output,
		Data: map[string]interface{}{
			"knowledge_id":    knowledgeID,
			"knowledge_title": knowledgeTitle,
			"total_chunks":    totalChunks,
			"fetched_chunks":  fetched,
			"page":            pagination.Page,
			"page_size":       pagination.PageSize,
			"chunks":          formattedChunks,
		},
	}, nil
}

// lookupKnowledgeTitle looks up the title of a knowledge document
func (t *ListKnowledgeChunksTool) lookupKnowledgeTitle(ctx context.Context, knowledgeID string) string {
	if t.knowledgeService == nil {
		return ""
	}
	knowledge, err := t.knowledgeService.GetKnowledgeByID(ctx, knowledgeID)
	if err != nil || knowledge == nil {
		return ""
	}
	return strings.TrimSpace(knowledge.Title)
}

// buildOutput builds the output for the list knowledge chunks tool
func (t *ListKnowledgeChunksTool) buildOutput(
	knowledgeID string,
	knowledgeTitle string,
	total int64,
	fetched int,
	chunks []*types.Chunk,
) string {
	builder := &strings.Builder{}
	builder.WriteString("=== 지식 문서 분할 ===\n\n")

	if knowledgeTitle != "" {
		fmt.Fprintf(builder, "문서: %s (%s)\n", knowledgeTitle, knowledgeID)
	} else {
		fmt.Fprintf(builder, "문서 ID: %s\n", knowledgeID)
	}
	fmt.Fprintf(builder, "총 분할 수: %d\n", total)

	if fetched == 0 {
		builder.WriteString("분할을 찾지 못했습니다. 문서 파싱이 완료되었는지 확인하세요.\n")
		if total > 0 {
			builder.WriteString("문서는 존재하지만 현재 페이지가 비어 있습니다. 페이지 매개변수를 확인하세요.\n")
		}
		return builder.String()
	}
	fmt.Fprintf(
		builder,
		"이번 조회: %d 개, 범위: %d - %d\n\n",
		fetched,
		chunks[0].ChunkIndex,
		chunks[len(chunks)-1].ChunkIndex,
	)

	builder.WriteString("=== 분할 콘텐츠 미리보기 ===\n\n")
	for idx, c := range chunks {
		fmt.Fprintf(builder, "Chunk #%d (Index %d)\n", idx+1, c.ChunkIndex+1)
		fmt.Fprintf(builder, "  chunk_id: %s\n", c.ID)
		fmt.Fprintf(builder, "  유형: %s\n", c.ChunkType)
		fmt.Fprintf(builder, "  내용: %s\n", summarizeContent(c.Content))

		// 관련 이미지 정보 출력
		if c.ImageInfo != "" {
			var imageInfos []types.ImageInfo
			if err := json.Unmarshal([]byte(c.ImageInfo), &imageInfos); err == nil && len(imageInfos) > 0 {
				fmt.Fprintf(builder, "  관련 이미지 (%d):\n", len(imageInfos))
				for imgIdx, img := range imageInfos {
					fmt.Fprintf(builder, "    이미지 %d:\n", imgIdx+1)
					if img.URL != "" {
						fmt.Fprintf(builder, "      URL: %s\n", img.URL)
					}
					if img.Caption != "" {
						fmt.Fprintf(builder, "      설명: %s\n", img.Caption)
					}
					if img.OCRText != "" {
						fmt.Fprintf(builder, "      OCR 텍스트: %s\n", img.OCRText)
					}
				}
			}
		}
		builder.WriteString("\n")
	}

	if int64(fetched) < total {
		builder.WriteString("알림: 더 많은 분할이 있습니다. offset을 조정하거나 여러 번 호출하세요.\n")
	}

	return builder.String()
}

// summarizeContent summarizes the content of a chunk
func summarizeContent(content string) string {
	cleaned := strings.TrimSpace(content)
	if cleaned == "" {
		return "(내용 없음)"
	}

	return strings.TrimSpace(string(cleaned))
}
