import { WsMessage, WsMessageType } from "../../common/ws.mts";
// import rpc from 'json-rpc-protocol';

type Resolver = {
    resolve: (result: any) => void;
    reject: (reason: any) => void;
};

const ws = new WebSocket(new URL(`ws://${location.hostname}:${location.port}/ws`));

let id = 0;

ws.onopen = () => {
    console.log('Websocket connection opened!');
    send('Hello Server', false, 'rpc');
};

ws.onerror = (event) => {
    console.warn('WS error', event);
}

ws.onclose = (event) => {
    console.warn('WS closed', event);
};

const responses = new Map<number, Resolver>();

ws.onmessage = (event) => {
    const message = JSON.parse(event.data as string) as WsMessage;
    const resolver = responses.get(message.id);
    responses.delete(message.id);
    console.log('RESOLVER FOR:', message, resolver);
    if (resolver) {
        if (message.ok) {
            resolver.resolve(message.data);
        } else {
            resolver.reject(message.data);
        }
    }
};

export function callBackend(message: string): Promise<any> {
    return send(message, true, 'rpc');
}

function send(message: any, needsAnswer: boolean, type: WsMessageType): Promise<any> {
    if (ws.readyState !== ws.OPEN) {
        throw new Error('Cannot send message, websocket state: ' + ws.readyState);
    }
    const mid = needsAnswer ? id++ : -1;
    let answer: Promise<any>;
    if (needsAnswer) {
        const resolver: Resolver = { resolve: () => { }, reject: () => { } };
        answer = new Promise((resolve, reject) => {
            resolver.resolve = resolve;
            resolver.reject = reject;
        });
        responses.set(mid, resolver);
    } else {
        answer = Promise.resolve();
    }
    const msg: WsMessage = { type, ok: true, id: mid, data: message };
    console.log('Sending WsMessage: ' + JSON.stringify(msg));
    try {
        ws.send(JSON.stringify(msg));
    } catch (e) {
        responses.delete(mid);
        throw e;
    }
    return answer;
}