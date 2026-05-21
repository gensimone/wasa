import { createRouter, createWebHashHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import HomeView from '@/views/HomeView.vue'
import SettingsView from '@/views/SettingsView.vue'
import AddConversationView from '@/views/AddConversationView.vue'
import ConversationView from '@/views/ConversationView.vue'

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', component: LoginView },
        { path: '/home', component: HomeView, meta: { requiresAuth: true } },
        { path: '/settings', component: SettingsView, meta: { requiresAuth: true } },
        { path: '/add', component: AddConversationView, meta: { requiresAuth: true } },
        { path: '/conversation', component: ConversationView, meta: { requiresAuth: true } },
    ]
})

router.beforeEach((to, _, next) => {
    const username = localStorage.getItem("name");

    if (to.path === "/" && username) {
        next("/home");
        return;
    }

    if (to.meta.requiresAuth && !username) {
        next("/");
    } else {
        next();
    }
});

export default router
