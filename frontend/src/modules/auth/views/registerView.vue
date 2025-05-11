<template>
  <div class="p-8 max-w-md mx-auto">
    <h1 class="text-2xl font-bold mb-4">Регистрация</h1>
    <form @submit.prevent="submit">
      <input v-model="form.username" placeholder="Username" class="input" required /><br />
      <input v-model="form.email" type="email" placeholder="Email" class="input" required /><br />
      <input v-model="form.password" type="password" placeholder="Password" class="input" required /><br />
      <button class="btn" type="submit">Зарегистрироваться</button>
    </form>
    <p v-if="error" class="text-red-500 mt-2">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { registerUser } from '@/api/auth'

const form = ref({
  username: '',
  email: '',
  password: ''
})

const error = ref('')

async function submit() {
  error.value = ''
  try {
    await registerUser(form.value)
    alert('Регистрация прошла успешно!')
  } catch (err) {
    error.value = err.message || 'Ошибка регистрации'
  }
}
</script>

<style scoped>
.input {
  margin-bottom: 10px;
  padding: 8px;
  width: 100%;
  border: 1px solid #ddd;
  border-radius: 6px;
}
.btn {
  padding: 10px 16px;
  background-color: #4f46e5;
  color: white;
  border-radius: 6px;
}
</style>