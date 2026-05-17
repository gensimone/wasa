<script>
import { user } from "@/state/user.js"
import { chat } from "@/state/chat.js"
import backIcon from "@/assets/icons/back.svg"
import homeIcon from "@/assets/icons/home.svg"
import LoggedAs from "@/components/LoggedAs.vue"
export default {
    components: {
        LoggedAs
    },

    data() {
        return {
            userId: user.userId, // Logged user userId (used to determine css in messages)
            messages: [],        // Conversation messages.
            text: null,          // Message text to send.
            loading: false,
            error: null,

            backIcon,
            homeIcon,
        }
    },
    methods: {
        async sendMessage() {
            this.loading = true
            try {
                const formData = new FormData()
                formData.append("text", this.text)

                const response = await this.$axios.post(`/users/${chat.userId}/message`,
                    formData,
                    {
                        headers: {
                            Authorization: user.userId,
                            "Content-Type": "multipart/form-data"
                        }
                    }
                )

                const data = response.data
                this.messages.push({
                    messageId: data.messageId,
                    text: data.text,
                    senderId: data.senderId,
                    conversationid: data.conversationId,
                    createdAt: data.createdAt,
                    isForwarded: data.isForwarded,
                    commentTo: data.commentTo,
                    attachmentUrl: data.attachmentUrl,
                    mediaType: data.mediaType
                })

                this.text = null

            } catch (e) {
                this.message = e?.response?.data?.error || "Unexpected error"
            }

            this.loading = false
        },

        async fetchMessage(messageId) {
            try {
                console.log("Message fetch. ID:", messageId)
                const response = await this.$axios.get(`/messages/${messageId}`, {
                    headers: { Authorization: user.userId }
                })
                console.log("Message response:", response.data)
                return response.data
            } catch (e) {
                console.log(e)
            }
        },

        async fetchMessages() {
            this.loading = true
            this.messages = []

            try {
                console.log("Getting message of conversation between ID:", chat.userId)
                const response = await this.$axios.get(`/users/${chat.userId}/messages`, {
                    headers: { Authorization: user.userId }
                })

                console.log("Message Ids:", response.data.messageIds)

                const messagePromises = response.data.messageIds.map(messageId =>
                    this.fetchMessage(messageId)
                )

                const messagesData = await Promise.all(messagePromises)

                this.messages = messagesData
                    .map(data => ({
                        messageId: Number(data.messageId),
                        text: data.text,
                        senderId: Number(data.senderId),
                        conversationid: Number(data.conversationId),
                        createdAt: data.createdAt,
                        isForwarded: Boolean(data.isForwarded),
                        commentTo: Number(data.commentTo),
                        attachmentUrl: data.attachmentUrl,
                        mediaType: data.mediaType
                    }))

            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
                this.messages = []
            }
            this.loading = false
        },

        getImageUrl(url) {
            return url ? `${__API_URL__}${url}` : null;
        }
    },
    mounted() {
        this.fetchMessages()
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
                <!-- MESSAGES -->
                <div class="messages">
                    <div v-for="msg in messages" :key="msg.messageId" class="message"
                        :class="{ mine: msg.senderId === userId }">
                        <div class="bubble">
                            <!-- TEXT -->
                            <div v-if="msg.text">
                                {{ msg.text }}
                            </div>
                            <!-- IMAGE (optional) -->
                            <img v-if="msg.attachmentUrl" :src="getImageUrl(msg.attachmentUrl)" class="message-image" />
                            <!-- META (time + status) -->
                            <!-- <div class="message-meta"> -->
                            <!--     <span class="time"> -->
                            <!--         {{ formatTime(msg.sentAt) }} -->
                            <!--     </span> -->
                            <!--     <span v-if="msg.senderId === userId" class="status"> -->
                            <!--         {{ statusIcon(msg.status) }} -->
                            <!--     </span> -->
                            <!-- </div> -->
                        </div>
                    </div>
                </div>
                <!-- INPUT BAR -->
                <div class="input-bar">
                    <input v-model="text" placeholder="Type a message..." @keyup.enter="sendMessage" />
                    <button @click="sendMessage">Send</button>
                </div>
            </div>
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

.chat-container {
    margin-top: 10px;
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
    height: 70vh;
    display: flex;
    flex-direction: column;
    border-radius: 22px;
    padding: 18px;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);
    position: relative;
    overflow-y: auto;
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

    transition: transform 0.25s ease, box-shadow 0.25s ease;
}

.bubble:hover {
    transform: translateY(-3px);
    box-shadow: 0 18px 60px rgba(0, 255, 120, 0.08);
}

.message.mine .bubble {
    background: rgba(0, 255, 120, 0.12);
    border: 1px solid rgba(0, 255, 120, 0.2);
}

.bubble::before {
    content: "";
    position: absolute;
    top: -40%;
    left: -90%;
    width: 130%;
    height: 200%;

    background: linear-gradient(90deg,
            transparent,
            rgba(0, 255, 120, 0.08),
            rgba(255, 255, 255, 0.05),
            transparent);

    transform: rotate(18deg);
    transition: left 0.6s ease;
}

.bubble:hover::before {
    left: 70%;
}

.input-bar {
    width: min(720px, 100%);
    margin: 12px auto 0;
    display: flex;
    gap: 10px;
    padding: 12px;
    border-radius: 18px;
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
</style>
