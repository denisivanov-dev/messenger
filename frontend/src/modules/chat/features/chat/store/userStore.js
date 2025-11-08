import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAllUsers } from '../api/chatApi'

export const useUserStore = defineStore('users', () => {
  const users = ref({})

  async function fetchUsers() {
    const data = await getAllUsers()
    
    const userMap = {}
    for (const user of data) {
      userMap[user.id] = user
    }

    users.value = userMap
    console.info("users:", JSON.stringify(users.value, null, 2))
  }

  function applyStatus(userId, status) {
    const u = users.value[userId]
    if (u) u.status = status
  }

  return {
    users,
    fetchUsers,
    applyStatus,
  }
})