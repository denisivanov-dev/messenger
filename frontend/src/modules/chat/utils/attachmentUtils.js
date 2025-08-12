import { getImagefromKey } from '../api/chatApi'

export async function loadAttachmentUrls(messages, imageUrlCache, attachmentUrlsRef) {
  const now = Date.now()
  
  if (!messages || messages.length === 0) return
  if (!imageUrlCache) return
  if (!attachmentUrlsRef?.value || typeof attachmentUrlsRef.value !== 'object') {
    attachmentUrlsRef.value = {}
  }

  for (const msg of messages) {
    if (!msg?.attachments?.length) continue

    for (const att of msg.attachments) {
      if (!att?.key || att.type !== 'image') continue

      const cached = imageUrlCache[att.key]
      const isValid = cached && cached.expiresAt > now

      if (isValid) {
        attachmentUrlsRef.value[att.key] = cached.url
        continue
      }

      try {
        const { url } = await getImagefromKey(att.key)
        if (url) {
          const expiresAt = now + 59 * 60 * 1000
          imageUrlCache[att.key] = { url, expiresAt }
          attachmentUrlsRef.value[att.key] = url

          setTimeout(() => {
            if (imageUrlCache[att.key]?.expiresAt <= Date.now()) {
              delete imageUrlCache[att.key]
            }
          }, expiresAt - now + 1000)
        }
      } catch (err) {
        console.error("Ошибка загрузки:", att.key, err)
      }
    }
  }
}