<script>
import MessageItem from "@/components/Messages/MessageItem.vue";

export default {
  components: { MessageItem },

  props: {
    messages: { type: Array, required: true },
    scrollTick: { type: Number, required: true },
    id: { type: Number, required: true },
    direct: { type: Boolean, required: true },
  },

  emits: ["replyToMessage", "forwardMessage", "deleteMessage"],

  data() {
    return {
      inShowInfoMessage: false,
      receipts: [],
      reactions: [],
    };
  },

  watch: {
    scrollTick() {
      this.scrollToBottomIfNeeded();
    },
  },

  mounted() {
    this.scrollToBottomInstant();
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

    jumpToMessage(message) {
      const el = document.getElementById(`message-${message.messageId}`);

      if (!el) return;

      el.scrollIntoView({
        behavior: "smooth",
        block: "center",
      });

      el.classList.add("message-highlight");

      setTimeout(() => {
        el.classList.remove("message-highlight");
      }, 1500);
    },
  },
};
</script>

<template>
  <div ref="container" class="message-list">
    <MessageItem
      v-for="m in messages"
      :id="id"
      :key="m.messageId"
      :message="m"
      :direct="direct"
      @reply-to-message="$emit('replyToMessage', $event)"
      @forward-message="$emit('forwardMessage', $event)"
      @delete-message="$emit('deleteMessage', $event)"
      @jump-to-message="jumpToMessage"
    />
    <div ref="bottom" />
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
