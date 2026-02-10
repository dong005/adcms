import { defineConfig } from '@vben/vite-config';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      server: {
        proxy: {
          '/api': {
            changeOrigin: true,
            // 代理到Gin后端
            target: 'http://localhost:8004',
            ws: true,
          },
          '/uploads': {
            changeOrigin: true,
            target: 'http://localhost:8004',
          },
        },
      },
    },
  };
});
