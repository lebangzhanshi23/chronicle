<script setup>
defineProps({
  isOpen: Boolean,
  title: { type: String, default: 'Confirm Action' },
  message: { type: String, default: 'Are you sure?' },
  confirmText: { type: String, default: 'Confirm' },
  cancelText: { type: String, default: 'Cancel' },
  confirmDanger: { type: Boolean, default: false }
})

const emit = defineEmits(['confirm', 'cancel'])
</script>

<template>
  <div v-show="isOpen" class="fixed inset-0 z-50">
    <div class="fixed inset-0 bg-black/70 backdrop-blur-sm transition-opacity" @click="emit('cancel')"></div>
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-center justify-center p-4">
        <div class="relative transform overflow-hidden rounded-2xl text-left shadow-2xl transition-all sm:w-full sm:max-w-md glass-card border border-dark-border">>
          <div class="px-6 py-5">
            <div class="flex items-start gap-4">
              <!-- Warning Icon -->
              <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center">
                <svg class="w-5 h-5 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                </svg>
              </div>
              <div class="flex-1">
                <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
                <p class="mt-2 text-sm text-slate-400">{{ message }}</p>
              </div>
            </div>
          </div>
          <div class="px-6 py-4 bg-dark-bg/50 flex justify-end gap-3 border-t border-dark-border">
            <button 
              @click="emit('cancel')" 
              class="px-4 py-2 text-sm font-medium text-slate-300 hover:text-white hover:bg-white/5 rounded-xl transition-all border border-white/10"
            >
              {{ cancelText }}
            </button>
            <button 
              @click="emit('confirm')" 
              :class="confirmDanger 
                ? 'bg-red-500 hover:bg-red-600 text-white shadow-lg shadow-red-500/20' 
                : 'bg-indigo-500 hover:bg-indigo-600 text-white shadow-lg shadow-indigo-500/20'"
              class="px-4 py-2 text-sm font-medium rounded-xl transition-all"
            >
              {{ confirmText }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
