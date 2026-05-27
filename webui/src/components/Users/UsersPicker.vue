<script>
import UserItem from "@/components/Users/UserItem.vue";

export default {
  components: { UserItem },

  props: {
    users: { type: Array, required: true },
    canSelectMultiple: { type: Boolean, required: true },
  },

  emits: ["select"],

  data() {
    return {
      selectedUsers: [],
      query: "",
    };
  },

  computed: {
    usersToShow() {
      const unselectedUsers = this.users.filter(
        (u) =>
          !this.selectedUsers.some((selected) => selected.userId === u.userId),
      );

      if (!this.query.trim()) return unselectedUsers;

      return unselectedUsers.filter((u) =>
        u.name.toLowerCase().includes(this.query.toLowerCase()),
      );
    },
  },

  methods: {
    onUserSelect(user) {
      if (!this.canSelectMultiple) {
        this.$emit("select", user);
        return;
      }

      const alreadySelected = this.selectedUsers.some(
        (u) => u.userId === user.userId,
      );

      if (alreadySelected) {
        this.selectedUsers = this.selectedUsers.filter(
          (u) => u.userId !== user.userId,
        );
      } else {
        this.selectedUsers.push(user);
      }
    },

    confirmSelection() {
      this.$emit("select", this.selectedUsers);
    },

    isSelected(user) {
      return this.selectedUsers.some((u) => u.userId === user.userId);
    },
  },
};
</script>

<template>
  <div class="user-picker">
    <input
      class="input-bar"
      placeholder="Search.."
      @input="query = $event.target.value"
    />
    <div class="user-picker-list">
      <UserItem
        v-for="user in usersToShow"
        :key="user.userId"
        :user="user"
        :selected="isSelected(user)"
        @select="onUserSelect"
      />
    </div>

    <div v-if="canSelectMultiple && !query" class="user-picker-sidebar">
      <div v-if="this.selectedUsers.length" class="selected-title">
        Selected users
      </div>

      <UserItem
        v-for="user in selectedUsers"
        :key="user.userId"
        :user="user"
        :selected="isSelected(user)"
        @select="onUserSelect"
      />

      <button
        v-if="selectedUsers.length"
        class="submit-button"
        @click="confirmSelection"
      >
        Confirm
      </button>
    </div>
  </div>
</template>

<style scoped>
.selected-title {
  font-size: 1.5rem;
  color: var(--text);
  text-transform: uppercase;
}
</style>
