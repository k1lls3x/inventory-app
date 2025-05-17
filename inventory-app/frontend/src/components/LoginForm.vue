<template>
  <div class="auth-bg">
    <div class="auth-glass">
      <div class="avatar">
  <svg width="54" height="54" viewBox="0 0 54 54" fill="none">
    <circle cx="27" cy="27" r="27" fill="#f2f3f6"/>
    <path d="M27 27c4.418 0 8-3.582 8-8s-3.582-8-8-8-8 3.582-8 8 3.582 8 8 8zm0 2c-5.33 0-16 2.67-16 8v3c0 1.105.895 2 2 2h28c1.105 0 2-.895 2-2v-3c0-5.33-10.67-8-16-8z" fill="#d1d5db"/>
  </svg>
</div>
      <h1 class="app-title">Складская система</h1>
      <div class="subtitle">Добро пожаловать! Войдите в свой аккаунт.</div>
      <form class="login-form" @submit.prevent="handleLogin">
        <div class="input-group">
          <span class="input-icon">
            <svg width="18" height="18" fill="none" viewBox="0 0 20 20"><path d="M10 10a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm0 2c-3.32 0-6 1.34-6 3v1a1 1 0 0 0 1 1h10a1 1 0 0 0 1-1v-1c0-1.66-2.68-3-6-3Z" fill="#ccc"/></svg>
          </span>
          <input
            type="text"
            v-model.trim="login"
            :class="{ 'invalid': loginError }"
            placeholder="Логин"
            @blur="validateLogin"
            @input="loginError = ''"
            autocomplete="username"
          />
          <transition name="fade">
            <div class="custom-error" v-if="loginError">{{ loginError }}</div>
          </transition>
        </div>
        <div class="input-group">
          <span class="input-icon">
            <svg width="18" height="18" fill="none" viewBox="0 0 20 20"><path d="M10 13a2 2 0 1 0 0-4 2 2 0 0 0 0 4Zm6-2a6 6 0 1 0-12 0 6 6 0 0 0 12 0Zm-2 0a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z" fill="#ccc"/></svg>
          </span>
          <input
            :type="showPassword ? 'text' : 'password'"
            v-model.trim="password"
            :class="{ 'invalid': passwordError }"
            placeholder="Пароль"
            @blur="validatePassword"
            @input="passwordError = ''"
            autocomplete="current-password"
          />

          <transition name="fade">
            <div class="custom-error" v-if="passwordError">{{ passwordError }}</div>
          </transition>
        </div>
        <button type="submit" class="login-btn" :disabled="loading">
          {{ loading ? 'Входим...' : 'Войти' }}
        </button>
        <div class="forgot-password">
          <a href="#" @click.prevent>Забыли пароль?</a>
        </div>
        <transition name="fade">
          <div v-if="error" class="error" :class="{ shake: error }">{{ error }}</div>
        </transition>
      </form>
      <div class="copyright">© 2025 Складская система</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const login = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)
const loginError = ref('')
const passwordError = ref('')
const showPassword = ref(false)

const emit = defineEmits(['login-success'])

function validateLogin() {
  if (!login.value.trim()) {
    loginError.value = 'Заполните это поле';
    return false
  }
  return true
}

function validatePassword() {
  if (!password.value.trim()) {
    passwordError.value = 'Заполните это поле';
    return false
  }
  return true
}

function handleLogin() {
  error.value = ''
  if (!validateLogin() | !validatePassword()) return

  loading.value = true
  setTimeout(() => {
    if (login.value === 'admin' && password.value === '1234') {
      localStorage.setItem('loggedIn', 'true')
      emit('login-success')
    } else {
      error.value = 'Неверный логин или пароль'
      setTimeout(() => error.value = '', 1200)
    }
    loading.value = false
  }, 700)
}

function toggleShowPassword() {
  showPassword.value = !showPassword.value
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;600&display=swap');

.auth-bg {
  min-height: 100vh;
  min-width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at 50% 40%, #f5f6fa 0%, #ececec 100%);
}
.auth-glass {
  background: rgba(255,255,255,0.92);
  border-radius: 34px; /* Было 26px */
  box-shadow: 0 16px 48px 0 rgba(40, 40, 40, 0.18), 0 6px 24px 0 rgba(0,0,0,0.09);
  padding: 2.7rem 2.5rem 2.1rem 2.5rem;
  width: 420px;
  max-width: 98vw;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2.7rem 2.5rem 2.1rem 2.5rem;
  backdrop-filter: blur(15px) saturate(140%);
  border: 2px solid #ececec;
  animation: fadeIn 0.7s cubic-bezier(.77,0,.18,1) both;
  position: relative;
  transition: box-shadow 0.25s cubic-bezier(.6,0,.4,1);
}
.auth-glass:hover {
  box-shadow: 0 24px 64px 0 rgba(50, 60, 90, 0.22), 0 10px 36px 0 rgba(0,0,0,0.12);
}
.avatar {
 /* больше отступ */
  width: 70px;
  height: 70px;
  border-radius: 50%;
  margin-bottom: 1.8rem; /* больше отступ после аватара */
  background: linear-gradient(145deg,#f7f7f7 60%,#e8e8e8 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  /* box-shadow: none; */ /* УБРАТЬ ЭТУ СТРОКУ! */
  box-shadow: none;      /* ← новая строка: нет тени! */
}
.app-title {
  font-family: 'Inter', Arial, sans-serif;
  margin-bottom: 0.25rem;        /* уменьшить нижний маргин */
  margin-top: -0.8rem;           /* поднять выше! */
  font-weight: 700;
  font-size: 2.13rem;
  letter-spacing: .13px;
  color: #111;
  text-align: center;
  background: linear-gradient(90deg, #111 5%, #838383 85%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  user-select: none;
}


.subtitle {
  margin-bottom: 1.5rem;      /* чуть больше отступ вниз */
  margin-top: 0.3rem;         /* и сверху, чтобы не липло */
  color: #888;
  font-size: 1.09rem;
  font-weight: 500;
  text-align: center;
  user-select: none;
}

.login-form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.65rem;
  font-family: 'Inter', Arial, sans-serif;
}
.input-group {
  margin-bottom: 1.25rem;
  position: relative;
  width: 100%;
  padding-bottom: 0.2rem;
  display: flex;
  align-items: center;
}
.input-icon {
  position: absolute;
  left: 1.3rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  opacity: 0.63;
  z-index: 2;
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.input-icon svg {
  width: 22px;
  height: 22px;
}
.eye-icon {
  position: absolute;
  right: 1.15rem;
  top: 50%;
  transform: translateY(-50%);
  cursor: pointer;
  opacity: 1;
  z-index: 3;
  transition: color 0.17s;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;

  color: #bdbdbd;
}

.eye-icon svg {
  width: 28px;
  height: 28px;
  fill: currentColor;
}


.eye-icon:hover { opacity: 1; }
input {
  width: 100%;
  padding: 1.15rem 3.1rem 1.15rem 2.8rem;
    font-size: 1.13rem;
  border: 1.5px solid #dadada;
  border-radius: 16px;
  outline: none;
  background: #fafafa;
  color: #444;
  transition: border 0.17s, box-shadow 0.18s, background 0.15s;
  font-family: inherit;
  box-shadow: 0 2px 8px 0 rgba(50,50,50,0.06);
  margin: 0;
  box-sizing: border-box;
}
input:focus {
  border-color: #3578ea;
  background: #f3f6fc;
  box-shadow: 0 0 0 2px #3578ea22, 0 4px 16px 0 rgba(70,100,220,0.11);
}
input.invalid {
  border-color: #e53935 !important;
  background: #fff6f6;
  box-shadow: 0 0 0 2px #e5393522, 0 4px 16px 0 #e5393540;
}


.custom-error {
  position: absolute;
  left: 0; right: 0;
  top: 100%;
  margin-top: 0.30rem;
  padding: 0.42em 1.12em;
  background: #e53935;
  color: #fff;
  font-size: 1.04em;
  border-radius: 7px;
  box-shadow: 0 2px 10px 0 #d5bcbc2e;
  text-align: left;
  z-index: 10;
  animation: fadeIn .26s;
  pointer-events: none;
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.login-btn {
  margin-top: 0.2em;
  padding: 1.02rem 1.21rem;
  background: linear-gradient(90deg,#3578ea 25%, #111 100%);
  color: #fff;
  border-radius: 18px; /* Было 9px */
  border: none;
  font-size: 1.13rem;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.14s, box-shadow 0.18s, transform 0.11s;
  box-shadow: 0 4px 16px 0 rgba(40,40,40,0.13);
  letter-spacing: .03em;
  outline: none;
}
.login-btn:hover:not(:disabled) {
  background: linear-gradient(90deg,#255fdc 15%, #232323 100%);
  box-shadow: 0 8px 32px 0 rgba(40,40,40,0.18);
  transform: translateY(-2px) scale(1.017);
}
.login-btn:active:not(:disabled) {
  background: #0e3476;
  transform: scale(0.98);
}
.login-btn:disabled {
  opacity: 0.68;
  cursor: wait;
}
.forgot-password {
  margin-top: 1.1em;
  text-align: right;
  font-size: 0.97rem;
}
.forgot-password a {
  color: #3578ea;
  text-decoration: none;
  transition: color 0.13s;
}
.forgot-password a:hover { color: #174487; }
.error {
  color: #fff;
  background: #e53935;
  text-align: center;
  font-size: 1.04rem;
  border-radius: 7px;
  padding: 0.58em 0;
  margin-top: -0.5em;
  letter-spacing: 0.02em;
  box-shadow: 0 2px 8px 0 #d5bcbc22;
  animation: shake 0.24s;
}
@keyframes shake {
  0% { transform: translateX(0);}
  23% { transform: translateX(-6px);}
  46% { transform: translateX(6px);}
  69% { transform: translateX(-5px);}
  92% { transform: translateX(3px);}
  100% { transform: translateX(0);}
}
@keyframes fadeIn {
  0% { transform: translateY(32px) scale(0.98); opacity: 0; }
  100% { transform: none; opacity: 1; }
}
.copyright {
  margin-top: 1.2em;
  font-size: 0.93em;
  color: #aaa;
  text-align: center;
  user-select: none;
  opacity: 0.7;
}
</style>
