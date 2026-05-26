import { watch } from "vue"
import { userConversations, groupConversations } from "@/state/conversations"
import { getNotifier } from "@/notifier"
import router from "@/router"

let previousUserConversations = new Map()
let previousGroupConversations = new Map()

export function startMessageNotifier() {
    watch(
        userConversations,
        (newMap) => {
            handleNewMessages(newMap, previousUserConversations, true)
        },
        { deep: true }
    )

    watch(
        groupConversations,
        (newMap) => {
            handleNewMessages(newMap, previousGroupConversations, false)
        },
        { deep: true }
    )
}

function matchConversation(id, direct) {
    const route = router.currentRoute.value

    const isDirect = route.query.direct === "true"
    const inConversation = route.name === "conversation"
    const idConversation = route.params.id

    return (
        inConversation &&
        isDirect === direct &&
        idConversation == id
    )
}

function handleNewMessages(newMap, oldMap, direct) {
    for (const [id, conversationData] of newMap.entries()) {
        if (matchConversation(id, direct)) continue

        const oldConversation = oldMap.get(id)

        if (!oldConversation) continue

        const oldMessageIds = (oldConversation.messages || [])
            .map(m => m.messageId)

        const newMessages = (conversationData.messages || [])
            .filter(m => !oldMessageIds.includes(m.messageId))

        newMessages.forEach(msg => {
            getNotifier()?.message({
                text: msg.text,
                id: id,
                isGroup: !direct,
                thumbnailUrl: conversationData.photoUrl,
                attachmentUrl: msg.attachmentUrl
            })
        })
    }

    if (direct) {
        previousUserConversations = new Map(newMap)
    } else {
        previousGroupConversations = new Map(newMap)
    }
}
