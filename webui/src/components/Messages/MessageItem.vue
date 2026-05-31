<script>
import logger from "@/utils/logger";
import ContextMenu from "@/components/Messages/ContextMenu/ContextMenu.vue";
import MessageReactions from "@/components/Messages/MessageReactions.vue";
import MessageCheckIcon from "@/components/Messages/MessageCheckIcon.vue";

import { getTime } from "@/utils/messages";
import { getIcon } from "@/state/theme";
import { user } from "@/state/user";
import { expandUrl } from "@/utils/media";
import { setMessageStatusAsRead } from "@/services/messages";
import { setImageModal } from "@/state/imageModal";

export default {
  components: { ContextMenu, MessageReactions, MessageCheckIcon },

  props: {
    message: { type: Object, required: true },
  },

  data() {
    return {
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

    onRightClick(e) {
      e.preventDefault();
      this.menu = {
        visible: true,
        canClose: false,
        x: e.clientX,
        y: e.clientY,
      };

      this.$nextTick(() => {
        setTimeout(() => {
          this.menu.canClose = true;
        }, 0);
      });
    },

    closeMenu() {
      this.menu.visible = false;
    },
  },

  emits: [
    "reactToMessage",
    "replyToMessage",
    "forwardMessage",
    "showInfoMessage",
    "deleteMessage",
  ],

  async mounted() {
    if (this.message.senderId !== user.userId) {
      try {
        await setMessageStatusAsRead(this.message.messageId);
      } catch (e) {
        logger.error(e);
      } finally {
        return;
      }
    }
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

      <MessageCheckIcon v-if="isMine" :messageId="message.messageId" />

      <MessageReactions :messageId="message.messageId" />

      <Transition name="context">
        <ContextMenu
          v-if="menu.visible"
          :message="message"
          :canClose="menu.canClose"
          :x="menu.x"
          :y="menu.y"
          @close="closeMenu"
          @reactToMessage="$emit('reactToMessage', $event)"
          @replyToMessage="$emit('replyToMessage', $event)"
          @forwardMessage="$emit('forwardMessage', $event)"
          @showInfoMessage="$emit('showInfoMessage', $event)"
          @deleteMessage="$emit('deleteMessage', $event)"
        />
      </Transition>
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

.context-enter-active,
.context-leave-active {
  transition: all 180ms ease;
  transform-origin: top;
}

.context-enter-from {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}

.context-enter-to {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.context-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.context-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.98);
}
</style>
