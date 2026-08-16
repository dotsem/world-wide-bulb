import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import process from 'node:process';
import path from 'node:path';

export default defineConfig(({ mode }) => {
	const envDir = path.resolve(process.cwd(), '..');
	const env = loadEnv(mode, envDir, '');

	return {
		envDir,
		plugins: [tailwindcss(), sveltekit()],
		server: {
			port: env.FRONTEND_PORT ? Number(env.FRONTEND_PORT) : 5001
		}
	};
});
