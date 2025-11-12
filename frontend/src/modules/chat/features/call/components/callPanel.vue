<template>
  <transition name="slide-fade">
    <div
      v-if="Object.keys(callStore.callMembers).length > 0"
      class="fixed top-[80px] left-1/2 z-40 transform -translate-x-1/2 
             w-full max-w-[1300px] bg-gray-50 border-b border-gray-200 
             px-4 py-3 shadow-md flex flex-col gap-4 items-center justify-center"
    >
      <div
        v-if="anyoneWithCam || anyoneWithScreen"
        class="w-full flex flex-col gap-4 items-center"
      >
        <transition-group
          name="participant"
          tag="div"
          class="grid gap-3 w-full"
          :class="{
            'grid-cols-1 place-items-center': effectiveCount === 1,
            'grid-cols-2': effectiveCount === 2,
            'grid-cols-3': effectiveCount === 3,
            'grid-cols-4': effectiveCount >= 4
          }"
        >
          <div
            v-for="user in joinedParticipants"
            :key="`tile-${String(user.id)}`"
            class="relative w-full h-48 bg-black rounded-lg overflow-hidden flex items-center justify-center"
            :class="[
              mediaStore.speakingUsers.has(Number(user.id))
                ? 'ring-4 ring-green-500 ring-opacity-80'
                : 'ring-0 ring-transparent',
              joinedParticipants.length === 1 ? 'max-w-[400px] mx-auto' : 'w-full'
            ]"
          >
            <!-- Видео экран -->
            <video
              v-if="hasJoined && callStore.screenStatusMap[user.id]"
              :ref="el => registerScreenEl(user.id, el)"
              autoplay
              playsinline
              muted
              class="absolute inset-0 w-full h-full object-cover cursor-pointer"
              @click="openFullscreen(user.id)"
            ></video>

            <!-- Видео камера -->
            <video
              v-else-if="hasJoined && cameraStatusMap[user.id]"
              :ref="el => registerCameraEl(user.id, el)"
              autoplay
              playsinline
              muted
              class="absolute inset-0 w-full h-full object-cover cursor-pointer"
              @click="openFullscreen(user.id)"
            ></video>

            <!-- Заглушка -->
            <div
              v-else
              class="absolute inset-0 flex flex-col items-center justify-center bg-gray-800 text-white"
            >
              <img
                :src="user.avatar_url || '/default-avatar.png'"
                class="w-16 h-16 rounded-full border border-white object-cover mb-2"
              />
              <span class="text-xs mb-2">{{ user.username }}</span>
              <span v-if="String(user.id) === String(currentUserID)" class="text-xs text-gray-400">(вы)</span>

              <!-- Иконки состояния для тех, кто не в звонке -->
              <div v-if="!hasJoined" class="flex gap-2 mt-2">
                <div
                  v-if="cameraStatusMap[user.id]"
                  class="flex items-center gap-1 text-xs bg-gray-700 px-2 py-1 rounded"
                >
                  <Video class="w-3 h-3" />
                  <span>Камера</span>
                </div>
                <div
                  v-if="callStore.screenStatusMap[user.id]"
                  class="flex items-center gap-1 text-xs bg-gray-700 px-2 py-1 rounded"
                >
                  <Monitor class="w-3 h-3" />
                  <span>Экран</span>
                </div>
                <div
                  v-if="callStore.micStatusMap[user.id] === false"
                  class="flex items-center gap-1 text-xs bg-gray-700 px-2 py-1 rounded"
                >
                  <MicOff class="w-3 h-3" />
                  <span>Микрофон выключен</span>
                </div>
              </div>
            </div>

            <!-- Лоадер -->
            <div
              v-if="videoLoadingMap[user.id]"
              class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-60"
            >
              <div class="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin"></div>
            </div>
          </div>
        </transition-group>
      </div>

      <div v-else class="w-full flex flex-col gap-3 items-center justify-center min-h-[140px]">
        <transition-group
          name="participant"
          tag="div"
          class="flex gap-4 flex-wrap justify-center"
        >
          <div
            v-for="user in joinedParticipants"
            :key="'joined-' + user.id"
            class="flex flex-col items-center text-xs relative"
          >
            <div class="relative w-[100px] h-[100px]">
              <img
                :src="user.avatar_url || '/default-avatar.png'"
                class="w-full h-full rounded-full border object-cover ring transition-all duration-200"
                :class="mediaStore.speakingUsers.has(Number(user.id))
                  ? 'ring-4 ring-green-500 ring-opacity-80'
                  : 'ring-0 ring-transparent'"
              />

              <!-- 🔊 Микрофон поверх аватарки -->
              <div class="absolute bottom-1 right-1 bg-white rounded-full p-[2px] shadow-md">
                <Mic 
                  v-if="callStore.micStatusMap[user.id]" 
                  class="w-4 h-4 text-green-500" 
                  title="Микрофон включён" 
                />
                <MicOff 
                  v-else 
                  class="w-4 h-4 text-red-500" 
                  title="Микрофон выключен" 
                />
              </div>
            </div>

            <button
              v-if="String(user.id) !== String(currentUserID)"
              @click="openSettings(user)"
              class="hover:underline mt-1"
            >
              {{ user.username }}
            </button>
            <span v-else class="mt-1">{{ user.username }} (вы)</span>
          </div>
        </transition-group>
      </div>

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
          @click="hasJoined ? mediaStore.toggleScreenShare() : null"
          :disabled="!hasJoined"
          :title="!hasJoined 
            ? 'Вы не в звонке' 
            : (mediaStore.screenSettings.enabled 
                ? 'Остановить демонстрацию экрана' 
                : 'Начать демонстрацию экрана')"
          class="p-2 rounded-full bg-gray-200 hover:bg-gray-300 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <MonitorOff 
            v-if="mediaStore.screenSettings.enabled" 
            class="w-5 h-5 text-red-600" 
          />
          <Monitor 
            v-else 
            class="w-5 h-5 text-blue-600" 
          />
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

  <div
    v-if="fullscreenVideoUserId && (fullscreenScreenStream || fullscreenCameraStream)"
    class="fixed inset-0 z-[100] bg-black flex items-center justify-center"
    @click.self="closeFullscreen"
  >
    <!-- Главное видео -->
    <div class="relative w-full h-full flex items-center justify-center">
      <video
        ref="fullscreenVideoRef"
        autoplay
        playsinline
        class="w-full h-full object-contain"
      ></video>
      <div
        v-if="videoLoadingMap[fullscreenVideoUserId]"
        class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-60"
      >
        <div class="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin"></div>
      </div>
    </div>

    <!-- Мини-окно (только если есть и демка, и камера) -->
    <div
      v-if="fullscreenScreenStream && fullscreenCameraStream"
      class="absolute w-48 h-28 bg-black rounded-lg shadow-lg cursor-move overflow-hidden pip-container"
      :class="[pipClass, 
              pipUser && mediaStore.speakingUsers.has(Number(pipUser.id))
                ? 'ring-4 ring-green-500 ring-opacity-80'
                : 'ring-0 ring-transparent']"
      :style="{ transform: `translate(${pipOffset.x}px, ${pipOffset.y}px)` }"
      @mousedown="startDrag"
      @click.stop="onPipClick"
    >
      <div class="relative w-full h-full">
        <video
          ref="pipVideoRef"
          autoplay
          playsinline
          muted
          class="w-full h-full object-cover"
        ></video>
        <!-- Лоадер -->
        <div
          v-if="videoLoadingMap[fullscreenVideoUserId]"
          class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-60"
        >
          <div class="w-8 h-8 border-4 border-white border-t-transparent rounded-full animate-spin"></div>
        </div>
      </div>
    </div>

    <button
      @click="closeFullscreen"
      class="absolute top-4 right-4 text-white bg-black bg-opacity-60 hover:bg-opacity-80 rounded-full p-2"
      title="Закрыть"
    >
      ✕
    </button>
  </div>

  <div
    v-if="callStore.showCamPopup"
    class="fixed inset-0 z-[200] flex items-center justify-center bg-black bg-opacity-50"
  >
    <div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full text-center">
      <p class="text-lg mb-4">У вас была включена камера в прошлом звонке.<br />Включить её снова?</p>
      <div class="flex justify-center gap-4">
        <button
          @click="enableCameraWithDelay"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
        >
          Да
        </button>
        <button
          @click="() => callStore.showCamPopup = false"
          class="px-4 py-2 bg-gray-300 rounded hover:bg-gray-400 transition"
        >
          Нет
        </button>
      </div>
    </div>
  </div>

  <!-- <div class="fixed bottom-2 left-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">🎙️ speakingUsers:</div>
    <pre class="whitespace-pre-wrap break-words">{{ Array.from(mediaStore.speakingUsers) }}</pre>
  </div>
  <div class="fixed bottom-20 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[400px] z-50 overflow-y-auto max-h-[200px]">
    <div class="font-semibold mb-1">🔉 Remote audio elements:</div>
    <pre class="whitespace-pre-wrap break-words">{{ debugAudioElements }}</pre>
  </div>
  <div class="fixed bottom-5 left-1/2 transform -translate-x-1/2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[400px] z-50">
    <div class="font-semibold mb-1">Все video-треки:</div>
    <pre class="whitespace-pre-wrap break-words">{{ allVideoTracksDebug }}</pre>
  </div>
  <div class="fixed bottom-24 left-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">remoteStreams.getVideoTracks():</div>
    <pre class="whitespace-pre-wrap break-words">{{ remoteTracksDebug }}</pre>
  </div>
  <div class="fixed bottom-24 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[400px] z-50">
    <div class="font-semibold mb-1">🎧 remoteAudioStreams.getAudioTracks():</div>
    <pre class="whitespace-pre-wrap break-words">{{ audioTracksDebug }}</pre>
  </div> -->
  <!-- <div class="fixed bottom-24 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">remoteStreams.getVideoTracks():</div>
    <pre class="whitespace-pre-wrap break-words">{{ remoteScreenTracksDebug }}</pre>
  </div>
  <div class="fixed bottom-2 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">screenStatusMap:</div>
    <pre class="whitespace-pre-wrap break-words">{{ callStore.screenStatusMap }}</pre>
  </div>
  <div class="fixed bottom-2 right-2 bg-white border border-gray-300 shadow-lg rounded-lg p-2 text-xs text-black max-w-[300px] z-50">
    <div class="font-semibold mb-1">cameraStatusMap:</div>
    <pre class="whitespace-pre-wrap break-words">
      {{ callStore.callMembers }}
    </pre>
  </div> -->
</template>

<script setup>
import { computed, ref, watchEffect, onMounted, onBeforeUnmount, watch} from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '../../../../auth/store/authStore'
import { useChatStore } from '../../chat/store/chatStore'
import { useCallStore } from '../store/callStore'
import { useMediaStore } from '../store/mediaStore'
import { useWebRTCStore } from '../store/webrtcStore'
import { PhoneOff, UserCheck, Mic, MicOff, Video, VideoOff, Monitor, MonitorOff } from 'lucide-vue-next'

const authStore = useAuthStore()
const chatStore = useChatStore()
const callStore = useCallStore()
const webrtcStore = useWebRTCStore()
const mediaStore = useMediaStore()

const currentUserID = computed(() => authStore.getUserId)
const { cameraStatusMap, callMembers } = storeToRefs(callStore)

// ИЗМЕНЕНИЕ: Локальные refs для хранения DOM-элементов
const cameraElements = ref({})
const screenElements = ref({})
const audioElements = ref({})

const userPreferredView = ref({})
const effectiveCount = ref(0)

const joinedParticipants = computed(() =>
  Object.entries(callMembers.value)
    .filter(([_, status]) => status === 'joined')
    .map(([id]) => chatStore.users[id] ?? { id, username: 'неизвестно', avatar_url: null })
    .filter(user => user && user.id !== undefined)
)

watch(joinedParticipants, (newVal = [], oldVal = []) => {
  if (newVal.length < oldVal.length) {
    setTimeout(() => {
      effectiveCount.value = newVal.length
    }, 300)
  } else {
    effectiveCount.value = newVal.length
  }
}, { immediate: true })

const participantsWithCam = computed(() =>
  joinedParticipants.value.filter(p => cameraStatusMap.value[String(p.id)] === true)
)

const participantsWithoutCam = computed(() =>
  joinedParticipants.value.filter(p => cameraStatusMap.value[String(p.id)] !== true)
)

const participantsWithScreen = computed(() =>
  joinedParticipants.value.filter(p => callStore.screenStatusMap[String(p.id)])
)

const participantsWithoutScreen = computed(() =>
  joinedParticipants.value.filter(p => !callStore.screenStatusMap[String(p.id)])
)

const hasJoined = computed(() => callStore.hasJoined)
const callMembersCount = computed(() => Object.keys(callMembers.value).length)

const buttonLabel = computed(() =>
  !hasJoined.value ? 'Присоединиться' : callMembersCount.value > 1 ? 'Выйти' : 'Завершить'
)

// Ваши debug-блоки оставлены без изменений
const remoteTracksDebug = computed(() => {
  const debug = {}
  for (const [userId, stream] of Object.entries(mediaStore.remoteCameraStreams)) {
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

const audioTracksDebug = computed(() => {
  const result = {}
  for (const [userId, stream] of Object.entries(mediaStore.remoteAudioStreams)) {
    result[userId] = stream?.getAudioTracks().map(t => ({
      id: t.id,
      readyState: t.readyState,
      enabled: t.enabled,
      label: t.label
    }))
  }
  return result
})

const remoteScreenTracksDebug = computed(() => {
  const debug = {}
  for (const [userId, stream] of Object.entries(mediaStore.remoteScreenStreams)) {
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

const allVideoTracksDebug = computed(() => {
  const result = {}

  Object.entries(mediaStore.remoteCameraStreams).forEach(([userId, stream]) => {
    result[userId] = result[userId] || []
    stream.getVideoTracks().forEach(track => result[userId].push({
      id: track.id,
      label: track.label,
      enabled: track.enabled,
      readyState: track.readyState,
      kind: track.kind
    }))
  })

  Object.entries(mediaStore.remoteScreenStreams).forEach(([userId, stream]) => {
    result[userId] = result[userId] || []
    stream.getVideoTracks().forEach(track => result[userId].push({
      id: track.id,
      label: track.label,
      enabled: track.enabled,
      readyState: track.readyState,
      kind: track.kind
    }))
  })

  return JSON.stringify(result, null, 2)
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

  for (const [userId, audioEl] of Object.entries(audioElements.value)) { // ИЗМЕНЕНИЕ: Используем локальный ref
    if (!audioEl) continue
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

const anyoneWithScreen = computed(() => {
  return Object.values(callStore.screenStatusMap).some(v => v === true)
})

function handleButtonClick() {
  hasJoined.value ? callStore.leaveCall() : callStore.joinCall()
}

function openSettings(user) {
  chatStore.openUserSettings(user.id)
}

// ИЗМЕНЕНИЕ: Функции регистрации теперь обновляют локальные refs
function registerAudioEl(userId, el) {
  audioElements.value[userId] = el
}

function registerCameraEl(userId, el) {
  cameraElements.value[userId] = el
}

function registerScreenEl(userId, el) {
  screenElements.value[userId] = el
}

// ИЗМЕНЕНИЕ: Главное исправление - реактивное связывание потоков
const videoLoadingMap = ref({})

// ИЗМЕНЕНИЕ: Главное исправление - реактивное связывание потоков
watchEffect(() => {
  if (!joinedParticipants.value) return

  for (const user of joinedParticipants.value) {
    const userId = String(user.id)
    const isSelf = String(userId) === String(authStore.user.id)

    // --- Камера ---
    const camEl = cameraElements.value[userId]
    const camStream = isSelf
      ? (mediaStore.localCameraTrack ? new MediaStream([mediaStore.localCameraTrack]) : null)
      : mediaStore.remoteCameraStreams[userId] || null
    if (camEl) {
      if (camStream && camEl.srcObject !== camStream) {
        camEl.srcObject = camStream
        videoLoadingMap.value[userId] = true
        camEl.onloadeddata = () => {
          videoLoadingMap.value[userId] = false
        }
        camEl.play?.().catch(() => {})
      } else if (!camStream && camEl.srcObject) {
        camEl.srcObject = null
        videoLoadingMap.value[userId] = false
      }
    }

    // --- Экран ---
    const screenEl = screenElements.value[userId]
    const screenStream = mediaStore.remoteScreenStreams[userId]
    if (screenEl) {
      if (screenStream && screenEl.srcObject !== screenStream) {
        screenEl.srcObject = screenStream
        videoLoadingMap.value[userId] = true
        screenEl.onloadeddata = () => {
          videoLoadingMap.value[userId] = false
        }
        screenEl.play?.().catch(() => {})
      } else if (!screenStream && screenEl.srcObject) {
        screenEl.srcObject = null
        videoLoadingMap.value[userId] = false
      }
    }

    // --- Аудио ---
    const audioEl = audioElements.value[userId]
    const audioStream = mediaStore.remoteAudioStreams[userId]
    if (audioEl) {
      if (audioStream && audioEl.srcObject !== audioStream) {
        audioEl.srcObject = audioStream
        audioEl.play?.().catch(() => {})
      } else if (!audioStream && audioEl.srcObject) {
        audioEl.srcObject = null
      }
    }
  }
})

// --- состояние fullscreen ---
const fullscreenVideoUserId = ref(null)
const fullscreenPrimary = ref("screen") // "screen" | "camera"

// --- функции открытия/закрытия ---
function openFullscreen(userId) {
  fullscreenVideoUserId.value = userId

  const hasScreen = !!callStore.screenStatusMap[userId]
  const hasCamera = !!callStore.cameraStatusMap[userId]

  if (hasScreen && hasCamera) {
    fullscreenPrimary.value = "screen"
  } else if (hasScreen) {
    fullscreenPrimary.value = "screen"
  } else if (hasCamera) {
    fullscreenPrimary.value = "camera"
  } else {
    fullscreenPrimary.value = null
  }
}

function closeFullscreen() {
  fullscreenVideoUserId.value = null
}

// --- отдельные стримы ---
const fullscreenScreenStream = computed(() => {
  const uid = fullscreenVideoUserId.value
  if (!uid) return null
  if (!callStore.screenStatusMap[uid]) return null   // 🔥 если статус false → убираем сразу
  return mediaStore.remoteScreenStreams[uid] || null
})

const fullscreenCameraStream = computed(() => {
  const uid = fullscreenVideoUserId.value
  if (!uid) return null
  if (!callStore.cameraStatusMap[uid]) return null   // 🔥 если статус false → убираем сразу
  return mediaStore.remoteCameraStreams[uid] || null
})

// --- главный поток ---
const fullscreenStream = computed(() => {
  if (!fullscreenVideoUserId.value) return null
  if (fullscreenPrimary.value === "screen" && fullscreenScreenStream.value) {
    return fullscreenScreenStream.value
  }
  if (fullscreenPrimary.value === "camera" && fullscreenCameraStream.value) {
    return fullscreenCameraStream.value
  }
  return fullscreenScreenStream.value || fullscreenCameraStream.value
})

// --- картинка-в-картинке ---
const pipStream = computed(() => {
  if (!fullscreenVideoUserId.value) return null

  const hasScreen = !!fullscreenScreenStream.value
  const hasCamera = !!fullscreenCameraStream.value

  if (hasScreen && hasCamera) {
    // только если оба трека есть → тогда второй идёт в PIP
    return fullscreenPrimary.value === "screen"
      ? fullscreenCameraStream.value
      : fullscreenScreenStream.value
  }

  return null
})

// --- refs для видео ---
const fullscreenVideoRef = ref(null)
const pipVideoRef = ref(null)


watchEffect(() => {
  // --- Fullscreen ---
  const fsEl = fullscreenVideoRef.value
  if (fsEl) {
    if (fullscreenStream.value) {
      if (fsEl.srcObject !== fullscreenStream.value) {
        fsEl.srcObject = fullscreenStream.value
        videoLoadingMap.value[fullscreenVideoUserId.value] = true

        fsEl.onloadeddata = () => {
          videoLoadingMap.value[fullscreenVideoUserId.value] = false
        }

        fsEl.play?.().catch(() => {})
      }
    } else {
      fsEl.srcObject = null
      videoLoadingMap.value[fullscreenVideoUserId.value] = false
    }
  }

  // --- PIP ---
  const pipEl = pipVideoRef.value
  if (pipEl) {
    if (pipStream.value) {
      if (pipEl.srcObject !== pipStream.value) {
        pipEl.srcObject = pipStream.value
        pipEl.play?.().catch(() => {})
      }
    } else {
      pipEl.srcObject = null
    }
  }
})

// --- переключение местами ---
function toggleFullscreenView() {
  fullscreenPrimary.value =
    fullscreenPrimary.value === "screen" ? "camera" : "screen"
}

// --- ESC закрывает fullscreen ---
function handleKeydown(e) {
  if (e.key === "Escape") {
    closeFullscreen()
  }
}
onMounted(() => {
  window.addEventListener("keydown", handleKeydown)
})
onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleKeydown)
})

// --- позиция pip ---
const pipPosition = ref("top-left")

const pipClass = computed(() => {
  return {
    "top-4 left-4": pipPosition.value === "top-left",
    "top-4 right-4": pipPosition.value === "top-right",
    "bottom-4 left-4": pipPosition.value === "bottom-left",
    "bottom-4 right-4": pipPosition.value === "bottom-right",
  }
})

const pipUser = computed(() => {
  return fullscreenVideoUserId.value
    ? chatStore.users[fullscreenVideoUserId.value]
    : null
})

const dragState = ref({ dragging: false, moved: false, startX: 0, startY: 0 })
const pipOffset = ref({ x: 0, y: 0 })
const pipDragging = ref(false)

function startDrag(e) {
  dragState.value.dragging = true
  dragState.value.moved = false
  pipDragging.value = true
  dragState.value.startX = e.clientX
  dragState.value.startY = e.clientY
  pipOffset.value = { x: 0, y: 0 }
  window.addEventListener("mousemove", onDrag)
  window.addEventListener("mouseup", stopDrag)
}

function onDrag(e) {
  if (!dragState.value.dragging) return
  const dx = e.clientX - dragState.value.startX
  const dy = e.clientY - dragState.value.startY
  pipOffset.value = { x: dx, y: dy }
  if (Math.abs(dx) > 5 || Math.abs(dy) > 5) {
    dragState.value.moved = true
  }
}

function stopDrag(e) {
  dragState.value.dragging = false
  pipDragging.value = false
  window.removeEventListener("mousemove", onDrag)
  window.removeEventListener("mouseup", stopDrag)

  if (dragState.value.moved) {
    pipOffset.value = { x: 0, y: 0 }

    const w = window.innerWidth
    const h = window.innerHeight
    const left = e.clientX < w / 2
    const top = e.clientY < h / 2
    if (top && left) pipPosition.value = "top-left"
    else if (top && !left) pipPosition.value = "top-right"
    else if (!top && left) pipPosition.value = "bottom-left"
    else pipPosition.value = "bottom-right"
  }
}

function onPipClick() {
  if (!dragState.value.moved) {
    toggleFullscreenView()
  }
}

function enableCameraWithDelay() {
  callStore.showCamPopup = false
  mediaStore.toggleCamera().then(() => {
    setTimeout(() => {
      webrtcStore.renegotiateWithAll()
    }, 200)
  })
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

.pip-video {
  transition: all 0.2s ease-in-out;
}

:deep(.participant-enter-active) {
  transition: all 0.35s ease;
}
:deep(.participant-enter-from) {
  opacity: 0;
  transform: scale(0.8);
}
:deep(.participant-enter-to) {
  opacity: 1;
  transform: scale(1);
}

:deep(.participant-leave-active) {
  transition: all 0.25s ease;
}
:deep(.participant-leave-from) {
  opacity: 1;
  transform: scale(1);
}
:deep(.participant-leave-to) {
  opacity: 0;
  transform: scale(0.8);
}
</style>