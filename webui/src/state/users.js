import { ref } from "vue";
import { getUsers } from "@/services/users";
import Poller from "@/services/poller";

export const users = ref(new Map());
export const userId = ref(null);

let poller = null;

async function fetchUsers() {
  const fetchedUsers = await getUsers();
  users.value = new Map(fetchedUsers.map((u) => [u.userId, u]));
}

export function startPollingUsers(interval = 5000) {
  stopPollingUsers();

  poller = new Poller(async () => {
    await fetchUsers();
  }, interval);

  poller.startPolling();
}

export function stopPollingUsers() {
  if (poller) {
    poller.stopPolling();
    poller = null;
  }
}

export function setUserId(id) {
  userId.value = Number(id);
  localStorage.setItem("userId", id);
}

export function hydrateUserId() {
  const rawUserId = localStorage.getItem("userId");
  if (rawUserId) setUserId(rawUserId);
}

export function clearUsers() {
  users.value.clear();
  userId.value = null;
  localStorage.removeItem("userId");
}
