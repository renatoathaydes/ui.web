import { createCommandInput } from './modules/commands.mts';
import { callBackend } from './modules/ws.mts';

createCommandInput();

setTimeout(async () => {
    const result = await callBackend({method: 'foo', args: ['hello from the frontend']});
    console.log('backend call returned: ', result);
}, 1000);
