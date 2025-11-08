import { nextTick } from 'vue'

function getChatContainer () {
  return document.querySelector('.chat-window')
}

/**
 * Проверяет, находится ли пользователь рядом с нижним краем чата.
 * @param {number} threshold — запас в пикселях
 */
export function isNearBottom (threshold = 120) {
  const el = getChatContainer()
  if (!el) return false
  return el.scrollTop + el.clientHeight >= el.scrollHeight - threshold
}

export function waitForImagesAndThen(callback, containerEl) {
  if (!containerEl) return callback()

  const images = containerEl.querySelectorAll('img')
  if (images.length === 0) {
    callback()
    return
  }

  let loaded = 0

  images.forEach((img) => {
    const isLoaded = img.complete && img.naturalHeight !== 0

    if (isLoaded) {
      loaded++
    } else {
      img.onload = img.onerror = () => {
        loaded++
        if (loaded === images.length) callback()
      }
    }
  })

  if (loaded === images.length) {
    callback()
  }
}