import api from "@/services/axios"
import { user } from "@/state/user"
import { chat } from "@/state/chat"

class ChatService {
    constructor() {
        this.poller = null
        this.lock = false
    }

    async fetchMessage(messageId) {
        const res = await api.get(`/messages/${messageId}`, {
            headers: { Authorization: user.userId }
        })
        return res.data
    }

    async fetchMessages() {
        const res = await api.get(`/users/${chat.userId}/messages`, {
            headers: { Authorization: user.userId }
        })

        const messages = await Promise.all(
            res.data.messageIds.map(id => this.fetchMessage(id))
        )

        messages.sort(
            (a, b) => new Date(a.createdAt) - new Date(b.createdAt)
        )

        return messages
    }

    async sendMessage(text, attachment) {
        const formData = new FormData()
        formData.append("text", text)

        if (attachment) {
            formData.append("file", attachment)
            formData.append("mediaType", "image")
        }

        return await api.post(
            `/users/${chat.userId}/message`,
            formData,
            {
                headers: {
                    Authorization: user.userId,
                    "Content-Type": "multipart/form-data"
                }
            }
        )
    }

    startPolling(callback, interval = 3000) {
        this.stopPolling()

        this.poller = setInterval(async () => {
            if (this.lock) return
            this.lock = true

            try {
                callback({
                    messages: await this.fetchMessages()
                })
            } catch (e) {
                console.error("Polling error:", e)
            } finally {
                this.lock = false
            }
        }, interval)
    }

    stopPolling() {
        if (this.poller) {
            clearInterval(this.poller)
            this.poller = null
        }
    }
}

export default ChatService
