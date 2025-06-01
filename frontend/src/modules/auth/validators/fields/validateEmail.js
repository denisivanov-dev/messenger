export const validateEmail = (email) => {
   email = email.toLowerCase()

   const emailRegex = /^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/
   if (!emailRegex.test(email)) {
      return { email: '- Неверный формат email' }
   }

   const parts = email.split('@')
   if (parts.length !== 2) {
      return { email: '- Некорректный email' }
   }

   return null
}