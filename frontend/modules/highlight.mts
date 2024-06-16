import hljs from 'highlight.js/lib/core';
import js from 'highlight.js/lib/languages/javascript';
import 'highlight.js/styles/github.css';

// Load any languages you need
hljs.registerLanguage('javascript', js);

export function highlightJs(value: string): string {
    return hljs.highlight(value, { language: 'javascript' }).value;
}
