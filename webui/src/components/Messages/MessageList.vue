<script>
import MessageItem from "@/components/Messages/MessageItem.vue";
import MessageInfo from "@/components/Messages/MessageInfo.vue";

export default {
  components: { MessageItem, MessageInfo },

  props: {
    messages: { type: Array, required: true },
    scrollTick: { type: Number, required: true },
  },

  watch: {
    scrollTick() {
      this.scrollToBottomIfNeeded();
    },
  },

  data() {
    return {
      inShowInfoMessage: false,
      receipts: [],
      reactions: [],
    };
  },

  methods: {
    isNearBottom() {
      const el = this.$refs.container;
      return el.scrollHeight - el.scrollTop - el.clientHeight < 80;
    },

    // FIXME: It doesn't work.
    scrollToBottomIfNeeded() {
      this.$nextTick(() => {
        const bottom = this.$refs.bottom;
        if (!bottom) return;

        bottom.scrollIntoView({
          behavior: "smooth",
          block: "end",
        });
      });
    },

    scrollToBottomInstant() {
      const bottom = this.$refs.bottom;
      if (!bottom) return;

      bottom.scrollIntoView({
        behavior: "auto",
        block: "end",
      });
    },

    showInfoMessage(messageInfo) {
      this.receipts = messageInfo.receipts;
      this.reactions = messageInfo.reactions;
      this.inShowInfoMessage = true;
    },
  },

  emits: ["replyToMessage", "forwardMessage", "deleteMessage"],

  mounted() {
    this.scrollToBottomInstant();
  },
};
</script>

<template>
  <div v-if="inShowInfoMessage" class="message-list">
    <MessageInfo
      @closeMessageInfo="inShowInfoMessage = false"
      :receipts="receipts"
      :reactions="reactions"
    />
  </div>
  <div v-else class="message-list" ref="container">
    <MessageItem
      v-for="m in messages"
      :key="m.messageId"
      :message="m"
      @replyToMessage="$emit('replyToMessage', $event)"
      @forwardMessage="$emit('forwardMessage', $event)"
      @deleteMessage="$emit('deleteMessage', $event)"
      @showInfoMessage="showInfoMessage"
    />
    <div ref="bottom"></div>
  </div>
</template>

<style scoped>
.message-list {
  height: 100vh;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
}
</style>
