import Watcher from "watcher";
function ignorePaths(path) {
  return path.startsWith(".git");
}
async function startWacher(dir) {
  const watcher = new Watcher(dir, {
    recursive: true,
    ignoreInitial: true,
    ignore: ignorePaths
  });
  watcher.on("all", (event, path) => {
    console.log(`event type is: ${event}, file='${path}'`);
  });
  return watcher;
}
startWacher("./src/modules").then((w) => {
  setTimeout(w.close, 3e4);
});
