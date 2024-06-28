import { CommandResponse } from "../../common/command.mts";
import { callBackend } from "./ws.mts";

type CommandMode = "FE-JS" | "BE-JS";

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
    @import url( '/startup.css' )
  </style>
  <div class="command-input">
    <span>
      <input type="text" size="50" placeholder="Enter UI.WEB command"></input>
      <span class="command-mode">
        <select class="mode" name="mode">
          <option value="FE-JS">FE-JS</option>
          <option value="BE-JS">BE-JS</option>
        </select>
      </span>
    </span>
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
    modeSelector!: HTMLSelectElement
    mode: CommandMode = "FE-JS"
    verbose: boolean = false

    constructor() {
        super();
        console.log('Creating CommandInput Component');
        const shadow = this.attachShadow({ mode: "open" });
        shadow.appendChild(template.content.cloneNode(true));
        const mySheet = new CSSStyleSheet();
        mySheet.replaceSync(css);
        shadow.adoptedStyleSheets = [mySheet];
    }

    connectedCallback() {
        console.log('Attaching CommandInput Component');
        this.textInput = this.shadowRoot!.querySelector('input')!;
        this.codeView = this.shadowRoot!.querySelector('div.codeview')!;
        this.output = this.shadowRoot!.querySelector('div.output')!;
        this.modeSelector = this.shadowRoot!.querySelector('select.mode')! as HTMLSelectElement;
        this.modeSelector.onchange = this.onModeChange;
        this.textInput.onkeyup = this.onInputKeyup;
    }

    disconnectedCallback() {
        console.log('Removing CommandInput Component');
    }

    onModeChange = () => {
        this.mode = this.modeSelector.value as CommandMode;
        console.log('CommandInput mode changed to', this.mode);
    }

    onInputKeyup = async (e: KeyboardEvent) => {
        if (e.key === 'ArrowUp') {
            this.historyIndex = Math.max(0, this.historyIndex - 1);
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
                // TODO show code in a read-only editor?
                this.codeView.innerHTML = cmd;
                try {
                    await this.runCommand(cmd);
                } catch (e) {
                    this.showOutput(e, true);
                } finally {
                    this.textInput.value = '';
                }
            }
        }
    }

    async runCommand(cmd: string) {
        switch (this.mode) {
            case "BE-JS":
                const response = await callBackend(cmd);
                return this.handleBeResponse(response);
            case "FE-JS":
                // indirect eval exposing this as self
                return this.showOutput(this.eval(cmd));
        }
    }

    handleBeResponse(resp: CommandResponse) {
        if ("error" in resp) {
            this.showOutput(resp.error, true);
        } else {
            let value = resp.value;
            this.showOutput(value);
            if (resp.feCmd) {
                this.eval(resp.feCmd, value);
            }
        }
    }

    async showOutput(result: any, isError: boolean = false) {
        if (isError) {
            this.output.innerText = result?.toString() ?? 'ERROR';
            this.output.classList.add('error');
        } else try {
            const value = await result;
            if (this.verbose) {
                console.log('command result', value);
            }
            this.output.innerText = JSON.stringify(value);
            this.output.classList.remove('error');
        } catch (e) {
            this.output.innerText = e.toString();
            this.output.classList.add('error');
        }
    }

    /// Eval the given cmd in the global scope, but expose this component
    /// as the "me" variable.
    async eval(cmd: string, value?: any): Promise<any> {
        return evalWith(cmd, this, value);
    }
}

window.customElements.define('command-input', CommandInputElement);

export function createCommandInput() {
    const component = document.createElement('command-input');
    document.body.appendChild(component);
}
