import api from "@/services/axios"
import { user } from "@/state/user"
import { fetchUser } from "@/services/usersService"

class ConversationsService {
    constructor() {
        this.poller = null
        this.lock = false
    }

    async fetchConversationIds() {
        const response = await api.get(`/conversations`, {
            headers: { Authorization: user.userId }
        })

        return response.data.conversations || []
    }

    async fetchGroup(groupId) {
        const response = await api.get(`/groups/${groupId}`, {
            headers: { Authorization: user.userId }
        })

        return response.data
    }

    async fetchLastMessage(conversationId) {
        const response = await api.get(`/conversations/${conversationId}/last`, {
            headers: { Authorization: user.userId }
        })

        return response.data
    }

    async fetchMembers(conversationId) {
        const response = await api.get(`/conversations/${conversationId}/members`, {
            headers: { Authorization: user.userId }
        })

        return response.data.userIds
    }

    async fetchConversations() {
        const conversationIds = await this.fetchConversationIds()

        const conversations = []
        for (const { conversationId, isGroup } of conversationIds) {
            const lastMessage = await this.fetchLastMessage(conversationId)
            if (isGroup) {
                const group = await this.fetchGroup(conversationId)
                conversations.push({
                    id: conversationId,
                    name: group.name,
                    photoUrl: group.photoUrl,
                    lastMessage: lastMessage,
                    isGroup: true
                })
            } else {
                const members = await this.fetchMembers(conversationId)
                let member = members[0]
                if (member === user.userId) member = members[1]
                const otherUser = await fetchUser(member)
                conversations.push({
                    id: member,
                    name: otherUser.name,
                    photoUrl: otherUser.photoUrl,
                    lastMessage: lastMessage,
                    isGroup: false
                })
            }
        }

        return conversations
    }

    startPolling(callback, interval = 3000) {
        this.stopPolling()

        this.poller = setInterval(async () => {
            if (this.lock) return
            this.lock = true

            try {
                callback({
                    conversations: await this.fetchConversations()
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

export default ConversationsService
