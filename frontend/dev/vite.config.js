import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vite.dev/config/
// TODO: Implement this: https://vite.dev/guide/build.html#multi-page-app
export default defineConfig({
  plugins: [svelte()],
});
