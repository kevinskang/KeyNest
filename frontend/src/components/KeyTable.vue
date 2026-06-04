<script lang="ts" setup>
import type { APIKeyDTO } from '../types'

defineProps<{ keys: APIKeyDTO[] }>()

const emit = defineEmits<{
  edit:   [key: APIKeyDTO]
  delete: [id: number]
  copy:   [value: string]
}>()

function expiryClass(status: number) {
  if (status === 3) return 'badge badge--expired'
  if (status === 2) return 'badge badge--soon'
  return ''
}

function expiryLabel(status: number, date: string) {
  if (!date) return '-'
  if (status === 3) return `${date} (만료됨)`
  if (status === 2) return `${date} (임박)`
  return date
}

function truncate(val: string, max = 40) {
  return val.length > max ? val.slice(0, max) + '…' : val
}
</script>

<template>
  <table class="key-table">
    <thead>
      <tr>
        <th class="col-no">#</th>
        <th class="col-name">Key Name</th>
        <th class="col-value">Key Value</th>
        <th class="col-url">URL</th>
        <th class="col-expiry">만료예정일</th>
        <th class="col-reg">등록일</th>
        <th class="col-memo">메모</th>
        <th class="col-action">액션</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(key, idx) in keys" :key="key.id" :class="{ 'row--expired': key.expiryStatus === 3, 'row--soon': key.expiryStatus === 2 }">
        <td class="col-no">{{ idx + 1 }}</td>
        <td class="col-name">
          <span class="key-name" :title="key.keyName">{{ key.keyName }}</span>
        </td>
        <td class="col-value">
          <span class="key-value" :title="key.keyValue">{{ truncate(key.keyValue) }}</span>
        </td>
        <td class="col-url">
          <a v-if="key.url" :href="key.url" target="_blank" class="url-link" :title="key.url">
            {{ truncate(key.url, 30) }}
          </a>
          <span v-else class="muted">-</span>
        </td>
        <td class="col-expiry">
          <span :class="expiryClass(key.expiryStatus)">
            {{ expiryLabel(key.expiryStatus, key.expiryDate) }}
          </span>
        </td>
        <td class="col-reg">{{ key.registeredDate || '-' }}</td>
        <td class="col-memo">
          <span :title="key.memo">{{ truncate(key.memo, 30) || '-' }}</span>
        </td>
        <td class="col-action">
          <div class="action-group">
            <button class="btn btn-ghost btn-sm" @click="emit('copy', key.keyValue)" title="클립보드 복사">
              복사
            </button>
            <button class="btn btn-ghost btn-sm" @click="emit('edit', key)" title="수정">
              수정
            </button>
            <button class="btn btn-danger btn-sm" @click="emit('delete', key.id)" title="삭제">
              삭제
            </button>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style scoped>
.key-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  table-layout: fixed;
}

thead th {
  background: var(--bg-card);
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 10px 12px;
  text-align: left;
  border-bottom: 2px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 1;
}

tbody tr {
  border-bottom: 1px solid var(--border);
  transition: background 0.1s;
}
tbody tr:hover    { background: rgba(74, 158, 237, 0.05); }
tbody tr.row--expired { background: rgba(231, 76, 60, 0.06); }
tbody tr.row--soon    { background: rgba(243, 156, 18, 0.06); }

td {
  padding: 10px 12px;
  vertical-align: middle;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Column widths */
.col-no     { width: 44px;  text-align: center; color: var(--text-muted); }
.col-name   { width: 16%; }
.col-value  { width: 22%; }
.col-url    { width: 14%; }
.col-expiry { width: 14%; }
.col-reg    { width: 10%; }
.col-memo   { width: 12%; }
.col-action { width: 160px; }

.key-name  { font-weight: 600; }
.key-value { font-family: monospace; font-size: 12px; color: #a8d8ff; }
.muted     { color: var(--text-muted); }

.url-link  { color: var(--primary); text-decoration: none; }
.url-link:hover { text-decoration: underline; }

.badge {
  display: inline-block;
  padding: 2px 7px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
}
.badge--expired { background: rgba(231,76,60,0.2); color: #e74c3c; }
.badge--soon    { background: rgba(243,156,18,0.2); color: #f39c12; }

.action-group { display: flex; gap: 6px; }
</style>
