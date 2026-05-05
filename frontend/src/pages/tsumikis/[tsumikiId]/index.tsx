import { useEffect, useMemo, useRef, useState } from "react";
import { useParams, useNavigate } from "react-router";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetBlocks,
  useGetTsumiki,
  useDeleteTsumiki,
  useAddBlock,
  useUploadTsumikiMedia,
  useSetFavorite,
} from "../../../api";
import NotFound from "../../../components/NotFound";
import { useAuth } from "../../../contexts";

const MAX_FAVORITE_COUNT = 10;

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
    [tsumikiIdRaw]
  );
  const tsumikiId = useMemo(() => Number(tsumikiIdRaw), [tsumikiIdRaw]);
  const { data: tsumiki, isError } = useGetTsumiki(tsumikiId, isValidId);
  const { data: blocks, refetch: refetchBlocks } = useGetBlocks(
    tsumikiId,
    isValidId
  );
  const { mutateAsync: deleteTsumiki } = useDeleteTsumiki();
  const { mutateAsync: addBlock, isPending: isAddingBlock } =
    useAddBlock(tsumikiId);
  const { mutateAsync: uploadMedia } = useUploadTsumikiMedia(tsumikiId);
  const { mutateAsync: setFavorite } = useSetFavorite(tsumikiId);
  const { me } = useAuth();

  const [localFavoriteCount, setLocalFavoriteCount] = useState(0);
  const pendingCountRef = useRef(0);
  const favoriteTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (tsumiki) {
      setLocalFavoriteCount(tsumiki.favorite.myFavoriteCount);
      pendingCountRef.current = tsumiki.favorite.myFavoriteCount;
    }
  }, [tsumiki?.favorite.myFavoriteCount]);

  useEffect(() => {
    return () => {
      if (favoriteTimerRef.current) clearTimeout(favoriteTimerRef.current);
    };
  }, []);

  const handleFavorite = () => {
    if (pendingCountRef.current >= MAX_FAVORITE_COUNT) return;
    pendingCountRef.current += 1;
    setLocalFavoriteCount(pendingCountRef.current);
    if (favoriteTimerRef.current) clearTimeout(favoriteTimerRef.current);
    favoriteTimerRef.current = setTimeout(() => {
      setFavorite(pendingCountRef.current);
    }, 1000);
  };

  const [form, setForm] = useState<BlockFormValues>({
    message: "",
    percentage: 0,
    condition: 3,
  });
  const [mediaFiles, setMediaFiles] = useState<File[]>([]);
  const mediaInputRef = useRef<HTMLInputElement>(null);
  const [showAddBlock, setShowAddBlock] = useState(false);

  useEffect(() => {
    if (tsumiki && typeof tsumiki.percentage === "number") {
      const percentage = tsumiki.percentage; // tsの型推論のため一度変数に代入
      setForm((v) => ({ ...v, percentage }));
    }
  }, [tsumiki, setForm]);

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
    if (!form.message && mediaFiles.length === 0) return;

    const mediaIds: number[] = [];
    if (mediaFiles.length > 0) {
      for (const file of mediaFiles) {
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
    setMediaFiles([]);
    setShowAddBlock(false);
    await refetchBlocks();
    queryClient.invalidateQueries({ queryKey: ["tsumikis", tsumikiId] });
  };

  return !isValidId || isError ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>積み木詳細</h1>
      {tsumiki && blocks && (
        <div style={{ maxWidth: 1024, margin: "auto" }}>
          <img
            style={{
              aspectRatio: "16 / 9",
              width: "100%",
              objectFit: "contain",
              backgroundColor: "gray",
            }}
            src={tsumiki.thumbnailUrl || ""}
          />
          <div style={{ display: "flex", justifyContent: "space-between" }}>
            <h2>{tsumiki.title}</h2>
            <div style={{ display: "flex" }}>
              {me && (
                <button onClick={handleFavorite}>
                  いいね: {localFavoriteCount}
                </button>
              )}
            </div>
          </div>
          <div style={{ display: "flex", gap: 8 }}>
            {tsumiki.percentage != null && (
              <span>進捗度: {tsumiki.percentage}%</span>
            )}
            <span>いいね: {tsumiki.favorite.totalFavoriteCount}</span>
          </div>
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
              <button type="button" onClick={handleDelete}>
                削除
              </button>
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
                    {block.medias && block.medias.length > 0 && (
                      <div
                        style={{
                          maxWidth: 800,
                          margin: "auto",
                          borderRadius: 16,
                          overflow: "hidden",
                          aspectRatio: "16 / 9",
                          display: "grid",
                          gap: 2,
                          gridTemplateColumns:
                            block.medias.length === 1 ? "1fr" : "1fr 1fr",
                          gridTemplateRows:
                            block.medias.length <= 2 ? "1fr" : "1fr 1fr",
                        }}
                      >
                        {block.medias.map((media, index) => (
                          <div
                            key={media.id}
                            style={{
                              overflow: "hidden",
                              gridColumn:
                                block.medias!.length === 4
                                  ? index < 2
                                    ? 1
                                    : 2
                                  : undefined,
                              gridRow:
                                block.medias!.length === 3 && index === 0
                                  ? "1 / 3"
                                  : block.medias!.length === 4
                                    ? (index % 2) + 1
                                    : undefined,
                            }}
                          >
                            {media.type === "image" ? (
                              <img
                                style={{
                                  display: "block",
                                  width: "100%",
                                  height: "100%",
                                  objectFit: "contain",
                                  backgroundColor: "gray",
                                }}
                                src={media.url}
                              />
                            ) : media.type === "audio" ? (
                              <audio
                                controls
                                style={{
                                  display: "block",
                                  width: "100%",
                                  height: "100%",
                                  backgroundColor: "gray",
                                }}
                                src={media.url}
                              />
                            ) : media.type === "video" ? (
                              <video
                                controls
                                style={{
                                  display: "block",
                                  width: "100%",
                                  height: "100%",
                                  objectFit: "contain",
                                  backgroundColor: "gray",
                                }}
                                src={media.url}
                              />
                            ) : null}
                          </div>
                        ))}
                      </div>
                    )}
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
                        setForm((v) => ({
                          ...v,
                          percentage: Number(e.target.value),
                        }))
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
                        setForm((v) => ({
                          ...v,
                          condition: Number(e.target.value),
                        }))
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
                    <label>メディア（{mediaFiles.length} / 4件）</label>
                    <br />
                    {mediaFiles.length > 0 && (
                      <div style={{ display: "flex", flexWrap: "wrap", gap: 6, marginBottom: 6 }}>
                        {mediaFiles.map((file, i) => (
                          <div key={i} style={{ position: "relative" }}>
                            {file.type.startsWith("image/") ? (
                              <img
                                src={URL.createObjectURL(file)}
                                style={{ width: 64, height: 64, objectFit: "contain", borderRadius: 4, backgroundColor: "gray", display: "block" }}
                              />
                            ) : (
                              <div style={{ width: 64, height: 64, borderRadius: 4, background: "#eee", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 11, textAlign: "center", padding: 4, wordBreak: "break-all" }}>
                                {file.name}
                              </div>
                            )}
                            <button
                              type="button"
                              onClick={() => setMediaFiles((prev) => prev.filter((_, j) => j !== i))}
                              style={{ position: "absolute", top: -6, right: -6, width: 18, height: 18, borderRadius: "50%", background: "rgba(0,0,0,0.6)", color: "white", border: "none", cursor: "pointer", fontSize: 11, lineHeight: 1, padding: 0 }}
                            >
                              ×
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                    <input
                      ref={mediaInputRef}
                      type="file"
                      multiple
                      accept="image/jpeg,image/png,image/gif,audio/mpeg,audio/wav,audio/ogg,audio/aac,video/mp4,video/webm,video/quicktime"
                      style={{ display: "none" }}
                      onChange={(e) => {
                        if (!e.target.files) return;
                        const added = Array.from(e.target.files);
                        setMediaFiles((prev) => {
                          const merged = [...prev, ...added];
                          return merged.slice(0, 4);
                        });
                        e.target.value = "";
                      }}
                    />
                    <button
                      type="button"
                      disabled={mediaFiles.length >= 4}
                      onClick={() => mediaInputRef.current?.click()}
                    >
                      ファイルを追加
                    </button>
                  </div>
                  <div style={{ display: "flex", gap: 8 }}>
                    <button type="submit" disabled={isAddingBlock}>
                      {isAddingBlock ? "追加中..." : "追加する"}
                    </button>
                    <button
                      type="button"
                      onClick={() => setShowAddBlock(false)}
                    >
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
