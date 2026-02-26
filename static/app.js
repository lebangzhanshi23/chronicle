document.addEventListener('DOMContentLoaded', () => {
    loadTasks();
});

function loadTasks() {
    fetch('/api/v1/tasks?status=todo,in-progress')
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                renderTasks(data.data || []);
            }
        });
    fetch('/api/v1/tasks?status=done')
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                renderDoneTasks(data.data || []);
            }
        });
}

function renderTasks(tasks) {
    const todoList = document.getElementById('todo-list');
    const inprogressList = document.getElementById('inprogress-list');

    todoList.innerHTML = '';
    inprogressList.innerHTML = '';

    let todoCount = 0;
    let inprogressCount = 0;

    tasks.forEach(task => {
        const card = createTaskCard(task);
        if (task.status === 'todo') {
            todoList.appendChild(card);
            todoCount++;
        } else if (task.status === 'in-progress') {
            inprogressList.appendChild(card);
            inprogressCount++;
        }
    });

    document.getElementById('todo-count').textContent = todoCount;
    document.getElementById('inprogress-count').textContent = inprogressCount;
}

function renderDoneTasks(tasks) {
    const doneList = document.getElementById('done-list');
    doneList.innerHTML = '';

    let doneCount = 0;
    tasks.forEach(task => {
        const card = createTaskCard(task);
        doneList.appendChild(card);
        doneCount++;
    });

    document.getElementById('done-count').textContent = doneCount;
}

function createTaskCard(task) {
    const div = document.createElement('div');
    div.className = 'glass-card rounded-xl p-4 border border-dark-border hover:border-indigo-500/50 transition-all cursor-pointer group hover:shadow-lg hover:shadow-indigo-500/10 task-enter' + (task.status === 'done' ? ' opacity-75 grayscale-[0.2]' : '');
    div.onclick = () => openTaskDetailModal(task.id);

    let statusBadge = '';
    if (task.status === 'in-progress') {
        statusBadge = '<span class="flex h-2 w-2 relative ml-auto"><span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span><span class="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span></span>';
    } else if (task.status === 'done') {
        statusBadge = '<span class="ml-auto flex items-center text-emerald-500"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg></span>';
    }

    div.innerHTML = `
        <div class="flex justify-between items-start mb-2">
            <h3 class="font-medium text-slate-200 group-hover:text-white transition-colors leading-tight ${task.status === 'done' ? 'line-through text-slate-400' : ''}">${escapeHtml(task.title)}</h3>
        </div>
        <div class="flex items-center gap-2 mt-3">
            <span class="text-[10px] uppercase tracking-wider font-semibold text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20">${escapeHtml(task.category)}</span>
            ${statusBadge}
        </div>
    `;
    return div;
}

// Modals
function openNewTaskModal() {
    document.getElementById('new-task-form').reset();

    // Set default deadline to 7 days from now at 20:30
    const dlInput = document.getElementById('task-deadline');
    const d = new Date();
    d.setDate(d.getDate() + 7);
    d.setHours(20, 30, 0, 0); // local time 20:30

    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');

    dlInput.value = `${year}-${month}-${day}T${hours}:${minutes}`;

    document.getElementById('new-task-modal').classList.remove('hidden');
}

function openProgressModal(task) {
    document.getElementById('progress-form').reset();
    document.getElementById('progress-task-id').value = task.id;
    document.getElementById('progress-task-title').textContent = task.title;

    // Select default radio based on current status
    const radios = document.getElementsByName('progress_action');
    radios[0].checked = true; // Default to update log
    document.getElementById('progress-modal').classList.remove('hidden');
}

function openTaskDetailModal(taskId) {
    fetch(`/api/v1/tasks/${taskId}`)
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                const task = data.data;
                document.getElementById('detail-task-title').textContent = task.title;
                document.getElementById('detail-task-category').textContent = task.category;

                const statusEl = document.getElementById('detail-task-status');
                statusEl.textContent = task.status;
                if (task.status === 'done') {
                    statusEl.className = 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-emerald-400 bg-emerald-500/10 px-2 py-1 rounded border border-emerald-500/20';
                } else if (task.status === 'in-progress') {
                    statusEl.className = 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20';
                } else {
                    statusEl.className = 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-slate-400 bg-slate-500/10 px-2 py-1 rounded border border-slate-500/20';
                }

                if (task.deadline) {
                    document.getElementById('detail-task-deadline-container').classList.remove('hidden');
                    document.getElementById('detail-task-deadline').textContent = new Date(task.deadline).toLocaleString();
                } else {
                    document.getElementById('detail-task-deadline-container').classList.add('hidden');
                }

                if (task.actual_completed_at) {
                    document.getElementById('detail-task-completed-container').classList.remove('hidden');
                    document.getElementById('detail-task-completed').textContent = new Date(task.actual_completed_at).toLocaleString();
                } else {
                    document.getElementById('detail-task-completed-container').classList.add('hidden');
                }

                if (task.description) {
                    document.getElementById('detail-task-desc-container').classList.remove('hidden');
                    document.getElementById('detail-task-desc').textContent = task.description;
                } else {
                    document.getElementById('detail-task-desc-container').classList.add('hidden');
                }

                if (task.targets) {
                    document.getElementById('detail-task-targets-container').classList.remove('hidden');
                    document.getElementById('detail-task-targets').textContent = task.targets;
                } else {
                    document.getElementById('detail-task-targets-container').classList.add('hidden');
                }

                const logsContainer = document.getElementById('detail-task-logs');
                logsContainer.innerHTML = '';
                if (task.logs && task.logs.length > 0) {
                    task.logs.forEach(log => {
                        const logDiv = document.createElement('div');
                        logDiv.className = 'bg-dark-bg p-3 rounded-lg border border-dark-border text-sm';
                        const timeStr = new Date(log.created_at).toLocaleString();
                        logDiv.innerHTML = `<div class="text-xs text-slate-500 mb-1">${timeStr}</div><div class="text-slate-200 whitespace-pre-wrap">${escapeHtml(log.log_text)}</div>`;
                        logsContainer.appendChild(logDiv);
                    });
                } else {
                    logsContainer.innerHTML = '<div class="text-sm text-slate-500 italic">No worklogs yet.</div>';
                }

                const btnUpdate = document.getElementById('btn-update-progress');
                if (task.status !== 'done') {
                    btnUpdate.classList.remove('hidden');
                    btnUpdate.onclick = () => {
                        closeModals();
                        openProgressModal(task);
                    };
                } else {
                    btnUpdate.classList.add('hidden');
                }

                const btnDelete = document.getElementById('btn-delete-task');
                btnDelete.classList.remove('hidden');
                btnDelete.onclick = () => confirmDeleteTask(task.id);

                document.getElementById('task-detail-modal').classList.remove('hidden');
            } else {
                showToast('Failed to load task details');
            }
        });
}

function closeModals() {
    document.getElementById('new-task-modal').classList.add('hidden');
    document.getElementById('progress-modal').classList.add('hidden');
    document.getElementById('task-detail-modal')?.classList.add('hidden');
}

// Actions
function handleCreateTask(e) {
    e.preventDefault();
    const title = document.getElementById('task-title').value;
    const category = document.getElementById('task-category').value;
    const desc = document.getElementById('task-desc').value;
    const targets = document.getElementById('task-targets').value;
    const deadlineInput = document.getElementById('task-deadline').value;

    let deadline = null;
    if (deadlineInput) {
        deadline = new Date(deadlineInput).toISOString();
    }

    fetch('/api/v1/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            title: title,
            category: category,
            description: desc,
            targets: targets,
            deadline: deadline
        })
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                closeModals();
                loadTasks();
                showToast('Task created successfully');
            }
        });
}

function handleUpdateProgress(e) {
    e.preventDefault();
    const taskId = document.getElementById('progress-task-id').value;
    const logText = document.getElementById('progress-log').value;

    // Find which radio is checked
    let markAsDone = false;
    const radios = document.getElementsByName('progress_action');
    for (const radio of radios) {
        if (radio.checked && radio.value === 'done') {
            markAsDone = true;
            break;
        }
    }

    fetch(`/api/v1/tasks/${taskId}/progress`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            log_text: logText,
            mark_as_done: markAsDone,
            new_status: markAsDone ? 'done' : 'in-progress'
        })
    })
        .then(res => res.json())
        .then(data => {
            if (data.code === 0) {
                closeModals();
                loadTasks();
                showToast(markAsDone ? 'Task marked as done! 🎉' : 'Progress logged');
            }
        });
}

function confirmDeleteTask(taskId) {
    if (confirm("Are you sure you want to delete this task? All of its worklogs will also be permanently deleted.")) {
        fetch(`/api/v1/tasks/${taskId}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(data => {
                if (data.code === 0) {
                    closeModals();
                    loadTasks();
                    showToast('Task deleted successfully');
                } else {
                    showToast('Failed to delete task: ' + data.msg);
                }
            });
    }
}

// Utils
function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

let toastTimeout;
function showToast(msg) {
    const toast = document.getElementById('toast');
    document.getElementById('toast-msg').textContent = msg;

    toast.classList.remove('translate-y-20', 'opacity-0');

    clearTimeout(toastTimeout);
    toastTimeout = setTimeout(() => {
        toast.classList.add('translate-y-20', 'opacity-0');
    }, 3000);
}
