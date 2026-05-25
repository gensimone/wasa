<script>
import { sendMessage } from "@/services/users";
import { sendMessageToConversation } from "@/services/conversations";
import { handleError } from "@/utils/errors"

export default {
    props: {
        id: { type: Number, required: true },
        routeType: {
            type: String,
            required: true,
            validator: (v) =>
                ["user", "conversation"].includes(v)
        }
    },

    data() {
        return {
            text: "",
            attachment: null,
            attachmentUrl: null,
            sending: false
        }
    },

    emits: ["reportConversationId", "triggerScrolldown", "pushMessage"],

    methods: {
        async send() {
            const cleanText = this.text?.trim()
            if (!cleanText && !this.attachment) return

            this.sending = true

            try {
                let message
                if (this.routeType == "user") {
                    message = await sendMessage(
                        this.id,
                        this.text,
                        this.attachment
                    )

                    this.$emit("reportConversationId", message.conversationId)

                } else {
                    message = await sendMessageToConversation(
                        this.id,
                        this.text,
                        this.attachment
                    )
                }

                this.$emit("pushMessage", message)
                this.$emit("triggerScrolldown")

                this.text = ""
                if (this.attachment) {
                    this.removeAttachment()
                }

            } catch (e) {
                handleError(e)

            } finally {
                this.sending = false
            }
        },

        revokeAttachment() {
            if (this.attachmentUrl) {
                URL.revokeObjectURL(this.attachmentUrl)
            }
        },

        addAttachment(event) {
            const file = event.target.files[0]
            if (!file) return

            this.revokeAttachment()

            this.attachment = file
            this.attachmentUrl = URL.createObjectURL(file)

            event.target.value = ""
        },

        removeAttachment() {
            this.revokeAttachment()

            this.attachmentUrl = null
            this.attachment = null
        }
    },

    beforeUnmount() {
        this.revokeAttachment()
    }
}
</script>

<template>
    <div>
        <div v-if="attachment" class="conversation-input-attachment-preview">
            <img class="conversation-input-preview-img" :src="attachmentUrl" />
            <button class="icon-btn" @click="removeAttachment">
                <img src="/icons/remove.svg" class="icon-img">
            </button>
        </div>

        <div class="conversation-input">
            <button class="icon-btn" @click="$refs.fileInput.click()">
                <img src="/icons/upload.svg" class="icon-img">
            </button>

            <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="addAttachment($event)" />

            <input class="input-bar" :value="text" @input="text = $event.target.value" placeholder="Type a message..."
                @keydown.enter.prevent="send" />

            <button class="icon-btn" :disabled="sending" @click="send">
                <img src="/icons/send.svg" class="icon-img">
            </button>
        </div>
    </div>
</template>

<style scoped>
.conversation-input {
    display: flex;
    gap: 10px;
}

.conversation-input-attachment-preview {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 0px auto 0;
    padding: 8px 10px;
    width: min(720px, 100%);
    border-radius: 0px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(10px);
}

.conversation-input-preview-img {
    width: 50px;
    height: 50px;
    border-radius: 12px;
    object-fit: cover;
    border: 1px solid rgba(255, 255, 255, 0.12);
}
</style>
