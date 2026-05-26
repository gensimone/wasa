<script>
import { getIcon } from '@/state/theme';
import { expandUrl } from '@/utils/media';
import { defaultGroupPhotoUrl, defaultUserPhotoUrl } from '@/assets/default';
export default {
    props: {
        photoUrl: { type: String, required: true },
        photoChanged: { type: Boolean, required: true },
        loading: { type: Boolean, required: true },
        enableEditing: { type: Boolean, required: true }
    },

    emits: [
        "revertPhoto",
        "deletePhoto",
        "uploadPhoto"
    ],

    methods: {
        expandUrl,
        getIcon
    },

    computed: {
        isDefault() {
            return (
                this.photoUrl === defaultUserPhotoUrl ||
                this.photoUrl === defaultGroupPhotoUrl
            )
        },
    },

}
</script>

<template>
    <div class="photo-editor">
        <div v-if="enableEditing">
            <button v-if="photoChanged" type="button" class="icon-btn" :disabled="loading"
                @click="$emit('revertPhoto')">
                <img :src="getIcon('revert')" class="icon-img">
            </button>
            <button v-else-if="!isDefault" type="button" class="icon-btn" :disabled="loading"
                @click="$emit('deletePhoto')">
                <img :src="getIcon('trash')" class="icon-img">
            </button>
            <button v-else type="button" class="icon-btn invisible-placeholder" disabled>
                <img :src="getIcon('trash')" class="icon-img">
            </button>
        </div>

        <label class="photo-editor-photo-wrapper">
            <img :src="expandUrl(photoUrl)" class="avatar-big" />
            <input type="file" accept="image/*" :disabled="!enableEditing" hidden @change="$emit('uploadPhoto', $event)" />
        </label>

        <div v-if="enableEditing">
            <button type="button" class="icon-btn invisible-placeholder" disabled>
                <img :src="getIcon('revert')" class="icon-img">
            </button>
        </div>
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

.avatar-big {
    width: 170px;
    height: 170px;

    border-radius: 50%;
    object-fit: cover;

    border: 2px solid var(--accent);
}
</style>
