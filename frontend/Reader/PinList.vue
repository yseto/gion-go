<template>
  <div>
    <p class="pt-1 pb-1">
      <a class="btn btn-sm btn-info" @click.prevent="togglePinList">Pin List</a>
    </p>
    <div v-if="visibleState" class="card pin__list">
      <div class="card-header">
        Pin List
        <span class="badge badge-info">{{ list.length }}</span>
      </div>
      <div class="list-group">
        <a v-for="(item, index) in list" :key="index" class="list-group-item" :href="item.url">{{ item.title }}</a>
      </div>
      <div class="card-footer text-center">
        <a class="btn btn-sm btn-outline-dark" :class="{ disabled: list.length == 0 }"
          @click.prevent="pinlistClean">Remove All Pin</a>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted, ref } from "vue"
import { openapiFetchClient } from "../UserAgent"
import { PinList } from "../types"
export default defineComponent({
  setup: () => {
    const visibleState = ref(false)
    const list = ref<PinList[]>([])

    const togglePinList = async () => {
      visibleState.value = visibleState.value ? false : true
      if (visibleState.value) {
        const { data } = await openapiFetchClient.GET("/api/pin")
        if (data === undefined) {
          return
        }
        list.value = data
      }
    }

    const pinlistClean = async () => {
      if (!confirm("ピンをすべて外しますか?")) {
        return
      }
      visibleState.value = false
      await openapiFetchClient.DELETE("/api/pin")
      list.value = []
    }

    onUnmounted(() => {
      document.removeEventListener("keypress", keypressHandler)
    })
    onMounted(() => {
      document.addEventListener("keypress", keypressHandler)
    })

    const keypressHandler = function (e: KeyboardEvent) {
      e.preventDefault()
      switch (e.code) {
        case "KeyO":
          togglePinList()
          break
      }
    }
    return {
      pinlistClean,
      togglePinList,
      visibleState,
      list,
    }
  },
})
</script>
