import { reactive } from "vue"

export const chat = reactive({
    userId: null,
    name: null,
    photoUrl: null
})

export function loadChatMetadataFromStorage() {
    chat.userId = Number(localStorage.getItem("chatUserId"))
    chat.name = localStorage.getItem("chatName")
    chat.photoUrl = localStorage.getItem("chatPhotoUrl")
}

export function setChatName(name) {
    chat.name = name;
    localStorage.setItem("chatName", name);
}

export function setChatPhotoUrl(photoUrl) {
    chat.photoUrl = photoUrl
    localStorage.setItem("chatPhotoUrl", chat.photoUrl);
}

export function setChatUserId(userId) {
    chat.userId = Number(userId);
    localStorage.setItem("chatUserId", userId);
}

export function clearChatState() {
    chat.userId = null;
    chat.name = null;
    chat.photoUrl = null;

    localStorage.removeItem("chatUserId")
    localStorage.removeItem("chatName")
    localStorage.removeItem("chatPhotoUrl")
}
