import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import dts from 'vite-plugin-dts';
import { visualizer } from 'rollup-plugin-visualizer';

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  // 加载环境变量 - 第二个参数是目录路径，不是文件名
  const env = loadEnv(mode, process.cwd(), '');
  const shouldAnalyze =
    process.argv.includes('--analyze') || env.ANALYZE === 'true';

  return {
    plugins: [
      react(),
      dts({
        insertTypesEntry: true,
        rollupTypes: true,
        tsconfigPath: './tsconfig.app.json',
        outDir: 'dist',
        include: ['src'],
      }),
      ...(command === 'build' && shouldAnalyze
        ? [
            visualizer({
              open: true, // 在默认浏览器中自动打开报告
              gzipSize: true, // 显示 gzip 格式下的包大小
              brotliSize: true, // 显示 brotli 格式下的包大小
              filename: 'dist/stats.html', // 分析图生成的文件名
            }),
          ]
        : []),
    ],
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
    // 优化构建配置
    build: {
      outDir: 'dist',
      lib: {
        entry: 'src/index.ts',
        formats: ['cjs', 'es'],
        fileName: (format) => `index.${format === 'cjs' ? 'js' : 'es.js'}`,
      },
      rollupOptions: {
        external: [
          '@ctzhian/ui',
          '@emotion/react',
          '@emotion/styled',
          '@mui/icons-material',
          'react',
          'react-dom',
          'react/jsx-runtime',
        ],
        output: {
          globals: {
            react: 'React',
            'react-dom': 'ReactDOM',
          },
        },
      },
      sourcemap: true,
    },
  };
});
