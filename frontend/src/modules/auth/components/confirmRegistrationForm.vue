<template>
   <div class="w-full min-h-screen bg-[#f9f5f0] flex items-center justify-center">
      <div id="confirmBox" class="w-[500px] max-w-md bg-white text-[#2e2e2e] p-8 rounded-2xl shadow-lg">
         <h2 class="text-2xl font-bold mb-6 text-center">Подтверждение регистрации</h2>

         <form @submit.prevent="onSubmit" class="flex flex-col">
            <!-- Code -->
            <div>
               <div class="flex items-center mb-2">
                  <label class="text-sm text-[#2e2e2e]">Код подтверждения</label>
                  <p v-if="codeError" class="text-sm text-red-500 ml-1">{{ codeError }}</p>
               </div>
               <input
                  v-model="code"
                  type="text"
                  maxlength="6"
                  class="w-full bg-white border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b] transition-all"
                  placeholder="Введите 6-значный код" />
               <p class="text-xs text-gray-500 mt-1">Код был отправлен на вашу почту. Проверьте папку «Спам».</p>
            </div>

            <!-- Submit -->
            <button
               type="submit"
               class="bg-[#e07a5f] hover:bg-[#cf6951] transition-colors p-3 rounded-xl font-medium box-border mt-7 text-white">
               Подтвердить
            </button>

            <!-- Resend Code -->
            <button
               type="button"
               @click="onResendCode"
               class="border border-[#e07a5f] text-[#e07a5f] hover:bg-[#fcebe8] transition-colors p-3 rounded-xl font-medium box-border mt-3">
               Отправить код ещё раз
            </button>

            <!-- Optional Timer -->
            <p class="text-xs text-red-500 text-center mt-1"></p>

            <!-- Back to login -->
            <div class="mt-4 flex items-center gap-1">
               <span class="text-sm text-[#2e2e2e]">Уже подтвердили?</span>
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
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'

const store = useAuthStore()
const router = useRouter()

const code = ref('')
const codeError = ref('')

const onSubmit = async () => {
   codeError.value = ''
   if (code.value.length !== 6) {
      codeError.value = '- Введите 6 цифр'
      return
   }

   try {
      const response = await store.confirmRegistration({
         username: store.getUsername,
         email: store.getEmail,
         code: code.value
      })

      if (response === 'success') {
         router.push('/global-chat')
      }
   } catch (error) {
      console.info(error)
      codeError.value = error.message
   }
}

const onResendCode = async () => {
   codeError.value = ''
   
   const email = computed(() => store.user?.email)

   try {
      const response = await store.resendCode(email.value)
   } catch (error) {

   }
}

const onGoToLogin = () => {
   router.push('/login')
}
</script>
