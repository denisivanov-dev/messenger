import { axiosInstance } from "../../connection/http/axiosInstance"

export async function getAllUsers() {
  const response = await axiosInstance.get('/api/users/all')
  
  return response.data
}
