import { ref } from "vue"

export function usePhotoManager(initialUrl, defaultUrl) {
    const photoUrl = ref(initialUrl)
    const oldPhotoUrl = ref(null)
    const photo = ref(null)
    const photoChanged = ref(false)

    function revoke(url) {
        if (url && url.startsWith("blob:")) {
            URL.revokeObjectURL(url)
        }
    }

    function uploadPhoto(event) {
        const file = event.target.files?.[0]
        if (!file) return

        revoke(photoUrl.value)

        oldPhotoUrl.value = photoUrl.value
        photoUrl.value = URL.createObjectURL(file)

        photo.value = file
        photoChanged.value = true

        event.target.value = ""
    }

    function revertPhoto() {
        revoke(photoUrl.value)

        photoUrl.value = oldPhotoUrl.value
        photo.value = null
        photoChanged.value = false
    }

    function deletePhoto() {
        revoke(photoUrl.value)

        oldPhotoUrl.value = photoUrl.value
        photoUrl.value = defaultUrl

        photo.value = null
        photoChanged.value = true
    }

    return {
        photoUrl,
        oldPhotoUrl,
        photo,
        photoChanged,

        uploadPhoto,
        revertPhoto,
        deletePhoto,
        revoke
    }
}
