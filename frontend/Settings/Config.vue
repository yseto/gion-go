<template>
  <div role="form">
    <div class="row form-group">
      <label class="col-form-label col-sm-4" for="numentry"
        >表示件数の上限</label
      >
      <div class="col-sm-8">
        <input
          v-model="profile.numentry"
          type="number"
          placeholder="0で無制限"
          class="form-control"
        />
        <span class="form-text">一度に表示する件数の上限を設定できます。</span>
      </div>
    </div>
    <div class="row form-group">
      <label class="col-form-label col-sm-4" for="numsubstr"
        >概要の文字数制限</label
      >
      <div class="col-sm-8">
        <input
          v-model="profile.numsubstr"
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
            id="nopinlist"
            v-model="profile.nopinlist"
            type="checkbox"
            class="form-check-input"
          />
          <label class="form-check-label" for="nopinlist"
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
import { Agent } from "../UserAgent";
import { useUserStore } from "../UserStore";
import { Profile } from "../types";

class CProfile {
  autoseen = false;
  nopinlist = false;
  numentry = 0;
  numsubstr = 0;
}

export default defineComponent({
  setup: () => {
    const store = useUserStore();
    const profile = reactive(new CProfile());
    const finished = ref(false);

    const apply = () => {
      Agent({
        url: "/api/set_profile",
        jsonRequest: true,
        data: {
          autoseen: !!profile.autoseen,
          nopinlist: !!profile.nopinlist,
          numentry: profile.numentry,
          numsubstr: profile.numsubstr,
        },
      }).then(() => {
        store.user.autoSeen = profile.autoseen;
        finished.value = true;
        setTimeout(function () {
          finished.value = false;
        }, 1000);
      });
    };

    Agent<Profile>({ url: "/api/profile" }).then((data) => {
      profile.autoseen = !!data.autoseen;
      profile.nopinlist = !!data.nopinlist;
      profile.numentry = data.numentry;
      profile.numsubstr = data.numsubstr;
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
