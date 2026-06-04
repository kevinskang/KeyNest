<script lang="ts" setup>
import { useToast } from './composables/useToast'
import { computed } from 'vue'

const { message, type, visible } = useToast()

const toastClass = computed(() => ({
  toast: true,
  'toast--success': type.value === 'success',
  'toast--error': type.value === 'error',
  'toast--info': type.value === 'info',
}))
</script>

<template>
  <RouterView />
  <Transition name="toast-fade">
    <div v-if="visible" :class="toastClass">{{ message }}</div>
  </Transition>
</template>

<style>
.toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 12px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}
.toast--success { background: #27ae60; color: #fff; }
.toast--error   { background: #e74c3c; color: #fff; }
.toast--info    { background: #4a9eed; color: #fff; }

.toast-fade-enter-active,
.toast-fade-leave-active { transition: opacity 0.3s, transform 0.3s; }
.toast-fade-enter-from,
.toast-fade-leave-to     { opacity: 0; transform: translateY(8px); }
</style>
