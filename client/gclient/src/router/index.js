import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/views/layout/LayoutIndex.vue'
import Login from '@/views/login/LoginIndex.vue'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/',
    component: Layout,
    children: [
      {
        path: '',
        redirect: '/home'
      },
      {
        path: 'home',
        name: 'home',
        component: () => import('@/views/layout/MyContent.vue')
      },
      {
        path: 'kb/:index_id',
        name: 'knowledge-base',
        component: () => import('@/views/knowledge/KnowledgeBase.vue'),
        props: route => ({
          ...route.params,
          file_id: route.params.file_id,
          file_name: route.params.file_name
        })
      },
      {
        path: 'shared',
        name: 'shared',
        component: () => import('@/views/shared/SharedDocs.vue')
      },
      {
        path: 'recycle',
        name: 'recycle',
        component: () => import('@/views/recycle/RecycleBin.vue')
      }
    ]
  },
  {
    path: '/shared-doc/:id',
    name: 'shared-doc-view',
    component: () => import('@/views/shared/SharedDocView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router