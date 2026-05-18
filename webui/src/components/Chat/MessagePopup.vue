<script>
export default {
    props: {
        visible: Boolean,
        x: Number,
        y: Number,
    },
    emits: ["react"]
}
</script>

<template>
    <div v-if="visible" class="message-popup" :style="{ top: y + 'px', left: x + 'px' }">
        <span @click="$emit('react', 'like')"><img src="/icons/reactions/like.svg" class="icon-img"></span>
        <span @click="$emit('react', 'love')"><img src="/icons/reactions/love.svg" class="icon-img"></span>
        <span @click="$emit('react', 'laugh')"><img src="/icons/reactions/laugh.svg" class="icon-img"></span>
        <span @click="$emit('react', 'sad')"><img src="/icons/reactions/sad.svg" class="icon-img"></span>
        <span @click="$emit('react', 'angry')"><img src="/icons/reactions/angry.svg" class="icon-img"></span>
    </div>
</template>

<style scoped>
.message-popup {
    position: absolute;
    display: flex;
    gap: 8px;
    padding: 8px 10px;
    background: rgba(20, 20, 20, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 14px;
    z-index: 99999;

    transform: translateX(-50%);
}
</style>
