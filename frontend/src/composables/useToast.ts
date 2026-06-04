import { ref } from 'vue'

type ToastType = 'success' | 'error' | 'info'

const message = ref('')
const type = ref<ToastType>('success')
const visible = ref(false)
let timer: ReturnType<typeof setTimeout>

export function useToast() {
  function show(msg: string, t: ToastType = 'success', duration = 2500) {
    message.value = msg
    type.value = t
    visible.value = true
    clearTimeout(timer)
    timer = setTimeout(() => { visible.value = false }, duration)
  }
  return { message, type, visible, show }
}
