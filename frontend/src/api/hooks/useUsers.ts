import { useQuery } from "@tanstack/react-query";
import type { GetMyTsumikisRequest, GetUserTsumikisRequest } from "../../generated/api";
import { usersApi } from "../client";

export const useGetMe = () => {
  return useQuery({
    queryKey: ["users", "me"],
    queryFn: () => usersApi.getMyInfo(),
  });
};

export const useGetUser = (userID: number) => {
  return useQuery({
    queryKey: ["users", userID],
    queryFn: () => usersApi.getUserInfo({ userID }),
  });
};

export const useGetMyTsumikis = (params: GetMyTsumikisRequest = {}) => {
  return useQuery({
    queryKey: ["users", "me", "tsumikis", params],
    queryFn: () => usersApi.getMyTsumikis(params),
  });
};

export const useGetUserTsumikis = (
  userID: number,
  params: Omit<GetUserTsumikisRequest, "userID"> = {},
) => {
  return useQuery({
    queryKey: ["users", userID, "tsumikis", params],
    queryFn: () => usersApi.getUserTsumikis({ userID, ...params }),
  });
};
