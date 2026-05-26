import { reactive } from "vue";

const imageModal = reactive({
    visible: false,
    photoUrl: null
})

export function setImageModal(photoUrl) {
    imageModal.visible = true
    imageModal.photoUrl = photoUrl
}

export default imageModal
