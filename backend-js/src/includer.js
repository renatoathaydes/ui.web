// esbuild entrypoint: include everything in the modules directory
const kind = 'files';
require('./modules/' + kind + '.mts');
