# Tsumiki API設計

※ あくまで下書きで、固まったら後でopenapiに移行する

## 必須

- GET /api/v1/auth/discord
  - ログイン画面にリダイレクトする
- GET /api/v1/auth/discord/callback
  - 認証後の情報をサーバで処理するコールバックエンドポイント
  - API用トークンとリフレッシュトークンを返す
- POST /api/v1/token/refresh
  - リフレッシュトークンを受け取り、問題なければAPIトークンを返す
- POST /api/v1/tsumikis
  - 積み木＆1つ目のブロック作成
- PUT /api/v1/tsumikis/{tsumiki_id}
  - 積み木の編集
- DELETE /api/v1/tsumikis/{tsumiki_id}
  - 積み木の削除
- POST /api/v1/tsumikis/{tsumiki_id}/blocks
  - ブロック追加
  - メディアの紐づけもこのエンドポイントでやる
- PUT /api/v1/tsumikis/{tsumiki_id}/blocks/{block_id}
  - ブロック編集
- DELETE /api/v1/tsumikis/{tsumiki_id}/blocks/{block_id}
  - ブロック削除
  - 前後関係を維持するためにソフトデリートになる
- POST /api/v1/tsumikis/{tsumiki_id}/blocks/{block_id}/medias
  - ブロックに紐づけるメディアの投稿
  - リレーションはこの段階ではつけない
- GET /api/v1/works
  - 作品一覧取得
- GET /api/v1/works/{work_id}
  - 作品取得
- GET /api/v1/works/{work_id}/tsumikis
  - 作品に関わる積み木一覧取得
- POST /api/v1/works
  - 作品作成
- PUT /api/v1/works/{work_id}
  - 作品編集
- DELETE /api/v1/works/{work_id}
  - 作品削除
- GET /api/v1/tsumikis
  - 積み木一覧取得
  - フィルタ、検索機能
- GET /api/v1/tsumikis/{tsumiki_id}
  - 積み木とそのブロック一覧取得

## v1.0までにほしいやつ

- PUT /api/v1/tsumikis/{tsumiki_id}/favorites
  - 積み木に送るいいね数を更新する
  - クライアントで連打する可能性を考えて、debounceして最終値を送る
  - GETは積み木取得でついでに付ける
- POST /api/v1/tsumikis/{tsumiki_id}/comments
  - 積み木にコメントする
  - ブロックにメンションする事もできる(ブロック向けコメント)
- GET /api/v1/tsumikis/{tsumiki_id}/comments
  - 積み木へのコメントを取得する
- PUT /api/v1/tsumikis/{tsumiki_id}/comments/{comment_id}
  - コメントを編集する
- DELETE /api/v1/tsumikis/{tsumiki_id}/comments/{comment_id}
  - コメントを削除する
- GET /api/v1/tsumikis/{tsumiki_id}/reactions
  - 積み木へのリアクションを取得する
- PUT /api/v1/tsumikis/{tsumiki_id}/reactions
  - 積み木へ自分が送ったリアクション一覧を更新する
- GET /api/v1/tsumikis/{tsumiki_id}/blocks/{block_id}/reactions
  - 積み木のブロックへのリアクションを取得する
- PUT /api/v1/tsumikis/{tsumiki_id}/blocks/{block_id}/reactions
  - 積み木のブロックへ自分が送ったリアクション一覧を更新する
