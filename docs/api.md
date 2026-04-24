# Tsumiki API設計

※ あくまで下書きで、固まったら後でopenapiに移行する

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
