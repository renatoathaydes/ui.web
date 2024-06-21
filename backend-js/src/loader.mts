import Watcher from 'watcher';

function ignorePaths(path: string): boolean {
    return path.startsWith('.git');
}

async function startWacher(dir: string): Promise<Watcher> {
    const watcher = new Watcher(dir, {
        recursive: true,
        ignoreInitial: true,
        ignore: ignorePaths,
    });
    watcher.on('all', (event, path) => {
        console.log(`event type is: ${event}, file='${path}'`);
    });
    return watcher;
}

startWacher('./src/modules').then((w) => {
    setTimeout(w.close, 30_000);
});
