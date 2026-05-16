import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { message } from 'ant-design-vue'
import { checkInit } from '@/api/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/views/auth/Init.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/ledger',
    name: 'Ledger',
    component: () => import('@/views/ledger/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/todo',
    name: 'Todo',
    component: () => import('@/views/todo/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/wish',
    name: 'Wish',
    component: () => import('@/views/wish/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/forum',
    name: 'Forum',
    component: () => import('@/views/forum/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/members',
    name: 'Members',
    component: () => import('@/views/member/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/categories',
    name: 'Categories',
    component: () => import('@/views/category/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/Index.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, _from, next) => {
  const token = localStorage.getItem('token')

  // No token: only allow /init, otherwise redirect to /login
  if (!token) {
    if (to.path === '/init') {
      next()
    } else if (to.meta.guest) {
      next()
    } else {
      next('/login')
    }
    return
  }

  // Has token: if accessing guest pages, check init status and redirect
  if (to.meta.guest) {
    try {
      const res: any = await checkInit()
      if (res.data.need_init) {
        if (to.path !== '/init') {
          next('/init')
          return
        }
      } else {
        next('/')
        return
      }
    } catch {
      message.error('网络错误，请刷新重试')
      next(false)
      return
    }
    next()
    return
  }

  // Protected routes: must have token (already checked above)
  next()
})

export default router
