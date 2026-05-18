import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { loadUserFromStorage } from './state/user.js'
import { loadChatMetadataFromStorage } from './state/chat.js'

import './main.css'

loadUserFromStorage()
loadChatMetadataFromStorage()

const app = createApp(App)
app.use(router)
app.mount('#app')
