<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['click', 'start'])

const isDone = computed(() => props.task.status === 'done')
const isInProgress = computed(() => props.task.status === 'in-progress')
const isTodo = computed(() => props.task.status === 'todo')

const cardClasses = computed(() => {
  let base = 'glass-card rounded-2xl p-5 hover:-translate-y-1 hover:border-indigo-500/30 transition-all duration-300 cursor-pointer group hover:shadow-2xl hover:shadow-indigo-500/20 task-enter'
  if (isDone.value) {
    base += ' opacity-50 grayscale-[0.7] hover:grayscale-0'
  }
  return base
})

const titleClasses = computed(() => {
  let base = 'font-semibold text-slate-200 group-hover:text-white transition-colors leading-snug break-words'
  if (isDone.value) {
    base += ' line-through text-slate-400'
  }
  return base
})

const isStarting = ref(false)

async function handleStart() {
  if (isStarting.value) return
  isStarting.value = true
  try {
    await fetch(`/api/v1/tasks/${props.task.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status: 'in-progress' })
    })
    emit('start')
  } catch (err) {
    console.error(err)
  } finally {
    isStarting.value = false
  }
}
</script>

<template>
  <div :class="cardClasses" @click="emit('click', task)">
    <div class="flex justify-between items-start mb-2">
      <h3 :class="titleClasses">{{ task.title }}</h3>
    </div>
    
    <div class="flex items-center gap-2 mt-4 flex-wrap">
      <span class="text-[10px] uppercase tracking-widest font-bold text-indigo-300 bg-indigo-500/10 px-2.5 py-1 rounded-full border border-indigo-500/20 shadow-[0_0_8px_rgba(99,102,241,0.1)]">
        {{ task.category }}
      </span>

      <!-- Todo Badge / Action -->
      <button v-if="isTodo" @click.stop="handleStart" title="Start Task" class="ml-auto flex items-center justify-center w-7 h-7 rounded-full bg-emerald-500/10 hover:bg-emerald-500 text-emerald-500 hover:text-white transition-all duration-300 border border-emerald-500/20 opacity-0 group-hover:opacity-100 flex-shrink-0 shadow-lg shadow-emerald-500/0 hover:shadow-emerald-500/30 disabled:opacity-50" :disabled="isStarting">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
      </button>

      <!-- In Progress Badge -->
      <span v-if="isInProgress" class="flex h-2 w-2 relative ml-auto flex-shrink-0">
        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-indigo-400 opacity-75"></span>
        <span class="relative inline-flex rounded-full h-2 w-2 bg-indigo-500"></span>
      </span>

      <!-- Done Badge -->
      <span v-else-if="isDone" class="ml-auto flex items-center text-emerald-500 flex-shrink-0">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
        </svg>
      </span>
    </div>
  </div>
</template>
