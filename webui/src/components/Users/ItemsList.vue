<script>
import ItemPicker from "@/components/Users/ItemPicker.vue";
import { groups } from "@/state/conversations";
import { users, userId } from "@/state/user";

export default {
  components: { ItemPicker },

  props: {
    canSelectMultiple: { type: Boolean, required: true },
    excludedUsers: { type: Array, required: false },
    excludedGroups: { type: Array, required: false },
    includeUsers: { type: Boolean, required: true },
    includeGroups: { type: Boolean, required: true },
  },

  emits: ["select"],

  computed: {
    items() {
      const fetchedUsers = [...users.value.values()]
        .filter((u) => u.userId != userId.value)
        .filter((u) => !this.excludedUsers?.some((e) => e.userId === u.userId))
        .map((u) => ({ ...u, id: u.userId, isGroup: false }));

      const fetchedGroups = [...groups.value.values()]
        .filter(
          (g) =>
            !this.excludedGroups?.some(
              (e) => e.conversationId === g.conversationId,
            ),
        )
        .map((g) => ({ ...g, id: g.conversationId, isGroup: true }));

      if (this.includeUsers && this.includeGroups)
        this.items = [...fetchedUsers, ...fetchedGroups];
      else if (this.includeUsers) this.items = fetchedUsers;
      else if (this.includeGroups) this.items = fetchedGroups;
      else this.items = [];
    },
  },
};
</script>

<template>
  <div class="users-list">
    <ItemPicker
      :canSelectMultiple="canSelectMultiple"
      :items="items"
      @select="$emit('select', $event)"
    />
  </div>
</template>

<style scoped>
.users-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 14px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);
}
</style>
