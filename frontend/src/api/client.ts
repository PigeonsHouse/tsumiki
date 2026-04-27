import {
  Configuration,
  UsersApi,
  WorksApi,
  TsumikisApi,
  BlocksApi,
  ThumbnailsApi,
} from "../generated/api";

const config = new Configuration({ credentials: "include" });

export const usersApi = new UsersApi(config);
export const worksApi = new WorksApi(config);
export const tsumikisApi = new TsumikisApi(config);
export const blocksApi = new BlocksApi(config);
export const thumbnailsApi = new ThumbnailsApi(config);
