import { CommandResponse } from "../../common/command.mts";
import * as fs from 'node:fs/promises';

export async function open(name: string): Promise<CommandResponse> {
    console.log('Attempting to open: ' + name);
    const value = await fs.readFile(name, { encoding: 'utf-8' });
    return {
        value, 
        feCmd: `editor.openEditor("${name}", value)`
    };
}
