<script>
import { user } from "@/state/user.js"
import { chat } from "@/state/chat.js"
import backIcon from "@/assets/icons/back.svg"
import homeIcon from "@/assets/icons/home.svg"
import likeIcon from "@/assets/icons/like.svg"
import loveIcon from "@/assets/icons/love.svg"
import laughIcon from "@/assets/icons/laugh.svg"
import sadIcon from "@/assets/icons/sad.svg"
import angryIcon from "@/assets/icons/angry.svg"
import uploadIcon from "@/assets/icons/upload.svg"
import sendIcon from "@/assets/icons/send.svg"
import deleteIcon from "@/assets/icons/delete.svg"
import LoggedAs from "@/components/LoggedAs.vue"
export default {
    components: {
        LoggedAs
    },

    data() {
        return {
            userId: user.userId, // Logged user userId (used to determine css in messages)
            messages: [],        // Conversation messages.
            lock: false,         // Simple look used to synchronize access to this.messages
            loading: false,
            error: null,
            messageIds: [],
            poller: null,

            chatName: chat.name,      // Chat name (group or user name)
            avatarUrl: chat.photoUrl, // Chat photo (group or user photo)

            zoomedImage: null,     // Image showed when the user click on an attachment.
            showImageModal: false,

            reactionPopup: {
                visible: false,
                messageId: null,
                x: 0,
                y: 0
            },

            text: null,          // Message text to send.
            attachment: null,    // Attachment image
            attachmentUrl: null, // Attachment preview

            // Icons
            backIcon,
            deleteIcon,
            homeIcon,
            sendIcon,
            uploadIcon,
            likeIcon,
            loveIcon,
            laughIcon,
            sadIcon,
            angryIcon,
        }
    },
    methods: {
        async sendMessage() {
            this.loading = true
            try {
                const formData = new FormData()
                formData.append("text", this.text)

                if (this.attachment) {
                    formData.append("file", this.attachment)
                    formData.append("mediaType", "image")
                    this.attachment = null
                }

                const response = await this.$axios.post(`/users/${chat.userId}/message`,
                    formData,
                    {
                        headers: {
                            Authorization: user.userId,
                            "Content-Type": "multipart/form-data"
                        }
                    }
                )

                this.text = null
                this.scrollToBottom()
            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
            }
            this.loading = false
        },

        async fetchMessage(messageId) {
            const response = await this.$axios.get(`/messages/${messageId}`, {
                headers: { Authorization: user.userId }
            })
            return response.data
        },

        async fetchMessages() {
            if (this.lock) return
            this.lock = true
            try {
                const response = await this.$axios.get(`/users/${chat.userId}/messages`, {
                    headers: { Authorization: user.userId }
                })

                const newMessageIds = response.data.messageIds
                const oldIdsSet = new Set(this.messageIds)
                const newIdsSet = new Set(newMessageIds)

                const changed =
                    newMessageIds.some(id => !oldIdsSet.has(id)) ||  // Is there a new message? Or
                    this.messages.some(id => !newIdsSet.has(id))     // Or has an old message been deleted?

                if (!changed) {
                    this.lock = false
                    return
                }

                const diffMessageIds = [...newIdsSet].filter(x => !oldIdsSet.has(x))

                const results = await Promise.all(
                    diffMessageIds.map(messageId => this.fetchMessage(messageId))
                )

                const newMessages = results
                    .map(data => ({
                        messageId: Number(data.messageId),
                        text: data.text,
                        senderId: Number(data.senderId),
                        conversationid: Number(data.conversationId),
                        createdAt: data.createdAt,
                        isForwarded: Boolean(data.isForwarded),
                        commentTo: Number(data.commentTo),
                        attachmentUrl: data.attachmentUrl == "" ? "" : `${__API_URL__}${data.attachmentUrl}`,
                        mediaType: data.mediaType
                    }))

                this.messageIds = newMessageIds
                const merged = [...this.messages, ...newMessages]
                merged.sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt))
                this.messages = merged

                const el = this.$el.querySelector(".messages")
                const nearBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 100
                if (nearBottom) {
                    el.scrollTop = el.scrollHeight
                }
            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.lock = false
            }
        },

        scrollToBottom() {
            this.$nextTick(() => {
                const el = this.$el.querySelector(".messages")
                if (el) el.scrollTop = el.scrollHeight
            })
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

        getMessageSentTime(sentAt) {
            const date = new Date(sentAt)
            const hh = String(date.getHours()).padStart(2, "0")
            const mm = String(date.getMinutes()).padStart(2, "0")
            return `${hh}:${mm}`
        },

        openImage(url) {
            this.zoomedImage = url
            this.showImageModal = true
        },

        closeImage() {
            this.showImageModal = false
            this.zoomedImage = null
        },

        openReactionPopup(event, msg) {
            this.reactionPopup = {
                visible: true,
                messageId: msg.messageId,
                x: event.clientX,
                y: event.clientY
            }
        },

        closeReactionPopup() {
            this.reactionPopup.visible = false
            this.reactionPopup.messageId = null
        },

        react(emoji) {
            if (!this.reactionPopup.messageId) return

            this.addReaction(this.reactionPopup.messageId, emoji)
            this.closeReactionPopup()
        },

        addReaction(messageId, emoji) {
            console.log("reaction:", messageId, emoji)

            // TODO: API call
            // return this.$axios.post(...)
        }
    },

    async mounted() {
        await this.fetchMessages()
        this.scrollToBottom()
        this.poller = setInterval(() => {
            this.fetchMessages()
        }, 2000)
    },

    beforeUnmount() {
        clearInterval(this.poller)
        this.poller = null
    }
}
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <button class="icon-btn" @click="$router.push('/home')">
                    <img :src="homeIcon" class="icon-img">
                </button>
                <button class="icon-btn" @click="$router.back()">
                    <img :src="backIcon" class="icon-img">
                </button>
            </div>
        </header>
        <!-- CHAT AREA -->
        <div class="chat-container">
            <div class="chat-box">
                <!-- CHAT HEADER -->
                <div class="chat-header">
                    <img class="avatar" :src="avatarUrl" />
                    <div class="chat-name">
                        {{ chatName }}
                    </div>
                </div>
                <!-- MESSAGES -->
                <div class="messages">
                    <div v-for="msg in messages" :key="msg.messageId" class="message"
                        :class="{ mine: Number(msg.senderId) === Number(userId) }">
                        <div class="bubble" @contextmenu.prevent="openReactionPopup($event, msg)">
                            <!-- TEXT -->
                            <div v-if="msg.text" class="message-text">
                                {{ msg.text }}
                            </div>
                            <!-- IMAGE (optional) -->
                            <img v-if="msg.attachmentUrl" :src="msg.attachmentUrl" class="message-image"
                                @click="openImage(msg.attachmentUrl)" />
                            <!-- SENT AT -->
                            <div class="message-meta">
                                <span class="time">
                                    {{ getMessageSentTime(msg.createdAt) }}
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
                <!-- ATTACHMENT PREVIEW -->
                <div v-if="attachment" class="attachment-preview">
                    <img class="preview-img" :src="attachmentUrl" />
                    <button class="icon-btn" @click="removeAttachment">
                        <img :src="deleteIcon" class="icon-img">
                    </button>
                </div>
                <!-- INPUT BAR -->
                <div class="input-bar">
                    <!-- UPLOAD BUTTON -->
                    <button class="icon-btn" @click="$refs.fileInput.click()">
                        <img :src="uploadIcon" class="icon-img">
                    </button>
                    <input ref="fileInput" type="file" accept="image/*" style="display: none" @change="addAttachment" />
                    <!-- PROMPT -->
                    <input class="prompt" v-model="text" placeholder="Type a message..." @keyup.enter="sendMessage" />
                    <!-- SEND BUTTON -->
                    <button class="icon-btn" @click="sendMessage">
                        <img :src="sendIcon" class="icon-img">
                    </button>
                </div>
            </div>
        </div>
        <!-- REACTION POPUP -->
        <div v-if="reactionPopup.visible" class="reaction-popup"
            :style="{ top: reactionPopup.y + 'px', left: reactionPopup.x + 'px' }">
            <span @click="react('like')">
                <img :src="likeIcon" class="icon-img">
            </span>
            <span @click="react('love')">
                <img :src="loveIcon" class="icon-img">
            </span>
            <span @click="react('laugh')">
                <img :src="laughIcon" class="icon-img">
            </span>
            <span @click="react('sad')">
                <img :src="sadIcon" class="icon-img">
            </span>
            <span @click="react('angry')">
                <img :src="angryIcon" class="icon-img">
            </span>
        </div>
        <div v-if="showImageModal" class="image-modal" @click="closeImage">
            <img :src="zoomedImage" class="image-modal-content" />
        </div>
        <LoggedAs />
    </div>
</template>

<style scoped>
.topbar,
.chat-container {
    position: relative;
    z-index: 1;
}

.chat-header {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    margin-bottom: 0px;

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

.chat-container {
    margin-top: 0px;
    padding: 18px;
    display: flex;
    justify-content: center;
    height: calc(100vh - 50px - 60px);
    /* topbar + loggedAs approx */
    box-sizing: border-box;
}

.chat-box {
    width: min(720px, 100%);
    height: 100%;
    display: flex;
    flex-direction: column;
    border-radius: 22px;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);
}

.messages {
    width: min(720px, 100%);
    overflow-x: hidden;
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding: 0px 18px;
    position: relative;
}

.messages::before {
    content: "";
    position: absolute;
    inset: -1px;
    border-radius: 22px;
    background: linear-gradient(120deg,
            rgba(0, 255, 120, 0.14),
            rgba(255, 255, 255, 0.05),
            rgba(0, 255, 120, 0.14));

    -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.25;
    pointer-events: none;
    animation: borderGlow 4s ease-in-out infinite alternate;
}

@keyframes borderGlow {
    0% {
        opacity: 0.18;
    }

    100% {
        opacity: 0.35;
    }
}

.message {
    display: flex;
    margin: 6px 0;
    animation: fadeInUp 0.35s ease both;
}

.message.mine {
    justify-content: flex-end;
}

.message-image {
    width: 220px;
    height: 160px;

    border-radius: 12px;
    margin-top: 8px;

    object-fit: cover;
    /* key part */
    display: block;
}

.message-text {
    font-size: 0.95rem;
    line-height: 1.35;
    word-wrap: break-word;
    white-space: pre-wrap;
    padding-right: 40px;
    /* spazio per l'orario a destra */
}

.message-content {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.message-meta {
    display: flex;
    justify-content: flex-end;
    margin-top: 2px;
}

.time {
    font-size: 0.72rem;
    opacity: 0.55;
    margin-left: auto;
}

.bubble {
    max-width: 70%;
    padding: 10px 14px;

    border-radius: 18px;

    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);

    color: rgba(245, 245, 245, 0.92);

    backdrop-filter: blur(10px);

    position: relative;
    overflow: hidden;
}

.message.mine .bubble {
    background: rgba(0, 255, 120, 0.12);
    border: 1px solid rgba(0, 255, 120, 0.2);
}

.input-bar {
    width: min(720px, 100%);
    margin: 0px auto 0;
    display: flex;
    gap: 10px;
    padding: 10px;
    border-radius: 0px;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
}

.input-bar input {
    flex: 1;
    padding: 10px 12px;
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(0, 0, 0, 0.2);
    color: white;
    outline: none;
}

.prompt {
    width: 100%;
    padding: 14px;

    border-radius: 14px;

    background: rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.06);

    color: rgba(245, 245, 245, 0.92);
    outline: none;

    transition: all 0.25s ease;
}

.prompt:focus {
    border: 1px solid rgba(0, 255, 120, 0.25);
    box-shadow: 0 0 20px rgba(0, 255, 120, 0.08);
}

.input-bar button {
    padding: 10px 14px;
    border-radius: 12px;
    border: 1px solid rgba(0, 255, 120, 0.2);
    background: rgba(0, 255, 120, 0.12);
    color: white;
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.input-bar button:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(0, 255, 120, 0.15);
}

.input-bar button:active {
    transform: scale(0.97);
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(12px) scale(0.98);
    }

    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

@media (max-width: 600px) {
    .messages {
        height: 65vh;
        padding: 14px;
    }

    .bubble {
        max-width: 85%;
    }

    .input-bar {
        padding: 10px;
    }
}

.attachment-preview {
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

.preview-img {
    width: 50px;
    height: 50px;
    border-radius: 12px;
    object-fit: cover;
    border: 1px solid rgba(255, 255, 255, 0.12);
}

.image-modal {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.85);

    display: flex;
    align-items: center;
    justify-content: center;

    z-index: 9999;
    cursor: zoom-out;
}

.image-modal-content {
    max-width: 90%;
    max-height: 90%;

    border-radius: 12px;
    box-shadow: 0 20px 80px rgba(0, 0, 0, 0.6);

    transform: scale(1);
    animation: zoomIn 0.2s ease;
}

@keyframes zoomIn {
    from {
        transform: scale(0.85);
        opacity: 0;
    }

    to {
        transform: scale(1);
        opacity: 1;
    }
}

.reaction-popup {
    position: fixed;
    display: flex;
    gap: 8px;

    padding: 8px 10px;

    background: rgba(20, 20, 20, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 14px;

    z-index: 99999;

    transform: translate(-50%, -120%);
    animation: popIn 0.15s ease;
}

.reaction-popup span {
    font-size: 18px;
    cursor: pointer;
    transition: transform 0.15s ease;
}

.reaction-popup span:hover {
    transform: scale(1.3);
}

@keyframes popIn {
    from {
        opacity: 0;
        transform: translate(-50%, -100%) scale(0.8);
    }

    to {
        opacity: 1;
        transform: translate(-50%, -120%) scale(1);
    }
}
</style>
