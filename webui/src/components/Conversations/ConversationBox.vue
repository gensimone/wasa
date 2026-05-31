<script>
import MessageList from "@/components/Messages/MessageList.vue";
import ConversationInput from "@/components/Conversations/ConversationInput.vue";
import { defaultUserPhotoUrl, defaultGroupPhotoUrl } from "@/assets/default";
import { expandUrl } from "@/utils/media";
import { getIcon } from "@/state/theme";
import { setImageModal } from "@/state/imageModal";
import { users } from "@/state/users";
import { groups, directMessages, groupMessages } from "@/state/conversations";
import { deleteMessage } from "@/services/messages";

export default {
  components: { MessageList, ConversationInput },

  data() {
    return {
      scrollTick: 0,
    };
  },

  props: {
    id: { type: Number, required: true },
    direct: { type: Boolean, required: true },
  },

  emits: ["react", "replyToMessage", "showInfoMessage"],

  computed: {
    groups() {
      return groups.value;
    },

    messages() {
      const messages = this.direct
        ? directMessages.value.get(this.id)
        : groupMessages.value.get(this.id);

      return messages || [];
    },

    name() {
      const name = this.direct
        ? users.value.get(this.id)?.name
        : groups.value.get(this.id)?.name;

      return name || "";
    },

    photoUrl() {
      return this.direct
        ? users.value.get(this.id)?.photoUrl || defaultUserPhotoUrl
        : groups.value.get(this.id)?.photoUrl || defaultGroupPhotoUrl;
    },
  },

  watch: {
    groups: {
      handler(groups) {
        if (!this.direct && !groups.has(this.id)) {
          this.$router.push("/home");
        }
      },
      deep: true,
    },
  },

  methods: {
    expandUrl,
    getIcon,
    setImageModal,

    openGroupSettings() {
      this.$router.push({
        name: "group",
        params: { id: this.id },
      });
    },

    forwardMessage(message) {
      this.$router.push({
        name: "message-forward",
        params: { id: message.messageId },
      });
    },

    async deleteMessage(message) {
      try {
        await deleteMessage(message.messageId);
        if (this.direct) {
          directMessages.value.set(
            this.id,
            directMessages.value
              .get(this.id)
              .filter((m) => m.messageId !== message.messageId),
          );
        } else {
          groupMessages.value.set(
            this.id,
            groupMessages.value
              .get(this.id)
              .filter((m) => m.messageId !== message.messageId),
          );
        }
      } catch (e) {
        handleError(e);
      }
    },
  },
};
</script>

<template>
  <div class="conversation-box">
    <div class="conversation-box-header">
      <img
        class="conversation-box-photo"
        :src="expandUrl(photoUrl)"
        @click="setImageModal(expandUrl(photoUrl))"
      />
      <div class="conversation-box-name">
        {{ name }}
      </div>
      <div class="conversation-box-info-button">
        <button
          v-if="!direct"
          class="conversation-box-info-btn"
          @click="openGroupSettings"
        >
          <img :src="getIcon('info')" class="icon-img" />
        </button>
      </div>
    </div>

    <MessageList
      :messages="messages"
      :scrollTick="scrollTick"
      @react="$emit('react', $event)"
      @replyToMessage="$emit('replyToMessage', $event)"
      @forwardMessage="forwardMessage"
      @showInfoMessage="$emit('showInfoMessage', $event)"
      @deleteMessage="deleteMessage"
    />

    <ConversationInput
      @triggerScrolldown="scrollTick++"
      :direct="direct"
      :id="id"
    />
  </div>
</template>

<style scoped>
.conversation-box {
  width: min(800px, 75%);
  height: calc(100vh - 150px);
  display: flex;
  flex-direction: column;
  border-radius: 22px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);
}

.conversation-box-header {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;

  padding: 15px 16px;
  border-bottom: 1px solid var(--border);
  background: transparent;
}

.conversation-box-photo {
  position: absolute;
  left: 16px;

  width: 45px;
  height: 45px;
  border-radius: 50%;
  object-fit: cover;

  border: 1px solid var(--border);
}

.conversation-box-name {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text);
}

.conversation-box-info-button {
  position: absolute;
  right: 16px;
}

.conversation-box-info-btn {
  width: 46px;
  height: 46px;
  border-radius: 100px;
  border: 0px;
  background: rgba(0, 0, 0, 0);

  display: flex;
  justify-content: center;
  align-items: center;

  cursor: pointer;
  position: relative;
  overflow: hidden;

  transition:
    transform 0.25s ease,
    border 0.25s ease,
    box-shadow 0.25s ease;
}

.conversation-box-info-btn:hover::before {
  transform: translateX(140%);
}

.conversation-box-info-btn:hover {
  transform: translateY(-4px) scale(1.05);
}

.conversation-box-info-btn:active {
  transform: scale(0.95);
}
</style>
