import { highlightJs } from "./highlight.mts";

const template = document.createElement('template');

const css = `
  .command-input {
    margin: 5px;
  }
  .codeview {
    min-height: 1.1em;
    border: solid gray 1px;
    padding: 3px;
  }
  .output {
    height: 100px;
    border: solid darkgray 1px;
    padding: 3px;
  }
`;

template.innerHTML = `
  <style>
    @import url( '/index.css' )
  </style>
  <div class="command-input">
    <input type="text" size="50" placeholder="Enter UI.WEB command"></input>
    <div class="small-label">code:</div>
    <div class="codeview"></div>
    <div class="small-label">output:</div>
    <div class="output" height="100px"></div>
  </div>
`;

class CommandInputElement extends HTMLElement {

    history: Array<string> = [];
    historyIndex = 0;
    textInput!: HTMLInputElement
    codeView!: HTMLDivElement
    output!: HTMLDivElement

    constructor() {
        super();
        console.log('Creating My Web Component');
        const shadow = this.attachShadow({ mode: "open" });
        shadow.appendChild(template.content.cloneNode(true));
        const mySheet = new CSSStyleSheet();
        mySheet.replaceSync(css);
        shadow.adoptedStyleSheets = [mySheet];
    }

    connectedCallback() {
        console.log('Attaching My Web Component');
        this.textInput = this.shadowRoot!.querySelector('input')!;
        this.codeView = this.shadowRoot!.querySelector('div.codeview')!;
        this.output = this.shadowRoot!.querySelector('div.output')!;
        this.textInput.onkeyup = this.onInputKeyup;
    }

    disconnectedCallback() {
        console.log('Removing My Web Component');
    }

    onInputKeyup = (e: KeyboardEvent) => {
        if (e.key === 'ArrowUp') {
            this.historyIndex = Math.max(0, this.historyIndex - 1);
            console.log('historyIndex:', this.historyIndex, 'history.length', this.history.length);
            if (this.historyIndex < this.history.length) {
                this.textInput.value = this.history[this.historyIndex];
            }
        } else if (e.key === 'ArrowDown') {
            this.historyIndex = Math.min(this.history.length, this.historyIndex + 1);
            if (this.historyIndex < this.history.length) {
                this.textInput.value = this.history[this.historyIndex];
            } else {
                this.textInput.value = '';
            }
        } else if (e.key === 'Enter') {
            const cmd = this.textInput.value;
            if (cmd) {
                this.history.push(cmd);
                this.historyIndex = this.history.length;
                this.codeView.innerHTML = highlightJs(cmd);
                try {
                    const result = eval(cmd);
                    showOutput(this.output, result);
                } catch (e) {
                    showOutput(this.output, e, true);
                } finally {
                    this.textInput.value = '';
                }
            }
        }
    }
}

window.customElements.define('command-input', CommandInputElement);

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

export function createCommandInput() {
    const component = document.createElement('command-input');
    document.body.appendChild(component);
}
