<template>
  <div class="container">
    <div class="row">
      <div class="col-md-8">
        <h4>Subscription</h4>
        <div class="form-horizontal">
          <div class="row form-group">
            <label class="col-sm-3 col-form-label" for="inputURL"
              >URL(Web Page)</label
            >
            <div class="col-sm-6">
              <input
                v-model="site.url"
                type="text"
                placeholder="URL"
                class="form-control"
                @blur="feedDetail"
              />
            </div>
            <div class="col-sm-3">
              <a class="btn btn-info" @click.prevent="feedDetail">Get Detail</a>
              <div
                v-if="searchState"
                class="spinner-border spinner-border-sm"
                role="status"
              >
                <span class="sr-only">Loading...</span>
              </div>
            </div>
          </div>
          <div class="row form-group">
            <label class="col-sm-3 col-form-label" for="inputTitle"
              >Title</label
            >
            <div class="col-sm-6">
              <input
                v-model="site.title"
                type="text"
                placeholder="Title"
                class="form-control"
              />
            </div>
          </div>
          <div class="row form-group">
            <label class="col-sm-3 col-form-label" for="inputRSS"
              >URL(Subscription)</label
            >
            <div class="col-sm-6">
              <input
                v-model="site.rss"
                type="text"
                placeholder="RSS"
                class="form-control"
              />
            </div>
          </div>
          <div class="row form-group">
            <label class="col-sm-3 col-form-label" for="selectCat"
              >Categories</label
            >
            <div class="col-sm-6">
              <select
                v-model="category"
                class="form-control"
                placeholder="Choose Category"
              >
                <option
                  v-for="item in categories"
                  :key="item.id"
                  :value="item.id"
                >
                  {{ item.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="row form-group">
            <div class="col-sm-3" />
            <div class="col-sm-6">
              <button
                type="button"
                class="btn"
                :class="success ? 'btn-outline-primary' : 'btn-primary'"
                :disabled="!canRegister"
                @click.prevent="registerFeed"
              >
                {{ success ? "Saved!..." : "Register" }}
              </button>
            </div>
          </div>
        </div>
      </div>
      <div class="col-md-4">
        <div v-if="previewFeed" class="card previewFeed">
          <div class="card-header">Preview</div>
          <ul class="list-group">
            <li
              v-for="item in previewFeed"
              :key="item.title"
              class="list-group-item"
            >
              {{ item.title }}<br />{{ item.date }}
            </li>
          </ul>
        </div>
      </div>
    </div>
    <hr />
    <CategoryRegister @fetch-list="fetchList" />
    <BackToTop />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive } from "vue";
import BackToTop from "./BackToTop.vue";
import { Agent } from "./UserAgent";
import CategoryRegister from "./addSubscription/Category.vue";

type Categories = {
  id: number;
  name: string;
};

type PreviewFeed = {
  title: string;
  date: string;
};

type ExamineSubscription = null | {
  url: string;
  title: string;
  preview_feed: PreviewFeed[];
};

class Site {
  url = "";
  title = "";
  rss = "";
}

export default defineComponent({
  components: {
    BackToTop,
    CategoryRegister,
  },
  setup: () => {
    const previewFeed = ref<PreviewFeed[]>([]);
    const searchState = ref(false);
    const site = reactive(new Site());

    const categories = ref<Categories[]>([]);
    const success = ref(false);
    const category = ref(0);

    const feedDetail = () => {
      if (site.url === "") {
        return;
      }
      if (!site.url.startsWith("http")) {
        return;
      }
      searchState.value = true;
      Agent<ExamineSubscription>({
        url: "/api/examine_subscription",
        data: {
          url: site.url,
        },
      }).then((data) => {
        if (data === null) {
          alert("Failure: Get information.\n please check url... :(");
          return;
        }
        site.rss = data.url;
        site.title = data.title;
        previewFeed.value = data.preview_feed;
        setTimeout(function () {
          searchState.value = false;
        }, 500);
      });
    };

    const clear = () => {
      site.url = "";
      site.title = "";
      site.rss = "";
      setTimeout(function () {
        success.value = false;
      }, 750);
    };

    const fetchList = () => {
      Agent<Categories[]>({ url: "/api/categories" }).then((data) => {
        categories.value = data;
        if (data.length > 0) {
          category.value = data[0].id;
        }
      });
    };

    const registerFeed = () => {
      Agent<{ result: string }>({
        url: "/api/register_subscription",
        data: {
          ...site,
          category: category.value,
        },
      }).then((data) => {
        if (data === null) {
          alert("情報の取得に失敗しました。URLを確認してください");
          return;
        }
        if (data.result === "ERROR_ALREADY_REGISTER") {
          alert("すでに登録されています。");
          return;
        }
        success.value = true;
        clear();
      });
    };

    fetchList();

    return {
      previewFeed,
      searchState,
      site,
      categories,
      success,
      category,
      fetchList,
      registerFeed,
      feedDetail,
    };
  },
  computed: {
    canRegister: function () {
      return (
        this.site.url &&
        this.site.url.match(/^https?:/g) &&
        this.site.rss &&
        this.site.rss.match(/^https?:/g) &&
        this.site.title &&
        this.category
      );
    },
  },
});
</script>
