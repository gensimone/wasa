<script>
import { conversation } from "@/state/conversation"
import { expandUrl } from "@/utils/media"
import ImageModal from "@/components/Shared/ImageModal.vue"
import MessageList from "@/components/Conversation/MessageList.vue"
import ConversationInput from "@/components/Conversation/ConversationInput.vue"

export default {
    components: {
        MessageList,
        ConversationInput,
        ImageModal
    },

    props: {
        messages: { type: Array, required: true },
        text: { type: String, default: "" },
        attachment: { type: File, default: null },
        attachmentUrl: { type: String, default: null },
        sending: { type: Boolean, default: false }
    },

    emits: [
        "send",
        "addAttachment",
        "removeAttachment",
        "update:text",
        "infoGroup"
    ],

    data() {
        return {
            conversation,

            zoomedImage: null,
            showImageModal: false,
            conversationName: conversation.name,
            photoUrl: expandUrl(conversation.photoUrl),
            isGroup: conversation.isGroup
        }
    },

    watch: {
        "conversation.photoUrl"(newPhotoUrl, _) {
            this.photoUrl = expandUrl(newPhotoUrl)
        },

        "conversation.name"(newName, _) {
            this.conversationName = newName
        }
    },

    methods: {
        expandUrl,

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
                <img class="photo" :src="photoUrl" @click="openImage(photoUrl)" />
                <div class="conversation-name">
                    {{ conversationName }}
                </div>
                <div class="info-button">
                    <button v-if="isGroup" class="info-btn" @click="$emit('infoGroup')">
                        <img src="/icons/info.svg" class="icon-img">
                    </button>
                </div>
            </div>

            <MessageList ref="messageList" :messages="messages" @openImage="openImage" />

            <ConversationInput :text="text" :sending="sending" :attachment="attachment" :attachmentUrl="attachmentUrl"
                @update:text="$emit('update:text', $event)" @send="$emit('send')"
                @addAttachment="$emit('addAttachment', $event)" @removeAttachment="$emit('removeAttachment')" />
        </div>
        <ImageModal :visible="showImageModal" :imageUrl="zoomedImage" @close="closeImage" />
    </div>
</template>

<style scoped>
.info-button {
    position: absolute;
    right: 16px;
}

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

.conversation-header .photo {
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

.info-btn {
    width: 46px;
    height: 46px;
    border-radius: 100px;
    border: 0px;
    background: rgba(0, 0, 0, 0.0);

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;
    position: relative;
    overflow: hidden;

    transition: transform 0.25s ease, border 0.25s ease, box-shadow 0.25s ease;
}

.info-btn:hover::before {
    transform: translateX(140%);
}

.info-btn:hover {
    transform: translateY(-4px) scale(1.05);
}

.info-btn:active {
    transform: scale(0.95);
}
</style>
