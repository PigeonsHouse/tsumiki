import { useMemo } from "react";
import { useParams } from "react-router";
import { useGetUser } from "../../../api";
import NotFound from "../../../components/NotFound";

const UserPage = () => {
  const { userId: userIdRaw } = useParams();
  const isValidId = useMemo(
    () =>
      userIdRaw !== "" &&
      userIdRaw !== undefined &&
      !Number.isNaN(Number(userIdRaw)),
    [userIdRaw]
  );
  const userId = useMemo(() => Number(userIdRaw), [userIdRaw]);
  const { data, isError } = useGetUser(userId, isValidId);

  return !isValidId || isError ? (
    <NotFound />
  ) : (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>ユーザページTOP</h1>
      {data && (
        <>
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
            }}
          >
            <img style={{ borderRadius: 9999 }} src={data.avatarUrl} />
            <h2>{data.name}</h2>
          </div>
          <div style={{ display: "flex", justifyContent: "center", gap: 16 }}>
            <a href={`/users/${userId}/tsumikis`}>積み木を見る</a>
            <a href={`/users/${userId}/works`}>関わった作品を見る</a>
          </div>
        </>
      )}
    </div>
  );
};

export default UserPage;
