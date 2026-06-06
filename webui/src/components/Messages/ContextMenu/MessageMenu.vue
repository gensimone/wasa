<script>
import { user } from "@/state/user";
import { getIcon } from "@/state/theme";

export default {
  props: {
    message: { type: Object, required: true },
  },

  emits: [
    "replyToMessage",
    "forwardMessage",
    "showInfoMessage",
    "deleteMessage",
  ],

  computed: {
    isMine() {
      return this.message.senderId === user.userId;
    },
  },

  methods: {
    getIcon,
  },
};
</script>

<template>
  <div class="menu">
    <div class="item" @click="$emit('replyToMessage', message)">
      <span class="icon">
        <img class="icon" :src="getIcon('reply')" alt="Reply">
      </span>
      <span class="label">Reply</span>
    </div>

    <div class="item" @click="$emit('forwardMessage', message)">
      <span class="icon">
        <img class="icon" :src="getIcon('forward')" alt="Forward">
      </span>
      <span class="label">Forward</span>
    </div>

    <div class="item" @click="$emit('showInfoMessage', message)">
      <span class="icon">
        <img class="icon" :src="getIcon('info')" alt="Info">
      </span>
      <span class="label">Info</span>
    </div>

    <div v-if="isMine" class="item" @click="$emit('deleteMessage', message)">
      <span class="icon">
        <img class="icon" :src="getIcon('trash')" alt="Delete">
      </span>
      <span class="label">Delete</span>
    </div>
  </div>
</template>

<style scoped>
.menu {
  padding: 6px;
}

.item {
  padding: 8px 10px;
  cursor: pointer;
  border-radius: 6px;
}

.item:hover {
  background: rgba(0, 0, 0, 0.08);
}

.icon {
  width: 30px;
  height: 30px;
  padding-right: 5px;
}
</style>
