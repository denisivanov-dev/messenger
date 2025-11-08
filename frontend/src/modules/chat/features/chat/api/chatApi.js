import { axiosInstance } from "../../connection/http/axiosInstance"

export async function startPrivateChat(targetId) {
  const response = await axiosInstance.post('/api/chats/private-chat', {
    target_id: targetId
  })

  return response.data
}
