<script>
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup,
    clearConversationState
} from "@/state/conversation"

import {
    setGroupFounderId,
    setGroupCreatedAt,
    clearGroupState
} from "@/state/group"

import { clearUserState } from "@/state/user"
import { Poller } from "@/services/poller"
import { getConversations } from "@/utils/conversations"
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
            poller: null,
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
            setGroupFounderId(conversation.founderId)
            setGroupCreatedAt(conversation.createdAt)

            this.$router.push("/conversation")
        },

        logout() {
            clearUserState()
            clearConversationState()
            clearGroupState()
            localStorage.clear()

            this.$router.replace("/")
        }
    },

    async mounted() {
        this.conversations = await getConversations()
        this.poller = new Poller(async () => {
            this.conversations = await getConversations()
        })

        this.poller.startPolling()
    },

    beforeUnmount() {
        this.poller?.stopPolling()
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
