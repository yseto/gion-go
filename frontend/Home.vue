<template>
  <div />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import { openapiFetchClient } from "./UserAgent";
export default defineComponent({
  setup: () => {
    const router = useRouter();
    openapiFetchClient.POST("/api/profile").then(data => {
      const next = data.data?.onLoginSkipPinList === true ? "/entry" : "/pin";
      router.push(next);
    });
  },
});
</script>
