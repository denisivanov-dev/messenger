import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useMediaStore } from './mediaStore'
import { sendMessage } from '../../api/chatApi'

export const useWebRTCStore = defineStore('webrtc', () => {
  const mediaStore = useMediaStore()

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

    pc.ontrack = event => {
      const remoteStream = event.streams[0]
      mediaStore.attachRemoteStream(userId, remoteStream)
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
    // await mediaStore.setMicDevice(mediaStore.micSettings.deviceId)

    // if (mediaStore.camSettings.enabled) {
    //   await mediaStore.setCamDevice(mediaStore.camSettings.deviceId)
    // }
    await mediaStore.initMediaTracks()

    await createPeerConnection(userId, isCaller)
  }

  async function handleOffer(from, offer) {
    // if (!mediaStore.localStream) {
    //   await mediaStore.setMicDevice(mediaStore.micSettings.deviceId)
    //   if (mediaStore.camSettings.enabled) {
    //     await mediaStore.setCamDevice(mediaStore.camSettings.deviceId)
    //   }
    // }
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
    pendingCandidates.value = {}
    mediaStore.cleanupAllRemoteStreams()
  }

  return {
    peerConnections,
    startCall,
    handleOffer,
    handleAnswer,
    handleIceCandidate,
    endCallWith,
    endAllCalls
  }
})