# 빌드 단계
FROM node:24-alpine as build-stage

WORKDIR /app

# 타입 체크 오류를 무시하도록 환경 변수 설정
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV VITE_IS_DOCKER=true

# 파일 크기 제한(MB), 빌드 인자로 전달
ARG MAX_FILE_SIZE_MB=50
ENV VITE_MAX_FILE_SIZE_MB=${MAX_FILE_SIZE_MB}

# 의존성 파일 복사
COPY package*.json ./
COPY packages/xlsx-0.20.2.tgz ./packages/xlsx-0.20.2.tgz

# 의존성 설치
RUN corepack enable
RUN pnpm install

# 프로젝트 파일 복사
COPY . .

# 애플리케이션 빌드
RUN pnpm run build

# 프로덕션 단계
FROM nginx:stable-alpine as production-stage

# 빌드 산출물을 nginx 서비스 디렉터리로 복사
COPY --from=build-stage /app/dist /usr/share/nginx/html

# nginx 설정 템플릿 복사
COPY nginx.conf /etc/nginx/templates/default.conf.template

# 기본 환경 변수 설정(MB)
ENV MAX_FILE_SIZE_MB=50

# 포트 노출
EXPOSE 80

# 시작 시 MAX_FILE_SIZE_MB를 단위를 포함한 MAX_FILE_SIZE로 변환해 nginx 설정에 주입
CMD ["/bin/sh", "-c", "export MAX_FILE_SIZE=${MAX_FILE_SIZE_MB}M && envsubst '${MAX_FILE_SIZE}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"] 
