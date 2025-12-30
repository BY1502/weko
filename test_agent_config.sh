#!/bin/bash

# Agent 설정 기능 테스트 스크립트

set -e

echo "========================================="
echo "Agent 설정 기능 테스트"
echo "========================================="
echo ""

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 설정
API_BASE_URL="http://localhost:8080"
KB_ID="kb-00000001"  # 자신의 지식베이스 ID로 변경
TENANT_ID="1"

echo "설정 정보:"
echo "  API 주소: ${API_BASE_URL}"
echo "  지식베이스 ID: ${KB_ID}"
echo "  테넌트 ID: ${TENANT_ID}"
echo ""

# 테스트 1: 현재 설정 조회
echo -e "${YELLOW}테스트 1: 현재 설정 조회${NC}"
echo "GET ${API_BASE_URL}/api/v1/initialization/config/${KB_ID}"
RESPONSE=$(curl -s -X GET "${API_BASE_URL}/api/v1/initialization/config/${KB_ID}")
echo "응답:"
echo "$RESPONSE" | jq '.data.agent' || echo "$RESPONSE"
echo ""

# 테스트 2: Agent 설정 저장
echo -e "${YELLOW}테스트 2: Agent 설정 저장${NC}"
echo "POST ${API_BASE_URL}/api/v1/initialization/initialize/${KB_ID}"

# 테스트 데이터 준비(전체 설정 포함)
TEST_DATA='{
  "llm": {
    "source": "local",
    "modelName": "qwen3:0.6b",
    "baseUrl": "",
    "apiKey": ""
  },
  "embedding": {
    "source": "local",
    "modelName": "nomic-embed-text:latest",
    "baseUrl": "",
    "apiKey": "",
    "dimension": 768
  },
  "rerank": {
    "enabled": false
  },
  "multimodal": {
    "enabled": false
  },
  "documentSplitting": {
    "chunkSize": 512,
    "chunkOverlap": 100,
    "separators": ["\n\n", "\n", "。", "！", "？", ";", "；"]
  },
  "nodeExtract": {
    "enabled": false
  },
  "agent": {
    "enabled": true,
    "maxIterations": 8,
    "temperature": 0.8,
    "allowedTools": ["knowledge_search", "multi_kb_search", "list_knowledge_bases"]
  }
}'

RESPONSE=$(curl -s -X POST "${API_BASE_URL}/api/v1/initialization/initialize/${KB_ID}" \
  -H "Content-Type: application/json" \
  -d "$TEST_DATA")

if echo "$RESPONSE" | grep -q '"success":true'; then
  echo -e "${GREEN}✓ Agent 설정 저장 성공${NC}"
  echo "$RESPONSE" | jq '.' || echo "$RESPONSE"
else
  echo -e "${RED}✗ Agent 설정 저장 실패${NC}"
  echo "$RESPONSE"
fi
echo ""

# 잠시 대기(저장 완료 보장)
sleep 1

# 테스트 3: 설정 저장 여부 검증
echo -e "${YELLOW}테스트 3: 설정 저장 확인${NC}"
echo "GET ${API_BASE_URL}/api/v1/initialization/config/${KB_ID}"
RESPONSE=$(curl -s -X GET "${API_BASE_URL}/api/v1/initialization/config/${KB_ID}")
AGENT_CONFIG=$(echo "$RESPONSE" | jq '.data.agent')

echo "Agent 설정:"
echo "$AGENT_CONFIG" | jq '.'

# 설정 검증
ENABLED=$(echo "$AGENT_CONFIG" | jq -r '.enabled')
MAX_ITER=$(echo "$AGENT_CONFIG" | jq -r '.maxIterations')
TEMP=$(echo "$AGENT_CONFIG" | jq -r '.temperature')

if [ "$ENABLED" == "true" ] && [ "$MAX_ITER" == "8" ] && [ "$TEMP" == "0.8" ]; then
  echo -e "${GREEN}✓ 설정 검증 성공 - 모든 값이 일치${NC}"
else
  echo -e "${RED}✗ 설정 검증 실패${NC}"
  echo "  enabled: $ENABLED (기대값: true)"
  echo "  maxIterations: $MAX_ITER (기대값: 8)"
  echo "  temperature: $TEMP (기대값: 0.8)"
fi
echo ""

# 테스트 4: Tenant API로 설정 조회
echo -e "${YELLOW}테스트 4: Tenant API로 설정 조회${NC}"
echo "GET ${API_BASE_URL}/api/v1/tenants/${TENANT_ID}/agent-config"
RESPONSE=$(curl -s -X GET "${API_BASE_URL}/api/v1/tenants/${TENANT_ID}/agent-config")
echo "응답:"
echo "$RESPONSE" | jq '.' || echo "$RESPONSE"
echo ""

# 테스트 5: DB 검증(접근 가능 시)
echo -e "${YELLOW}테스트 5: DB 검증${NC}"
echo "안내: 아래 SQL을 수동 실행해 데이터를 확인하세요:"
echo ""
echo "MySQL:"
echo "  mysql -u root -p weknora -e \"SELECT id, agent_config FROM tenants WHERE id = ${TENANT_ID};\""
echo ""
echo "PostgreSQL:"
echo "  psql -U postgres -d weknora -c \"SELECT id, agent_config FROM tenants WHERE id = ${TENANT_ID};\""
echo ""

echo "========================================="
echo "테스트 완료!"
echo "========================================="
echo ""
echo "모든 테스트를 통과했다면 Agent 설정 기능이 정상입니다."
echo "실패 시 다음을 확인하세요:"
echo "  1. 백엔드 서비스가 실행 중인지"
echo "  2. DB 마이그레이션이 완료되었는지"
echo "  3. 지식베이스 ID가 올바른지"
echo ""
