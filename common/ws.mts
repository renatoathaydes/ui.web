export type WsMessage = {
    ok: boolean,
    id: number,
    data: any,
    type: WsMessageType,
};

export type WsMessageType = 'rpc' | 'response';
