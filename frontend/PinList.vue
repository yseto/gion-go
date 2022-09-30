<template>
  <div class="container">
    <div class="card pin__list__page">
      <div class="card-header">
        Pin List <span class="badge badge-info">{{ list.length }}</span>
      </div>
      <ul class="list-group">
        <li v-for="(item, index) in list" :key="index" class="list-group-item">
          <a
            class="btn btn-sm btn-info"
            style="cursor: pointer"
            :data-serial="item.serial"
            :data-feed_id="item.feed_id"
            :data-index="index"
            @click="applyRead"
          >
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
import { defineComponent, ref } from "vue";
import { format, parse } from "date-fns";
import { Agent } from "./UserAgent";
import BackToTop from "./BackToTop.vue";
import { PinList } from "./types";
const localtime = (mysqlDT: string) => {
  return format(
    parse(`${mysqlDT} Z`, "yyyy-MM-dd HH:mm:ss X", new Date()),
    "yyyy-MM-dd HH:mm"
  );
};
export default defineComponent({
  components: {
    BackToTop,
  },
  setup: () => {
    const list = ref<PinList[]>([]);
    const applyRead = (event: Event) => {
      const target = event.target as HTMLElement;
      Agent({
        url: "/api/set_pin",
        data: {
          readflag: 2,
          serial: target.getAttribute("data-serial"),
          feed_id: target.getAttribute("data-feed_id"),
        },
      }).then(() => {
        const idx = target.getAttribute("data-index");
        if (idx) {
          list.value.splice(parseInt(idx, 10), 1);
        }
      });
    };

    Agent<PinList[]>({ url: "/api/pinned_items" }).then((data) => {
      list.value = data.map((x) => {
        return { ...x, update_at: localtime(x.update_at) };
      });
    });

    return {
      list,
      applyRead,
    };
  },
});
</script>
