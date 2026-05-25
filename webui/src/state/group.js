import { reactive } from "vue"
import { Poller } from "@/services/poller"
import { getGroup } from "@/services/groups"

export const group = reactive({
    conversationId: null,
    founderId: null,
    name: null,
    photoUrl: null,
    createdAt: null,
    poller: null
})

export function updateGroupState(partial = {}) {
    if (partial.conversationId !== undefined) {
        const id = partial.conversationId !== null ? Number(partial.conversationId) : null
        group.conversationId = Number.isNaN(id) ? null : id

        if (group.conversationId !== null) {
            localStorage.setItem("groupConversationId", String(group.conversationId))
        } else {
            localStorage.removeItem("groupConversationId")
        }
    }

    if (partial.founderId !== undefined) {
        const id = partial.founderId !== null ? Number(partial.founderId) : null
        group.founderId = Number.isNaN(id) ? null : id

        if (group.founderId !== null) {
            localStorage.setItem("groupFounderId", String(group.founderId))
        } else {
            localStorage.removeItem("groupFounderId")
        }
    }

    if (partial.name !== undefined) {
        group.name = partial.name

        if (partial.name !== null) {
            localStorage.setItem("groupName", partial.name)
        } else {
            localStorage.removeItem("groupName")
        }
    }

    if (partial.photoUrl !== undefined) {
        group.photoUrl = partial.photoUrl

        if (partial.photoUrl !== null) {
            localStorage.setItem("groupPhotoUrl", partial.photoUrl)
        } else {
            localStorage.removeItem("groupPhotoUrl")
        }
    }

    if (partial.createdAt !== undefined) {
        group.createdAt = partial.createdAt

        if (partial.createdAt !== null) {
            localStorage.setItem("groupCreatedAt", partial.createdAt)
        } else {
            localStorage.removeItem("groupCreatedAt")
        }
    }
}

export function hydrateGroupState() {
    const rawConversationId = localStorage.getItem("groupConversationId")
    const rawFounderId = localStorage.getItem("groupFounderId")

    updateGroupState({
        conversationId: rawConversationId ? Number(rawConversationId) : null,
        founderId: rawFounderId ? Number(rawFounderId) : null,
        name: localStorage.getItem("groupName"),
        photoUrl: localStorage.getItem("groupPhotoUrl"),
        createdAt: localStorage.getItem("groupCreatedAt")
    })
}

export function clearGroupState() {
    updateGroupState({
        conversationId: null,
        founderId: null,
        name: null,
        photoUrl: null,
        createdAt: null
    })

    localStorage.removeItem("groupConversationId")
    localStorage.removeItem("groupFounderId")
    localStorage.removeItem("groupName")
    localStorage.removeItem("groupPhotoUrl")
    localStorage.removeItem("groupCreatedAt")
}

export async function fetchGroupState() {
    if (!group.conversationId) {
        throw new Error("Cannot fetch group: missing conversationId")
    }

    const fetchedGroup = await getGroup(group.conversationId)

    updateGroupState({
        founderId: fetchedGroup.founderId,
        name: fetchedGroup.name,
        photoUrl: fetchedGroup.photoUrl,
        createdAt: fetchedGroup.createdAt
    })
}

export function startPollingGroup(interval = 5000) {
    stopPollingGroup()

    group.poller = new Poller(async () => {
        await fetchGroupState()
    }, interval)

    group.poller.startPolling()
}

export function stopPollingGroup() {
    if (group.poller) {
        group.poller.stopPolling()
        group.poller = null
    }
}
