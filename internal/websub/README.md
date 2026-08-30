# WebSubによるフィードエントリの受信

WebSubによる通知で、フィードエントリを生成する

## 現在の適用先

- https://developers.google.com/youtube/v3/guides/push_notifications?hl=ja

## 課題点

### Youtube限定となっていること

以下を克服すれば、利用箇所は増やせるはず

- フィード登録時に、フィードを取得してWebSub対応であるかを調べる
  - rel="hub", rel="self" の `<link>` タグを確認する
- WebSubのフィードサーバーをDBに記録する
  - Youtubeの場合、Googleのサーバーであることがわかっており、保存していない

## 参考

https://zenn.dev/tatsumin/articles/youtube-pubsub