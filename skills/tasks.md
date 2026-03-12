# Chronicle Tasks Skill

调用 Chronicle API 管理任务。

## 环境要求

- Chronicle 服务运行在 `http://localhost:8080`
- Python 3.x（用于调用 API，避免中文乱码）

## API 列表

### 1. 创建任务

```python
import urllib.request, json

data = {
    "title": "任务标题",
    "category": "分类",
    "description": "描述（可选）",
    "deadline": "2026-03-20T20:30:00.000Z"  # ISO 格式，可选
}

req = urllib.request.Request(
    'http://localhost:8080/api/v1/tasks',
    data=json.dumps(data).encode('utf-8'),
    headers={'Content-Type': 'application/json'}
)
print(urllib.request.urlopen(req).read().decode('utf-8'))
```

### 2. 查询任务

```python
import urllib.request

# 全部进行中任务
req = urllib.request.Request('http://localhost:8080/api/v1/tasks?status=todo,in-progress')
print(urllib.request.urlopen(req).read().decode('utf-8'))

# 已完成任务
req = urllib.request.Request('http://localhost:8080/api/v1/tasks?status=done')
print(urllib.request.urlopen(req).read().decode('utf-8'))
```

### 3. 更新任务（添加记录/完成）

```python
import urllib.request, json

# 添加记录（自动改为进行中）
data = {"log_text": "工作内容", "new_status": "in-progress"}

# 标记完成
data = {"log_text": "任务完成", "mark_as_done": True}

req = urllib.request.Request(
    'http://localhost:8080/api/v1/tasks/{task_id}/progress',
    data=json.dumps(data).encode('utf-8'),
    headers={'Content-Type': 'application/json'}
)
print(urllib.request.urlopen(req).read().decode('utf-8'))
```

### 4. 删除任务

```python
import urllib.request

req = urllib.request.Request(
    'http://localhost:8080/api/v1/tasks/{task_id}',
    method='DELETE'
)
print(urllib.request.urlopen(req).read().decode('utf-8'))
```

## 分类建议

- 工作
- 许可管理
- DevOps
- 学习
- 其他

## Deadline 格式

必须使用 ISO 8601 格式：`YYYY-MM-DDTHH:MM:SS.000Z`

例如：`2026-03-20T20:30:00.000Z`

## 状态说明

| 状态 | 说明 |
|-----|------|
| todo | 待处理 |
| in-progress | 进行中 |
| done | 已完成 |
