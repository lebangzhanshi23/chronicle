<script setup>
import { ref, onMounted } from 'vue'

const summaryData = ref(null)
const error = ref(null)

async function loadDailySummary() {
  try {
    const res = await fetch('/api/v1/reports/daily-summary').then(r => r.json())
    if (res.code === 0) {
      summaryData.value = res.data
    } else {
      error.value = res.msg
    }
  } catch (err) {
    error.value = 'Failed to load daily summary'
  }
}

onMounted(() => {
  loadDailySummary()
})
</script>

<template>
  <div class="fixed inset-0 z-50 bg-dark-bg/95 flex items-center justify-center p-4 backdrop-blur-md overflow-y-auto" @click.self="$emit('close')">
    <div class="max-w-4xl w-full">
      <div class="flex items-center justify-between mb-8">
        <h1 class="text-3xl font-bold text-white tracking-tight flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-amber-400 to-amber-600 flex items-center justify-center text-white shadow-lg shadow-amber-500/20">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          </div>
          Daily Standup Summary
        </h1>
        <button @click="$emit('close')" class="p-2 rounded-full hover:bg-white/10 text-slate-400 hover:text-white transition-colors">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
        </button>
      </div>

      <div v-if="error" class="bg-red-500/10 border border-red-500/20 text-red-400 p-4 rounded-xl text-sm mb-6">
        {{ error }}
      </div>

      <div v-else-if="summaryData" class="space-y-6">
        <div class="bg-dark-card rounded-3xl p-8 border border-dark-border shadow-2xl relative overflow-hidden">
          <div class="absolute -right-20 -top-20 w-64 h-64 bg-amber-500/5 rounded-full blur-3xl"></div>
          
          <div class="flex flex-col gap-2 mb-8 relative z-10 border-b border-dark-border/50 pb-6">
            <p class="text-amber-400 font-semibold tracking-wider uppercase text-sm">Today's Progress</p>
            <h2 class="text-4xl font-bold text-white">{{ summaryData.date }}</h2>
            <p class="text-slate-400 mt-2 text-sm max-w-xl">This is an aggregated view of all tasks and worklogs you've touched today. Perfect for daily standups or wrapping up your day.</p>
          </div>

          <div v-if="!summaryData.activities || summaryData.activities.length === 0" class="flex flex-col items-center justify-center py-12 text-center relative z-10">
            <div class="w-16 h-16 rounded-full bg-dark-bg border border-dark-border flex items-center justify-center mb-4 text-slate-500">
               <svg class="w-8 h-8 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path></svg>
            </div>
            <h3 class="text-lg font-medium text-slate-300">No activity recorded today</h3>
            <p class="text-sm text-slate-500 mt-2">Start updating your tasks to see your daily standup report here.</p>
          </div>

          <div v-else class="space-y-8 relative z-10">
            <div v-for="activity in summaryData.activities" :key="activity.task_id" class="group">
              <div class="flex items-center gap-3 mb-4">
                <span class="text-[10px] uppercase tracking-wider font-semibold px-2.5 py-1 rounded bg-dark-bg border border-dark-border text-slate-300">{{ activity.category || 'Uncategorized' }}</span>
                <h3 class="text-lg font-medium text-white group-hover:text-amber-300 transition-colors">{{ activity.task_title }}</h3>
                <span v-if="activity.status === 'done'" class="ml-auto text-[10px] uppercase tracking-wider font-semibold text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">Done</span>
                <span v-else-if="activity.status === 'in-progress'" class="ml-auto text-[10px] uppercase tracking-wider font-semibold text-indigo-400 bg-indigo-500/10 px-2 py-0.5 rounded border border-indigo-500/20">In Progress</span>
              </div>
              
              <div class="pl-4 ml-2 border-l-2 border-dark-border space-y-3">
                <div v-for="(log, idx) in activity.today_logs" :key="idx" class="relative group/log">
                  <div class="absolute -left-[23px] top-2 w-2.5 h-2.5 rounded-full bg-dark-card border-2 border-amber-500/50 group-hover/log:border-amber-400 transition-colors"></div>
                  <div class="text-slate-300 text-sm leading-relaxed bg-dark-bg/50 px-4 py-3 rounded-xl border border-transparent hover:border-dark-border/80 transition-colors">
                    <span class="text-amber-400/80 font-mono text-xs mr-2">{{ log.split(' - ')[0] }}</span> 
                    <span class="text-slate-200">{{ log.substring(log.indexOf(' - ') + 3) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="flex justify-center p-12">
        <div class="w-8 h-8 rounded-full border-2 border-amber-500/20 border-t-amber-500 animate-spin"></div>
      </div>
    </div>
  </div>
</template>
