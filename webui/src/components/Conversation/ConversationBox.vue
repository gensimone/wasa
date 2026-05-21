<script>
import { expandURL } from "@/utils/media"
import ImageModal from "@/components/Shared/ImageModal.vue"
import MessageList from "@/components/Conversation/MessageList.vue"
import ConversationInput from "@/components/Conversation/ConversationInput.vue"
import SendError from "@/components/Conversation/SendError.vue"

export default {
    components: {
        MessageList,
        ConversationInput,
        ImageModal,
        SendError
    },

    props: {
        conversationName: { type: String, required: true },
        avatarUrl: { type: String, required: true },

        messages: { type: Array, required: true },
        userId: { type: [String, Number], required: true },

        text: { type: String, default: "" },
        attachment: { type: File, default: null },
        attachmentUrl: { type: String, default: null },

        sending: { type: Boolean, default: false },
        sendError: { type: String, default: null }
    },

    emits: [
        "send",
        "addAttachment",
        "removeAttachment",
        "update:text",
    ],

    data() {
        return {
            zoomedImage: null,
            showImageModal: false,
        }
    },

    methods: {
        expandURL,

        openImage(url) {
            this.zoomedImage = url
            this.showImageModal = true
        },

        closeImage() {
            this.showImageModal = false
            this.zoomedImage = null
        }
    }
}
</script>

<template>
    <div class="conversation-container">
        <div class="conversation-box">

            <div class="conversation-header">
                <img class="avatar" :src="expandURL(avatarUrl)" @click="openImage(avatarUrl)" />
                <div class="conversation-name">
                    {{ conversationName }}
                </div>
            </div>

            <MessageList ref="messageList" :messages="messages" :userId="userId" @openImage="openImage" />

            <ConversationInput :text="text" :sending="sending" :attachment="attachment" :attachmentUrl="attachmentUrl"
                @update:text="$emit('update:text', $event)" @send="$emit('send')"
                @addAttachment="$emit('addAttachment', $event)" @removeAttachment="$emit('removeAttachment')" />

            <SendError :error="sendError" />

        </div>
        <ImageModal :visible="showImageModal" :imageUrl="zoomedImage" @close="closeImage" />
    </div>
</template>

<style scoped>
.conversation-container {
    display: flex;
    width: 100%;
    justify-content: center;
    height: calc(100vh - 50px - 60px);
    box-sizing: border-box;
}

.conversation-box {
    width: min(720px, 100%);
    height: 95%;
    display: flex;
    flex-direction: column;
    border-radius: 22px;
    position: relative;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);
}

.conversation-header {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;

    padding: 12px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(255, 255, 255, 0.02);
    backdrop-filter: blur(10px);
}

.conversation-header .avatar {
    position: absolute;
    left: 16px;

    width: 38px;
    height: 38px;
    border-radius: 50%;
    object-fit: cover;

    border: 1px solid rgba(255, 255, 255, 0.2);
}

.conversation-header .conversation-name {
    font-size: 1rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
}
</style>
