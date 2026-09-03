<script setup lang="ts">
import { computed, ref } from 'vue'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import json from 'highlight.js/lib/languages/json'
import { BookOpen, Check, Container, Copy, Download, Search, Settings2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import PageHero from '@/components/PageHero.vue'
import { copyText } from '@/lib/utils'

hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)

const configCopied = ref(false)
const copyHint = ref('')
const host = computed(() => window.location.host)

const dockerExamples = computed(() => [
  {
    id: 'official',
    label: '官方镜像',
    accelerated: `docker pull ${host.value}/nginx`,
  },
  {
    id: 'user',
    label: '用户镜像',
    accelerated: `docker pull ${host.value}/user/app:tag`,
  },
  {
    id: 'ghcr',
    label: 'GHCR',
    accelerated: `docker pull ${host.value}/ghcr.io/org/app`,
  },
])

const dockerConfig = computed(() => `{
  "registry-mirrors": [
    "https://${host.value}"
  ]
}`)

const restartDockerCommand = 'sudo systemctl restart docker'
const verifyDockerCommand = "docker info | sed -n '/Registry Mirrors/,+3p'"

function highlightCode(code: string, language: 'bash' | 'json') {
  return hljs.highlight(code, { language }).value
}

const highlightedDockerConfig = computed(() => highlightCode(dockerConfig.value, 'json'))
const highlightedRestartDockerCommand = highlightCode(restartDockerCommand, 'bash')
const highlightedVerifyDockerCommand = highlightCode(verifyDockerCommand, 'bash')

async function onCopyConfig() {
  configCopied.value = await copyText(dockerConfig.value)
  copyHint.value = configCopied.value ? 'Docker 配置已复制' : '复制失败，请手动选择配置内容'
  setTimeout(() => {
    copyHint.value = ''
    configCopied.value = false
  }, 2000)
}
</script>

<template>
  <div class="home-page mx-auto max-w-4xl">
    <PageHero
      eyebrow="Docker Registry"
      title="Docker 镜像加速"
      subtitle="为 Docker Hub、GHCR 和其他配置的 Registry 提供稳定的拉取入口"
    >
      <div class="flex flex-wrap justify-start gap-2 pt-1">
        <span class="feature-pill">
          <Container class="size-4" />
          Docker 镜像
        </span>
      </div>
    </PageHero>

    <nav class="quick-actions" aria-label="常用工具">
      <RouterLink to="/search" class="quick-action">
        <span class="quick-action__icon"><Search class="size-4" /></span>
        <span><strong>搜索镜像</strong><small>查看版本、架构与拉取命令</small></span>
      </RouterLink>
      <RouterLink to="/images" class="quick-action">
        <span class="quick-action__icon"><Download class="size-4" /></span>
        <span><strong>离线下载</strong><small>生成可供 docker load 的归档</small></span>
      </RouterLink>
      <a href="https://docs.52013120.xyz/" class="quick-action">
        <span class="quick-action__icon"><BookOpen class="size-4" /></span>
        <span><strong>部署文档</strong><small>反向代理、安全与配置说明</small></span>
      </a>
    </nav>

    <section class="content-section" aria-labelledby="docker-title">
      <div class="section-heading section-heading--stacked">
        <div>
          <p class="section-kicker">DOCKER</p>
          <h2 id="docker-title">配置镜像加速</h2>
        </div>
        <p class="section-description">配置一次后，Docker Hub 镜像会自动使用本站入口。</p>
      </div>

      <div class="config-grid">
        <div class="config-card">
          <div class="config-card-header">
            <span class="config-step">01</span>
            <div>
              <h3>写入 Docker 配置</h3>
              <p>编辑服务器上的 <code>/etc/docker/daemon.json</code></p>
            </div>
            <Button
              variant="ghost"
              size="icon"
              :aria-label="configCopied ? 'Docker 配置已复制' : '复制 Docker 配置'"
              :title="configCopied ? '已复制' : '复制 Docker 配置'"
              @click="onCopyConfig"
            >
              <Check v-if="configCopied" class="size-4" />
              <Copy v-else class="size-4" />
            </Button>
          </div>
          <pre class="code-block code-block--json" data-language="json"><code v-html="highlightedDockerConfig"></code></pre>
          <p class="sr-only" aria-live="polite">{{ copyHint }}</p>
        </div>

        <div class="config-card config-card--commands config-card--combined">
          <div class="config-command-row">
            <div class="config-card-header">
              <span class="config-step">02</span>
              <div>
                <h3>重启 Docker</h3>
                <p>让守护进程读取新配置</p>
              </div>
            </div>
            <pre class="code-block code-block--shell" data-language="shell"><code v-html="highlightedRestartDockerCommand"></code></pre>
          </div>
          <div class="config-command-row">
            <div class="config-card-header">
              <span class="config-step">03</span>
              <div>
                <h3>确认配置生效</h3>
                <p>查看 Docker 当前使用的镜像源</p>
              </div>
            </div>
            <pre class="code-block code-block--shell" data-language="shell"><code v-html="highlightedVerifyDockerCommand"></code></pre>
          </div>
        </div>
      </div>
    </section>

    <section class="content-section content-section--examples" aria-labelledby="examples-title">
      <div class="section-heading section-heading--stacked">
        <div>
          <p class="section-kicker">EXAMPLES</p>
          <h2 id="examples-title">直接使用</h2>
        </div>
        <p class="section-description">也可以在单次拉取时，把本站域名放到镜像名前。</p>
      </div>

      <div class="terminal-block">
        <div class="terminal-header">
          <Settings2 class="size-4" />
          <span>shell</span>
        </div>
        <div class="terminal-body">
          <div v-for="item in dockerExamples" :key="item.id" class="terminal-example">
            <span class="example-tag">{{ item.label }}</span>
            <p class="min-w-0 break-all font-mono leading-relaxed">
              <span class="text-muted-foreground">$ </span>
              <span class="text-primary">{{ item.accelerated }}</span>
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
