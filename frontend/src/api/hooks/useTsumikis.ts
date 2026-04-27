import { useQuery, useMutation } from "@tanstack/react-query";
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

export const useGetTsumiki = (tsumikiID: number) => {
  return useQuery({
    queryKey: ["tsumikis", tsumikiID],
    queryFn: () => tsumikisApi.getSpecifiedTsumiki({ tsumikiID }),
  });
};

export const useCreateTsumiki = () => {
  return useMutation({
    mutationFn: (createTsumikiRequest: CreateTsumikiRequest) =>
      tsumikisApi.createTsumiki({ createTsumikiRequest }),
  });
};

export const useEditTsumiki = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (editTsumikiRequest: EditTsumikiRequest) =>
      tsumikisApi.editTsumiki({ tsumikiID, editTsumikiRequest }),
  });
};

export const useDeleteTsumiki = (tsumikiID: number) => {
  return useMutation({
    mutationFn: () => tsumikisApi.deleteTsumiki({ tsumikiID }),
  });
};

export const useUpdateTsumikiThumbnail = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (updateTsumikiThumbnailRequest: UpdateTsumikiThumbnailRequest) =>
      tsumikisApi.updateTsumikiThumbnail({ tsumikiID, updateTsumikiThumbnailRequest }),
  });
};

export const useUploadTsumikiMedia = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (file: Blob) => tsumikisApi.postMedia({ tsumikiID, file }),
  });
};
