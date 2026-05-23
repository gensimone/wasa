<script>
import { getConversationMessages, sendMessage } from "@/utils/conversations"
import { conversation } from "@/state/conversation"
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

            text: null,
            attachment: null,
            attachmentUrl: null,
            sending: false,
            sendError: null,
        }
    },

    methods: {
        async sendMessage() {
            this.sending = true
            this.sendError = null

            try {
                const message = await sendMessage(
                    this.text,
                    this.attachment
                )

                this.messages.push(message)

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
        this.poller = new Poller(async () => {
            this.messages = await getConversationMessages()
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
            { icon: '/icons/home.svg', onClick: () => $router.push('/home') },
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <div class="content">
            <ConversationBox ref="conversationBox" :messages="messages" :text="text" :attachment="attachment"
                :attachmentUrl="attachmentUrl" @infoGroup="infoGroup" @update:text="text = $event" @send="sendMessage"
                @addAttachment="addAttachment" @removeAttachment="removeAttachment" :sending="sending"
                :sendError="sendError" />
        </div>
        <Bottombar />
    </div>
</template>
