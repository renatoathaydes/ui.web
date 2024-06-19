import { methodHandlers } from '../../../common/rpc.mjs';

console.warn('Installing openFile function');

methodHandlers.set('openFile', (name: string) => {
    console.log('openFile called', name);
    return `Going to open file ${name}`;
});

export function run() {
    console.warn('files run()');
}
