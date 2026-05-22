<script>
import { expandUrl } from "@/utils/media"
export default {
    props: {
        message: Object,
        isMine: Boolean
    },

    emits: ["openImage"],

    methods: {
        expandUrl,

        getTime(sentAt) {
            const date = new Date(sentAt)
            const hh = date.getHours()
            const mm = date.getMinutes()
            return `${hh}:${mm}`
        }
    }
}
</script>

<template>
    <div class="message" :class="{ mine: isMine }">
        <div class="bubble">
            <div v-if="message.text" class="message-text">
                {{ message.text }}
            </div>

            <img v-if="message.attachmentUrl" :src="expandUrl(message.attachmentUrl)" class="message-image"
                @click="$emit('openImage', message.attachmentUrl)" />

            <div class="message-meta">
                <span class="time">
                    {{ getTime(message.createdAt) }}
                </span>
            </div>
        </div>
    </div>
</template>

<style>
.message {
    display: flex;
    margin: 6px 0;
    animation: fadeInUp 0.35s ease both;
}

.message.mine {
    justify-content: flex-end;
}

.message-image {
    width: 220px;
    height: 160px;
    border-radius: 12px;
    margin-top: 8px;
    object-fit: cover;
    display: block;
}

.message-text {
    font-size: 0.95rem;
    line-height: 1.35;
    word-wrap: break-word;
    white-space: pre-wrap;
    padding-right: 40px;
}

.message-content {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.message-meta {
    display: flex;
    justify-content: flex-end;
    margin-top: 2px;
}

.time {
    font-size: 0.72rem;
    opacity: 0.55;
    margin-left: auto;
}

.bubble {
    max-width: 70%;
    padding: 10px 14px;
    border-radius: 18px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    color: rgba(245, 245, 245, 0.92);
    backdrop-filter: blur(10px);
    position: relative;
    overflow: hidden;
}

.message.mine .bubble {
    background: rgba(0, 255, 120, 0.12);
    border: 1px solid rgba(0, 255, 120, 0.2);
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(12px) scale(0.98);
    }

    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}
</style>
