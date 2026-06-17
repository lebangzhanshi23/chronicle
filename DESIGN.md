# Tech Design: Agent-Centric Task Management System

## 1. Context & Architecture Goal

系统旨在为 LLM Agent 提供一组稳定、原子化的后台 API，接管用户的日常任务拆解、进度追踪与日志记录。核心原则是：

1. **防止幻觉：** 通过极简的列表接口让 Agent 获取真实 ID；通过复合接口保证状态更新与日志记录的原子性（数据库事务）。
2. **自动化归档：** 提供机器友好的 JSON 日报接口（供 Agent 写总结），以及人类友好的 Markdown 导出接口（供 Obsidian 冷备）。

## 2. Database Design (GORM / SQLite)

采用物理外键或逻辑外键关联的扁平双表结构。

### 2.1 Table: `tasks` (主任务表)

* `id`: `varchar(36)` - 主键，建议使用短 UUID 或 NanoID (方便 Agent 识别)。
* `title`: `varchar(255)` - 任务标题 (例："【BCS】菜单栏隐藏掉部署管理")。
* `category`: `varchar(100)` - 分类/所属产品 (例："BCS")。
* `description`: `text` - 任务背景/受影响服务等详细长文本。
* `targets`: `text` - 任务目标（纯文本存储，可包含换行符）。
* `status`: `varchar(20)` - 状态枚举：`todo`, `in-progress`, `done`。默认 `todo`。
* `deadline`: `datetime` - 预期截止时间（可选）。
* `actual_completed_at`: `datetime` - 实际完成时间（由系统根据 status 自动维护）。
* `created_at`: `datetime` - 创建时间。
* `updated_at`: `datetime` - 更新时间。

### 2.2 Table: `task_logs` (工作记录/日志表)

* `id`: `varchar(36)` - 主键。
* `task_id`: `varchar(36)` - 关联 tasks.id，建普通索引。
* `log_text`: `text` - 工作记录正文 (例："评估完成工作量，修改不多...")。
* `progress_note`: `varchar(100)` - 可选，简短的进度备注。
* `created_at`: `datetime` - 记录发生的时间（精确到分）。

## 3. API Design (RESTful)

所有接口统一返回标准 JSON 结构：`{"code": 0, "msg": "success", "data": {...}}`

### 3.1 核心操作接口 (Agent 侧)

**1. 创建任务**

* **Endpoint:** `POST /api/v1/tasks`
* **Payload:** `{"title": "...", "category": "...", "description": "...", "targets": "..."}`
* **Response:** 返回新创建的 Task 对象（重点包含 `id`）。

**2. 获取活跃任务列表 (Agent "看" 列表的眼睛)**

* **Endpoint:** `GET /api/v1/tasks?status=in-progress,todo`
* **Response:** 必须裁剪字段！只返回 `id`, `title`, `category`, `status`。绝对不要返回 `description`，以节省 Token。

**3. 复合更新接口 (记录日志并更新状态)**

* **Endpoint:** `POST /api/v1/tasks/:id/progress`
* **Payload:** ```json
{
"log_text": "完成编写代码构建之后验证，并且交付",
"mark_as_done": true,  // 如果为 true，则触发状态机将 status 改为 done
"new_status": "done"   // 可选，兼容直接指定状态
}
```

```


* **核心逻辑 (State Machine):** 必须开启 DB 事务。插入一条 `task_logs`；如果检测到状态流转至 `done`，更新 `tasks.status = 'done'` 且 `actual_completed_at = NOW()`。

### 3.2 报表与导出接口 (归档侧)

**4. 每日工作总结 (供 Agent 生成日报文本)**

* **Endpoint:** `GET /api/v1/reports/daily-summary?date=2026-02-26` (默认为今天)
* **Response:** 返回高聚合、防 Token 爆炸的精简 JSON。
```json
{
  "date": "2026-02-26",
  "activities": [
    {
      "task_id": "xxx",
      "task_title": "【BCS】菜单栏隐藏...",
      "status": "done",
      "today_logs": [
        "10:12 - 评估完成工作量...",
        "17:35 - 完成编写代码..."
      ]
    }
  ]
}

```



**5. 每日完成任务 Markdown 导出 (供 Obsidian)**

* **Endpoint:** `GET /api/v1/exports/daily-markdown?date=2026-02-26`
* **核心逻辑:** 1. 查询 `actual_completed_at` 在指定日期（当天 00:00 到 23:59）的所有 `tasks`。
2. 获取这些 tasks 下的**所有历史** `task_logs`（按 `created_at` 升序，并在内存中按日期进行 Map 分组）。
3. 使用 Go `text/template` 将数据注入，直接返回 `Content-Type: text/markdown` 的原始文本流。

## 4. File Structure (Golang)

建议遵循标准 Go 工程规范：

```text
.
├── cmd/server/main.go            # 启动服务，初始化 DB
├── internal/
│   ├── model/                    # GORM 实体与 DTO (request/response struct)
│   ├── handler/                  # HTTP 路由控制器
│   ├── service/                  # 核心业务逻辑 (事务处理、状态机)
│   └── exporter/                 # Markdown 模板渲染引擎
├── templates/
│   └── obsidian_task.tmpl        # Markdown 模板文件
└── data/app.db                   # SQLite 本地库

```

## 5. Acceptance Criteria

* **AC1 (防呆):** `/api/v1/tasks` GET 列表接口的响应体中，坚决不包含长文本字段。
* **AC2 (原子性):** 调用 `/progress` 接口时，如果 `log_text` 写入成功但状态更新失败，必须整体 Rollback，不能产生脏数据。
* **AC3 (归档正确性):** `/exports/daily-markdown` 仅导出今天被标记为 `done` 的任务，但该任务生成的 Markdown 内部，必须完整呈现该任务从创建第一天到今天的**所有**历史日志，并按天分块。
