import LoginView from "@/views/LoginView.vue"
import HomeView from "@/views/HomeView.vue"
import UserSettingsView from "@/views/UserSettingsView.vue"
import AddConversationView from "@/views/AddConversationView.vue"
import ConversationView from "@/views/ConversationView.vue"
import GroupCreationView from "@/views/GroupCreationView.vue"
import GroupSettingsView from "@/views/GroupSettingsView.vue"
import GroupAddMemberView from "@/views/GroupAddMemberView.vue"

import { createRouter, createWebHashHistory } from "vue-router"
import { stopPollingUser } from "@/state/user"
import { stopPollingConversations } from "@/state/conversations"

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        { path: "/", component: LoginView },
        { path: "/home", component: HomeView, meta: { requiresAuth: true } },
        { path: "/settings", component: UserSettingsView, meta: { requiresAuth: true } },
        { path: "/conversation/add", component: AddConversationView, meta: { requiresAuth: true } },
        { path: "/conversation/:conversationId", component: ConversationView, meta: { requiresAuth: true } },
        { path: "/group/create", component: GroupCreationView, meta: { requiresAuth: true } },
        { path: "/group/settings", component: GroupSettingsView, meta: { requiresAuth: true } },
        { path: "/group/add", component: GroupAddMemberView, meta: { requiresAuth: true } },
        { path: "/:pathMatch(.*)*", redirect: "/" }
    ]
})

router.beforeEach((to, _, next) => {
    const name = localStorage.getItem("name");

    if (to.path === "/" && name) {
        next("/home");
        return;
    }

    if (to.meta.requiresAuth && !name) {
        stopPollingUser()
        stopPollingConversations()
        next("/");

    } else {
        if (to.path === "/") {
            stopPollingUser()
            stopPollingConversations()
        }

        next();
    }
});

export default router
