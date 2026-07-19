import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss(), react()],
  clearScreen: false,
  base: "/admin/",
  server: {
    port: 5174,
    strictPort: true,
  },
});
