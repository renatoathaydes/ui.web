// esbuild entrypoint: include everything in the modules directory
const kind = 'commands';
require('./modules/' + kind + '.mts');
import {} from './entrypoint';
