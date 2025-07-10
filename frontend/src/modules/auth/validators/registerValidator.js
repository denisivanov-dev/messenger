import { validateEmail } from "./fields/validateEmail"
import { validateUsername } from "./fields/validateUsername"
import { validatePassword } from "./fields/validatePassword"

export const validateRegisterForm = ({email, username, password, confirmPassword}) => {
   email = email.trim()
   username = username.trim()
   password = password.trim()
   confirmPassword = confirmPassword.trim()

   if (!email) { 
      return { email: '- Введите email' }
   }

   if (!username) {
      return { username: '- Введите имя' }
   }

   if (!password) {
      return { password: '- Введите пароль' }
   }

   if (!confirmPassword) {
      return { confirmPassword: '- Подтвердите пароль' }
   }

   let error = validateEmail(email)
   if (error) return error

   error = validateUsername(username)
   if (error) return error

   error = validatePassword(password, 'password')
   if (error) return error

   error = validatePassword(confirmPassword, 'confirmPassword')
   if (error) return error

   if (password !== confirmPassword) {
      return { password: '- Пароли не совпадают', confirmPassword: '- Пароли не совпадают' }
   }

   return null
}