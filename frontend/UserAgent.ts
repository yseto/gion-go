import { useUserStore } from "./UserStore";
import createClient, { type Middleware } from "openapi-fetch"
import type { paths } from "./api/schema"

const fetchRequestInterceptor: Middleware = {
  async onRequest({ request }) {

    const store = useUserStore();

    if (store.isLogin === null) {
      return request
    }
    if (!store.user.token) {
      return request
    }
    request.headers.set("Authorization", `Bearer ${store.user.token}`)
    return request
  }
}

const fetchResponseInterceptor: Middleware = {
  async onResponse({ response }) {
    const store = useUserStore();
    if (response.status === 401) {
      if (store.isLogin()) {
        store.Logout();
        location.reload();
      }
    }
    return response
  }
}

export const openapiFetchClient = createClient<paths>({})

openapiFetchClient.use(fetchRequestInterceptor)
openapiFetchClient.use(fetchResponseInterceptor)

export type APIClient = typeof openapiFetchClient