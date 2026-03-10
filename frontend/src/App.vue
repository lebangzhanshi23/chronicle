<script setup>
import { ref, computed, onMounted } from 'vue'

import TaskBoard from './components/TaskBoard.vue'
import TaskFormModal from './components/TaskFormModal.vue'
import TaskDetailModal from './components/TaskDetailModal.vue'
import ProgressModal from './components/ProgressModal.vue'
import StatsModal from './components/StatsModal.vue'
import StatsBar from './components/StatsBar.vue'
import DailySummaryModal from './components/DailySummaryModal.vue'
import ConfirmModal from './components/ConfirmModal.vue'

// -- GLOBAL STATE --
const tasks = ref([])
const activeTask = ref(null)

// -- MODAL VISIBILITY STATE --
const isFormModalOpen = ref(false)
const isDetailModalOpen = ref(false)
const isProgressModalOpen = ref(false)
const isStatsModalOpen = ref(false)
const isSummaryModalOpen = ref(false)

// Confirm modal state
const isConfirmModalOpen = ref(false)
const confirmModalConfig = ref({
  title: 'Confirm Delete',
  message: 'Are you sure you want to delete this task? This action cannot be undone.',
  confirmText: 'Delete',
  cancelText: 'Cancel',
  confirmDanger: true,
  onConfirm: () => {}
})

const isEditMode = ref(false)

// -- COMPUTED COLUMNS --
const todos = computed(() => tasks.value.filter(t => t.status === 'todo'))
const inProgress = computed(() => tasks.value.filter(t => t.status === 'in-progress'))
const done = computed(() => tasks.value.filter(t => t.status === 'done'))

// -- DATA FETCHING --
async function loadTasks() {
  const [activeRes, doneRes] = await Promise.all([
    fetch('/api/v1/tasks?status=todo,in-progress').then(r => r.json()),
    fetch('/api/v1/tasks?status=done').then(r => r.json())
  ])
  
  let allTasks = []
  if (activeRes.code === 0 && activeRes.data) {
    allTasks = [...allTasks, ...activeRes.data]
  }
  if (doneRes.code === 0 && doneRes.data) {
    allTasks = [...allTasks, ...doneRes.data]
  }
  tasks.value = allTasks
}

async function loadTaskDetail(id) {
  const res = await fetch(`/api/v1/tasks/${id}`).then(r => r.json())
  if (res.code === 0) {
    activeTask.value = res.data
  }
}

onMounted(() => {
  loadTasks()
})

// -- HANDLERS --
function handleNewTaskClick() {
  activeTask.value = null
  isEditMode.value = false
  isFormModalOpen.value = true
}

async function handleTaskClick(task) {
  await loadTaskDetail(task.id)
  isDetailModalOpen.value = true
}

function handleOpenEdit() {
  isDetailModalOpen.value = false
  isEditMode.value = true
  isFormModalOpen.value = true
}

async function handleDeleteTask() {
  isDetailModalOpen.value = false
  confirmModalConfig.value = {
    title: 'Delete Task',
    message: `Are you sure you want to delete "${activeTask.value.title}"? This action cannot be undone.`,
    confirmText: 'Delete',
    cancelText: 'Cancel',
    confirmDanger: true,
    onConfirm: async () => {
      await fetch(`/api/v1/tasks/${activeTask.value.id}`, { method: 'DELETE' })
      isConfirmModalOpen.value = false
      loadTasks()
    }
  }
  isConfirmModalOpen.value = true
}

function handleOpenProgress() {
  isDetailModalOpen.value = false
  isProgressModalOpen.value = true
}

async function handleProgressSubmit() {
  isProgressModalOpen.value = false
  loadTasks()
}

async function handleDeleteWorklog(worklogId) {
  if (!worklogId) return
  isDetailModalOpen.value = false
  confirmModalConfig.value = {
    title: 'Delete Worklog',
    message: 'Are you sure you want to delete this worklog? This action cannot be undone.',
    confirmText: 'Delete',
    cancelText: 'Cancel',
    confirmDanger: true,
    onConfirm: async () => {
      await fetch(`/api/v1/worklogs/${worklogId}`, { method: 'DELETE' })
      isConfirmModalOpen.value = false
      await loadTaskDetail(activeTask.value.id)
    }
  }
  isConfirmModalOpen.value = true
}
</script>

<template>
  <nav class="sticky top-0 z-40 w-full backdrop-blur-md flex-none transition-colors duration-500 lg:z-50 border-b border-dark-border bg-dark-bg/85">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-xl bg-gradient-to-tr from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-indigo-500/30 ring-1 ring-white/20">
            C
          </div>
          <span class="font-bold text-xl tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-slate-100 to-slate-400">Chronicle</span>
        </div>
        <button @click="handleNewTaskClick" class="flex items-center gap-2 bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-400 hover:to-indigo-500 text-white px-5 py-2 rounded-xl text-sm font-semibold transition-all shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/40 active:scale-95 border border-indigo-400/20">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          New Task
        </button>
      </div>
    </div>
  </nav>

  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <StatsBar />
    <TaskBoard 
      :todos="todos" 
      :inProgress="inProgress" 
      :done="done" 
      @task-click="handleTaskClick" 
      @start="loadTasks"
      @view-stats="isStatsModalOpen = true"
      @view-summary="isSummaryModalOpen = true"
    />
  </main>

  <!-- Modals -->
  <TaskFormModal 
    :is-open="isFormModalOpen"
    :task="isEditMode ? activeTask : null"
    @close="isFormModalOpen = false"
    @refresh="loadTasks"
  />

  <TaskDetailModal 
    :is-open="isDetailModalOpen"
    :task="activeTask"
    @close="isDetailModalOpen = false"
    @edit="handleOpenEdit"
    @delete="handleDeleteTask"
    @progress="handleOpenProgress"
    @start="loadTasks"
    @delete-worklog="handleDeleteWorklog"
    @refresh="loadTasks"
  />

  <ProgressModal 
    :is-open="isProgressModalOpen"
    :task="activeTask"
    @close="isProgressModalOpen = false"
    @submit="handleProgressSubmit"
  />

  <StatsModal
    v-if="isStatsModalOpen"
    @close="isStatsModalOpen = false"
  />

  <DailySummaryModal
    v-if="isSummaryModalOpen"
    @close="isSummaryModalOpen = false"
  />

  <ConfirmModal
    :is-open="isConfirmModalOpen"
    :title="confirmModalConfig.title"
    :message="confirmModalConfig.message"
    :confirm-text="confirmModalConfig.confirmText"
    :cancel-text="confirmModalConfig.cancelText"
    :confirm-danger="confirmModalConfig.confirmDanger"
    @confirm="confirmModalConfig.onConfirm"
    @cancel="isConfirmModalOpen = false"
  />
</template>
