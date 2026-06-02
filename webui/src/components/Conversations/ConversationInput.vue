<script>
import { expandUrl } from "@/utils/media";
import { getIcon } from "@/state/theme";
import { handleError } from "@/utils/errors";
import { sendMessage } from "@/services/conversations";
import { commentMessage } from "@/services/comments";
import { directMessages, groupMessages } from "@/state/conversations";

export default {
  props: {
    id: { type: Number, required: true },
    direct: { type: Boolean, required: true },

    message: { type: Object, required: false },
  },

  watch: {
    message(v) {
      this.messageToComment = v;
    },
  },

  data() {
    return {
      text: "",
      attachment: null,
      attachmentUrl: null,
      sending: false,

      messageToComment: null,
    };
  },

  emits: ["triggerScrolldown"],

  methods: {
    expandUrl,
    getIcon,

    async handleSendMessage() {
      const cleanText = this.text?.trim();
      if (!cleanText && !this.attachment) return;

      try {
        let message;

        if (this.messageToComment) {
          message = await commentMessage(
            this.messageToComment.messageId,
            this.text,
            this.attachment,
          );
          this.messageToComment = null;
        } else {
          message = await sendMessage(
            this.id,
            this.direct,
            this.text,
            this.attachment,
          );
        }

        if (this.direct) {
          if (directMessages.value.has(this.id)) {
            directMessages.value.get(this.id).push(message);
          } else {
            directMessages.value.set(this.id, [message]);
          }
        } else {
          groupMessages.value.get(this.id).push(message);
        }

        this.$emit("triggerScrolldown");

        this.text = "";
        if (this.attachment) {
          this.removeAttachment();
        }
      } catch (e) {
        // The message to comment has been deleted.
        if (e.response?.status === 404 && this.messageToComment) {
          this.messageToComment = null;
        }

        handleError(e);
      } finally {
        this.sending = true;
      }
    },

    revokeAttachment() {
      if (this.attachmentUrl) {
        URL.revokeObjectURL(this.attachmentUrl);
      }
    },

    addAttachment(event) {
      const file = event.target.files[0];
      if (!file) return;

      this.revokeAttachment();

      this.attachment = file;
      this.attachmentUrl = URL.createObjectURL(file);

      event.target.value = "";
    },

    removeAttachment() {
      this.revokeAttachment();

      this.attachmentUrl = null;
      this.attachment = null;
    },

    abortComment() {
      this.messageToComment = null;
    },
  },

  beforeUnmount() {
    this.revokeAttachment();
  },
};
</script>

<template>
  <div>
    <div v-if="messageToComment" class="conversation-input-data-preview">
      <div class="preview-label">Comment to</div>
      <div v-if="messageToComment.attachmentUrl">
        <img
          class="conversation-input-preview-img"
          :src="expandUrl(messageToComment.attachmentUrl)"
        />
      </div>
      <div v-if="messageToComment.text">
        {{ messageToComment.text }}
      </div>

      <button class="icon-btn" @click="abortComment">
        <img :src="getIcon('trash')" class="icon-img" />
      </button>
    </div>

    <div v-if="attachment" class="conversation-input-data-preview">
      <div class="preview-label">Attachment</div>
      <img class="conversation-input-preview-img" :src="attachmentUrl" />

      <button class="icon-btn" @click="removeAttachment">
        <img :src="getIcon('trash')" class="icon-img" />
      </button>
    </div>

    <div class="conversation-input">
      <button class="icon-btn" @click="$refs.fileInput.click()">
        <img :src="getIcon('upload')" class="icon-img" />
      </button>

      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="addAttachment($event)"
      />

      <input
        name="conversation-input-bar"
        class="input-bar"
        :value="text"
        @input="text = $event.target.value"
        placeholder="Type a message..."
        @keydown.enter.prevent="handleSendMessage"
      />

      <button class="icon-btn" :disabled="sending" @click="handleSendMessage">
        <img :src="getIcon('send')" class="icon-img" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.conversation-input {
  display: flex;
  gap: 10px;
  align-items: center;
}

.conversation-input-data-preview {
  display: flex;
  gap: 10px;
  align-items: flex-start;

  margin: 0 auto 10px;
  padding: 22px 44px 10px 12px;

  width: min(720px, 100%);

  border-radius: 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(14px);

  color: rgba(255, 255, 255, 0.75);
  font-size: 13px;
  line-height: 1.4;

  position: relative;
}

.conversation-input-comment-preview div {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.conversation-input-preview-img {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  flex-shrink: 0;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.conversation-input-data-preview .icon-btn {
  position: absolute;
  top: 6px;
  right: 6px;

  width: 28px;
  height: 28px;

  display: flex;
  align-items: center;
  justify-content: center;

  border-radius: 8px;
  background: rgba(0, 0, 0, 0.35);
  backdrop-filter: blur(8px);

  transition:
    transform 0.15s ease,
    background 0.15s ease;
}

.conversation-input-comment-preview .icon-btn:hover {
  transform: scale(1.05);
  background: rgba(0, 0, 0, 0.5);
}

.input-bar {
  flex: 1;
  padding: 10px 12px;

  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);

  background: rgba(255, 255, 255, 0.03);
  color: rgba(255, 255, 255, 0.9);

  outline: none;
  transition:
    border 0.15s ease,
    background 0.15s ease;
}

.input-bar:focus {
  border: 1px solid rgba(255, 255, 255, 0.22);
  background: rgba(255, 255, 255, 0.05);
}

.preview-label {
  position: absolute;
  top: 0;
  left: 12px;

  transform: translateY(-50%);

  padding: 3px 10px;
  border-radius: 999px;

  font-size: 11px;
  font-weight: 500;
  letter-spacing: 0.3px;

  background: rgba(0, 0, 0, 0.55);
  color: rgba(255, 255, 255, 0.85);

  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);

  pointer-events: none;
}

.conversation-input-comment-preview {
  position: relative;
}
</style>
