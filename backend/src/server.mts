import http from 'http';
import serveStatic from 'serve-static';
import finalHandler from 'finalhandler';
import { WsServer } from './ws.mjs';

const serveFile = serveStatic('../frontend', { index: ['index.html', 'index.htm'] });

const server = http.createServer((req, res) => {
    if (req.url === '/') {
        res.setHeader('Content-Type', 'text/html');
        res.write('<h1>Hello world</h1>');
        res.end();
        return;
    }
    serveFile(req, res, finalHandler(req, res));
});

const ws = new WsServer();

server.on('upgrade', (req, socket, head) => ws.handleUpgrade(req, socket, head));
server.listen();
console.log(`Server running at http://localhost:${server.address().port}/`);
