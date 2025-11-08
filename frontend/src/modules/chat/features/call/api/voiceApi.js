import { axiosInstance } from "../../connection/http/axiosInstance"

export async function getCallRoomStatus(roomId) {
  const response = await axiosInstance.get(`/api/voice/${roomId}`)
  
  return response.data
}
