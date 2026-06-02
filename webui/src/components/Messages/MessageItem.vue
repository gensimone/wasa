<script>
import logger from "@/utils/logger";
import Poller from "@/services/poller";
import ContextMenu from "@/components/Messages/ContextMenu/ContextMenu.vue";
import MessageReactions from "@/components/Messages/MessageReactions.vue";
import MessageCheckIcon from "@/components/Messages/MessageCheckIcon.vue";

import { getTime } from "@/utils/messages";
import { getIcon } from "@/state/theme";
import { user } from "@/state/user";
import { expandUrl } from "@/utils/media";
import { addReaction, getReactions } from "@/services/reactions";
import { setMessageStatusAsRead } from "@/services/messages";
import { setImageModal } from "@/state/imageModal";
import { handleError } from "@/utils/errors";
import { getCheckIcon } from "@/utils/messages";
import { getReceipts } from "@/services/messages";

export default {
  components: { ContextMenu, MessageReactions, MessageCheckIcon },

  props: {
    message: { type: Object, required: true },
  },

  emits: [
    "replyToMessage",
    "forwardMessage",
    "deleteMessage",
    "showInfoMessage",
  ],

  data() {
    return {
      checkIcon: null,

      reactionsPoller: null,
      receiptsPoller: null,

      reactions: [],
      receipts: [],

      menu: {
        visible: false,
        canClose: false,
        x: 0,
        y: 0,
      },
    };
  },

  computed: {
    isMine() {
      return this.message.senderId === user.userId;
    },
  },

  methods: {
    expandUrl,
    getIcon,
    setImageModal,
    getTime,

    async reactToMessage(reactionData) {
      const message = reactionData.message;
      const emoji = reactionData.emoji;

      try {
        await addReaction(message.messageId, emoji);
      } catch (e) {
        handleError(e);
      }
    },

    onRightClick(e) {
      e.preventDefault();

      this.menu = {
        visible: true,
        canClose: false,
        x: e.clientX,
        y: e.clientY,
      };

      this.$nextTick(() => {
        const menuEl = document.querySelector(".context-menu");

        if (!menuEl) return;

        const { offsetWidth, offsetHeight } = menuEl;

        let x = e.clientX;
        let y = e.clientY;

        const margin = 8;

        if (y + offsetHeight > window.innerHeight - margin) {
          y = e.clientY - offsetHeight;
        }

        if (x + offsetWidth > window.innerWidth - margin) {
          x = e.clientX - offsetWidth;
        }

        x = Math.max(margin, x);
        y = Math.max(margin, y);

        this.menu.x = x;
        this.menu.y = y;

        setTimeout(() => {
          this.menu.canClose = true;
        }, 0);
      });
    },

    closeMenu() {
      this.menu.visible = false;
    },

    async fetchReactions() {
      try {
        this.reactions = (await getReactions(this.message.messageId)) || [];
      } catch (e) {
        if (e.response?.status === 404) {
          // The message has been deleted
          this.reactionsPoller?.stopPolling();
          this.reactions = [];
        } else {
          handleError(e);
        }
      }
    },

    async fetchReceipts() {
      this.receipts = await getReceipts(this.message.messageId);
      const icon = getCheckIcon(this.receipts);
      if (!icon || icon === "check-read") this.receiptsPoller.stopPolling();

      this.checkIcon = icon;
    },

    showInfoMessage() {
      this.$emit("showInfoMessage", {
        receipts: this.receipts,
        reactions: this.reactions,
      });
    },
  },

  async mounted() {
    this.reactionsPoller = new Poller();
    this.reactionsPoller.callback = this.fetchReactions;
    this.reactionsPoller.startPolling();

    if (this.message.senderId !== user.userId) {
      try {
        await setMessageStatusAsRead(this.message.messageId);
      } catch (e) {
        logger.error(e);
      } finally {
        return;
      }
    }

    this.receiptsPoller = new Poller();
    this.receiptsPoller.callback = this.fetchReceipts;
    this.receiptsPoller.startPolling();
  },

  beforeUnmount() {
    this.reactionsPoller?.stopPolling();
    this.receiptsPoller?.stopPolling();
  },
};
</script>
<template>
  <div class="message-item" :class="{ mine: isMine }">
    <div class="message-item-bubble" @contextmenu="onRightClick">
      <img
        v-if="message.isForwarded"
        class="message-item-forward-icon"
        :src="getIcon('forward')"
        alt="forwarded"
      />

      <div v-if="message.text" class="message-item-text">
        {{ message.text }}
      </div>

      <img
        v-if="message.attachmentUrl"
        :src="expandUrl(message.attachmentUrl)"
        class="message-item-image"
        @click="setImageModal(expandUrl(message.attachmentUrl))"
      />

      <div class="message-item-meta">
        <span class="message-item-time">
          {{ getTime(message.createdAt) }}
        </span>
      </div>

      <MessageCheckIcon v-if="isMine" :checkIcon="checkIcon" />

      <MessageReactions :reactions="reactions" />

      <ContextMenu
        v-if="menu.visible"
        :message="message"
        :canClose="menu.canClose"
        :x="menu.x"
        :y="menu.y"
        @close="closeMenu"
        @reactToMessage="reactToMessage"
        @showInfoMessage="showInfoMessage"
        @replyToMessage="$emit('replyToMessage', $event)"
        @forwardMessage="$emit('forwardMessage', $event)"
        @deleteMessage="$emit('deleteMessage', $event)"
      />
    </div>
  </div>
</template>

<style scoped>
.message-item {
  display: flex;
  margin: 15px 20px;
}

.message-item.mine {
  justify-content: flex-end;
}

.message-item-bubble {
  position: relative;

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

.message-item-forward-icon {
  position: absolute;
  top: 6px;
  left: 6px;

  width: 20px;
  height: 20px;

  opacity: 0.55;
  pointer-events: none;
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
  padding-left: 16px;
  font-size: 0.95rem;
  line-height: 1.35;
  word-wrap: break-word;
  white-space: pre-wrap;
  padding-right: 40px;
  color: var(--text);
}
</style>
