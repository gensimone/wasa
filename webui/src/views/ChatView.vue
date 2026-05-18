<script>
import { user } from "@/state/user.js"
import { chat } from "@/state/chat.js"

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
            userId: user.userId,
            messages: [],
            lock: false,
            error: null,
            messageIds: [],
            poller: null,

            sending: false,
            sendError: null,

            fetchError: null,

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
                const formData = new FormData()
                formData.append("text", this.text)

                if (this.attachment) {
                    formData.append("file", this.attachment)
                    formData.append("mediaType", "image")
                }

                const response = await this.$axios.post(`/users/${chat.userId}/message`, formData, {
                    headers: {
                        Authorization: user.userId,
                        "Content-Type": "multipart/form-data"
                    }
                })

                const data = response.data
                this.messages.push({
                    messageId: Number(data.messageId),
                    text: data.text,
                    senderId: Number(data.senderId),
                    conversationid: Number(data.conversationId),
                    createdAt: data.createdAt,
                    isForwarded: Boolean(data.isForwarded),
                    commentTo: Number(data.commentTo),
                    attachmentUrl: data.attachmentUrl == "" ? "" : `${__API_URL__}${data.attachmentUrl}`,
                    mediaType: data.mediaType
                })

                this.messageIds.push(Number(data.messageId))

                this.text = null
                if (this.attachment)
                    this.removeAttachment()

                this.scrollToBottom()
            } catch (e) {
                this.sendError = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.sending = false
            }
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
                    newMessageIds.some(id => !oldIdsSet.has(id)) ||
                    this.messageIds.some(id => !newIdsSet.has(id))

                if (!changed) {
                    this.lock = false
                    return
                }

                const diffMessageIds = [...newIdsSet].filter(x => !oldIdsSet.has(x))

                const results = await Promise.all(
                    diffMessageIds.map(messageId => this.fetchMessage(messageId))
                )

                const newMessages = results.map(data => ({
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

            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.lock = false
            }
        },

        scrollToBottom() {
            this.$nextTick(() => {
                this.$refs.chatBox?.scrollToBottom()
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

        react(messageId, emoji) {
            // TODO: API call
            console.log("React to message:", messageId, emoji)
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
        <Topbar :actions="[
            { icon: '/icons/home.svg', onClick: () => $router.push('/home') },
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <div class="content">
            <ChatBox ref="chatBox" :chatName="chatName" :avatarUrl="avatarUrl" :messages="messages" :userId="userId"
                :text="text" :attachment="attachment" :attachmentUrl="attachmentUrl" @update:text="text = $event"
                @send="sendMessage" @addAttachment="addAttachment" @removeAttachment="removeAttachment"
                :sending="sending" :sendError="sendError" @react="react" />
        </div>
        <Bottombar />
    </div>
</template>
