import { useQuery, useMutation } from "@tanstack/react-query";
import type {
  GetBlocksRequest,
  AddBlockRequest,
  EditBlockRequest,
} from "../../generated/api";
import { blocksApi } from "../client";

export const useGetBlocks = (
  tsumikiID: number,
  params: Omit<GetBlocksRequest, "tsumikiID"> = {},
) => {
  return useQuery({
    queryKey: ["tsumikis", tsumikiID, "blocks", params],
    queryFn: () => blocksApi.getBlocks({ tsumikiID, ...params }),
  });
};

export const useAddBlock = (tsumikiID: number) => {
  return useMutation({
    mutationFn: (addBlockRequest: AddBlockRequest) =>
      blocksApi.addBlock({ tsumikiID, addBlockRequest }),
  });
};

export const useEditBlock = (tsumikiID: number, blockID: number) => {
  return useMutation({
    mutationFn: (editBlockRequest: EditBlockRequest) =>
      blocksApi.editBlock({ tsumikiID, blockID, editBlockRequest }),
  });
};

export const useDeleteBlock = (tsumikiID: number, blockID: number) => {
  return useMutation({
    mutationFn: () => blocksApi.omitBlock({ tsumikiID, blockID }),
  });
};
