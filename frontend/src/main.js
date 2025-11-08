import { createApp } from 'vue'
import App from './app.vue'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import { createPinia } from 'pinia'

import LoginPage from './modules/auth/pages/loginPage.vue'
import ForgotPasswordPage from './modules/auth/pages/forgotPasswordPage.vue' 
import registerPage from './modules/auth/pages/registerPage.vue' 
import confirmRegistrationPage from './modules/auth/pages/confirmRegistrationPage.vue' 
import ChatPage from './modules/chat/features/chat/pages/chatPage.vue'
import PrivateChatPage from './modules/chat/features/chat/pages/privateChatPage.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginPage },
  { path: '/forgot-password', component: ForgotPasswordPage },
  { path: '/register', component: registerPage},
  { path: '/confirm-registration', component: confirmRegistrationPage},
  { path: '/global-chat', component: ChatPage},
  { path: '/private-chat', component: PrivateChatPage }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)

app.use(createPinia()) 
app.use(router)

app.mount('#app')