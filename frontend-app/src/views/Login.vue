<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h1 class="auth-title">登入</h1>
      <form @submit.prevent="login" class="auth-form">
        <input v-model="username" type="text" placeholder="使用者名稱" class="auth-input" />
        <input v-model="password" type="password" placeholder="密碼" class="auth-input" />
        <button type="submit" class="auth-button bg-green-500 hover:bg-green-600">登入</button>
      </form>
      <p class="auth-link">
        還沒有帳號？
        <button @click="router.push('/register')" class="link-button">前往註冊</button>
      </p>
      <p v-if="error" class="auth-error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()

const login = async () => {
  try {
    const res = await axios.post('http://localhost:8080/login', {
      username: username.value,
      password: password.value
    })
    
    const token = res.data.token
    localStorage.setItem('token', token)
    localStorage.setItem('username', res.data.username)
    router.push('/tasks')
  } catch (err: any) {
    error.value = err.response?.data?.error || '登入失敗'
  }
}
</script>
