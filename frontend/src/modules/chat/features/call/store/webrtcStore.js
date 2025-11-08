import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMediaStore } from './mediaStore'
import { useCallStore } from './callStore'
import { useAuthStore } from '../../../../auth/store/authStore'
import { sendSocketPayload } from '../../api/chatApi'

export const useWebRTCStore = defineStore('webrtc', () => {
  const authStore = useAuthStore()
  const mediaStore = useMediaStore()
  const callStore = useCallStore()

  const peerConnections = ref({})
  const pendingCandidates = ref({})

  const makingOfferMap = ref({})
  const politeMap = ref({})

  const servers = {
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' },
      {
        urls: 'turn:46.109.9.240:3478',
        username: 'test',
        credential: '1234'
      }
    ]
  }

  async function createPeerConnection(userId, isCaller) {
    if (peerConnections.value[userId]) {
      try { peerConnections.value[userId].close() } catch {}
      delete peerConnections.value[userId]
    }

    const pc = new RTCPeerConnection({ iceServers: servers.iceServers })
    peerConnections.value[userId] = pc

    politeMap.value[userId] = String(authStore.getUserId) > String(userId)
    makingOfferMap.value[userId] = false

    // === audio ===
    const audioTransceiver = pc.addTransceiver("audio", { direction: "sendrecv" })
    mediaStore.micSenderMap.set(userId, audioTransceiver.sender)
    if (mediaStore.localMicTrack) {
      await audioTransceiver.sender.replaceTrack(mediaStore.localMicTrack)
    }

    // === camera ===
    const camTransceiver = pc.addTransceiver("video", { direction: "sendrecv" })
    mediaStore.camSenderMap.set(userId, camTransceiver.sender)
    if (mediaStore.localCameraTrack) {
      await camTransceiver.sender.replaceTrack(mediaStore.localCameraTrack)
    }

    // === screen ===
    const screenTransceiver = pc.addTransceiver("video", { direction: "sendrecv" })
    mediaStore.screenSenderMap.set(userId, screenTransceiver.sender)
    if (mediaStore.localScreenTrack) {
      mediaStore.localScreenTrack.contentHint = "detail"
      await screenTransceiver.sender.replaceTrack(mediaStore.localScreenTrack)
    }

    // === ontrack ===
    pc.ontrack = event => {
      if (!callStore.hasJoined) return
      if (!peerConnections.value[userId]) return 
      
      const { track } = event
      if (track.readyState === "ended") return

      const settings = track.getSettings?.() || {}
      const displaySurface = (settings.displaySurface || "").toLowerCase()
      const isScreen =
        track.kind === "video" &&
        (track.contentHint === "detail" ||
          ["monitor", "window", "browser", "application"].includes(displaySurface))

      if (track.kind === "audio") {
        mediaStore.remoteAudioStreams[userId] ??= new MediaStream()

        mediaStore.remoteAudioStreams[userId].getTracks().forEach(t => {
          if (t.readyState === 'ended') {
            mediaStore.remoteAudioStreams[userId].removeTrack(t)
          }
        })

        mediaStore.remoteAudioStreams[userId].addTrack(track)
        console.log(`[ontrack] 🔊 audio track от ${userId}, stream now has ${mediaStore.remoteAudioStreams[userId].getAudioTracks().length} audio tracks`)

        mediaStore.startVoiceDetection(userId, mediaStore.remoteAudioStreams[userId], true)
      } else if (track.kind === "video") {
        const alreadyHasCam = mediaStore.remoteCameraStreams[userId]?.getVideoTracks().length > 0

        if (isScreen || alreadyHasCam) {
          mediaStore.remoteScreenStreams[userId] ??= new MediaStream()

          mediaStore.remoteScreenStreams[userId].getTracks().forEach(t => {
            if (t.readyState === 'ended') {
              mediaStore.remoteScreenStreams[userId].removeTrack(t)
            }
          })

          mediaStore.remoteScreenStreams[userId].addTrack(track)
        } else {
          mediaStore.remoteCameraStreams[userId] ??= new MediaStream()

          mediaStore.remoteCameraStreams[userId].getTracks().forEach(t => {
            if (t.readyState === 'ended') {
              mediaStore.remoteCameraStreams[userId].removeTrack(t)
            }
          })

          mediaStore.remoteCameraStreams[userId].addTrack(track)
        }
      }
    }

    pc.onicecandidate = e => {
      if (e.candidate) {
        sendSocketPayload({
          type: 'ice_candidate',
          chat_type: 'private',
          receiver_id: userId,
          candidate: e.candidate
        })
      }
    }

    pc.onnegotiationneeded = async () => {
      try {
        makingOfferMap.value[userId] = true
        const offer = await pc.createOffer()
        await pc.setLocalDescription(offer)
        sendSocketPayload({
          type: "webrtc_offer",
          chat_type: "private",
          receiver_id: userId,
          offer: pc.localDescription
        })
      } finally {
        makingOfferMap.value[userId] = false
      }
    }

    console.log('[createPeerConnection]', userId, {
      localMicTrack: !!mediaStore.localMicTrack,
      localCameraTrack: !!mediaStore.localCameraTrack,
      localScreenTrack: !!mediaStore.localScreenTrack
    })

    if (isCaller) {
      const offer = await pc.createOffer()
      await pc.setLocalDescription(offer)
      sendSocketPayload({
        type: "webrtc_offer",
        chat_type: "private",
        receiver_id: userId,
        offer: pc.localDescription
      })
    }
  }

  async function startCall(userId, isCaller) {
    await mediaStore.initMediaTracks()
    await createPeerConnection(userId, isCaller)
  }

  async function handleOffer(from, offer) {
    if (!callStore.hasJoined) return
    
    await mediaStore.initMediaTracks()
    if (!peerConnections.value[from]) {
      await createPeerConnection(from, false)
    }

    const pc = peerConnections.value[from]
    const desc = new RTCSessionDescription(offer)

    const makingOffer = makingOfferMap.value[from]
    const isPolite = politeMap.value[from]
    const offerCollision = (makingOffer || pc.signalingState !== "stable")

    if (!isPolite && offerCollision) {
      console.warn(`[webrtc] Игнорируем оффер от ${from} (collision)`)
      return
    }

    await pc.setRemoteDescription(desc)
    const answer = await pc.createAnswer()
    await pc.setLocalDescription(answer)

    sendSocketPayload({
      type: "webrtc_answer",
      chat_type: "private",
      receiver_id: from,
      answer: pc.localDescription
    })

    flushPendingCandidates(from)
  }

  async function handleAnswer(from, answer) {
    if (!callStore.hasJoined) return
      
    const pc = peerConnections.value[from]
    if (!pc) return

    try {
      if (pc.signalingState === "have-local-offer") {
        await pc.setRemoteDescription(new RTCSessionDescription(answer))
        flushPendingCandidates(from)
      } else if (pc.signalingState === "stable") {
        console.warn(`[webrtc] Игнорируем дубликат answer от ${from} (state = stable)`)
      } else {
        console.warn(`[webrtc] Неподдерживаемое состояние при answer: ${pc.signalingState}`)
      }
    } catch (err) {
      console.warn(`[webrtc] Ошибка handleAnswer от ${from}:`, err)
    }
  }

  async function handleIceCandidate(from, candidate) {
    if (!callStore.hasJoined) return
    
    const pc = peerConnections.value[from]

    if (!pc || !pc.remoteDescription) {
      pendingCandidates.value[from] ??= []
      pendingCandidates.value[from].push(candidate)
      return
    }
    try {
      await pc.addIceCandidate(candidate)
    } catch (err) {
      console.warn(`[webrtc] Ошибка addIceCandidate от ${from}:`, err)
    }
  }

  async function flushPendingCandidates(userId) {
    const pc = peerConnections.value[userId]
    const candidates = pendingCandidates.value[userId] || []
    for (const cand of candidates) {
      try { await pc.addIceCandidate(cand) } catch {}
    }
    delete pendingCandidates.value[userId]
  }

  async function renegotiateWithAll() {
    for (const [userId, pc] of Object.entries(peerConnections.value)) {
      try {
        const offer = await pc.createOffer({
          offerToReceiveAudio: true,
          offerToReceiveVideo: true
        })
        await pc.setLocalDescription(offer)

        sendSocketPayload({
          type: 'webrtc_offer',
          chat_type: 'private',
          receiver_id: userId,
          offer: pc.localDescription
        })
      } catch (err) {
        console.warn(`[webrtc] Не удалось выполнить renegotiation для ${userId}:`, err)
      }
    }
  }

  function endCallWith(userId) {
    const pc = peerConnections.value[userId]

    if (pc) {
      try {
        pc.close()
      } catch {}
      delete peerConnections.value[userId]
    }

    delete mediaStore.remoteAudioStreams[userId]
    delete mediaStore.remoteCameraStreams[userId]
    delete mediaStore.remoteScreenStreams[userId]

    mediaStore.stopVoiceDetection(userId)

    mediaStore.micSenderMap.delete(userId)
    mediaStore.camSenderMap.delete(userId)
    mediaStore.screenSenderMap.delete(userId)

    delete pendingCandidates.value[userId]
    delete makingOfferMap.value[userId]
    delete politeMap.value[userId]
  }


  function endAllCalls() {
    Object.keys(peerConnections.value).forEach(userId => {
      endCallWith(userId)
    })
  }

  function clearLocalTracks() {
    if (mediaStore.localMicTrack) {
      mediaStore.localMicTrack.stop()
      mediaStore.localMicTrack = null
    }
    if (mediaStore.localCameraTrack) {
      mediaStore.localCameraTrack.stop()
      mediaStore.localCameraTrack = null
    }
    if (mediaStore.localScreenTrack) {
      mediaStore.localScreenTrack.stop()
      mediaStore.localScreenTrack = null
    }
  }

  return {
    peerConnections,
    startCall,
    handleOffer,
    handleAnswer,
    handleIceCandidate,
    renegotiateWithAll,
    endCallWith,
    endAllCalls,
    clearLocalTracks
  }
})