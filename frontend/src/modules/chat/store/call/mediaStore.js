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
  const remoteStreams = ref({})
  const remoteAudioElements = ref({})
  const remoteVideoElements = ref({})

  const micSettings = ref({
    enabled: true,
    deviceId: null,
  })

  const camSettings = ref({
    enabled: false,
    deviceId: null,
  })

  const isMuted = computed(() => !micSettings.value.enabled)
  const isCamOff = computed(() => !camSettings.value.enabled)

  const micKeyForUser = () => `call:mic:${String(authStore.getUserId)}`
  const camKeyForUser = () => `call:cam:${String(authStore.getUserId)}`

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

  function applyMicStateToLocalStream() {
    if (!localStream.value) return
    const shouldEnable = micSettings.value.enabled
    localStream.value.getAudioTracks().forEach(track => {
      track.enabled = shouldEnable
    })
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

    // Удаляем старые треки
    localStream.value.getVideoTracks().forEach(t => {
      localStream.value.removeTrack(t)
      t.stop()
    })

    // Добавляем новые
    newStream.getVideoTracks().forEach(track => {
      track.enabled = true // ✅ всегда true
      localStream.value.addTrack(track)
      replaceVideoTrack(track)
    })

    applyCamStateToLocalStream() // ✅ применяем видимость
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
      audioEl.srcObject = stream
      audioEl.play().catch(() => {})
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
      el.srcObject = stream
      el.play().catch(() => {})
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
    isMuted,
    isCamOff,

    toggleMute,
    toggleCamera,
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
    hasLiveVideo
  }
})