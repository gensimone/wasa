import { reactive } from "vue"

export const user = reactive({
    userId: null,
    name: null,
    photoUrl: null
})

export function loadUserFromStorage() {
    user.userId = Number(localStorage.getItem("userId"))
    user.name = localStorage.getItem("name")
    user.photoUrl = localStorage.getItem("photoUrl")
}

export function setName(name) {
    user.name = name;
    localStorage.setItem("name", name);
}

export function setPhotoUrl(photoUrl) {
    user.photoUrl = photoUrl
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

    localStorage.removeItem("userId")
    localStorage.removeItem("name")
    localStorage.removeItem("photoUrl")
}
