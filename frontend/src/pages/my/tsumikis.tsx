import { useDeleteTsumiki, useGetMyTsumikis } from "../../api";

const MyTsumikis = () => {
  const { data } = useGetMyTsumikis();
  const { mutate: deleteTsumiki } = useDeleteTsumiki();

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>自分の積み木一覧</h1>
      <div style={{ maxWidth: 1024, margin: "auto" }}>
        <div style={{ display: "flex", flexWrap: "wrap", gap: 8 }}>
          {data?.map((tsumiki) => (
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
              {tsumiki.work && (
                <a href={`/works/${tsumiki.work.id}`}>
                  作品：{tsumiki.work.title}
                </a>
              )}
              <small style={{ display: "flex", gap: 8 }}>
                {tsumiki.percentage !== null && (
                  <span>進捗度: {tsumiki.percentage}%</span>
                )}
                <span>いいね: {tsumiki.favorite.totalFavoriteCount}</span>
              </small>
              <div style={{ display: "flex", justifyContent: "flex-end" }}>
                <a
                  style={{ textDecoration: "none" }}
                  href={`/tsumikis/${tsumiki.id}/edit`}
                >
                  <button>✏️</button>
                </a>
                <button
                  onClick={(e) => {
                    e.preventDefault();
                    if (confirm(`${tsumiki.title} を削除してもよいですか？`)) {
                      deleteTsumiki(tsumiki.id, {
                        onError: () => alert("削除に失敗しました"),
                      });
                    }
                  }}
                >
                  🗑️
                </button>
              </div>
              <small>{tsumiki.createdAt.toISOString()}</small>
            </a>
          ))}
        </div>
      </div>
    </div>
  );
};

export default MyTsumikis;
