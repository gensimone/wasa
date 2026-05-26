<script>
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import ConversationBox from "@/components/Conversations/ConversationBox.vue";

export default {
  components: { Topbar, Bottombar, ConversationBox },

  computed: {
    id() {
      return Number(this.$route.params.id);
    },

    direct() {
      return this.$route.query.direct === "true";
    },
  },

  methods: {
    react(data) {
      console.log("Message:", data.message);
      console.log("Emoji:", data.emoji);
    },

    replyToMessage(message) {
      console.log("Reply to message:", message);
    },

    forwardMessage(message) {
      console.log("Forward message:", message);
    },

    showInfoMessage(message) {
      console.log("Show info message:", message);
    },

    deleteMessage(message) {
      console.log("Delete message:", message);
    },
  },
};
</script>

<template>
  <div class="app">
    <Topbar :actions="[{ icon: 'back', onClick: () => $router.back() }]" />
    <div class="content">
      <ConversationBox
        :id="id"
        :direct="direct"
        @react="react"
        @replyToMessage="replyToMessage"
        @forwardMessage="forwardMessage"
        @showInfoMessage="showInfoMessage"
        @deleteMessage="deleteMessage"
      />
    </div>
    <Bottombar />
  </div>
</template>
