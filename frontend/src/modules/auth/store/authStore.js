import {
  loginApi,
  registerApi,
  confirmRegistrationApi,
  autoLoginApi,
  refreshTokenApi
} from '../api/authApi'
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const accessToken = ref(null)

  const getAccessToken = computed(() => accessToken.value)
  const getUserId      = computed(() => user.value?.id       ?? null)
  const getUsername    = computed(() => user.value?.username ?? null)
  const getEmail       = computed(() => user.value?.email    ?? null)

  const setAccessToken = token => { accessToken.value = token ?? null }
  const resetAuth      = ()   => { user.value = null; accessToken.value = null }

  const login = async ({ email, password }) => {
    try {
      const response = await loginApi({ email, password })
      user.value = response.user
      accessToken.value = response.accessToken
      return response.message
    } catch (error) {
      const detail = error.response?.data?.detail
      if (typeof detail === 'object') throw detail
    }
  }

  const register = async ({ email, username, password, confirmPassword }) => {
    try {
      const response = await registerApi({ email, username, password, confirmPassword })
      user.value = response.user
      console.info(user.value)
      return response.message
    } catch (error) {
      const detail = error.response?.data?.detail
      if (typeof detail === 'object') throw detail
    }
  }

  const confirmRegistration = async ({ username, email, code }) => {
    try {
      const response = await confirmRegistrationApi({ username, email, code })
      accessToken.value = response.accessToken
      return response.message
    } catch (error) {
      throw new Error(error.response?.data?.detail)
    }
  }

  const autoLogin = async () => {
    try {
      const response = await autoLoginApi()
      if (response?.message === 'success' && response?.accessToken) {
        accessToken.value = response.accessToken
        user.value = response.user
      }
    } catch {}
  }

  const refreshToken = async () => {
    try {
      const response = await refreshTokenApi()
      accessToken.value = response.access_token
      return true
    } catch {
      resetAuth()
      return false
    }
  }

  return {
    user,
    accessToken,
    login,
    register,
    confirmRegistration,
    autoLogin,
    refreshToken,
    getAccessToken,
    getUserId,
    getUsername,
    getEmail,
    setAccessToken,
    resetAuth
  }
})
