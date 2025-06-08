import { validateEmail } from "./fields/validateEmail"
import { validatePassword } from "./fields/validatePassword"

export const validateLoginForm = ({ email, password }) => {
   email = email.trim()
   password = password.trim()

   if (!email) {
      return { email: '- Введите email' }
   }

   if (!password) {
      return { password: '- Введите пароль' }
   }

   let error = validateEmail(email)
   if (error) return error
   
   error = validatePassword(password, 'password')
   if (error) return error

   return null
}