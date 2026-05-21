<script>
export default {
    props: {
        error: {
            type: String,
            default: null
        }
    }
}
</script>

<template>
    <div v-if="error" class="send-error">
        <span class="message">
            {{ error }}
        </span>
    </div>
</template>

<style scoped>
.send-error {
    margin: 8px 12px;
    padding: 10px 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    font-size: 0.85rem;
    font-weight: 500;
    color: rgba(255, 90, 90, 0.95);
    background: rgba(255, 80, 80, 0.08);
    border: 1px solid rgba(255, 80, 80, 0.2);
    border-radius: 12px;
    backdrop-filter: blur(10px);
    animation: fadeInError 0.25s ease both;
}

@keyframes fadeInError {
    from {
        opacity: 0;
        transform: translateY(6px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
