import { useMutation } from "@tanstack/react-query";
import { authApi } from "../client";

export const useLogout = () => {
  return useMutation({
    mutationFn: () => authApi.logout(),
  });
};
