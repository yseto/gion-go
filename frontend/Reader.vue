<template>
  <div class="container">
    <div class="row">
      <div class="col-sm-3">
        <PinList />
        <br />
        <CategoryList @content-update="contentUpdate" />
        <br />
      </div>
      <div class="col-sm-9">
        <div>
          <div v-if="contentStore.list.length === 0" class="tw card well">
            <h5 class="text-center">No unreading entries.</h5>
          </div>
          <div
            v-if="contentStore.list.length > 0 && autoSeen === false"
            class="sticky-top bg-white pt-1 pb-1 text-right"
          >
            <a class="btn btn-sm btn-dark" @click.prevent="contentReadIt"
              >Mark as read</a
            >
          </div>
          <div v-for="(item, index) in contentStore.list" :key="index">
            <div
              class="tw card"
              :class="{
                'tw--active border-info': index == contentStore.selected,
                'tw--pinned': item.readflag == 'Setpin',
              }"
            >
              <h5 class="viewpage">
                <a
                  :href="item.url"
                  target="blank"
                  rel="noreferrer"
                  class="text-dark"
                >
                  <span v-if="item.title.length > 0">{{ item.title }}</span>
                  <span v-else>[nothing title...]</span>
                </a>
              </h5>
              <p>{{ item.description }}</p>
              <div class="clearfix">
                <span class="float-left"
                  >{{ item.date }} - {{ item.site_title }}</span
                >
                <span class="float-right d-inline d-md-inline">
                  <span v-if="item.readflag == 'Seen'"> &#x2714; </span>
                  <span
                    :class="{
                      pinned: item.readflag == 'Setpin',
                      unpinned: item.readflag != 'Setpin',
                    }"
                    @click="togglePin(index)"
                    >&#x1f4cc;</span
                  >
                </span>
              </div>
              <div class="d-md-none d-lg-none">
                <br />
                <button
                  class="btn btn-info btn-sm btn-block"
                  @click="togglePin(index)"
                >
                  Pin!
                </button>
              </div>
            </div>
          </div>
        </div>
        <br />
      </div>
    </div>
    <CategoryPager @content-update="contentUpdate" />
    <br />
    <BackToTop />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, onUnmounted, onUpdated, ref } from "vue";
import { format, fromUnixTime } from "date-fns";
import BackToTop from "./BackToTop.vue";
import PinList from "./Reader/PinList.vue";
import { openapiFetchClient } from "./UserAgent";
import {
  useCategoryStore,
  Category,
  Entry,
  useContentStore,
  ReadFlag,
} from "./Reader/Store";
import { useUserStore } from "./UserStore";
import CategoryPager from "./Reader/CategoryPager.vue";
import CategoryList from "./Reader/Category.vue";

const epochToDateTime = (epoch: number) => {
  return format(fromUnixTime(epoch), "MM/dd HH:mm");
};

export default defineComponent({
  components: {
    BackToTop,
    PinList,
    CategoryPager,
    CategoryList,
  },
  setup: () => {
    const userStore = useUserStore();
    const contentStore = useContentStore();
    const categoryStore = useCategoryStore();

    const readItTimeoutID = ref(0);
    const autoSeen = ref(false);
    const doScroll = ref(false);

    onUnmounted(() => {
      document.removeEventListener("keypress", keypressHandler);
    });
    onMounted(() => {
      document.addEventListener("keypress", keypressHandler);
    });

    const contentNext = () => {
      contentStore.nextContent();
      doScroll.value = true;
    };
    const contentPrevious = () => {
      contentStore.previousContent();
      doScroll.value = true;
    };

    const contentUpdate = () => {
      doScroll.value = true;

      if (categoryStore.list.length == 0) {
        contentStore.setContents({ list: [], index: 0 });
        return;
      }

      openapiFetchClient.POST("/api/unread_entry", {
        body: {
          category: categoryStore.list[categoryStore.selected].id,
        }
      }).then(data => {
        if (data.data === undefined) {
          return
        }
        const list = data.data.map((x) => {
          return { ...x, date: epochToDateTime(x.date_epoch) };
        });
        contentStore.setContents({ list: list, index: 0 });
        categoryUpdate();
      }).then(() => {
        if (autoSeen.value) {
          readIt(500);
        }
      });
    };

    const categoryUpdate = () => {
      openapiFetchClient.POST("/api/category_with_count").then((data) => {
        if (data.data === undefined) {
          return
        }
        const list = data.data
        let updated = false;
        list.forEach(function (_, index) {
          // 現在選択している category_id が一致するものがある時
          // indexを任意の位置に移動する
          if (categoryStore.currentPointer() === list[index].id) {
            categoryStore.setCategories({ list: list, index: index });
            updated = true;
            return;
          }
        });

        // 表示していたカテゴリをすべて読み終えた場合などは、
        // カテゴリ一覧の一番上のカテゴリにindexを設定する
        // 選択可能なカテゴリがない場合は実行しない
        if (!updated) {
          categoryStore.setCategories({ list: list, index: 0 });
          contentUpdate();
        }
      });
    };

    const readIt = (sendDelay: number) => {
      // ダブルタップやキーボードの連続押下の場合にイベントをキャンセルする
      if (readItTimeoutID.value) {
        clearTimeout(readItTimeoutID.value);
      }

      // 未読ステータスのものだけ送るため、フィードのアイテム既読リストを作成をする
      const params = contentStore.list
        .filter((item) => item.readflag === "Unseen")
        .map((item) => {
          return { serial: item.serial, feed_id: item.feed_id };
        });

      if (params.length === 0) {
        return;
      }
      readItTimeoutID.value = window.setTimeout(function () {
        readItTimeoutID.value = 0;

        openapiFetchClient.POST("/api/set_asread", {
          body: params,
        }).then(() => {
          params.map((e) =>
            contentStore.setSeen({
              feed_id: e.feed_id,
              serial: e.serial,
            })
          );
        });
      }, sendDelay);
    };
    const contentReadIt = () => {
      readIt(0);
    };
    const itemView = () => {
      const url = contentStore.currentEntryURL();
      if (url) {
        window.open(url);
      }
    };

    const togglePin = async (index: number | undefined) => {
      // for touch device.
      if (typeof index !== "undefined") {
        contentStore.setIndex(index);
      }

      openapiFetchClient.POST("/api/set_pin", {
        body: contentStore.currentEntrySerialData(),
      }).then(data => {
        if (data.data === undefined) {
          return
        }
        contentStore.setReadflag(data.data.readflag);
      });
    };

    const keypressHandler = function (e: KeyboardEvent) {
      e.preventDefault();
      switch (e.code) {
        case "KeyI":
          contentReadIt();
          break;
        case "KeyP":
          togglePin(contentStore.selected);
          break;
        case "KeyR":
          contentUpdate();
          break;
        case "KeyK":
          contentPrevious();
          break;
        case "KeyJ":
          contentNext();
          break;
        case "KeyV":
          itemView();
          break;
      }
    };

    // https://stackoverflow.com/a/36673184
    const isMobile = ('ontouchstart' in document.documentElement);

    onUpdated(() => {
      const element = document.querySelectorAll(".tw--active");
      if (doScroll.value && element.length === 1) {
        const rect = element[0].getBoundingClientRect();
        if (!isMobile && contentStore.selected === 0) {
          window.scrollTo(0, 0);
        } else {
          const positionY = rect.top + window.pageYOffset - 40; // offset: 40
          window.scrollTo(0, positionY);
        }
      }
      doScroll.value = false;
    });

    autoSeen.value = userStore.user.autoSeen;
    openapiFetchClient.POST("/api/category_with_count").then((data) => {
      if (data.data === undefined) {
        return
      }
      categoryStore.setCategories({ list: data.data, index: 0 });
      contentUpdate();
    });

    return {
      contentStore,

      autoSeen,

      contentReadIt,
      contentUpdate,
      togglePin,
    };
  },
});
</script>
