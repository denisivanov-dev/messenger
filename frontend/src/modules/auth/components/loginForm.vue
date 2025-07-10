<template>
	<div class="w-full min-h-screen bg-[#f9f5f0] flex items-center justify-center">
		<div id="loginBox" class="w-[500px] max-w-md bg-white text-[#2e2e2e] p-8 rounded-2xl shadow-lg">
			<h2 class="text-2xl font-bold mb-6 text-center">Добро пожаловать!</h2>

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
						class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b] box-border" />
				</div>

				<!-- Password -->
				<div class="mt-7">
					<div class="flex items-center mb-2">
						<label class="text-sm text-[#6c6c6c]">Пароль</label>
						<p v-if="passwordError" class="text-sm text-red-500 mr-auto ml-1">{{ passwordError }}</p>
					</div>
					<input
						v-model="password"
						type="password"
						class="w-full bg-white text-[#2e2e2e] border border-[#ccc] p-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#3d405b] box-border" />
				</div>

				<!-- Forgot password -->
				<div class="mt-1 text-left">
					<button
						type="button"
						@click="onForgotPassword"
						class="text-sm text-[#a0625e] hover:underline cursor-pointer bg-transparent p-0 focus:outline-none">
						Забыли пароль?
					</button>
				</div>

				<!-- Submit -->
				<button
					type="submit"
					class="bg-[#e07a5f] hover:bg-[#cf6951] transition-colors p-3 rounded-xl font-medium box-border mt-7 text-white">
					Войти
				</button>

				<!-- Go to register -->
				<div class="mt-2 flex items-center gap-1">
					<span class="text-sm text-[#2e2e2e]">Нужна учётная запись?</span>
					<button
						type="button"
						@click="onGoToRegister"
						class="text-sm text-[#a0625e] hover:underline cursor-pointer bg-transparent p-0 focus:outline-none font-normal leading-none">
						Зарегистрироваться
					</button>
				</div>
			</form>
		</div>
	</div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../store/authStore'
import { validateLoginForm } from '../validators/loginValidator'
import { applyErrors } from '../utils/errorHandler'
import { autoLoginApi } from '../api/authApi'

const store = useAuthStore()
const router = useRouter()

const email = ref('')
const password = ref('')

const emailError = ref('')
const passwordError = ref('')

const onSubmit = async () => {
	emailError.value = ''
	passwordError.value = ''
	
	const data = {
      email: email.value,
      password: password.value
   }

   const errorFields = {
      email: emailError,
      password: passwordError
   }
	
	const validationErrors = validateLoginForm(data)
	if (validationErrors) {
      applyErrors(validationErrors, errorFields)
      return
   }

	try {
    	const response = await store.login(data)

		if (response === 'success') {
        	router.push('/global-chat')
			return
    	}

		if (response === 'need_code_verification') {
			router.push('/confirm-login')
			return
		}
   } catch (backendErrors) {
      applyErrors(backendErrors, errorFields)
   }
}

const onForgotPassword = () => {
	router.push('/forgot-password')
}

const onGoToRegister = () => {
	router.push('/register')
}
</script>
