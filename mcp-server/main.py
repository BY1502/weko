#!/usr/bin/env python3
"""
WeKnora MCP Server 메인 엔트리 포인트

WeKnora MCP 서버를 실행하기 위한 통합 진입점을 제공합니다.
다음과 같이 실행할 수 있습니다:
1. python main.py
2. python -m weknora_mcp_server
3. weknora-mcp-server (설치 후)
"""

import argparse
import asyncio
import os
import sys
from pathlib import Path


def setup_environment():
    """환경과 경로를 설정"""
    # 현재 디렉터리가 Python 경로에 포함되어 있는지 확인
    current_dir = Path(__file__).parent.absolute()
    if str(current_dir) not in sys.path:
        sys.path.insert(0, str(current_dir))


def check_dependencies():
    """필요한 의존성이 설치되었는지 확인"""
    try:
        import mcp
        import requests

        return True
    except ImportError as e:
        print(f"필요한 라이브러리가 없습니다: {e}")
        print("다음 명령을 실행하세요: pip install -r requirements.txt")
        return False


def check_environment_variables():
    """환경 변수 설정 확인"""
    base_url = os.getenv("WEKNORA_BASE_URL")
    api_key = os.getenv("WEKNORA_API_KEY")

    print("=== WeKnora MCP Server 환경 점검 ===")
    print(f"Base URL: {base_url or 'http://localhost:8080/api/v1 (기본값)'}")
    print(f"API Key: {'설정됨' if api_key else '미설정 (경고)'}")

    if not base_url:
        print("안내: WEKNORA_BASE_URL 환경 변수를 설정할 수 있습니다")

    if not api_key:
        print("경고: WEKNORA_API_KEY 환경 변수를 설정하는 것을 권장합니다")

    print("=" * 40)
    return True


def parse_arguments():
    """명령행 인수 파싱"""
    parser = argparse.ArgumentParser(
        description="WeKnora MCP Server - Model Context Protocol server for WeKnora API",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
사용 예시:
  python main.py                    # 기본 설정으로 시작
  python main.py --check-only       # 환경만 점검하고 서버는 시작하지 않음
  python main.py --verbose          # 상세 로그 활성화

환경 변수:
  WEKNORA_BASE_URL    WeKnora API 기본 URL (기본값: http://localhost:8080/api/v1)
  WEKNORA_API_KEY     WeKnora API 키
        """,
    )

    parser.add_argument(
        "--check-only", action="store_true", help="환경만 점검하고 서버는 시작하지 않음"
    )

    parser.add_argument("--verbose", "-v", action="store_true", help="상세 로그 출력 활성화")

    parser.add_argument(
        "--version", action="version", version="WeKnora MCP Server 1.0.0"
    )

    return parser.parse_args()


async def main():
    """메인 함수"""
    args = parse_arguments()

    # 환경 설정
    setup_environment()

    # 의존성 확인
    if not check_dependencies():
        sys.exit(1)

    # 환경 변수 확인
    check_environment_variables()

    # 환경만 점검하도록 요청한 경우 종료
    if args.check_only:
        print("환경 점검이 완료되었습니다.")
        return

    # 로그 레벨 설정
    if args.verbose:
        import logging

        logging.basicConfig(level=logging.DEBUG)
        print("상세 로그 모드를 활성화했습니다")

    try:
        print("WeKnora MCP Server를 시작하는 중입니다...")

        # 서버를 가져와 실행
        from weknora_mcp_server import run

        await run()

    except ImportError as e:
        print(f"모듈 로드 오류: {e}")
        print("모든 파일이 올바른 위치에 있는지 확인하세요")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\n서버가 중지되었습니다")
    except Exception as e:
        print(f"서버 실행 오류: {e}")
        if args.verbose:
            import traceback

            traceback.print_exc()
        sys.exit(1)


def sync_main():
    """entry_points에서 사용할 동기 버전 메인 함수"""
    asyncio.run(main())


if __name__ == "__main__":
    asyncio.run(main())
