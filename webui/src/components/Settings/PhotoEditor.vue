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
    <div class="avatar-row">

        <!-- DELETE / REVERT / PLACEHOLDER BUTTON -->
        <button v-if="photoChanged" type="button" class="icon-btn" :disabled="loading" @click="$emit('revertPhoto')">
            <img src="/icons/revert.svg" class="icon-img">
        </button>
        <button v-else-if="!isDefault" type="button" class="icon-btn" :disabled="loading" @click="$emit('deletePhoto')">
            <img src="/icons/remove.svg" class="icon-img">
        </button>
        <button v-else type="button" class="icon-btn invisible-placeholder" disabled>
            <img src="/icons/remove.svg" class="icon-img">
        </button>

        <!-- AVATAR -->
        <label class="avatar-wrapper">
            <img :src="photoUrl" class="avatar-big" />
            <input type="file" accept="image/*" hidden @change="$emit('uploadPhoto', $event)" />
        </label>

        <!-- PLACEHOLDER BUTTON -->
        <button type="button" class="icon-btn invisible-placeholder" disabled>
            <img src="/icons/revert.svg" class="icon-img">
        </button>

    </div>
</template>
