import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },
    server: {
      proxy: {
        '^/api/': env.VITE_API_BASE_URL || 'http://localhost:8080/',
      },
      host: '0.0.0.0',
      port: 3300,
    },
    // 手动配置 Monaco Editor 支持
    define: {
      // 禁用 Monaco Editor 从 CDN 加载
      'process.env.REACT_APP_MONACO_CDN': JSON.stringify('false'),
    },
    // 优化构建配置
    build: {
      emptyOutDir: false, // 不清空输出目录，保留 TypeScript 生成的声明文件
      lib: {
        entry: path.resolve(__dirname, 'src/index.ts'),
        name: 'MonkeyCodeUI',
        formats: ['es', 'cjs'],
        fileName: (format) => `index.${format === 'es' ? 'es.js' : 'js'}`,
      },
      rollupOptions: {
        external: [
          'react',
          'react-dom',
          'react-hook-form',
          '@mui/material',
          '@mui/icons-material',
          '@c-x/ui',
          'axios'
        ],
        output: {
          globals: {
            react: 'React',
            'react-dom': 'ReactDOM',
            'react-hook-form': 'ReactHookForm',
            '@mui/material': 'Material',
            '@mui/icons-material': 'MaterialIcons',
            '@c-x/ui': 'CXUi',
            axios: 'axios'
          },
        },
      },
    },
    // 确保 Monaco Editor 被正确优化
    // optimizeDeps: {
    //   include: ['monaco-editor', '@monaco-editor/react'],
    // },
    // // 处理 worker 文件
    // worker: {
    //   format: 'es',
    // },
    // 确保 Monaco Editor workers 能正确加载
    // assetsInclude: ['**/*.worker.js'],
  };
});
