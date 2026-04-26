import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import generouted from "@generouted/react-router/plugin";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), generouted()],
  server: {
    host: true,
    port: 5173,
    // Windowsのバインドマウントではinotifyイベントがコンテナに伝わらないためポーリングで代替
    watch: {
      usePolling: true,
    },
    proxy: {
      "/api": {
        target: "http://backend:8000",
        changeOrigin: true,
      },
    },
  },
});
