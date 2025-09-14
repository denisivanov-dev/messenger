import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '../../../auth/store/authStore'
import { useWebRTCStore } from './webrtcStore'
import { useCallStore } from './callStore'

export const useMediaStore = defineStore('media', () => {
  const authStore = useAuthStore()
  // Инициализируем webrtcStore лениво, чтобы избежать циклических зависимостей
  let webrtcStore
  import('./webrtcStore').then(module => {
    webrtcStore = module.useWebRTCStore()
  })
  
  const callStore = useCallStore()

  // === local streams ===
  const localStream = ref(null)
  const localScreenStream = ref(null)

  // === remote streams ===
  const remoteAudioStreams = ref({})
  const remoteCameraStreams = ref({})
  const remoteScreenStreams = ref({})

  // === settings ===
  const micSettings = ref({
    enabled: true,
    deviceId: null,
  })

  const camSettings = ref({
    enabled: false,
    deviceId: null,
  })

  const screenSettings = ref({
    enabled: false,
  })

  // === states ===
  const isMuted = computed(() => !micSettings.value.enabled)
  const isCamOff = computed(() => !camSettings.value.enabled)
  const speakingUsers = ref(new Set())

  // === VAD ===
  const vadContexts = {}
  const lastSpeakingMap = {}
  const speakingDebounceMs = 300

  // === Redis/localStorage keys ===
  const micKeyForUser = () => `call:mic:${String(authStore.getUserId)}`
  const camKeyForUser = () => `call:cam:${String(authStore.getUserId)}`
  const screenKeyForUser = () => `call:screen:${String(authStore.getUserId)}`
  
  // === local tracks ===
  const localMicTrack = ref(null)
  const localCameraTrack = ref(null)
  const localScreenTrack = ref(null)

  // === sender maps ===
  const micSenderMap = new Map()
  const camSenderMap = new Map()
  const screenSenderMap = new Map()

  function saveMicSettings() {
    try {
      localStorage.setItem(micKeyForUser(), JSON.stringify(micSettings.value))
    } catch {}
  }

  function loadMicSettings() {
    try {
      const raw = localStorage.getItem(micKeyForUser())
      if (raw) {
        const parsed = JSON.parse(raw)
        micSettings.value = {
          enabled: typeof parsed.enabled === 'boolean' ? parsed.enabled : true,
          deviceId: parsed.deviceId ?? null,
        }
      }
    } catch {}
  }

  function saveCamSettings() {
    try {
      localStorage.setItem(camKeyForUser(), JSON.stringify(camSettings.value))
    } catch {}
  }

  function loadCamSettings() {
    try {
      const raw = localStorage.getItem(camKeyForUser())
      if (raw) {
        const parsed = JSON.parse(raw)
        camSettings.value = {
          enabled: typeof parsed.enabled === 'boolean' ? parsed.enabled : false,
          deviceId: parsed.deviceId ?? null,
        }
      }
    } catch {}
  }

  function saveScreenSettings() {
    try {
      localStorage.setItem(screenKeyForUser(), JSON.stringify(screenSettings.value))
    } catch {}
  }

  function loadScreenSettings() {
    try {
      const raw = localStorage.getItem(screenKeyForUser())
      if (raw) {
        const parsed = JSON.parse(raw)
        screenSettings.value = {
          enabled: typeof parsed.enabled === 'boolean' ? parsed.enabled : false
        }
      }
    } catch {}
  }

  function startVoiceDetection(userId, stream, isRemote = false, el = null) {
    if (!stream) {
      console.warn(`❌ [VAD] Нет stream для userId=${userId}`)
      return
    }

    if (isRemote && el) {
      if (el.__vadInitialized) return
      el.__vadInitialized = true
    }

    console.log('✅ [VAD] Запуск для userId:', userId)

    if (!vadContexts[userId]) {
      vadContexts[userId] = new AudioContext()
    }

    const audioCtx = vadContexts[userId]
    audioCtx.resume().catch(err => {
      console.warn('⚠️ [VAD] Ошибка resume AudioContext:', err)
    })

    const source = audioCtx.createMediaStreamSource(stream)
    const analyser = audioCtx.createAnalyser()
    analyser.fftSize = 2048
    const data = new Uint8Array(analyser.fftSize)

    source.connect(analyser)

    const detect = () => {
      analyser.getByteTimeDomainData(data)

      let sum = 0
      for (let i = 0; i < data.length; i++) {
        const val = (data[i] - 128) / 128
        sum += val * val
      }

      const volume = Math.sqrt(sum / data.length)
      const isSpeaking = volume > 0.03
      const idNum = Number(userId)
      const now = Date.now()

      const last = lastSpeakingMap[idNum] || { state: null, changed: 0 }

      if (isSpeaking !== last.state && now - last.changed > speakingDebounceMs) {
        lastSpeakingMap[idNum] = { state: isSpeaking, changed: now }

        if (isSpeaking) {
          speakingUsers.value.add(idNum)
        } else {
          speakingUsers.value.delete(idNum)
        }
      }

      requestAnimationFrame(detect)
    }

    detect()
  }
  
  function applyMicStateToLocalStream() {
    const track = localStream.value?.getAudioTracks?.()[0]
    if (track) {
      track.enabled = micSettings.value.enabled
    }
  }

  function applyCamStateToLocalStream() {
    const track = localStream.value?.getVideoTracks?.()[0]
    if (track) {
      track.enabled = camSettings.value.enabled
    }
  }

  function toggleMute() {
    micSettings.value.enabled = !micSettings.value.enabled
    saveMicSettings()
    applyMicStateToLocalStream()
  }

  async function toggleCamera() {
    if (camSettings.value.enabled) {
      await stopCamera()
    } else {
      await startCamera()
    }
  }

  async function toggleScreenShare() {
    if (screenSettings.value.enabled) {
      await stopScreenShare()
    } else {
      await startScreenShare()
    }
  }
  
  async function startCamera() {
    camSettings.value.enabled = true
    saveCamSettings()

    await setCamDevice(camSettings.value.deviceId)

    const track = localCameraTrack.value
    if (!track) return

    // ИЗМЕНЕНИЕ: ВОЗВРАЩАЕМ ЛОГИКУ ДЛЯ ОТОБРАЖЕНИЯ СВОЕЙ КАМЕРЫ
    remoteCameraStreams.value[authStore.getUserId] = new MediaStream([track])

    let needRenegotiate = false

    for (const [userId, sender] of camSenderMap.entries()) {
      if (!sender.track) {
        needRenegotiate = true
      }
      await sender.replaceTrack(track)
    }

    if (needRenegotiate && webrtcStore) {
      webrtcStore.renegotiateWithAll()
    }

    callStore.sendCameraStatusUpdate(true)
  }

  async function stopCamera() {
    camSettings.value.enabled = false
    saveCamSettings()

    const track = localCameraTrack.value
    if (track) {
      track.stop()
    }
    localCameraTrack.value = null

    // ИЗМЕНЕНИЕ: УДАЛЯЕМ СВОЙ ПОТОК ИЗ СПИСКА
    delete remoteCameraStreams.value[authStore.getUserId]

    for (const sender of camSenderMap.values()) {
      await sender.replaceTrack(null)
    }

    callStore.sendCameraStatusUpdate(false)
  }

  async function startScreenShare() {
    try {
      const stream = await navigator.mediaDevices.getDisplayMedia({
        video: {
          frameRate: 30,
          width: { ideal: 1280 },
          height: { ideal: 720 },
        },
        audio: false,
      })

      const track = stream.getVideoTracks()[0]
      if (!track) {
        console.warn('❌ Нет видеотрека в screen stream')
        return
      }

      localScreenStream.value = stream
      localScreenTrack.value = track
      screenSettings.value.enabled = true
      saveScreenSettings()

      // ИЗМЕНЕНИЕ: ВОЗВРАЩАЕМ ЛОГИКУ ДЛЯ ОТОБРАЖЕНИЯ СВОЕЙ ДЕМКИ
      remoteScreenStreams.value[authStore.getUserId] = stream

      let needRenegotiate = false

      for (const [userId, sender] of screenSenderMap.entries()) {
        if (!sender.track) {
          needRenegotiate = true
        }
        await sender.replaceTrack(track)
      }

      if (needRenegotiate && webrtcStore) {
        webrtcStore.renegotiateWithAll()
      }

      callStore.sendScreenStatusUpdate(true)

      track.onended = () => {
        stopScreenShare()
      }
    } catch (err) {
      console.warn('❌ Ошибка при старте демонстрации экрана:', err)
    }
  }

  async function stopScreenShare() {
    const track = localScreenTrack.value
    if (track) {
      track.stop()
    }
    localScreenStream.value = null
    localScreenTrack.value = null

    screenSettings.value.enabled = false
    saveScreenSettings()
    callStore.sendScreenStatusUpdate(false)

    // ИЗМЕНЕНИЕ: УДАЛЯЕМ СВОЙ ПОТОК ИЗ СПИСКА
    delete remoteScreenStreams.value[authStore.getUserId]

    for (const sender of screenSenderMap.values()) {
      await sender.replaceTrack(null)
    }
  }
  
  async function setMicDevice(deviceId) {
    micSettings.value.deviceId = deviceId
    saveMicSettings()

    const constraints = {
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true,
        channelCount: 2,
        sampleRate: 48000,
        sampleSize: 16,
        ...(deviceId ? { deviceId: { exact: deviceId } } : {})
      }
    }

    const newStream = await navigator.mediaDevices.getUserMedia(constraints)
    const newTrack = newStream.getAudioTracks()[0]
    if (!newTrack) {
      console.warn('❌ Не удалось получить аудиотрек')
      return
    }

    newTrack.enabled = micSettings.value.enabled
    localMicTrack.value = newTrack

    localStream.value ??= new MediaStream()
    localStream.value.getAudioTracks().forEach(t => {
      localStream.value.removeTrack(t)
      t.stop()
    })
    localStream.value.addTrack(newTrack)

    for (const sender of micSenderMap.values()) {
      sender.replaceTrack(newTrack)
    }
  }

  async function setCamDevice(deviceId) {
    camSettings.value.deviceId = deviceId
    saveCamSettings()

    const constraints = {
      video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        frameRate: { ideal: 30 },
        ...(deviceId ? { deviceId: { exact: deviceId } } : {})
      }
    }

    const newStream = await navigator.mediaDevices.getUserMedia(constraints)
    const newTrack = newStream.getVideoTracks()[0]
    if (!newTrack) {
      console.warn('❌ Не удалось получить видеотрек')
      return
    }

    newTrack.enabled = camSettings.value.enabled
    localCameraTrack.value = newTrack

    localStream.value ??= new MediaStream()
    localStream.value.getVideoTracks().forEach(t => {
      localStream.value.removeTrack(t)
      t.stop()
    })
    localStream.value.addTrack(newTrack)

    for (const sender of camSenderMap.values()) {
      sender.replaceTrack(newTrack)
    }
  }

  function detachRemoteAudioStream(userId) {
    const stream = remoteAudioStreams.value[userId]
    if (stream) {
      stream.getTracks().forEach(t => t.stop())
      delete remoteAudioStreams.value[userId]
    }

    delete vadContexts[userId]
    delete lastSpeakingMap[userId]
    speakingUsers.value.delete(Number(userId))
  }

  function detachRemoteCameraStream(userId) {
    const stream = remoteCameraStreams.value[userId]
    if (stream) {
      stream.getTracks().forEach(t => t.stop())
      delete remoteCameraStreams.value[userId]
    }
  }

  function detachRemoteScreenStream(userId) {
    const stream = remoteScreenStreams.value[userId]
    if (stream) {
      stream.getTracks().forEach(t => t.stop())
      delete remoteScreenStreams.value[userId]
    }
  }

  function detachAllRemoteStreams(userId) {
    detachRemoteAudioStream(userId)
    detachRemoteCameraStream(userId)
    detachRemoteScreenStream(userId)
  }

  function cleanupAllRemoteStreams() {
    const userIds = new Set([
      ...Object.keys(remoteAudioStreams.value),
      ...Object.keys(remoteCameraStreams.value),
      ...Object.keys(remoteScreenStreams.value),
    ])

    for (const userId of userIds) {
      detachAllRemoteStreams(userId)
    }
  }

  async function initMediaTracks() {
    if (!localMicTrack.value) {
      await setMicDevice(micSettings.value.deviceId)

      for (const [userId, sender] of micSenderMap.entries()) {
        if (localMicTrack.value) {
          await sender.replaceTrack(localMicTrack.value)
          console.log(`[initMediaTracks] 🔄 Микрофон заменён у ${userId}`)
        }
      }
    }

    if (localStream.value) {
      startVoiceDetection(authStore.getUserId, localStream.value, false)
    }
  }

  function hasLiveVideo(userId) {
    const checkTrack = (track) =>
      track?.enabled && track.readyState === 'live'

    if (String(userId) === String(authStore.getUserId)) {
      return checkTrack(localCameraTrack.value) || checkTrack(localScreenTrack.value)
    }

    const checkTracks = (stream) =>
      stream?.getVideoTracks?.().some(checkTrack) ?? false

    return (
      checkTracks(remoteCameraStreams.value[userId]) ||
      checkTracks(remoteScreenStreams.value[userId])
    )
  }

  return {
    localStream,
    localScreenStream,
    remoteAudioStreams,
    remoteCameraStreams,
    remoteScreenStreams,
    micSettings,
    camSettings,
    screenSettings,
    isMuted,
    isCamOff,
    speakingUsers,
    toggleMute,
    toggleCamera,
    toggleScreenShare,
    setMicDevice,
    setCamDevice,
    applyMicStateToLocalStream,
    applyCamStateToLocalStream,
    detachRemoteAudioStream,
    detachRemoteCameraStream,
    detachRemoteScreenStream,
    detachAllRemoteStreams,
    startScreenShare,
    stopScreenShare,
    startCamera,
    stopCamera,
    startVoiceDetection,
    initMediaTracks,
    hasLiveVideo,
    cleanupAllRemoteStreams,
    saveMicSettings,
    saveCamSettings,
    saveScreenSettings,
    loadMicSettings,
    loadCamSettings,
    loadScreenSettings,
    localMicTrack,
    localCameraTrack,
    localScreenTrack,
    micSenderMap,
    camSenderMap,
    screenSenderMap
  }
})