<template>
  <div />
</template>
<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import { openapiFetchClient } from "./UserAgent";
import { useUserStore } from "./UserStore";

export default defineComponent({
  setup: () => {
    const router = useRouter();
    const store = useUserStore();
    openapiFetchClient.POST("/api/logout").then(() => {
      store.Logout();
      router.push("/login");
    });
  },
});
</script>
