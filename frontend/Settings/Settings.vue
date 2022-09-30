<template>
  <div class="container">
    <Config />
    <hr />
    <Password />
    <hr />
    <Opml />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

import Config from "./Config.vue";
import Opml from "./Opml.vue";
import Password from "./Password.vue";

export default defineComponent({
  components: {
    Config,
    Opml,
    Password,
  },
});
</script>
