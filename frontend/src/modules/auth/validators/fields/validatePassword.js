export const validatePassword = (password, field) => {
    if (password.length < 8 || password.length > 128) {
        return { [field]: '- Пароль должен быть от 8 до 128 символов' }
    }

    if (password === password[0].repeat(password.length)) {
        return { [field]: '- Пароль не может состоять из одного символа' }
    }

    if (!/[A-Z]/.test(password)) {
        return { [field]: '- Пароль должен содержать хотя бы одну заглавную букву' }
    }

    if (!/[a-z]/.test(password)) {
        return { [field]: '- Пароль должен содержать хотя бы одну строчную букву' }
    }

    if (!/[0-9]/.test(password)) {
        return { [field]: '- Пароль должен содержать хотя бы одну цифру' }
    }

    if (!/[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)) {
        return { [field]: '- Пароль должен содержать хотя бы один специальный символ' }
    }

    return null
}