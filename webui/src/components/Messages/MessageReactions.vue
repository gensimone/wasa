<script>
import MessageReaction from "@/components/Messages/MessageReaction.vue";

export default {
  components: { MessageReaction },

  props: {
    reactions: { type: Array },
  },

  data() {
    return {
      reactionsToShow: [],
    };
  },

  watch: {
    reactions(newReactions) {
      const reducedReactions = new Map();

      for (const { emoji } of newReactions) {
        reducedReactions.set(emoji, {
          emoji,
          count: (reducedReactions.get(emoji)?.count ?? 0) + 1,
        });
      }

      this.reactionsToShow = [...reducedReactions.values()];
    },
  },
};
</script>

<template>
  <div v-if="reactions.length" class="message-item-reactions">
    <MessageReaction
      v-for="(r, index) in reactionsToShow"
      :emoji="r.emoji"
      :count="r.count"
      :key="index"
      :style="{ zIndex: reactions.length - index }"
    />
  </div>
</template>

<style scoped>
.message-item-reactions {
  position: absolute;
  bottom: -13px;
  gap: 15px;

  display: flex;
  align-items: center;

  pointer-events: none;
}
</style>
