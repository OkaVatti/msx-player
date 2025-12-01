import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2025-11-30",
  ssr: false,
  devtools: { enabled: true },
  css: ["./app/assets/css/main.css"],
  modules: ["@pinia/nuxt"],
  devServer: {
    host: "localhost",
    port: 3000,
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || "http://localhost:1323",
    },
  },
  nitro: {
    preset: "deno",
    devProxy: {
      "/api": {
        target: "http://localhost:1323/api",
        changeOrigin: true,
      },
      "/music": {
        target: "http://localhost:1323/music",
        changeOrigin: true,
      },
    },
  },

  // Better Deno compatibility
  build: {
    transpile: [],
  },
  vite: {
    clearScreen: false,
    server: {
      watch: {
        usePolling: false,
      },
    },
    plugins: [tailwindcss()],
  },
  typescript: {
    typeCheck: false, // Disable during dev to avoid Deno issues
    strict: false, // Temporarily disable strict mode
    shim: false,
  },
});
