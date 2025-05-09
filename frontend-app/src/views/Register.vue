<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <h1 class="auth-title">註冊</h1>
      <form @submit.prevent="register" class="auth-form">
        <input v-model="username" type="text" placeholder="使用者名稱" class="auth-input" />
        <input v-model="password" type="password" placeholder="密碼" class="auth-input" />
        <div class="password-rules">
          <div :class="{'rule-ok': passwordLengthValid, 'rule-no': !passwordLengthValid}">
            {{ passwordLengthValid ? '✅' : '❌' }} 至少8個字
          </div>
          <div :class="{'rule-ok': passwordHasLetter, 'rule-no': !passwordHasLetter}">
            {{ passwordHasLetter ? '✅' : '❌' }} 包含英文
          </div>
          <div :class="{'rule-ok': passwordHasNumber, 'rule-no': !passwordHasNumber}">
            {{ passwordHasNumber ? '✅' : '❌' }} 包含數字
          </div>
        </div>

        <div class="password-strength">
          <div class="strength-bar" :style="{ width: passwordStrength + '%', backgroundColor: strengthColor }"></div>
        </div>
        <p class="strength-text">{{ strengthText }}</p>


        <button
          type="submit"
          class="auth-button bg-blue-500 hover:bg-blue-600"
          :disabled="!canRegister"
        >
          註冊
        </button>
      </form>
      <p class="auth-link">
        已經有帳號了？
        <button @click="router.push('/login')" class="link-button">前往登入</button>
      </p>
      <p v-if="error" class="auth-error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const error = ref('')
const router = useRouter()

const passwordLengthValid = computed(() => password.value.length >= 8)
const passwordHasLetter = computed(() => /[a-zA-Z]/.test(password.value))
const passwordHasNumber = computed(() => /\d/.test(password.value))

const canRegister = computed(() => {
  return passwordLengthValid.value && passwordHasLetter.value && passwordHasNumber.value && username.value.trim() !== ''
})

function isPasswordValid(pwd: string): boolean {
  const hasMinLength = pwd.length >= 8
  const hasLetter = /[a-zA-Z]/.test(pwd)
  const hasNumber = /\d/.test(pwd)
  return hasMinLength && hasLetter && hasNumber
}

const passwordStrength = computed(() => {
  const pwd = password.value
  const length = pwd.length
  if (length < 8) return 0

  let score = 0

  let lengthScore = 0
  if (length >= 8) {
    lengthScore = Math.min(((length - 8) / (20 - 8)) * 40, 40)
  }
  score += lengthScore

  let uppercase = 0
  let lowercase = 0
  let numbers = 0
  let symbols = 0

  for (const c of pwd) {
    if (/[A-Z]/.test(c)) uppercase++
    else if (/[a-z]/.test(c)) lowercase++
    else if (/\d/.test(c)) numbers++
    else symbols++
  }

  let typeScore = 0
  if (uppercase > 0) typeScore += 5
  if (lowercase > 0) typeScore += 5
  if (numbers > 0) typeScore += 5
  if (symbols > 0) typeScore += 5

  score += typeScore

  const totalChars = uppercase + lowercase + numbers + symbols
  if (totalChars === 0) return score

  const ratios = [
    uppercase / totalChars,
    lowercase / totalChars,
    numbers / totalChars,
    symbols / totalChars
  ]

  const idealRatio = 0.25

  let balanceScore = 0
  for (const ratio of ratios) {
    const diffPercent = Math.abs(ratio - idealRatio) * 100
    let charTypeScore = 0

    if (diffPercent <= 10) {
      charTypeScore = 10 - Math.floor(diffPercent / 2)
    } else if (diffPercent <= 25) {
      charTypeScore = 5 - Math.floor((diffPercent - 10) / 3)
    } else if (diffPercent <= 35) {
      charTypeScore = 5 - Math.floor((diffPercent - 25) / 2)
    } else if (diffPercent <= 50) {
      charTypeScore = 0 + Math.max(0, 5 - Math.floor((diffPercent - 35) / 3))
    } else {
      charTypeScore = 0
    }

    balanceScore += Math.max(0, charTypeScore)
  }

  score += balanceScore

  return Math.min(Math.round(score), 100)
})

const strengthColor = computed(() => {
  if (passwordStrength.value <= 40) return '#e53e3e'
  if (passwordStrength.value <= 80) return '#d69e2e'
  return '#38a169'
})

const strengthText = computed(() => {
  if (passwordStrength.value <= 40) return '弱'
  if (passwordStrength.value <= 80) return '中等'
  return '強'
})

const register = async () => {
  error.value = ''

  if (!username.value.trim()) {
    error.value = '請輸入使用者名稱'
    return
  }

  if (!isPasswordValid(password.value)) {
    error.value = '密碼至少8位數，且需包含英文和數字'
    return
  }

  try {
    const res = await axios.post('http://localhost:8080/register', {
      username: username.value,
      password: password.value,
    })
    
    const token = res.data.token
    localStorage.setItem('token', token)
    localStorage.setItem('username', res.data.username)
    router.push('/tasks')
  } catch (err: any) {
    error.value = err.response?.data?.error || '註冊失敗'
  }
}
</script>
