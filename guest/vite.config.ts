import { fileURLToPath, URL } from "node:url";
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import type { ImportMetaEnv } from "./env";

export default defineConfig(({ mode }) => {
  let env: Record<keyof ImportMetaEnv, string> = loadEnv(mode, process.cwd());

  const serverUrl = env.VITE_SERVER_URL;
  console.log("serverUrl",serverUrl);
  return {
    css: {
      preprocessorOptions: {
        scss: { api: 'modern-compiler' },
      }
    },
    plugins: [vue()],
    envDir: "./",
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    server: {
      host: "0.0.0.0",
      port: 80,
      proxy: {
        "/api": {
          target: serverUrl,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ""), // 添加这行，如果后端接口不需要 /api 前缀
        },
        "/uploads": {
          target: serverUrl,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ""), // 添加这行，如果后端接口不需要 /api 前缀
        },
      },
    },
  };
});
