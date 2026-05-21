<script>
import { clearUserState } from "@/state/user"
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup,
    clearConversationState
} from "@/state/conversation"
import ConversationsService from "@/services/conversationsService"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationsList from "@/components/Home/ConversationsList.vue"

export default {
    components: {
        Topbar,
        Bottombar,
        ConversationsList
    },

    data() {
        return {
            conversations: [],
            error: null,
            loading: false
        }
    },

    methods: {
        startConversation(conversation) {
            setConversationName(conversation.name)
            setConversationPhotoUrl(conversation.photoUrl)
            setConversationId(conversation.id)
            setConversationIsGroup(conversation.isGroup)
            this.$router.push("/conversation")
        },

        logout() {
            clearUserState()
            clearConversationState()
            localStorage.clear()
            this.$router.replace("/")
        }
    },

    async mounted() {
        this.service = new ConversationsService()

        this.conversations = await this.service.fetchConversations()

        this.service.startPolling(({ conversations }) => {
            this.conversations = conversations
        })
    },

    beforeUnmount() {
        this.service?.stopPolling()
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/settings.svg', onClick: () => $router.push('/settings') },
            { icon: '/icons/logout.svg', onClick: () => logout() }
        ]" />
        <div class="content">
            <ConversationsList :conversations="conversations" :loading="loading" @select="startConversation" />
        </div>
        <Bottombar />
    </div>
</template>
