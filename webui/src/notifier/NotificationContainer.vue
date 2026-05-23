<script>
import Notification from "./Notification.vue"

export default {
    components: { Notification },

    data() {
        return {
            notifications: []
        }
    },

    methods: {
        add(notification) {
            const newMessage = notification.message?.trim()
            if (!newMessage) throw new Error("Empty notification message")

            const lastNotification = this.notifications.at(-1)

            if (
                lastNotification?.message === newMessage &&
                lastNotification?.type === notification.type
            ) return

            const id = Date.now() + Math.random()

            this.notifications.push({
                id,
                message: newMessage,
                type: notification.type || "info",
                duration: notification.duration || 5000
            })
        },

        remove(id) {
            this.notifications = this.notifications.filter(n => n.id !== id)
        }
    }
}
</script>

<template>
    <TransitionGroup name="toast" tag="div" class="toast-container">
        <Notification v-for="n in notifications" :key="n.id" :message="n.message" :type="n.type" :duration="n.duration"
            @close="remove(n.id)" />
    </TransitionGroup>
</template>

<style scoped>
.toast-container {
    z-index: 999999;
    position: fixed;
    bottom: 20px;
    right: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.toast-enter-active,
.toast-leave-active {
    transition: all 0.50s cubic-bezier(0.2, 0.9, 0.2, 1);
}

.toast-move {
    transition: transform 0.25s ease;
}

.toast-enter-from {
    opacity: 0;
    transform: translateX(40px);
}

.toast-enter-to {
    opacity: 1;
    transform: translateX(0);
}

.toast-leave-from {
    opacity: 1;
    transform: translateX(0);
}

.toast-leave-to {
    opacity: 0;
    transform: translateX(40px);
}
</style>
