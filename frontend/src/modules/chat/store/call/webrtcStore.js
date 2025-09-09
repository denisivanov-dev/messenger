import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMediaStore } from './mediaStore'
import { useCallStore } from './callStore'
import { sendMessage } from '../../api/chatApi'

export const useWebRTCStore = defineStore('webrtc', () => {
  const mediaStore = useMediaStore()
  const callStore = useCallStore()
  
  const peerConnections = ref({})
  const pendingCandidates = ref({})

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
      try {
        peerConnections.value[userId].close()
      } catch (e) {}
      delete peerConnections.value[userId]
    }

    const pc = new RTCPeerConnection({ iceServers: servers.iceServers })
    peerConnections.value[userId] = pc

    pc.ontrack = (event) => {
      const { track } = event
      const stream = new MediaStream([track])

      const kind = track.kind
      const settings = typeof track.getSettings === 'function' ? track.getSettings() : {}
      const displaySurface = (settings.displaySurface || '').toLowerCase()
      const isScreen =
        kind === 'video' &&
        (track.contentHint === 'detail' ||
        ['monitor', 'window', 'browser', 'application'].includes(displaySurface))

      // 💡 fallback: если уже есть камера, а это второй video-трек — считаем его экраном
      const existingCamera = mediaStore.remoteCameraStreams[userId]
      const isSecondVideo = kind === 'video' && existingCamera

      const finalIsScreen = isScreen || isSecondVideo

      console.log('📥 Новый track:', {
        userId,
        kind,
        label: track.label,
        displaySurface,
        contentHint: track.contentHint,
        isScreen,
        isSecondVideo,
        finalIsScreen,
      })

      if (
        kind === 'video' &&
        finalIsScreen &&
        callStore.screenStatusMap[userId] === false
      ) {
        console.warn(`⛔ Получен экран ${track.label}, но он выключен — игнорим`)
        track.stop?.()
        return
      }

      if (kind === 'audio') {
        mediaStore.remoteAudioStreams[userId] = stream
        mediaStore.attachRemoteAudioStream(userId, stream)
      } else if (kind === 'video') {
        if (finalIsScreen) {
          mediaStore.remoteScreenStreams[userId] = stream
          mediaStore.attachRemoteScreenStream(userId, stream)
        } else {
          mediaStore.remoteCameraStreams[userId] = stream
          mediaStore.attachRemoteCameraStream(userId, stream)
        }
      }
    }

    pc.onicecandidate = e => {
      if (e.candidate) {
        sendMessage({
          type: 'ice_candidate',
          chat_type: 'private',
          receiver_id: userId,
          candidate: e.candidate
        })
      }
    }

    mediaStore.localStream?.getAudioTracks()?.forEach(track => {
      pc.addTrack(track, mediaStore.localStream)
    })

    mediaStore.localStream?.getVideoTracks()?.forEach(track => {
      pc.addTrack(track, mediaStore.localStream)
    })

    if (mediaStore.localScreenStream) {
      const screenTrack = mediaStore.localScreenStream.getVideoTracks()[0]
      if (screenTrack) {
        pc.addTrack(screenTrack, mediaStore.localScreenStream)
      }
    }

    if (isCaller) {
      const offer = await pc.createOffer()
      await pc.setLocalDescription(offer)

      sendMessage({
        type: 'webrtc_offer',
        chat_type: 'private',
        receiver_id: userId,
        offer
      })
    }
  }

  async function startCall(userId, isCaller) {
    await mediaStore.initMediaTracks()

    await createPeerConnection(userId, isCaller)
  }

  async function handleOffer(from, offer) {
    await mediaStore.initMediaTracks()

    await createPeerConnection(from, false)
    const pc = peerConnections.value[from]

    await pc.setRemoteDescription(new RTCSessionDescription(offer))

    const answer = await pc.createAnswer()
    await pc.setLocalDescription(answer)

    sendMessage({
      type: 'webrtc_answer',
      chat_type: 'private',
      receiver_id: from,
      answer
    })

    flushPendingCandidates(from)
  }

  async function handleAnswer(from, answer) {
    const pc = peerConnections.value[from]
    if (!pc) return

    await pc.setRemoteDescription(new RTCSessionDescription(answer))
    flushPendingCandidates(from)
  }

  async function handleIceCandidate(from, candidate) {
    const pc = peerConnections.value[from]

    if (!pc || !pc.remoteDescription) {
      if (!pendingCandidates.value[from]) {
        pendingCandidates.value[from] = []
      }
      pendingCandidates.value[from].push(candidate)
      return
    }

    try {
      await pc.addIceCandidate(candidate)
    } catch (err) {}
  }

  async function flushPendingCandidates(userId) {
    const pc = peerConnections.value[userId]
    const candidates = pendingCandidates.value[userId] || []

    for (const cand of candidates) {
      try {
        await pc.addIceCandidate(cand)
      } catch (err) {}
    }

    delete pendingCandidates.value[userId]
  }

  function endCallWith(userId) {
    const pc = peerConnections.value[userId]
    if (pc) {
      pc.close()
      delete peerConnections.value[userId]
    }

    delete pendingCandidates.value[userId]
    mediaStore.detachRemoteStream(userId)

    if (Object.keys(peerConnections.value).length === 0) {
      mediaStore.localStream?.getTracks().forEach(track => track.stop())
      mediaStore.localStream = null
    }
  }

  function endAllCalls() {
    Object.keys(peerConnections.value).forEach(endCallWith)

    Object.keys(pendingCandidates.value).forEach(key => {
      delete pendingCandidates.value[key]
    })

    mediaStore.cleanupAllRemoteStreams()
  }

  function renegotiateWithAll() {
    Object.entries(peerConnections.value).forEach(async ([userId, pc]) => {
      try {
        const offer = await pc.createOffer()
        await pc.setLocalDescription(offer)

        sendMessage({
          type: 'webrtc_offer',
          chat_type: 'private',
          receiver_id: userId,
          offer: {
            type: pc.localDescription.type,
            sdp: pc.localDescription.sdp
          }
        })
      } catch (err) {
        console.warn(`[webrtc] Не удалось выполнить renegotiation для ${userId}:`, err)
      }
    })
  }

  return {
    peerConnections,
    startCall,
    handleOffer,
    handleAnswer,
    handleIceCandidate,
    endCallWith,
    endAllCalls,
    renegotiateWithAll
  }
})