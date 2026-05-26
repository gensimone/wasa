<script>
import { expandUrl } from "@/utils/media"
import { getIcon } from "@/state/theme";
import imageModal from "@/state/imageModal"

export default {
    props: {
        type: {
            type: String,
            required: true,
            validator: (v) =>
                ["success", "error", "warning", "info", "message"].includes(v)
        },
        duration: { type: Number, default: 5000 },

        text: { type: String, required: true },
        thumbnailUrl: { type: String, required: false },
        attachmentUrl: { type: String, required: false },
        id: { type: Number, required: false },
        isGroup: { type: Boolean, required: false }
    },

    data() {
        return {
            progress: 100,
            startTime: 0,
            remaining: this.duration,
            timeoutId: null,
            frameId: null,
        }
    },

    computed: {
        photoUrl() {
            return this.type === "message"
                ? expandUrl(this.thumbnailUrl)
                : getIcon(this.type)
        }
    },

    mounted() {
        this.startTimer()
    },

    beforeUnmount() {
        this.clearTimers()
    },

    emits: ["close"],

    methods: {
        expandUrl,
        getIcon,

        startTimer() {
            this.clearTimers()
            this.startTime = Date.now()

            this.timeoutId = setTimeout(() => {
                this.close()
            }, this.remaining)

            this.updateProgress()
        },

        updateProgress() {
            const elapsed = Date.now() - this.startTime
            const timeLeft = this.remaining - elapsed

            this.progress = Math.max((timeLeft / this.duration) * 100, 0)

            if (timeLeft <= 0) {
                this.close()
                return
            }

            this.frameId = requestAnimationFrame(this.updateProgress)
        },

        pauseTimer() {
            this.clearTimers()
            const elapsed = Date.now() - this.startTime
            this.remaining -= elapsed
        },

        resumeTimer() {
            if (this.remaining > 0) {
                this.startTimer()
            }
        },

        handleClick() {
            if (this.type === "message") {
                this.$router.push({
                    name: "conversation",
                    params: { id: this.id },
                    query: { direct: !this.isGroup }
                })
                this.close()
                imageModal.visible = false
            } else {
                this.close()
            }
        },

        close() {
            this.clearTimers()
            this.$emit("close")
        },

        clearTimers() {
            clearTimeout(this.timeoutId)
            cancelAnimationFrame(this.frameId)
        }
    }
}
</script>

<template>
    <div class="notification" :class="type" @mouseenter="pauseTimer" @mouseleave="resumeTimer" @click="handleClick()">

        <img :src="photoUrl" class="notification-thumbnail" />

        <div class="notification-content">

            <div class="notification-message">
                {{ text }}
            </div>

            <img v-if="attachmentUrl" :src="expandUrl(attachmentUrl)" class="notification-attachment" />

        </div>

        <div class="progress" :style="{ width: progress + '%' }" />
    </div>
</template>

<style scoped>
.notification {
    position: relative;
    display: flex;
    align-items: flex-start;
    gap: 12px;

    width: 320px;
    padding: 12px 14px;
    border-radius: 12px;

    background: var(--surface);
    border: 1px solid var(--border);

    color: var(--text);

    backdrop-filter: blur(18px);
    -webkit-backdrop-filter: blur(18px);

    box-shadow:
        0 10px 30px var(--shadow),
        0 0 0 1px rgba(255, 255, 255, 0.02);

    overflow: hidden;
    z-index: 4;
    pointer-events: auto;
}

/* glow laterale */
.notification::before {
    content: "";
    position: absolute;
    inset: 0 auto 0 0;
    width: 4px;
    border-radius: inherit;
    background: var(--accent-color);
}

/* =========================
   CONTENT
========================= */

.notification-thumbnail {
    width: 46px;
    height: 46px;
    border-radius: 10px;
    object-fit: cover;
    flex-shrink: 0;
}

.notification-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-width: 0;
}

.notification-message {
    font-size: 14px;
    line-height: 1.35;
    word-break: break-word;
    color: var(--text);
}

.notification-attachment {
    width: 100%;
    max-width: 160px;
    max-height: 90px;
    border-radius: 8px;
    object-fit: cover;
    align-self: flex-start;
}

/* =========================
   PROGRESS
========================= */

.progress {
    position: absolute;
    bottom: 0;
    left: 0;
    height: 3px;
    width: 100%;

    background: var(--accent-color);

    transition: width 0.1s linear;
}

/* =========================
   TYPES
========================= */

.notification.success {
    --accent-color: var(--success);
}

.notification.error {
    --accent-color: var(--error);
}

.notification.warning {
    --accent-color: var(--warning);
}

.notification.info,
.notification.message {
    --accent-color: var(--info);
}
</style>
