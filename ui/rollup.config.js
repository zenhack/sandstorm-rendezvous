import node_resolve from 'rollup-plugin-node-resolve'

export default {
	input: './index.js',
	output: {
		file: './vnc-bundle.js',
		format: 'iife',
		name: 'spread',
	},
	plugins: [
		node_resolve({browser: true}),
	],
}
