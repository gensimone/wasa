import LoginView from "@/views/LoginView.vue"
import HomeView from "@/views/HomeView.vue"
import UserSettingsView from "@/views/UserSettingsView.vue"
import ConversationView from "@/views/ConversationView.vue"
import AddConversationView from "@/views/AddConversationView.vue"

import { createRouter, createWebHashHistory } from "vue-router"
import { stopPollingUser } from "@/state/user"
import { groupConversations, stopPollingConversations } from "@/state/conversations"
import { getUserById } from "@/services/users"

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        { path: "/:pathMatch(.*)*", redirect: "/" },
        { path: "/", component: LoginView },
        { path: "/home", component: HomeView, meta: { requiresAuth: true } },
        { path: "/user/settings", component: UserSettingsView, meta: { requiresAuth: true } },
        { path: "/conversation/add", component: AddConversationView, meta: { requiresAuth: true } },
        {
            path: "/conversation/:id",
            name: "conversation",
            component: ConversationView,
            meta: { requiresAuth: true },
            beforeEnter: async (to, _, next) => {
                const id = Number(to.params.id)

                const validId =
                    Number.isInteger(id) &&
                    id > 0

                const direct = String(to.query.direct).toLowerCase()

                const validDirect =
                    direct === "true" || direct === "false"

                if (!validId || !validDirect) {
                    next("/home")
                    return
                }

                let isValid = false
                if (direct == "true") {
                    try {
                        const user = await getUserById(id)
                        isValid = user?.userId === id
                    } catch (e) {
                        isValid = false
                    }
                } else {
                    // FIXME: Make sure groupConversations has been fetched
                    // before checking this!!
                    isValid = groupConversations.value.has(id)
                }

                if (!isValid) {
                    next("/home")
                } else {
                    next()
                }
            }
        }
    ]
})

router.beforeEach((to, _, next) => {
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
        next("/");
        return;
    }

    if (isLoginRoute) {
        stopPollingUser();
        stopPollingConversations();
    }

    next();
});

export default router
