<script>
import imageModal from '@/state/imageModal';

export default {
    computed: {
        visible() {
            return imageModal.visible
        },

        photoUrl() {
            return imageModal.photoUrl
        }
    },

    methods: {
        close() {
            imageModal.visible = false
        }
    }
}
</script>

<template>
    <div v-if="visible" class="image-modal" @click="close()">
        <img :src="photoUrl" class="image-modal-content" />
    </div>
</template>

<style scoped>
.image-modal {
    position: fixed;
    inset: 0;
    z-index: 2;

    display: flex;
    align-items: center;
    justify-content: center;

    background: rgba(0, 0, 0, 0.85);
}

.image-modal-content {
    max-width: 90%;
    max-height: 90%;

    border-radius: 12px;
}
</style>
