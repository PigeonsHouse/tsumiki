# 環境構築

開発・デプロイではDockerを利用しています。
Windows・MacではDockerDesktopやRancherDesktop、Linuxでは任意のパッケージマネージャから導入してください。

`docker -v` などのコマンドでDockerが動作していることを確認できたらOKです。

## 開発環境

1. コンテナをビルドする(初回のみ)

   ```sh
   docker compose --profile dev build
   ```

2. コンテナを起動する

   ```sh
   docker compose --profile dev up -d
   ```

3. (ライブラリを追加した際)再ビルドを行う
   package.json に追加したら `frontend`、go.mod に追加したら `backend`

   ```sh
   docker compose --profile dev build backend
   ```

## 本番環境

1. コンテナをビルドする(初回のみ)

   ```sh
   docker compose --profile prod build
   ```

2. コンテナを起動する

   ```sh
   docker compose --profile prod up -d
   ```
