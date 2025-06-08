import axios from 'axios'

const axiosInstance = axios.create({
  baseURL: 'http://localhost:8000', 
  withCredentials: true, 
})

export const loginApi = async (credentials) => {
  const response = await axiosInstance.post('/api/auth/login', credentials || {})
  return response.data
}

export const registerApi = async (credentials) => {
  const response = await axiosInstance.post('/api/auth/register', credentials || {})
  return response.data
}

export const confirmRegistrationApi = async (credentials) => {
  const response = await axiosInstance.post('/api/auth/confirm-registration', credentials || {})
  return response.data
}

export const autoLoginApi = async () => {
	const response = await axiosInstance.get('/api/auth/auto-login', {withCredentials: true })
	return response.data
}