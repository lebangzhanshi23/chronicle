<script setup>
import TaskCard from './TaskCard.vue'

defineProps({
  todos: { type: Array, required: true },
  inProgress: { type: Array, required: true },
  done: { type: Array, required: true }
})

const emit = defineEmits(['task-click', 'start', 'view-stats', 'view-summary'])
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 items-start">
    
    <!-- Todo Column -->
    <div class="flex flex-col gap-4 bg-dark-card/30 rounded-3xl p-4 border border-white/[0.02]">
      <div class="flex items-center gap-2 mb-2 px-1">
        <div class="w-2.5 h-2.5 rounded-full bg-slate-500 shadow-[0_0_8px_rgba(100,116,139,0.5)]"></div>
        <h2 class="font-semibold text-slate-200">To Do</h2>
        <span class="text-xs bg-dark-bg border border-dark-border px-2.5 py-0.5 rounded-full text-slate-400 font-medium">
          {{ todos.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px]">
        <TaskCard 
          v-for="task in todos" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
          @start="emit('start', $event)"
        />
        <div v-if="todos.length === 0" class="text-sm text-slate-500 italic text-center py-4">No tasks to do</div>
      </div>
    </div>

    <!-- In Progress Column -->
    <div class="flex flex-col gap-4 bg-indigo-950/10 rounded-3xl p-4 border border-indigo-500/5 ring-1 ring-inset ring-indigo-500/5">
      <div class="flex items-center gap-2 mb-2 px-1">
        <div class="w-2.5 h-2.5 rounded-full bg-indigo-500 animate-pulse shadow-[0_0_8px_rgba(99,102,241,0.8)]"></div>
        <h2 class="font-semibold text-indigo-100">In Progress</h2>
        <span class="text-xs bg-indigo-500/10 border border-indigo-500/20 px-2.5 py-0.5 rounded-full text-indigo-300 font-medium">
          {{ inProgress.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px]">
        <TaskCard 
          v-for="task in inProgress" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
        />
        <div v-if="inProgress.length === 0" class="text-sm text-slate-500 italic text-center py-4">No tasks in progress</div>
      </div>
    </div>

    <!-- Done Column -->
    <div class="flex flex-col gap-4 bg-emerald-950/5 rounded-3xl p-4 border border-emerald-500/5">
      <div class="flex items-center gap-2 mb-2 px-1">
        <div class="w-2.5 h-2.5 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></div>
        <h2 class="font-semibold text-emerald-100">History / Done</h2>
        <span class="text-xs bg-emerald-500/10 border border-emerald-500/20 px-2.5 py-0.5 rounded-full text-emerald-400 font-medium">
          {{ done.length }}
        </span>
      </div>
      <div class="flex flex-col gap-3 min-h-[200px] max-h-[70vh] overflow-y-auto pr-2 custom-scroll">
        <TaskCard 
          v-for="task in done" 
          :key="task.id" 
          :task="task" 
          @click="emit('task-click', $event)" 
        />
        <div v-if="done.length === 0" class="text-sm text-slate-500 italic text-center py-4">No completed tasks yet</div>
      </div>
    </div>

    <!-- Quick Actions Column -->
    <div class="flex flex-col gap-4 bg-dark-card/30 rounded-3xl p-4 border border-white/[0.02]">
      <div class="flex items-center gap-2 mb-2 px-1">
        <div class="w-2.5 h-2.5 rounded-full bg-amber-500 shadow-[0_0_8px_rgba(245,158,11,0.5)]"></div>
        <h2 class="font-semibold text-amber-100">Quick Actions</h2>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-indigo-500/30 transition-colors group cursor-pointer" onclick="window.open('/api/v1/exports/daily-markdown', '_blank')">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2 bg-indigo-500/10 rounded-lg text-indigo-400 group-hover:bg-indigo-500 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
          </div>
          <h3 class="font-medium text-slate-200">Export Daily Markdown</h3>
        </div>
        <p class="text-xs text-slate-400">Generate Obsidian compatible markdown format of completed tasks for today.</p>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-emerald-500/30 transition-colors group cursor-pointer" @click="emit('view-stats')">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2 bg-emerald-500/10 rounded-lg text-emerald-400 group-hover:bg-emerald-500 group-hover:text-white transition-colors shadow-lg shadow-emerald-500/0 group-hover:shadow-emerald-500/40">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path></svg>
          </div>
          <h3 class="font-medium text-slate-200">System Statistics</h3>
        </div>
        <p class="text-xs text-slate-400 leading-relaxed">View an interactive dashboard of your productivity metrics and task categories.</p>
      </div>

      <div class="glass-card rounded-xl p-5 border border-dark-border hover:border-amber-500/30 transition-colors group cursor-pointer" @click="emit('view-summary')">
        <div class="flex items-center gap-3 mb-2">
            <div class="p-2 bg-amber-500/10 rounded-lg text-amber-500 group-hover:bg-amber-500 group-hover:text-white transition-colors shadow-lg shadow-amber-500/0 group-hover:shadow-amber-500/40">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
            </div>
            <h3 class="font-medium text-slate-200">Daily Standup</h3>
        </div>
        <p class="text-xs text-slate-400">View a beautiful timeline-based aggregated summary of everything you've done today.</p>
      </div>

    </div>
  </div>
</template>
