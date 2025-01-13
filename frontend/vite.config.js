import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vite.dev/config/
// TODO: Implement this: https://vite.dev/guide/build.html#multi-page-app
export default defineConfig({
  plugins: [svelte()],
  build: {
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
        manualChunks: undefined,
        entryFileNames: 'app.js',
        assetFileNames: (assetInfo) => {
          if (assetInfo.names === 'style.css') {
            return 'app.css';
          }
          return assetInfo.name;
        },
      }
    }
  }
});
