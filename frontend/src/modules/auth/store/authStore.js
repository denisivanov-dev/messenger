import { loginApi, registerApi } from '../api/authApi'
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)

  const login = async ({ email, password }) => {
    try {
      const response = await loginApi({ email, password })
      user.value = response.user
    } catch (error) {
      throw error.response?.data?.message || 'Ошибка авторизации'
    }
  }

  const register = async ({ email, username, password, confirmPassword }) => {
    try {
      console.info({ email, username, password, confirmPassword })
      const response = await registerApi({ email, username, password, confirmPassword })
      
      console.info(response)
      user.value = response.user
    } catch (error) {
      const status = error.response.status
      const detail = error.response.data.detail

      if (typeof detail === 'object') {
        throw detail
      }
    }
  }

  return { user, login, register }
})