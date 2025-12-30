# uv를 사용하여 WeKnora MCP 서버 실행하기

> Python 기반 MCP 서비스를 실행할 때는 `uv`를 사용하는 것을 권장합니다.

## 1. uv 설치

```bash
# macOS/Linux
curl -LsSf [https://astral.sh/uv/install.sh](https://astral.sh/uv/install.sh) | sh

# 또는 Homebrew 사용 (macOS)
brew install uv

# Windows
powershell -ExecutionPolicy ByPass -c "irm [https://astral.sh/uv/install.ps1](https://astral.sh/uv/install.ps1) | iex"
```

## 2. MCP 클라이언트 설정

### Claude Desktop 설정

Claude Desktop 설정 파일에 다음 내용을 추가하세요:

```json
{
  "mcpServers": {
    "weknora": {
      "args": [
        "--directory",
        "/path/WeKnora/mcp-server",
        "run",
        "run_server.py"
      ],
      "command": "uv",
      "env": {
        "WEKNORA_API_KEY": "your_api_key_here",
        "WEKNORA_BASE_URL": "http://localhost:8080/api/v1"
      }
    }
  }
}
```

### Cursor 설정

Cursor에서 MCP 설정 파일(보통 ~/.cursor/mcp-config.json)을 편집하세요:

```json
{
  "mcpServers": {
    "weknora": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/WeKnora/mcp-server",
        "run",
        "run_server.py"
      ],
      "env": {
        "WEKNORA_API_KEY": "your_api_key_here",
        "WEKNORA_BASE_URL": "http://localhost:8080/api/v1"
      }
    }
  }
}
```

### KiloCode 설정

KiloCode 또는 MCP를 지원하는 다른 에디터의 경우:

```json
{
  "mcpServers": {
    "weknora": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/WeKnora/mcp-server",
        "run",
        "run_server.py"
      ],
      "env": {
        "WEKNORA_API_KEY": "your_api_key_here",
        "WEKNORA_BASE_URL": "http://localhost:8080/api/v1"
      }
    }
  }
}
```

### 기타 MCP 클라이언트

일반적인 MCP 클라이언트 설정:

```json
{
  "mcpServers": {
    "weknora": {
      "command": "uv",
      "args": [
        "--directory",
        "/path/WeKnora/mcp-server",
        "run",
        "run_server.py"
      ],
      "env": {
        "WEKNORA_API_KEY": "your_api_key_here",
        "WEKNORA_BASE_URL": "http://localhost:8080/api/v1"
      }
    }
  }
}
```
