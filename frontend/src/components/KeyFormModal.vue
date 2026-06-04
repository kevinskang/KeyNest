<script lang="ts" setup>
import { ref, watch } from 'vue'
import type { APIKeyDTO, CreateKeyRequest, UpdateKeyRequest } from '../types'
import AppDatePicker from './AppDatePicker.vue'

const props = defineProps<{ editKey: APIKeyDTO | null }>()
const emit = defineEmits<{
  close: []
  save:  [req: CreateKeyRequest | UpdateKeyRequest]
}>()

const keyName        = ref('')
const keyValue       = ref('')
const url            = ref('')
const expiryDate     = ref('')
const registeredDate = ref('')
const memo           = ref('')
const errorMsg       = ref('')

// 수정 모드 진입 시 폼 초기화
watch(() => props.editKey, (key) => {
  if (key) {
    keyName.value        = key.keyName
    keyValue.value       = key.keyValue
    url.value            = key.url
    expiryDate.value     = key.expiryDate
    registeredDate.value = key.registeredDate
    memo.value           = key.memo
  } else {
    keyName.value        = ''
    keyValue.value       = ''
    url.value            = ''
    expiryDate.value     = ''
    registeredDate.value = ''
    memo.value           = ''
  }
  errorMsg.value = ''
}, { immediate: true })

const isEdit = () => props.editKey !== null

function handleSubmit() {
  errorMsg.value = ''
  if (!keyName.value.trim()) {
    errorMsg.value = 'Key Name은 필수입니다.'
    return
  }
  if (!keyValue.value.trim()) {
    errorMsg.value = 'Key Value는 필수입니다.'
    return
  }

  const base = {
    keyName:        keyName.value.trim(),
    keyValue:       keyValue.value.trim(),
    url:            url.value.trim(),
    expiryDate:     expiryDate.value,
    registeredDate: registeredDate.value,
    memo:           memo.value.trim(),
  }

  if (isEdit()) {
    emit('save', { id: props.editKey!.id, ...base } as UpdateKeyRequest)
  } else {
    emit('save', base as CreateKeyRequest)
  }
}
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <h2 class="modal-title">{{ isEdit() ? '키 수정' : '키 등록' }}</h2>

      <div class="form-grid">
        <div class="form-group span-2">
          <label class="form-label">Key Name <span class="required">*</span></label>
          <input v-model="keyName" type="text" placeholder="GitHub Personal Token" />
        </div>
        <div class="form-group span-2">
          <label class="form-label">Key Value <span class="required">*</span></label>
          <textarea v-model="keyValue" rows="3" placeholder="ghp_xxxxxxxxxxxxxxxxxxxx" class="value-area" />
        </div>
        <div class="form-group span-2">
          <label class="form-label">URL</label>
          <input v-model="url" type="url" placeholder="https://github.com/settings/tokens" />
        </div>
        <div class="form-group">
          <label class="form-label">만료예정일</label>
          <AppDatePicker v-model="expiryDate" placeholder="만료 날짜 선택" />
        </div>
        <div class="form-group">
          <label class="form-label">등록일 (키 발급일)</label>
          <AppDatePicker v-model="registeredDate" placeholder="발급 날짜 선택" />
        </div>
        <div class="form-group span-2">
          <label class="form-label">메모</label>
          <textarea v-model="memo" rows="2" placeholder="용도, 권한 범위 등..." />
        </div>
      </div>

      <p v-if="errorMsg" class="form-error" style="margin-top: 12px;">{{ errorMsg }}</p>

      <div class="modal-actions">
        <button class="btn btn-ghost" @click="emit('close')">취소</button>
        <button class="btn btn-primary" @click="handleSubmit">{{ isEdit() ? '수정' : '등록' }}</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}
.span-2 { grid-column: span 2; }
.required { color: var(--danger); }
.value-area {
  font-family: monospace;
  font-size: 12px;
  resize: vertical;
}
</style>
