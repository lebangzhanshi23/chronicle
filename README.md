# Chronicle

Chronicle 这是一个专为 LLM Agent 打造的**任务管理与追踪系统 (Agent-Centric Task Management System)**。
系统提供了一组稳定、原子化的后台 API，接管用户的日常任务拆解、进度追踪与日志记录，并提供了一个简单易用的前端网页界面进行管理。

核心设计原则：
1. **防止幻觉：** 通过极简的列表接口让 Agent 获取真实 ID；通过复合接口保证状态更新与日志记录的原子性（数据库事务）。
2. **自动化归档：** 提供机器友好的 JSON 日报接口（供 Agent 写总结），以及人类友好的 Markdown 导出接口（供 Obsidian 或其他笔记软件冷备）。

## ✨ 核心特性

- 🤖 **Agent-Friendly API**: 专门设计的防止大模型产生幻觉的接口格式，严格限制 Token 消耗。
- 🔄 **原子化事务**: 状态流转和日志记录确保数据库层面的强一致性，避免脏数据。
- 📊 **简单易用的 UI**: 自带基于 TailwindCSS 构建的可视化前端界面，支持查看任务、录入进度、查看历史以及生成日报。
- 📦 **轻量级**: 使用 Go 和 SQLite 构建，仅需单一可执行文件即可运行，无需复杂的环境依赖。
- 🐳 **Docker 支持**: 提供 Dockerfile 和 docker-compose 配置，支持容器化部署。
- 📤 **数据导入导出**: 支持 `export` / `import` 命令，方便数据备份与迁移。

## 🛠 技术栈

- **后端**: Golang, Gin 框架, GORM
- **数据库**: SQLite (本地单文件存储, 纯 Go 驱动)
- **前端**: Vue 3, Vite, Tailwind CSS v4

## 📂 项目结构

```text
.
├── cmd/                          # 命令行工具逻辑
│   ├── root.go                   # 根命令定义
│   ├── server.go                 # 服务器启动命令
│   ├── tasks.go                  # 任务管理相关命令
│   ├── version.go                # 版本信息命令
│   └── export.go                 # 数据导入导出命令
├── internal/                     # 核心业务逻辑
│   ├── config/                   # 配置管理（支持环境变量和命令行参数）
│   ├── handler/                  # HTTP 请求处理与响应 (Controller)
│   ├── model/                    # GORM 实体与接口 DTO
│   └── service/                  # 核心业务逻辑 (事务处理、状态机等)
├── frontend/                     # 现代前端 Vue 3 工程目录
│   ├── src/                      # Vue 组件和主入口
│   ├── package.json              # Node.js 依赖配置
│   └── dist/                     # Vite 构建输出目录 (Go Server 通过静态代理服务该目录)
├── skills/                       # AI Agent 集成 Skill 文件目录
├── .github/workflows/            # GitHub Actions CI 配置
├── Dockerfile                    # Docker 镜像构建文件
├── docker-compose.yml            # Docker Compose 编排配置
├── Makefile                      # 构建 Makefile
├── DESIGN.md                     # 技术设计文档
├── main.go                       # 项目主入口
└── go.mod                        # Go 模块定义
```

## 🚀 快速开始

### 1. 环境准备

请确保您本地已经安装了：
- [Golang](https://go.dev/dl/) (版本要求 >= 1.24)
- [Node.js](https://nodejs.org/) (版本建议 >= 18) 和 npm

### 2. 安装 Chronicle

#### 方式一：go install（推荐）

```bash
go install github.com/yuyudeqiu/chronicle@latest
```

#### 方式二：手动编译

```bash
# 使用 Makefile（推荐）
make build

# 或直接使用 go build
go build -o chronicle main.go
```

#### 方式三：Docker

```bash
docker-compose up -d
```

### 3. 数据存储配置

默认情况下，数据存放在项目目录下的 `data/app.db`。

可通过以下方式自定义数据存储路径（优先级：命令行参数 > 环境变量 > 默认值）：

```bash
# 方式一：环境变量
export CHRONICLE_DATA_DIR=/custom/path

# 方式二：命令行参数
chronicle --data-dir /custom/path create "新任务"
```

### 4. 构建前端页面

因为使用了现代化的 Vite + Vue 3 架构，在运行 Web 服务器前，需要先编译前端资源：

```bash
cd frontend
npm install
npm run build
cd ..
```
*(编译后的产物会存放在 `frontend/dist` 记录中，并且自动由后端的 Go 服务提供代理访问)*

### 5. 启动服务

```bash
# 启动 Web 服务器，默认端口为 8080
chronicle server

# 或使用环境变量指定数据目录
CHRONICLE_DATA_DIR=/path/to/data chronicle server
```

服务启动后，可以直接通过浏览器访问主操作界面： http://localhost:8080/

### 6. 使用命令行工具 (CLI)

除了网页界面，你也可以直接在终端管理任务：

```bash
# 查看所有命令
chronicle --help

# 创建任务
chronicle create "完成项目重构" -c 开发 -d "使用 Cobra 优化 CLI"

# 列出进行中的任务
chronicle list in-progress

# 查看任务详情
chronicle get <task_id>

# 添加执行日志
chronicle log <task_id> "完成了 CLI 重构"

# 更新任务状态
chronicle update <task_id> --new-status done

# 更新任务信息
chronicle update <task_id> -c 运维 -d "更新描述" --deadline 2026-01-01T00:00:00Z

# 删除任务
chronicle delete <task_id>

# 查看每日总结
chronicle summary
chronicle summary 2026-01-15

# 查看周报
chronicle weekly-summary

# 查看任务统计
chronicle stats

# 查看版本信息
chronicle version

# 导出数据库
chronicle export /path/to/backup.db

# 导入数据库（会覆盖现有数据）
chronicle import /path/to/backup.db
chronicle import /path/to/backup.db --force  # 跳过确认
```

#### 全局参数

| 参数 | 说明 |
|------|------|
| `--data-dir <path>` | 指定数据目录（默认 `data/`，也可通过 `CHRONICLE_DATA_DIR` 环境变量设置） |
| `-o, --json` | 以 JSON 格式输出结果（便于 Agent 解析） |

### 📚 AI Agent 集成

如果你是 AI Agent 想接入 Chronicle，推荐使用 `skills/` 目录下的 Skill 示例代码。

**快速集成示例：**

```bash
# 列出任务（JSON 格式，便于解析）
chronicle list --json

# 创建任务
chronicle create "任务标题" -c "分类"

# 添加记录
chronicle log <id> "工作内容"

# 标记完成
chronicle update <id> --new-status done
```

详细文档见：[skills/](./skills/) 目录
