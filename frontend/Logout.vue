<template>
  <div />
</template>
<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import { Agent } from "./UserAgent";
import { useUserStore } from "./UserStore";

export default defineComponent({
  setup: () => {
    const router = useRouter();
    const store = useUserStore();
    Agent({ url: "/api/logout" }).then(() => {
      store.Logout();
      router.push("/login");
    });
  },
});
</script>
