const Index = () => {
  return (
    <>
      <h1 style={{ textAlign: "center" }}>TOP</h1>
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          gap: 16,
        }}
      >
        <a href="/tsumikis">積み木一覧</a>
        <a href="/works">作品一覧</a>
      </div>
    </>
  );
};

export default Index;
