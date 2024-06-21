import * as hello from './hello.mjs';
import './utils.js';

function lookupMod(mod) {
    if (!mod) {
        console.log('No mod provided, using global');
        return global;
    }
    switch(mod) {
        case "hello": return hello;
        default: throw Error('Module does not exist: ' + mod);
    }
}

function joe() {
    return 'JOE';
}

function call(name, mod) {
    console.log('calling ', name, ', mod: ', mod);
    const namespace = lookupMod(mod);
    const value = namespace[name];
    console.log('function: ', value);
    if (typeof(value) === 'function') {
        return value();
    }
    return ('no such function');
}

console.log(call('hello', 'hello'));
console.log(call('joe'));
console.log(call('util'));

