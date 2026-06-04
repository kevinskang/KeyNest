<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Register } from '../../wailsjs/go/main/App'
import { useToast } from '../composables/useToast'

const router = useRouter()
const { show } = useToast()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const errorMsg = ref('')
const loading = ref(false)

async function handleRegister() {
  errorMsg.value = ''
  if (!email.value || !password.value || !confirmPassword.value) {
    errorMsg.value = '모든 항목을 입력해주세요.'
    return
  }
  if (password.value !== confirmPassword.value) {
    errorMsg.value = '비밀번호가 일치하지 않습니다.'
    return
  }
  if (password.value.length < 8) {
    errorMsg.value = '비밀번호는 최소 8자 이상이어야 합니다.'
    return
  }
  loading.value = true
  try {
    const err = await Register(email.value, password.value)
    if (err) {
      errorMsg.value = err
      return
    }
    show('회원가입이 완료되었습니다. 로그인해주세요.', 'success')
    router.push('/login')
  } catch {
    errorMsg.value = '회원가입 중 오류가 발생했습니다.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <div class="auth-logo">🔑</div>
      <h1 class="auth-title">KeyNest</h1>
      <p class="auth-subtitle">새 계정 만들기</p>

      <form class="auth-form" @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">이메일</label>
          <input
            v-model="email"
            type="email"
            placeholder="user@example.com"
            autocomplete="username"
            :disabled="loading"
          />
        </div>
        <div class="form-group">
          <label class="form-label">비밀번호 <span class="hint">(최소 8자)</span></label>
          <input
            v-model="password"
            type="password"
            placeholder="••••••••"
            autocomplete="new-password"
            :disabled="loading"
          />
        </div>
        <div class="form-group">
          <label class="form-label">비밀번호 확인</label>
          <input
            v-model="confirmPassword"
            type="password"
            placeholder="••••••••"
            autocomplete="new-password"
            :disabled="loading"
            @keyup.enter="handleRegister"
          />
        </div>
        <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>
        <button class="btn btn-primary auth-submit" type="submit" :disabled="loading">
          {{ loading ? '처리 중...' : '회원가입' }}
        </button>
      </form>

      <p class="auth-footer">
        이미 계정이 있으신가요?
        <RouterLink to="/login">로그인</RouterLink>
      </p>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
.auth-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 40px 48px;
  width: 400px;
  max-width: 92vw;
  text-align: center;
}
.auth-logo  { font-size: 40px; margin-bottom: 8px; }
.auth-title { font-size: 26px; font-weight: 800; letter-spacing: -0.5px; }
.auth-subtitle { color: var(--text-muted); margin-bottom: 28px; }

.auth-form  { display: flex; flex-direction: column; gap: 14px; text-align: left; }
.auth-submit { width: 100%; justify-content: center; padding: 10px; font-size: 15px; margin-top: 4px; }

.hint { font-weight: 400; color: var(--text-muted); font-size: 11px; }
.auth-footer { margin-top: 20px; color: var(--text-muted); font-size: 13px; }
.auth-footer a { color: var(--primary); text-decoration: none; font-weight: 700; }
.auth-footer a:hover { text-decoration: underline; }
</style>
