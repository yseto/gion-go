<template>
  <div />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useRouter } from "vue-router";
import { Agent } from "./UserAgent";
import { Profile } from "./types";
export default defineComponent({
  setup: () => {
    const router = useRouter();
    Agent<Profile>({ url: "/api/profile" }).then((data) => {
      const next = data.onLoginSkipPinList === true ? "/entry" : "/pin";
      router.push(next);
    });
  },
});
</script>
