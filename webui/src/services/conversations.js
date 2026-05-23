import api from "@/services/axios"
import { user } from "@/state/user"

export async function getMemberIds(conversationId) {
    const response = await api.get(`/conversations/${conversationId}/members`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.userIds
}

export async function getLastMessage(conversationId) {
    const response = await api.get(`/conversations/${conversationId}/last`,
        { headers: { Authorization: user.userId } }
    )

    return response.data
}

export async function getConversation(conversationId) {
    const response = await api.get(`/conversations/${conversationId}`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.messageIds
}

export async function getMyConversations() {
    const response = await api.get(`/conversations`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.conversations
}

export async function sendMessageToConversation(conversationId, text, attachment, mediaType = "image") {
    const formData = new FormData()

    formData.append("text", text)

    if (attachment) {
        formData.append("file", attachment)
        formData.append("mediaType", mediaType)
    }

    const response = await api.post(`/conversations/${conversationId}/message`,
        formData,
        {
            headers: {
                Authorization: user.userId,
                "Content-Type": "multipart/form-data"
            }
        }
    )

    return response.data
}

export async function forwardMessageToConversation(conversationId, messageId) {
    const response = await api.post(`conversations/${conversationId}/fmessage`,
        { messageId: messageId },
        { headers: { Authorization: user.userId } }
    )

    return response.data
}
