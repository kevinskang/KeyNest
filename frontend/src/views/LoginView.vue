<script lang="ts" setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Login } from '../../wailsjs/go/main/App'
import { useAuthStore } from '../stores/auth'
import { useKeyStore } from '../stores/keys'

const router = useRouter()
const auth = useAuthStore()
const keyStore = useKeyStore()

const email = ref('')
const password = ref('')
const errorMsg = ref('')
const loading = ref(false)

async function handleLogin() {
  errorMsg.value = ''
  if (!email.value || !password.value) {
    errorMsg.value = '이메일과 비밀번호를 입력해주세요.'
    return
  }
  loading.value = true
  try {
    const err = await Login(email.value, password.value)
    if (err) {
      errorMsg.value = err
      return
    }
    auth.setLoggedIn(true)
    await keyStore.loadExpiringKeys()
    router.push('/')
  } catch {
    errorMsg.value = '로그인 중 오류가 발생했습니다.'
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
      <p class="auth-subtitle">API 키 관리 도구</p>

      <form class="auth-form" @submit.prevent="handleLogin">
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
          <label class="form-label">비밀번호</label>
          <input
            v-model="password"
            type="password"
            placeholder="••••••••"
            autocomplete="current-password"
            :disabled="loading"
            @keyup.enter="handleLogin"
          />
        </div>
        <p v-if="errorMsg" class="form-error">{{ errorMsg }}</p>
        <button class="btn btn-primary auth-submit" type="submit" :disabled="loading">
          {{ loading ? '로그인 중...' : '로그인' }}
        </button>
      </form>

      <p class="auth-footer">
        계정이 없으신가요?
        <RouterLink to="/register">회원가입</RouterLink>
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

.auth-footer { margin-top: 20px; color: var(--text-muted); font-size: 13px; }
.auth-footer a { color: var(--primary); text-decoration: none; font-weight: 700; }
.auth-footer a:hover { text-decoration: underline; }
</style>
