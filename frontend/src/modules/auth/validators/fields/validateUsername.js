export const validateUsername = (username) => {
    username = username.toLowerCase()

    if (username.length < 3 || username.length > 20) {
        return { username: '- Имя должно быть от 3 до 20 символов' }
    }

    if (!/^[a-z][a-z0-9_]*$/.test(username)) {
        return { username: '- Имя должно содержать только буквы и цифры' }
    }

    if (username === username[0].repeat(username.length)) {
        return { username: '- Имя не может состоять из одного символа' }
    }

    if (username.includes('__')) {
        return { username: '- Имя не должно содержать двойные подчёркивания' }
    }

    if (username.endsWith('_')) {
        return { username: '- Имя не должно заканчиваться подчёркиванием' }
    }

    return null
}