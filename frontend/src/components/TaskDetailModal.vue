<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  isOpen: Boolean,
  task: Object
})

const emit = defineEmits(['close', 'edit', 'delete', 'progress', 'delete-worklog', 'start', 'refresh'])

const toastMsg = ref('')
const showToast = ref(false)

function showToastMsg(msg) {
  toastMsg.value = msg
  showToast.value = true
  setTimeout(() => { showToast.value = false }, 3000)
}

async function handleArchive() {
  try {
    const res = await fetch(`/api/v1/tasks/${props.task.id}/archive`, { method: 'POST' }).then(r => r.json())
    if (res.code === 0) {
      showToastMsg('Task archived')
      emit('refresh')
      emit('close')
    } else {
      showToastMsg('Failed: ' + res.msg)
    }
  } catch (err) {
    showToastMsg('Network error')
  }
}

async function handleUnarchive() {
  try {
    const res = await fetch(`/api/v1/tasks/${props.task.id}/unarchive`, { method: 'POST' }).then(r => r.json())
    if (res.code === 0) {
      showToastMsg('Task unarchived')
      emit('refresh')
      emit('close')
    } else {
      showToastMsg('Failed: ' + res.msg)
    }
  } catch (err) {
    showToastMsg('Network error')
  }
}

const formattedDeadline = computed(() => {
  if (!props.task || !props.task.deadline) return ''
  return new Date(props.task.deadline).toLocaleString()
})

const formattedCompletedAt = computed(() => {
  if (!props.task || !props.task.actual_completed_at) return ''
  return new Date(props.task.actual_completed_at).toLocaleString()
})

const parsedLinks = computed(() => {
  if (!props.task || !props.task.links) return []
  return props.task.links.split('\n').filter(l => l.trim())
})

const statusClass = computed(() => {
  if (!props.task) return ''
  if (props.task.status === 'done') {
    return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-emerald-400 bg-emerald-500/10 px-2 py-1 rounded border border-emerald-500/20'
  } else if (props.task.status === 'in-progress') {
    return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20'
  }
  return 'text-[10px] uppercase tracking-wider font-semibold ml-2 text-slate-400 bg-slate-500/10 px-2 py-1 rounded border border-slate-500/20'
})

function formatTime(isoString) {
  return new Date(isoString).toLocaleString()
}
</script>

<template>
  <div v-show="isOpen && task" class="fixed inset-0 z-50" @click.self="emit('close')">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity pointer-events-none"></div>
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-2xl bg-dark-card border border-dark-border text-left shadow-2xl transition-all sm:my-8 sm:w-full sm:max-w-2xl">
          <div class="px-6 py-5 border-b border-dark-border flex justify-between items-center bg-gradient-to-r from-dark-card to-dark-bg">
            <h3 class="text-xl font-bold text-white truncate pr-4 drop-shadow-sm">{{ task?.title }}</h3>
            <div class="flex gap-2">
              <button v-if="task?.archived_at" @click="handleUnarchive" class="text-xs bg-amber-500/10 hover:bg-amber-500/20 border border-amber-500/30 text-amber-400 px-3.5 py-1.5 rounded-xl transition-all shadow-sm">
                Unarchive
              </button>
              <button v-else-if="task?.status === 'done'" @click="handleArchive" class="text-xs bg-amber-500/10 hover:bg-amber-500/20 border border-amber-500/30 text-amber-400 px-3.5 py-1.5 rounded-xl transition-all shadow-sm">
                Archive
              </button>
              <button @click="emit('edit')" class="text-xs bg-dark-bg hover:bg-white/5 border border-white/10 text-slate-300 hover:text-white px-3.5 py-1.5 rounded-xl transition-all shadow-sm">
                Edit
              </button>
              <button @click="emit('delete')" class="text-xs bg-dark-bg hover:bg-red-500/10 border border-red-500/20 text-red-400 hover:text-red-300 px-3.5 py-1.5 rounded-xl transition-all shadow-sm">
                Delete
              </button>
              <button v-if="task?.status === 'todo'" @click="emit('start')" class="text-xs bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-white px-3.5 py-1.5 rounded-xl transition-all shadow-lg shadow-emerald-500/20 active:scale-95 border border-emerald-400/20 font-medium">
                Start Task
              </button>
              <button v-if="task?.status !== 'done'" @click="emit('progress')" class="text-xs bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-400 hover:to-indigo-500 text-white px-3.5 py-1.5 rounded-xl transition-all shadow-lg shadow-indigo-500/20 active:scale-95 border border-indigo-400/20 font-medium">
                Update Progress
              </button>
              <button @click="emit('close')" class="flex-shrink-0 text-slate-400 hover:text-white transition-colors p-1 ml-1 rounded-full hover:bg-white/5">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
              </button>
            </div>
          </div>
          <div class="px-6 py-5 space-y-5 max-h-[70vh] overflow-y-auto">
            <div class="flex items-center">
              <span class="text-[10px] uppercase tracking-wider font-semibold text-indigo-400 bg-indigo-500/10 px-2 py-1 rounded border border-indigo-500/20">
                {{ task?.category }}
              </span>
              <span :class="statusClass">{{ task?.status }}</span>
            </div>
            
            <div class="flex flex-col gap-2 mt-1 mb-2">
              <div v-if="task?.deadline" class="flex items-center gap-2 text-xs text-slate-400">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path></svg>
                <span>Deadline: <span class="text-slate-300 font-medium">{{ formattedDeadline }}</span></span>
              </div>
              <div v-if="task?.actual_completed_at" class="flex items-center gap-2 text-xs text-slate-400">
                <svg class="w-3.5 h-3.5 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <span>Completed: <span class="text-emerald-400 font-medium">{{ formattedCompletedAt }}</span></span>
              </div>
            </div>

            <div v-if="task?.description">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Description</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed">{{ task.description }}</div>
            </div>

            <div v-if="task?.targets">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Targets</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed whitespace-pre-wrap">{{ task.targets }}</div>
            </div>

            <div v-if="parsedLinks.length > 0">
              <h4 class="text-sm font-medium text-slate-400 mb-2">Links</h4>
              <div class="text-sm text-slate-300 bg-dark-bg p-4 rounded-xl border border-dark-border leading-relaxed flex flex-col gap-1">
                <a v-for="(link, i) in parsedLinks" :key="i" :href="link" target="_blank" class="text-indigo-400 hover:text-indigo-300 underline break-all">{{ link }}</a>
              </div>
            </div>

            <div>
              <h4 class="text-sm font-medium text-slate-400 mb-3 flex items-center gap-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                Worklogs
              </h4>
                <div v-if="!task?.logs || task.logs.length === 0" class="text-sm text-slate-500 italic px-4 py-3 bg-dark-bg/50 rounded-xl border border-dashed border-dark-border text-center">No worklogs yet. Start tracking your progress!</div>
                
                <div v-else class="relative pl-7 space-y-5 before:absolute before:inset-y-0 before:left-[11px] before:w-px before:bg-dark-border mt-3 mb-2">
                  <div v-for="(log, index) in task?.logs" :key="log.id" class="relative group">
                    <!-- Clean Timeline Dot -->
                    <div class="absolute -left-[30px] top-[7px] w-2 h-2 rounded-full bg-slate-400 group-hover:bg-indigo-400 group-hover:shadow-[0_0_8px_rgba(99,102,241,0.6)] transition-all z-10 ring-4 ring-dark-card"></div>
                    
                    <div class="flex flex-col gap-1 -mt-0.5">
                      <div class="flex items-center gap-3">
                        <span class="text-[11px] font-medium text-slate-500 group-hover:text-indigo-300 transition-colors">{{ formatTime(log.created_at) }}</span>
                        <div class="flex-grow border-t border-dashed border-dark-border/50"></div>
                        <button @click="emit('delete-worklog', log.id)" title="Delete Worklog" class="text-xs text-red-500/50 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-all p-1 hover:bg-red-500/10 rounded">
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                        </button>
                      </div>
                      <div class="text-slate-300 whitespace-pre-wrap text-sm leading-relaxed pr-6">{{ log.log_text }}</div>
                    </div>
                  </div>
                </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Toast -->
  <div class="fixed bottom-4 right-4 z-50 transform transition-all duration-300" :class="[showToast ? 'translate-y-0 opacity-100' : 'translate-y-20 opacity-0']">
    <div class="bg-amber-500 text-white px-4 py-3 rounded-lg shadow-lg flex items-center gap-3">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
      </svg>
      <span class="text-sm font-medium">{{ toastMsg }}</span>
    </div>
  </div>
</template>
