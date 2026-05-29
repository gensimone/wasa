import { watch } from "vue";
import { directMessages, groupMessages } from "@/state/conversations";
import { getNotifier } from "@/notifier";
import { users, userId } from "@/state/users";
import router from "@/router";

let prevDirectMessages = new Map();
let prevGroupMessages = new Map();

export function startMessageNotifier() {
  watch(
    directMessages,
    (newMap) => {
      handleNewMessages(newMap, prevDirectMessages, true);
    },
    { deep: true },
  );

  watch(
    groupMessages,
    (newMap) => {
      handleNewMessages(newMap, prevGroupMessages, false);
    },
    { deep: true },
  );
}

function matchConversation(id, direct) {
  const route = router.currentRoute.value;

  const isDirect = route.query.direct === "true";
  const inConversation = route.name === "conversation";
  const idConversation = route.params.id;

  return inConversation && isDirect === direct && idConversation == id;
}

function handleNewMessages(newMap, oldMap, direct) {
  for (const [id, messages] of newMap.entries()) {
    if (matchConversation(id, direct)) continue;

    const oldMessages = oldMap.get(id);

    if (!oldMessages) continue;

    const oldMessageIds = (oldMessages.messages || []).map((m) => m.messageId);

    const newMessages = (messages || []).filter(
      (m) => !oldMessageIds.includes(m.messageId),
    );

    newMessages.forEach((message) => {
      if (message.senderId !== userId.value) {
        const thumbnailUrl = users.value.get(message.senderId).photoUrl;
        getNotifier()?.message({
          text: message.text,
          id: id,
          isGroup: !direct,
          thumbnailUrl: thumbnailUrl,
          attachmentUrl: message.attachmentUrl,
        });
      }
    });
  }

  if (direct) {
    prevDirectMessages = new Map(newMap);
  } else {
    prevGroupMessages = new Map(newMap);
  }
}
