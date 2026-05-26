<script>
import MemberItem from "@/components/Groups/MemberItem.vue";
import { getIcon } from "@/state/theme";

export default {
  components: { MemberItem },

  props: {
    members: { type: Array, required: true },
    founderId: { required: true },
  },

  emits: ["removeUser", "addToGroup", "selectUser"],

  methods: {
    getIcon,
  },
};
</script>

<template>
  <div class="members-list">
    <div class="members-list-header">
      <h2>Members</h2>

      <button class="icon-btn" @click="$router.push('/group/add')">
        <img :src="getIcon('plus')" class="icon-img" />
      </button>
    </div>

    <MemberItem
      v-for="m in members"
      :key="m.userId"
      :member="m"
      @removeUser="$emit('removeUser', $event)"
      :founderId="founderId"
      @selectUser="$emit('selectUser', $event)"
    />
  </div>
</template>

<style scoped>
.members-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 22px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);
}

.members-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

.members-list-header h2 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text);
}
</style>
