import { useMemo } from "react";
import { useParams } from "react-router";
import { useGetWork, useGetWorkTsumikis, useDeleteWork } from "../../../api";
import NotFound from "../../../components/NotFound";

const WorkDetail = () => {
  const { workId: workIdRaw } = useParams();
  const isValidId = useMemo(
    () =>
      workIdRaw !== "" &&
      workIdRaw !== undefined &&
      !Number.isNaN(Number(workIdRaw)),
    [workIdRaw],
  );
  const workId = useMemo(() => Number(workIdRaw), [workIdRaw]);
  const { data: work, isError } = useGetWork(workId, isValidId);
  const { data: tsumikis } = useGetWorkTsumikis(workId);
  const { mutateAsync: deleteWork } = useDeleteWork(workId);

  const handleDelete = async () => {
    if (!confirm("この作品を削除しますか？")) return;
    await deleteWork();
    location.href = "/works";
  };

  return !isValidId || isError ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>作品詳細</h1>
      {work && (
        <div style={{ maxWidth: 1024, margin: "auto" }}>
          <img
            style={{
              aspectRatio: "16 / 9",
              width: "100%",
              objectFit: "contain",
              backgroundColor: "gray",
            }}
            src={work.thumbnailUrl || ""}
          />
          <h2>{work.title}</h2>
          {work.description && <p>{work.description}</p>}
          <a
            style={{ display: "flex", alignItems: "center", gap: 8 }}
            href={`/users/${work.owner.id}`}
          >
            <img
              style={{ width: 40, height: 40, borderRadius: 9999 }}
              src={work.owner.avatarUrl}
            />
            <span>{work.owner.name}</span>
          </a>
          <small>{work.createdAt.toISOString()}</small>
          {work.isOwner && (
            <div style={{ display: "flex", gap: 8, marginTop: 12 }}>
              <a href={`/works/${workId}/edit`}>
                <button type="button">編集</button>
              </a>
              <button type="button" onClick={handleDelete}>
                削除
              </button>
              <a href={`/tsumikis/new?workId=${workId}`}>
                <button type="button">積み木を追加</button>
              </a>
            </div>
          )}
          <hr />
          <div style={{ display: "flex", flexWrap: "wrap", gap: 8 }}>
            {tsumikis?.map((tsumiki) => (
              <a
                key={tsumiki.id}
                href={`/tsumikis/${tsumiki.id}`}
                style={{
                  border: "1px solid black",
                  display: "flex",
                  flexDirection: "column",
                  width: 280,
                  gap: 16,
                  padding: 8,
                  borderRadius: 4,
                  textDecoration: "none",
                  color: "unset",
                  cursor: "pointer",
                }}
              >
                <h2 style={{ margin: 0 }}>{tsumiki.title}</h2>
                <img
                  style={{
                    aspectRatio: "16 / 9",
                    width: "100%",
                    objectFit: "contain",
                    backgroundColor: "gray",
                  }}
                  src={tsumiki.thumbnailUrl || ""}
                />
                {tsumiki.percentage != null && (
                  <div>進捗度: {tsumiki.percentage}%</div>
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
              </a>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default WorkDetail;
