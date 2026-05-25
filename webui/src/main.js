import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router'
import notifier from "@/notifier"

import { user, hydrateUserState, startPollingUser } from '@/state/user.js'
import { startPollingConversations } from '@/state/conversations.js'

import './styles/index.css'

hydrateUserState()
if (user.userId) {
    startPollingUser()
    startPollingConversations()
}

const app = createApp(App)
app.use(router)
app.use(notifier)
app.mount('#app')
