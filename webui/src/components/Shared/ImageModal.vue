<script>
import { expandUrl } from "@/utils/media"
export default {
    props: {
        visible: { type: Boolean, required: true },
        imageUrl: String
    },

    emits: ["close"],

    methods: {
        expandUrl
    }
}
</script>

<template>
    <div v-if="visible" class="image-modal" @click="$emit('close')">
        <img :src="imageUrl" class="image-modal-content" />
    </div>
</template>

<style scoped>
.image-modal {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.85);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    cursor: zoom-out;
}

.image-modal-content {
    max-width: 90%;
    max-height: 90%;
    border-radius: 12px;
    box-shadow: 0 20px 80px rgba(0, 0, 0, 0.6);
    transform: scale(1);
    animation: zoomIn 0.2s ease;
}

@keyframes zoomIn {
    from {
        transform: scale(0.85);
        opacity: 0;
    }

    to {
        transform: scale(1);
        opacity: 1;
    }
}
</style>
