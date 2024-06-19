import http from 'http';
import serveStatic from 'serve-static';
import finalHandler from 'finalhandler';
import { WsServer } from './ws.mjs';
import { buildFrontend } from './builder.mjs';
import { run } from './modules/files.mjs';

async function main() {
    // TODO remove this, but make sure the files.mjs runs!
    run();
    
    const stopFeWatcher = await buildFrontend();

    const serveFile = serveStatic('../frontend/assets', { index: ['index.html'] });

    const server = http.createServer((req, res) => {
        serveFile(req, res, finalHandler(req, res));
    });

    const ws = new WsServer();

    server.on('upgrade', (req, socket, head) => ws.handleUpgrade(req, socket, head));
    server.on('close', stopFeWatcher);
    server.listen(8001);
    console.log(`Server running at http://localhost:${server.address().port}/`);
}

main();
