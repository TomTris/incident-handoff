import { createRouter, createWebHistory } from 'vue-router'
import { isAuthenticated } from '@/api.ts'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '',
      redirect: "/entry",
    },
    {
      path: '/entry',
      name: 'entry',
      component: () => {},
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/incident-list',
      name: 'incident-list',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('@/views/IncidentsView.vue'),
    },
  ],
})


router.beforeEach(async (to) => {
  if (to.name == "entry") {
    const ok = await isAuthenticated()
    return ok == true ? "incident-list" : "login"
  }
})

export default router
