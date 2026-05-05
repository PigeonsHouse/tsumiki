import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  GetTsumikisRequest,
  CreateTsumikiRequest,
  EditTsumikiRequest,
  UpdateTsumikiThumbnailRequest,
} from "../../generated/api";
import { tsumikisApi } from "../client";

export const useGetTsumikis = (params: GetTsumikisRequest = {}) => {
  return useQuery({
    queryKey: ["tsumikis", params],
    queryFn: () => tsumikisApi.getTsumikis(params),
  });
};

export const useGetTsumiki = (tsumikiID: number, enabled: boolean) => {
  return useQuery({
    queryKey: ["tsumikis", tsumikiID],
    queryFn: () => tsumikisApi.getSpecifiedTsumiki({ tsumikiID }),
    enabled,
  });
};

export const useCreateTsumiki = () => {
  return useMutation({
    mutationFn: (createTsumikiRequest: CreateTsumikiRequest) =>
      tsumikisApi.createTsumiki({ createTsumikiRequest }),
  });
};

export const useEditTsumiki = (tsumikiID: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (editTsumikiRequest: EditTsumikiRequest) =>
      tsumikisApi.editTsumiki({ tsumikiID, editTsumikiRequest }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiID] });
    },
  });
};

export const useDeleteTsumiki = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (tsumikiID: number) => tsumikisApi.deleteTsumiki({ tsumikiID }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users", "me", "tsumikis"] });
    },
  });
};

export const useUpdateTsumikiThumbnail = (tsumikiID: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (
      updateTsumikiThumbnailRequest: UpdateTsumikiThumbnailRequest
    ) =>
      tsumikisApi.updateTsumikiThumbnail({
        tsumikiID,
        updateTsumikiThumbnailRequest,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiID] });
    },
  });
};

export const useUploadTsumikiMedia = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (file: Blob) => tsumikisApi.postMedia({ tsumikiID, file }),
  });
};
