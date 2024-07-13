import { Extension } from '@codemirror/state';
import { EditorView, basicSetup } from "codemirror"

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

type Language = {
    "name": string,
    "extensions": Set<string>,
};

export const languages: Language[] = [
    { "name": "javascript", "extensions": new Set([".js", ".ts", ".json", ".mjs", ".mts"]) }
];

class EditorParent extends HTMLElement {

    readonly languages: Language[]
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

export async function openEditor(name: string, contents: string): Promise<string> {
    const ext = extension(name);
    const parentElement = document.createElement('editor-parent') as EditorParent;
    parentElement.fileNameElement.textContent = name;
    document.body.appendChild(parentElement);
    const extensions: Extension[] = [basicSetup];
    if (ext) {
        const lang = languages.find(lang => lang.extensions.has(ext));
        if (lang) {
            const langMod = await import(`@codemirror/lang-${lang.name}`);
            extensions.push(langMod[lang.name]());
        } else {
            console.log('No Codemirror language extension found for:', ext);
        }
    }
    new EditorView({
        doc: contents,
        extensions,
        parent: parentElement.rootElement,
    });
    return `Read ${contents.length} bytes.`;
}

function extension(path: string): string | null {
    const index = path.lastIndexOf('.');
    if (index < 0) return null;
    return path.substring(index);
}
