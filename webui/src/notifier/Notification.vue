<script>
export default {
    props: {
        type: {
            type: String,
            required: true,
            validator: (v) =>
                ["success", "error", "warning", "info"].includes(v)
        },
        message: { type: String, required: true },
        duration: { type: Number, default: 5000 }
    },

    data() {
        return {
            progress: 100,
            startTime: 0,
            remaining: this.duration,
            timeoutId: null,
            frameId: null
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
    <div class="notification" :class="type" @mouseenter="pauseTimer" @mouseleave="resumeTimer" @click="close">
        <div class="notification-message">
            {{ message }}
        </div>

        <div class="progress" :style="{ width: progress + '%' }" />
    </div>
</template>

<style scoped>
.notification {
    pointer-events: auto;
    position: relative;
    width: 320px;
    padding: 16px;
    border-radius: 10px;
    cursor: pointer;
    overflow: hidden;
    color: white;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.25);
    z-index: 999999;
}

.notification.success {
    background: #16251a;
    border-left: 5px solid #4ade80;
}

.notification.success .progress {
    background: #4ade80;
}

.notification.error {
    background: #2a1a1a;
    border-left: 5px solid #ef4444;
}

.notification.error .progress {
    background: #ef4444;
}

.notification.warning {
    background: #2a241a;
    border-left: 5px solid #f59e0b;
}

.notification.warning .progress {
    background: #f59e0b;
}

.notification.info {
    background: #1b1f2a;
    border-left: 5px solid #60a5fa;
}

.notification.info .progress {
    background: #60a5fa;
}

.progress {
    position: absolute;
    bottom: 0;
    left: 0;
    height: 4px;
    transition: width 0.1s linear;
}
</style>
