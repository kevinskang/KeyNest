import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GetKeys, GetExpiringKeys } from '../../wailsjs/go/main/App'
import type { APIKeyDTO, KeyFilter } from '../types'

export const useKeyStore = defineStore('keys', () => {
  const keys = ref<APIKeyDTO[]>([])
  const expiringKeys = ref<APIKeyDTO[]>([])
  const loading = ref(false)
  const filter = reactive<KeyFilter>({ keyName: '', dateFrom: '', dateTo: '' })

  async function loadKeys() {
    loading.value = true
    try {
      const result = await GetKeys({ ...filter })
      keys.value = result ?? []
    } finally {
      loading.value = false
    }
  }

  async function loadExpiringKeys() {
    const result = await GetExpiringKeys(30)
    expiringKeys.value = result ?? []
  }

  function resetFilter() {
    filter.keyName = ''
    filter.dateFrom = ''
    filter.dateTo = ''
  }

  return { keys, expiringKeys, loading, filter, loadKeys, loadExpiringKeys, resetFilter }
})
