import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'
import {
  apiSendFriendRequest,
  apiCancelFriendRequest,
  apiAcceptFriendRequest,
  apiDeclineFriendRequest,
  apiDeleteFriend,
  apiGetFriends
} from '../api/chatApi'
import { eventBus } from '../utils/eventBus'

export const useFriendsStore = defineStore('friends', () => {
  const friendStatusCache = useStorage('friend-status-cache', {})

  function applyFriendRequestUpdate(msg) {
    console.info('🔥 🔥 🔥 Friend request received:', msg)
    friendStatusCache.value[msg.user_id] = msg.status
    eventBus.emit('friend-update', msg)
    eventBus.emit('friend-panel-update', msg)
  }

  async function getFriends(userId) {
    const statusMap = await apiGetFriends(userId)

    Object.keys(friendStatusCache.value).forEach(id => {
      if (!(id in statusMap)) {
        delete friendStatusCache.value[id]
      }
    })

    Object.entries(statusMap).forEach(([id, status]) => {
      friendStatusCache.value[id] = status
    })
    console.info(statusMap)
    return statusMap
  }


  async function sendFriendRequest(fromId, toId) { return apiSendFriendRequest(fromId, toId) }
  async function cancelFriendRequest(fromId, toId) { return apiCancelFriendRequest(fromId, toId) }
  async function acceptFriendRequest(fromId, toId) { return apiAcceptFriendRequest(fromId, toId) }
  async function declineFriendRequest(fromId, toId) { return apiDeclineFriendRequest(fromId, toId) }
  async function deleteFriend(fromId, toId) { return apiDeleteFriend(fromId, toId) }

  return {
    friendStatusCache,
    applyFriendRequestUpdate,
    getFriends,
    sendFriendRequest,
    cancelFriendRequest,
    acceptFriendRequest,
    declineFriendRequest,
    deleteFriend
  }
})