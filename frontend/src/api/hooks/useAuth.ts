import { useMutation, useQueryClient } from "@tanstack/react-query";
import { authApi } from "../client";

export const useLogout = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      // ログイン情報キャッシュを削除してヘッダーを即時更新
      queryClient.removeQueries({ queryKey: ["users", "me"] });
    },
  });
};
