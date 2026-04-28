import { Outlet } from "react-router";
import { useGetMe } from "../api";

const Layout = () => {
  const { data, isSuccess } = useGetMe();

  return (
    <>
      <div
        style={{
          height: 64,
          borderBottom: "1px solid black",
          position: "sticky",
          backgroundColor: "white",
          top: 0,
          display: "flex",
          alignItems: "center",
          padding: "0 16px",
          justifyContent: "space-between",
        }}
      >
        <h1 style={{ margin: 0 }}>
          <a href="/" style={{ textDecoration: "none", color: "black" }}>
            TSUMIKI
          </a>
        </h1>
        {isSuccess && data ? (
          <div
            style={{
              height: "85%",
              display: "flex",
              alignItems: "center",
              gap: 8,
            }}
          >
            <a href="/upload" style={{ height: "100%" }}>
              <button style={{ height: "100%" }}>投稿する</button>
            </a>
            <div
              style={{
                height: "100%",
                boxSizing: "border-box",
                border: "1px solid black",
                padding: "0 8px",
                borderRadius: 4,
                cursor: "pointer",
              }}
            >
              <a
                href="/my"
                style={{
                  height: "100%",
                  display: "flex",
                  alignItems: "center",
                  gap: 8,
                  textDecoration: "none",
                  color: "black",
                }}
              >
                <img
                  style={{ height: "80%", borderRadius: 9999 }}
                  src={data.avatarUrl}
                />
                <span>{data.name}</span>
              </a>
            </div>
          </div>
        ) : (
          <a href="/login" style={{ height: "85%" }}>
            <button style={{ height: "100%" }}>ログインする</button>
          </a>
        )}
      </div>
      <Outlet />
    </>
  );
};

export default Layout;
