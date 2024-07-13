import { exec } from 'node:child_process';
import { fileURLToPath } from 'url';
import path from "node:path";
import * as api from '../../common/command.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const projectDirectoryUrl = new URL('../../frontend', path.dirname(import.meta.url));

export function install(lib: string): Promise<api.CommandResponse> {
    return new Promise((resolve) => {
        exec('npm install ' + lib, (err, stdout, stderr) => {
            if (err) {
                resolve({ "error": `${err}` })
                return;
            }
            console.log(stdout);
            console.log(stderr);
            updateImportMap().then((_) => {
                resolve({ "value": "Installation succeeded, re-generated import map!" });
            }, (err) => {
                resolve({ "error": `${err}` });
            })
        });
    });
}

async function updateImportMap() {
     // TODO
}
