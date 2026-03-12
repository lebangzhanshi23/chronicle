---
name: start-my-day
description: 生成每日工作计划，汇总 Chronicle 任务管理系统的待办事项和 DDL 提醒。
---

# start-my-day

生成每日工作计划，汇总 Chronicle 任务管理系统的待办事项和 DDL 提醒。

## CLI 命令

```powershell
# 获取任务列表
chronicle list --json

# 获取统计
chronicle stats --json
```

## Formatter 输出格式

### 紧急任务 (🚨)
格式：`- {任务标题} (⏰ DDL: {deadline})`

### 待办事项 (📝)
格式：`- {任务标题} ({分类})`

### 已完成 (✅)
格式：`- {任务标题}`

## 执行步骤

### 1. 获取任务列表

```powershell
chronicle list --json
```

### 2. 分析任务

解析 JSON，按 deadline 分类：

- **紧急**: deadline < 3 天
- **近期**: deadline 3-7 天
- **普通**: 无 deadline 或 > 7 天

### 3. 输出格式

```
📅 {日期} 工作计划

🚨 紧急任务 ({数量})
- {任务标题} (⏰ DDL: {deadline})

📝 待办事项 ({数量})
- {任务标题} ({分类})

✅ 已完成 ({数量})
- {任务标题}
```

## 使用触发

- "start my day"
- "今日计划"
- "工作汇总"
- "看看有什么任务"
- "待办"
