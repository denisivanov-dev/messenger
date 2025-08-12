import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAllUsers } from '../api/chatApi'

export const useUserStore = defineStore('users', () => {
  const users = ref([])

  async function fetchUsers() {
    const data = await getAllUsers()
    users.value = data
    console.info(JSON.stringify(users.value, null, 2))
  }

  function applyStatus(userId, status) {
    const u = users.value.find(x => x.id === userId)
    if (u) u.status = status
  }

  return {
    users,
    fetchUsers,
    applyStatus,
  }
})
