<script>
import UsersPicker from "@/components/Users/UsersPicker.vue";
import Poller from "@/services/poller";
import { getUsers } from "@/services/users";
import { user } from "@/state/user";

export default {
  components: { UsersPicker },

  data() {
    return {
      users: [],
      poller: null,
    };
  },

  props: {
    canSelectMultiple: { type: Boolean, required: true },
    excludeUsers: { type: Array, required: false },
  },

  emits: ["select"],

  async mounted() {
    this.poller = new Poller(async () => {
      const users = await getUsers();
      this.users = users
        .filter((u) => u.userId != user.userId)
        .filter((u) => !this.excludeUsers?.some((e) => e.userId === u.userId));
    });

    this.poller.startPolling();
  },

  beforeUnmount() {
    this.poller?.stopPolling();
  },
};
</script>

<template>
  <UsersPicker
    :canSelectMultiple="canSelectMultiple"
    :users="users"
    @select="$emit('select', $event)"
  />
</template>
