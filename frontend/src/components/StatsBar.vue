<script setup>
import { ref, onMounted } from 'vue'

const summary = ref(null)

async function loadSummary() {
  try {
    const res = await fetch('/api/v1/stats/summary').then(r => r.json())
    if (res.code === 0) {
      summary.value = res.data
    }
  } catch (err) {
    console.error('Failed to load summary', err)
  }
}

onMounted(() => {
  loadSummary()
})
</script>

<template>
  <div v-if="summary" class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
    <div class="bg-dark-card rounded-2xl p-4 border border-dark-border">
      <p class="text-xs text-slate-400 mb-1">Total Tasks</p>
      <p class="text-2xl font-bold text-white">{{ summary.total_tasks }}</p>
    </div>
    
    <div class="bg-dark-card rounded-2xl p-4 border border-dark-border">
      <p class="text-xs text-slate-400 mb-1">Completed</p>
      <p class="text-2xl font-bold text-emerald-400">{{ summary.completed_tasks }}</p>
    </div>

    <div class="bg-dark-card rounded-2xl p-4 border border-dark-border">
      <p class="text-xs text-slate-400 mb-1">In Progress</p>
      <p class="text-2xl font-bold text-indigo-400">{{ summary.in_progress_tasks }}</p>
    </div>

    <div class="bg-dark-card rounded-2xl p-4 border border-dark-border">
      <p class="text-xs text-slate-400 mb-1">Completion Rate</p>
      <p class="text-2xl font-bold text-purple-400">{{ Math.round(summary.completion_rate * 100) }}%</p>
    </div>
  </div>
</template>
