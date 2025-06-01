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