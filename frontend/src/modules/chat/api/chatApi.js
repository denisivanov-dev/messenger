import axios from 'axios'
import { useAuthStore } from '../../auth/store/authStore'

const axiosInstance = axios.create({
  baseURL: 'http://localhost:8000',
  withCredentials: true,
})

axiosInstance.interceptors.request.use(
  config => {
    const authStore = useAuthStore()
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    return config
  },
  error => Promise.reject(error)
)

// axiosInstance.interceptors.response.use(
//   res => res,
//   async err => {
//     if (err.response?.status === 401) {
//       try {
//         const authStore = useAuthStore()
//         const { access_token } = await authStore.refreshToken()
//         authStore.accessToken = access_token
//         err.config.headers.Authorization = `Bearer ${access_token}`
//         return axiosInstance.request(err.config)
//       } catch (refreshError) {
//         console.error('Не удалось обновить токен', refreshError)
//         // authStore.logout() или редирект
//       }
//     }
//     return Promise.reject(err)
//   }
// )

let socket = null
export function connect(token, onMessage, mode = 'global', receiverId = null) {
  if (!token) {
    console.error('WebSocket connection aborted: no token provided')
    return
  }

  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.close(1000, 'Reconnecting')
  }

  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`)

  socket.onopen = () => {
    console.log('WebSocket connected')

    if (mode === 'global') {
      socket.send(JSON.stringify({
        type: 'init_global',
        chat_type: 'global'
      }))
    } else if (mode === 'private' && receiverId) {
      socket.send(JSON.stringify({
        type: 'init_private',
        chat_type: 'private',
        receiver_id: receiverId
      }))
    }
  }

  socket.onmessage = event => onMessage(JSON.parse(event.data))
  socket.onclose = () => console.log('WebSocket disconnected')
  socket.onerror = err => console.error('WebSocket error', err)
}

export function sendMessage(message) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message))
  }
}

export function disconnect() {
  if (socket) socket.close()
}


// ---------- REST -----------------
export async function getAllUsers() {
  const response = await axiosInstance.get('/api/users/all')
  return response.data
}

export async function startPrivateChat(targetId) {
  const response = await axiosInstance.post('/api/chats/private-chat', {
    target_id: targetId
  })

  return response.data
}

export async function uploadFileToR2(file) {
  const formData = new FormData()
  formData.append('file', file)

  const response = await axiosInstance.post('/api/messages/upload-image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export async function getImagefromKey(key) {
  const response = await axiosInstance.get(`api/messages/${key}`)

  return response.data
}

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

export async function getCallRoomStatus(roomId) {
  const response = await axiosInstance.get(`/api/voice/${roomId}`)

  return response.data
}