import { highlightJs } from "./highlight.mts";

export function createCommandInput() {
    const id = `cmd-input-${Date.now()}`;
    const history = [] as string[];
    let historyIndex = 0;
    const div = document.createElement('div');
    div.id = id;
    const el = document.createElement('input');
    const codeView = document.createElement('div');
    codeView.style.margin = '6px 4px';
    codeView.style.minHeight = '1.1em';
    const out = document.createElement('div');
    out.style.height = '100px';
    out.style.border = 'solid black 1px';
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
                    showOutput(out, result);
                } catch (e) {
                    showOutput(out, e, true);
                } finally {
                    el.value = '';
                }
            }
        }
        codeView.innerHTML = highlightJs(el.value);
    };
    el.placeholder = 'Enter UI.WEB command here';
    document.body.appendChild(div);
    div.appendChild(el);
    div.appendChild(codeView);
    div.appendChild(out);
    return div;
}

function showOutput(out: HTMLDivElement, result: any, isError: boolean = false) {
    const show = (value: any, error: boolean = false) => {
        console.log('command result', value);
        if (error) {
            out.innerText = value.toString();
            out.classList.add('error');
        } else {
            out.innerText = JSON.stringify(value);
            out.classList.remove('error');
        }
    };
    if (result instanceof Promise) {
        result.catch(t => show(t, true)).then(show);
    } else {
        show(result, isError);
    }
}

