import { axiosInstance } from "../../connection/http/axiosInstance"

export async function apiSendFriendRequest(fromId, toId) {
  const response = await axiosInstance.post('/api/friends/request', {
    from_id: fromId,
    to_id: toId,
  })

  return response.data
}

export async function apiCancelFriendRequest(fromId, toId) {
  const response = await axiosInstance.post('/api/friends/cancel', {
    from_id: fromId,
    to_id: toId,
  })

  return response.data
}

export async function apiAcceptFriendRequest(fromId, toId) {
  const response = await axiosInstance.post('/api/friends/accept', {
    from_id: fromId,
    to_id: toId,
  })

  return response.data
}

export async function apiDeclineFriendRequest(fromId, toId) {
  const response = await axiosInstance.post('/api/friends/decline', {
    from_id: fromId,
    to_id: toId,
  })

  return response.data
}

export async function apiDeleteFriend(fromId, toId) {
  const response = await axiosInstance.post('/api/friends/delete', {
    from_id: fromId,
    to_id: toId,
  })

  return response.data
}

export async function apiGetFriends(userId) {
  const response = await axiosInstance.get(`/api/friends/list/${userId}`)

  return response.data
}