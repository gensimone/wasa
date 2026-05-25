<script>
export default {
    props: {
        photoUrl: String,
        photoChanged: Boolean,
        loading: Boolean
    },

    emits: [
        "revertPhoto",
        "deletePhoto",
        "uploadPhoto"
    ],

    computed: {
        isDefault() {
            const filename = this.photoUrl.split("/").pop()
            return filename === "default-user-photo.jpg" || filename === "default-group-photo.jpg"
        },
    },

}
</script>

<template>
    <div class="photo-editor">
        <button v-if="photoChanged" type="button" class="icon-btn" :disabled="loading" @click="$emit('revertPhoto')">
            <img src="/icons/revert.svg" class="icon-img">
        </button>
        <button v-else-if="!isDefault" type="button" class="icon-btn" :disabled="loading" @click="$emit('deletePhoto')">
            <img src="/icons/remove.svg" class="icon-img">
        </button>
        <button v-else type="button" class="icon-btn invisible-placeholder" disabled>
            <img src="/icons/remove.svg" class="icon-img">
        </button>

        <label class="photo-editor-photo-wrapper">
            <img :src="photoUrl" class="avatar-big" />
            <input type="file" accept="image/*" hidden @change="$emit('uploadPhoto', $event)" />
        </label>

        <button type="button" class="icon-btn invisible-placeholder" disabled>
            <img src="/icons/revert.svg" class="icon-img">
        </button>
    </div>
</template>

<style scoped>
.photo-editor {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}

.photo-editor-photo-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.invisible-placeholder {
    visibility: hidden;
}
</style>
