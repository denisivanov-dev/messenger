import axios from 'axios'

export const loginApi = async (credentials) => {
  const response = await axios.post('/api/auth/login', credentials)
  return response.data
}