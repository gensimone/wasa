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
      // First of all the reactions must be visible on the message.
      // We could overlap all the emojis on the message. If more than one of the same
      // reaction is on the message, than only one is showed (this way we don't end up
      // with too much emojies on the message that will eventually overflow).
      // All the emojies (and their respective senders) are showed in the info message
      // anyway so we should not be worried about not showing all the reactions on the
      // message.
      // Anyway, reacting to a message should be visible immediatly so we must find a
      // way to add the reaction in the data structure that holds them and that is
      // constantly updated by a mean of polling. The mini component that show the reactions
      // on the message must watch that data structure so it is updated automatically.
      console.log("Message:", data.message);
      console.log("Emoji:", data.emoji);
    },

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
        @react="react"
        @replyToMessage="replyToMessage"
        @showInfoMessage="showInfoMessage"
      />
    </div>
    <Bottombar />
  </div>
</template>
