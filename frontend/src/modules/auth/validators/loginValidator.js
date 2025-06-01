export const validateLoginForm = ({ email, password }) => {
   email = email.trim()
   password = password.trim()

   if (!email) {
      return { email: '- Введите email' }
   }

   const emailRegex = /^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/
   if (!emailRegex.test(email)) {
      return { email: '- Некорректный email' }
   }

   if (!password) {
      return { password: '- Введите пароль' }
   }

   return null
}