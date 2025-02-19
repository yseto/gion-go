<template>
  <div role="form">
    <div class="row form-group">
      <label class="col-form-label col-sm-4" for="entryCount"
        >表示件数の上限</label
      >
      <div class="col-sm-8">
        <input
          v-model="profile.entryCount"
          type="number"
          placeholder="0で無制限"
          class="form-control"
        />
        <span class="form-text">一度に表示する件数の上限を設定できます。</span>
      </div>
    </div>
    <div class="row form-group">
      <label class="col-form-label col-sm-4" for="substringLength"
        >概要の文字数制限</label
      >
      <div class="col-sm-8">
        <input
          v-model="profile.substringLength"
          type="number"
          placeholder="0で無制限"
          class="form-control"
        />
        <span class="form-text">概要の文字数の上限を設定できます。</span>
      </div>
    </div>
    <div class="row form-group">
      <label class="col-form-label col-sm-4">その他の設定</label>
      <div class="col-sm-8">
        <div class="form-check">
          <input
            id="onLoginSkipPinList"
            v-model="profile.onLoginSkipPinList"
            type="checkbox"
            class="form-check-input"
          />
          <label class="form-check-label" for="onLoginSkipPinList"
            >ログインしたらすぐにエントリ一覧を表示する</label
          >
        </div>

        <div class="form-check">
          <input
            id="autoseen"
            v-model="profile.autoseen"
            type="checkbox"
            class="form-check-input"
          />
          <label class="form-check-label" for="autoseen"
            >エントリーを自動既読にする</label
          >
        </div>
      </div>
    </div>
    <div class="row form-group">
      <div class="col-sm-4" />
      <div class="col-sm-8">
        <button
          class="btn"
          :class="finished ? 'btn-outline-primary' : 'btn-primary'"
          @click="apply"
        >
          {{ finished ? "Saved!..." : "OK" }}
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from "vue";
import { openapiFetchClient } from "../UserAgent";
import { useUserStore } from "../UserStore";

class CProfile {
  autoseen = false;
  onLoginSkipPinList = false;
  entryCount = 0;
  substringLength = 0;
}

export default defineComponent({
  setup: () => {
    const store = useUserStore();
    const profile = reactive(new CProfile());
    const finished = ref(false);

    const apply = () => {
      openapiFetchClient.POST("/api/set_profile", {
        body: {
          autoseen: !!profile.autoseen,
          onLoginSkipPinList: !!profile.onLoginSkipPinList,
          entryCount: profile.entryCount,
          substringLength: profile.substringLength,
        },
      }).then(() => {
        store.user.autoSeen = profile.autoseen;
        finished.value = true;
        setTimeout(function () {
          finished.value = false;
        }, 1000);
      });
    };

    openapiFetchClient.POST("/api/profile").then((data) => {
      if (data.data === undefined) {
        return
      }
      profile.autoseen = !!data.data.autoseen;
      profile.onLoginSkipPinList = !!data.data.onLoginSkipPinList;
      profile.entryCount = data.data.entryCount;
      profile.substringLength = data.data.substringLength;
    });

    return {
      store,
      profile,
      finished,
      apply,
    };
  },
});
</script>
