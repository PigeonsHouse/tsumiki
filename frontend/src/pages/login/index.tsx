import { useCallback } from "react";

const Login = () => {
  const onClick = useCallback(() => {
    location.href = "/api/v1/auth/discord?mode=redirect";
  }, []);

  return (
    <center>
      <p>
        TSUMIKIのログインは、DiscordのOAuthで認証しています。
        <br />
        現在は認証可能なDiscordサーバに制限をかけています。サーバ管理者から案内を受けていない場合は、ゲストのままTSUMIKIをお楽しみください。
      </p>
      <button onClick={onClick}>Discord認証でログインする</button>
    </center>
  );
};

export default Login;
