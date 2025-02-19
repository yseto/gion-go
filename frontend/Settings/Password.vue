<template>
  <div>
    <h4>パスワード設定</h4>
    <div role="form">
      <div class="row form-group">
        <label class="col-form-label col-sm-4" for="passwordOld">今のパスワード</label>
        <div class="col-sm-8">
          <input id="passwordOld" v-model="passwordOld" type="password" placeholder="8文字以上" class="form-control" />
        </div>
      </div>
      <div class="row form-group">
        <label class="col-form-label col-sm-4" for="password">新しいパスワード</label>
        <div class="col-sm-8">
          <input id="password" v-model="password" type="password" placeholder="8文字以上" class="form-control" />
        </div>
      </div>
      <div class="row form-group">
        <label class="col-form-label col-sm-4" for="passwordChecked">新しいパスワード(確認)</label>
        <div class="col-sm-8">
          <input id="passwordChecked" v-model="passwordChecked" type="password" placeholder="8文字以上"
            class="form-control" />
        </div>
      </div>
      <div class="row form-group">
        <div class="col-sm-4" />
        <div class="col-sm-8">
          <a class="btn btn-primary" @click.prevent="updatePassword">Password Change.</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { openapiFetchClient } from "../UserAgent";
export default defineComponent({
  setup: () => {
    const password = ref("");
    const passwordOld = ref("");
    const passwordChecked = ref("");

    const updatePassword = () => {
      openapiFetchClient.POST("/api/update_password", {
        body: {
          password_old: passwordOld.value,
          password: password.value,
          passwordc: passwordChecked.value,
        },
      }).then((data) => alert(data.data?.result));
    };

    return {
      password,
      passwordOld,
      passwordChecked,
      updatePassword,
    };
  },
});
</script>
