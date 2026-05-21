import { reactive } from "vue"

export const conversation = reactive({
    id: null,        // User or Conversation ID
    name: null,      // User or Group name
    photoUrl: null,  // User or Group photo URL
    isGroup: null    // Is it a group?
})

export function loadConversationMetadataFromStorage() {
    conversation.id = Number(localStorage.getItem("conversationId"))
    conversation.name = localStorage.getItem("conversationName")
    conversation.photoUrl = localStorage.getItem("conversationPhotoUrl")
    conversation.isGroup = localStorage.getItem("conversationIsGroup")
}

export function setConversationName(name) {
    conversation.name = name;
    localStorage.setItem("conversationName", name);
}

export function setConversationPhotoUrl(photoUrl) {
    conversation.photoUrl = photoUrl
    localStorage.setItem("conversationPhotoUrl", conversation.photoUrl);
}

export function setConversationId(id) {
    conversation.id = Number(id);
    localStorage.setItem("conversationid", id);
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
