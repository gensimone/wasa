<script>
import Bottombar from "@/components/Shared/Bottombar.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import UsersList from "@/components/Users/UsersList.vue";

export default {
  components: { Bottombar, Topbar, UsersList },
};
</script>

<template>
  <div class="app">
    <Topbar
      :actions="[{ icon: '/icons/back.svg', onClick: () => $router.back() }]"
    />
    <div class="content">
      <UsersList :users="users" @select="startConversation" />
    </div>
    <Bottombar />
  </div>
</template>
