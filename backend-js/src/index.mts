import { CommandResponse } from '../../common/command.mjs';

// esbuild entrypoint: include everything in the modules directory
const kind = 'files';
require('./modules/' + kind + '.mts');

export async function runScript(js: string): Promise<CommandResponse> {
    console.log('Runing js', js);
    try {
        const value = await eval(js);
        console.log(`Success: ${value}`);
        return asCommandResponse(value);
    } catch (e) {
        console.warn(e);
        return { error: e.toString() };
    }
}

function asCommandResponse(value: any): CommandResponse {
    if (typeof value === 'object') {
        if ("error" in value) {
            return value as CommandResponse;
        }
        if ("value" in value) {
            return value as CommandResponse;
        }
    }
    return { value };
}
