import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  GetWorksRequest,
  GetWorkTsumikiRequest,
  CreateWorkRequest,
  EditWorkRequest,
  UpdateTsumikiThumbnailRequest,
} from "../../generated/api";
import { worksApi } from "../client";

export const useGetWorks = (params: GetWorksRequest = {}) => {
  return useQuery({
    queryKey: ["works", params],
    queryFn: () => worksApi.getWorks(params),
  });
};

export const useGetWork = (workID: number, enabled: boolean) => {
  return useQuery({
    queryKey: ["works", workID],
    queryFn: () => worksApi.getSpecifiedWork({ workID }),
    enabled,
  });
};

export const useCreateWork = () => {
  return useMutation({
    mutationFn: (createWorkRequest: CreateWorkRequest) =>
      worksApi.createWork({ createWorkRequest }),
  });
};

export const useEditWork = (workID: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (editWorkRequest: EditWorkRequest) =>
      worksApi.editWork({ workID, editWorkRequest }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["works", workID] });
    },
  });
};

export const useDeleteWork = (workID: number) => {
  return useMutation({
    mutationFn: () => worksApi.deleteWork({ workID }),
  });
};

export const useUpdateWorkThumbnail = (workID: number) => {
  return useMutation({
    mutationFn: (
      updateTsumikiThumbnailRequest: UpdateTsumikiThumbnailRequest,
    ) =>
      worksApi.updateWorkThumbnail({ workID, updateTsumikiThumbnailRequest }),
  });
};

export const useGetWorkTsumikis = (
  workID: number,
  params: Omit<GetWorkTsumikiRequest, "workID"> = {},
) => {
  return useQuery({
    queryKey: ["works", workID, "tsumikis", params],
    queryFn: () => worksApi.getWorkTsumiki({ workID, ...params }),
  });
};
