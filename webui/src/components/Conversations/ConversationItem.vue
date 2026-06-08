<script>
import Poller from "@/services/poller";
import { expandUrl } from "@/utils/media";
import { user } from "@/state/user";
import { getTime, getCheckIcon } from "@/utils/messages";
import { getIcon } from "@/state/theme";
import { getReceipts } from "@/services/messages";

export default {
  props: {
    conversation: { type: Object, required: true }
  },

  emits: ["select"],

  data() {
    return {
      checkIcon: "check-sent",
      poller: null,
    };
  },

  computed: {
    isMine() {
      return this.conversation.lastMessage?.senderId === user.userId;
    },
  },

  async mounted() {
    if (!this.isMine) return;

    this.poller = new Poller();

    this.poller.callback = async () => {
      const receipts = await getReceipts(
        this.conversation.lastMessage.messageId,
      );
      const icon = getCheckIcon(receipts);
      if (!icon || icon === "check-read") this.poller.stopPolling();

      this.checkIcon = icon;
    };

    this.poller.startPolling();
  },

  beforeUnmount() {
    this.poller?.stopPolling();
  },

  methods: {
    expandUrl,
    getTime,
    getIcon,
  },
};
</script>
<template>
  <div class="conversation-item" @click="$emit('select', conversation)">
    <div class="conversation-item-photo-wrapper">
      <img
        :src="expandUrl(conversation.photoUrl)"
        class="conversation-item-photo"
      >
    </div>

    <div class="conversation-item-info">
      <div class="conversation-item-header">
        <div class="conversation-item-name">
          {{ conversation.name }}
        </div>

        <div v-if="conversation.lastMessage" class="conversation-item-right">
          <div
            v-if="conversation.lastMessage.attachmentUrl"
            class="conversation-item-attachment-inline"
          >
            <img
              :src="expandUrl(conversation.lastMessage?.attachmentUrl)"
              class="conversation-item-photo"
            >
          </div>

          <div class="conversation-item-time">
            {{ getTime(conversation.lastMessage.createdAt) }}
          </div>
        </div>
      </div>

      <div class="conversation-item-text-box">
        <img
          v-if="isMine"
          class="message-item-check-icon"
          :src="getIcon(checkIcon)"
          alt=""
        >

        <div
          v-if="conversation.lastMessage"
          class="conversation-item-last-message"
        >
          {{ conversation.lastMessage.text }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conversation-item {
  display: flex;
  align-items: center;
  gap: 14px;

  padding: 14px;
  margin-bottom: 10px;

  border-radius: 18px;

  background: var(--surface-2);
  border: 1px solid var(--border);

  cursor: pointer;
  overflow: hidden;

  animation: fadeInUp 0.35s ease both;
}

.conversation-item-text-box {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
}

.conversation-item:hover {
  transform: translateY(-6px) scale(1.02);
  border: 1px solid var(--accent);
}

.conversation-item-photo-wrapper {
  width: 75px;
  height: 75px;
  border-radius: 16px;
  overflow: hidden;
  flex-shrink: 0;
}

.conversation-item-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conversation-item-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.conversation-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.conversation-item-name {
  font-size: 1.05rem;
  font-weight: 800;

  color: var(--text);

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-item-right {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.conversation-item-attachment-inline {
  width: 50px;
  height: 50px;

  border-radius: 12px;
  overflow: hidden;

  flex-shrink: 0;
}

.conversation-item-attachment-inline img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conversation-item-time {
  font-size: 0.72rem;
  color: var(--text-muted);

  white-space: nowrap;
}

.conversation-item-last-message {
  margin-top: 4px;

  font-size: 0.92rem;
  color: var(--text-muted);

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.message-item-check-icon {
  width: 16px;
  height: 16px;
  opacity: 0.9;
}
</style>
