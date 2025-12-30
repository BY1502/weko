# 자주 묻는 질문 (FAQ)

## 1. 로그는 어떻게 확인하나요?

```bash
docker compose logs -f app docreader postgres
```

# 서비스 시작

```bash
./scripts/start_all.sh

# 서비스 중지
./scripts/start_all.sh --stop

# 데이터베이스 초기화 (주의: 데이터 삭제됨)
./scripts/start_all.sh --stop && make clean-db
```

```bash
# LLM 모델 이름
INIT_LLM_MODEL_NAME=your_llm_model
# 임베딩 모델 이름
INIT_EMBEDDING_MODEL_NAME=your_embedding_model
# 임베딩 모델 차원 수 (예: 768, 1024 등)
INIT_EMBEDDING_MODEL_DIMENSION=your_embedding_model_dimension
# 임베딩 모델 ID (보통 문자열)
INIT_EMBEDDING_MODEL_ID=your_embedding_model_id
```

# LLM 모델 접속 주소 및 API 키 설정

```bash
INIT_LLM_MODEL_BASE_URL=your_llm_model_base_url
# LLM 모델 API 키
INIT_LLM_MODEL_API_KEY=your_llm_model_api_key
# 임베딩 모델 접속 주소
INIT_EMBEDDING_MODEL_BASE_URL=your_embedding_model_base_url
# 임베딩 모델 API 키
INIT_EMBEDDING_MODEL_API_KEY=your_embedding_model_api_key

# Rerank 모델 이름
INIT_RERANK_MODEL_NAME=your_rerank_model_name
# Rerank 모델 접속 주소
INIT_RERANK_MODEL_BASE_URL=your_rerank_model_base_url
# Rerank 모델 API 키
INIT_RERANK_MODEL_API_KEY=your_rerank_model_api_key
```

# VLM(시각-언어) 모델 이름 설정

```bash
VLM_MODEL_NAME=your_vlm_model_name

# VLM 모델 접속 주소
VLM_MODEL_BASE_URL=your_vlm_model_base_url

# VLM 모델 API 키
VLM_MODEL_API_KEY=your_vlm_model_api_key
```

# 텐센트 클라우드 COS Secret ID

```bash
COS_SECRET_ID=your_cos_secret_id

# 텐센트 클라우드 COS Secret Key
COS_SECRET_KEY=your_cos_secret_key

# COS 리전 (예: ap-guangzhou)
COS_REGION=your_cos_region

# COS 버킷 이름
COS_BUCKET_NAME=your_cos_bucket_name

# COS 앱 ID
COS_APP_ID=your_cos_app_id

# 파일 저장 경로 접두사
COS_PATH_PREFIX=your_cos_path_prefix
```

```bash
중요: COS 내 파일 권한을 반드시 **'공개 읽기(Public Read)'**로 설정해야 문서 파싱 모듈이 정상적으로 파일을 읽을 수 있습니다.

문서 파싱 모듈(docreader) 로그를 확인하여 OCR과 캡션(Caption) 생성이 정상적으로 되는지 확인하세요.

5. 데이터 분석 기능은 어떻게 사용하나요?
데이터 분석 기능을 사용하려면 에이전트 도구 설정이 필요합니다:

지능형 추론 (Agent 모드): 도구 설정에서 다음 두 가지를 체크하세요.

데이터 메타 정보 조회 (View Data Metadata)

데이터 분석 (Data Analysis)

빠른 문답 (Chat 모드): 별도 도구 선택 없이 바로 간단한 데이터 질의가 가능합니다.

주의사항 및 사용 규범
지원 파일 형식

현재 CSV (.csv) 및 Excel (.xlsx, .xls) 형식만 지원합니다.

복잡한 Excel 파일의 경우 읽기 실패 시 표준 CSV 형식으로 변환하여 업로드하는 것을 권장합니다.

쿼리 제한

읽기 전용 쿼리만 지원합니다: SELECT, SHOW, DESCRIBE, EXPLAIN, PRAGMA 등.

데이터 수정 작업(INSERT, UPDATE, DELETE, CREATE, DROP 등)은 금지됩니다.

P.S.
위 방법으로 해결되지 않는 경우, Issue에 문제 상황과 관련 로그 정보를 남겨주시면 문제 해결을 도와드리겠습니다.
```
