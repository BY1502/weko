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
# Embedding模型向量维度
INIT_EMBEDDING_MODEL_DIMENSION=your_embedding_model_dimension
# Embedding模型的ID，通常是一个字符串
INIT_EMBEDDING_MODEL_ID=your_embedding_model_id
```

원격 API를 통해 모델에 접근한다면 `BASE_URL`과 `API_KEY`도 설정해야 합니다:
```bash
# LLM模型的访问地址
INIT_LLM_MODEL_BASE_URL=your_llm_model_base_url
# LLM模型的API密钥，如果需要身份验证，可以设置
INIT_LLM_MODEL_API_KEY=your_llm_model_api_key
# Embedding模型的访问地址
INIT_EMBEDDING_MODEL_BASE_URL=your_embedding_model_base_url
# Embedding模型的API密钥，如果需要身份验证，可以设置
INIT_EMBEDDING_MODEL_API_KEY=your_embedding_model_api_key
```

Rerank 기능이 필요하면 Rerank 모델도 설정합니다:
```bash
# 使用的Rerank模型名称
INIT_RERANK_MODEL_NAME=your_rerank_model_name
# Rerank模型的访问地址
INIT_RERANK_MODEL_BASE_URL=your_rerank_model_base_url
# Rerank模型的API密钥，如果需要身份验证，可以设置
INIT_RERANK_MODEL_API_KEY=your_rerank_model_api_key
```

2. 메인 서비스 로그에 `ERROR`가 있는지 확인하세요.

## 4. 멀티모달 기능을 켜려면?
1. `.env`에서 아래 설정을 맞춥니다:
```bash
# VLM_MODEL_NAME 使用的多模态模型名称
VLM_MODEL_NAME=your_vlm_model_name

# VLM_MODEL_BASE_URL 使用的多模态模型访问地址
VLM_MODEL_BASE_URL=your_vlm_model_base_url

# VLM_MODEL_API_KEY 使用的多模态模型API密钥
VLM_MODEL_API_KEY=your_vlm_model_api_key
```
참고: 멀티모달 모델은 현재 remote API만 지원하므로 `VLM_MODEL_BASE_URL`, `VLM_MODEL_API_KEY`가 필요합니다.

2. 파싱된 파일을 COS에 업로드해야 하므로 `.env`의 COS 설정을 올바르게 입력합니다:
```bash
# 腾讯云COS的访问密钥ID
COS_SECRET_ID=your_cos_secret_id

# 腾讯云COS的密钥
COS_SECRET_KEY=your_cos_secret_key

# 腾讯云COS的区域，例如 ap-guangzhou
COS_REGION=your_cos_region

# 腾讯云COS的桶名称
COS_BUCKET_NAME=your_cos_bucket_name

# 腾讯云COS的应用ID
COS_APP_ID=your_cos_app_id

# 腾讯云COS的路径前缀，用于存储文件
COS_PATH_PREFIX=your_cos_path_prefix
```
重要：务必将COS中文件的权限设置为**公有读**，否则文档解析模块无法正常解析文件

3. 查看文档解析模块日志，查看OCR和Caption是否正确解析和打印


## 5. 如何使用数据分析功能？

在使用数据分析功能前，请确保智能体已配置相关工具：

1. **智能推理**：需在工具配置中勾选以下两个工具：
   - 查看数据元信息
   - 数据分析

2. **快速问答智能体**：无需手动选择工具，即可直接进行简单的数据查询操作。

### 注意事项与使用规范

1. **支持的文件格式**
   - 目前仅支持 **CSV** (`.csv`) 和 **Excel** (`.xlsx`, `.xls`) 格式的文件。
   - 对于复杂的 Excel 文件，如果读取失败，建议将其转换为标准的 CSV 格式后重新上传。

2. **查询限制**
   - 仅支持 **只读查询**，包括 `SELECT`, `SHOW`, `DESCRIBE`, `EXPLAIN`, `PRAGMA` 等语句。
   - 禁止执行任何修改数据的操作，如 `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP` 等。


## P.S.
如果以上方式未解决问题，请在issue中描述您的问题，并提供必要的日志信息辅助我们进行问题排查
