<script>
import { getIcon } from "@/state/theme";

export default {
  props: {
    checkIcon: { type: String },
  },

  methods: {
    getIcon,
  },
};
</script>

<template>
  <img class="message-item-check-icon" :src="getIcon(checkIcon)" alt="">
</template>

<style scoped>
.message-item-check-icon {
  width: 16px;
  height: 16px;
  opacity: 0.9;
}
</style>
