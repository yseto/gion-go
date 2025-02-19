<template>
  <div style="max-width: 330px; padding: 80px 15px 0; margin: 0 auto">
    <h3>Gion</h3>
    <input ref="focus" v-model="creds.id" type="text" class="form-control" placeholder="ID" required
      @keydown.enter="login" />
    <input v-model="creds.password" type="password" class="form-control" placeholder="Password" required
      @keydown.enter="login" />
    <button class="btn btn-primary" style="margin: 20px 0" @click="login">
      Sign in
    </button>
    <div v-if="failed" class="alert alert-warning" role="alert">
      failed sign in
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent, reactive, ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { openapiFetchClient } from "./UserAgent"
import { useUserStore } from "./UserStore"
class Credentials {
  id = ""
  password = ""
}
export default defineComponent({
  setup: () => {
    const router = useRouter()
    const route = useRoute()
    const store = useUserStore()
    const creds = reactive(new Credentials())
    const focus = ref<HTMLInputElement | null>()
    const failed = ref<boolean>(false)
    onMounted(() => {
      if (focus.value) {
        focus.value.focus()
      }
    })
    const login = async () => {
      failed.value = false
      const { data, response } = await openapiFetchClient.POST("/api/login", {
        body: {
          id: creds.id,
          password: creds.password,
        },
      })
      if (response.status === 401) {
        failed.value = true
        return
      }
      if (data === undefined) {
        failed.value = true
        return
      }
      store.Login({
        autoSeen: data.autoseen,
        token: data.token,
      })

      if (route.query.redirect) {
        router.push(route.query.redirect.toString())
      } else {
        router.push("/")
      }
    }
    return {
      creds,
      login,
      focus,
      failed,
    }
  },
})
</script>
