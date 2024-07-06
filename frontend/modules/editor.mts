import { EditorView, basicSetup } from "codemirror"
import { javascript } from "@codemirror/lang-javascript"

const template = document.createElement('template');

const css = `
.editor-file-name {
	padding: 2px;
	background: black;
	color: white;
	font-family: monospace;
	font-size: small;
}

.cm-editor {
	max-height: 20em;
	border: solid lightgray 1px;
}
`;

template.innerHTML = `
<div class='editor-parent'>
    <div class='editor-file-name'></div>
<div>
`;

class EditorParent extends HTMLElement {

    readonly rootElement: HTMLDivElement
    readonly fileNameElement: HTMLDivElement

    constructor() {
        super();
        const shadow = this.attachShadow({ mode: "open" });
        shadow.appendChild(template.content.cloneNode(true));
        const mySheet = new CSSStyleSheet();
        mySheet.replaceSync(css);
        shadow.adoptedStyleSheets = [mySheet];
        this.rootElement = shadow.querySelector('.editor-parent')!
        this.fileNameElement = shadow.querySelector('.editor-file-name')!
    }
}

window.customElements.define('editor-parent', EditorParent);

export function openEditor(path: string, contents: string): string {
    const parentElement = document.createElement('editor-parent') as EditorParent;
    parentElement.fileNameElement.textContent = path;
    document.body.appendChild(parentElement);
    new EditorView({
        doc: contents,
        extensions: [
            basicSetup,
            javascript(),
        ],
        parent: parentElement.rootElement,
    });
    return `Read ${contents.length} bytes.`;
}
