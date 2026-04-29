import { useParams, useNavigate } from "react-router";
import {
  useGetBlocks,
  useGetTsumiki,
  useDeleteTsumiki,
  useAddBlock,
  useUploadTsumikiMedia,
} from "../../../api";
import { useMemo, useState } from "react";
import NotFound from "../../../components/NotFound";
import { useQueryClient } from "@tanstack/react-query";

type BlockFormValues = {
  message: string;
  percentage: number;
  condition: number;
};

const TsumikiDetail = () => {
  const { tsumikiId: tsumikiIdRaw } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isValidId = useMemo(
    () =>
      tsumikiIdRaw !== "" &&
      tsumikiIdRaw !== undefined &&
      !Number.isNaN(Number(tsumikiIdRaw)),
    [tsumikiIdRaw],
  );
  const tsumikiId = useMemo(() => Number(tsumikiIdRaw), [tsumikiIdRaw]);
  const { data: tsumiki, isError } = useGetTsumiki(tsumikiId, isValidId);
  const { data: blocks, refetch: refetchBlocks } = useGetBlocks(tsumikiId, isValidId);
  const { mutateAsync: deleteTsumiki } = useDeleteTsumiki();
  const { mutateAsync: addBlock, isPending: isAddingBlock } = useAddBlock(tsumikiId);
  const { mutateAsync: uploadMedia } = useUploadTsumikiMedia(tsumikiId);

  const [form, setForm] = useState<BlockFormValues>({
    message: "",
    percentage: 0,
    condition: 3,
  });
  const [mediaFiles, setMediaFiles] = useState<FileList | null>(null);
  const [showAddBlock, setShowAddBlock] = useState(false);

  const latestBlockId = useMemo(() => {
    if (!blocks || blocks.length === 0) return null;
    return blocks[blocks.length - 1].id;
  }, [blocks]);

  const handleDelete = async () => {
    if (!confirm("この積み木を削除しますか？")) return;
    await deleteTsumiki(tsumikiId);
    navigate("/tsumikis");
  };

  const handleAddBlock = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.message && (!mediaFiles || mediaFiles.length === 0)) return;

    const mediaIds: number[] = [];
    if (mediaFiles && mediaFiles.length > 0) {
      const files = Array.from(mediaFiles).slice(0, 4);
      for (const file of files) {
        const media = await uploadMedia(file);
        mediaIds.push(media.id);
      }
    }

    await addBlock({
      message: form.message || null,
      percentage: Number(form.percentage),
      condition: Number(form.condition),
      mediaIds: mediaIds.length > 0 ? mediaIds : undefined,
      latestBlockId,
    });

    setForm({ message: "", percentage: 0, condition: 3 });
    setMediaFiles(null);
    setShowAddBlock(false);
    await refetchBlocks();
    queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiId] });
  };

  return (!isValidId || isError) ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>積み木詳細</h1>
      {tsumiki && blocks && (
        <div style={{ width: 1024, margin: "auto" }}>
          <img
            style={{
              aspectRatio: "16 / 9",
              width: "100%",
              objectFit: "contain",
              backgroundColor: "gray",
            }}
            src={tsumiki.thumbnailUrl || ""}
          />
          <h2>{tsumiki.title}</h2>
          {tsumiki.percentage != null && (
            <div>進捗度: {tsumiki.percentage}%</div>
          )}
          {tsumiki.work && (
            <a href={`/works/${tsumiki.work.id}`}>
              <div
                style={{
                  borderRadius: 8,
                  padding: 8,
                  border: "1px solid black",
                  display: "flex",
                  gap: 16,
                  alignItems: "center",
                }}
              >
                <img
                  style={{
                    height: 48,
                    aspectRatio: "16 / 9",
                    objectFit: "contain",
                    backgroundColor: "gray",
                  }}
                  src={tsumiki.work.thumbnailUrl || ""}
                />
                作品：{tsumiki.work.title}
              </div>
            </a>
          )}
          <a
            style={{ display: "flex", alignItems: "center", gap: 8 }}
            href={`/users/${tsumiki.user.id}`}
          >
            <img
              style={{ width: 40, height: 40, borderRadius: 9999 }}
              src={tsumiki.user.avatarUrl}
            />
            <span>{tsumiki.user.name}</span>
          </a>
          <small>{tsumiki.createdAt.toISOString()}</small>
          {tsumiki.isOwner && (
            <div style={{ display: "flex", gap: 8, marginTop: 12 }}>
              <a href={`/tsumikis/${tsumikiId}/edit`}>
                <button type="button">編集</button>
              </a>
              <button type="button" onClick={handleDelete}>削除</button>
            </div>
          )}
          {blocks.map((block) => (
            <div key={block.id}>
              <hr />
              <div>
                {block.isDeleted ? (
                  <>削除済み</>
                ) : (
                  <>
                    <small>
                      コンディション: {block.condition}/5 完成度:{" "}
                      {block.percentage}%
                    </small>
                    {block.medias &&
                      block.medias.map((media) => (
                        <img
                          key={media.id}
                          style={{
                            aspectRatio: "16 / 9",
                            width: "100%",
                            objectFit: "contain",
                            backgroundColor: "gray",
                          }}
                          src={media.url}
                        />
                      ))}
                    {block.message && <p>{block.message}</p>}
                    <small>{block.createdAt?.toISOString()}</small>
                  </>
                )}
              </div>
            </div>
          ))}
          {tsumiki.isOwner && (
            <>
              <hr />
              {showAddBlock ? (
                <form
                  onSubmit={handleAddBlock}
                  style={{
                    display: "flex",
                    flexDirection: "column",
                    gap: 12,
                    paddingBottom: 24,
                  }}
                >
                  <div>
                    <label>進捗率: {form.percentage}%</label>
                    <br />
                    <input
                      type="range"
                      min={0}
                      max={100}
                      step={1}
                      style={{ width: "100%" }}
                      value={form.percentage}
                      onChange={(e) =>
                        setForm((v) => ({ ...v, percentage: Number(e.target.value) }))
                      }
                    />
                  </div>
                  <div>
                    <label>コンディション: {form.condition} / 5</label>
                    <br />
                    <input
                      type="range"
                      min={1}
                      max={5}
                      step={1}
                      style={{ width: "100%" }}
                      value={form.condition}
                      onChange={(e) =>
                        setForm((v) => ({ ...v, condition: Number(e.target.value) }))
                      }
                    />
                  </div>
                  <div>
                    <label>メッセージ（200文字以内）</label>
                    <br />
                    <textarea
                      rows={3}
                      style={{ width: "100%" }}
                      maxLength={200}
                      value={form.message}
                      onChange={(e) =>
                        setForm((v) => ({ ...v, message: e.target.value }))
                      }
                    />
                  </div>
                  <div>
                    <label>メディア（最大4件）</label>
                    <br />
                    <input
                      type="file"
                      multiple
                      accept="image/jpeg,image/png,image/gif,audio/mpeg,audio/wav,audio/ogg,audio/aac,video/mp4,video/webm,video/quicktime"
                      onChange={(e) => setMediaFiles(e.target.files)}
                    />
                  </div>
                  <div style={{ display: "flex", gap: 8 }}>
                    <button type="submit" disabled={isAddingBlock}>
                      {isAddingBlock ? "追加中..." : "追加する"}
                    </button>
                    <button type="button" onClick={() => setShowAddBlock(false)}>
                      キャンセル
                    </button>
                  </div>
                </form>
              ) : (
                <button
                  type="button"
                  style={{ marginBottom: 24 }}
                  onClick={() => setShowAddBlock(true)}
                >
                  ブロックを追加
                </button>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default TsumikiDetail;
