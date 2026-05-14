import { createRouter, createWebHashHistory } from 'vue-router'
import Login    from '../views/Login.vue'
import Home     from '../views/Home.vue'
import Settings from '../views/Settings.vue'

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        {path: '/',         component: Login},
        {path: '/Home',     component: Home,     meta: {requiresAuth: true}},
        {path: '/settings', component: Settings, meta: {requiresAuth: true}},
    ]
})

router.beforeEach((to, from, next) => {
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
