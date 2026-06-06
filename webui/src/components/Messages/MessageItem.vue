<script>
import logger from "@/utils/logger";
import Poller from "@/services/poller";
import ContextMenu from "@/components/Messages/ContextMenu/ContextMenu.vue";
import MessageInfo from "@/components/Messages/MessageInfo.vue";
import MessageReactions from "@/components/Messages/MessageReactions.vue";
import MessageCheckIcon from "@/components/Messages/MessageCheckIcon.vue";
import { getTime } from "@/utils/messages";
import { getIcon } from "@/state/theme";
import { directMessages, groupMessages } from "@/state/conversations";
import { user } from "@/state/user";
import { expandUrl } from "@/utils/media";
import { setMessageStatusAsRead } from "@/services/messages";
import { setImageModal } from "@/state/imageModal";
import { handleError } from "@/utils/errors";
import { getCheckIcon } from "@/utils/messages";
import { getReceipts } from "@/services/messages";
import {
  addReaction,
  getReactions,
  deleteReaction,
} from "@/services/reactions";

export default {
  components: { ContextMenu, MessageReactions, MessageCheckIcon, MessageInfo },

  props: {
    message: { type: Object, required: true },
    id: { type: Number, required: true },
    direct: { type: Boolean, required: true },
  },

  emits: ["replyToMessage", "forwardMessage", "deleteMessage", "jumpToMessage"],

  data() {
    return {
      checkIcon: null,

      reactionsPoller: null,
      receiptsPoller: null,

      reactions: [],
      receipts: [],

      showInfoMessage: false,

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

    commentedMessage() {
      if (!this.message.commentTo) return null;

      const messages = this.direct
        ? directMessages.value.get(this.id)
        : groupMessages.value.get(this.id);

      const commentedMessage = messages.find(
        (m) => m.messageId === this.message.commentTo,
      );

      return commentedMessage;
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

  methods: {
    expandUrl,
    getIcon,
    setImageModal,
    getTime,

    async reactToMessage(reactionData) {
      const message = reactionData.message;
      const emoji = reactionData.emoji;

      try {
        const reaction = await addReaction(message.messageId, emoji);
        this.reactions = [
          reaction,
          ...this.reactions.filter((r) => r.senderId !== user.userId),
        ];
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

    async deleteReaction() {
      try {
        await deleteReaction(this.message.messageId);
        this.reactions = this.reactions.filter(
          (r) => r.senderId !== user.userId,
        );
      } catch (e) {
        if (e.response?.status !== 404) {
          handleError(e);
        }
      }
    },
  },
};
</script>
<template>
  <div
    :id="`message-${message.messageId}`"
    class="message-item"
    :class="{ mine: isMine }"
  >
    <div class="message-item-bubble" @contextmenu="onRightClick">
      <div
        v-if="commentedMessage"
        class="message-item-comment-preview"
        @click="$emit('jumpToMessage', commentedMessage)"
      >
        <div v-if="commentedMessage.text" class="message-item-comment-text">
          {{ commentedMessage.text }}
        </div>

        <div
          v-if="commentedMessage.attachmentUrl"
          class="message-item-comment-attachment"
        >
          <img :src="expandUrl(commentedMessage.attachmentUrl)">
        </div>
      </div>

      <img
        v-if="message.isForwarded"
        class="message-item-forward-icon"
        :src="getIcon('forward')"
        alt="forwarded"
      >

      <div v-if="message.text" class="message-item-text">
        {{ message.text }}
      </div>

      <img
        v-if="message.attachmentUrl"
        :src="expandUrl(message.attachmentUrl)"
        class="message-item-image"
        @click="setImageModal(expandUrl(message.attachmentUrl))"
      >

      <div class="message-item-meta">
        <span class="message-item-time">
          {{ getTime(message.createdAt) }}
        </span>
      </div>

      <MessageCheckIcon v-if="isMine" :check-icon="checkIcon" />

      <MessageReactions :reactions="reactions" />

      <ContextMenu
        v-if="menu.visible"
        :message="message"
        :can-close="menu.canClose"
        :x="menu.x"
        :y="menu.y"
        @close="closeMenu"
        @react-to-message="reactToMessage"
        @show-info-message="showInfoMessage = true"
        @reply-to-message="$emit('replyToMessage', $event)"
        @forward-message="$emit('forwardMessage', $event)"
        @delete-message="$emit('deleteMessage', $event)"
      />
    </div>
  </div>
  <MessageInfo
    v-if="showInfoMessage"
    :receipts="receipts"
    :reactions="reactions"
    @close-message-info="showInfoMessage = false"
    @delete-reaction="deleteReaction"
  />
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

.message-item-comment-preview {
  display: flex;
  align-items: center;
  gap: 8px;

  border-left: 3px solid var(--accent-strong);
  padding: 6px 10px;
  margin-bottom: 8px;
  border-radius: 8px;

  background: rgba(255, 255, 255, 0.05);
  cursor: pointer;
}

.message-item-comment-text {
  flex: 1;

  font-size: 0.8rem;
  opacity: 0.7;

  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;

  color: var(--text);
}

.message-item-comment-attachment img {
  width: 48px;
  height: 48px;

  object-fit: cover;
  border-radius: 6px;

  flex-shrink: 0;
}

.message-highlight {
  animation: message-highlight 1.5s ease;
}

@keyframes message-highlight {
  0% {
    background: var(--accent);
  }

  40% {
    background: var(--accent-strong);
  }

  100% {
    background: inherit;
  }
}
</style>
