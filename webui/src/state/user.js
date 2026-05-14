import { reactive } from "vue"

export const userState = reactive({
    userId: localStorage.getItem("userId") || null,
    name: localStorage.getItem("name") || null,
    photoUrl: localStorage.getItem("photoUrl") || null,
})

export function setName(name) {
    userState.name = name;
    localStorage.setItem("name", name);
}

export function setPhotoUrl(photoUrl) {
    userState.photoUrl = `${__API_URL__}${photoUrl}`
    localStorage.setItem("photoUrl", userState.photoUrl);
}

export function setUserId(userId) {
    userState.userId = userId;
    localStorage.setItem("userId", userId);
}

export function clearUserState() {
    userState.userId = null;
    userState.name = null;
    userState.photoUrl = null;
    localStorage.clear();
}
