import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '../../../auth/store/authStore'
import { useWebRTCStore } from './webrtcStore'
import { useCallStore } from './callStore'

export const useMediaStore = defineStore('media', () => {
  const authStore = useAuthStore()
  const webrtcStore = useWebRTCStore()
  const callStore = useCallStore()

  const localStream = ref(null)
  const localScreenStream = ref(null)
  const remoteStreams = ref({})
  const remoteAudioElements = ref({})
  const remoteVideoElements = ref({})

  const speakingUsers = ref(new Set())
  const vadContexts = {}
  const lastSpeakingMap = {}
  const speakingDebounceMs = 300

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

  const isMuted = computed(() => !micSettings.value.enabled)
  const isCamOff = computed(() => !camSettings.value.enabled)

  const micKeyForUser = () => `call:mic:${String(authStore.getUserId)}`
  const camKeyForUser = () => `call:cam:${String(authStore.getUserId)}`
  const screenKeyForUser = () => `call:screen:${String(authStore.getUserId)}`

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

  function saveMicSettings() {
    try {
      localStorage.setItem(micKeyForUser(), JSON.stringify(micSettings.value))
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

  function saveCamSettings() {
    try {
      localStorage.setItem(camKeyForUser(), JSON.stringify(camSettings.value))
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

  function saveScreenSettings() {
    try {
      localStorage.setItem(screenKeyForUser(), JSON.stringify(screenSettings.value))
    } catch {}
  }

  function applyMicStateToLocalStream() {
    if (!localStream.value) return
    const shouldEnable = micSettings.value.enabled
    localStream.value.getAudioTracks().forEach(track => {
      track.enabled = shouldEnable
    })
  }

  function startLocalVoiceDetection(userId) {
    if (!localStream.value) {
      console.warn('❌ Нет localStream, не запускаем VAD')
      return
    }

    console.log('✅ VAD запущен для local user', userId)

    if (!vadContexts[userId]) {
      vadContexts[userId] = new AudioContext()
    }

    const audioCtx = vadContexts[userId]
    audioCtx.resume().catch(err => {
      console.warn('⚠️ [VAD] Ошибка resume AudioContext:', err)
    })

    const source = audioCtx.createMediaStreamSource(localStream.value)
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

  function startVoiceDetectionForUser(userId, audioElement) {
    if (!audioElement || audioElement.__vadInitialized) return
    audioElement.__vadInitialized = true

    console.log('✅ [VAD] Запуск для userId:', userId)

    const stream = audioElement.srcObject
    if (!stream) {
      console.warn(`❌ [VAD] Нет srcObject для audioElement (userId=${userId})`)
      return
    }

    const tracks = stream.getAudioTracks?.() || []
    if (tracks.length === 0) {
      console.warn(`❌ [VAD] Нет аудиотреков в srcObject для userId=${userId}`)
      return
    }

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

  function applyCamStateToLocalStream() {
    if (!localStream.value) return
    const shouldEnable = camSettings.value.enabled
    localStream.value.getVideoTracks().forEach(track => {
      track.enabled = shouldEnable
    })
  }

  function toggleMute() {
    micSettings.value.enabled = !micSettings.value.enabled
    saveMicSettings()
    applyMicStateToLocalStream()
  }

  async function toggleCamera() {
    camSettings.value.enabled = !camSettings.value.enabled
    saveCamSettings()

    if (camSettings.value.enabled) {
      if (!localStream.value || !localStream.value.getVideoTracks().length) {
        await setCamDevice(camSettings.value.deviceId)
      }
      const videoTrack = localStream.value.getVideoTracks()[0]
      if (videoTrack) replaceVideoTrack(videoTrack)
    }
    applyCamStateToLocalStream()

    callStore.sendCameraStatusUpdate(camSettings.value.enabled)
  }

  async function toggleScreenShare() {
    if (screenSettings.value.enabled) {
      await stopScreenShare()
    } else {
      await startScreenShare()
    }
  }

  async function startScreenShare() {
    try {
      console.log('📺 Запрашиваем экран для демонстрации…')
      const stream = await navigator.mediaDevices.getDisplayMedia({
        video: {
          frameRate: 30,
          width: { ideal: 1280 },
          height: { ideal: 720 },
        },
        audio: false,
      })

      console.log('✅ Экран получен:', stream)

      localScreenStream.value = stream
      screenSettings.value.enabled = true
      saveScreenSettings()

      const screenTrack = stream.getVideoTracks()[0]
      if (!screenTrack) {
        console.warn('❌ Нет видеотрека в screen stream')
        return
      }

      // 💥 Проверяем, что webrtcStore и peerConnections доступны
      if (!webrtcStore || !webrtcStore.peerConnections?.value) {
        console.warn('❌ webrtcStore или peerConnections не определены')
      }

      console.log('📦 webrtcStore', webrtcStore)
      console.log(localScreenStream)
      console.log(toggleScreenShare)
      console.log(webrtcStore?.peerConnections?.value)

      const pcs = Object.values(webrtcStore.peerConnections.value)
      console.log(`📡 Подключено peer-соединений: ${pcs.length}`)

      pcs.forEach((pc, index) => {
        try {
          pc.addTrack(screenTrack, stream)
          console.log(`🧩 Добавлен screenTrack в peerConnection #${index}`)
        } catch (err) {
          console.warn(`⚠️ Не удалось добавить трек в pc #${index}`, err)
        }
      })

      screenTrack.onended = () => {
        console.log('📴 Демонстрация экрана завершена пользователем')
        stopScreenShare()
      }

      callStore.sendScreenStatusUpdate(true)
      console.log('📨 Отправлен статус "демонстрация экрана включена"')

    } catch (err) {
      console.warn('❌ Ошибка при старте демонстрации экрана:', err)
    }
  }

  async function stopScreenShare() {
    if (!localScreenStream.value) return

    localScreenStream.value.getTracks().forEach(t => t.stop())
    localScreenStream.value = null

    screenSettings.value.enabled = false
    saveScreenSettings()
    callStore.sendScreenStatusUpdate(false)
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
    localStream.value ??= new MediaStream()

    localStream.value.getAudioTracks().forEach(t => {
      localStream.value.removeTrack(t)
      t.stop()
    })

    newStream.getAudioTracks().forEach(track => {
      track.enabled = micSettings.value.enabled
      localStream.value.addTrack(track)
    })
  }

  async function setCamDevice(deviceId) {
    camSettings.value.deviceId = deviceId
    camSettings.value.enabled = true
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
    localStream.value ??= new MediaStream()

    localStream.value.getVideoTracks().forEach(t => {
      localStream.value.removeTrack(t)
      t.stop()
    })

    newStream.getVideoTracks().forEach(track => {
      track.enabled = true
      localStream.value.addTrack(track)
      replaceVideoTrack(track)
    })

    applyCamStateToLocalStream()
  }

  function replaceVideoTrack(newTrack) {
    if (!newTrack) return
    if (!webrtcStore?.peerConnections?.value) return

    Object.values(webrtcStore.peerConnections.value).forEach(pc => {
      const senders = pc.getSenders().filter(s => s.track && s.track.kind === 'video')

      if (senders.length > 0) {
        senders[0].replaceTrack(newTrack)
      } else {
        if (localStream.value) {
          pc.addTrack(newTrack, localStream.value)
        }
      }
    })
  }

  function attachRemoteStream(userId, stream) {
    remoteStreams.value[userId] = stream

    const audioEl = remoteAudioElements.value[userId]
    if (audioEl) {
      const isNew = audioEl.srcObject !== stream
      if (isNew) {
        audioEl.srcObject = stream
        audioEl.__vadInitialized = false
      }
      audioEl.play().catch(() => {})
      if (!audioEl.__vadInitialized) {
        setTimeout(() => {
          startVoiceDetectionForUser(userId, audioEl)
        }, 300)
      }
    }

    const videoEl = remoteVideoElements.value[userId]
    if (videoEl) {
      videoEl.srcObject = stream
      videoEl.play().catch(() => {})
    }
  }

  function detachRemoteStream(userId) {
    remoteStreams.value[userId]?.getTracks().forEach(t => t.stop())
    delete remoteStreams.value[userId]
    const audioEl = remoteAudioElements.value[userId]
    if (audioEl) {
      audioEl.pause()
      audioEl.srcObject = null
    }
    delete remoteAudioElements.value[userId]
    const videoEl = remoteVideoElements.value[userId]
    if (videoEl) {
      videoEl.pause()
      videoEl.srcObject = null
    }
    delete remoteVideoElements.value[userId]
  }

  function registerAudioElement(userId, el) {
    if (!el) return
    remoteAudioElements.value[userId] = el

    const stream = remoteStreams.value[userId]

    if (stream) {
      const isNew = el.srcObject !== stream
      if (isNew) {
        el.srcObject = stream
        el.__vadInitialized = false
      }
      el.play().catch(() => {})
      if (!el.__vadInitialized) {
        setTimeout(() => {
          startVoiceDetectionForUser(userId, el)
        }, 300)
      }
    } else {
      const interval = setInterval(() => {
        const stream = remoteStreams.value[userId]
        if (stream) {
          clearInterval(interval)
          el.srcObject = stream
          el.__vadInitialized = false
          el.play().catch(() => {})
          setTimeout(() => {
            startVoiceDetectionForUser(userId, el)
          }, 300)
        }
      }, 200)
    }
  }

  function registerVideoElement(userId, el) {
    if (!el) return
    remoteVideoElements.value[userId] = el
    const stream = remoteStreams.value[userId]
    if (stream) {
      el.srcObject = stream
      el.play().catch(() => {})
    }
  }

  function cleanupAllRemoteStreams() {
    for (const userId in remoteStreams.value) {
      detachRemoteStream(userId)
    }
  }

  async function initMediaTracks() {
    await setMicDevice(micSettings.value.deviceId)
    startLocalVoiceDetection(authStore.getUserId)

    const deviceId = camSettings.value.deviceId
    const constraints = {
      video: {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        frameRate: { ideal: 30 },
        ...(deviceId ? { deviceId: { exact: deviceId } } : {})
      }
    }

    try {
      const newStream = await navigator.mediaDevices.getUserMedia({ video: constraints.video })
      localStream.value ??= new MediaStream()

      localStream.value.getVideoTracks().forEach(t => {
        localStream.value.removeTrack(t)
        t.stop()
      })

      newStream.getVideoTracks().forEach(track => {
        track.enabled = true
        localStream.value.addTrack(track)
        replaceVideoTrack(track)
      })

      applyCamStateToLocalStream()

      remoteStreams.value[authStore.getUserId] = localStream.value
    } catch (err) {
      console.warn('[initMediaTracks] Камера не получена:', err)
    }
  }

  function hasLiveVideo(userId) {
    const checkTracks = (stream) =>
      stream?.getVideoTracks?.().some(track => track.enabled && track.readyState === 'live') ?? false

    if (String(userId) === String(authStore.getUserId)) {
      return checkTracks(localStream.value)
    }
    return checkTracks(remoteStreams.value[userId])
  }

  loadMicSettings()
  loadCamSettings()

  return {
    localStream,
    remoteStreams,
    remoteAudioElements,
    remoteVideoElements,
    micSettings,
    camSettings,
    screenSettings,
    isMuted,
    isCamOff,
    speakingUsers,
    localScreenStream,

    toggleMute,
    toggleCamera,
    toggleScreenShare,
    setMicDevice,
    setCamDevice,
    applyMicStateToLocalStream,
    applyCamStateToLocalStream,

    attachRemoteStream,
    detachRemoteStream,
    registerAudioElement,
    registerVideoElement,
    cleanupAllRemoteStreams,
    initMediaTracks,
    hasLiveVideo,

    startLocalVoiceDetection,
    startVoiceDetectionForUser,
    startScreenShare,
    stopScreenShare,

    loadCamSettings,
    loadMicSettings,
    loadScreenSettings
  }
})