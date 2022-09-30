<template>
  <div style="max-width: 330px; padding: 80px 15px 0; margin: 0 auto">
    <h3>Gion</h3>
    <input
      ref="focus"
      v-model="creds.id"
      type="text"
      class="form-control"
      placeholder="ID"
      required
      @keydown.enter="login"
    />
    <input
      v-model="creds.password"
      type="password"
      class="form-control"
      placeholder="Password"
      required
      @keydown.enter="login"
    />
    <button class="btn btn-primary" style="margin-top: 20px" @click="login">
      Sign in
    </button>
  </div>
</template>
<script lang="ts">
import { defineComponent, reactive, ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Agent } from "./UserAgent";
import { useUserStore } from "./UserStore";
class Credentials {
  id = "";
  password = "";
}
export default defineComponent({
  setup: () => {
    const router = useRouter();
    const route = useRoute();
    const store = useUserStore();
    const creds = reactive(new Credentials());
    const focus = ref<HTMLInputElement | null>();
    onMounted(() => {
      if (focus.value) {
        focus.value.focus();
      }
    });
    const login = () => {
      Agent<{ autoseen: boolean; token: string }>({
        url: "/api/login",
        data: {
          id: creds.id,
          password: creds.password,
        },
      }).then((payload) => {
        store.Login({
          autoSeen: payload.autoseen,
          token: payload.token,
        });

        if (route.query.redirect) {
          router.push(route.query.redirect.toString());
        } else {
          router.push("/");
        }
      });
    };
    return {
      creds,
      login,
      focus,
    };
  },
});
</script>
