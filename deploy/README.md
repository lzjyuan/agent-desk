# agent-service 非 Docker 部署

发布程序名和 systemd 服务名为 `agent-service`。Go 模块名、仓库目录不改。只支持 Linux/amd64、非 LanceDB 构建；前端和工作流编辑器在本地构建后嵌入 Go 程序。

## 安装 Qdrant

默认向量后端为 Qdrant，和 `agent-service` 安装在同一台服务器，不需要单独域名。当前服务器是 Linux x86_64、glibc 2.38；Qdrant v1.20.0 的 GNU 包要求更高的 glibc，改用已验证可运行的 [官方 MUSL 发布包](https://github.com/qdrant/qdrant/releases/tag/v1.20.0)。在服务器上执行：

```bash
tmp=$(mktemp -d)
curl -fL --retry 3 \
  https://github.com/qdrant/qdrant/releases/download/v1.20.0/qdrant-x86_64-unknown-linux-musl.tar.gz \
  -o "$tmp/qdrant.tar.gz"
tar -xzf "$tmp/qdrant.tar.gz" -C "$tmp"
"$tmp/qdrant" --version
sudo install -m 755 "$tmp/qdrant" /usr/local/bin/qdrant
id qdrant >/dev/null 2>&1 || sudo useradd --system --user-group --home-dir /var/lib/qdrant --shell /sbin/nologin qdrant
sudo install -d -m 750 -o qdrant -g qdrant /var/lib/qdrant /var/lib/qdrant/storage /var/lib/qdrant/snapshots
sudo install -d -m 755 /etc/qdrant
```

创建 `/etc/qdrant/config.yaml`，固定数据目录和监听地址；`6333` 用于本机健康检查，`agent-service` 通过 gRPC `6334` 连接。不要把这两个端口开放到公网。[Qdrant 配置](https://qdrant.tech/documentation/operations/configuration/)、[网络安全说明](https://qdrant.tech/documentation/security/)：

```bash
sudo tee /etc/qdrant/config.yaml >/dev/null <<'YAML'
storage:
  storage_path: /var/lib/qdrant/storage
  snapshots_path: /var/lib/qdrant/snapshots
service:
  host: 127.0.0.1
  http_port: 6333
  grpc_port: 6334
YAML
sudo chown root:qdrant /etc/qdrant/config.yaml
sudo chmod 640 /etc/qdrant/config.yaml
```

创建 systemd 服务并启动：

```bash
sudo tee /etc/systemd/system/qdrant.service >/dev/null <<'UNIT'
[Unit]
Description=Qdrant Vector Database
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=qdrant
Group=qdrant
WorkingDirectory=/var/lib/qdrant
ExecStart=/usr/local/bin/qdrant --config-path /etc/qdrant/config.yaml
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target
UNIT
sudo systemctl daemon-reload
sudo systemctl enable --now qdrant
curl -fsS http://127.0.0.1:6333/healthz
ss -lnt '( sport = :6334 )' | grep LISTEN
```

健康检查失败时查看 `sudo journalctl -u qdrant -n 50 --no-pager`。应用的 `config/config.yaml` 应保持 `vectorDB.type: qdrant`，并设置 `vectorDB.qdrant.host: 127.0.0.1`、`grpcPort: 6334`、`useTls: false`；同机连接不需要 Qdrant 域名。将 `/var/lib/qdrant` 纳入备份，不能只备份应用的 `shared/data`。

## 首次准备服务器

先按上节安装并启动 Qdrant（或配置其他受支持的向量后端）。配置文件必须放在 `deploy/target.local.json` 的 `remoteRoot/shared/config/config.yaml`，不能放在其他目录。以下以 `remoteRoot=/www/wwwroot/agent-service`、`serviceUser=www` 为例，在服务器上准备目录：

```bash
root=/www/wwwroot/agent-service
service_user=www
service_group=$(id -gn "$service_user")
sudo install -d -m 755 "$root" "$root/releases"
sudo install -d -m 750 -o "$service_user" -g "$service_group" "$root/shared" "$root/shared/config" "$root/shared/data"
```

在本地确认 `config/config.yaml` 的数据库、Qdrant 地址和 `customerSession.secret` 适用于服务器，然后从项目根目录用 PowerShell 上传：

```powershell
$t = Get-Content .\deploy\target.local.json -Raw | ConvertFrom-Json
scp -i $t.identityFile -P $t.sshPort .\config\config.yaml "$($t.user)@$($t.host):$($t.remoteRoot)/shared/config/config.yaml"
```

上传后在服务器执行 `sudo chown www:www /www/wwwroot/agent-service/shared/config/config.yaml` 和 `sudo chmod 600 /www/wwwroot/agent-service/shared/config/config.yaml`；使用其他 `serviceUser` 或 `remoteRoot` 时相应替换。默认配置使用 SQLite，其 DSN 和本地文件存储保持在相对路径 `data/` 下；程序工作目录为 `remoteRoot/shared`。健康检查固定访问 `127.0.0.1:8083/api/config`，因此服务端口保持 `8083`。`shared/data` 和 Qdrant 数据分别备份。

## 从 Windows PowerShell 发布

本机需要 Go、Node.js 和 OpenSSH 客户端。将 `deploy/target.example.json` 复制为 `deploy/target.local.json`，填写 SSH 地址、账号、端口、服务运行用户和服务器目录。此文件已被 Git 忽略。SSH 账号需为 root，或能免密码执行 `sudo`。

在项目根目录执行：

```powershell
.\deploy\deploy_backend.ps1
```

脚本使用 `npx pnpm@10.30.2` 构建前端，再交叉编译 `linux/amd64` 程序；上传发布包并核对 SHA-256。服务器按版本保存程序，通过 `current` 链接切换，启动 `agent-service.service`，健康检查失败时切回上一版本。若已构建好 Linux 程序，可用 `-SkipBuild` 只重新上传。

查看运行状态：`sudo systemctl status agent-service`；查看日志：`sudo journalctl -u agent-service -f`。旧版本保存在 `remoteRoot/releases`，确认不再需要时手动清理。

## 配置域名和 Nginx

前端已嵌入 Go 程序，服务器不需要单独部署 `web`。把 `admin.example.com` 和 `chat.example.com` 的 DNS A 记录指向服务器，为两个域名分别配置 HTTPS 证书。下面的后台规则应放入后台域名的 HTTPS `server` 块；客户域名另有完整的建站示例。使用宝塔时不要覆盖面板生成的证书配置，站点已有 `location /` 时应替换原规则，不要重复添加。两个域名都反代到本机的 `127.0.0.1:8083`。

后台域名 `admin.example.com`：

```nginx
location / {
    proxy_pass http://127.0.0.1:8083;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_read_timeout 3600s;
}
```

后台站点若由宝塔“添加站点”生成，还要删除原配置中针对 `gif|jpg|jpeg|png|bmp|swf` 和 `js|css` 的本地缓存 `location`；这些正则规则会优先于 `location /`，使 `/_next/static/` 的脚本或样式在空站点目录返回 404，页面表现为白屏。保留宝塔的 `well-known`、SSL 证书和证书续期规则。可先在服务器比较 `curl -I http://127.0.0.1:8083/_next/static/实际资源路径` 与 `curl -I https://admin.example.com/_next/static/实际资源路径`，两者应返回相同的成功状态。

客户域名 `chat.example.com`：多个渠道共用同一域名，首页进入公共支持中心 `/support/`，不是渠道列表。每个渠道分别使用 `https://chat.example.com/support/chat/?channelId=对应渠道ID`。在宝塔创建“纯静态”站点（不创建 PHP、FTP、数据库），站点根目录使用独立空目录，不要指向 `agent-service` 程序目录。**宝塔“网站目录”和 Nginx 的 `root` 必须是同一个实际存在的目录**，否则宝塔写入的 Let's Encrypt 验证文件会返回 404。下面是申请 HTTPS **之前**可用的完整 Nginx `server` 配置；将域名和站点目录替换为实际值：

```nginx
server {
    listen 80;
    server_name chat.example.com;
    root /www/wwwroot/chat.example.com;

    #CERT-APPLY-CHECK--START
    # 保留宝塔的证书验证目录，供申请和续期使用。
    include /www/server/panel/vhost/nginx/well-known/chat.example.com.conf;
    #CERT-APPLY-CHECK--END
    include /www/server/panel/vhost/nginx/extension/chat.example.com/*.conf;

    #SSL-START SSL相关配置，请勿删除或修改下一行带注释的404规则
    #error_page 404/404.html;
    #SSL-END

    location ~ \.well-known {
        allow all;
    }

    location = / {
        return 302 /support/;
    }

    location = /index.html {
        return 404;
    }

    location ~ ^/(dashboard|workbench)(/|\.html$|$) {
        return 404;
    }

    location ~ ^/api/(dashboard|auth)(/|$) {
        return 404;
    }

    location ~ ^/api/ws/dashboard(/|$) {
        return 404;
    }

    location / {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 3600s;
    }

    access_log /www/wwwlogs/chat.example.com.log;
    error_log /www/wwwlogs/chat.example.com.error.log;
}
```

申请证书前，先确认域名解析到本机、宝塔网站目录与上面的 `root` 一致，并保留宝塔的 `well-known` include、本地处理规则及 `#SSL-START`、`#error_page 404/404.html;`、`#SSL-END` 标记。验证 URL 返回 404 时先检查两处目录，不要反复申请。若证书已申请成功但宝塔提示“未找到标识信息 `#error_page 404/404.html;`”，将这三个 SSL 标记补回站点配置并在宝塔的证书列表点击“部署”，无需重新申请。如果站点已经部署证书，**不要用上面的 HTTP 配置覆盖宝塔生成的 HTTPS/证书配置**；只把六个 `location` 规则放入实际处理 HTTPS 的 `server` 块，替换原有 `location /`。原站点配置中针对图片和 JS/CSS 的本地缓存 `location` 必须移除，否则会绕过反代。客户域名仍需转发 `/api/ws/open` 等客户接口；通用反代会保留这些请求和 WebSocket 升级。只跳转首页不能隐藏后台，必须同时拦截后台页面和接口。`8083` 不要对公网开放，否则可通过服务器 IP 绕过 Nginx 限制；仅允许本机 Nginx 和部署脚本的本机健康检查访问。

修改后执行 `sudo nginx -t` 检查配置，再重载 Nginx。验证 `https://chat.example.com/` 跳转客户页，`https://chat.example.com/index.html`、`/dashboard`、`/workbench` 和 `/api/dashboard/` 返回 404；`https://admin.example.com/workbench` 可正常打开。嵌入其他网站且由浏览器跨域直连 API 时，把该网站的 HTTPS 来源加入 `config.yaml` 的 `server.cors.allowedOrigins`；同源访问无需额外添加。
