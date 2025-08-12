import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sendMessage, startPrivateChat } from '../api/chatApi'
import { isNearBottom } from '../utils/chatWindowUtils'

export const useMessagesStore = defineStore('messages', () => {
  const messages = ref([])
  const shouldScroll = ref(false)

  function clear() {
    messages.value = []
  }

  function pushFromWs(msg) {
    try {
      if (
        typeof msg !== 'object' ||
        !msg.message_id ||
        (!msg.text && msg.attachments?.length === 0)
      ) return
      messages.value.push(msg)
    } catch (err) {
      console.error('PUSH CRASH:', err, msg)
    }
    shouldScroll.value = isNearBottom()
  }

  function handleDeleted(messageId) {
    messages.value = messages.value.filter(m => m.message_id !== messageId)
  }

  function handleEdited(messageId, newText, editedAt) {
    const edited = messages.value.find(m => m.message_id === messageId)
    if (edited) {
      edited.text = newText
      edited.timestamp = editedAt
      edited.edited_at = editedAt
    }
  }

  function handlePinned(messageId, action) {
    const pinned = messages.value.find(m => m.message_id === messageId)
    if (pinned) pinned.pinned = action === 'pin'
  }

  function sendMessageData({ text = '', attachments = [], replyToMessage = null }, chatType, receiverID) {
    if (!text.trim() && attachments.length === 0) return

    const safeAttachments = attachments.map(att => ({
      key: att.key,
      type: att.type,
      size: att.size,
      original_name: att.original_name,
    }))

    const msg = {
      type: 'send_message',
      text: text.trim(),
      attachments: safeAttachments,
      timestamp: Date.now(),
      chat_type: chatType,
      receiver_id: receiverID,
    }

    if (replyToMessage) {
      msg.reply_to = replyToMessage.message_id
      msg.reply_to_text = replyToMessage.text
      msg.reply_to_user = replyToMessage.username
    }
    const json = JSON.stringify(msg)
    console.log('📦 WS payload size:', json.length)

    sendMessage(msg)
    shouldScroll.value = true
  }

  function deleteMessage(message, chatType, receiverID) {
    sendMessage({
      type: 'delete_message',
      message_id: message.message_id,
      receiver_id: receiverID,
      chat_type: chatType
    })
  }

  function editMessage(message, newText, chatType, receiverID) {
    sendMessage({
      type: 'edit_message',
      message_id: message.message_id,
      new_text: newText,
      receiver_id: receiverID,
      chat_type: chatType
    })
  }

  function pinMessage(message, shouldPin, chatType, receiverID) {
    sendMessage({
      type: 'pin_message',
      message_id: message.message_id,
      chat_id: message.chat_id,
      action: shouldPin ? 'pin' : 'unpin',
      receiver_id: receiverID,
      chat_type: chatType
    })
  }

  async function openOrCreatePrivateChat(targetId, receiverIDRef) {
    try {
      messages.value = []
      const response = await startPrivateChat(targetId)
      receiverIDRef.value = targetId
      localStorage.setItem('receiverId', targetId)
      console.info(response)
      console.info(receiverIDRef.value)
      messages.value = response.messages
      shouldScroll.value = true
      return response.success
    } catch (err) {
      console.error('openOrCreatePrivateChat error:', err)
      throw err
    }
  }

  return {
    messages,
    shouldScroll,

    clear,
    pushFromWs,
    handleDeleted,
    handleEdited,
    handlePinned,

    sendMessageData,
    deleteMessage,
    editMessage,
    pinMessage,
    openOrCreatePrivateChat,
  }
})