<script>
export default {
    props: {
        message: String,
        error: Boolean
    }
}
</script>

<template>
    <div v-if="message" class="status-area">
        <p :class="error ? 'error' : 'success'">
            {{ message }}
        </p>
    </div>
</template>

<style scoped>
.status-area {
    width: 100%;
    margin-top: 14px;

    display: flex;
    flex-direction: column;
    text-align: center;
}

.error {
    color: rgba(255, 80, 80, 0.85);
    font-size: 0.9rem;

    align-items: center;
    animation: fadeInError 0.25s ease both;
    backdrop-filter: blur(10px);
    background: rgba(255, 80, 80, 0.08);
    border-radius: 12px;
    border: 1px solid rgba(255, 80, 80, 0.2);
    display: flex;
    font-size: 0.85rem;
    font-weight: 500;
    gap: 12px;
    justify-content: space-between;
    margin: 8px 12px;
    padding: 10px 12px;
}

.success {
    margin-top: 10px;
    color: rgba(0, 255, 120, 0.7);
    font-size: 0.9rem;

    align-items: center;
    animation: fadeInError 0.25s ease both;
    backdrop-filter: blur(10px);
    background: rgba(0, 255, 120, 0.08);
    border-radius: 12px;
    border: 1px solid rgba(0, 255, 120, 0.2);
    display: flex;
    font-size: 0.85rem;
    font-weight: 500;
    gap: 12px;
    justify-content: space-between;
    margin: 8px 12px;
    padding: 10px 12px;
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
