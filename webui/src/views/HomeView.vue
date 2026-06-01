<script>
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import ConversationsList from "@/components/Conversations/ConversationsList.vue";

import { clearUserState, stopPollingUser } from "@/state/user";

export default {
  components: { Topbar, Bottombar, ConversationsList },

  methods: {
    select(conversation) {
      this.$router.push({
        name: "conversation",
        params: { id: conversation.id },
        query: { direct: conversation.isDirect },
      });
    },

    logout() {
      stopPollingUser();
      clearUserState();
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
