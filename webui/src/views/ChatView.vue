<script>
import { user } from "@/state/user.js"
import { chat } from "@/state/chat.js"
import ChatService from "@/services/chatService"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ChatBox from "@/components/Chat/ChatBox.vue"

export default {
    components: {
        Topbar,
        Bottombar,
        ChatBox
    },

    data() {
        return {
            service: null,
            userId: user.userId,
            messages: [],
            sending: false,
            sendError: null,
            chatName: chat.name,
            avatarUrl: chat.photoUrl,
            text: null,
            attachment: null,
            attachmentUrl: null
        }
    },

    methods: {
        async sendMessage() {
            this.sending = true
            this.sendError = null

            try {
                const response = await this.service.sendMessage(
                    this.text,
                    this.attachment
                )

                this.messages.push(response.data)

                this.text = null
                if (this.attachment)
                    this.removeAttachment()

            } catch (e) {
                this.sendError = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.sending = false
            }
        },

        addAttachment(event) {
            const file = event.target.files[0]
            if (!file) return

            if (this.attachmentUrl) {
                URL.revokeObjectURL(this.attachmentUrl)
            }

            this.attachment = file
            this.attachmentUrl = URL.createObjectURL(file)

            event.target.value = ""
        },

        removeAttachment() {
            if (this.attachmentUrl) {
                URL.revokeObjectURL(this.attachmentUrl)
            }

            this.attachmentUrl = null
            this.attachment = null
        },

        expandURL(url) {
            return `${__API_URL__}${url}`
        }
    },

    async mounted() {
        this.service = new ChatService()

        this.messages = await this.service.fetchMessages()

        this.service.startPolling(({ messages }) => {
            this.messages = messages
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
            { icon: '/icons/home.svg', onClick: () => $router.push('/home') },
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <div class="content">
            <ChatBox ref="chatBox" :chatName="chatName" :avatarUrl="expandURL(avatarUrl)" :messages="messages" :userId="userId"
                :text="text" :attachment="attachment" :attachmentUrl="attachmentUrl" @update:text="text = $event"
                @send="sendMessage" @addAttachment="addAttachment" @removeAttachment="removeAttachment"
                :sending="sending" :sendError="sendError" />
        </div>
        <Bottombar />
    </div>
</template>
