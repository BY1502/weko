import torch
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel, Field
from transformers import AutoModelForSequenceClassification, AutoTokenizer
from typing import List

# --- 1. API 요청/응답 데이터 구조 정의 ---

# 요청 본문 구조는 그대로 유지
class RerankRequest(BaseModel):
    query: str
    documents: List[str]

# --- 수정 시작: 테스트용 응답 구조 정의, 필드명은 "score" ---

# DocumentInfo 구조는 그대로 유지
class DocumentInfo(BaseModel):
    text: str

# 기존 GoRankResult 를 TestRankResult 로 변경
# 핵심 변경: "relevance_score" 필드를 "score"로 이름 변경
class TestRankResult(BaseModel):
    index: int
    document: DocumentInfo
    score: float  # <--- 핵심 변경: 필드명이 relevance_score → score

# 최종 응답 구조, "results" 리스트에 TestRankResult 포함
class TestFinalResponse(BaseModel):
    results: List[TestRankResult]

# --- 수정 종료 ---


# --- 2. 모델 로드(서비스 시작 시 한 번 실행) ---
print("모델을 로드하는 중입니다. 잠시만 기다려주세요...")
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print(f"사용 중인 디바이스: {device}")
try:
    # 여기 경로가 올바른지 확인하세요
    model_path = '/data1/home/lwx/work/Download/rerank_model_weight'
    tokenizer = AutoTokenizer.from_pretrained(model_path)
    model = AutoModelForSequenceClassification.from_pretrained(model_path)
    model.to(device)
    model.eval()
    print("모델 로드 성공!")
except Exception as e:
    print(f"모델 로드 실패: {e}")
    # 테스트 환경에서 모델 로드 실패 시 잘못된 서비스를 실행하지 않도록 종료
    exit()

# --- 3. FastAPI 앱 생성 ---
app = FastAPI(
    title="Reranker API (Test Version)",
    description="'score' 필드를 반환해 Go 클라이언트 호환성을 테스트하는 API 서비스",
    version="1.0.1"
)

# --- 4. API 엔드포인트 정의 ---
# --- 수정 시작: response_model 을 새로운 테스트 응답 구조로 지정 ---
@app.post("/rerank", response_model=TestFinalResponse) # <--- 【핵심 변경】response_model을 TestFinalResponse로
def rerank_endpoint(request: RerankRequest):
    # --- 수정 종료 ---

    pairs = [[request.query, doc] for doc in request.documents]

    with torch.no_grad():
        inputs = tokenizer(pairs, padding=True, truncation=True, return_tensors='pt', max_length=1024).to(device)
        scores = model(**inputs, return_dict=True).logits.view(-1, ).float()

    # --- 수정 시작: 테스트용 구조에 맞게 결과 생성 ---
    results = []
    for i, (text, score_val) in enumerate(zip(request.documents, scores)):
        
        # 1. 중첩 document 객체 생성
        doc_info = DocumentInfo(text=text)
        
        # 2. TestRankResult 객체 생성
        #    필드명: index, document, score
        test_result = TestRankResult(
            index=i,
            document=doc_info,
            score=score_val.item()  # <--- 【핵심 변경】"score" 필드에 값 할당
        )
        results.append(test_result)

    # 3. 정렬(key를 score로 변경)
    sorted_results = sorted(results, key=lambda x: x.score, reverse=True)
    # --- 수정 종료 ---
    
    # 딕셔너리를 반환하면 FastAPI가 response_model(TestFinalResponse)에 따라 검증/직렬화함
    # 최종 JSON: {"results": [{"index": ..., "document": ..., "score": ...}]}
    return {"results": sorted_results}

@app.get("/")
def read_root():
    return {"status": "Reranker API (Test Version) is running"}

# --- 5. 서비스 시작 ---
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
    
