export type BackendCall = {
    method: string,
    args: any[],
};

export type WsMessage = {
    'ok': boolean,
    'id': number,
    'data': any,
};
