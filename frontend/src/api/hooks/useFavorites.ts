import { useMutation } from "@tanstack/react-query";
import { tsumikisApi } from "../client";

export const useSetFavorite = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (count: number) =>
      tsumikisApi.setFavorite({ tsumikiID, setFavoriteRequest: { count } }),
  });
};
