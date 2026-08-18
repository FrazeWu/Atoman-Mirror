# HubProxy

 **Docker 镜像加速代理服务器**

一个轻量级、高性能的 Docker 镜像代理服务，提供 Registry 镜像加速、离线镜像下载与在线搜索功能。


<p align="center">
  <img src="https://count.getloli.com/get/@sky22333.hubproxy?theme=rule34" alt="Visitors">
</p>


## 特性

- 🐳 **Docker 镜像加速** — 兼容 Registry API v2，支持 Docker Hub、GHCR、Quay、GCR、registry.k8s.io；流式传输，Manifest / Token 缓存
- 📦 **离线镜像包** — 无需本地 Docker，在线打包单镜像或批量 tar；流式下载 + 防抖设计
- 🔍 **镜像搜索** — Web 界面与 API 搜索 Docker Hub 镜像、浏览标签
- 🛡️ **智能限流** — 按真实客户端 IP 令牌桶限流（IPv6 按 `/64`）；可配置周期与配额
- 🚫 **仓库访问控制** — IP 黑白名单（限流豁免 / 封禁）+ 镜像黑白名单，支持通配符
- 🌐 **上游 SOCKS5 代理** — 可选配置出站代理，适配特殊网络环境
- 🖥️ **Web 界面** — 内置 Vue SPA，镜像搜索、离线包下载、标签浏览
- ⚡ **轻量高效** — Go 单二进制，支持 `deb` / `rpm` / `apk` 与 Docker 多架构镜像
- 🔧 **统一配置** — `config.toml` + 环境变量覆盖，开箱即用
- 🚀 **统一镜像入口** — 单个程序覆盖 Docker Hub 与多个 OCI Registry，简化部署
- ☁️ **完全自托管** — 不依赖第三方免费 CDN 代理，数据与带宽自主可控

## 快速开始

### Go 原生部署（推荐）

构建前端并将其嵌入 Go 单二进制文件：

```bash
./scripts/build-go.sh
```

启动服务：

```bash
CONFIG_PATH=/etc/hubproxy/config.toml ./build/hubproxy
```

生产环境建议使用仓库中的 systemd 服务文件运行：

```bash
sudo install -Dm755 build/hubproxy /usr/bin/hubproxy
sudo install -Dm644 src/config.toml /etc/hubproxy/config.toml
sudo install -Dm644 packaging/hubproxy.service /etc/systemd/system/hubproxy.service
sudo systemctl daemon-reload
sudo systemctl enable --now hubproxy
```

验证服务：

```bash
curl http://127.0.0.1:5000/ready
```

### 脚本安装

自动识别 `amd64` / `arm64` 与 `apt`、`dnf`、`apk` 等包管理器：

```bash
curl -fsSL https://raw.githubusercontent.com/sky22333/hubproxy/main/install.sh | sh
```

安装后配置文件位于 `/etc/hubproxy/config.toml`，服务自动启动。

### 快速上手

将 `yourdomain.com` 换成你的 `Atoman` 镜像加速地址：
```bash
# Docker Hub 官方镜像
docker pull yourdomain.com/nginx

# Docker Hub 用户镜像
docker pull yourdomain.com/user/app:tag

# GHCR 镜像
docker pull yourdomain.com/ghcr.io/owner/app:tag
```

> **生产环境建议**：绑定自有域名，通过 Caddy / Nginx 反代并开启 HTTPS，不要长期暴露裸 `http://IP:5000`。详见 [文档](https://docs.52013120.xyz/getting-started/quick-start/)。

## 详细文档

部署架构、完整配置、K8s / NAS、传输特性与 FAQ 见官方文档站：

- [**中文文档**](https://docs.52013120.xyz/)
- [**English**](https://docs.52013120.xyz/en/)


## 界面预览

![demo](.github/demo/demo.png)

## 免责声明

- 本程序仅供学习交流使用，请勿用于非法用途
- 使用本程序需遵守当地法律法规
- 作者不对使用者的任何行为承担责任

---

**如果这个项目对你有帮助，请给个 Star ⭐**