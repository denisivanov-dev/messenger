import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: '/api',
  withCredentials: true
})

export const loginApi = async (credentials) => {
  const response = await axiosInstance.post('/auth/login', credentials || {})
  return response.data
}

export const registerApi = async (credentials) => {
  const response = await axiosInstance.post('/auth/register', credentials || {})
  return response.data
}

export const confirmRegistrationApi = async (credentials) => {
  const response = await axiosInstance.post('/auth/confirm-registration', credentials || {})
  return response.data
}

export const autoLoginApi = async () => {
  const response = await axiosInstance.get('/auth/auto-login', { withCredentials: true })
  return response.data
}

export const refreshTokenApi = async () => {
  const response = await axiosInstance.post('/auth/refresh', {}, { withCredentials: true })
  return response.data
}