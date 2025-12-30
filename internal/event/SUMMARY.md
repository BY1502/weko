# WeKnora 이벤트 시스템 요약

## 개요

WeKnora 프로젝트에 대한 이벤트 송신/수신 메커니즘을 구축해, 사용자 질의 처리 흐름의 각 단계를 이벤트로 처리할 수 있습니다.

## 핵심 기능

### ✅ 구현 완료

1. **이벤트 버스 (EventBus)**
   - `Emit(ctx, event)` - 이벤트 전송
   - `On(eventType, handler)` - 이벤트 리스너 등록
   - `Off(eventType)` - 리스너 제거
   - `EmitAndWait(ctx, event)` - 이벤트 전송 후 모든 핸들러 완료 대기
   - 동기/비동기 모드 지원

2. **이벤트 유형**
   - 질의 처리(수신, 검증, 전처리, 재작성)
   - 검색(시작, 벡터 검색, 키워드 검색, 엔티티 검색, 완료)
   - 정렬(시작, 완료)
   - 병합(시작, 완료)
   - 채팅 생성 이벤트(시작, 완료, 스트리밍 출력)
   - 에러 이벤트

3. **이벤트 데이터 구조**
   - `QueryData` - 질의 데이터
   - `RetrievalData` - 검색 데이터
   - `RerankData` - 재정렬 데이터
   - `MergeData` - 병합 데이터
   - `ChatData` - 채팅 데이터
   - `ErrorData` - 에러 데이터

4. **미들웨어 지원**
   - `WithLogging()` - 로그 미들웨어
   - `WithTiming()` - 타이밍 미들웨어
   - `WithRecovery()` - 에러 복구 미들웨어
   - `Chain()` - 미들웨어 체이닝

5. **글로벌 이벤트 버스**
   - 싱글턴 글로벌 이벤트 버스
   - 글로벌 편의 함수(`On`, `Emit`, `EmitAndWait` 등)

6. **예제 및 테스트**
   - 완전한 단위 테스트
   - 성능 벤치마크 테스트
   - 전체 사용 예제
   - 실제 시나리오 데모

## 파일 구조

```
internal/event/
├── event.go                    # 코어 이벤트 버스 구현
├── event_data.go              # 이벤트 데이터 구조 정의
├── middleware.go              # 미들웨어 구현
├── global.go                  # 글로벌 이벤트 버스
├── integration_example.go     # 통합 예제(모니터링, 분석 핸들러)
├── example_test.go            # 테스트 및 예제
├── demo/
│   └── main.go               # 완전한 RAG 플로우 데모
├── README.md                 # 상세 문서
├── usage_example.md          # 사용 예제 문서
└── SUMMARY.md                # 본 문서
```

## 성능 지표

- **이벤트 전송 성능**: ~9ns/건 (벤치마크 기준)
- **동시성 안전**: `sync.RWMutex`로 스레드 안전 보장
- **메모리 오버헤드**: 매우 적음, 핸들러 함수 참조만 저장

## 사용 시나리오

### 1. 모니터링 및 지표 수집

```go
bus.On(event.EventRetrievalComplete, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.RetrievalData)
    // Prometheus 등 모니터링 시스템으로 전송
    metricsCollector.RecordRetrievalDuration(data.Duration)
    return nil
})
```

### 2. 로그 기록

```go
bus.On(event.EventQueryRewritten, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.QueryData)
    logger.Infof(ctx, "Query rewritten: %s -> %s", 
        data.OriginalQuery, data.RewrittenQuery)
    return nil
})
```

### 3. 사용자 행동 분석

```go
bus.On(event.EventQueryReceived, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.QueryData)
    // 분석 플랫폼으로 전송
    analytics.TrackQuery(data.UserID, data.OriginalQuery)
    return nil
})
```

### 4. 에러 추적

```go
bus.On(event.EventError, func(ctx context.Context, e event.Event) error {
    data := e.Data.(event.ErrorData)
    // 에러 추적 시스템으로 전송
    sentry.CaptureException(data.Error)
    return nil
})
```

## 통합 방식

### 1단계: 이벤트 시스템 초기화

애플리케이션 시작 시(`main.go`, `container.go` 등):

```go
import "github.com/Tencent/WeKnora/internal/event"

func Initialize() {
    // 글로벌 이벤트 버스 가져오기
    bus := event.GetGlobalEventBus()
    
    // 모니터링 및 분석 설정
    event.NewMonitoringHandler(bus)
    event.NewAnalyticsHandler(bus)
}
```

### 2단계: 처리 단계별 이벤트 전송

질의 처리 플러그인마다 이벤트 전송 추가:

```go
// search.go 에서
event.Emit(ctx, event.NewEvent(event.EventRetrievalStart, event.RetrievalData{
    Query:           chatManage.ProcessedQuery,
    KnowledgeBaseID: chatManage.KnowledgeBaseID,
    TopK:            chatManage.EmbeddingTopK,
}).WithSessionID(chatManage.SessionID))

// rerank.go 에서
event.Emit(ctx, event.NewEvent(event.EventRerankComplete, event.RerankData{
    Query:       chatManage.ProcessedQuery,
    InputCount:  len(chatManage.SearchResult),
    OutputCount: len(rerankResults),
    Duration:    time.Since(startTime).Milliseconds(),
}).WithSessionID(chatManage.SessionID))
```

### 3단계: 사용자 정의 핸들러 등록

필요에 따라 커스텀 핸들러 등록:

```go
event.On(event.EventQueryRewritten, func(ctx context.Context, e event.Event) error {
    // 커스텀 처리 로직
    return nil
})
```

## 장점

1. **낮은 결합도**: 이벤트 송신자와 리스너가 완전 분리되어 유지보수/확장에 용이
2. **고성능**: 매우 낮은 오버헤드(~9ns/건)
3. **유연성**: 동기/비동기, 단일/다중 리스너 지원
4. **확장성**: 새 이벤트 유형과 핸들러 추가가 용이
5. **타입 안전**: 사전 정의된 이벤트 데이터 구조
6. **미들웨어 지원**: 로깅/타이밍/에러 처리 등 횡단 관심사 추가 용이
7. **테스트 친화적**: 테스트에서 이벤트 동작 검증이 쉬움

## 테스트 결과

✅ 모든 단위 테스트 통과
✅ 성능 테스트 통과(~9ns/건)
✅ 비동기 처리 테스트 통과
✅ 다중 핸들러 테스트 통과
✅ 전체 플로우 데모 성공

## 후속 제안

### 선택적 개선 기능

1. **이벤트 영속화**: 중요 이벤트를 DB/메시지 큐에 저장
2. **이벤트 재생**: 이벤트 재생으로 디버깅/분석 지원
3. **이벤트 필터링**: 더 복잡한 이벤트 필터링과 라우팅 지원
4. **우선순위 큐**: 이벤트 우선순위 처리 지원
5. **분산 이벤트**: 메시지 큐로 서비스 간 이벤트 지원

### 통합 제안

1. **모니터링 통합**: Prometheus와 연동해 지표 수집
2. **로그 통합**: 통합된 구조화 로그 기록
3. **트레이싱 통합**: 기존 트레이싱 시스템과 통합
4. **알림 통합**: 이벤트 기반 알림 메커니즘

## 출력 예시

`go run ./internal/event/demo/main.go`를 실행하면 전체 RAG 플로우 이벤트 출력을 볼 수 있습니다:

```
Step 1: Query Received
[MONITOR] Query received - Session: session-xxx, Query: RAG 기술이란?
[ANALYTICS] Query tracked - User: user-123, Session: session-xxx

Step 2: Query Rewriting
[MONITOR] Query rewrite started
[MONITOR] Query rewritten - Original: RAG 기술이란?, Rewritten: Retrieval-Augmented Generation 기술...
[CUSTOM] Query Transformation: ...

Step 3: Vector Retrieval
[MONITOR] Retrieval started - Type: vector, TopK: 20
[MONITOR] Retrieval completed - Results: 18, Duration: 301ms
[CUSTOM] Retrieval Efficiency: Rate: 90.00%

Step 4: Result Reranking
[MONITOR] Rerank started - Input: 18
[MONITOR] Rerank completed - Output: 5, Duration: 201ms
[CUSTOM] Rerank Statistics: Reduction: 72.22%

Step 5: Chat Completion
[MONITOR] Chat generation started
[MONITOR] Chat generation completed - Tokens: 256, Duration: 801ms
[ANALYTICS] Chat metrics - Model: gpt-4, Tokens: 256
```

## 요약

이벤트 시스템은 구현·테스트를 마쳤으며, WeKnora에 바로 통합해 모니터링, 로그 기록, 분석, 디버깅에 활용할 수 있습니다. 설계가 단순하고 성능이 뛰어나 사용과 확장이 용이합니다.

