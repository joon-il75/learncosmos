import type { Editor } from '@tiptap/react';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import {
  AlignCenterIcon, AlignLeftIcon, AlignRightIcon,
  BlockquoteIcon, BoldIcon, CodeBlockIcon, CodeIcon,
  HighlightIcon, ImagePlusIcon, ItalicIcon, LinkIcon,
  ListIcon, ListOrderedIcon, Redo2Icon, StrikeIcon,
  SubscriptIcon, SuperscriptIcon,
  TableColumnMinusIcon, TableColumnPlusIcon, TableDeleteIcon,
  TableIcon, TableRowMinusIcon, TableRowPlusIcon,
  UnderlineIcon, Undo2Icon,
} from './tiptapIcons';
import {
  createToolbarButtonStyle,
  createToolbarGroupStyle,
  createToolbarSelectStyle,
} from './tiptapEditorStyles';

type Props = {
  editor: Editor | null;
  editable: boolean;
  themeTokens: PointPageThemeTokens;
  isLinkPanelOpen: boolean;
  isImagePanelOpen: boolean;
  copy: PointLearningCopy['workspace']['editor'];
  onOpenLinkPanel: () => void;
  onOpenImagePanel: () => void;
  runCommand: (command: () => boolean) => void;
};

export default function TiptapToolbar({
  editor,
  editable,
  themeTokens,
  isLinkPanelOpen,
  isImagePanelOpen,
  copy,
  onOpenLinkPanel,
  onOpenImagePanel,
  runCommand,
}: Props) {
  const toolbarButtonStyle = (active = false) => createToolbarButtonStyle(themeTokens, editable, active);
  const toolbarGroupStyle = (withDivider = true) => createToolbarGroupStyle(themeTokens, withDivider);
  const toolbarSelectStyle = createToolbarSelectStyle(themeTokens, editable);
  const currentBlockStyle = editor?.isActive('heading', { level: 1 })
    ? 'h1'
    : editor?.isActive('heading', { level: 2 })
      ? 'h2'
      : editor?.isActive('heading', { level: 3 })
        ? 'h3'
        : 'paragraph';
  const isInsideTable = Boolean(editor?.isActive('table'));
  const tableButtonStyle = {
    ...toolbarButtonStyle(false),
    opacity: editable && isInsideTable ? 1 : 0.42,
  };

  const renderHistoryGroup = () => (
    <span style={toolbarGroupStyle()}>
      <button type="button" title={copy.undo} aria-label={copy.undo} disabled={!editable || !editor?.can().undo()} onClick={() => runCommand(() => editor!.chain().focus().undo().run())} style={toolbarButtonStyle(false)}><Undo2Icon /></button>
      <button type="button" title={copy.redo} aria-label={copy.redo} disabled={!editable || !editor?.can().redo()} onClick={() => runCommand(() => editor!.chain().focus().redo().run())} style={toolbarButtonStyle(false)}><Redo2Icon /></button>
    </span>
  );

  const renderBlockGroup = () => (
    <span style={toolbarGroupStyle()}>
      <select
        aria-label={copy.paragraphStyle}
        title={copy.paragraphStyle}
        disabled={!editable}
        value={currentBlockStyle}
        onChange={(event) => {
          const nextStyle = event.target.value;
          runCommand(() => {
            const chain = editor!.chain().focus();
            if (nextStyle === 'h1') return chain.setHeading({ level: 1 }).run();
            if (nextStyle === 'h2') return chain.setHeading({ level: 2 }).run();
            if (nextStyle === 'h3') return chain.setHeading({ level: 3 }).run();
            return chain.setParagraph().run();
          });
        }}
        style={toolbarSelectStyle}
      >
        <option value="paragraph">{copy.paragraph}</option>
        <option value="h1">{copy.heading1}</option>
        <option value="h2">{copy.heading2}</option>
        <option value="h3">{copy.heading3}</option>
      </select>
      <button type="button" title={copy.quote} aria-label={copy.quote} aria-pressed={Boolean(editor?.isActive('blockquote'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleBlockquote().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('blockquote')))}><BlockquoteIcon /></button>
      <button type="button" title={copy.codeBlock} aria-label={copy.codeBlock} aria-pressed={Boolean(editor?.isActive('codeBlock'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleCodeBlock().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('codeBlock')))}><CodeBlockIcon /></button>
    </span>
  );

  const renderInlineGroup = () => (
    <span style={toolbarGroupStyle()}>
      <button type="button" title={copy.bold} aria-label={copy.bold} aria-pressed={Boolean(editor?.isActive('bold'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleBold().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('bold')))}><BoldIcon /></button>
      <button type="button" title={copy.italic} aria-label={copy.italic} aria-pressed={Boolean(editor?.isActive('italic'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleItalic().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('italic')))}><ItalicIcon /></button>
      <button type="button" title={copy.underline} aria-label={copy.underline} aria-pressed={Boolean(editor?.isActive('underline'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleUnderline().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('underline')))}><UnderlineIcon /></button>
      <button type="button" title={copy.strike} aria-label={copy.strike} aria-pressed={Boolean(editor?.isActive('strike'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleStrike().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('strike')))}><StrikeIcon /></button>
      <button type="button" title={copy.highlight} aria-label={copy.highlight} aria-pressed={Boolean(editor?.isActive('highlight'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleHighlight().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('highlight')))}><HighlightIcon /></button>
      <button type="button" title={copy.inlineCode} aria-label={copy.inlineCode} aria-pressed={Boolean(editor?.isActive('code'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleCode().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('code')))}><CodeIcon /></button>
      <button type="button" title={copy.superscript} aria-label={copy.superscript} aria-pressed={Boolean(editor?.isActive('superscript'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleSuperscript().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('superscript')))}><SuperscriptIcon /></button>
      <button type="button" title={copy.subscript} aria-label={copy.subscript} aria-pressed={Boolean(editor?.isActive('subscript'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleSubscript().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('subscript')))}><SubscriptIcon /></button>
    </span>
  );

  const renderListGroup = () => (
    <span style={toolbarGroupStyle()}>
      <button type="button" title={copy.bulletList} aria-label={copy.bulletList} aria-pressed={Boolean(editor?.isActive('bulletList'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleBulletList().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('bulletList')))}><ListIcon /></button>
      <button type="button" title={copy.orderedList} aria-label={copy.orderedList} aria-pressed={Boolean(editor?.isActive('orderedList'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleOrderedList().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('orderedList')))}><ListOrderedIcon /></button>
    </span>
  );

  const renderAlignGroup = () => (
    <span style={toolbarGroupStyle()}>
      <button type="button" title={copy.alignLeft} aria-label={copy.alignLeft} aria-pressed={Boolean(editor?.isActive({ textAlign: 'left' }))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().setTextAlign('left').run())} style={toolbarButtonStyle(Boolean(editor?.isActive({ textAlign: 'left' })))}><AlignLeftIcon /></button>
      <button type="button" title={copy.alignCenter} aria-label={copy.alignCenter} aria-pressed={Boolean(editor?.isActive({ textAlign: 'center' }))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().setTextAlign('center').run())} style={toolbarButtonStyle(Boolean(editor?.isActive({ textAlign: 'center' })))}><AlignCenterIcon /></button>
      <button type="button" title={copy.alignRight} aria-label={copy.alignRight} aria-pressed={Boolean(editor?.isActive({ textAlign: 'right' }))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().setTextAlign('right').run())} style={toolbarButtonStyle(Boolean(editor?.isActive({ textAlign: 'right' })))}><AlignRightIcon /></button>
    </span>
  );

  const renderInsertGroup = () => (
    <span style={toolbarGroupStyle(false)}>
      <button type="button" title={copy.link} aria-label={copy.link} aria-pressed={Boolean(editor?.isActive('link'))} disabled={!editable} onClick={onOpenLinkPanel} style={toolbarButtonStyle(Boolean(editor?.isActive('link') || isLinkPanelOpen))}><LinkIcon /></button>
      <button type="button" title={copy.image} aria-label={copy.image} disabled={!editable} onClick={onOpenImagePanel} style={toolbarButtonStyle(isImagePanelOpen)}><ImagePlusIcon /></button>
      <button type="button" title={copy.insertTable} aria-label={copy.insertTable} aria-pressed={isInsideTable} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run())} style={toolbarButtonStyle(isInsideTable)}><TableIcon /></button>
      <button type="button" title={copy.addColumn} aria-label={copy.addColumn} disabled={!editable || !isInsideTable} onClick={() => runCommand(() => editor!.chain().focus().addColumnAfter().run())} style={tableButtonStyle}><TableColumnPlusIcon /></button>
      <button type="button" title={copy.addRow} aria-label={copy.addRow} disabled={!editable || !isInsideTable} onClick={() => runCommand(() => editor!.chain().focus().addRowAfter().run())} style={tableButtonStyle}><TableRowPlusIcon /></button>
      <button type="button" title={copy.deleteColumn} aria-label={copy.deleteColumn} disabled={!editable || !isInsideTable} onClick={() => runCommand(() => editor!.chain().focus().deleteColumn().run())} style={tableButtonStyle}><TableColumnMinusIcon /></button>
      <button type="button" title={copy.deleteRow} aria-label={copy.deleteRow} disabled={!editable || !isInsideTable} onClick={() => runCommand(() => editor!.chain().focus().deleteRow().run())} style={tableButtonStyle}><TableRowMinusIcon /></button>
      <button type="button" title={copy.deleteTable} aria-label={copy.deleteTable} disabled={!editable || !isInsideTable} onClick={() => runCommand(() => editor!.chain().focus().deleteTable().run())} style={tableButtonStyle}><TableDeleteIcon /></button>
    </span>
  );

  return (
    <div data-testid="tiptap-sticky-toolbar" className="lw-tiptap-toolbar">
      <div className="lw-tiptap-toolbar-desktop">
        {renderHistoryGroup()}
        {renderBlockGroup()}
        {renderInlineGroup()}
        {renderListGroup()}
        {renderAlignGroup()}
        {renderInsertGroup()}
      </div>
      <div className="lw-tiptap-toolbar-mobile">
        {renderHistoryGroup()}
        <span style={toolbarGroupStyle()}>
          <button type="button" title={copy.bold} aria-label={copy.bold} aria-pressed={Boolean(editor?.isActive('bold'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleBold().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('bold')))}><BoldIcon /></button>
          <button type="button" title={copy.italic} aria-label={copy.italic} aria-pressed={Boolean(editor?.isActive('italic'))} disabled={!editable} onClick={() => runCommand(() => editor!.chain().focus().toggleItalic().run())} style={toolbarButtonStyle(Boolean(editor?.isActive('italic')))}><ItalicIcon /></button>
          <button type="button" title={copy.link} aria-label={copy.link} aria-pressed={Boolean(editor?.isActive('link'))} disabled={!editable} onClick={onOpenLinkPanel} style={toolbarButtonStyle(Boolean(editor?.isActive('link') || isLinkPanelOpen))}><LinkIcon /></button>
        </span>
        <details className="lw-tiptap-toolbar-menu">
          <summary aria-label="more editor tools" title="more editor tools" style={toolbarButtonStyle(false)}>...</summary>
          <div className="lw-tiptap-toolbar-menu-panel">
            {renderBlockGroup()}
            {renderInlineGroup()}
            {renderListGroup()}
            {renderAlignGroup()}
            {renderInsertGroup()}
          </div>
        </details>
      </div>
    </div>
  );
}
