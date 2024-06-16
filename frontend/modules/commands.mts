import { highlightJs } from "./highlight.mts";

export function createCommandInput() {
    const id = `cmd-input-${Date.now()}`;
    const history = [] as string[];
    let historyIndex = 0;
    const div = document.createElement('div');
    div.id = id;
    const el = document.createElement('input');
    const out = document.createElement('div');
    el.type = 'text';
    el.size = 50;
    el.onkeyup = (e) => {
        console.log(`Key: ${e.key}`);
        if (e.metaKey) {
            console.log('meta key');
        }
        if (e.ctrlKey) {
            console.log('ctrl key');
        }
        if (e.key === 'ArrowUp') {
            historyIndex = Math.max(0, historyIndex);
            if (historyIndex < history.length) {
                el.value = history[historyIndex];
            }
        }
        if (e.key === 'ArrowDown') {
            historyIndex = Math.min(history.length - 1, historyIndex + 1);
            if (historyIndex < history.length) {
                el.value = history[historyIndex];
            }
        }
    };
    el.placeholder = 'Enter UI.WEB command here';
    el.onchange = () => {
        const cmd = el.value;
        if (cmd) {
            history.push(cmd);
            pushHtmlTo(out, cmd);
            historyIndex = history.length;
            try {
                const result = eval(cmd);
                console.log('command result', result);
            } finally {
                el.value = '';
            }
        }
    };

    document.body.appendChild(div);
    div.appendChild(el);
    div.appendChild(out);
    return div;
}

function pushHtmlTo(out: HTMLElement, text: string) {
    const html = highlightJs(text);
    const el = document.createElement('div');
    el.innerHTML = html;
    out.appendChild(el);
}
