<script>
import { getIcon } from "@/state/theme";
import { handleError } from "@/utils/errors";
import { sendMessage } from "@/services/conversations";
import { directMessages, groupMessages } from "@/state/conversations";

export default {
  props: {
    id: { type: Number, required: true },
    direct: { type: Boolean, required: true },
  },

  data() {
    return {
      text: "",
      attachment: null,
      attachmentUrl: null,
      sending: false,
    };
  },

  emits: ["triggerScrolldown"],

  methods: {
    getIcon,

    async sendMessage() {
      const cleanText = this.text?.trim();
      if (!cleanText && !this.attachment) return;

      this.sending = true;

      try {
        const message = await sendMessage(
          this.id,
          this.direct,
          this.text,
          this.attachment,
        );

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
        handleError(e);
      } finally {
        this.sending = false;
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
  },

  beforeUnmount() {
    this.revokeAttachment();
  },
};
</script>

<template>
  <div>
    <div v-if="attachment" class="conversation-input-attachment-preview">
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
        @keydown.enter.prevent="sendMessage"
      />

      <button class="icon-btn" :disabled="sending" @click="sendMessage">
        <img :src="getIcon('send')" class="icon-img" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.conversation-input {
  display: flex;
  gap: 10px;
}

.conversation-input-attachment-preview {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 0px auto 0;
  padding: 8px 10px;
  width: min(720px, 100%);
  border-radius: 0px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.conversation-input-preview-img {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  object-fit: cover;
  border: 1px solid rgba(255, 255, 255, 0.12);
}
</style>
