package tools

// Tool names constants
const (
	ToolThinking            = "thinking"
	ToolTodoWrite           = "todo_write"
	ToolGrepChunks          = "grep_chunks"
	ToolKnowledgeSearch     = "knowledge_search"
	ToolListKnowledgeChunks = "list_knowledge_chunks"
	ToolQueryKnowledgeGraph = "query_knowledge_graph"
	ToolGetDocumentInfo     = "get_document_info"
	ToolDatabaseQuery       = "database_query"
	ToolDataAnalysis        = "data_analysis"
	ToolDataSchema          = "data_schema"
	ToolWebSearch           = "web_search"
	ToolWebFetch            = "web_fetch"
)

// AvailableTool defines a simple tool metadata used by settings APIs.
type AvailableTool struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// AvailableToolDefinitions returns the list of tools exposed to the UI.
// Keep this in sync with registered tools in this package.
func AvailableToolDefinitions() []AvailableTool {
	return []AvailableTool{
		{Name: ToolThinking, Label: "생각하기", Description: "동적·반성적 문제 해결을 위한 사고 도구"},
		{Name: ToolTodoWrite, Label: "계획 수립", Description: "구조화된 리서치 계획을 작성"},
		{Name: ToolGrepChunks, Label: "키워드 검색", Description: "특정 키워드를 포함한 문서/분할을 빠르게 찾기"},
		{Name: ToolKnowledgeSearch, Label: "의미 검색", Description: "질문을 이해하고 의미적으로 연관된 내용을 탐색"},
		{Name: ToolListKnowledgeChunks, Label: "문서 분할 보기", Description: "문서의 전체 분할 콘텐츠를 확인"},
		{Name: ToolQueryKnowledgeGraph, Label: "지식 그래프 조회", Description: "지식 그래프에서 관계를 조회"},
		{Name: ToolGetDocumentInfo, Label: "문서 정보 조회", Description: "문서 메타데이터 확인"},
		{Name: ToolDatabaseQuery, Label: "데이터베이스 쿼리", Description: "데이터베이스 정보를 조회"},
		{Name: ToolDataAnalysis, Label: "데이터 분석", Description: "데이터 파일을 이해하고 분석"},
		{Name: ToolDataSchema, Label: "데이터 메타정보", Description: "테이블 파일의 메타정보 조회"},
	}
}

// DefaultAllowedTools returns the default allowed tools list.
func DefaultAllowedTools() []string {
	return []string{
		ToolThinking,
		ToolTodoWrite,
		ToolKnowledgeSearch,
		ToolGrepChunks,
		ToolListKnowledgeChunks,
		ToolQueryKnowledgeGraph,
		ToolGetDocumentInfo,
		ToolDatabaseQuery,
		ToolDataAnalysis,
		ToolDataSchema,
	}
}
