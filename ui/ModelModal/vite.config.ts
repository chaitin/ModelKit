import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import dts from 'vite-plugin-dts';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '');

  return {
    plugins: [react(), dts({
      insertTypesEntry: true,
      rollupTypes: true,
      tsconfigPath: './tsconfig.app.json',
      outDir: 'dist',
      include: ['src'],
    })],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },
    server: {
      proxy: {
        '^/api/': env.VITE_BASE_URL_URL || 'http://localhost:8080/',
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
      outDir: 'dist',
      lib: {
        entry: 'src/index.ts',
        formats: ['cjs', 'es'],
        fileName: (format) => `index.${format === 'cjs' ? 'js' : 'es.js'}`
      },
      rollupOptions: {
        external: ['react', 'react-dom'],
        output: {
          globals: {
            react: 'React',
            'react-dom': 'ReactDOM'
          }
        }
      },
      sourcemap: true,
    },
    // 确保 Monaco Editor 被正确优化
    optimizeDeps: {
      include: ['monaco-editor', '@monaco-editor/react'],
    },
    // // 处理 worker 文件
    // worker: {
    //   format: 'es',
    // },
    // 确保 Monaco Editor workers 能正确加载
    // assetsInclude: ['**/*.worker.js'],
  };
});
