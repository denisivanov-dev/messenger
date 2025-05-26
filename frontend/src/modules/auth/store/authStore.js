import { loginApi } from '../api/authApi'
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
  }),
  actions: {
    async login(credentials) {
      try {
        const data = await loginApi(credentials)
        this.user = data.user
        this.token = data.token
        console.log('Успешный вход:', data)
      } catch (err) {
        console.error('Ошибка входа:', err)
      }
    },
  },
})