<script>
import { user } from "@/state/user.js"
import Message from "@/components/Conversation/Message.vue"

export default {
    data() {
        return {
            user,
        }
    },

    components: {
        Message
    },

    props: {
        messages: { type: Array, required: true }
    },

    emits: ["openImage"]
}
</script>

<template>
    <div class="messages" ref="container">
        <Message v-for="message in messages" :key="message.messageId" :message="message"
            :isMine="message.senderId === user.userId" @openImage="$emit('openImage', $event)" />
    </div>
</template>

<style scoped>
.messages {
    width: min(720px, 100%);
    overflow-x: hidden;
    height: 100vh;
    display: flex;
    flex-direction: column;
    padding: 0px 18px;
    position: relative;
}

.messages::before {
    content: "";
    position: absolute;
    inset: -1px;
    border-radius: 22px;
    background: linear-gradient(120deg,
            rgba(0, 255, 120, 0.14),
            rgba(255, 255, 255, 0.05),
            rgba(0, 255, 120, 0.14));

    -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0.25;
    pointer-events: none;
    animation: borderGlow 4s ease-in-out infinite alternate;
}

@keyframes borderGlow {
    0% {
        opacity: 0.18;
    }

    100% {
        opacity: 0.35;
    }
}

@media (max-width: 600px) {
    .messages {
        height: 65vh;
        padding: 14px;
    }
}
</style>
