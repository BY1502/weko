````markdown
# WeKnora 개발 가이드

## 빠른 개발 모드 (권장)

`app`(백엔드)이나 `frontend`(프론트엔드) 코드를 자주 수정해야 한다면, **매번 Docker 이미지를 다시 빌드할 필요 없이** 로컬 개발 모드를 사용하세요.

### 방법 1: Make 명령어 사용 (권장)

#### 1. 인프라 서비스 시작

```bash
make dev-start
```
````

다음 서비스들의 Docker 컨테이너를 시작합니다:

PostgreSQL (데이터베이스)

Redis (캐시)

MinIO (객체 스토리지)

Neo4j (그래프 데이터베이스)

DocReader (문서 파싱 서비스)

Jaeger (분산 트레이싱)

2. 백엔드 앱 시작 (새 터미널)

```bash
make dev-app
```

로컬에서 Go 애플리케이션을 직접 실행합니다. 코드를 수정한 뒤 Ctrl+C로 멈추고 다시 실행하면 반영됩니다.

#### 3. 프론트엔드 시작 (새 터미널)

```bash
make dev-frontend
```

Vite 개발 서버를 시작합니다. 핫 리로드(Hot Reload)를 지원하여 코드 수정 시 자동으로 브라우저가 새로고침 됩니다.

#### 4. 서비스 상태 확인

```bash
make dev-status
```

#### 5. 모든 서비스 중지

```bash
make dev-stop
```

### 방법 2: 스크립트 직접 사용

Make 명령어가 안 된다면 스크립트를 직접 실행하세요

```bash
# 인프라 시작
./scripts/dev.sh start

# 백엔드 시작 (새 터미널)
./scripts/dev.sh app

# 프론트엔드 시작 (새 터미널)
./scripts/dev.sh frontend

# 로그 확인
./scripts/dev.sh logs

# 모든 서비스 중지
./scripts/dev.sh stop
```

## 접속 주소

### 개발 환경

- **프론트엔드 개발 서버**: http://localhost:5173
- **백엔드 API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **MinIO**: 콘솔: http://localhost:9001
- **Neo4j**: 브라우저: http://localhost:7474
- **Jaeger UI**: http://localhost:16686

## 개발 워크플로우 비교

### ❌ 기존 방식 (느림)

```bash
# 코드 수정할 때마다:
sh scripts/build_images.sh -p      # 이미지 재빌드 (오래 걸림)
sh scripts/start_all.sh --no-pull  # 컨테이너 재시작
```

**소요 시간**: 수정 시마다 2~5분

### ✅ 새로운 방식 (빠름)

```bash
# 최초 1회 실행:
make dev-start

# 각각 다른 터미널에서 실행:
make dev-app       # Go 코드 수정 후 재시작 (수 초 내 완료)
make dev-frontend  # 프론트엔드 수정 시 즉시 반영 (재시작 불필요)
```

**소요 시간**：

- 최초 기동: 1~2분
- 백엔드 수정: 5~10초 (Go 앱 재시작)
- 프론트엔드 수정: 실시간 (Hot Reload)

## Air를 이용한 백엔드 핫 리로드 (선택 사항)

백엔드 코드 수정 시 수동 재시작 없이 자동으로 반영되게 하려면 air를 설치하세요:

### 1. Air 설치

```bash
go install github.com/cosmtrek/air@latest
```

### 2. 설정 파일 생성

프로젝트 루트 경로에 `.air.toml` 파일이 이미 있습니다. (번역된 파일 참고)

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "frontend", "migrations"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### 3. Air로 시작

```bash
# 프로젝트 루트에서
air
```

이제 Go 코드를 저장하면 자동으로 다시 컴파일하고 재시작합니다!

## 기타 개발 팁

### 프론트만 수정

```bash
cd frontend
npm run dev
```

프론트엔드는 자동으로 http://localhost:8080 의 백엔드 API에 연결됩니다.

### 백엔드만 수정

```bash
# 인프라 시작
make dev-start

# 백엔드 실행
make dev-app
```

### 디버깅 모드

#### 백엔드 디버깅

VS Code나 GoLand의 디버깅 기능을 사용하여 로컬에서 실행 중인 Go 앱에 연결할 수 있습니다.

VS Code 설정 예시 (`.vscode/launch.json`):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Server",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {
        "DB_HOST": "localhost",
        "DOCREADER_ADDR": "localhost:50051",
        "MINIO_ENDPOINT": "localhost:9000",
        "REDIS_ADDR": "localhost:6379",
        "OTEL_EXPORTER_OTLP_ENDPOINT": "localhost:4317",
        "NEO4J_URI": "bolt://localhost:7687"
      },
      "args": []
    }
  ]
}
```

#### 프론트엔드 디버깅

브라우저 개발자 도구(F12)를 사용하세요. Vite가 소스 맵(Source Map)을 제공하므로 원본 코드를 보며 디버깅할 수 있습니다.

## 운영 환경 배포

개발이 끝나고 실제로 배포할 때만 이미지를 빌드하면 됩니다:

```bash
# 모든 이미지 빌드
sh scripts/build_images.sh

# 또는 특정 이미지만 빌드
sh scripts/build_images.sh -p  # 백엔드만 빌드
sh scripts/build_images.sh -f  # 프론트엔드만 빌드

# 운영 환경 시작
sh scripts/start_all.sh
```

## 문제 해결

### Q: dev-app 시작 시 DB 연결 오류가 발생해요.

A: make dev-start를 먼저 실행했는지 확인하고, 모든 컨테이너가 완전히 뜰 때까지(약 30초) 기다려 주세요.

### Q: 프론트엔드에서 API 호출 시 CORS 에러가 떠요.

A: 프론트엔드의 vite.config.ts 파일에서 프록시 설정이 올바르게 되어 있는지 확인하세요.

### Q: DocReader 서비스를 수정해야 하면 어떻게 하나요?

A: DocReader는 여전히 Docker 이미지를 사용합니다. 코드를 수정했다면 다시 빌드해야 합니다:

```bash
sh scripts/build_images.sh -d
make dev-restart
```

## 요약

- **일상 개발**: `make dev-*` 명령어로 빠르게 반복 개발
- **통합 테스트**: `sh scripts/start_all.sh --no-pull`로 전체 환경 테스트
- **배포**: `sh scripts/build_images.sh` + `sh scripts/start_all.sh`
