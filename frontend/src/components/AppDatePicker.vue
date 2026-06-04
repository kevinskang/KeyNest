<script lang="ts" setup>
import VueDatePicker from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  clearable?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// model-type="format" + format="yyyy-MM-dd" 설정 시
// VueDatePicker가 내부에서 Date ↔ YYYY-MM-DD 문자열 변환을 처리합니다.
// 빈 문자열을 null로 변환해야 달력이 빈 상태로 표시됩니다.
function handleUpdate(val: string | null) {
  emit('update:modelValue', val ?? '')
}
</script>

<template>
  <VueDatePicker
    :model-value="modelValue || null"
    @update:model-value="handleUpdate"
    model-type="format"
    format="yyyy-MM-dd"
    :enable-time-picker="false"
    :auto-apply="true"
    :clearable="clearable ?? true"
    :dark="true"
    locale="ko"
    :placeholder="placeholder ?? '날짜 선택'"
    input-class-name="dp-input-override"
  />
</template>

<style>
/* ── 입력 필드 ─────────────────────────────────────── */
.dp-input-override {
  background: var(--bg-input) !important;
  border: 1px solid var(--border) !important;
  color: var(--text) !important;
  border-radius: var(--radius) !important;
  font-size: 14px !important;
  font-family: inherit !important;
  height: 36px !important;
  cursor: pointer !important;
}
.dp-input-override:focus {
  border-color: var(--primary) !important;
  outline: none !important;
}

/* ── 달력 팝업 다크 테마 ────────────────────────────── */
.dp__theme_dark {
  --dp-background-color:   #243447;
  --dp-text-color:         #dde6f0;
  --dp-hover-color:        #2e4a65;
  --dp-hover-text-color:   #dde6f0;
  --dp-primary-color:      #4a9eed;
  --dp-primary-text-color: #ffffff;
  --dp-secondary-color:    #8aa8c0;
  --dp-border-color:       #2e4a65;
  --dp-menu-border-color:  #2e4a65;
  --dp-border-color-hover: #4a9eed;
  --dp-disabled-color:     #1a2c40;
  --dp-icon-color:         #8aa8c0;
  --dp-danger-color:       #e74c3c;
  --dp-highlight-color:    rgba(74, 158, 237, 0.18);
  --dp-range-between-dates-background-color: rgba(74, 158, 237, 0.10);
  --dp-range-between-border-color:           rgba(74, 158, 237, 0.10);
  --dp-font-family:        "Nunito", sans-serif;
  --dp-font-size:          13px;
  --dp-border-radius:      6px;
  --dp-cell-border-radius: 4px;
  --dp-cell-size:          34px;
  --dp-input-padding:      8px 12px;
  --dp-menu-min-width:     260px;
}
</style>
