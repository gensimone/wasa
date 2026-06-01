import { getNotifier } from "@/notifier";
import { getMessage } from "@/services/messages";
import { getUserById } from "@/services/users";
import { handleError } from "@/utils/errors";
import { groupMessages } from "@/state/conversations";
import { users } from "@/state/users";
import { getMyConversations } from "@/services/conversations";

export async function isValidMessageId(id) {
  let isValid = false;

  try {
    const message = await getMessage(id);
    isValid = message?.messageId === id;
  } catch (e) {
    if (e.response?.status === 404) {
      getNotifier().error("Message not found");
    } else {
      handleError(e);
    }
  }

  return isValid;
}

export async function isValidUserId(id) {
  let isValid = users.value.has(id);

  if (!isValid) {
    // We must make sure the user exists.
    try {
      const user = await getUserById(id);
      isValid = user?.userId === id;
    } catch (e) {
      if (e.response?.status !== 404) {
        handleError(e);
      }

      isValid = false;
    }
  }

  return isValid;
}

export async function isValidConversationId(id) {
  let isValid = groupMessages.value.has(id);

  if (groupMessages.value.size !== 0) return isValid;

  // groupMessages may not have been loaded yet.
  try {
    const fetched = (await getMyConversations()) || [];
    isValid = fetched
      .filter((c) => !c.isDirect)
      .map((c) => c.id)
      .includes(id);
  } catch (e) {
    handleError(e);
    isValid = false;
  } finally {
    return isValid;
  }
}

export async function isValidMessageRoute(route) {
  const id = Number(route.params.id);
  const validId = Number.isInteger(id) && id > 0;

  if (!validId) return false;

  return isValidMessageId(id);
}

export async function isValidGroupRoute(route) {
  const id = Number(route.params.id);
  const validId = Number.isInteger(id) && id > 0;

  return validId ? isValidConversationId(id) : false;
}

export async function isValidConversationRoute(route) {
  const id = Number(route.params.id);
  const validId = Number.isInteger(id) && id > 0;

  const direct = String(route.query.direct).toLowerCase();
  const validDirect = direct === "true" || direct === "false";

  if (!validId || !validDirect) return false;

  return direct === "true" ? isValidUserId(id) : isValidConversationId(id);
}
