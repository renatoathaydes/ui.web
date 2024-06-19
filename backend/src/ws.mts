import http from 'http';
import { WebSocketServer, WebSocket } from 'ws';
import { WsMessage } from '../../common/ws.mjs';
import { MethodCall, methodHandlers, methodCallType } from '../../common/rpc.mjs';
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
                let message: WsMessage;
                try {
                    message = JSON.parse(event.toString()) as WsMessage;
                } catch (e) {
                    console.warn(e);
                    self.sendMessage(ws, -1, 'Server Error', false);
                    return;
                }
                try {
                    let result: any, ok: boolean = false;
                    if (message.type === methodCallType) {
                        const rpcCall = message.data as MethodCall;
                        const method = methodHandlers.get(rpcCall.name);
                        if (method) {
                            result = method(rpcCall.arg);
                            ok = true;
                        } else {
                            result = 'Unknown method: ' + rpcCall.name;
                        }
                    } else {
                        result = 'Unknown message type: ' + message.type;
                    }    
                    self.sendMessage(ws, message.id, result, ok);
                } catch (e) {
                    console.warn(e);
                    self.sendMessage(ws, message.id, 'Server Error', false);
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

