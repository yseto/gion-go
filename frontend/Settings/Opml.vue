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
import fileDownload from "js-file-download";
import { defineComponent, ref } from "vue";
import { Agent } from "../UserAgent";
export default defineComponent({
  setup: () => {
    const fileElement = ref<HTMLInputElement | null>(null);
    const opmlImport = () => {
      if (fileElement.value === null || fileElement.value?.files === null) {
        return;
      }
      const reader = new FileReader();
      reader.addEventListener(
        "load",
        function () {
          Agent({
            url: "/api/opml_import",
            data: {
              xml: reader.result,
            },
          }).then(() => alert("sending done."));
        },
        false
      );
      if (fileElement.value.files[0]) {
        reader.readAsText(fileElement.value.files[0]);
      }
    };
    const opmlExport = () => {
      Agent<{ xml: string }>({ url: "/api/opml_export" }).then((data) => {
        fileDownload(data.xml, "opml.xml");
      });
    };
    return {
      fileElement,
      opmlImport,
      opmlExport,
    };
  },
});
</script>
