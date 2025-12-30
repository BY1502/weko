import { createRouter, createWebHistory } from 'vue-router'
import { listKnowledgeBases } from '@/api/knowledge-base'
import { useAuthStore } from '@/stores/auth'
import { validateToken } from '@/api/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/platform/knowledge-bases",
    },
    {
      path: "/login",
      name: "login",
      component: () => import("../views/auth/Login.vue"),
      meta: { requiresAuth: false, requiresInit: false }
    },
    {
      path: "/knowledgeBase",
      name: "home",
      component: () => import("../views/knowledge/KnowledgeBase.vue"),
      meta: { requiresInit: true, requiresAuth: true }
    },
    {
      path: "/platform",
      name: "Platform",
      redirect: "/platform/knowledge-bases",
      component: () => import("../views/platform/index.vue"),
      meta: { requiresInit: true, requiresAuth: true },
      children: [
        {
          path: "tenant",
          redirect: "/platform/settings"
        },
        {
          path: "settings",
          name: "settings",
          component: () => import("../views/settings/Settings.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "knowledge-bases",
          name: "knowledgeBaseList",
          component: () => import("../views/knowledge/KnowledgeBaseList.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "knowledge-bases/:kbId",
          name: "knowledgeBaseDetail",
          component: () => import("../views/knowledge/KnowledgeBase.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "agents",
          name: "agentList",
          component: () => import("../views/agent/AgentList.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "creatChat",
          name: "globalCreatChat",
          component: () => import("../views/creatChat/creatChat.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "knowledge-bases/:kbId/creatChat",
          name: "kbCreatChat",
          component: () => import("../views/creatChat/creatChat.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
        {
          path: "chat/:chatid",
          name: "chat",
          component: () => import("../views/chat/index.vue"),
          meta: { requiresInit: true, requiresAuth: true }
        },
      ],
    },
  ],
});

// 라우터 가드: 인증 상태와 시스템 초기화 상태 확인
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 로그인/초기화 페이지 접근 시 바로 통과
  if (to.meta.requiresAuth === false || to.meta.requiresInit === false) {
    // 로그인한 사용자가 로그인 페이지로 접근하면 지식베이스 목록으로 리다이렉트
    if (to.path === '/login' && authStore.isLoggedIn) {
      next('/platform/knowledge-bases')
      return
    }
    next()
    return
  }

  // 사용자 인증 상태 확인
  if (to.meta.requiresAuth !== false) {
    if (!authStore.isLoggedIn) {
      // 미로그인 시 로그인 페이지로 이동
      next('/login')
      return
    }

    // 토큰 유효성 검증
    // try {
    //   const { valid } = await validateToken()
    //   if (!valid) {
    //     // 토큰이 무효면 인증 정보를 삭제하고 로그인 페이지로 이동
    //     authStore.logout()
    //     next('/login')
    //     return
    //   }
    // } catch (error) {
    //   console.error('토큰 검증 실패:', error)
    //   authStore.logout()
    //   next('/login')
    //   return
    // }
  }

  next()
});

export default router
