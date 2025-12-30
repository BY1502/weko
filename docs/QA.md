# 자주 묻는 질문(FAQ)

## 1. 로그는 어떻게 보나요?
```bash
docker compose logs -f app docreader postgres
```

## 2. 서비스를 시작/중지하는 방법은?
```bash
# 시작
./scripts/start_all.sh

# 중지
./scripts/start_all.sh --stop

# DB 초기화
./scripts/start_all.sh --stop && make clean-db
```

## 3. 서비스 기동 후 문서 업로드가 되지 않나요?

대부분 Embedding/LLM 모델 설정이 누락된 경우입니다. 다음을 확인하세요.

1. `.env`의 모델 정보가 모두 채워졌는지 확인. 로컬 ollama를 사용할 경우 서비스가 실행 중인지 확인하고 아래 변수를 올바르게 설정합니다:
```bash
# LLM Model
INIT_LLM_MODEL_NAME=your_llm_model
# Embedding Model
INIT_EMBEDDING_MODEL_NAME=your_embedding_model
# Embedding 모델 벡터 차원
INIT_EMBEDDING_MODEL_DIMENSION=your_embedding_model_dimension
# Embedding 모델 ID(일반적으로 문자열)
INIT_EMBEDDING_MODEL_ID=your_embedding_model_id
```

원격 API를 통해 모델에 접근한다면 `BASE_URL`과 `API_KEY`도 설정해야 합니다:
```bash
# LLM 모델 접근 주소
INIT_LLM_MODEL_BASE_URL=your_llm_model_base_url
# LLM 모델 API 키(인증 필요 시 설정)
INIT_LLM_MODEL_API_KEY=your_llm_model_api_key
# Embedding 모델 접근 주소
INIT_EMBEDDING_MODEL_BASE_URL=your_embedding_model_base_url
# Embedding 모델 API 키(인증 필요 시 설정)
INIT_EMBEDDING_MODEL_API_KEY=your_embedding_model_api_key
```

Rerank 기능이 필요하면 Rerank 모델도 설정합니다:
```bash
# 사용 중인 Rerank 모델 이름
INIT_RERANK_MODEL_NAME=your_rerank_model_name
# Rerank 모델 접근 주소
INIT_RERANK_MODEL_BASE_URL=your_rerank_model_base_url
# Rerank 모델 API 키(인증 필요 시 설정)
INIT_RERANK_MODEL_API_KEY=your_rerank_model_api_key
```

2. 메인 서비스 로그에 `ERROR`가 있는지 확인하세요.

## 4. 멀티모달 기능을 켜려면?
1. `.env`에서 아래 설정을 맞춥니다:
```bash
# VLM_MODEL_NAME 사용 중인 멀티모달 모델 이름
VLM_MODEL_NAME=your_vlm_model_name

# VLM_MODEL_BASE_URL 사용 중인 멀티모달 모델 접근 주소
VLM_MODEL_BASE_URL=your_vlm_model_base_url

# VLM_MODEL_API_KEY 사용 중인 멀티모달 모델 API 키
VLM_MODEL_API_KEY=your_vlm_model_api_key
```
참고: 멀티모달 모델은 현재 remote API만 지원하므로 `VLM_MODEL_BASE_URL`, `VLM_MODEL_API_KEY`가 필요합니다.

2. 파싱된 파일을 COS에 업로드해야 하므로 `.env`의 COS 설정을 올바르게 입력합니다:
```bash
# 텐센트 COS 접근 키 ID
COS_SECRET_ID=your_cos_secret_id

# 텐센트 COS 비밀 키
COS_SECRET_KEY=your_cos_secret_key

# 텐센트 COS 리전(예: ap-guangzhou)
COS_REGION=your_cos_region

# 텐센트 COS 버킷 이름
COS_BUCKET_NAME=your_cos_bucket_name

# 텐센트 COS 앱 ID
COS_APP_ID=your_cos_app_id

# 텐센트 COS 경로 프리픽스(파일 저장용)
COS_PATH_PREFIX=your_cos_path_prefix
```
중요: COS 파일 권한을 반드시 **공개 읽기**로 설정하세요. 그렇지 않으면 문서 파싱 모듈이 정상 동작하지 않습니다.

3. 문서 파싱 모듈 로그에서 OCR/Caption이 올바르게 파싱·출력되는지 확인하세요.


## 5. 데이터 분석 기능 사용법

데이터 분석 기능을 사용하기 전에 에이전트에 관련 도구가 설정되어 있는지 확인하세요.

1. **지능형 추론**：도구 설정에서 다음 두 도구를 선택합니다：
   - 데이터 메타정보 보기
   - 데이터 분석

2. **빠른 질의응답 에이전트**：별도 도구 선택 없이 간단한 데이터 조회를 바로 사용할 수 있습니다.

### 주의사항 및 사용 규칙

1. **지원 파일 형식**
   - 현재 **CSV**(`.csv`)와 **Excel**(`.xlsx`, `.xls`)만 지원합니다.
   - 복잡한 Excel 파일이 읽히지 않으면 표준 CSV로 변환 후 다시 업로드하는 것이 좋습니다.

2. **쿼리 제한**
   - `SELECT`, `SHOW`, `DESCRIBE`, `EXPLAIN`, `PRAGMA` 등 **읽기 전용 쿼리**만 지원합니다.
   - `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP` 등 데이터 변경 작업은 금지됩니다.


## P.S.
위 방법으로 해결되지 않으면 issue에 문제를 상세히 적고 필요한 로그를 함께 제공해 주세요.
