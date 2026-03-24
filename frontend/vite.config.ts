/// <reference types="vitest" />
import react from '@vitejs/plugin-react';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), 'VITE_');

    return {
        base: '/',
        css: {
            preprocessorOptions: {
                scss: {
                    api: 'modern-compiler',
                },
            },
        },
        plugins: [
            react({ jsxImportSource: '@emotion/react' }),
        ],
        test: {
            environment: 'happy-dom',
            setupFiles: ['./src/test-setup.ts'],
        },
        server: {
            host: true,
            port: 3001,
            proxy: env.VITE_API_URL
                ? {
                    '/api': {
                        target: env.VITE_API_URL,
                        changeOrigin: true,
                        rewrite: (path) => path.replace(/^\/api/, ''),
                    },
                }
                : {},
        },
    };
});
