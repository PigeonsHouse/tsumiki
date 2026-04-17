# 環境構築

開発・デプロイではDockerを利用しています。
Windows・MacではDockerDesktopやRancherDesktop、Linuxでは任意のパッケージマネージャから導入してください。

`docker -v` などのコマンドでDockerが動作していることを確認できたらOKです。

## 開発環境

1. サンプルをコピーしてenvファイルを作成し、必要な値を入力する

   ```
   backend/.env.sample -> backend/.env
   database/.env.sample -> database/.env
   redis/.env.sample -> redis/.env
   ```

1. DBを起動してマイグレーションする(初回のみ)
   テーブル定義を更新した際もこのコマンドを実行する

   ```sh
   docker compose up -d db
   docker compose --profile migrate up --build migrate
   ```

1. コンテナをビルドする(初回のみ)

   ```sh
   docker compose --profile dev build
   ```

1. コンテナを起動する

   ```sh
   docker compose --profile dev up -d
   ```

1. (ライブラリを追加した際)再ビルドを行う
   package.json に追加したら `frontend`、go.mod に追加したら `backend`

   ```sh
   docker compose --profile dev build backend
   ```

## 本番環境

1. サンプルをコピーしてenvファイルを作成し、必要な値を入力する

   ```
   backend/.env.sample -> backend/.env.dev
   database/.env.sample -> database/.env
   redis/.env.sample -> redis/.env
   ```

1. DBを起動してマイグレーションする(初回のみ)

   ```sh
   docker compose up -d db
   docker compose --profile migrate up --build migrate
   ```

1. コンテナをビルドする

   ```sh
   docker compose --profile prod build
   ```

1. コンテナを起動する

   ```sh
   docker compose --profile prod up -d
   ```
