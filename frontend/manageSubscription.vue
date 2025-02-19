<template>
  <div class="container">
    <table class="table table-condensed" style="table-layout: fixed">
      <tbody>
        <div v-for="category in subscription" :key="category.id">
          <tr class="row">
            <th class="col-9 text-truncate">
              <span>{{ category.name }}</span>
            </th>
            <td class="col-3 text-right">
              <button class="btn btn-danger btn-sm" @click="removeCategory(category.id, category.name)">
                削除
              </button>
            </td>
          </tr>
          <tr v-for="item in category.subscription" :key="item.id" class="row">
            <td class="col-9 text-truncate">
              <a class="btn btn-link btn-sm" :href="item.siteurl" :title="item.title" target="blank">
                <span v-if="item.http_status >= '400'" class="badge badge-dark">取得に失敗</span>
                <span>{{ item.title }}</span>
              </a>
            </td>
            <td class="col-3 text-right">
              <button class="btn btn-info btn-sm" @click="changeCategory(item.id, item.category_id)">
                移動
              </button>
              <button class="btn btn-danger btn-sm" @click="removeSubscription(item.id, item.title)">
                削除
              </button>
            </td>
          </tr>
        </div>
      </tbody>
    </table>

    <div id="categoryModal" :class="{ 'd-block': categoryModal }" class="modal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h4 class="modal-title">Change: Categories</h4>
          </div>
          <div class="modal-body">
            <label class="col-form-label" for="selectCat">Categories</label>
            <select v-model="fieldCategory" class="form-control" placeholder="Choose Category">
              <option v-for="item in categories" :key="item.id" :value="item.id">
                {{ item.name }}
              </option>
            </select>
          </div>
          <div class="modal-footer">
            <a class="btn btn-success" @click="submit">OK</a>
            <button class="btn btn-light" @click="categoryModal = false">
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
    <BackToTop />
  </div>
  <!--/.container-->
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from "vue"
import BackToTop from "./BackToTop.vue"
import { openapiFetchClient } from "./UserAgent"

type Site = {
  id: number
  category_id: number
  http_status: string
  siteurl: string
  title: string
}

type Subscription = {
  id: number
  name: string
  subscription: Site[]
}

type Category = {
  id: number
  name: string
}

export default defineComponent({
  components: {
    BackToTop,
  },
  setup: () => {
    const categories = ref<Category[]>([])
    const subscription = ref<Subscription[]>([])
    const fieldCategory = ref(0)
    const fieldId = ref(0)
    const categoryModal = ref(false)

    const changeCategory = (id: number, category: number) => {
      fieldCategory.value = category
      fieldId.value = id
      categoryModal.value = true
    }

    const fetchList = async () => {
      const { data } = await openapiFetchClient.GET("/api/subscription")
      if (data === undefined) {
        return
      }
      subscription.value = data
      categories.value = data.map((x) => {
        return { id: x.id, name: x.name }
      })
    }

    const submit = async () => {
      await openapiFetchClient.PUT("/api/subscription/{id}", {
        params: {
          path: {
            id: fieldId.value,
          }
        },
        body: {
          category: fieldCategory.value,
        },
      })
      categoryModal.value = false
      await fetchList()
    }

    const removeCategory = async (id: number, name: string) => {
      if (!confirm("カテゴリ:" + name + " を削除しますか?\n内包されている購読もすべて削除されます")) {
        return
      }
      await openapiFetchClient.DELETE("/api/category/{id}", {
        params: {
          path: { id }
        },
      })
      await fetchList()
    }

    const removeSubscription = async (id: number, name: string) => {
      if (!confirm(name + " を削除しますか")) {
        return
      }
      await openapiFetchClient.DELETE("/api/subscription/{id}", {
        params: {
          path: { id }
        },
      })
      await fetchList()
    }

    onMounted(async () => await fetchList())

    return {
      categories,
      subscription,
      fieldCategory,
      fieldId,
      categoryModal,

      changeCategory,
      removeCategory,
      removeSubscription,
      submit,
    }
  },
})
</script>
