import { EditorView, basicSetup } from "codemirror"

export function openEditor(path: string, contents: string) {
    let view = new EditorView({
        doc: contents,
        extensions: [
            basicSetup,
        ],
        parent: document.body
    });
}
