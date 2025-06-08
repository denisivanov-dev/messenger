<template>
   <div class="w-full min-h-screen bg-[#f9f5f0] flex items-center justify-center">
      <div id="registerBox" class="w-[500px] max-w-md bg-white text-[#2e2e2e] p-8 rounded-2xl shadow-lg">
         <h2 class="text-2xl font-bold mb-6 text-center">Регистрация</h2>

         <form @submit.prevent="onSubmit" class="flex flex-col">
            <!-- Email -->
            <div>
               <div class="flex items-center mb-2">
                  <label class="text-sm text-[#2e2e2e]">Email</label>
                  <p v-if="emailError" class="text-sm text-red-500 mr-auto ml-1">{{ emailError }}</p>
               </div>
               <input
                  v-model="email"
                  type="text"
                  class="w-full bg-white border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]" />
            </div>

            <!-- Name -->
            <div class="mt-5">
               <div class="flex items-center mb-2">
                  <label class="text-sm text-[#2e2e2e]">Имя пользователя</label>
                  <p v-if="usernameError" class="text-sm text-red-500 mr-auto ml-1">{{ usernameError }}</p>
               </div>
               <input
                  v-model="username"
                  type="text"
                  class="w-full bg-white border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]" />
            </div>

            <!-- Password -->
            <div class="mt-5">
               <div class="flex items-center mb-2">
                  <label class="text-sm text-[#6c6c6c]">Пароль</label>
                  <p v-if="passwordError" class="text-sm text-red-500 mr-auto ml-1">{{ passwordError }}</p>
               </div>
               <input
                  v-model="password"
                  type="password"
                  class="w-full bg-white border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]" />
            </div>

            <!-- Confirm Password -->
            <div class="mt-5">
               <div class="flex items-center mb-2">
                  <label class="text-sm text-[#6c6c6c]">Подтвердите пароль</label>
                  <p v-if="confirmPasswordError" class="text-xs text-red-500 mr-auto ml-1">{{ confirmPasswordError }}</p>
               </div>
               <input
                  v-model="confirmPassword"
                  type="password"
                  class="w-full bg-white border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]" />
            </div>

            <!-- Submit -->
            <button
               type="submit"
               class="bg-[#e07a5f] hover:bg-[#cf6951] transition-colors p-3 rounded-xl font-medium box-border mt-7 text-white">
               Зарегистрироваться
            </button>

            <!-- Go to login -->
            <div class="mt-2 flex items-center gap-1">
               <span class="text-sm text-[#2e2e2e]">Уже зарегистрированы?</span>
               <button
                  type="button"
                  @click="onGoToLogin"
                  class="text-sm text-[#a0625e] hover:underline cursor-pointer bg-transparent p-0 focus:outline-none font-normal leading-none">
                  Войти
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
import { validateRegisterForm } from '../validators/registerValidator'
import { applyErrors } from '../utils/errorHandler'

const store = useAuthStore()
const router = useRouter()

const email = ref('')
const username = ref('')
const password = ref('')
const confirmPassword = ref('')

const emailError = ref('')
const usernameError = ref('')
const passwordError = ref('')
const confirmPasswordError = ref('')

const onSubmit = async () => {
   const data = {
      email: email.value,
      username: username.value,
      password: password.value,
      confirmPassword: confirmPassword.value 
   }

   const errorFields = {
      email: emailError,
      username: usernameError,
      password: passwordError,
      confirmPassword: confirmPasswordError
   }

   // const validationErrors = validateRegisterForm(data)

   emailError.value = ''
   usernameError.value = ''
   passwordError.value = ''
   confirmPasswordError.value = ''

   // if (validationErrors) {
   //    applyErrors(validationErrors, errorFields)
   //    return
   // }

   try {
      const response = await store.register(data)
      if (response === 'success') {
         router.push('/confirm-registration')
      }
   } catch (backendErrors) {
      applyErrors(backendErrors, errorFields)
   }
}

const onGoToLogin = () => {
   router.push('/login')
}
</script>
