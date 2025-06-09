import { loginApi, registerApi, confirmRegistrationApi } from '../api/authApi'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const accessToken = ref(null)

  const getAccessToken = computed(() => accessToken.value)
  const getUserId = computed(() => user.value?.id || null)
  const getUsername = computed(() => user.value?.username || null)
  const getEmail = computed(() => user.value?.email || null)

  const setAccessToken = (token) => {
    accessToken.value = token
  }
  const setUserId = (id) => {
    if (!user.value) user.value = {}
    user.value.id = id
  }
  const setUsername = (username) => {
    if (!user.value) user.value = {}
    user.value.username = username
  }
  const setEmail = (email) => {
    if (!user.value) user.value = {}
    user.value.email = email
  }

  const login = async ({ email, password }) => {
    try {
      const response = await loginApi({ email, password })
      user.value = response.user

      return response.message
    } catch (error) {
      const status = error.response?.status
		  const detail = error.response?.data?.detail

      if (typeof detail === 'object') {
        throw detail
      }
    }
  }

  const register = async ({ email, username, password, confirmPassword }) => {
    try {
      const response = await registerApi({ email, username, password, confirmPassword })
      user.value = response.user
      
      return response.message
    } catch (error) {
      const status = error.response?.status
		  const detail = error.response?.data?.detail

      if (typeof detail === 'object') {
        throw detail
      }
    }
  }

  const confirmRegistration = async ({ username, email, code }) => {
    try {
      const response = await confirmRegistrationApi({ username, email, code })
      accessToken.value = response.access_token

      return response.message
    } catch (error) {
      const status = error.response?.status
		  const detail = error.response?.data?.detail
      
      throw new Error(detail)
    }
  }

   return {
    user,
    accessToken,
    login,
    register,
    confirmRegistration,

    getAccessToken,
    getUserId,
    getUsername,
    getEmail,

    setAccessToken,
    setUserId,
    setUsername,
    setEmail
  }
})  