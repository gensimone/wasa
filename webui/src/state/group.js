import { reactive } from "vue"

export const group = reactive({
    founderId: localStorage.getItem("groupFounderId") || null,
    name: localStorage.getItem("groupName") || null,
    photoUrl: localStorage.getItem("groupPhotoUrl") || null,
})

export function setGroupName(name) {
    group.name = name;
    localStorage.setItem("groupName", name);
}

export function setGroupPhotoUrl(photoUrl) {
    group.photoUrl = `${__API_URL__}${photoUrl}`
    localStorage.setItem("groupPhotoUrl", group.photoUrl);
}

export function setGroupFounderId(founderId) {
    group.founderId = founderId;
    localStorage.setItem("groupFounderId", founderId);
}

export function clearGroupState() {
    group.founderId = null;
    group.name = null;
    group.photoUrl = null;
}
