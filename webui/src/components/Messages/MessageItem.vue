<script>
import { user } from "@/state/user"
import { expandUrl } from "@/utils/media"
import { Poller } from "@/services/poller"
import { getReceipts } from "@/services/messages"

export default {
    props: {
        message: { type: Object, required: true }
    },

    data() {
        return {
            poller: null,
            checkIcon: "check-sent"
        }
    },

    computed: {
        isMine() {
            return this.message.senderId === user.userId
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

        updateCheckIcon(receipts) {
            if (!receipts) {
                this.poller?.stopPolling()
                return
            }

            const statuses = receipts.map(r => r.status)

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
            const receipts = await getReceipts(this.message.messageId)
            this.updateCheckIcon(receipts)
        })

        this.poller.startPolling()
    },

    emits: ["openImage"],

    beforeUnmount() {
        this.poller?.stopPolling()
    }
}
</script>
<template>
    <div class="message-item" :class="{ mine: isMine }">
        <div class="message-item-bubble">
            <div v-if="message.text" class="message-item-text">
                {{ message.text }}
            </div>

            <img v-if="message.attachmentUrl" :src="expandUrl(message.attachmentUrl)" class="message-item-image"
                @click="$emit('openImage', expandUrl(message.attachmentUrl))" />

            <div class="message-item-meta">
                <span class="message-item-time">
                    {{ getTime(message.createdAt) }}
                </span>
            </div>
            <img v-if="isMine" class="message-item-check-icon" :src="`/icons/${checkIcon}.svg`" alt="" />
        </div>
    </div>
</template>

<style scoped>
.message-item {
    display: flex;
    margin: 6px 20px;
}

.message-item.mine {
    justify-content: flex-end;
}

.message-item-bubble {
    max-width: 70%;
    padding: 10px 14px;

    border-radius: 18px;

    background: var(--surface);
    border: 1px solid var(--border);
}

.message-item.mine .message-item-bubble {
    background: var(--accent);
    border: 1px solid var(--accent-strong);
}

.message-item-meta {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 4px;
    margin-top: 2px;
}

.message-item-meta .message-item-time {
    font-size: 0.72rem;
    opacity: 0.55;
    margin-left: auto;
}

.message-item-image {
    width: 220px;
    height: 160px;
    border-radius: 12px;
    margin-top: 8px;
    object-fit: cover;
    display: block;
}

.message-item-text {
    font-size: 0.95rem;
    line-height: 1.35;
    word-wrap: break-word;
    white-space: pre-wrap;
    padding-right: 40px;
    color: var(--text);
}

.message-item-check-icon {
    width: 16px;
    height: 16px;
    opacity: 0.9;
}
</style>
