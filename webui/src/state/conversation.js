import { reactive } from "vue"
import { Poller } from "@/services/poller"
import { getGroup } from "@/services/groups"
import { getUserById } from "@/services/users"

export const conversation = reactive({
    id: null,
    name: null,
    photoUrl: null,
    isGroup: null,
    poller: null
})

export function loadConversationState() {
    conversation.id = Number(localStorage.getItem("conversationId"))
    conversation.name = localStorage.getItem("conversationName")
    conversation.photoUrl = localStorage.getItem("conversationPhotoUrl")
    conversation.isGroup = localStorage.getItem("conversationIsGroup") === "true"
}

export function setConversationName(name) {
    conversation.name = name;
    localStorage.setItem("conversationName", name);
}

export function setConversationPhotoUrl(photoUrl) {
    conversation.photoUrl = photoUrl
    localStorage.setItem("conversationPhotoUrl", photoUrl);
}

export function setConversationId(id) {
    conversation.id = Number(id);
    localStorage.setItem("conversationId", id);
}

export function setConversationIsGroup(isGroup) {
    conversation.isGroup = isGroup
    localStorage.setItem("conversationIsGroup", isGroup);
}

export function clearConversationState() {
    conversation.id = null;
    conversation.name = null;
    conversation.photoUrl = null;
    conversation.isGroup = null;

    localStorage.removeItem("conversationId")
    localStorage.removeItem("conversationName")
    localStorage.removeItem("conversationPhotoUrl")
    localStorage.removeItem("conversationIsGroup")
}

export async function fetchConversationState() {
    if (!conversation.id) throw new Error("Cannot fetch conversation")

    let fetchedConversation
    if (conversation.isGroup) {
        fetchedConversation = await getGroup(conversation.id)
    } else {
        fetchedConversation = await getUserById(conversation.id)
    }

    if (fetchedConversation.name !== conversation.name)
        setConversationName(fetchedConversation.name)
    if (fetchedConversation.photoUrl !== conversation.photoUrl)
        setConversationPhotoUrl(fetchedConversation.photoUrl)
}

export function startPollingConversation() {
    conversation.poller?.stopPolling()
    conversation.poller = new Poller(async () => {
        await fetchConversationState()
    })

    conversation.poller.startPolling()
}
