import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router'
import { loadUserFromStorage } from '@/state/user.js'
import { loadConversationMetadataFromStorage } from '@/state/conversation.js'

import './main.css'

loadUserFromStorage()
loadConversationMetadataFromStorage()

const app = createApp(App)
app.use(router)
app.mount('#app')
