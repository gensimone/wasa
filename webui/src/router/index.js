import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "@/views/LoginView.vue";
import HomeView from "@/views/HomeView.vue";
import UserSettingsView from "@/views/UserSettingsView.vue";
import ConversationView from "@/views/ConversationView.vue";
import AddConversationView from "@/views/AddConversationView.vue";
import GroupCreationView from "../views/GroupCreationView.vue";
import GroupSettingsView from "../views/GroupSettingsView.vue";
import ForwardMessageToUsersView from "@/views/ForwardMessageToUsersView.vue";
import { stopMessageNotifier } from "@/notifier/messageNotifier";
import {
  isValidMessageRoute,
  isValidGroupRoute,
  isValidConversationRoute,
} from "@/router/validate";
import { stopPollingUsers, clearUsers } from "@/state/users";
import {
  stopPollingConversations,
  clearMessages,
  clearGroups,
} from "@/state/conversations";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: "/:pathMatch(.*)*", redirect: "/" },
    { path: "/", component: LoginView },
    { path: "/home", component: HomeView, meta: { requiresAuth: true } },
    {
      path: "/user/settings",
      component: UserSettingsView,
      meta: { requiresAuth: true },
    },
    {
      path: "/group/create",
      component: GroupCreationView,
      meta: { requiresAuth: true },
    },
    {
      path: "/conversation/add",
      component: AddConversationView,
      meta: { requiresAuth: true },
    },
    {
      path: "/group/:id",
      name: "group",
      component: GroupSettingsView,
      meta: { requiresAuth: true },
    },
    {
      path: "/conversation/:id",
      name: "conversation",
      component: ConversationView,
      meta: { requiresAuth: true },
    },
    {
      path: "/message/:id/forward",
      name: "message-forward",
      component: ForwardMessageToUsersView,
      meta: { requiresAuth: true },
    },
  ],
});

function cleanup() {
  stopPollingUsers();
  stopPollingConversations();
  stopMessageNotifier();

  clearUsers();
  clearMessages();
  clearGroups();
}

router.beforeEach(async (to, _, next) => {
  const userId = localStorage.getItem("userId");

  const isLoginRoute = to.path === "/";
  const requiresAuth = to.meta.requiresAuth;

  if (isLoginRoute && userId) {
    next("/home");
    return;
  }

  if (requiresAuth && !userId) {
    cleanup();

    next("/");
    return;
  }

  if (to.name === "conversation") {
    (await isValidConversationRoute(to)) ? next() : next("/home");
    return;
  }

  if (to.name === "group") {
    (await isValidGroupRoute(to)) ? next() : next("/home");
    return;
  }

  if (to.name === "message" || to.name === "message-forward") {
    (await isValidMessageRoute(to)) ? next() : next("/home");
    return;
  }

  if (isLoginRoute) cleanup();

  next();
});

export default router;
