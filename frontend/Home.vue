<template>
  <div />
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue"
import { useRouter } from "vue-router"
import { openapiFetchClient } from "./api"
export default defineComponent({
  setup: () => {
    const router = useRouter()
    onMounted(async () => {
      const { data } = await openapiFetchClient.GET("/api/profile")
      const next = data?.onLoginSkipPinList === true ? "/entry" : "/pin"
      router.push(next)
    })
  },
})
</script>
