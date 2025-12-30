#!/usr/bin/env python3
"""
WeKnora MCP Server 시작 스크립트
"""

import asyncio
import os
import sys


def check_environment():
    """환경 설정 확인"""
    base_url = os.getenv("WEKNORA_BASE_URL")
    api_key = os.getenv("WEKNORA_API_KEY")

    if not base_url:
        print(
            "경고: WEKNORA_BASE_URL 환경 변수가 설정되지 않아 기본값 http://localhost:8080/api/v1 을 사용합니다"
        )

    if not api_key:
        print("경고: WEKNORA_API_KEY 환경 변수가 설정되지 않았습니다")

    print(f"WeKnora Base URL: {base_url or 'http://localhost:8080/api/v1'}")
    print(f"API Key: {'설정됨' if api_key else '미설정'}")


def main():
    """메인 함수"""
    print("WeKnora MCP Server를 시작합니다...")
    check_environment()

    try:
        from weknora_mcp_server import run

        asyncio.run(run())
    except ImportError as e:
        print(f"모듈 로드 오류: {e}")
        print("필요한 의존성을 설치하세요: pip install -r requirements.txt")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\n서버가 중지되었습니다")
    except Exception as e:
        print(f"서버 실행 오류: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
