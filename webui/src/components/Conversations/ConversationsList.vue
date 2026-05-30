<script>
import { getIcon } from "@/state/theme";
import { users } from "@/state/users";
import { groupMessages, directMessages, groups } from "@/state/conversations";
import ConversationItem from "@/components/Conversations/ConversationItem.vue";

export default {
  components: { ConversationItem },

  emits: ["select"],

  methods: {
    getIcon,
  },

  computed: {
    conversations() {
      let conversations = [];
      for (const [id, messages] of directMessages.value) {
        const user = users.value.get(id);
        conversations.push({
          id: id,
          name: user.name,
          photoUrl: user.photoUrl,
          lastMessage: messages?.at(-1),
          isGroup: false,
        });
      }

      for (const [id, messages] of groupMessages.value) {
        const group = groups.value.get(id);
        conversations.push({
          id: id,
          name: group.name,
          photoUrl: group.photoUrl,
          lastMessage: messages?.at(-1),
          isGroup: true,
        });
      }

      return conversations;
    },
  },
};
</script>

<template>
  <div class="conversations-list">
    <div class="conversations-list-header">
      <h2>Conversations</h2>

      <button class="icon-btn" @click="$router.push('/conversation/add')">
        <img :src="getIcon('plus')" class="icon-img" />
      </button>
    </div>

    <ConversationItem
      v-for="c in conversations"
      :key="`${c.isGroup ? 'g' : 'u'}-${c.id}`"
      :conversation="c"
      @select="$emit('select', c)"
    />
  </div>
</template>

<style scoped>
.conversations-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 22px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);
}

.conversations-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.conversations-list-header h2 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text);
}
</style>
