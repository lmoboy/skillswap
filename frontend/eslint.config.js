import js from '@eslint/js';
import ts from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';
import svelte from 'eslint-plugin-svelte';
import svelteParser from 'svelte-eslint-parser';
import globals from 'globals';

export default [
	js.configs.recommended,
	{
		files: ['**/*.js'],
		languageOptions: {
			ecmaVersion: 2022,
			sourceType: 'module',
			globals: {
				...globals.browser,
				...globals.node
			}
		},
		rules: {
			'no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
			'no-console': 'off'
		}
	},
	{
		files: ['**/*.ts', '**/*.tsx'],
		languageOptions: {
			parser: tsParser,
			parserOptions: {
				ecmaVersion: 2022,
				sourceType: 'module',
				project: './tsconfig.json'
			},
			globals: {
				...globals.browser,
				...globals.node
			}
		},
		plugins: {
			'@typescript-eslint': ts
		},
		rules: {
			...ts.configs.recommended.rules,
			'@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
			'@typescript-eslint/no-explicit-any': 'warn'
		}
	},
	{
		files: ['**/*.svelte'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				parser: tsParser,
				extraFileExtensions: ['.svelte']
			},
			globals: {
				...globals.browser
			}
		},
		plugins: {
			svelte
		},
		rules: {
			...svelte.configs.recommended.rules,
			'svelte/no-unused-svelte-ignore': 'warn',
			'svelte/valid-compile': 'off',
			'no-unused-vars': ['warn', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
			'no-undef': 'off'
		}
	},
	{
		ignores: [
			'build/',
			'.svelte-kit/',
			'dist/',
			'node_modules/',
			'**/*.config.js',
			'**/*.config.ts'
		]
	}
];

