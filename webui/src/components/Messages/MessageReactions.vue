<script>
import Poller from "@/services/poller";
import MessageReaction from "@/components/Messages/MessageReaction.vue";
import { getReactions } from "@/services/reactions";
import logger from "@/utils/logger";

export default {
  components: { MessageReaction },

  props: {
    messageId: { type: Number, required: true },
  },

  data() {
    return {
      reactions: [],
      poller: null,
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

    async fetchReactions() {
      try {
        const reactions = (await getReactions(this.messageId)) || [];
        this.reactions = this.reduceEmojies(reactions);
      } catch (e) {
        if (e.response?.status === 404) {
          // The message has been deleted
          console.log("Message deleted!");
          this.poller?.stopPolling();
          this.reactions = [];
        } else {
          logger.error(e);
        }
      }
    },
  },

  async mounted() {
    this.poller = new Poller();
    this.poller.callback = this.fetchReactions;

    this.poller.startPolling();
  },

  beforeUnmount() {
    this.poller?.stopPolling();
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
