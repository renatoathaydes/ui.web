import { CommandResponse } from "../../common/command.mts";
import * as fs from 'node:fs/promises';

export async function openFile(name: string): Promise<CommandResponse> {
    console.log('Attempting to open file: ' + name);
    const value = await fs.readFile(name, { encoding: 'utf-8' });
    return {
        value, 
        feCmd: `openEditor("${name}", value)`
    };
}
