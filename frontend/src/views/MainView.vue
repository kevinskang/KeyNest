<script lang="ts" setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Logout, CreateKey, UpdateKey, DeleteKey } from '../../wailsjs/go/main/App'
import { useAuthStore } from '../stores/auth'
import { useKeyStore } from '../stores/keys'
import { useToast } from '../composables/useToast'
import KeyTable from '../components/KeyTable.vue'
import KeyFormModal from '../components/KeyFormModal.vue'
import ExpiryModal from '../components/ExpiryModal.vue'
import VueDatePicker from '@vuepic/vue-datepicker'
import type { APIKeyDTO, CreateKeyRequest, UpdateKeyRequest } from '../types'

const router = useRouter()
const auth = useAuthStore()
const keyStore = useKeyStore()
const { show } = useToast()

const showKeyForm = ref(false)
const editingKey = ref<APIKeyDTO | null>(null)
const showExpiry = ref(false)
const appVersion = 'v1.0.0'
const appAuthor = 'KANGPOLE'

// model-type="format" 사용 시 VueDatePicker가 YYYY-MM-DD 문자열을 직접 다룹니다.
const rangeModelValue = computed<string[] | null>(() => {
  if (!keyStore.filter.dateFrom && !keyStore.filter.dateTo) return null
  return [keyStore.filter.dateFrom, keyStore.filter.dateTo]
})

function handleRangeUpdate(val: string[] | null) {
  keyStore.filter.dateFrom = Array.isArray(val) && val[0] ? val[0] : ''
  keyStore.filter.dateTo   = Array.isArray(val) && val[1] ? val[1] : ''
}

let debounceTimer: ReturnType<typeof setTimeout>
let suppressWatch = false // 초기화 시 이중 조회 방지

// 자동 검색: 필터 변경 시 300ms debounce
watch(() => ({ ...keyStore.filter }), () => {
  if (suppressWatch) return
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => keyStore.loadKeys(), 300)
}, { deep: true })

onMounted(async () => {
  await keyStore.loadKeys()
  if (keyStore.expiringKeys.length > 0) {
    showExpiry.value = true
  }
})

// ── 검색 / 초기화 ─────────────────────────
async function handleSearch() {
  clearTimeout(debounceTimer)
  await keyStore.loadKeys()
}

async function handleReset() {
  clearTimeout(debounceTimer)
  suppressWatch = true
  keyStore.resetFilter()
  suppressWatch = false
  await keyStore.loadKeys()
}

// ── Auth ──────────────────────────────────
async function handleLogout() {
  await Logout()
  auth.setLoggedIn(false)
  keyStore.resetFilter()
  router.push('/login')
}

// ── CRUD ─────────────────────────────────
function openCreate() {
  editingKey.value = null
  showKeyForm.value = true
}

function openEdit(key: APIKeyDTO) {
  editingKey.value = key
  showKeyForm.value = true
}

async function handleSave(req: CreateKeyRequest | UpdateKeyRequest) {
  try {
    if ('id' in req && req.id) {
      await UpdateKey(req as UpdateKeyRequest)
      show('키가 수정되었습니다.', 'success')
    } else {
      await CreateKey(req as CreateKeyRequest)
      show('키가 등록되었습니다.', 'success')
    }
    showKeyForm.value = false
    await keyStore.loadKeys()
  } catch (e: unknown) {
    show(e instanceof Error ? e.message : '저장 중 오류가 발생했습니다.', 'error')
  }
}

async function handleDelete(id: number) {
  if (!window.confirm('이 키를 삭제하시겠습니까? 삭제된 데이터는 복구할 수 없습니다.')) return
  try {
    await DeleteKey(id)
    show('키가 삭제되었습니다.', 'success')
    await keyStore.loadKeys()
  } catch (e: unknown) {
    show(e instanceof Error ? e.message : '삭제 중 오류가 발생했습니다.', 'error')
  }
}

async function handleCopy(value: string) {
  try {
    await navigator.clipboard.writeText(value)
    show('클립보드에 복사되었습니다.', 'success')
  } catch {
    show('복사에 실패했습니다.', 'error')
  }
}
</script>

<template>
  <div class="main-layout">
    <!-- ── Header ── -->
    <header class="header">
      <div class="header-left">
        <span class="header-icon">🔑</span>
        <div class="header-texts">
          <span class="header-title">KeyNest</span>
          <span class="header-meta">{{ appVersion }} · {{ appAuthor }}</span>
        </div>
      </div>
      <div class="header-right">
        <span class="key-count">총 {{ keyStore.keys.length }}개</span>
        <button class="btn btn-ghost btn-sm" @click="handleLogout">로그아웃</button>
      </div>
    </header>

    <!-- ── Search bar ── -->
    <div class="toolbar">
      <div class="search-row">
        <input
          v-model="keyStore.filter.keyName"
          type="text"
          placeholder="Key Name 검색..."
          class="search-input"
        />
        <div class="date-range">
          <VueDatePicker
            :model-value="rangeModelValue"
            @update:model-value="handleRangeUpdate"
            range
            model-type="format"
            format="yyyy-MM-dd"
            :enable-time-picker="false"
            :auto-apply="true"
            :clearable="true"
            :dark="true"
            locale="ko"
            placeholder="등록일 범위 선택"
            input-class-name="dp-input-override dp-range-input"
          />
        </div>
      </div>
      <div class="toolbar-actions">
        <button class="btn btn-ghost" @click="handleReset" :disabled="keyStore.loading">
          초기화
        </button>
        <button class="btn btn-ghost" @click="handleSearch" :disabled="keyStore.loading">
          조회
        </button>
        <button class="btn btn-primary" @click="openCreate">+ 키 등록</button>
      </div>
    </div>

    <!-- ── Table ── -->
    <div class="table-wrapper">
      <div v-if="keyStore.loading" class="state-msg">불러오는 중...</div>
      <div v-else-if="keyStore.keys.length === 0" class="state-msg">등록된 키가 없습니다.</div>
      <KeyTable
        v-else
        :keys="keyStore.keys"
        @edit="openEdit"
        @delete="handleDelete"
        @copy="handleCopy"
      />
    </div>
  </div>

  <!-- ── Modals ── -->
  <KeyFormModal
    v-if="showKeyForm"
    :edit-key="editingKey"
    @close="showKeyForm = false"
    @save="handleSave"
  />

  <ExpiryModal
    v-if="showExpiry"
    :keys="keyStore.expiringKeys"
    @close="showExpiry = false"
  />
</template>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

/* ── Header ── */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  min-height: 64px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}
.header-left  { display: flex; align-items: center; gap: 10px; }
.header-icon  { font-size: 20px; }
.header-texts { display: flex; flex-direction: column; gap: 2px; }
.header-title { font-size: 18px; font-weight: 800; letter-spacing: -0.5px; }
.header-meta  { font-size: 12px; color: var(--text-muted); }
.header-right { display: flex; align-items: center; gap: 12px; }
.key-count    { font-size: 12px; color: var(--text-muted); }

/* ── Toolbar ── */
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
  flex-wrap: wrap;
}
.search-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 300px;
}
.search-input { max-width: 260px; }
.date-range { display: flex; align-items: center; }
.toolbar-actions { display: flex; gap: 8px; }
.toolbar-actions .btn { min-width: 72px; justify-content: center; }

/* ── Table area ── */
.table-wrapper {
  flex: 1;
  overflow: auto;
  padding: 16px;
}
.state-msg {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: var(--text-muted);
  font-size: 15px;
}
</style>
