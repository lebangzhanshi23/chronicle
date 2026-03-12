---
name: create-task
description: 在 Chronicle 任务管理系统中创建新任务。
---

# create-task

在 Chronicle 任务管理系统中创建新任务。

## CLI 命令

```powershell
chronicle create <title> [flags]
```

## Formatter 输出格式

创建成功后输出：
```
✅ 任务已创建
- 标题: {title}
- 分类: {category}
- 截止: {deadline 或 "无"}
- 状态: {status}
- ID: {id}
```

## 参数

| 参数 | 简写 | 说明 |
|-----|-----|-----|
| --category | -c | 分类 |
| --desc | -d | 描述 |
| --deadline | | 截止时间 (ISO8601) |
| --target | -t | 目标 |
| --links | -l | 链接 |

## 示例

```powershell
# 创建任务
chronicle create "新任务" -c "工作"

# 带描述和截止时间
chronicle create "重要任务" -c "工作" -d "详细描述" --deadline "2026-03-20T20:30:00.000Z"
```

## 使用触发

- "create task"
- "创建任务"
- "新建任务"
- "添加任务"
