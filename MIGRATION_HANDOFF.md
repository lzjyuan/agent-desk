# agent-desk / agent-service 路径迁移交接

本文根据当前仓库和此前排障记录整理。**已核实**指本地代码与脚本；服务器、数据库和域名的**当前运行状态未在本次交接中远程核实**。不要把聊天记录中的历史截图当作生产现状。

## 1. 项目位置与边界

- 历史本地路径：`D:\Project-PHP\MexicoChat\agent-desk`；该目录当前已不存在。
- 当前本地仓库：`D:\Project-Go\agent-service`。以下本地命令均从该目录执行。迁移源码目录本身不要求修改 `go.mod` 的 `module agent-desk` 或 Go import；这是模块标识，不是磁盘路径。
- 历史服务器发布根目录：`/www/wwwroot/agent-service`；实际值以 `deploy/target.local.json` 的 `remoteRoot` 和服务器上的 systemd unit 为准。
- 代码约定：先读仓库根目录 `AGENTS.md`；一般操作看 `README.md`，服务器部署与 Nginx 规则看 `deploy/README.md`。

## 2. 代码和运行结构

| 位置 | 作用 |
| --- | --- |
| `cmd/server/`、`internal/` | Go 服务、API、业务、GORM 模型和启动迁移 |
| `web/` | Next.js 16 管理后台、工作台、客户页及 SDK |
| `flowgram-editor/` | Workflow 编辑器源码，构建产物进入 `web/public/flowgram-editor/` |
| `config/config.example.yaml` | 配置模板；实际 `config/config.yaml` 被 Git 忽略 |
| `deploy/deploy_backend.ps1`、`deploy/deploy_release.sh` | Windows 构建上传、Linux 版本切换和健康检查 |

生产前端不是独立 Node 服务。发布时先构建 Workflow 编辑器、SDK 和 `web/out`，再由 `web/embed.go` 将静态资源编入 Go 二进制；服务器运行 Go 程序即可。`-tags dev` 下则使用 `web/embed_dev.go` 读取本地文件。缺少 `web/out` 时不要直接做普通 Go 发布构建。

## 3. 只迁本地源码目录

1. 确认新目录包含 `.git`、`go.mod`、`web/`、`flowgram-editor/`、`deploy/`。在新目录执行 `git status --short`；不要覆盖其他仓库的改动。
2. 从旧目录单独迁移**未纳入 Git**、仍需使用的文件：`config/config.yaml`、`deploy/target.local.json`、本地数据 `data/`，以及 SSH 密钥等外置依赖。不要把密钥或生产配置提交到 Git；若这些文件已经在新目录，先比较内容而不是直接覆盖。
3. 从新目录运行下面的构建命令。`npx` 缓存报 `ENOENT` 时，先按 `deploy/deploy_backend.ps1` 使用独立 npm 缓存，不要把该报错误判为 Go 服务问题。
4. 检查 IDE、终端、任务计划、CI 和自写脚本里指向旧路径的**外部配置**。仓库内未发现硬编码的旧绝对路径；不必批量替换 `agent-desk` 模块名。

```powershell
Set-Location 'D:\Project-Go\agent-service'
git status --short
Push-Location flowgram-editor
npx --yes pnpm@10.30.2 install --frozen-lockfile
npx --yes pnpm@10.30.2 build
Pop-Location
Push-Location web
npx --yes pnpm@10.30.2 install --frozen-lockfile
npx --yes pnpm@10.30.2 build:sdk
npx --yes pnpm@10.30.2 build
Pop-Location
go test ./...
go build -o dist/agent-service.exe ./cmd/server
```

本地启动：`go run -tags dev ./cmd/server -config config/config.yaml`。配置中的相对路径按**进程工作目录**解析；从仓库根目录启动时，默认 SQLite 为 `data/app.db`，本地附件为 `data/storage`。本机如改用 MySQL，先创建目标数据库，并确认 `db.dsn` 指向它；GORM 启动迁移负责建表，不负责建库。

## 4. 服务器发布与服务器路径迁移

如果**只改本地源码路径**、服务器 `remoteRoot` 不变，现有 Nginx、systemd、数据库和 Qdrant 地址通常无需因本地搬家而修改。确认 `deploy/target.local.json` 和本地配置后，从仓库根目录执行 `./deploy/deploy_backend.ps1`（PowerShell 写法：`.\deploy\deploy_backend.ps1`）。该脚本从自身位置推算仓库根目录，不依赖旧的 Windows 绝对路径。

如果**服务器发布根目录也要改**，先备份并核实以下内容，再切换：

- `remoteRoot/shared/config/config.yaml`：生产配置和密钥；`remoteRoot/shared/data/`：SQLite 数据库及 WAL/SHM、上传文件等持久数据。用 MySQL 时另行备份 MySQL；用 Qdrant 时另行备份其数据目录（部署文档示例为 `/var/lib/qdrant`）。不能只迁 `releases/` 或 `current` 链接。
- `remoteRoot/releases/`：版本二进制；`remoteRoot/current`：当前版本符号链接。程序的 systemd `WorkingDirectory` 应为 `remoteRoot/shared`，`ExecStart` 应指向新根目录的 `current/agent-service` 与 `shared/config/config.yaml`。
- **现有 unit 不会被发布脚本自动重写**：`deploy/deploy_release.sh` 只有在 `/etc/systemd/system/agent-service.service` 不存在时才创建它。变更服务器根目录必须手动核对并更新 unit，然后执行 `systemctl daemon-reload` 和 `systemctl restart agent-service`。
- 核对 `shared` 目录及配置文件的运行用户权限。发布脚本使用 `target.local.json` 的 `serviceUser`；此前部署文档示例是 `www`，不可假定新服务器相同。
- Nginx 通常仍反代本机 `127.0.0.1:8083`，不直接依赖发布根目录；若同时改域名或端口，再检查证书、站点根目录、反代规则和 `server.cors.allowedOrigins`。不要用示例 Nginx `server` 覆盖宝塔现有 SSL 配置。

发布脚本会校验包的 SHA-256、切换 `current`、重启服务，并循环检查 `http://127.0.0.1:8083/api/config`；健康检查失败会尝试回滚二进制。服务启动后 GORM 会执行 `AutoMigrate`，**二进制回滚不等于数据库回滚**；变更数据库结构前先备份。

## 5. 历史配置与故障线索

以下来自此前讨论，迁移时逐项确认，不代表当前生产一定如此：

- 后台域名曾使用 `adminsar.honertop.com`，客户域名曾使用 `chat.honertop.com`；两者反代同一 Go 端口。客户域名必须拦截后台页面及后台 API，Web 渠道的 `publicBaseUrl` 应为客户域名且为必填。渠道 ID 不随本地路径迁移而改变。
- 曾从 SQLite 转向 MySQL；以服务器实际 `db.type`、`db.dsn` 为准。DMS 是数据库管理入口，不是数据库引擎。迁库和迁目录是两件事。
- 默认向量后端是 Qdrant，Go 客户端走本机 gRPC `6334`；连接拒绝先查 Qdrant 进程，`Model not exist` 则先查 AI 模型名、区域和密钥，不能一概归咎于向量库。
- 服务器曾从 `Asia/Shanghai` 切到沙特时区；核实 `timedatectl`、systemd 环境和页面时间，不凭服务器物理位置推断时区。
- MCP Server 目前从 `config.yaml` 的 `mcp.servers` 读取；后台 “MCP tools” 页只负责测试、列工具和调用。此前讨论过改成页面新增并持久化到数据库，**尚不能把它当作已实现功能**。
- 证书验证 404、`/_next/static/` 404 或白屏时，先查宝塔站点 `root`、`.well-known` include，以及是否有本地图片/JS/CSS 正则 `location` 抢在反代前处理请求。浏览器扩展的 `Unchecked runtime.lastError` 不等同于本服务报错。

## 6. 迁移后验收

```powershell
Set-Location 'D:\Project-Go\agent-service'
git status --short
Test-Path .\config\config.yaml
Test-Path .\deploy\target.local.json
Test-Path .\web\out\index.html
```

```bash
systemctl cat agent-service
systemctl status agent-service --no-pager
journalctl -u agent-service -n 80 --no-pager
curl -fsS http://127.0.0.1:8083/api/config
nginx -t
```

最后从浏览器分别核对后台登录、客户聊天页、SDK 静态资源、WebSocket、知识库检索和一次实际消息回复。服务健康检查通过只证明 HTTP 入口可用，不证明 MySQL、Qdrant、模型 API 或全部渠道可用。

交接时只传递**配置文件所在位置和密钥保管方式**，不要把数据库密码、API Key、SSH 私钥写进本文。服务器当前根目录、数据库类型、备份位置及责任人仍需接手人现场确认。
