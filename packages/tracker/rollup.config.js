import terser from '@rollup/plugin-terser';
import typescript from '@rollup/plugin-typescript';

/**
 * @type {import('rollup').RollupOptions}
 */
export default {
    input: 'src/main.ts',
    output: [
        {
            file: 'dist/tracker.js',
            format: 'iife',
            plugins: [
                terser(),
            ],
        }    
    ],
    plugins: [
        typescript(),
    ]
}
