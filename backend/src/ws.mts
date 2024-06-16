import http from 'http';
import { WebSocketServer, WebSocket } from 'ws';
import { WsMessage } from '../../common/ws.mjs';
import stream from "node:stream";

export class WsServer {
    readonly server: WebSocketServer;

    constructor(server?: WebSocketServer) {
        this.server = server ?? new WebSocketServer({ noServer: true });
        const self = this;
        this.server.on('connection', function connection(ws) {
            ws.on('error', console.error);

            ws.on('message', (event) => {
                console.log('received message %s', event);
                try {
                    const message = JSON.parse(event.toString()) as WsMessage;
                    self.sendMessage(ws, message.id, `I received this: ${message.data}`);
                } catch (e) {
                    console.warn(e);
                    self.sendMessage(ws, -1, 'Server Error');
                }
            });
            self.sendMessage(ws, -1, 'Confirming connection works');
        });
    }

    handleUpgrade(request: http.IncomingMessage, socket: stream.Duplex, head: Buffer) {
        const { pathname } = new URL(request.url as string, 'ws://url');
        if (pathname === '/ws') {
            this.server.handleUpgrade(request, socket, head, (ws) => {
                this.server.emit('connection', ws, request);
            });
        } else {
            socket.destroy();
        }
    }

    private sendMessage(ws: WebSocket, id: number, message: any, ok: boolean = true) {
        ws.send(JSON.stringify({ ok, id, data: message } as WsMessage));
    }
}

