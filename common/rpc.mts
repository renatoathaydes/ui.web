export type MethodCall = {
    name: string,
    arg?: any,
};

export const methodCallType = 'rpc';

/// Add method handlers to this Map to allow the frontend to call methods
/// on the backend.
export const methodHandlers: Map<string, (arg: any) => any> = new Map();
