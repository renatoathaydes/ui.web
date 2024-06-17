import http from 'http';
import serveStatic from 'serve-static';
import finalHandler from 'finalhandler';
import { WsServer } from './ws.mjs';
import { buildFrontend } from './builder.mjs';

async function main() {
    const stopWatcher = await buildFrontend();

    const serveFile = serveStatic('../frontend/assets', { index: ['index.html'] });

    const server = http.createServer((req, res) => {
        serveFile(req, res, finalHandler(req, res));
    });

    const ws = new WsServer();

    server.on('upgrade', (req, socket, head) => ws.handleUpgrade(req, socket, head));
    server.on('close', stopWatcher);
    server.listen();
    console.log(`Server running at http://localhost:${server.address().port}/`);
}

main();
