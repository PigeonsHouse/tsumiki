import { useGetWorks } from "../../api";

const WorkList = () => {
  const { data } = useGetWorks();
  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>作品一覧</h1>
      <div style={{ width: 1024, margin: "auto" }}>
        <div style={{ display: "flex", flexWrap: "wrap", gap: 8 }}>
          {data?.map((tsumiki) => (
            <a
              href={`/works/${tsumiki.id}`}
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
                href={`/users/${tsumiki.owner.id}`}
              >
                <img
                  style={{ width: 40, height: 40, borderRadius: 9999 }}
                  src={tsumiki.owner.avatarUrl}
                />
                <span>{tsumiki.owner.name}</span>
              </a>
              <small>{tsumiki.createdAt.toISOString()}</small>
            </a>
          ))}
        </div>
      </div>
    </div>
  );
};

export default WorkList;
