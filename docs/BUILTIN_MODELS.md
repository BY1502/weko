# 내장 모델 관리 가이드

## 개요

내장 모델은 시스템 수준의 모델 설정으로 모든 테넌트에서 볼 수 있지만, 민감 정보는 숨겨지며 수정·삭제할 수 없습니다. 시스템 기본 모델 구성을 제공해 모든 테넌트가 동일한 모델 서비스를 사용할 수 있도록 합니다.

## 내장 모델 특징

- **모든 테넌트에 노출**: 별도 설정 없이 공용으로 노출
- **보안 보호**: API Key, Base URL 등 민감 정보는 숨김 처리
- **읽기 전용**: 편집/삭제 불가, 기본 모델로만 설정 가능
- **통합 관리**: 시스템 관리자가 일관된 설정을 유지

## 내장 모델 추가 방법

내장 모델은 DB에 직접 삽입해야 합니다. 아래 절차를 따르세요.

### 1. 모델 데이터 준비

내장 모델로 지정할 설정 정보를 준비합니다.
- 모델 이름(name)
- 모델 타입(type): `KnowledgeQA`, `Embedding`, `Rerank`, `VLLM`
- 모델 소스(source): `local` 또는 `remote`
- 모델 파라미터(parameters): base_url, api_key, provider 등
- 테넌트 ID(tenant_id): 충돌을 피하기 위해 10000 미만 권장

**지원 provider**: `generic`(custom), `openai`, `aliyun`, `zhipu`, `volcengine`, `hunyuan`, `deepseek`, `minimax`, `mimo`, `siliconflow`, `jina`, `openrouter`, `gemini`

### 2. SQL 삽입 실행

아래 SQL로 내장 모델을 추가합니다.

```sql
-- 예시: LLM 내장 모델 삽입
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
    'builtin-llm-001',                    -- 고정 ID, builtin- 접두사 권장
    10000,                                -- 테넌트 ID(첫 테넌트 사용)
    'GPT-4',                              -- 모델 이름
    'KnowledgeQA',                        -- 모델 타입
    'remote',                             -- 모델 소스
    '내장 LLM 모델',                       -- 설명
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai"}'::jsonb,  -- 파라미터(JSON)
    false,                                -- 기본 여부
    'active',                             -- 상태
    true                                  -- 내장 모델 표시
) ON CONFLICT (id) DO NOTHING;

-- 예시: Embedding 내장 모델 삽입
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
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai", "embedding_parameters": {"dimension": 1536, "truncate_prompt_tokens": 0}}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 예시: ReRank 내장 모델 삽입
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
    '{"base_url": "https://api.jina.ai/v1", "api_key": "jina-xxx", "provider": "jina"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 예시: VLLM 내장 모델 삽입
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

### 3. 验证插入结果

执行以下 SQL 查询验证内置模型是否成功插入：

```sql
SELECT id, name, type, is_builtin, status 
FROM models 
WHERE is_builtin = true
ORDER BY type, created_at;
```

## 注意事项

1. **ID 命名规范**：建议使用 `builtin-{type}-{序号}` 的格式，例如 `builtin-llm-001`、`builtin-embedding-001`
2. **租户ID**：内置模型可以属于任意租户，但建议使用第一个租户ID（通常是 10000）
3. **参数格式**：`parameters` 字段必须是有效的 JSON 格式
4. **幂等性**：使用 `ON CONFLICT (id) DO NOTHING` 确保重复执行不会报错
5. **安全性**：内置模型的 API Key 和 Base URL 在前端会被自动隐藏，但数据库中的原始数据仍然存在，请妥善保管数据库访问权限

## 将现有模型设置为内置模型

如果你已经有一个模型，想将其设置为内置模型，可以使用 UPDATE 语句：

```sql
UPDATE models 
SET is_builtin = true 
WHERE id = '模型ID' AND name = '模型名称';
```

## 移除内置模型

如果需要移除内置模型标记（恢复为普通模型），执行：

```sql
UPDATE models 
SET is_builtin = false 
WHERE id = '模型ID';
```

注意：移除内置模型标记后，该模型将恢复为普通模型，可以被编辑和删除。
