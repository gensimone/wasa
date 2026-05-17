import { reactive } from "vue"

export const chat = reactive({
    userId: localStorage.getItem("chatUserId")
        ? Number(localStorage.getItem("chatUserId"))
        : null,
    name: localStorage.getItem("chatName") || null,
    photoUrl: localStorage.getItem("chatPhotoUrl") || null,
})

export function setChatName(name) {
    chat.name = name;
    localStorage.setItem("chatName", name);
}

export function setChatPhotoUrl(photoUrl) {
    chat.photoUrl = `${__API_URL__}${photoUrl}`
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
}
