<script>
import ItemPicker from "@/components/Users/ItemPicker.vue";
import { groups } from "@/state/conversations";
import { users } from "@/state/users";
import { user } from "@/state/user";

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
        .filter((u) => u.userId != user.userId)
        .filter((u) => !this.excludedUsers?.some((e) => e.userId === u.userId))
        .map((u) => ({ ...u, id: u.userId, isDirect: true }));

      const fetchedGroups = [...groups.value.values()]
        .filter(
          (g) =>
            !this.excludedGroups?.some(
              (e) => e.conversationId === g.conversationId,
            ),
        )
        .map((g) => ({ ...g, id: g.conversationId, isDirect: false }));

      if (this.includeUsers && this.includeGroups)
        return [...fetchedUsers, ...fetchedGroups];
      else if (this.includeUsers) return fetchedUsers;
      else if (this.includeGroups) return fetchedGroups;
      else return [];
    },
  },
};
</script>

<template>
  <div class="users-list">
    <ItemPicker
      :can-select-multiple="canSelectMultiple"
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

.users-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 14px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);

  max-height: 70vh; /* oppure un valore fisso tipo 400px */
  overflow-y: auto; /* abilita scroll verticale */
}
</style>
