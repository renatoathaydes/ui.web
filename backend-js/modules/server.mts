import http from 'http';
import serveStatic from 'serve-static';
import finalHandler from 'finalhandler';
import { BackendCommandRunner } from './command.mts';
import { jsonFromRequest } from './request.mts';

export async function start() {
    const runner = new BackendCommandRunner;
    const serveFile = serveStatic('./assets', { index: ['index.html'] });

    const server = http.createServer(async (req, res) => {
        const url = new URL(`http://${process.env.HOST ?? 'localhost'}${req.url ?? '/'}`);
        if (req.method === 'POST' && url.pathname === '/command') {
            try {
                const json = await jsonFromRequest(req);
                console.log(`Received POST: '${json}'`);
                const result = await runner.runScript(json);
                res.end(JSON.stringify(result));
            } catch (error) {
                console.warn(error);
                res.end(JSON.stringify({ error: error.toString() }));
            }
        } else {
            serveFile(req, res, finalHandler(req, res));
        }
    });
    server.listen(8001);
    console.log(`JS backend running at http://localhost:${server.address().port}/`);
}
