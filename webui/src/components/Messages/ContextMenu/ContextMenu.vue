<script>
import ReactionBar from "./ReactionBar.vue";
import MessageMenu from "./MessageMenu.vue";

export default {
  components: { ReactionBar, MessageMenu },

  props: {
    x: { type: Number, required: true },
    y: { type: Number, required: true },
    message: { type: Object, required: true },
    canClose: { type: Boolean, required: true },
  },

  emits: [
    "close",

    "reactToMessage",
    "replyToMessage",
    "forwardMessage",
    "showInfoMessage",
    "deleteMessage",
  ],

  mounted() {
    window.addEventListener("click", this.onOutsideClick);
    window.addEventListener("contextmenu", this.onOutsideClick);
    window.addEventListener("keydown", this.onKeyDown);
  },

  beforeUnmount() {
    window.removeEventListener("click", this.onOutsideClick);
    window.removeEventListener("contextmenu", this.onOutsideClick);
    window.removeEventListener("keydown", this.onKeyDown);
  },

  methods: {
    onOutsideClick() {
      if (this.canClose) this.$emit("close");
    },

    onKeyDown(e) {
      if (e.key === "Escape") {
        this.$emit("close");
      }
    },
  },
};
</script>

<template>
  <div
    v-if="true"
    class="context-menu"
    :style="{ top: y + 'px', left: x + 'px' }"
  >
    <ReactionBar
      :message="message"
      @react-to-message="$emit('reactToMessage', $event)"
    />
    <MessageMenu
      :message="message"
      @reply-to-message="$emit('replyToMessage', $event)"
      @forward-message="$emit('forwardMessage', $event)"
      @show-info-message="$emit('showInfoMessage', $event)"
      @delete-message="$emit('deleteMessage', $event)"
    />
  </div>
</template>

<style scoped>
.context-menu {
  position: fixed;
  background: var(--surface-3);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  min-width: 180px;
  z-index: 2;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 1);
}
</style>
