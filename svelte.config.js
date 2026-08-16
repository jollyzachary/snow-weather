import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      pages: 'backend/static/dist',
      assets: 'backend/static/dist',
      fallback: 'index.html',
      precompress: true,
      strict: true
    })
  }
};

export default config;
