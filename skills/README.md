# Chronicle Skills

本目录包含 Chronicle CLI 的 AI Agent 集成示例，供其他 AI Agent 使用。

## 环境要求

- Chronicle 服务运行在 `http://localhost:8080`
- CLI 路径：`{chronicle项目目录}/bin/chronicle`

## Skills

### start-my-day
生成每日工作计划，汇总待办事项和 DDL 提醒。

```powershell
# 获取任务列表
.\bin\chronicle.exe list --json

# 获取统计
.\bin\chronicle.exe stats --json
```

### create-task
创建新任务。

```powershell
.\bin\chronicle.exe create "任务标题" -c "分类" -d "描述" --deadline "2026-03-20T20:30:00.000Z"
```

### update-task
更新任务：添加工作记录、标记完成。

```powershell
# 添加记录
.\bin\chronicle.exe log <id> "工作内容"

# 标记完成
.\bin\chronicle.exe update <id> --new-status done
```

### log-daily-work
获取今日工作记录，写入 Obsidian Daily Log。

```powershell
.\bin\chronicle.exe summary --json
```

## 通用参数

- `--json` / `-o`: 输出 JSON 格式（方便 AI 解析）
- `--help`: 查看帮助

## 状态说明

| 状态 | 说明 |
|-----|------|
| todo | 待处理 |
| in-progress | 进行中 |
| done | 已完成 |

## Deadline 格式

ISO 8601: `YYYY-MM-DDTHH:MM:SS.000Z`

例如：`2026-03-20T20:30:00.000Z`
