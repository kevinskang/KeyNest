import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const isLoggedIn = ref(false)

  function setLoggedIn(value: boolean) {
    isLoggedIn.value = value
  }

  return { isLoggedIn, setLoggedIn }
})
