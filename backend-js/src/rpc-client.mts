import Peer from 'json-rpc-peer';

const peer = Peer();

export function invoke(request: any): Promise<any> {
    return peer.request(request);
}
