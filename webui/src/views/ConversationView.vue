<script>
import { conversation, startPollingConversation } from "@/state/conversation"
import { getConversationMessages, sendMessage } from "@/utils/conversations"
import { setMessageStatusAsRead } from "@/services/messages"
import { user } from "@/state/user"
import { Poller } from "@/services/poller"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationBox from "@/components/Conversation/ConversationBox.vue"
import { handleError } from "@/utils/errors"

export default {
    components: {
        Topbar,
        Bottombar,
        ConversationBox
    },

    data() {
        return {
            poller: null,
            messages: [],
            scrollTick: 0, // Used to trigger the scroll down.

            text: null,
            attachment: null,
            attachmentUrl: null,
            sending: false
        }
    },

    watch: {
        // Temporary solution for message receipts (read).
        async messages(newValue, oldValue) {
            const equal = JSON.stringify(newValue) === JSON.stringify(oldValue)

            if (equal) return

            try {
                const requests = this.messages
                    .filter(m => m.senderId !== user.userId)
                    .map(m => setMessageStatusAsRead(m.messageId))

                await Promise.all(requests)
            } catch (e) {
                handleError(e)
            }
        }
    },

    methods: {
        async sendMessage() {
            const cleanText = this.text?.trim()
            if (!cleanText && !this.attachment) return

            this.sending = true

            try {
                const message = await sendMessage(
                    this.text,
                    this.attachment
                )

                this.messages.push(message)

                this.scrollTick++

                this.text = null
                if (this.attachment)
                    this.removeAttachment()

            } catch (e) {
                handleError(e)

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

        infoGroup() {
            this.$router.push('/group/info')
        }
    },

    async mounted() {
        startPollingConversation()
        this.poller = new Poller(async () => {
            this.messages = await getConversationMessages()
        })

        this.poller.startPolling()
    },

    beforeUnmount() {
        conversation.poller?.stopPolling()
        this.poller?.stopPolling()
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/home.svg', onClick: () => $router.push('/home') }
        ]" />
        <div class="content">
            <ConversationBox ref="conversationBox" :messages="messages" :scrollTick="scrollTick" :text="text"
                :attachment="attachment" :attachmentUrl="attachmentUrl" @infoGroup="infoGroup"
                @update:text="text = $event" @send="sendMessage" @addAttachment="addAttachment"
                @removeAttachment="removeAttachment" :sending="sending" />
        </div>
        <Bottombar />
    </div>
</template>
