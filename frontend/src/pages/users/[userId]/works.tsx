import { useMemo } from "react";
import { useParams } from "react-router";

const MyWorks = () => {
  const { userId: userIdRaw } = useParams();
  const isValidId = useMemo(
    () =>
      userIdRaw !== "" &&
      userIdRaw !== undefined &&
      !Number.isNaN(Number(userIdRaw)),
    [userIdRaw],
  );
  const userId = useMemo(() => Number(userIdRaw), [userIdRaw]);

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>
        ユーザの関わった作品一覧
      </h1>
    </div>
  );
};

export default MyWorks;
