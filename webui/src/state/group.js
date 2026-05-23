import { reactive } from "vue"

export const defaultGroupPhotoUrl = "/media/default-group-photo.jpg"

export const group = reactive({
    founderId: null,
    createdAt: null
})

export function loadGroupState() {
    group.founderId = Number(localStorage.getItem("groupFounderId"))
    group.createdAt = localStorage.getItem("groupCreatedAt")
}

export function setGroupFounderId(founderId) {
    group.founderId = Number(founderId)
    localStorage.setItem("groupFounderId", founderId)
}

export function setGroupCreatedAt(createdAt) {
    group.createdAt = createdAt
    localStorage.setItem("groupCreatedAt", createdAt)
}

export function clearGroupState() {
    group.founderId = null
    group.createdAt = null

    localStorage.removeItem("groupFounderId")
    localStorage.removeItem("groupCreatedAt")
}
