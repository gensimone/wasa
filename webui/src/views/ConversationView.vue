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
      // Replying to e message basically put the message to reply in the same spot in
      // which we have already put the attachment window. After that the user can choose
      // to delete the "message attachment", or to write a comment and then send the
      // message. The commented message must appear above the comment with a snippet
      // (if it as an image, that image must be also present). Clicking on the comment
      // snippet must change the MessageList position so that the commented message is
      // clearly visible.
      console.log("Reply to message:", message);
    },

    forwardMessage(message) {
      // Forwaring a message must involve the selection of a user/gruop from a picker.
      // The user can select multiple entries from the picker, in which case the message
      // is forwarded to all them. If the user, while on the picker, decide that is
      // doesn't want to forward the message anymore, than we must go back to the previous
      // route (which, in this case is of course that conversation in which the user
      // decided to forward the message).
      // NOTE: We also need to implemente some kind of indicator on the message that
      // clearly shows that the message was forwarded.
      console.log("Forward message:", message);
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

    deleteMessage(message) {
      // Removing a message involves two things:
      // 1) Call API to delete the message.
      // 2) Remove the message from the userConversations or from the
      //    groupConversations (depending from the value of direct).
      // 3) Optional: Removing the message must trigger an animation.
      //    NOTE: We could implement this with a Transition (enter and leave).
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
