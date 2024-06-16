import { WsMessage, BackendCall } from "../../common/ws.mts";

type Resolver = {
    resolve: (result: any) => void;
    reject: (reason: any) => void;
};

const ws = new WebSocket(new URL(`ws://${location.hostname}:${location.port}/ws`));

let id = 0;

ws.onopen = () => {
    send('Hello Server');
};

ws.onerror = (event) => {
    console.warn('WS error', event);
}

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

export function callBackend(message: BackendCall): Promise<any> {
    return send(message, true);
}

function send(message: any, needsAnswer: boolean = false): Promise<any> {
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
    try {
        ws.send(JSON.stringify({ ok: true, id: mid, data: message } as WsMessage));
    } catch (e) {
        responses.delete(mid);
        throw e;
    }
    return answer;
}