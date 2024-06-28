import * as foo from './foo.mts';

export function hello() {
    return 'index.hello' + foo.hello();
}
