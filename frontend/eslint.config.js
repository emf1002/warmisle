import pluginVue from 'eslint-plugin-vue'
import vueTsConfig from '@vue/eslint-config-typescript'
import eslintConfigPrettier from 'eslint-config-prettier'

export default [
  // Ignore patterns
  {
    ignores: ['node_modules/', 'dist/', '*.d.ts'],
  },
  // Vue 3 essential rules (flat config)
  ...pluginVue.configs['flat/essential'],
  // Vue + TypeScript
  ...vueTsConfig(),
  // Project-level overrides
  {
    rules: {
      // 项目大量使用 any，V1 不强制禁止
      '@typescript-eslint/no-explicit-any': 'off',
      // 页面组件用单词命名（如 Index、Login）是项目约定
      'vue/multi-word-component-names': 'off',
    },
  },
  // Prettier (must be last to override formatting rules)
  eslintConfigPrettier,
]
