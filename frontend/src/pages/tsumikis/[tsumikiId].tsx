import { useParams } from "react-router";
import { useGetBlocks, useGetTsumiki } from "../../api";
import { useMemo } from "react";

const TsumikiDetail = () => {
  const { tsumikiId: tsumikiIdRaw } = useParams();
  const isValidId = useMemo(
    () =>
      tsumikiIdRaw !== "" &&
      tsumikiIdRaw !== undefined &&
      !Number.isNaN(Number(tsumikiIdRaw)),
    [tsumikiIdRaw],
  );
  const tsumikiId = useMemo(() => Number(tsumikiIdRaw), [tsumikiIdRaw]);
  const { data: tsumiki } = useGetTsumiki(tsumikiId, isValidId);
  const { data: blocks } = useGetBlocks(tsumikiId, isValidId);

  return (
    <div>
      <h1 style={{ marginTop: 0 }}>積み木詳細</h1>
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
          {blocks.map((block) => (
            <>
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
            </>
          ))}
        </div>
      )}
    </div>
  );
};

export default TsumikiDetail;
