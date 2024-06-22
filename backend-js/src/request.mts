import { IncomingMessage } from "http";

export function jsonFromRequest(req: IncomingMessage): Promise<string> {
    let body = '';
    req.on('data', (chunk) => {
        body += chunk.toString();
    });
    return new Promise((resolve, reject) => {
        req.on('end', () => resolve(body));
        req.on('error', reject);
    });
}
