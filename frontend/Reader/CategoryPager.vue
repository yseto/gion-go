<template>
  <div class="d-lg-none d-md-none clearfix">
    <div class="float-left">
      <a class="btn btn-dark btn-sm" @click.prevent="categoryPrevious"
        >&lt;&lt; Category</a
      >
    </div>
    <div class="float-right">
      <a class="btn btn-dark btn-sm" @click.prevent="categoryNext"
        >Category &gt;&gt;</a
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted } from "vue";
import { useCategoryStore } from "./Store";
export default defineComponent({
  emits: ["content-update"],
  setup: (props, context) => {
    const categoryStore = useCategoryStore();

    const categoryNext = () => {
      categoryStore.nextCategory();
      context.emit("content-update");
    };
    const categoryPrevious = () => {
      categoryStore.previousCategory();
      context.emit("content-update");
    };

    const keypressHandler = function (e: KeyboardEvent) {
      e.preventDefault();
      switch (e.code) {
        case "KeyA":
          categoryPrevious();
          break;
        case "KeyS":
          categoryNext();
          break;
      }
    };
    onUnmounted(() => {
      document.removeEventListener("keypress", keypressHandler);
    });
    onMounted(() => {
      document.addEventListener("keypress", keypressHandler);
    });

    return {
      categoryNext,
      categoryPrevious,
    };
  },
});
</script>
