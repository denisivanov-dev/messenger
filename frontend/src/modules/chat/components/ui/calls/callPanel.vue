<template>
  <transition name="slide-fade">
    <div
      v-if="Object.keys(callStore.callMembers).length > 0"
      class="fixed top-[80px] left-1/2 z-40 transform -translate-x-1/2 
             w-full max-w-[1300px] bg-gray-50 border-b border-gray-200 
             px-4 py-3 shadow-md"
    >
      <div
        v-if="anyoneWithCam"
        class="mx-auto w-full max-w-[1300px] flex flex-col gap-4 items-center"
      >
        <!-- Сетка участников -->
        <div
          class="grid gap-3 w-full"
          :class="{
            'grid-cols-1': joinedParticipants.length === 1,
            'grid-cols-2': joinedParticipants.length === 2,
            'grid-cols-3': joinedParticipants.length === 3,
            'grid-cols-4': joinedParticipants.length >= 4
          }"
        >
          <div
            v-for="user in joinedParticipants"
            :key="'cam-' + user.id"
            class="relative flex flex-col items-center justify-center bg-black rounded-lg overflow-hidden"
          >
            <!-- ВИДЕО если камера включена -->
            <video
              v-if="cameraStatusMap[user.id]"
              :ref="el => registerVideoEl(user.id, el)"
              autoplay
              playsinline
              muted
              class="w-full h-48 object-cover"
            ></video>

            <!-- АВАТАРКА если камеры нет -->
            <div
              v-else
              class="flex flex-col items-center justify-center w-full h-48 bg-gray-800 text-white"
            >
              <img
                :src="user.avatar_url || '/default-avatar.png'"
                class="w-16 h-16 rounded-full border border-white object-cover mb-2"
              />
              <span class="text-xs">{{ user.username }}</span>
              <span v-if="String(user.id) === String(currentUserID)"> (вы) </span>
            </div>
          </div>
        </div>

        <!-- кнопки управления -->
        <div class="flex gap-2 items-center mt-2">
          <!-- Мьют -->
          <button
            @click="mediaStore.toggleMute"
            :title="mediaStore.isMuted ? 'Включить микрофон' : 'Выключить микрофон'"
            class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
          >
            <MicOff v-if="mediaStore.isMuted" class="w-5 h-5 text-red-600" />
            <Mic v-else class="w-5 h-5 text-green-600" />
          </button>

          <!-- Камера -->
          <button
            @click="mediaStore.toggleCamera"
            :title="mediaStore.isCamOff ? 'Включить камеру' : 'Выключить камеру'"
            class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
          >
            <VideoOff v-if="mediaStore.isCamOff" class="w-5 h-5 text-red-600" />
            <Video v-else class="w-5 h-5 text-green-600" />
          </button>

          <!-- Выйти/Присоединиться -->
          <button
            @click="handleButtonClick"
            :class="[
              'px-4 py-2 rounded-full flex items-center gap-2 transition text-white',
              hasJoined ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'
            ]"
          >
            <PhoneOff class="w-4 h-4" />
            {{ buttonLabel }}
          </button>
        </div>
      </div>

      <div
        v-else
        class="mx-auto w-full max-w-[1300px] flex justify-between text-sm text-gray-800 font-medium min-h-[130px]"
      >
        <!-- классический режим с аватарками -->
        <div class="flex flex-col gap-2">
          <div class="flex gap-6">
            <!-- JOINED -->
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <UserCheck class="w-4 h-4 text-green-600" />
                <span>Участники:</span>
              </div>
              <div class="flex gap-3 flex-wrap max-w-[250px]">
                <div
                  v-for="user in joinedParticipants"
                  :key="'joined-' + user.id"
                  class="flex flex-col items-center text-xs"
                >
                  <div class="relative">
                    <img
                      :src="user.avatar_url || '/default-avatar.png'"
                      class="w-10 h-10 rounded-full border object-cover"
                    />
                  </div>
                  <button
                    v-if="String(user.id) !== String(currentUserID)"
                    @click="openSettings(user)"
                    class="mt-1 hover:underline"
                  >
                    {{ user.username }}
                  </button>
                  <span v-else class="mt-1">{{ user.username }} (вы)</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Кнопки справа -->
        <div class="flex flex-col justify-end items-end gap-2">
          <div class="flex gap-2 items-center">
            <!-- Мьют -->
            <button
              @click="mediaStore.toggleMute"
              :title="mediaStore.isMuted ? 'Включить микрофон' : 'Выключить микрофон'"
              class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
            >
              <MicOff v-if="mediaStore.isMuted" class="w-5 h-5 text-red-600" />
              <Mic v-else class="w-5 h-5 text-green-600" />
            </button>

            <!-- Камера -->
            <button
              @click="mediaStore.toggleCamera"
              :title="mediaStore.isCamOff ? 'Включить камеру' : 'Выключить камеру'"
              class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
            >
              <VideoOff v-if="mediaStore.isCamOff" class="w-5 h-5 text-red-600" />
              <Video v-else class="w-5 h-5 text-green-600" />
            </button>

            <!-- Выйти/Присоединиться -->
            <button
              @click="handleButtonClick"
              :class="[
                'px-4 py-2 rounded-full flex items-center gap-2 transition text-white',
                hasJoined ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'
              ]"
            >
              <PhoneOff class="w-4 h-4" />
              {{ buttonLabel }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>

  <!-- Аудио теги для всех участников (кроме себя) -->
  <div class="hidden">
    <template v-for="user in joinedParticipants" :key="'audio-' + user.id">
      <audio
        v-if="String(user.id) !== String(currentUserID)"
        :ref="el => registerAudioEl(user.id, el)"
        autoplay
        playsinline
      ></audio>
    </template>
  </div>

  <!-- DEBUG: Камера статус-мап -->
  <div class="fixed bottom-2 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">cameraStatusMap:</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ callStore.cameraStatusMap }}
    </pre>
  </div>

  <!-- DEBUG: Видео-треки localStream -->
  <div class="fixed bottom-2 left-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">localStream.getVideoTracks():</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ localTracksDebug }}
    </pre>
  </div>

  <!-- DEBUG: Видео-треки remoteStreams -->
  <div class="fixed bottom-24 left-2 mb-40 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">remoteStreams.getVideoTracks():</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ remoteTracksDebug }}
    </pre>
  </div>

  <!-- DEBUG: hasLiveVideo для всех участников -->
  <div class="fixed bottom-24 right-2 mb-40 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">hasLiveVideo(userId):</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ debugLiveVideoMap }}
    </pre>
  </div>
</template>

<script setup>
import { computed, ref, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useChatStore } from '../../../store/chatStore'
import { useCallStore } from '../../../store/call/callStore'
import { useMediaStore } from '../../../store/call/mediaStore'
import { PhoneOff, UserCheck, Mic, MicOff, Video, VideoOff } from 'lucide-vue-next'

const authStore = useAuthStore()
const chatStore = useChatStore()
const callStore = useCallStore()
const mediaStore = useMediaStore()

const currentUserID = computed(() => authStore.getUserId)
const { cameraStatusMap, callMembers } = storeToRefs(callStore)

const localTracksDebug = computed(() => {
  const tracks = mediaStore.localStream?.getVideoTracks?.() ?? []
  return tracks.map((t, i) => ({
    id: t.id,
    label: t.label,
    enabled: t.enabled,
    readyState: t.readyState,
    kind: t.kind
  }))
})

const remoteTracksDebug = computed(() => {
  const debug = {}
  for (const [userId, stream] of Object.entries(mediaStore.remoteStreams)) {
    const tracks = stream?.getVideoTracks?.() ?? []
    debug[userId] = tracks.map(track => ({
      id: track.id,
      label: track.label,
      enabled: track.enabled,
      readyState: track.readyState,
      kind: track.kind
    }))
  }
  return debug
})

const debugLiveVideoMap = computed(() => {
  const map = {}
  for (const user of joinedParticipants.value) {
    map[user.id] = userHasLiveVideo(user.id)
  }
  return map
})

const localVideoRef = ref(null)

const joinedParticipants = computed(() => {
  return Object.entries(callStore.callMembers)
    .filter(([id, status]) => id && status === 'joined')
    .map(([id]) => {
      const user = chatStore.users[id]
      return user ?? { id, username: 'неизвестно', avatar_url: null }
    })
    .filter(user => user && user.id !== undefined)
})

const hasJoined = computed(() => callStore.hasJoined)
const callMembersCount = computed(() => Object.keys(callStore.callMembers).length)

const buttonLabel = computed(() => {
  if (!hasJoined.value) return 'Присоединиться'
  return callMembersCount.value > 1 ? 'Выйти' : 'Завершить'
})

function handleButtonClick() {
  if (!hasJoined.value) {
    callStore.joinCall()
  } else {
    callStore.leaveCall()
  }
}

function openSettings(user) {
  chatStore.openUserSettings(user.id)
}

function registerAudioEl(userId, el) {
  if (!el) return
  mediaStore.registerAudioElement(userId, el)
}

function registerVideoEl(userId, el) {
  if (!el) return
  mediaStore.registerVideoElement(userId, el)
}

const userHasLiveVideo = mediaStore.hasLiveVideo

// const userHasLiveVideo = (userId) => {
//   return !!(cameraStatusMap.value && cameraStatusMap.value[String(userId)])
// }

const anyoneWithCam = computed(() => {
  const members = callMembers.value || {}
  const map = cameraStatusMap.value || {}

  return Object.entries(members)
    .some(([id, status]) => status === 'joined' && map[String(id)] === true)
})

watchEffect(() => {
  const videoEl = localVideoRef.value
  const stream = mediaStore.localStream
  if (videoEl && stream) {
    videoEl.srcObject = stream
  }
})
</script>

<style scoped>
.slide-fade-enter-active {
  transition: all 0.3s ease;
}
.slide-fade-leave-active {
  transition: all 0.2s ease;
  opacity: 0;
  transform: translateY(-10px);
}
.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}
</style>