import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  compatibilityDate: "2025-11-30",
  ssr: false,
  devtools: { enabled: true },
  css: ["./app/assets/css/main.css"],
  modules: ["@pinia/nuxt"],
  runtimeConfig: {
    // Private keys (only available server-side)
    apiSecret: "",
    // Public keys (exposed to client-side)
    public: {
      apiBase: "", // This can be overridden by NUXT_PUBLIC_API_BASE env var
    },
  },
  devServer: {
    host: "localhost",
    port: 3000,
  },
  vite: { plugins: [tailwindcss()] },
});
