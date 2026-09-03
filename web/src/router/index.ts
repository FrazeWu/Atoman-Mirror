import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import ImagesPage from '@/pages/ImagesPage.vue'
import SearchPage from '@/pages/SearchPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: HomePage,
      meta: { title: 'Docker 镜像加速' },
    },
    {
      path: '/images',
      component: ImagesPage,
      meta: { title: '离线镜像下载' },
    },
    {
      path: '/search',
      component: SearchPage,
      meta: { title: '镜像搜索' },
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.path !== from.path) return { top: 0, left: 0 }
    return false
  },
})

router.afterEach((to) => {
  const title = (to.meta.title as string) || 'Atoman Mirror'
  document.title = `${title} · Atoman Mirror`
})

export default router
