#!/usr/bin/env node
import { mkdir, readFile, rm, writeFile } from "node:fs/promises";
import { resolve } from "node:path";

const publicDir = resolve(import.meta.dirname, "../public");
const canonicalOrigin = "https://mirror.atoman.org";
const services = [
  { id: "docker-hub", slug: "docker", name: "Docker Hub", description: "配置 Atoman Mirror 加速 Docker Hub 公开镜像拉取。", notes: "仅代理公开镜像。配置后可用 docker pull mirror.atoman.org/library/alpine:latest 验证；删除 daemon.json 中的 registry-mirrors 并重启 Docker 即可恢复默认源。" },
  { id: "ghcr", slug: "ghcr", name: "GitHub Container Registry", description: "通过 Atoman Mirror 拉取公开 GHCR 容器镜像。", notes: "仅支持公开 GHCR 镜像；私有镜像认证未开放。用公开镜像完成 docker pull 验证，项目不再使用本站前缀即可恢复默认拉取方式。" },
  { id: "pypi", slug: "pypi", name: "PyPI", description: "为 pip 配置 Atoman Mirror 的 PyPI 公开包索引。", notes: "适用于 pip 及兼容 PyPI 的工具。执行 python -m pip install --index-url https://mirror.atoman.org/pypi/simple pip 可验证；删除 ~/.config/pip/pip.conf 中的 index-url 即可恢复。" },
  { id: "npm", slug: "npm", name: "npm", description: "为 npm、pnpm 与 Yarn Classic 配置 Atoman Mirror registry。", notes: "仅代理公开 npm 包。用 npm view npm --registry=https://mirror.atoman.org/npm/ 验证；删除 ~/.npmrc 中的 registry 配置即可恢复默认 registry。" },
  { id: "go", slug: "go", name: "Go Modules", description: "为 Go Modules 配置 Atoman Mirror 模块代理与校验服务。", notes: "保留 sum.golang.org 的完整性校验。执行 go env GOPROXY GOSUMDB 确认配置；使用 go env -u GOPROXY GOSUMDB 可恢复 Go 默认设置。" },
  { id: "ubuntu", slug: "ubuntu", name: "Ubuntu APT", description: "为 Ubuntu 24.04 LTS 配置 Atoman Mirror APT 源。", notes: "本页命令适用于 Ubuntu 24.04 (Noble)。先保留自动创建的 .bak 备份；执行 sudo apt update 验证，恢复时将 ubuntu.sources.bak 还原为 ubuntu.sources。" },
  { id: "debian", slug: "debian", name: "Debian APT", description: "为 Debian 12 Bookworm 配置 Atoman Mirror APT 源。", notes: "本页命令适用于 Debian 12 (Bookworm)。先保留自动创建的 sources.list.bak；执行 sudo apt update 验证，恢复时将该备份还原为 /etc/apt/sources.list。" },
];

function extractArticle(source, id) {
  const start = source.indexOf(`<article id="${id}"`);
  if (start < 0) throw new Error(`Missing ${id} tutorial`);
  const end = source.indexOf("</article>", start);
  return source.slice(start, end + "</article>".length).replace(/ hidden(?=>)/, "");
}

function nav(current) {
  return services.map(({ slug, name }) => `<a href="/${slug}"${slug === current ? ' aria-current="page"' : ""}>${name}</a>`).join("\n            ");
}

function documentFor(service, article) {
  const canonical = `${canonicalOrigin}/${service.slug}`;
  return `<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="${service.description}" />
    <link rel="canonical" href="${canonical}" />
    <meta property="og:type" content="article" />
    <meta property="og:title" content="${service.name} 镜像配置 | Atoman Mirror" />
    <meta property="og:description" content="${service.description}" />
    <meta property="og:url" content="${canonical}" />
    <meta property="og:image" content="https://www.atoman.org/atoman-share.png" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content="${service.name} 镜像配置 | Atoman Mirror" />
    <meta name="twitter:description" content="${service.description}" />
    <meta name="twitter:image" content="https://www.atoman.org/atoman-share.png" />
    <script type="application/ld+json">{"@context":"https://schema.org","@type":"TechArticle","headline":"${service.name} 镜像配置","description":"${service.description}","url":"${canonical}","publisher":{"@type":"Organization","name":"Atoman","url":"https://www.atoman.org/"}}</script>
    <title>${service.name} 镜像配置 | Atoman Mirror</title>
    <link rel="stylesheet" href="/styles.css?v=20260819-2" />
  </head>
  <body data-static-guide="true">
    <a class="skip-link" href="#content">跳到内容</a>
    <header class="topbar"><div class="topbar-inner"><a class="brand-link" href="/" aria-label="Atoman Mirror 首页"><span class="logo-box" aria-hidden="true"><span class="logo-inner"></span></span><span class="logo-copy"><span class="logo-text">ATOMAN</span><span class="logo-meta">MIRROR <span>beta</span></span></span></a><nav class="top-nav" aria-label="Atoman 导航"><a class="top-nav-link is-active" href="/">镜像</a></nav><div class="topbar-actions"><a class="site-link" href="https://www.atoman.org">返回主站</a></div></div></header>
    <div class="app-shell"><aside class="sidebar" aria-label="镜像服务"><div class="sidebar-inner"><p class="sidebar-label">镜像教程</p><nav class="service-nav">${nav(service.slug)}</nav></div></aside><main id="content" class="main-content" tabindex="-1">${article}<section class="config-section" aria-labelledby="guide-notes"><div class="section-heading"><h2 id="guide-notes">验证与恢复</h2><span>02</span></div><p class="hint">${service.notes}</p></section></main></div>
    <footer class="site-footer"><span>ATOMAN MIRROR</span><span>仅代理公开上游资源。</span></footer>
    <script src="/app.js?v=20260819-2" defer></script>
  </body>
</html>`;
}

const source = await readFile(resolve(publicDir, "index.html"), "utf8");
await rm(resolve(publicDir, "guides"), { recursive: true, force: true });
for (const service of services) {
  const directory = resolve(publicDir, service.slug);
  await mkdir(directory, { recursive: true });
  await writeFile(resolve(directory, "index.html"), documentFor(service, extractArticle(source, service.id)));
}
