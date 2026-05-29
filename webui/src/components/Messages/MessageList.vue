<script>
import MessageItem from "@/components/Messages/MessageItem.vue";

export default {
  components: { MessageItem },

  props: {
    messages: { type: Array, required: true },
    scrollTick: { type: Number, required: true },
  },

  watch: {
    scrollTick() {
      this.scrollToBottomIfNeeded();
    },
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
  },

  emits: [
    "react",
    "replyToMessage",
    "forwardMessage",
    "showInfoMessage",
    "deleteMessage",
  ],

  mounted() {
    this.scrollToBottomInstant();
  },
};
</script>

<template>
  <div class="message-list" ref="container">
    <MessageItem
      v-for="m in messages"
      :key="m.messageId"
      :message="m"
      @react="$emit('react', $event)"
      @replyToMessage="$emit('replyToMessage', $event)"
      @forwardMessage="$emit('forwardMessage', $event)"
      @showInfoMessage="$emit('showInfoMessage', $event)"
      @deleteMessage="$emit('deleteMessage', $event)"
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
