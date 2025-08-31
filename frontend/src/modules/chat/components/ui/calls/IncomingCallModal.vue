<template>
  <transition name="fade">
    <div
      v-if="username"
      class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center"
    >
      <div class="bg-white rounded-2xl shadow-2xl p-6 w-[320px] flex flex-col items-center gap-4 animate-fade-in-up">
        <!-- Иконка звонка -->
        <PhoneIncoming class="w-10 h-10 text-green-600" />

        <!-- Имя пользователя -->
        <div class="text-lg font-semibold text-gray-800 text-center">
          Входящий звонок от<br /><span class="text-blue-600">{{ username }}</span>
        </div>

        <!-- Кнопки -->
        <div class="flex gap-4 mt-2">
          <button
            class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-full flex items-center gap-2 transition"
            @click="accept"
          >
            <Phone class="w-4 h-4" />
            Принять
          </button>
          <button
            class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-full flex items-center gap-2 transition"
            @click="decline"
          >
            <PhoneOff class="w-4 h-4" />
            Отклонить
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { onMounted, onBeforeUnmount } from 'vue'
import { PhoneIncoming, Phone, PhoneOff } from 'lucide-vue-next'

defineProps({ username: String })
const emit = defineEmits(['accept', 'decline'])

let ringtone = null

function playRingtone() {
  ringtone = new Audio('/sounds/ring_call.mp3')
  ringtone.volume = 0.6
  ringtone.play().catch(console.error)
}

function stopRingtone() {
  if (ringtone instanceof Audio) {
    ringtone.pause()
    ringtone.currentTime = 0
    ringtone = null
  }
}

function accept() {
  stopRingtone()
  emit('accept')
}

function decline() {
  stopRingtone()
  emit('decline')
}

onMounted(() => {
  playRingtone()
})

onBeforeUnmount(() => {
  stopRingtone()
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fadeInUp 0.3s ease-out;
}
</style>