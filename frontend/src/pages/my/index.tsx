import { useNavigate } from "react-router";
import { useGetMe, useLogout } from "../../api";

const MyPage = () => {
  const { data } = useGetMe();
  const navigate = useNavigate();
  const { mutate: logout } = useLogout();

  const handleLogout = () => {
    logout(undefined, {
      onSuccess: () => navigate("/"),
    });
  };

  return (
    <div>
      <h1 style={{ marginTop: 0, textAlign: "center" }}>マイページTOP</h1>
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
            <a href="/my/tsumikis">積み木を見る</a>
            <a href="/my/works">関わった作品を見る</a>
          </div>
          <div style={{ display: "flex", justifyContent: "center", marginTop: 24 }}>
            <button onClick={handleLogout}>ログアウト</button>
          </div>
        </>
      )}
    </div>
  );
};

export default MyPage;
