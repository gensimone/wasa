import { reactive } from "vue"

export const user = reactive({
    userId: localStorage.getItem("userId")
        ? Number(localStorage.getItem("userId"))
        : null,
    name: localStorage.getItem("name") || null,
    photoUrl: localStorage.getItem("photoUrl") || null,
})

export function setName(name) {
    user.name = name;
    localStorage.setItem("name", name);
}

export function setPhotoUrl(photoUrl) {
    user.photoUrl = `${__API_URL__}${photoUrl}`
    localStorage.setItem("photoUrl", user.photoUrl);
}

export function setUserId(userId) {
    user.userId = Number(userId);
    localStorage.setItem("userId", userId);
}

export function clearUserState() {
    user.userId = null;
    user.name = null;
    user.photoUrl = null;
}
