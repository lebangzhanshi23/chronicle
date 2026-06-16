<script setup>
import { ref, watch, computed } from 'vue'

const emit = defineEmits(['close', 'refresh'])

const props = defineProps({
  isOpen: Boolean,
  task: Object
})

const isEditMode = computed(() => !!props.task)
const modalTitle = computed(() => isEditMode.value ? 'Edit Task' : 'Create New Task')
const submitButtonText = computed(() => isEditMode.value ? 'Save Changes' : 'Create Task')

const toastMsg = ref('')
const showToast = ref(false)

function showToastMsg(msg) {
  toastMsg.value = msg
  showToast.value = true
  setTimeout(() => { showToast.value = false }, 3000)
}

async function handleSubmit() {
  const payload = { ...formData.value }
  if (payload.deadline) {
    payload.deadline = new Date(payload.deadline).toISOString()
  } else {
    payload.deadline = null
  }

  try {
    let res
    if (isEditMode.value) {
      res = await fetch(`/api/v1/tasks/${payload.id}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      }).then(r => r.json())
    } else {
      res = await fetch('/api/v1/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      }).then(r => r.json())
    }

    if (res.code === 0) {
      showToastMsg(isEditMode.value ? 'Task updated' : 'Task created')
      emit('refresh')
      emit('close')
    } else {
      showToastMsg('Failed: ' + res.msg)
    }
  } catch (err) {
    showToastMsg('Network error')
  }
}

const formData = ref({
  id: '',
  title: '',
  category: '',
  deadline: '',
  description: '',
  targets: '',
  links: ''
})

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    if (props.task) {
      // Edit mode: populate form
      formData.value = {
        id: props.task.id,
        title: props.task.title,
        category: props.task.category,
        description: props.task.description || '',
        targets: props.task.targets || '',
        links: props.task.links || '',
        deadline: ''
      }
      if (props.task.deadline) {
        const d = new Date(props.task.deadline)
        const year = d.getFullYear()
        const month = String(d.getMonth() + 1).padStart(2, '0')
        const day = String(d.getDate()).padStart(2, '0')
        const hours = String(d.getHours()).padStart(2, '0')
        const minutes = String(d.getMinutes()).padStart(2, '0')
        formData.value.deadline = `${year}-${month}-${day}T${hours}:${minutes}`
      }
    } else {
      // Create mode: reset form and set default deadline
      formData.value = { id: '', title: '', category: '', deadline: '', description: '', targets: '', links: '' }
      const d = new Date()
      d.setDate(d.getDate() + 7)
      d.setHours(20, 30, 0, 0)
      const year = d.getFullYear()
      const month = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      const hours = String(d.getHours()).padStart(2, '0')
      const minutes = String(d.getMinutes()).padStart(2, '0')
      formData.value.deadline = `${year}-${month}-${day}T${hours}:${minutes}`
    }
  }
})
</script>

<template>
  <div v-show="isOpen" class="fixed inset-0 z-50" @click.self="emit('close')">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity pointer-events-none"></div>
    <div class="fixed inset-0 z-10 overflow-y-auto">
      <div class="flex min-h-full items-end justify-center p-4 text-center sm:items-center sm:p-0">
        <div class="relative transform overflow-hidden rounded-2xl bg-dark-card border border-dark-border text-left shadow-2xl transition-all sm:my-8 sm:w-full sm:max-w-lg">
          <div class="px-6 py-5 border-b border-dark-border flex justify-between items-center bg-gradient-to-r from-dark-card to-dark-bg">
            <h3 class="text-xl font-bold text-white drop-shadow-sm">{{ modalTitle }}</h3>
            <button @click="emit('close')" class="text-slate-400 hover:text-white transition-colors p-1 ml-1 rounded-full hover:bg-white/5">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path></svg>
            </button>
          </div>
          <form @submit.prevent="handleSubmit" class="px-6 py-5 space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Title</label>
              <input v-model="formData.title" type="text" required class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500" placeholder="e.g. Develop login API">
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Category</label>
              <input v-model="formData.category" type="text" required class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500" placeholder="e.g. Backend">
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Deadline <span class="text-slate-500 text-xs font-normal">(Optional)</span></label>
              <input v-model="formData.deadline" type="datetime-local" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500">
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Description <span class="text-slate-500 text-xs font-normal">(Optional)</span></label>
              <textarea v-model="formData.description" rows="3" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500 resize-none" placeholder="Provide additional background context..."></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-300 mb-1">Targets <span class="text-slate-500 text-xs font-normal">(Optional)</span></label>
              <textarea v-model="formData.targets" rows="2" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500 resize-none" placeholder="Acceptance criteria..."></textarea>
            </div>
            <div v-show="isEditMode">
              <label class="block text-sm font-medium text-slate-300 mb-1">Links <span class="text-slate-500 text-xs font-normal">(Optional, one URL per line)</span></label>
              <textarea v-model="formData.links" rows="2" class="w-full bg-dark-bg border border-dark-border rounded-lg px-4 py-2.5 text-slate-200 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition-all placeholder:text-slate-500 resize-none" placeholder="https://..."></textarea>
            </div>
            <div class="mt-6 flex justify-end gap-3 pt-2">
              <button type="button" @click="emit('close')" class="px-5 py-2.5 text-sm font-medium text-slate-400 hover:text-white transition-colors rounded-xl hover:bg-white/5">Cancel</button>
              <button type="submit" class="bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-400 hover:to-indigo-500 text-white px-6 py-2.5 rounded-xl text-sm font-semibold transition-all shadow-lg shadow-indigo-500/20 active:scale-95 border border-indigo-400/20">{{ submitButtonText }}</button>
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
