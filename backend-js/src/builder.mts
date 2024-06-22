import * as esbuild from 'esbuild'
import { exit } from 'process';
import * as paths from 'path';

export async function buildBackendModules(): Promise<() => Promise<void>> {
    console.log('Runnning backend on ', paths.resolve(__dirname, '../src'));
    const startTime = Date.now();
    const buildResult = await esbuild.build({
        entryPoints: ['./includer.js'],
        outdir: 'backend-modules-out',
        write: true,
        bundle: true,
        logLevel: 'silent',
        absWorkingDir: paths.resolve(__dirname, '../src'),
    });
    // await context.watch();
    (async () => {
            if (buildResult.errors.length > 0) {
        console.warn('Frontend build errors', buildResult.errors);
        exit(1);
    }
    if (buildResult.warnings.length > 0) {
        buildResult.warnings
            .filter((w) => w.id != 'direct-eval')
            .forEach((w) => console.warn('Warning', w));
    }
    })();

    return async () => {
        // await context.dispose();
        console.log('Stopped watching backend modules');
    };
}

export async function buildFrontend(): Promise<() => Promise<void>> {
    const startTime = Date.now();
    const context = await esbuild.context({
        entryPoints: ['./index.js'],
        outdir: 'assets/out',
        write: true,
        bundle: true,
        logLevel: 'silent',
        absWorkingDir: paths.resolve(__dirname, '../../frontend/'),
    });
    await context.watch();

    return async () => {
        await context.dispose();
        console.log('Stopped watching frontend');
    };

    // if (buildResult.errors.length > 0) {
    //     console.warn('Frontend build errors', buildResult.errors);
    //     exit(1);
    // }
    // if (buildResult.warnings.length > 0) {
    //     buildResult.warnings
    //         .filter((w) => w.id != 'direct-eval')
    //         .forEach((w) => console.warn('Warning', w));
    // }
    // console.log('Frontend built in ', Date.now() - startTime, 'ms');
}