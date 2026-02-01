import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
	server: {
		host: "0.0.0.0",
		cors: true,
	},
	plugins: [svelte(), tailwindcss()],
	build: {
		manifest: true,
	},
});
