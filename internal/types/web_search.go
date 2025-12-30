package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// WebSearchConfig represents the web search configuration for a tenant
type WebSearchConfig struct {
	Provider          string   `json:"provider"`           // 검색 엔진 제공자 ID
	APIKey            string   `json:"api_key"`            // API 키(필요한 경우)
	MaxResults        int      `json:"max_results"`        // 최대 검색 결과 수
	IncludeDate       bool     `json:"include_date"`       // 날짜 포함 여부
	CompressionMethod string   `json:"compression_method"` // 압축 방법: none, summary, extract, rag
	Blacklist         []string `json:"blacklist"`          // 블랙리스트 규칙 목록
	// RAG 압축 관련 설정
	EmbeddingModelID   string `json:"embedding_model_id,omitempty"`  // 임베딩 모델 ID(RAG 압축용)
	EmbeddingDimension int    `json:"embedding_dimension,omitempty"` // 임베딩 차원(RAG 압축용)
	RerankModelID      string `json:"rerank_model_id,omitempty"`     // 재정렬 모델 ID(RAG 압축용)
	DocumentFragments  int    `json:"document_fragments,omitempty"`  // 문서 조각 수(RAG 압축용)
}

// Value implements driver.Valuer interface for WebSearchConfig
func (c WebSearchConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan implements sql.Scanner interface for WebSearchConfig
func (c *WebSearchConfig) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, c)
}

// WebSearchResult represents a single web search result
type WebSearchResult struct {
	Title       string     `json:"title"`                  // 검색 결과 제목
	URL         string     `json:"url"`                    // 결과 URL
	Snippet     string     `json:"snippet"`                // 요약 스니펫
	Content     string     `json:"content"`                // 전체 내용(선택 사항, 추가 수집 필요)
	Source      string     `json:"source"`                 // 출처(예: duckduckgo 등)
	PublishedAt *time.Time `json:"published_at,omitempty"` // 게시 시각(있는 경우)
}

// WebSearchProviderInfo represents information about a web search provider
type WebSearchProviderInfo struct {
	ID             string `json:"id"`                // 제공자 ID
	Name           string `json:"name"`              // 제공자 이름
	Free           bool   `json:"free"`              // 무료 여부
	RequiresAPIKey bool   `json:"requires_api_key"`  // API 키 필요 여부
	Description    string `json:"description"`       // 설명
	APIURL         string `json:"api_url,omitempty"` // API 주소(선택)
}
