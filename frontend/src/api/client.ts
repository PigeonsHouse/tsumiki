import {
  AuthApi,
  Configuration,
  UsersApi,
  WorksApi,
  TsumikisApi,
  BlocksApi,
  ThumbnailsApi,
  type Middleware,
} from "../generated/api";

// アクセストークン期限切れ時に自動でリフレッシュしてリトライするミドルウェア
const tokenRefreshMiddleware: Middleware = {
  post: async ({ fetch, url, init, response }) => {
    // リフレッシュエンドポイント自体の401はループを避けるためスキップ
    if (response.status !== 401 || url.includes("/auth/token-refresh")) {
      return response;
    }

    const refreshResponse = await fetch("/api/v1/auth/token-refresh", {
      credentials: "include",
    });

    if (!refreshResponse.ok) {
      return response;
    }

    // 新しいトークンで元のリクエストをリトライ
    return fetch(url, init);
  },
};

const config = new Configuration({
  credentials: "include",
  middleware: [tokenRefreshMiddleware],
});

export const authApi = new AuthApi(config);
export const usersApi = new UsersApi(config);
export const worksApi = new WorksApi(config);
export const tsumikisApi = new TsumikisApi(config);
export const blocksApi = new BlocksApi(config);
export const thumbnailsApi = new ThumbnailsApi(config);
