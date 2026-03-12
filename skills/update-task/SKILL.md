---
name: update-task
description: 更新 Chronicle 任务：添加工作记录、标记完成。
---

# update-task

更新 Chronicle 任务：添加工作记录、标记完成。

## CLI 命令

```powershell
# 添加工作记录
chronicle log <id> <message>

# 更新任务
chronicle update <id> [flags]

# 删除任务
chronicle delete <id>
```

## Formatter 输出格式

更新成功后输出：
```
✅ 任务已更新
- 标题: {title}
- 状态: {status}
- 截止: {deadline 或 "无"}

📝 最新工作记录：
- {时间} {内容}
```

## 参数

| 命令 | 参数 | 说明 |
|-----|-----|-----|
| log | | 添加工作记录 |
| update | --new-status | 状态 (todo/in-progress/done) |
| update | --deadline | 截止时间 |
| update | -c, --category | 分类 |
| delete | | 删除任务 |

## 示例

```powershell
# 添加工作记录
chronicle log <id> "完成了一半"

# 标记完成
chronicle update <id> --new-status done

# 改为进行中
chronicle update <id> --new-status in-progress

# 修改截止时间
chronicle update <id> --deadline "2026-03-20T20:30:00.000Z"

# 删除任务
chronicle delete <id>
```

## 使用触发

- "update task"
- "更新任务"
- "记录进度"
- "标记完成"
- "任务完成"
