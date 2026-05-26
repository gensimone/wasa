<script>
import { defaultUserPhotoUrl } from "@/assets/default"
import { expandUrl } from "@/utils/media"
export default {
    props: {
        user: { type: Object, required: true }
    },

    emits: ["select"],

    methods: {
        expandUrl,

        getPhotoUrl(user) {
            return user.photoUrl || defaultUserPhotoUrl
        }
    }
}
</script>

<template>
    <div class="user-item" @click="$emit('select', user)">
        <div class="user-item-photo-wrapper">
            <img :src="expandUrl(getPhotoUrl(user))" class="user-item-photo" />
        </div>
        <div class="user-item-info">
            <div class="user-item-name">
                {{ user.name }}
            </div>
        </div>
    </div>
</template>

<style scoped>
.user-item {
    display: flex;
    align-items: center;
    gap: 14px;

    padding: 14px;
    margin-bottom: 10px;

    border-radius: 18px;

    background: var(--surface-2);
    border: 1px solid var(--border);

    cursor: pointer;
    overflow: hidden;

    animation: fadeInUp 0.35s ease both;
}

.user-item:hover {
    transform: translateY(-6px) scale(1.02);
    border: 1px solid var(--accent);
}

.user-item-photo-wrapper {
    width: 75px;
    height: 75px;
    border-radius: 16px;
    overflow: hidden;
    flex-shrink: 0;
}

.user-item-photo {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.user-item-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 0;
}

.user-item-name {
    font-size: 1.05rem;
    font-weight: 800;

    color: var(--text);

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
