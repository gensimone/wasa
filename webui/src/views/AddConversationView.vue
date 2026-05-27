<script>
import { getIcon } from "@/state/theme";
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import UsersList from "@/components/Users/UsersList.vue";

export default {
  components: { Topbar, Bottombar, UsersList },

  methods: {
    getIcon,

    async select(user) {
      this.$router.push({
        name: "conversation",
        params: { id: user.userId },
        query: { direct: true },
      });
    },
  },
};
</script>

<template>
  <div class="app">
    <Topbar :actions="[{ icon: 'back', onClick: () => $router.back() }]" />
    <div class="content">
      <div class="items-list">
        <div class="list-item" @click="$router.push('/group/create')">
          <img :src="getIcon('plus')" class="icon-img" />
          <div class="item-info">
            <div class="item-name">Create a new group</div>
          </div>
        </div>
        <UsersList @select="select" :canSelectMultiple="false" />
      </div>
    </div>
    <Bottombar />
  </div>
</template>

<style scoped>
.items-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 22px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);
}

.list-item {
  display: flex;
  align-items: center;
  gap: 14px;

  padding: 14px;
  margin-bottom: 10px;

  border-radius: 18px;

  background: var(--surface-2);
  border: 1px solid var(--border);

  cursor: pointer;
  overflow: hidden;

  animation: fadeInUp 0.35s ease both;
}

.list-item:hover {
  transform: translateY(-6px) scale(1.02);
  border: 1px solid var(--accent);
}

.item-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 1.05rem;
  font-weight: 800;

  color: var(--text);

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
