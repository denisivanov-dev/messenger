<template>
  <!-- Верхняя панель -->
  <div
    class="call-bar w-full max-w-[1300px] h-[50px] px-5 py-2 bg-gray-100 rounded-2xl shadow-md mb-3 flex items-center justify-between"
  >
    <!-- Слева: имя собеседника + статус звонка -->
    <div class="flex items-center gap-4 text-sm text-gray-700 font-medium">
      <div>Собеседник: {{ username }}</div>

      <!-- Статус вызова -->
      <div v-if="isCalling" class="flex items-center gap-1 text-yellow-600">
        <PhoneIcon class="w-4 h-4" />
        <span class="flex items-center gap-1">
          Вызов<span class="dots"><span>.</span><span>.</span><span>.</span></span>
        </span>
      </div>

      <div v-else-if="inCall" class="flex items-center gap-1 text-green-600">
        <PhoneIcon class="w-4 h-4" />
        <span>Звонок в процессе</span>
      </div>
    </div>

    <!-- Кнопка звонка -->
    <div class="relative group">
      <button
        @click="startCall"
        :disabled="inCall || isCalling"
        class="p-2 rounded-full flex items-center justify-center transition
               text-white bg-green-600 hover:bg-green-700
               disabled:bg-gray-400 disabled:cursor-not-allowed"
      >
        <PhoneIcon class="w-4 h-4" />
      </button>
      <div
        class="absolute bottom-full mb-1 left-1/2 -translate-x-1/2 px-2 py-1 text-xs text-white bg-black rounded opacity-0 group-hover:opacity-100 transition pointer-events-none"
      >
        Позвонить
      </div>
    </div>
  </div>

  <!-- Панель звонка -->
  <CallPanel v-if="hasCallMembers" />

  <!-- Модалка входящего звонка -->
  <IncomingCallModal
    v-if="incomingUsername"
    :username="incomingUsername"
    @accept="acceptCall"
    @decline="declineCall"
  />
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '../../../store/chatStore'
import { useCallStore } from '../../../store/call/callStore'
import { useAuthStore } from '../../../../auth/store/authStore'
import { PhoneIcon } from 'lucide-vue-next'
import IncomingCallModal from '../calls/IncomingCallModal.vue'
import CallPanel from '../calls/callPanel.vue'

const chatStore = useChatStore()
const callStore = useCallStore()
const authStore = useAuthStore()

const currentUserID = computed(() => authStore.getUserId)

const username = computed(() => {
  const id = chatStore.receiverID
  return Object.values(chatStore.users).find(u => String(u.id) === String(id))?.username || 'неизвестно'
})

const inCall = computed(() => callStore.inCall)
const isCalling = computed(() => callStore.isCalling)

const hasCallMembers = computed(() => Object.keys(callStore.callMembers).length > 0)

const incomingFromID = computed(() => callStore.incomingCallFrom)
const incomingUsername = computed(() => {
  if (!incomingFromID.value) return null
  return Object.values(chatStore.users).find(
    u => String(u.id) === String(incomingFromID.value)
  )?.username || 'неизвестно'
})

async function startCall() {
  const myId = String(currentUserID.value)
  const inRoom = Object.prototype.hasOwnProperty.call(callStore.callMembers, myId)

  if (Object.keys(callStore.callMembers).length > 0) {
    if (!inRoom) {
      await callStore.joinCall()
    }
  } else {
    await callStore.startRequestCall()
  }
}

function acceptCall() {
  if (incomingFromID.value) callStore.acceptCall(incomingFromID.value)
}

function declineCall() {
  if (incomingFromID.value) callStore.declineCall(incomingFromID.value)
}
</script>

<style scoped>
.dots span {
  animation: blink 1.4s infinite;
  font-weight: bold;
}
.dots span:nth-child(2) {
  animation-delay: 0.2s;
}
.dots span:nth-child(3) {
  animation-delay: 0.4s;
}
@keyframes blink {
  0% { opacity: 0.2; }
  20% { opacity: 1; }
  100% { opacity: 0.2; }
}
</style>