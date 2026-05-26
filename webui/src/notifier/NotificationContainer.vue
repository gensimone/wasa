<script>
import Notification from "./Notification.vue"
import { setNotifier } from ".";

export default {
    components: { Notification },

    data() {
        return {
            notifications: []
        }
    },

    mounted() {
        setNotifier(this.add)
    },

    methods: {
        add(notification) {
            const text = notification.text?.trim()

            const notificationId = Date.now() + Math.random()

            this.notifications.push({
                notificationId: notificationId,
                type: notification.type || "info",
                duration: notification.duration || 5000,
                text: text,

                // Conversation messages.
                thumbnailUrl: notification.thumbnailUrl || null,
                attachmentUrl: notification.attachmentUrl || null,
                id: notification.id || null,
                isGroup: notification.isGroup || null
            })
        },

        remove(notificationId) {
            this.notifications = this.notifications
                .filter(n => n.notificationId !== notificationId)
        }
    }
}
</script>

<template>
    <TransitionGroup name="toast" tag="div" class="toast-container">
        <Notification v-for="n in notifications" :key="n.notificationId" :text="n.text" :id="n.id" :isGroup="n.isGroup"
            :attachmentUrl="n.attachmentUrl" :thumbnailUrl="n.thumbnailUrl" :type="n.type" :duration="n.duration"
            @close="remove(n.notificationId)" />
    </TransitionGroup>
</template>

<style scoped>
.toast-container {
    z-index: 3;
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
