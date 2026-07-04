import type { PointPageThemeTokens } from '../pointPageUtils';

type Props = {
  themeTokens: PointPageThemeTokens;
};

export default function TiptapEditorGlobalStyles({ themeTokens }: Props) {
  return (
    <style jsx global>{`
      .lw-tiptap-editor {
        width: 100%;
        max-width: 860px;
        margin-left: auto;
        margin-right: auto;
        min-height: 170px;
        padding: 16px;
        color: ${themeTokens.inputText};
        font-size: 16px;
        line-height: 1.75;
        outline: none;
      }
      .lw-tiptap-toolbar {
        padding: 8px;
        max-width: 100%;
        scrollbar-width: thin;
        scrollbar-color: ${themeTokens.secondaryButtonBorder} transparent;
      }
      .lw-tiptap-toolbar-desktop,
      .lw-tiptap-toolbar-mobile {
        display: flex;
        flex-wrap: nowrap;
        align-items: center;
        gap: 5px;
        max-width: 100%;
        overflow-x: auto;
        overflow-y: hidden;
      }
      .lw-tiptap-toolbar-mobile {
        display: none;
        overflow: visible;
      }
      .lw-tiptap-toolbar-menu {
        position: relative;
        flex: 0 0 auto;
      }
      .lw-tiptap-toolbar-menu > summary {
        list-style: none;
        user-select: none;
      }
      .lw-tiptap-toolbar-menu > summary::-webkit-details-marker {
        display: none;
      }
      .lw-tiptap-toolbar-menu-panel {
        position: absolute;
        right: 0;
        top: calc(100% + 8px);
        z-index: 40;
        display: grid;
        grid-template-columns: 1fr;
        gap: 8px;
        width: min(280px, calc(100vw - 32px));
        max-width: calc(100vw - 32px);
        max-height: min(420px, calc(100vh - 220px));
        overflow: auto;
        padding: 10px;
        border-radius: 8px;
        border: 1px solid ${themeTokens.surfaceBorder};
        background: ${themeTokens.inputBackground};
        box-shadow: 0 18px 42px rgba(15, 23, 42, 0.18);
      }
      .lw-tiptap-toolbar-menu-panel > span {
        border-right: 0 !important;
        margin-right: 0 !important;
        padding-right: 0 !important;
        flex-wrap: wrap;
      }
      @media (max-width: 640px) {
        .lw-tiptap-toolbar {
          padding: 7px;
        }
        .lw-tiptap-toolbar-desktop {
          display: none;
        }
        .lw-tiptap-toolbar-mobile {
          display: flex;
          justify-content: space-between;
        }
        .lw-tiptap-toolbar-menu-panel {
          grid-template-columns: 1fr;
        }
        .lw-tiptap-editor {
          min-height: 280px;
          padding: 13px;
          font-size: 15px;
          line-height: 1.72;
        }
      }
      .lw-tiptap-toolbar::-webkit-scrollbar {
        height: 6px;
      }
      .lw-tiptap-toolbar::-webkit-scrollbar-track {
        background: transparent;
      }
      .lw-tiptap-toolbar::-webkit-scrollbar-thumb {
        border-radius: 999px;
        background: ${themeTokens.secondaryButtonBorder};
      }
      .lw-tiptap-editor p {
        margin: 0 0 0.85em;
      }
      .lw-tiptap-editor h1,
      .lw-tiptap-editor h2,
      .lw-tiptap-editor h3 {
        margin: 1.1em 0 0.45em;
        color: ${themeTokens.inputText};
        font-weight: 900;
        line-height: 1.28;
      }
      .lw-tiptap-editor h1 {
        font-size: 1.55em;
      }
      .lw-tiptap-editor h2 {
        font-size: 1.35em;
      }
      .lw-tiptap-editor h3 {
        font-size: 1.18em;
      }
      .lw-tiptap-editor h1:first-child,
      .lw-tiptap-editor h2:first-child,
      .lw-tiptap-editor h3:first-child {
        margin-top: 0;
      }
      .lw-tiptap-editor ul,
      .lw-tiptap-editor ol {
        margin: 0 0 0.85em;
        padding-left: 1.7em;
      }
      .lw-tiptap-editor ul {
        list-style-type: disc;
      }
      .lw-tiptap-editor ol {
        list-style-type: decimal;
      }
      .lw-tiptap-editor li {
        margin: 0.2em 0;
        padding-left: 0.1em;
      }
      .lw-tiptap-editor li > p {
        margin: 0.15em 0;
      }
      .lw-tiptap-editor blockquote {
        margin: 0 0 0.85em;
        padding-left: 12px;
        border-left: 3px solid ${themeTokens.primaryButtonBorder};
        color: ${themeTokens.description};
      }
      .lw-tiptap-editor a {
        color: ${themeTokens.metaValue};
        text-decoration: underline;
        text-underline-offset: 3px;
      }
      .lw-tiptap-editor mark {
        border-radius: 4px;
        background: rgba(250, 204, 21, 0.36);
        color: inherit;
        padding: 0 2px;
      }
      .lw-tiptap-editor code {
        border-radius: 5px;
        background: ${themeTokens.placeholderBackground};
        border: 1px solid ${themeTokens.surfaceBorder};
        color: ${themeTokens.inputText};
        font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace;
        font-size: 0.92em;
        padding: 0.08em 0.32em;
      }
      .lw-tiptap-editor pre {
        margin: 0 0 0.85em;
        overflow-x: auto;
        border-radius: 10px;
        border: 1px solid ${themeTokens.surfaceBorder};
        background: ${themeTokens.placeholderBackground};
        color: ${themeTokens.inputText};
        padding: 12px 14px;
      }
      .lw-tiptap-editor pre code {
        display: block;
        border: 0;
        background: transparent;
        padding: 0;
        white-space: pre;
      }
      .lw-tiptap-editor sup,
      .lw-tiptap-editor sub {
        line-height: 0;
      }
      .lw-tiptap-editor img {
        display: block;
        max-width: 100%;
        height: auto;
        margin: 12px 0;
        border-radius: 10px;
        border: 1px solid ${themeTokens.surfaceBorder};
      }
      .lw-tiptap-editor .tableWrapper {
        margin: 0 0 0.9em;
        overflow-x: auto;
      }
      .lw-tiptap-editor table {
        width: 100%;
        border-collapse: collapse;
        table-layout: fixed;
        font-size: 0.94em;
      }
      .lw-tiptap-editor th,
      .lw-tiptap-editor td {
        min-width: 88px;
        border: 1px solid ${themeTokens.surfaceBorder};
        padding: 8px 10px;
        vertical-align: top;
      }
      .lw-tiptap-editor th {
        background: ${themeTokens.placeholderBackground};
        color: ${themeTokens.inputText};
        font-weight: 900;
      }
      .lw-tiptap-editor td {
        background: ${themeTokens.inputBackground};
      }
      .lw-tiptap-editor th p,
      .lw-tiptap-editor td p {
        margin: 0;
      }
      .lw-tiptap-editor .selectedCell::after {
        content: "";
        position: absolute;
        inset: 0;
        pointer-events: none;
        background: rgba(59, 130, 246, 0.16);
      }
      .lw-tiptap-editor th,
      .lw-tiptap-editor td {
        position: relative;
      }
      .lw-tiptap-editor:empty::before {
        content: attr(data-placeholder);
        color: ${themeTokens.mutedText};
        float: left;
        pointer-events: none;
        height: 0;
      }
    `}</style>
  );
}
