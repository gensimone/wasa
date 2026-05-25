import api from "@/services/axios"
import { user } from "@/state/user"

export async function addReaction(messageId, emoji) {
    const response = await api.post(`/reactions/${messageId}`,
        { emoji: emoji },
        { headers: { Authorization: user.userId } }
    )

    return response.data
}

export async function deleteReaction(messageId) {
    return await api.delete(`/reactions/${messageId}`,
        { headers: { Authorization: user.userId } }
    )
}

export async function getReactions(messageId) {
    const response = await api.get(`/reactions/${messageId}`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.reactions
}
