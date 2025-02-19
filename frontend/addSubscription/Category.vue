<template>
  <div class="row">
    <div class="col-md-8">
      <h4>Categories</h4>
      <div class="form-horizontal">
        <div class="row form-group">
          <label class="col-sm-3 col-form-label" for="inputCategoryName">Name</label>
          <div class="col-sm-6">
            <input v-model="inputCategoryName" type="text" class="form-control" placeholder="Name" />
          </div>
        </div>
        <div class="row form-group">
          <div class="col-sm-3" />
          <div class="col-sm-9">
            <button type="button" class="btn" :class="categorySuccess ? 'btn-outline-primary' : 'btn-primary'"
              :disabled="!!!inputCategoryName" @click.prevent="registerCategory">
              {{ categorySuccess ? "Saved!..." : "Register" }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue"
import { openapiFetchClient } from "../UserAgent"
export default defineComponent({
  emits: ["fetch-list"],
  setup: (_, context) => {
    const inputCategoryName = ref("")
    const categorySuccess = ref(false)

    const registerCategory = async () => {
      const { response } = await openapiFetchClient.POST("/api/category", {
        body: {
          name: inputCategoryName.value,
        },
      })
      if (response.status == 409) {
        alert("すでに登録されています。")
        return
      }
      if (!response.ok) {
        return
      }
      context.emit("fetch-list")
      inputCategoryName.value = ""
      categorySuccess.value = true
      setTimeout(function () {
        categorySuccess.value = false
      }, 750)
    }

    return {
      inputCategoryName,
      categorySuccess,
      registerCategory,
    }
  },
})
</script>
