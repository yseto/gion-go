import { createApp } from "vue"
import { createRouter, createWebHashHistory } from "vue-router"
import { createPinia } from "pinia"
import { useUserStore } from "./UserStore"
import GionHeader from "./Header.vue"
import Home from "./Home.vue"
import Login from "./Login.vue"
import Logout from "./Logout.vue"
import PinList from "./PinList.vue"
import addSubscription from "./addSubscription.vue"
import Reader from "./Reader.vue"
import Settings from "./Settings/Settings.vue"
import manageSubscription from "./manageSubscription.vue"
import NotFound from "./NotFound.vue"

const app = createApp({
  components: {
    GionHeader,
  },
})

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", component: Home, meta: { requiresAuth: true } },
    { path: "/pin", component: PinList, meta: { requiresAuth: true } },
    { path: "/add", component: addSubscription, meta: { requiresAuth: true } },
    { path: "/entry", component: Reader, meta: { requiresAuth: true } },
    { path: "/settings", component: Settings, meta: { requiresAuth: true } },
    {
      path: "/subscription",
      component: manageSubscription,
      meta: { requiresAuth: true },
    },

    { path: "/login", component: Login, meta: { anonymous: true } },
    { path: "/logout", component: Logout, meta: { anonymous: true } },
    {
      path: "/:pathMatch(.*)*",
      name: "not-found",
      component: NotFound,
      meta: { anonymous: true },
    },
  ],
})

router.beforeEach((to, from, next) => {
  const store = useUserStore()
  if (
    to.matched.some((record) => record.meta.requiresAuth) &&
    !store.isLogin()
  ) {
    next({ path: "/login", query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

const pinia = createPinia()

app.use(router).use(pinia).mount("#app")
