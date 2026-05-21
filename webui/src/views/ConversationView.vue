<script>
import { user } from "@/state/user.js"
import { conversation } from "@/state/conversation.js"
import ConversationService from "@/services/conversationService"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationBox from "@/components/Conversation/ConversationBox.vue"

export default {
    components: {
        Topbar,
        Bottombar,
        ConversationBox
    },

    data() {
        return {
            service: null,                       // Service for polling/send messages.
            userId: user.userId,                 // The user ID of the logged user.
            messages: [],                        // Messages

            conversationName: conversation.name, // Group or user name.
            avatarUrl: conversation.photoUrl,    // Group or user photo.

            text: null,                          // Text for the InputBox.
            attachment: null,                    // Attachment for the InputBox.
            attachmentUrl: null,                 // Attachment URL for the InputBox.
            sending: false,                      // Enable/Disable the InputBox.
            sendError: null,                     // Error message in SendError.
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
        }
    },

    async mounted() {
        this.service = new ConversationService()

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
            <ConversationBox ref="conversationBox" :conversationName="conversationName" :avatarUrl="avatarUrl"
                :messages="messages" :userId="userId" :text="text" :attachment="attachment"
                :attachmentUrl="attachmentUrl" @update:text="text = $event" @send="sendMessage"
                @addAttachment="addAttachment" @removeAttachment="removeAttachment" :sending="sending"
                :sendError="sendError" />
        </div>
        <Bottombar />
    </div>
</template>
