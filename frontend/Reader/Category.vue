<template>
  <div class="card categories__list">
    <div class="card-header">Categories</div>
    <div class="list-group">
      <a
        v-for="(item, index) in categoryStore.list"
        :key="index"
        class="list-group-item"
        :class="{ active: index == categoryStore.selected }"
        @click.prevent="
          categoryStore.setIndex(index);
          $emit('content-update');
        "
        >{{ item.name }}
        <span class="badge badge-light d-sm-none d-md-inline">{{
          item.count
        }}</span></a
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useCategoryStore } from "./Store";

export default defineComponent({
  emits: ["content-update"],
  setup: () => {
    const categoryStore = useCategoryStore();
    return {
      categoryStore,
    };
  },
});
</script>
