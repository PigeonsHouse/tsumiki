import { useMutation } from "@tanstack/react-query";
import { thumbnailsApi } from "../client";

export const useUploadThumbnail = () => {
  return useMutation({
    mutationFn: (thumbnail: Blob) => thumbnailsApi.postThumbnail({ thumbnail }),
  });
};
