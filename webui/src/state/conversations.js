import { ref } from "vue"
import { Poller } from "@/services/poller"
import { getUserById, getConversationByUserId } from "@/services/users"
import { getMessage } from "@/services/messages"
import { getGroup } from "@/services/groups"
import { getMyConversations, getConversation } from "@/services/conversations"

export const userConversations = ref(new Map())
export const groupConversations = ref(new Map())

let poller = null

const CACHE_TTL = 30_000
const messageCache = new Map()

async function fetchMessageCached(id) {
    const cached = messageCache.get(id)

    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
        return cached.data
    }

    const msg = await getMessage(id)

    messageCache.set(id, {
        data: msg,
        timestamp: Date.now()
    })

    return msg
}

export async function getConversationMessages(id, isGroup) {
    const messageIds = isGroup
        ? await getConversation(id)
        : await getConversationByUserId(id)

    if (!messageIds?.length) return []

    const messages = await Promise.all(
        messageIds.map(fetchMessageCached)
    )

    return messages

    // return messages.sort(
    //   (a, b) => new Date(a.createdAt) - new Date(b.createdAt)
    // )
}

async function fetchGroupConversation(id) {
    const [messages, group] = await Promise.all([
        getConversationMessages(id, true),
        getGroup(id)
    ])

    return [
        id,
        {
            id: id,
            founderId: group.founderId,
            name: group.name,
            photoUrl: group.photoUrl,
            createdAt: group.createdAt,
            isGroup: true,
            messages
        }
    ]
}

async function fetchUserConversation(id) {
    const [messages, user] = await Promise.all([
        getConversationMessages(id, false),
        getUserById(id)
    ])

    return [
        id,
        {
            id: user.userId,
            name: user.name,
            photoUrl: user.photoUrl,
            isGroup: false,
            messages
        }
    ]
}

export async function fetchConversations() {
    const fetched = (await getMyConversations()) || []

    const groupIds = fetched.filter(c => c.isGroup).map(c => c.id)
    const userIds = fetched.filter(c => !c.isGroup).map(c => c.id)

    const [groups, users] = await Promise.all([
        Promise.all(groupIds.map(fetchGroupConversation)),
        Promise.all(userIds.map(fetchUserConversation))
    ])

    groupConversations.value = new Map(groups)
    userConversations.value = new Map(users)
}

export function startPollingConversations(interval = 5000) {
    stopPollingConversations()

    poller = new Poller(async () => {
        if (document.visibilityState !== "visible") return
        await fetchConversations()
    }, interval)

    poller.startPolling()
}

export function stopPollingConversations() {
    if (poller) {
        poller.stopPolling()
        poller = null
    }
}
