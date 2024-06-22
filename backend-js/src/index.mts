import { CommandResponse } from '../../common/command.mjs';

// esbuild entrypoint: include everything in the modules directory
const kind = 'files';
require('./modules/' + kind + '.mts');

export async function runScript(js: string): Promise<CommandResponse> {
    console.log('Runing js', js);
    try {
        const value = await eval(js);
        console.log(`Success: ${value}`);
        return { value };
    } catch (e) {
        console.warn(e);
        return { error: e.toString() };
    }
}
