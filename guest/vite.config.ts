import { fileURLToPath, URL } from "node:url";
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import type { ImportMetaEnv } from "./env";

export default defineConfig(({ mode }) => {
  let env: Record<keyof ImportMetaEnv, string> = loadEnv(mode, process.cwd());

  const serverUrl = env.VITE_SERVER_URL;
  console.log(serverUrl);
  return {
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
    },
  };
});
