<template>
  <div class="container">
    <div class="card pin__list__page">
      <div class="card-header">
        Pin List <span class="badge badge-info">{{ list.length }}</span>
      </div>
      <ul class="list-group">
        <li v-for="(item, index) in list" :key="index" class="list-group-item">
          <a class="btn btn-sm btn-info" style="cursor: pointer" :data-serial="item.serial" :data-feed_id="item.feed_id"
            :data-index="index" @click="applyRead">
            削除
          </a>
          <span>{{ item.update_at }}</span>
          <a :href="item.url" target="blank">{{ item.title }}</a>
        </li>
      </ul>
    </div>
    <BackToTop />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from "vue"
import { format, parse } from "date-fns"
import { openapiFetchClient } from "./api"
import BackToTop from "./BackToTop.vue"
import { PinList } from "./types"
const localtime = (mysqlDT: string) => {
  return format(
    parse(`${mysqlDT} Z`, "yyyy-MM-dd HH:mm:ss X", new Date()),
    "yyyy-MM-dd HH:mm"
  )
}
export default defineComponent({
  components: {
    BackToTop,
  },
  setup: () => {
    const list = ref<PinList[]>([])
    const applyRead = async (event: Event) => {
      const target = event.target as HTMLElement
      const serial = target.getAttribute("data-serial")
      if (serial === null) {
        return
      }
      const feed_id = target.getAttribute("data-feed_id")
      if (feed_id === null) {
        return
      }
      await openapiFetchClient.POST("/api/pin", {
        body: {
          readflag: "Setpin",
          serial: parseInt(serial, 10),
          feed_id: parseInt(feed_id, 10),
        }
      })
      const idx = target.getAttribute("data-index")
      if (idx) {
        list.value.splice(parseInt(idx, 10), 1)
      }
    }

    onMounted(async () => {
      const { data } = await openapiFetchClient.GET("/api/pin")
      if (data === undefined) {
        return
      }
      list.value = data.map((x) => {
        return { ...x, update_at: localtime(x.update_at) }
      })
    })

    return {
      list,
      applyRead,
    }
  },
})
</script>
