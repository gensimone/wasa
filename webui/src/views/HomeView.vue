<script>
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationsList from "@/components/Conversations/ConversationsList.vue"
import { clearUserState } from "@/state/user"

export default {
    components: { Topbar, Bottombar, ConversationsList },

    methods: {
        select(conversation) {
            this.$router.push({
                name: "conversation",
                params: { id: conversation.id },
                query: { direct: !conversation.isGroup }
            })
        },

        logout() {
            clearUserState()
            this.$router.push("/")
        }
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/settings.svg', onClick: () => $router.push('/user/settings') },
            { icon: '/icons/logout.svg', onClick: () => logout() }
        ]" />
        <div class="content">
            <ConversationsList @select="select" />
        </div>
        <Bottombar />
    </div>
</template>
