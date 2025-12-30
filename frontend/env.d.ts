/// <reference types="vite/client" />
// 이 파일을 설정해 "@/views/login/index.vue" 모듈이나 타입 선언을 찾지 못하는 오류(ts(2307))를 해결합니다.
// 아래 선언으로 TypeScript에 ".vue"로 끝나는 모든 파일이 Vue 컴포넌트임을 알려 import 시 타입 인식을 돕습니다.
declare module '*.vue' {
    import { Component } from 'vue'; const component: Component; export default component;
}
