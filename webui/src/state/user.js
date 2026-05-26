import Poller from "@/services/poller"
import { reactive } from "vue"
import { getUserById } from "@/services/users"
import { defaultUserPhotoUrl } from "@/assets/default"

export const user = reactive({
    userId: null,
    name: null,
    photoUrl: defaultUserPhotoUrl,
    poller: null
})

export function updateUserState(partial = {}) {
    if (partial.userId !== undefined) {
        const id = partial.userId !== null ? Number(partial.userId) : null
        user.userId = Number.isNaN(id) ? null : id

        if (user.userId !== null) {
            localStorage.setItem("userId", user.userId)
        } else {
            localStorage.removeItem("userId")
        }
    }

    if (partial.name !== undefined) {
        user.name = partial.name

        if (partial.name !== null) {
            localStorage.setItem("name", partial.name)
        } else {
            localStorage.removeItem("name")
        }
    }

    if (partial.photoUrl !== undefined) {
        user.photoUrl = partial.photoUrl || defaultUserPhotoUrl
        localStorage.setItem("photoUrl", user.photoUrl)
    }
}

export function hydrateUserState() {
    const rawUserId = localStorage.getItem("userId")

    updateUserState({
        userId: rawUserId ? Number(rawUserId) : null,
        name: localStorage.getItem("name"),
        photoUrl: localStorage.getItem("photoUrl") || defaultUserPhotoUrl
    })
}

export function clearUserState() {
    stopPollingUser()

    updateUserState({
        userId: null,
        name: null,
        photoUrl: defaultUserPhotoUrl
    })

    localStorage.removeItem("userId")
    localStorage.removeItem("name")
    localStorage.removeItem("photoUrl")
}

export async function fetchUserState() {
    if (!user.userId) throw new Error("Cannot fetch user: missing userId")

    const fetchedUser = await getUserById(user.userId)

    updateUserState({
        name: fetchedUser.name,
        photoUrl: fetchedUser.photoUrl
    })
}

export function startPollingUser(interval = 5000) {
    stopPollingUser()

    user.poller = new Poller(async () => {
        await fetchUserState()
    }, interval)

    user.poller.startPolling()
}

export function stopPollingUser() {
    if (user.poller) {
        user.poller.stopPolling()
        user.poller = null
    }
}
