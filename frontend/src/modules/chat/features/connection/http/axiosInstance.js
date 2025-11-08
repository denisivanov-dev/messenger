import axios from 'axios'
import { useAuthStore } from '../../../../auth/store/authStore'

export const axiosInstance = axios.create({
  baseURL: 'http://localhost:8000',
  withCredentials: true,
})

axiosInstance.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    return config
  },
  (error) => Promise.reject(error)
)