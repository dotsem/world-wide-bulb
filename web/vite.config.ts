import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';
import process from 'node:process';
import path from 'node:path';

export default defineConfig(({ mode }) => {
	const envDir = path.resolve(process.cwd(), '..');
	const env = loadEnv(mode, envDir, '');

	return {
		envDir,
		plugins: [
			tailwindcss(),
			sveltekit({
				compilerOptions: {
					runes: ({ filename }) =>
						filename.split(/[/\\]/).includes('node_modules') ? undefined : true
				},
				adapter: adapter()
			})
		],
		server: {
			port: env.FRONTEND_PORT ? Number(env.FRONTEND_PORT) : 5173
		}
	};
});
