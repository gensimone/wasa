<script>
import { expandURL } from '@/utils/media';
export default {
    methods: {
        expandURL
    },

    props: {
        avatarUrl: String,
        isDefault: Boolean,
        photoChanged: Boolean,
        loading: Boolean
    },

    emits: ["upload", "remove", "save"]
}
</script>

<template>
    <div class="avatar-row">
        <button type="button" class="icon-btn" :disabled="isDefault" @click="$emit('remove')">
            <img src="/icons/remove.svg" class="icon-img">
        </button>

        <label class="avatar-wrapper">
            <img :src="expandURL(avatarUrl)" class="avatar-big" />
            <input type="file" accept="image/*" hidden @change="$emit('upload', $event)" />
        </label>

        <button type="button" class="icon-btn" :disabled="!photoChanged" @click="$emit('save')">
            <img src="/icons/upload.svg" class="icon-img">
        </button>

    </div>
</template>

<style scoped>
.avatar-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}

.avatar-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.avatar-big {
    width: 170px;
    height: 170px;
    border-radius: 50%;
    object-fit: cover;

    border: 2px solid rgba(0, 255, 120, 0.18);
    box-shadow: 0 0 50px rgba(0, 255, 120, 0.08);

    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.avatar-big:hover {
    transform: scale(1.04);
    box-shadow: 0 0 70px rgba(0, 255, 120, 0.15);
}
</style>
