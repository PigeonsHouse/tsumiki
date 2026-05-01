import { useMemo } from "react";
import { useParams } from "react-router";
import { useGetUserTsumikis } from "../../../api";
import NotFound from "../../../components/NotFound";

const UserTsumikis = () => {
  const { userId: userIdRaw } = useParams();
  const isValidId = useMemo(
    () =>
      userIdRaw !== "" &&
      userIdRaw !== undefined &&
      !Number.isNaN(Number(userIdRaw)),
    [userIdRaw],
  );
  const userId = useMemo(() => Number(userIdRaw), [userIdRaw]);
  const { data, isError } = useGetUserTsumikis(userId, isValidId);

  return !isValidId || isError ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>ユーザの積み木一覧</h1>
      <div style={{ maxWidth: 1024, margin: "auto" }}>
        <div style={{ display: "flex", flexWrap: "wrap", gap: 8 }}>
          {data?.map((tsumiki) => (
            <a
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
              <small>{tsumiki.createdAt.toISOString()}</small>
            </a>
          ))}
        </div>
      </div>
    </div>
  );
};

export default UserTsumikis;
