import { highlightJs } from "./highlight.mts";

export function createCommandInput() {
    const id = `cmd-input-${Date.now()}`;
    const history = [] as string[];
    let historyIndex = 0;
    const div = document.createElement('div');
    div.id = id;
    const el = document.createElement('input');
    const out = document.createElement('div');
    out.style.margin = '6px 4px';
    el.type = 'text';
    el.size = 50;
    el.onkeyup = (e) => {
        if (e.key === 'ArrowUp') {
            historyIndex = Math.max(0, historyIndex - 1);
            if (historyIndex < history.length) {
                el.value = history[historyIndex];
            }
        } else if (e.key === 'ArrowDown') {
            historyIndex = Math.min(history.length, historyIndex + 1);
            if (historyIndex < history.length) {
                el.value = history[historyIndex];
            } else {
                el.value = '';
            }
        } else if (e.key === 'Enter') {
            const cmd = el.value;
            if (cmd) {
                history.push(cmd);
                historyIndex = history.length;
                try {
                    const result = eval(cmd);
                    console.log('command result', result);
                } finally {
                    el.value = '';
                }
            }
        }
        out.innerHTML = highlightJs(el.value);
    };
    el.placeholder = 'Enter UI.WEB command here';
    document.body.appendChild(div);
    div.appendChild(el);
    div.appendChild(out);
    return div;
}
