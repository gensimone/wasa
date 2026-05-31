<script>
export default {
  props: {
    emoji: { type: String, required: true },
    count: { type: Number, required: true },
  },

  methods: {
    getIconFromEmoji(emoji) {
      return `/icons/reactions/${emoji}.svg`;
    },
  },
};
</script>

<template>
  <div class="message-item-reaction">
    <img :src="getIconFromEmoji(emoji)" class="message-item-emoji" />
    <span v-if="count > 1" class="message-item-reaction-count">
      {{ count }}
    </span>
  </div>
</template>

<style scoped>
.message-item-reaction {
  position: relative;
  width: 24px;
  height: 24px;

  display: flex;
  align-items: center;
  justify-content: center;

  margin-left: -6px;

  border-radius: 999px;

  background: var(--surface);
  border: 1px solid var(--border);

  font-size: 14px;

  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.message-item-reaction:first-child {
  margin-left: 0;
}

.message-item-reaction-count {
  position: absolute;
  bottom: -15px;
  right: -12px;

  width: 10px;
  height: 10px;
  padding: 10px;

  display: flex;
  align-items: center;
  justify-content: center;

  border-radius: 999px;
  font-size: 5px;

  background: var(--surface-3);
  color: var(--text);

  border: 1px solid var(--text);
}

.message-item-emoji {
  width: 35px;
  height: 35px;
}
</style>
