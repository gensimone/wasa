<script>
import { users } from "@/state/users";
import { user } from "@/state/user";
import { expandUrl } from "@/utils/media";
import { getIcon } from "@/state/theme";

export default {
  name: "MessageInfo",

  props: {
    receipts: {
      type: Array,
      default: () => [],
    },
    reactions: {
      type: Array,
      default: () => [],
    },
  },

  emits: ["deleteReaction", "closeMessageInfo"],

  methods: {
    expandUrl,
    getIcon,

    formatDate(date) {
      if (!date) return "-";
      return new Date(date).toLocaleString();
    },

    messageInfoGetUserPhotoUrl(userId) {
      return users.value.get(userId)?.photoUrl || "";
    },

    messageInfoGetUserName(userId) {
      return users.value.get(userId)?.name || `User #${userId}`;
    },

    messageInfoGetReceiptDate(receipt) {
      if (receipt.status === "read") return receipt.readAt;
      if (receipt.status === "received") return receipt.receivedAt;
      if (receipt.status === "sent") return receipt.sentAt;
      return null;
    },

    isMine(userId) {
      return user.userId === userId;
    },
  },
};
</script>

<template>
  <div class="message-info-overlay">
    <div class="message-info">
      <div class="message-info-header">
        <div class="message-info-title">Message Info</div>

        <button
          class="message-info-close-btn"
          @click="$emit('closeMessageInfo')"
        >
          Close
        </button>
      </div>

      <div class="message-info-section">
        <div class="message-info-section-title">
          Reactions ({{ reactions.length }})
        </div>

        <div class="message-info-scroll">
          <div
            v-for="r in reactions"
            :key="`${r.messageId}-${r.senderId}-${r.emoji}`"
            class="message-info-row"
          >
            <div class="message-info-left">
              <img
                class="message-info-avatar"
                :src="expandUrl(messageInfoGetUserPhotoUrl(r.senderId))"
              />
              <span class="message-info-name">
                {{ messageInfoGetUserName(r.senderId) }}
              </span>
            </div>

            <button
              v-if="isMine(r.senderId)"
              class="message-info-delete-reaction"
              @click="$emit('deleteReaction')"
            >
              <img :src="getIcon('remove')" class="icon-img" />
            </button>

            <div class="message-info-right">
              <img
                class="message-info-emoji"
                :src="`/icons/reactions/${r.emoji}.svg`"
              />
            </div>
          </div>

          <div v-if="!reactions.length" class="message-info-empty">
            No reactions
          </div>
        </div>
      </div>

      <div class="message-info-section">
        <div class="message-info-section-title">
          Receipts ({{ receipts.length }})
        </div>

        <div class="message-info-scroll">
          <div v-for="r in receipts" :key="r.userId" class="message-info-row">
            <div class="message-info-left">
              <img
                class="message-info-avatar"
                :src="expandUrl(messageInfoGetUserPhotoUrl(r.userId))"
              />
              <span class="message-info-name">
                {{ messageInfoGetUserName(r.userId) }}
              </span>
            </div>

            <div class="message-info-right">
              <span class="message-info-status">
                {{ r.status }}
              </span>

              <span class="message-info-timestamp">
                {{ formatDate(messageInfoGetReceiptDate(r)) }}
              </span>
            </div>
          </div>

          <div v-if="!receipts.length" class="message-info--empty">
            No receipts
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.message-info-overlay {
  position: fixed;
  inset: 0;
  z-index: 2;

  display: flex;
  justify-content: center;
  align-items: center;

  background: rgba(0, 0, 0, 0.85);
}

.message-info {
  width: min(700px, 90vw);
  max-height: 80vh;

  display: flex;
  flex-direction: column;

  background: var(--surface-3);
  border-radius: 14px;
  border: 1px solid var(--text);
  padding: 12px;
}

.message-info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-info-title {
  font-weight: 600;
}

.message-info-close-btn {
  padding: 6px 10px;
  border: 1px solid var(--border);
  color: var(--text);
  background: var(--surface);
  border-radius: 6px;
  cursor: pointer;
}

.message-info-section {
  flex: 1;
  min-height: 0;
  border: 1px solid var(--border);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
}

.message-info-section-title {
  padding: 12px;
  font-weight: 600;
  border-bottom: 1px solid var(--border);
}

.message-info-scroll {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.message-info-row {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border);
}

.message-info-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.message-info-delete-reaction {
  display: flex;
  align-items: center;
  justify-self: end;
  border-radius: 20px;
  border: 1px solid var(--border);
  margin-right: 7px;
}

.message-info-right {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-self: end;
}

.message-info-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.message-info-name {
  font-weight: 500;
}

.message-info-emoji {
  width: 40px;
  height: 40px;
}

.message-info-status {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
}

.message-info-timestamp {
  font-size: 12px;
  color: #888;
}

.message-info-empty {
  padding: 16px;
  color: #888;
}
</style>
