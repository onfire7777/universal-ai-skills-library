import adapter from "@sveltejs/adapter-auto";
import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss(), sveltekit({ adapter: adapter(), alias: { "$lib": "src/lib" } })],
  ssr: {
    noExternal: ["@openuidev/svelte-lang", "@openuidev/lang-core"],
  },
});
