import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router'
import notifier from "@/notifier"

import { loadUserState } from '@/state/user.js'
import { loadConversationState } from '@/state/conversation.js'
import { loadGroupState } from '@/state/group'

import './main.css'

// Load local storage data.
loadUserState()
loadConversationState()
loadGroupState()

const app = createApp(App)
app.use(router)
app.use(notifier)
app.mount('#app')
