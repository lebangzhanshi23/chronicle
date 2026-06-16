<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  isOpen: Boolean,
  task: Object
})

const emit = defineEmits(['close', 'refresh'])

const formData = ref({
  logText: '',
  deadline: '',
  action: 'update'
})

const toastMsg = ref('')
const showToast = ref(false)

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    formData.value.logText = ''
    formData.value.action = 'update'
    formData.value.deadline = ''
    if (props.task && props.task.deadline) {
      const d = new Date(props.task.deadline)
      const year = d.getFullYear()
      const month = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      const hours = String(d.getHours()).padStart(2, '0')
      const minutes = String(d.getMinutes()).padStart(2, '0')
      formData.value.deadline = `${year}-${month}-${day}T${hours}:${minutes}`
    }
  }
})

function showToastMsg(msg) {
  toastMsg.value = msg
  showToast.value = true
  setTimeout(() => { showToast.value = false }, 3000)
}

async function handleSubmit() {
  const payload = {
    log_text: formData.value.logText,
    mark_as_done: formData.value.action === 'done',
    new_status: formData.value.action === 'done' ? 'done' : 'in-progress',
    deadline: formData.value.deadline ? new Date(formData.value.deadline).toISOString() : null
  }

  try {
    const res = await fetch(`/api/v1/tasks/${props.task.id}/progress`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(r => r.json())

    if (res.code === 0) {
      showToastMsg(formData.value.action === 'done' ? 'Task done!' : 'Progress logged')
      emit('refresh')
      emit('close')
    } else {
      showToastMsg('Failed: ' + res.msg)
    }
  } catch (err) {
    showToastMsg('Network error')
  }
}
</script>

<template>
  <div v-show="isOpen" class="fixed inset-0 z-50" @click.self="emit('close')">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity pointer-events-none"></div>
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-2xl bg-dark-card border border-dark-border text-left shadow-2xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
          <div class="px-6 py-5 border-b border-dark-border flex justify-between items-center bg-gradient-to-r from-dark-card to-dark-bg">
            <h3 class="text-xl font-bold text-white truncate pr-4 drop-shadow-sm">Update Task: {{ task?.title }}</h3>
            <button @click="emit('close')" class="flex-shrink-0 text-slate-400 hover:text-white transition-colors p-1 ml-1 rounded-full hover:bg-white/5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
          </div>
          
          <form @submit.prevent="handleSubmit" class="px-6 py-5 space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Work Log <span class="text-red-400">*</span></label>
              <textarea v-model="formData.logText" required rows="3" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500 resize-none" placeholder="What did you just accomplish? This text will be saved to the database."></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Update Deadline <span class="text-slate-500 text-xs font-normal">(Optional)</span></label>
              <input v-model="formData.deadline" type="datetime-local" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500">
            </div>

            <div class="bg-dark-bg rounded-lg p-4 border border-dark-border">
              <label class="font-medium text-slate-200 mb-3 block text-sm">Action</label>
              <div class="flex flex-col gap-3">
                <label class="flex items-center gap-3 cursor-pointer group">
                  <div class="relative flex items-center justify-center w-5 h-5">
                    <input type="radio" v-model="formData.action" value="update" class="peer appearance-none w-5 h-5 rounded-full border-2 border-slate-500 checked:border-indigo-500 transition-colors cursor-pointer">
                    <div class="absolute w-2.5 h-2.5 rounded-full bg-indigo-500 scale-0 peer-checked:scale-100 transition-transform"></div>
                  </div>
                  <span class="text-sm text-slate-300 group-hover:text-slate-200 transition-colors">Just update log (stay In Progress)</span>
                </label>
                <label class="flex items-center gap-3 cursor-pointer group">
                  <div class="relative flex items-center justify-center w-5 h-5">
                    <input type="radio" v-model="formData.action" value="done" class="peer appearance-none w-5 h-5 rounded-full border-2 border-emerald-500 checked:border-emerald-500 transition-colors cursor-pointer">
                    <div class="absolute w-2.5 h-2.5 rounded-full bg-emerald-500 scale-0 peer-checked:scale-100 transition-transform"></div>
                  </div>
                  <span class="text-sm text-emerald-400 group-hover:text-emerald-300 transition-colors font-medium">Mark as Done!</span>
                </label>
              </div>
            </div>

            <div class="mt-6 flex justify-end gap-3 pt-2">
              <button type="button" @click="emit('close')" class="px-5 py-2.5 text-sm font-medium text-slate-400 hover:text-white transition-colors rounded-xl hover:bg-white/5">Cancel</button>
              <button type="submit" class="bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-400 hover:to-indigo-500 text-white px-6 py-2.5 rounded-xl text-sm font-semibold transition-all shadow-lg shadow-indigo-500/20 active:scale-95 border border-indigo-400/20">Save Log</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>

  <!-- Toast -->
  <div class="fixed bottom-4 right-4 z-50 transform transition-all duration-300" :class="[showToast ? 'translate-y-0 opacity-100' : 'translate-y-20 opacity-0']">
    <div class="bg-indigo-500 text-white px-4 py-3 rounded-lg shadow-lg flex items-center gap-3">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
      </svg>
      <span class="text-sm font-medium">{{ toastMsg }}</span>
    </div>
  </div>
</template>
