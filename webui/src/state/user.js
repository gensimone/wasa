import { reactive } from "vue"
import { getUserById } from "@/services/users"
import { Poller } from "@/services/poller"

export const defaultUserPhotoUrl = "/media/default-user-photo.jpg"

export const user = reactive({
    userId: null,
    name: null,
    photoUrl: null,
    poller: null
})

export function loadUserState() {
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

export async function fetchUserState() {
    if (!user.userId) throw new Error("Cannot fetch user")

    const fetchedUser = await getUserById(user.userId)
    if (fetchedUser.name !== user.name) setName(fetchedUser.name)
    if (fetchedUser.photoUrl !== user.photoUrl) setPhotoUrl(fetchedUser.photoUrl)
}

export function startPollingUser() {
    user.poller?.stopPolling()
    user.poller = new Poller(async () => {
        await fetchUserState()
    })

    user.poller.startPolling()
}
