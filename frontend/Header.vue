<template>
  <div class="bg-light my-navigation">
    <nav v-if="$route.meta.anonymous ? false : true" class="navbar navbar-expand-lg navbar-light container">
      <span class="navbar-brand">Gion</span>
      <button class="navbar-toggler" type="button" @click="navbar">
        <span class="navbar-toggler-icon" />
      </button>
      <div class="collapse navbar-collapse" :class="{ show: navbarState }">
        <ul class="navbar-nav mr-auto">
          <li v-for="item in items" :key="item.caption" class="nav-item"
            :class="{ active: $root?.$route.path === item.route }">
            <a class="nav-link" @click="go(item.route)">{{ item.caption }}</a>
          </li>
        </ul>
        <ul class="navbar-nav">
          <li class="nav-item" :class="{ active: $root?.$route.path === '/settings' }">
            <a class="nav-link" @click="go('/settings')">Settings</a>
          </li>
          <li class="nav-item hidden-sm">
            <a class="nav-link" style="cursor: pointer" @click="helpModal = true">Help</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" @click="go('/logout')">Logout</a>
          </li>
        </ul>
      </div>
    </nav>

    <div id="helpModal" :class="{ 'd-block': helpModal }" class="modal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Help</h5>
          </div>
          <div class="modal-body">
            <p>Pin Listでは、ピン止めしたエントリの一覧が表示されます。</p>

            <hr />

            <p>
              URLに巡回したいWebサイトのアドレスを入力して、<a class="btn btn-info">Get Detail</a>をクリックすると、必要な情報を取得します。
            </p>
            <p>
              この方法で取得できない場合は、Webページのデータを確認してください。RSSを配信していない可能性があります。
            </p>
            <hr />
            <p>
              カテゴリを増やしたい場合は、下の入力欄にカテゴリの名前を入力してください。この時追加したカテゴリはすぐに、サイト登録でお使いになれます。
            </p>

            <hr />

            <table class="table table-striped table-bordered">
              <tr>
                <th>Key</th>
                <th>Description</th>
              </tr>
              <tr>
                <th>A</th>
                <td>一つ前のカテゴリ</td>
              </tr>
              <tr>
                <th>S</th>
                <td>一つ次のカテゴリ</td>
              </tr>
              <tr>
                <th>K</th>
                <td>一つ前のアイテムを選択する</td>
              </tr>
              <tr>
                <th>J</th>
                <td>一つ次のアイテムを選択する</td>
              </tr>
              <tr>
                <th>I</th>
                <td>ピンを立ててない表示中のアイテムを既読にする</td>
              </tr>
              <tr>
                <th>O</th>
                <td>ピンリストを開く、閉じる</td>
              </tr>
              <tr>
                <th>P</th>
                <td>ピンを立てる、外す</td>
              </tr>
              <tr>
                <th>V</th>
                <td>アイテムを開く</td>
              </tr>
              <tr>
                <th>R</th>
                <td>ページを再読み込みする</td>
              </tr>
            </table>

            <hr />

            <p>変更ボタンで巡回先のカテゴリを移動させることができます。</p>
            <p>削除ボタンで巡回先を巡回対象から外すことができます。</p>
            <p>
              カテゴリを削除した場合、カテゴリ以下に登録されている巡回先も同時に削除されます。
            </p>

            <hr />
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="helpModal = false">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue"
import { useRoute, useRouter } from "vue-router"
export default defineComponent({
  setup: () => {
    const router = useRouter()
    const route = useRoute()

    const items = [
      {
        caption: "Pin List",
        route: "/pin",
      },
      {
        caption: "Read entries",
        route: "/entry",
      },
      {
        caption: "Add a new subscription",
        route: "/add",
      },
      {
        caption: "Manage subscription",
        route: "/subscription",
      },
    ]
    const navbarState = ref(false)
    const helpModal = ref(false)
    const navbar = () => {
      navbarState.value = navbarState.value ? false : true
    }
    const go = (to: string) => {
      navbarState.value = false

      if (route.path !== to) {
        router.push(to)
      }
    }
    return {
      items,
      navbar,
      go,
      navbarState,
      helpModal,
    }
  },
})
</script>
