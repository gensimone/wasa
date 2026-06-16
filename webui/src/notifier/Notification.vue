<script>
import { expandUrl } from "@/utils/media";
import { getIcon } from "@/state/theme";
import imageModal from "@/state/imageModal";

export default {
  props: {
    type: {
      type: String,
      required: true,
      validator: (v) =>
        ["success", "error", "warning", "info", "message"].includes(v),
    },
    duration: { type: Number, default: 5000 },

    text: { type: String, required: true },
    thumbnailUrl: { type: String, required: false },
    attachmentUrl: { type: String, required: false },
    id: { type: Number, required: false },
    isDirect: { type: Boolean, required: false },
  },

  emits: ["close"],

  data() {
    return {
      startTime: 0,
      remaining: this.duration,
      timeoutId: null,
      frameId: null,
    };
  },

  computed: {
    photoUrl() {
      return this.type === "message"
        ? expandUrl(this.thumbnailUrl)
        : getIcon(this.type);
    },
  },

  mounted() {
    this.startTimer();
  },

  beforeUnmount() {
    clearTimeout(this.timeoutId);
  },

  methods: {
    expandUrl,
    getIcon,

    startTimer() {
      clearTimeout(this.timeoutId);
      this.startTime = Date.now();

      this.timeoutId = setTimeout(() => {
        this.close();
      }, this.remaining);
    },

    pauseTimer() {
      clearTimeout(this.timeoutId);
      const elapsed = Date.now() - this.startTime;
      this.remaining -= elapsed;
    },

    resumeTimer() {
      if (this.remaining > 0) {
        this.startTimer();
      }
    },

    handleClick() {
      if (this.type === "message") {
        this.$router.push({
          name: "conversation",
          params: { id: this.id },
          query: { direct: this.isDirect },
        });
        this.close();
        imageModal.visible = false;
      } else {
        this.close();
      }
    },

    close() {
      clearTimeout(this.timeoutId);
      this.$emit("close");
    },

  },
};
</script>

<template>
  <div
    class="notification"
    :class="type"
    @mouseenter="pauseTimer"
    @mouseleave="resumeTimer"
    @click="handleClick()"
  >
    <img :src="photoUrl" class="notification-thumbnail">

    <div class="notification-content">
      <div class="notification-message">
        {{ text }}
      </div>

      <img
        v-if="attachmentUrl"
        :src="expandUrl(attachmentUrl)"
        class="notification-attachment"
      >
    </div>
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

.notification::before {
  content: "";
  position: absolute;
  inset: 0 auto 0 0;
  width: 4px;
  border-radius: inherit;
  background: var(--accent-color);
}

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
