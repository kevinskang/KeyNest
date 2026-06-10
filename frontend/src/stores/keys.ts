import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GetKeys, GetExpiringKeys } from '../../wailsjs/go/main/App'
import type { APIKeyDTO, KeyFilter } from '../types'

export const useKeyStore = defineStore('keys', () => {
  const keys = ref<APIKeyDTO[]>([])
  const expiringKeys = ref<APIKeyDTO[]>([])
  const loading = ref(false)
  const error = ref('')
  const filter = reactive<KeyFilter>({ keyName: '', dateFrom: '', dateTo: '' })

  async function loadKeys() {
    loading.value = true
    error.value = ''
    try {
      const result = await GetKeys({ ...filter })
      keys.value = result ?? []
    } catch (e) {
      error.value = e instanceof Error ? e.message : '키 목록을 불러오는 중 오류가 발생했습니다.'
      keys.value = []
    } finally {
      loading.value = false
    }
  }

  async function loadExpiringKeys() {
    try {
      const result = await GetExpiringKeys(30)
      expiringKeys.value = result ?? []
    } catch {
      expiringKeys.value = []
    }
  }

  function resetFilter() {
    filter.keyName = ''
    filter.dateFrom = ''
    filter.dateTo = ''
  }

  return { keys, expiringKeys, loading, error, filter, loadKeys, loadExpiringKeys, resetFilter }
})
