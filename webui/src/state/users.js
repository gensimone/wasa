import { ref } from "vue";
import { getUsers } from "@/services/users";
import { defaultUserPhotoUrl } from "@/assets/default";
import Poller from "@/services/poller";

export const users = ref(new Map());

let poller = null;

async function fetchUsers() {
  const fetchedUsers = await getUsers();
  users.value = new Map(
    fetchedUsers.map((u) => [
      u.userId,
      {
        ...u,
        photoUrl: u.photoUrl || defaultUserPhotoUrl,
      },
    ]),
  );
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

export function clearUsers() {
  users.value.clear();
}
