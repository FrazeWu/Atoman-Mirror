<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { BookOpen, Container, Download, ExternalLink, Search, Zap } from 'lucide-vue-next'
import ThemeToggle from '@/components/ThemeToggle.vue'

const STORAGE_KEY = 'theme'
const route = useRoute()
const isDark = ref(false)
const serviceStatus = ref<'checking' | 'online' | 'offline'>('checking')

const links = [
  { to: '/', label: 'Docker 加速', icon: Container },
  { to: '/images', label: '离线镜像', icon: Download },
  { to: '/search', label: '镜像搜索', icon: Search },
] as const

const currentPath = computed(() => route.path)
const serviceStatusLabel = computed(() => ({
  checking: '服务检测中',
  online: 'Registry 在线',
  offline: 'Registry 暂不可用',
})[serviceStatus.value])

function applyTheme(dark: boolean) {
  isDark.value = dark
  document.documentElement.classList.toggle('dark', dark)
  localStorage.setItem(STORAGE_KEY, dark ? 'dark' : 'light')
}

function toggleTheme() {
  applyTheme(!isDark.value)
}

async function checkService() {
  try {
    const response = await fetch('/v2/', { cache: 'no-store' })
    serviceStatus.value = response.ok ? 'online' : 'offline'
  } catch {
    serviceStatus.value = 'offline'
  }
}

onMounted(() => {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') {
    applyTheme(saved === 'dark')
  } else {
    applyTheme(window.matchMedia('(prefers-color-scheme: dark)').matches)
  }
  void checkService()
})
</script>

<template>
  <div class="shell-atmosphere flex min-h-screen flex-col text-foreground">
    <header class="sticky top-0 z-50 border-b border-border bg-background">
      <div class="mx-auto flex h-16 max-w-6xl items-center justify-between gap-3 px-4 sm:px-8">
        <RouterLink
          to="/"
          class="flex min-w-0 items-center gap-3 font-display text-lg font-semibold transition-opacity hover:opacity-80"
        >
          <span class="brand-mark flex size-9 shrink-0 items-center justify-center">
            <Zap class="size-[18px]" />
          </span>
          <span class="truncate">Atoman Mirror</span>
        </RouterLink>

        <div class="flex shrink-0 items-center gap-1.5">
          <span
            class="service-status hidden items-center gap-2 text-xs text-muted-foreground sm:inline-flex"
            :data-state="serviceStatus"
            aria-live="polite"
          >
            <span class="service-status__dot" aria-hidden="true" />
            {{ serviceStatusLabel }}
          </span>
          <a
            href="https://www.atoman.org/"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden h-10 items-center gap-1.5 px-3 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground lg:inline-flex"
          >
            <ExternalLink class="size-4" />
            Atoman 主站
          </a>
          <ThemeToggle :is-dark="isDark" @toggle="toggleTheme" />
        </div>
      </div>

      <div class="border-t border-border">
        <nav aria-label="服务导航" class="mx-auto flex max-w-6xl gap-1 overflow-x-auto px-4 py-1.5 sm:px-8">
          <RouterLink
            v-for="link in links"
            :key="link.to"
            :to="link.to"
            class="inline-flex h-10 shrink-0 items-center gap-2 px-3 text-sm transition-colors"
            :class="currentPath === link.to ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:bg-accent hover:text-foreground'"
          >
            <component :is="link.icon" class="size-4" />
            {{ link.label }}
          </RouterLink>
          <a
            href="https://docs.52013120.xyz/"
            class="inline-flex h-10 shrink-0 items-center gap-2 px-3 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          >
            <BookOpen class="size-4" />
            使用文档
          </a>
        </nav>
      </div>
    </header>

    <main class="mx-auto w-full max-w-6xl flex-1 px-4 py-8 text-base sm:px-8 sm:py-12">
      <slot />
    </main>

    <footer class="border-t border-border px-4 py-6 text-sm text-muted-foreground sm:px-8">
      <div class="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-3">
        <span>Atoman Mirror · Registry Gateway</span>
        <span class="font-mono text-xs">{{ serviceStatusLabel }}</span>
      </div>
    </footer>
  </div>
</template>
