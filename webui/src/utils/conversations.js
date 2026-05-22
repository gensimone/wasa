import {
    getUserById,
    getConversationByUserId,
    sendMessage as sendMessageToUser,
    forwardMessage as forwardMessageToUser
} from "@/services/users"

import {
    getMemberIds,
    getLastMessage,
    getConversation,
    getMyConversations,
    sendMessageToConversation,
    forwardMessageToConversation,
} from "@/services/conversations"

import { getMessage } from "@/services/messages"
import { getGroup } from "@/services/groups"
import { conversation } from "@/state/conversation"
import { user } from "@/state/user"

export async function getConversations() {
    const conversationIds = await getMyConversations() || []
    const conversations = []
    for (const { conversationId, isGroup } of conversationIds) {
        const lastMessage = await getLastMessage(conversationId)
        if (isGroup) {
            const group = await getGroup(conversationId)
            conversations.push({
                id: conversationId,
                name: group.name,
                photoUrl: group.photoUrl,
                lastMessage: lastMessage,
                isGroup: true,
                founderId: group.founderId,
                createdAt: group.createdAt
            })
        } else {
            const members = await getMemberIds(conversationId)
            let member = members[0]
            if (member === user.userId) member = members[1]
            const otherUser = await getUserById(member)
            conversations.push({
                id: member,
                name: otherUser.name,
                photoUrl: otherUser.photoUrl,
                lastMessage: lastMessage,
                isGroup: false,
                founderId: null,
                createdAt: null
            })
        }
    }

    return conversations
}

export async function getMembers(conversationId) {
    const memberIds = await getMemberIds(conversationId)
    return await Promise.all(
        memberIds.map(id => getUserById(id))
    )
}

export async function getConversationMessages() {
    let messageIds
    if (conversation.isGroup) {
        messageIds = await getConversation(conversation.id)
    } else {
        messageIds = await getConversationByUserId(conversation.id)
    }

    if (!messageIds) {
        return []
    }

    const messages = await Promise.all(
        messageIds.map(id => getMessage(id))
    )

    messages.sort(
        (a, b) => new Date(a.createdAt) - new Date(b.createdAt)
    )

    return messages
}

export async function sendMessage(text, photo) {
    if (conversation.isGroup) {
        return sendMessageToConversation(conversation.id, text, photo)
    } else {
        return sendMessageToUser(conversation.id, text, photo)
    }
}

export async function forwardMessage(messageId) {
    if (conversation.isGroup) {
        return forwardMessageToConversation(conversation.id, messageId)
    } else {
        return forwardMessageToUser(conversation.id, messageId)
    }
}
