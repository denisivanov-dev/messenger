<template>
  <div class="w-full min-h-screen bg-[#f9f5f0] flex items-center justify-center">
    <div class="w-[500px] max-w-md bg-white text-[#2e2e2e] p-8 rounded-2xl shadow-lg">
      <h2 class="text-2xl font-bold mb-6 text-center">Восстановление пароля</h2>

      <form @submit.prevent="onSubmit" class="flex flex-col space-y-6">
        <div>
          <label class="block text-sm mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]"
          />
        </div>

        <div>
          <label class="block text-sm mb-1">Код подтверждения</label>
          <input
            v-model="code"
            type="text"
            class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]"
          />
        </div>

        <template v-if="code.length >= 4">
          <div>
            <label class="block text-sm mb-1">Новый пароль</label>
            <input
              v-model="newPassword"
              type="password"
              class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]"
            />
          </div>

          <div>
            <label class="block text-sm mb-1">Подтвердите пароль</label>
            <input
              v-model="confirmPassword"
              type="password"
              class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b]"
            />
          </div>
        </template>

        <button
          type="submit"
          class="bg-[#e07a5f] hover:bg-[#cf6951] transition-colors p-3 rounded-xl font-medium text-white"
        >
          Сбросить пароль
        </button>

        <button
          type="button"
          @click="router.push('/login')"
          class="text-sm text-[#a0625e] hover:underline mt-2"
        >
          Вернуться к логину
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const email = ref('')
const code = ref('')
const newPassword = ref('')
const confirmPassword = ref('')

const router = useRouter()

const onSubmit = () => {
  if (newPassword.value !== confirmPassword.value) {
    alert('Пароли не совпадают')
    return
  }

  console.log('Сброс пароля:', {
    email: email.value,
    code: code.value,
    newPassword: newPassword.value,
  })

  // здесь можно вызвать API или отправить форму
}
</script>