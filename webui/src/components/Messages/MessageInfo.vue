<script>
import { users } from "@/state/users";
import { expandUrl } from "@/utils/media";

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

  methods: {
    expandUrl,

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
  },
};
</script>

<template>
  <div class="message-info">
    <div class="message-info__header">
      <div class="message-info__title">Message Info</div>

      <button
        class="message-info__close-btn"
        @click="$emit('closeMessageInfo')"
      >
        Close
      </button>
    </div>

    <div class="message-info__section">
      <div class="message-info__section-title">
        Reactions ({{ reactions.length }})
      </div>

      <div class="message-info__scroll">
        <div
          v-for="r in reactions"
          :key="`${r.messageId}-${r.senderId}-${r.emoji}`"
          class="message-info__row"
        >
          <div class="message-info__left">
            <img
              class="message-info__avatar"
              :src="expandUrl(messageInfoGetUserPhotoUrl(r.senderId))"
            />
            <span class="message-info__name">
              {{ messageInfoGetUserName(r.senderId) }}
            </span>
          </div>

          <div class="message-info__right">
            <img
              class="message-info__emoji"
              :src="`/icons/reactions/${r.emoji}.svg`"
            />
          </div>
        </div>

        <div v-if="!reactions.length" class="message-info__empty">
          No reactions
        </div>
      </div>
    </div>

    <div class="message-info__section">
      <div class="message-info__section-title">
        Receipts ({{ receipts.length }})
      </div>

      <div class="message-info__scroll">
        <div v-for="r in receipts" :key="r.userId" class="message-info__row">
          <div class="message-info__left">
            <img
              class="message-info__avatar"
              :src="expandUrl(messageInfoGetUserPhotoUrl(r.userId))"
            />
            <span class="message-info__name">
              {{ messageInfoGetUserName(r.userId) }}
            </span>
          </div>

          <div class="message-info__right">
            <span class="message-info__status">
              {{ r.status }}
            </span>

            <span class="message-info__timestamp">
              {{ formatDate(messageInfoGetReceiptDate(r)) }}
            </span>
          </div>
        </div>

        <div v-if="!receipts.length" class="message-info__empty">
          No receipts
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.message-info {
  height: 100vh;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 12px;
}

.message-info__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-info__title {
  font-weight: 600;
}

.message-info__close-btn {
  padding: 6px 10px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 6px;
  cursor: pointer;
}

.message-info__section {
  flex: 1;
  min-height: 0;
  border: 1px solid #ddd;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
}

.message-info__section-title {
  padding: 12px;
  font-weight: 600;
  border-bottom: 1px solid #ddd;
}

.message-info__scroll {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.message-info__row {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid #f3f3f3;
}

.message-info__left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.message-info__right {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-self: end;
}

.message-info__avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
}

.message-info__name {
  font-weight: 500;
}

.message-info__emoji {
  width: 40px;
  height: 40px;
}

.message-info__status {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
}

.message-info__timestamp {
  font-size: 12px;
  color: #888;
}

.message-info__empty {
  padding: 16px;
  color: #888;
}
</style>
