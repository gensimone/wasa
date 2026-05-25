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
    z-index: 3;

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
