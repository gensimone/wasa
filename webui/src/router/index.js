import LoginView from "@/views/LoginView.vue";
import HomeView from "@/views/HomeView.vue";
import UserSettingsView from "@/views/UserSettingsView.vue";
import ConversationView from "@/views/ConversationView.vue";
import AddConversationView from "@/views/AddConversationView.vue";
import GroupCreationView from "../views/GroupCreationView.vue";
import GroupSettingsView from "../views/GroupSettingsView.vue";
import { isValidGroupRoute, isValidConversationRoute } from "@/router/validate";
import { createRouter, createWebHashHistory } from "vue-router";
import { stopPollingUser } from "@/state/user";
import {
  stopPollingConversations,
  clearConversations,
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
  ],
});

router.beforeEach(async (to, _, next) => {
  const name = localStorage.getItem("name");

  const isLoginRoute = to.path === "/";
  const requiresAuth = to.meta.requiresAuth;

  if (isLoginRoute && name) {
    next("/home");
    return;
  }

  if (requiresAuth && !name) {
    stopPollingUser();
    stopPollingConversations();
    clearConversations();
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

  if (isLoginRoute) {
    stopPollingUser();
    stopPollingConversations();
    clearConversations();
  }

  next();
});

export default router;
