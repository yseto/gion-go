<template>
  <div>
    <h4>OPML</h4>

    <div class="row form-group">
      <label class="col-form-label col-sm-4">エクスポート</label>
      <div class="col-sm-8">
        <a class="btn btn-info" @click="opmlExport">エクスポート</a>
      </div>
    </div>

    <div class="row form-group">
      <label class="col-form-label col-sm-4">インポート</label>
      <div class="col-sm-8">
        <label>
          <a class="btn btn-light">ファイルの選択</a>
          <input ref="fileElement" type="file" class="d-none" />
        </label>
      </div>
    </div>
    <div class="row form-group">
      <div class="col-sm-4" />
      <div class="col-sm-8">
        <button class="btn btn-dark" @click="opmlImport">インポート</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import fileDownload from "js-file-download"
import { defineComponent, ref } from "vue"
import { openapiFetchClient } from "../UserAgent"
export default defineComponent({
  setup: () => {
    const fileElement = ref<HTMLInputElement | null>(null)
    const opmlImport = async () => {
      if (fileElement.value === null || fileElement.value?.files === null) {
        return
      }
      const reader = new FileReader()
      await reader.addEventListener(
        "load",
        async function () {
          if (reader.result === null) {
            return
          }
          await openapiFetchClient.POST("/api/opml", {
            body: {
              xml: reader.result.toString(),
            },
          })
          alert("sending done.")
        },
        false
      )
      if (fileElement.value.files[0]) {
        reader.readAsText(fileElement.value.files[0])
      }
    }
    const opmlExport = async () => {
      const { data } = await openapiFetchClient.GET("/api/opml")
      if (data === undefined) {
        return
      }
      fileDownload(data.xml, "opml.xml")
    }
    return {
      fileElement,
      opmlImport,
      opmlExport,
    }
  },
})
</script>
