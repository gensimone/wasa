<script>
export default {
  props: {
    message: { type: Object, required: true },
  },

  emits: ["react"],

  methods: {
    react(emoji) {
      this.$emit("react", {
        message: this.message,
        emoji: emoji,
      });
    },
  },
};
</script>

<template>
  <div class="reaction-bar">
    <button @click="react('angry')">
      <img class="reaction-btn" src="/icons/reactions/angry.svg" />
    </button>
    <button @click="react('laugh')">
      <img class="reaction-btn" src="/icons/reactions/laugh.svg" />
    </button>
    <button @click="react('like')">
      <img class="reaction-btn" src="/icons/reactions/like.svg" />
    </button>
    <button @click="react('love')">
      <img class="reaction-btn" src="/icons/reactions/love.svg" />
    </button>
    <button @click="react('sad')">
      <img class="reaction-btn" src="/icons/reactions/sad.svg" />
    </button>
  </div>
</template>

<style scoped>
.reaction-bar {
  display: flex;
  gap: 6px;
  padding: 6px;
}

.reaction-bar button {
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
}

.reaction-btn {
  width: 35px;
  height: 35px;
  display: block;

  transition: transform 0.18s ease;
}

.reaction-bar button:hover .reaction-btn {
  transform: translateY(-4px) scale(1.08);
}

.reaction-bar button:active .reaction-btn {
  transform: scale(0.95);
}
</style>
