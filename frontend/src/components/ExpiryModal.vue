<script lang="ts" setup>
import type { APIKeyDTO } from '../types'

defineProps<{ keys: APIKeyDTO[] }>()
const emit = defineEmits<{ close: [] }>()

function rowClass(status: number) {
  if (status === 3) return 'row--expired'
  if (status === 2) return 'row--soon'
  return ''
}
function statusLabel(status: number) {
  if (status === 3) return '만료됨'
  if (status === 2) return '임박'
  return ''
}
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <h2 class="modal-title">⚠️ 만료 임박 키 알림</h2>
      <p class="desc">아래 키들이 30일 이내에 만료되거나 이미 만료되었습니다. 갱신 여부를 확인해주세요.</p>

      <table class="expiry-table">
        <thead>
          <tr>
            <th>Key Name</th>
            <th>만료예정일</th>
            <th>상태</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="key in keys" :key="key.id" :class="rowClass(key.expiryStatus)">
            <td>{{ key.keyName }}</td>
            <td>{{ key.expiryDate || '-' }}</td>
            <td>
              <span class="badge" :class="key.expiryStatus === 3 ? 'badge--expired' : 'badge--soon'">
                {{ statusLabel(key.expiryStatus) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="modal-actions">
        <button class="btn btn-primary" @click="emit('close')">확인</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.desc {
  color: var(--text-muted);
  font-size: 13px;
  margin-bottom: 16px;
  line-height: 1.6;
}
.expiry-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  margin-bottom: 4px;
}
.expiry-table th {
  background: var(--bg-input);
  padding: 8px 12px;
  text-align: left;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
}
.expiry-table td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--border);
}
.expiry-table tr.row--expired { background: rgba(231,76,60,0.07); }
.expiry-table tr.row--soon    { background: rgba(243,156,18,0.07); }

.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
}
.badge--expired { background: rgba(231,76,60,0.2); color: #e74c3c; }
.badge--soon    { background: rgba(243,156,18,0.2); color: #f39c12; }
</style>
