# 내장 모델(Built-in Model) 관리 가이드

## 개요

내장 모델은 시스템 레벨의 모델 설정으로, 모든 테넌트(사용자 그룹)에게 보이지만 민감한 정보는 숨겨지며 편집이나 삭제가 불가능합니다. 보통 모든 사용자가 공통으로 사용할 기본 모델 서비스를 제공할 때 사용합니다.

## 내장 모델의 특징

- **모든 테넌트에게 공개**: 별도 설정 없이 모든 테넌트가 볼 수 있습니다.
- **보안**: API Key나 Base URL 같은 민감한 정보는 숨겨져 있어 상세 내용을 볼 수 없습니다.
- **읽기 전용**: 내장 모델은 수정하거나 삭제할 수 없으며, '기본 모델'로 설정만 가능합니다.
- **통합 관리**: 시스템 관리자가 일괄 관리하므로 설정의 일관성과 안전성을 보장합니다.

## 내장 모델 추가 방법

내장 모델은 데이터베이스에 직접 SQL을 실행하여 추가해야 합니다.

### 1. 모델 데이터 준비

추가할 모델의 설정 정보를 준비합니다:

- 모델 이름 (name)
- 모델 유형 (type): `KnowledgeQA`, `Embedding`, `Rerank`, `VLLM`
- 모델 출처 (source): `local` 또는 `remote`
- 모델 매개변수 (parameters): base_url, api_key, provider 등
- 테넌트 ID (tenant_id): 충돌 방지를 위해 10000 미만의 시스템용 ID 사용 권장

**지원하는 제공자 (provider)**: `generic` (사용자 정의), `openai`, `aliyun`, `zhipu`, `volcengine`, `hunyuan`, `deepseek`, `minimax`, `mimo`, `siliconflow`, `jina`, `openrouter`, `gemini`

### 2. SQL 삽입 실행

아래 SQL 예시를 참고하여 모델을 추가하세요.

```sql
-- 예시: LLM 내장 모델 추가
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-llm-001',                    -- 고정 ID 사용 권장 (builtin- 접두사)
    10000,                                -- 테넌트 ID (첫 번째 테넌트 사용)
    'GPT-4',                              -- 모델 이름
    'KnowledgeQA',                        -- 모델 유형
    'remote',                             -- 모델 출처
    '내장 LLM 모델',                       -- 설명
    '{"base_url": "[https://api.openai.com/v1](https://api.openai.com/v1)", "api_key": "sk-xxx", "provider": "openai"}'::jsonb,  -- 파라미터 (JSON)
    false,                                -- 기본값 여부
    'active',                             -- 상태
    true                                  -- 내장 모델로 표시
) ON CONFLICT (id) DO NOTHING;

-- 예시: Embedding 내장 모델 추가
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-embedding-001',
    10000,
    'text-embedding-ada-002',
    'Embedding',
    'remote',
    '내장 Embedding 모델',
    '{"base_url": "[https://api.openai.com/v1](https://api.openai.com/v1)", "api_key": "sk-xxx", "provider": "openai", "embedding_parameters": {"dimension": 1536, "truncate_prompt_tokens": 0}}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 예시: ReRank 내장 모델 추가
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-rerank-001',
    10000,
    'bge-reranker-base',
    'Rerank',
    'remote',
    '내장 ReRank 모델',
    '{"base_url": "[https://api.jina.ai/v1](https://api.jina.ai/v1)", "api_key": "jina-xxx", "provider": "jina"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 示例：插入一个 VLLM 内置模型 ??
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-vllm-001',
    10000,
    'gpt-4-vision',
    'VLLM',
    'remote',
    '内置 VLLM 模型',
    '{"base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1", "api_key": "sk-xxx", "provider": "aliyun"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;
```

### 3. 결과 확인

다음 쿼리로 내장 모델이 잘 들어갔는지 확인합니다:

```sql
SELECT id, name, type, is_builtin, status
FROM models
WHERE is_builtin = true
ORDER BY type, created_at;
```

## 주의사항

ID 명명 규칙: builtin-{type}-{번호} 형식을 권장합니다. (예: builtin-llm-001)

테넌트 ID: 보통 첫 번째 테넌트 ID인 10000을 사용합니다.

파라미터 형식: parameters 필드는 반드시 유효한 JSON 형식이어야 합니다.

멱등성: ON CONFLICT (id) DO NOTHING 구문을 사용하여 중복 실행 시 에러를 방지합니다.

보안: 프론트엔드에서는 키가 가려지지만 DB에는 원문이 저장되므로 DB 접근 권한 관리에 유의하세요.

## 기존 모델을 내장 모델로 변경하기

이미 등록된 모델을 내장 모델로 바꾸려면:

```sql
UPDATE models
SET is_builtin = true
WHERE id = '모델ID' AND name = '모델이름';
```

## 내장 모델 해제하기

내장 모델 표시를 제거하고 일반 모델로 되돌리려면:

```sql
UPDATE models
SET is_builtin = false
WHERE id = '모델ID';
```

(주의: 해제 후에는 편집 및 삭제가 가능해집니다.)
