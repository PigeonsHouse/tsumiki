import { Outlet } from "react-router";

const Layout = () => {
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
        }}
      >
        <h1 style={{ margin: 0 }}>
          <a href="/" style={{ textDecoration: "none", color: "black" }}>
            TSUMIKI
          </a>
        </h1>
      </div>
      <Outlet />
    </>
  );
};

export default Layout;
