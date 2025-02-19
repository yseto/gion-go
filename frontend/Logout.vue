<template>
  <div />
</template>
<script lang="ts">
import { defineComponent, onMounted } from "vue"
import { useRouter } from "vue-router"
import { openapiFetchClient } from "./api"
import { useUserStore } from "./UserStore"

export default defineComponent({
  setup: () => {
    const router = useRouter()
    const store = useUserStore()
    onMounted(async ()=>{
      await openapiFetchClient.POST("/api/logout")
      store.Logout()
      router.push("/login")
    })
  },
})
</script>
