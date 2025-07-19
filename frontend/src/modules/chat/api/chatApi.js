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

  // sendMessage({
  //   type:        "init_private",
  //   receiver_id: targetId
  // })

  return response.data
}