import api from "@/services/axios"
import { user } from "@/state/user"

export async function getMessage(messageId) {
    const response = await api.get(`/messages/${messageId}`,
        { headers: { Authorization: user.userId } }
    )
    return response.data
}

export async function deleteMessage(messageId) {
    await api.delete(`/messages/${messageId}`,
        { headers: { Authorization: user.userId } }
    )
}

export async function setMessageStatusAsRead(messageId) {
    await api.put(`/messages/${messageId}/receipts`,
        { headers: { Authorization: user.userId } }
    )
}

export async function getReceipts(messageId) {
    const response = await api.get(`/messages/${messageId}/receipts`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.receipts
}
