import { createApp } from "vue";
import App from "@/App.vue";
import router from "@/router";
import notifier from "@/notifier";
import { startMessageNotifier } from "@/notifier/messageNotifier";
import { startPollingUsers } from "@/state/users.js";
import { startPollingConversations } from "@/state/conversations.js";
import { user, hydrateUserState, startPollingUser } from "@/state/user";
import "./styles/index.css";

hydrateUserState();

if (user.userId) {
  startPollingUser();
  startPollingUsers();
  startPollingConversations();

  startMessageNotifier();
}

const app = createApp(App);
app.use(router);
app.use(notifier);
app.mount("#app");
