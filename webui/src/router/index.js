import { createRouter, createWebHashHistory } from "vue-router"
import { user, startPollingUser } from "@/state/user"
import LoginView from "@/views/LoginView.vue"
import HomeView from "@/views/HomeView.vue"
import SettingsView from "@/views/SettingsView.vue"
import AddConversationView from "@/views/AddConversationView.vue"
import ConversationView from "@/views/ConversationView.vue"
import CreateGroupView from "@/views/CreateGroupView.vue"
import InfoGroupView from "@/views/InfoGroupView.vue"

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        { path: "/", component: LoginView },
        { path: "/home", component: HomeView, meta: { requiresAuth: true } },
        { path: "/settings", component: SettingsView, meta: { requiresAuth: true } },
        { path: "/conversation/add", component: AddConversationView, meta: { requiresAuth: true } },
        { path: "/conversation", component: ConversationView, meta: { requiresAuth: true } },
        { path: "/group/create", component: CreateGroupView, meta: { requiresAuth: true } },
        { path: "/group/info", component: InfoGroupView, meta: { requiresAuth: true } },
        { path: "/:pathMatch(.*)*", redirect: "/" }
    ]
})

router.beforeEach((to, _, next) => {
    const username = localStorage.getItem("name");

    if (to.path === "/" && username) {
        startPollingUser()
        next("/home");
        return;
    }

    if (to.meta.requiresAuth && !username) {
        user.poller?.stopPolling()
        next("/");
    } else {
        if (to.path !== "/") startPollingUser()
        next();
    }
});

export default router
