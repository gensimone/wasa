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
    replyToMessage(message) {
      // Replying to a message basically put the message to reply in the same spot in
      // which we already have put the attachment window. After that the user can choose
      // to delete the "message attachment", or to write a comment and then send the
      // message. The commented message must appear above the comment with a snippet
      // (if it as an image, that image must be also present). Clicking on the comment
      // snippet must change the MessageList position so that the commented message is
      // clearly visible.
      console.log("Reply to message:", message);
    },

    showInfoMessage(message) {
      // Showing the info of the message is perhaps the most problematic one.
      // We must first decide how to show the informations.
      // First of all the informations that we want to show are the following:
      // 1) Reactions:
      //    - emoji.
      //    - list of users that sended that emoji.
      // 2) Receipts (only if user.userId == message.senderId):
      //    - user.
      //    - one of "seen at $DATE_TIME" and "received at $DATE_TIME"
      console.log("Show info message:", message);
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
        @replyToMessage="replyToMessage"
        @showInfoMessage="showInfoMessage"
      />
    </div>
    <Bottombar />
  </div>
</template>
