<template>
  <transition name="slide-fade">
    <div
      v-if="Object.keys(callStore.callMembers).length > 0"
      class="fixed top-[80px] left-1/2 z-40 transform -translate-x-1/2 
             w-full max-w-[1300px] bg-gray-50 border-b border-gray-200 
             px-4 py-3 shadow-md flex flex-col gap-4 items-center justify-center"
    >
      <!-- Видео режим -->
      <div
        v-if="anyoneWithCam"
        class="w-full flex flex-col gap-4 items-center"
      >
        <!-- Видео-сетка -->
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
            <video
              v-if="cameraStatusMap[user.id]"
              :ref="el => registerVideoEl(user.id, el)"
              autoplay
              playsinline
              muted
              class="w-full h-48 object-cover"
            ></video>

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
      </div>

      <!-- Режим без камер -->
      <div
        v-else
        class="w-full flex flex-col gap-3 items-center justify-center min-h-[140px]"
      >
        <div class="flex gap-4 flex-wrap justify-center">
          <div
            v-for="user in joinedParticipants"
            :key="'joined-' + user.id"
            class="flex flex-col items-center text-xs"
          >
            <img
              :src="user.avatar_url || '/default-avatar.png'"
              :class="[
                'w-[100px] h-[100px] rounded-full border object-cover mb-1',
                'ring transition-all duration-200',
                mediaStore.speakingUsers.has(Number(user.id))
                  ? 'ring-4 ring-green-500 ring-opacity-80'
                  : 'ring-0 ring-transparent'
              ]"
            />
            <button
              v-if="String(user.id) !== String(currentUserID)"
              @click="openSettings(user)"
              class="hover:underline"
            >
              {{ user.username }}
            </button>
            <span v-else>{{ user.username }} (вы)</span>
          </div>
        </div>
      </div>

      <!-- Централизованные кнопки -->
      <div class="flex gap-2 items-center mt-2 justify-center">
        <button
          @click="mediaStore.toggleMute"
          :title="mediaStore.isMuted ? 'Включить микрофон' : 'Выключить микрофон'"
          class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
        >
          <MicOff v-if="mediaStore.isMuted" class="w-5 h-5 text-red-600" />
          <Mic v-else class="w-5 h-5 text-green-600" />
        </button>

        <button
          @click="mediaStore.toggleCamera"
          :title="mediaStore.isCamOff ? 'Включить камеру' : 'Выключить камеру'"
          class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition"
        >
          <VideoOff v-if="mediaStore.isCamOff" class="w-5 h-5 text-red-600" />
          <Video v-else class="w-5 h-5 text-green-600" />
        </button>

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
  </transition>

  <!-- Аудио -->
  <div style="width:0; height:0; overflow:hidden; position:fixed; pointer-events:none;">
    <template v-for="user in joinedParticipants" :key="'audio-' + user.id">
      <audio
        v-if="String(user.id) !== String(currentUserID)"
        :ref="el => registerAudioEl(user.id, el)"
        autoplay
        playsinline
      ></audio>
    </template>
  </div>

  <!-- DEBUG: Speaking users -->
  <div class="fixed bottom-2 left-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">🎙️ speakingUsers:</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ Array.from(mediaStore.speakingUsers) }}
    </pre>
  </div>

  <!-- DEBUG: Аудио элементы -->
  <div class="fixed bottom-20 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[400px] z-50 overflow-y-auto max-h-[200px]">
    <div class="font-semibold mb-1">🔉 Remote audio elements:</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ debugAudioElements }}
    </pre>
  </div>

</template>

  <!-- DEBUG: Камера статус-мап -->
  <!-- <div class="fixed bottom-2 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">cameraStatusMap:</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ callStore.cameraStatusMap }}
    </pre>
  </div> -->

  <!-- DEBUG: Видео-треки remoteStreams -->
  <!-- <div class="fixed bottom-24 left-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">remoteStreams.getVideoTracks():</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ remoteTracksDebug }}
    </pre>
  </div> -->

  <!-- DEBUG: hasLiveVideo для всех участников -->
  <!-- <div class="fixed bottom-24 right-2 mb-40 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">hasLiveVideo(userId):</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ debugLiveVideoMap }}
    </pre>
  </div> -->

<script setup>
import { computed, ref, watchEffect, onMounted } from 'vue'
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

const localVideoRef = ref(null)

const joinedParticipants = computed(() =>
  Object.entries(callMembers.value)
    .filter(([_, status]) => status === 'joined')
    .map(([id]) => chatStore.users[id] ?? { id, username: 'неизвестно', avatar_url: null })
    .filter(user => user && user.id !== undefined)
)

const participantsWithCam = computed(() =>
  joinedParticipants.value.filter(p => cameraStatusMap.value[String(p.id)] === true)
)

const participantsWithoutCam = computed(() =>
  joinedParticipants.value.filter(p => cameraStatusMap.value[String(p.id)] !== true)
)

const hasJoined = computed(() => callStore.hasJoined)
const callMembersCount = computed(() => Object.keys(callMembers.value).length)

const buttonLabel = computed(() =>
  !hasJoined.value ? 'Присоединиться' : callMembersCount.value > 1 ? 'Выйти' : 'Завершить'
)

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
    map[user.id] = mediaStore.hasLiveVideo(user.id)
  }
  return map
})

const debugAudioElements = computed(() => {
  const result = {}

  for (const [userId, audioEl] of Object.entries(mediaStore.remoteAudioElements)) {
    const stream = audioEl?.srcObject
    const tracks = stream?.getAudioTracks?.() ?? []
    result[userId] = {
      speaking: mediaStore.speakingUsers.has(Number(userId)),
      tracks: tracks.map(track => ({
        id: track.id,
        kind: track.kind,
        label: track.label,
        enabled: track.enabled,
        readyState: track.readyState
      })),
      elReady: !!audioEl,
      streamReady: !!stream
    }
  }

  return result
})

const anyoneWithCam = computed(() => participantsWithCam.value.length > 0)

function handleButtonClick() {
  hasJoined.value ? callStore.leaveCall() : callStore.joinCall()
}

function openSettings(user) {
  chatStore.openUserSettings(user.id)
}

function registerAudioEl(userId, el) {
  if (el) mediaStore.registerAudioElement(userId, el)
}

function registerVideoEl(userId, el) {
  if (el) mediaStore.registerVideoElement(userId, el)
}
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