<script>
import { user } from "@/state/user"
import { expandUrl } from "@/utils/media"
import { Poller } from "@/services/poller"
import { getReceipts } from "@/services/messages"

export default {
    props: {
        message: Object,
        isMine: Boolean
    },

    emits: ["openImage"],

    data() {
        return {
            poller: null,
            receipts: [],
            checkIcon: "check-sent"
        }
    },

    methods: {
        expandUrl,

        getTime(sentAt) {
            const date = new Date(sentAt)

            const hh = String(date.getHours()).padStart(2, "0")
            const mm = String(date.getMinutes()).padStart(2, "0")

            return `${hh}:${mm}`
        },

        updateCheckIcon() {
            const statuses = this.receipts.map(r => r.status)

            const hasSent = statuses.includes("sent")
            const hasReceived = statuses.includes("received")
            const allRead =
                statuses.length > 0 &&
                statuses.every(s => s === "read")

            if (hasSent) {
                this.checkIcon = "check-sent"
            }
            else if (hasReceived) {
                this.checkIcon = "check-received"
            }
            else if (allRead) {
                this.checkIcon = "check-read"
                this.poller?.stopPolling()
            }
        }
    },

    mounted() {
        if (this.message.senderId !== user.userId) return

        this.poller = new Poller(async () => {
            this.receipts = await getReceipts(
                this.message.messageId
            )

            this.updateCheckIcon()
        })

        this.poller.startPolling()
    },

    beforeUnmount() {
        this.poller?.stopPolling()
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
                @click="$emit('openImage', expandUrl(message.attachmentUrl))" />

            <div class="message-meta">
                <span class="time">
                    {{ getTime(message.createdAt) }}
                </span>
            </div>
            <img v-if="isMine" class="check-icon" :src="`/icons/${checkIcon}.svg`" alt="" />
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
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
    margin-top: 2px;
}

.check-icon {
    width: 16px;
    height: 16px;
    opacity: 0.9;
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
