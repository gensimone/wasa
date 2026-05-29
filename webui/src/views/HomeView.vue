<script>
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import ConversationsList from "@/components/Conversations/ConversationsList.vue";

import { userId } from "@/state/users";

export default {
  components: { Topbar, Bottombar, ConversationsList },

  methods: {
    select(conversation) {
      this.$router.push({
        name: "conversation",
        params: { id: conversation.id },
        query: { direct: !conversation.isGroup },
      });
    },

    logout() {
      userId.value = null;
      this.$router.push("/");
    },
  },
};
</script>

<template>
  <div class="app">
    <Topbar
      :actions="[
        { icon: 'settings', onClick: () => $router.push('/user/settings') },
        { icon: 'logout', onClick: () => logout() },
      ]"
    />
    <div class="content">
      <ConversationsList @select="select" />
    </div>
    <Bottombar />
  </div>
</template>
