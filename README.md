# emby-proxy

Go 实现的 Emby 反向代理，零配置自动改写响应中的 URL，运行于 Docker，前面套 Nginx Proxy Manager 即可使用。

## 原理

通过 URL 路径编码目标地址，自动反代任意 Emby 服务器：

```
https://你的域名/{scheme}/{domain}/{port}/{path}
```

程序自动完成：
- 响应体中所有绝对 URL 改写为代理地址（替代 nginx `sub_filter`）
- 302 重定向 Location 头改写（替代 nginx `proxy_redirect`）
- Referer/Origin 还原为真实上游地址（绕过防盗链）
- 请求/响应中 hop-by-hop 头正确处理，不跨连接转发
- 请求特征头抹除（X-Real-IP / X-Forwarded-* / Via）
- 响应特征头抹除（Server / X-Powered-By）
- Host 伪装 + TLS SNI 自动匹配
- 视频流式转发，支持 Range 断点续传

## 部署

```bash
docker-compose up -d --build
```

服务监听 `8080` 端口，在 Nginx Proxy Manager 中将你的域名反代到 `http://localhost:8080` 即可。

## 使用示例

假设你的代理域名为 `https://proxy.example.com`：

| 场景 | 访问地址 |
|------|---------|
| HTTPS 默认端口 | `https://proxy.example.com/https/emby.example.com/443/` |
| HTTP 自定义端口 | `https://proxy.example.com/http/emby.example.com/2052/` |
| 前后端分离（主站） | `https://proxy.example.com/https/emby.example.com/443/` |

推流地址会被自动改写——客户端无需关心推流域名，所有流量都经过代理。

## 日志

```
[API]    200 GET emby.example.com/emby/Items/123 (45ms)
[STREAM] 206 GET cdn.example.com/videos/xxx.mp4 | 1.2 GB | 3m22s
[PROXY]  200 GET emby.example.com/web/index.html | 156 KB | 89ms
[ERROR]  GET cdn.example.com/videos/fail.mp4 : connection refused
```

- `[STREAM]` 推流/视频，带传输量和耗时
- `[API]` JSON/HTML 等需要改写 body 的请求
- `[PROXY]` 其他二进制透传
- `[ERROR]` 上游连接失败

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `LISTEN_ADDR` | `:8080` | 监听地址 |

## 项目结构

```
├── main.go            # 入口
├── handler.go         # 反代核心：URL 解析、连接池、流式转发
├── rewriter.go        # 响应改写：URL 正则替换、Header 改写
├── Dockerfile         # 多阶段构建
└── docker-compose.yml
```
