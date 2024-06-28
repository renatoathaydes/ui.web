import { createCommandInput } from './commands.mts';

// load the Websockets module, which bootstraps the communication with the backend
import './ws.mts';

// include the CSS file in the output dir
import './styles.css';

/// UI.WEB startup function called by the index.html file.
export function startup() {
    createCommandInput();
}
