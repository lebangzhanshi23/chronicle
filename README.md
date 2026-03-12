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
- 📝 **Markdown 导出**: 每日任务可直接导出为 Markdown 格式，完美兼容 Obsidian 等本地知识库软件。
- 📦 **轻量级**: 使用 Go 和 SQLite 构建，仅需单一可执行文件即可运行，无需复杂的环境依赖。

## 🛠 技术栈

- **后端**: Golang, Gin 框架, GORM
- **数据库**: SQLite (本地单文件存储)
- **前端**: Vue 3, Vite, Tailwind CSS v4

## 📂 项目结构

```text
.
├── cmd/server/main.go            # 启动控制点，存放 HTTP 路由和应用初始化逻辑
├── internal/                     # 核心业务逻辑
│   ├── handler/                  # HTTP 请求处理与响应 (Controller)
│   ├── model/                    # GORM 实体与接口 DTO
│   └── service/                  # 核心业务逻辑 (事务处理、状态机等)
├── frontend/                     # 现代前端 Vue 3 工程目录
│   ├── src/                      # Vue 组件和主入口
│   ├── package.json              # Node.js 依赖配置
│   └── dist/                     # Vite 构建输出目录 (Go Server 通过静态代理服务该目录)
├── templates/                    # Go Template 模板目录
│   └── obsidian_task.tmpl        # 导出为 Obsidian Markdown 的渲染模板
├── data/                         # 数据库文件存放目录 (运行时自动生成)
│   └── app.db                    # SQLite 本地库
└── bin/                          # 二进制编译输出目录
```

## 🚀 快速开始

### 1. 环境准备

请确保您本地已经安装了：
- [Golang](https://go.dev/dl/) (版本建议 >= 1.20)
- [Node.js](https://nodejs.org/) (版本建议 >= 18) 和 npm

### 2. 构建前端页面

因为使用了现代化的 Vite + Vue 3 架构，在运行 Go 服务前，必须先编译前端资源：

```bash
cd frontend
npm install
npm run build
cd ..
```
*(编译后的产物会存放在 `frontend/dist` 记录中，并且自动由后端的 Go 服务提供代理访问)*

### 3. 下载后端依赖并运行服务

```bash
go mod tidy
# 启动本地开发服务器，默认端口为 8080
go run cmd/server/main.go
```

服务启动后，可以直接通过浏览器访问主操作界面： http://localhost:8080/

### 4. 编译可执行文件

如果您想将系统编译为一个独立的二进制文件在后台运行：

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## 🤖 接口说明 (供 Agent 使用)

系统主要提供了以下几类核心接口（详细 Schema 请参考 `DESIGIN.md`）：

1. **获取任务列表**: `GET /api/v1/tasks?status=in-progress,todo` (仅返回精简信息，防止 Token 爆炸)
2. **创建新任务**: `POST /api/v1/tasks`
3. **追加执行日志并标记进度**: `POST /api/v1/tasks/:id/progress` (复合更新，保证原子性)
4. **获取每日 JSON 总结**: `GET /api/v1/reports/daily-summary?date=YY-MM-DD`
5. **获取 Markdown 导出**: `GET /api/v1/exports/daily-markdown?date=YY-MM-DD`

### 📚 AI Agent 集成

如果你是 AI Agent 想接入 Chronicle，推荐使用 `skills/tasks.md` 中的 Python 示例代码，支持创建任务、查询进度、记录日志等操作。

**快速集成示例：**

```python
import urllib.request, json

# 查询进行中的任务
req = urllib.request.Request('http://localhost:8080/api/v1/tasks?status=todo,in-progress')
print(urllib.request.urlopen(req).read().decode('utf-8'))
```

详细文档见：[skills/tasks.md](./skills/tasks.md)
