import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  GetBlocksRequest,
  AddBlockRequest,
  EditBlockRequest,
} from "../../generated/api";
import { blocksApi } from "../client";

export const useGetBlocks = (
  tsumikiID: number,
  enabled: boolean,
  params: Omit<GetBlocksRequest, "tsumikiID"> = {}
) => {
  return useQuery({
    queryKey: ["tsumikis", tsumikiID, "blocks", params],
    queryFn: () => blocksApi.getBlocks({ tsumikiID, ...params }),
    enabled,
  });
};

export const useAddBlock = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (addBlockRequest: AddBlockRequest) =>
      blocksApi.addBlock({ tsumikiID, addBlockRequest }),
  });
};

export const useEditBlock = (tsumikiID: number, blockID: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (editBlockRequest: EditBlockRequest) =>
      blocksApi.editBlock({ tsumikiID, blockID, editBlockRequest }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiID, "blocks"] });
    },
  });
};

export const useDeleteBlock = (tsumikiID: number, blockID: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => blocksApi.omitBlock({ tsumikiID, blockID }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiID, "blocks"] });
    },
  });
};
