<script>
import Poller from "@/services/poller";
import MessageReaction from "@/components/Messages/MessageReaction.vue";
import { getReactions } from "@/services/reactions";

export default {
  components: { MessageReaction },

  props: {
    messageId: { type: Number, required: true },
  },

  data() {
    return {
      reactions: [],
    };
  },

  methods: {
    reduceEmojies(reactions) {
      const reducedReactions = new Map();

      for (const { emoji } of reactions) {
        reducedReactions.set(emoji, {
          emoji,
          count: (reducedReactions.get(emoji)?.count ?? 0) + 1,
        });
      }

      return [...reducedReactions.values()];
    },
  },

  async mounted() {
    this.reactionsPoller = new Poller(async () => {
      const reactions = (await getReactions(this.messageId)) || [];

      this.reactions = this.reduceEmojies(reactions);
    });

    this.reactionsPoller.startPolling();
  },

  beforeUnmount() {
    this.reactionsPoller?.stopPolling();
  },
};
</script>

<template>
  <div v-if="reactions.length" class="message-item-reactions">
    <MessageReaction
      v-for="(r, index) in reactions"
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
  bottom: -14px;

  display: flex;
  align-items: center;

  pointer-events: none;
}
</style>
