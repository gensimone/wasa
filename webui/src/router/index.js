import { createRouter, createWebHashHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import Settings from '../views/Settings.vue'
import Users from '../views/Users.vue'
import GroupCreate from '../views/GroupCreate.vue'

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        { path: '/', component: Login },
        { path: '/Home', component: Home, meta: { requiresAuth: true } },
        { path: '/settings', component: Settings, meta: { requiresAuth: true } },
        { path: '/users', component: Users, meta: { requiresAuth: true } },
        { path: '/create/group', component: GroupCreate, meta: { requiresAuth: true } },
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
