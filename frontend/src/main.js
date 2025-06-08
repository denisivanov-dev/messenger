import { createApp } from 'vue'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router'
import './style.css'
import { createPinia } from 'pinia'

import LoginPage from './modules/auth/pages/loginPage.vue'
import ForgotPasswordPage from './modules/auth/pages/forgotPasswordPage.vue' 
import registerPage from './modules/auth/pages/registerPage.vue' 
import confirmRegistration from './modules/auth/pages/confirmRegistration.vue' 

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: LoginPage },
  { path: '/forgot-password', component: ForgotPasswordPage },
  { path: '/register', component: registerPage},
  { path: '/confirm-registration', component: confirmRegistration}
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const app = createApp(App)

app.use(createPinia()) 
app.use(router)

app.mount('#app')