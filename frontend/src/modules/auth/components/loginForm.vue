<template>
  <div class="w-full min-h-screen bg-[#f9f5f0] flex items-center justify-center">
    <div id="loginBox" class="w-[500px] max-w-md bg-white text-[#2e2e2e] p-8 rounded-2xl shadow-lg">
      <h2 class="text-2xl font-bold mb-6 text-center">Добро пожаловать!</h2>

      <form @submit.prevent="onSubmit" class="flex flex-col">
        <div>
          <label class="block text-sm text-[#2e2e2e] mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b] box-border"
          />
        </div>

        <div class="mt-7">
          <label class="block text-sm text-[#6c6c6c] mb-1">Пароль</label>
          <input
            v-model="password"
            type="password"
            class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b] box-border"
          />
        </div>

        <div class="text-left">
          <button
            type="button"
            @click="onForgotPassword"
            class="text-sm text-[#a0625e] hover:underline cursor-pointer bg-transparent p-0 focus:outline-none"
          >
            Забыли пароль?
          </button>
        </div>

        <button 
          type="submit"
          class="bg-[#e07a5f] hover:bg-[#cf6951] 
          transition-colors p-3 rounded-xl font-medium box-border mt-7 text-white"
        >
          Войти
        </button>

        <div class="mt-2 flex items-center gap-1">
          <span class="text-sm text-[#2e2e2e]">Нужна учётная запись?</span>
          <button
            type="button"
            @click="onRegister"
            class="text-sm text-[#a0625e] hover:underline cursor-pointer bg-transparent p-0 focus:outline-none font-normal leading-none"
          >
            Зарегистрироваться
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'

const email = ref('')
const password = ref('')

const store = useAuthStore()
const router = useRouter()

const onSubmit = () => {
  store.login({ email: email.value, password: password.value })
}

const onForgotPassword = () => {
  router.push('/forgot-password')
}
</script>
