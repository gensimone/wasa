<script>
import MessageList from "@/components/Chat/MessageList.vue"
import ChatInput from "@/components/Chat/ChatInput.vue"
import ImageModal from "@/components/Shared/ImageModal.vue"
import SendError from "@/components/Chat/SendError.vue"

export default {
    components: {
        MessageList,
        ChatInput,
        ImageModal,
        SendError
    },

    props: {
        chatName: { type: String, required: true },
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
        "react"
    ],

    data() {
        return {
            zoomedImage: null,
            showImageModal: false,
        }
    },

    methods: {
        scrollToBottom() {
            this.$refs.messageList?.scrollToBottom()
        },

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
    <div class="chat-container">
        <div class="chat-box">

            <div class="chat-header">
                <img class="avatar" :src="avatarUrl" @click="openImage(avatarUrl)" />
                <div class="chat-name">
                    {{ chatName }}
                </div>
            </div>

            <MessageList ref="messageList" :messages="messages" :userId="userId" @openImage="openImage" />

            <ChatInput :text="text" :sending="sending" :attachment="attachment" :attachmentUrl="attachmentUrl"
                @update:text="$emit('update:text', $event)" @send="$emit('send')"
                @addAttachment="$emit('addAttachment', $event)" @removeAttachment="$emit('removeAttachment')" />

            <SendError :error="sendError" />

        </div>
        <ImageModal :visible="showImageModal" :imageUrl="zoomedImage" @close="closeImage" />
    </div>
</template>

<style scoped>
.chat-container {
    display: flex;
    width: 100%;
    justify-content: center;
    height: calc(100vh - 50px - 60px);
    box-sizing: border-box;
}

.chat-box {
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

.chat-header {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;

    padding: 12px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(255, 255, 255, 0.02);
    backdrop-filter: blur(10px);
}

.chat-header .avatar {
    position: absolute;
    left: 16px;

    width: 38px;
    height: 38px;
    border-radius: 50%;
    object-fit: cover;

    border: 1px solid rgba(255, 255, 255, 0.2);
}

.chat-header .chat-name {
    font-size: 1rem;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
}
</style>
