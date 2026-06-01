import Poller from "@/services/poller";
import { ref } from "vue";
import { getMessage } from "@/services/messages";
import { getGroup } from "@/services/groups";
import { getMyConversations, getConversation } from "@/services/conversations";
import { defaultGroupPhotoUrl } from "@/assets/default";

export const directMessages = ref(new Map());
export const groupMessages = ref(new Map());
export const groups = ref(new Map());

let poller = null;

const CACHE_TTL = 30_000;
export const messageCache = new Map();

async function fetchMessageCached(id) {
  const cached = messageCache.get(id);

  if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
    return cached.data;
  }

  const msg = await getMessage(id);

  messageCache.set(id, {
    data: msg,
    timestamp: Date.now(),
  });

  return msg;
}

export async function getConversationMessages(id, direct) {
  const messageIds = await getConversation(id, direct);

  if (!messageIds?.length) return [];

  const messages = await Promise.all(messageIds.map(fetchMessageCached));

  return messages;
}

async function fetchGroup(id) {
  const group = await getGroup(id);
  return [
    id,
    {
      ...group,
      photoUrl: group.photoUrl || defaultGroupPhotoUrl,
    },
  ];
}

async function fetchGroupMessages(id) {
  const messages = await getConversationMessages(id, false);
  return [id, messages];
}

async function fetchDirectMessages(id) {
  const messages = await getConversationMessages(id, true);
  return [id, messages];
}

export async function fetchConversations() {
  const fetched = (await getMyConversations()) || [];

  const groupIds = fetched.filter((c) => !c.isDirect).map((c) => c.id);
  const directIds = fetched.filter((c) => c.isDirect).map((c) => c.id);

  const [fetchedGroups, fetchedGroupMessages, fetchedDirectMessages] =
    await Promise.all([
      Promise.all(groupIds.map(fetchGroup)),
      Promise.all(groupIds.map(fetchGroupMessages)),
      Promise.all(directIds.map(fetchDirectMessages)),
    ]);

  groups.value = new Map(fetchedGroups);
  groupMessages.value = new Map(fetchedGroupMessages);
  directMessages.value = new Map(fetchedDirectMessages);
}

export function startPollingConversations(interval = 5000) {
  stopPollingConversations();

  poller = new Poller(async () => {
    await fetchConversations();
  }, interval);

  poller.startPolling();
}

export function stopPollingConversations() {
  if (poller) {
    poller.stopPolling();
    poller = null;
  }
}

export function clearMessages() {
  directMessages.value.clear();
  groupMessages.value.clear();
  messageCache.clear();
}

export function clearGroups() {
  groups.value.clear();
}
