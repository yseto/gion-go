<template>
  <div class="d-lg-none d-md-none clearfix">
    <p class="float-right">
      <a class="btn btn-light" @click="returntop">Back to Top</a>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue"
export default defineComponent({
  setup: () => {
    const returntop = () => {
      window.scrollTo(0, 0)
    }
    return {
      returntop
    }
  },
})
</script>
