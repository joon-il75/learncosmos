'use client';

import { useEffect } from 'react';
import { useState } from 'react';
import { EditorContent, useEditor } from '@tiptap/react';
import Highlight from '@tiptap/extension-highlight';
import Image from '@tiptap/extension-image';
import Subscript from '@tiptap/extension-subscript';
import Superscript from '@tiptap/extension-superscript';
import { Table } from '@tiptap/extension-table';
import { TableCell } from '@tiptap/extension-table-cell';
import { TableHeader } from '@tiptap/extension-table-header';
import { TableRow } from '@tiptap/extension-table-row';
import TextAlign from '@tiptap/extension-text-align';
import Underline from '@tiptap/extension-underline';
import StarterKit from '@tiptap/starter-kit';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import TiptapEditorGlobalStyles from './TiptapEditorGlobalStyles';
import TiptapImagePanel from './TiptapImagePanel';
import TiptapLinkPanel from './TiptapLinkPanel';
import TiptapToolbar from './TiptapToolbar';
import {
  isHttpsUrl,
  normalizeContent,
  validateInlineImageFile,
} from './tiptapEditorUtils';

interface Props {
  value: string;
  editable: boolean;
  placeholder: string;
  themeTokens: PointPageThemeTokens;
  copy: PointLearningCopy['workspace']['editor'];
  onChange: (value: string) => void;
  onUploadImage?: (file: File) => Promise<{ src: string; alt?: string }>;
  toolbarStickyTop?: string;
}

export default function TiptapResearchEditor({ value, editable, placeholder, themeTokens, copy, onChange, onUploadImage, toolbarStickyTop = '0px' }: Props) {
  const [linkInput, setLinkInput] = useState('');
  const [linkError, setLinkError] = useState('');
  const [isLinkPanelOpen, setIsLinkPanelOpen] = useState(false);
  const [imageInput, setImageInput] = useState('');
  const [imageAltInput, setImageAltInput] = useState('');
  const [selectedImageFile, setSelectedImageFile] = useState<File | null>(null);
  const [imageError, setImageError] = useState('');
  const [isUploadingImage, setIsUploadingImage] = useState(false);
  const [isImagePanelOpen, setIsImagePanelOpen] = useState(false);

  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        link: {
          autolink: true,
          defaultProtocol: 'https',
          enableClickSelection: true,
          linkOnPaste: true,
          openOnClick: false,
          protocols: ['https'],
          validate: isHttpsUrl,
          shouldAutoLink: isHttpsUrl,
          isAllowedUri: (url, { defaultValidate }) => isHttpsUrl(url) && defaultValidate(url),
          HTMLAttributes: {
            rel: 'noopener noreferrer',
            target: '_blank',
          },
        },
      }),
      Underline,
      Highlight,
      Superscript,
      Subscript,
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
      Table.configure({
        resizable: true,
        HTMLAttributes: {
          class: 'lw-tiptap-table',
        },
      }),
      TableRow,
      TableHeader,
      TableCell,
      Image.configure({
        allowBase64: false,
        HTMLAttributes: {
          loading: 'lazy',
        },
      }),
    ],
    content: normalizeContent(value),
    editable,
    immediatelyRender: false,
    editorProps: {
      attributes: {
        class: 'lw-tiptap-editor',
        'data-placeholder': placeholder,
      },
    },
    onUpdate: ({ editor: currentEditor }) => {
      onChange(currentEditor.getHTML());
    },
  });

  useEffect(() => {
    if (!editor) return;
    editor.setEditable(editable);
  }, [editable, editor]);

  useEffect(() => {
    if (!editor) return;
    const next = normalizeContent(value);
    if (editor.getHTML() !== next) {
      editor.commands.setContent(next, { emitUpdate: false });
    }
  }, [editor, value]);

  const runCommand = (command: () => boolean) => {
    if (!editor || !editable) return;
    command();
  };

  const openLinkPanel = () => {
    if (!editor || !editable) return;
    setLinkInput(String(editor.getAttributes('link').href ?? ''));
    setLinkError('');
    setIsImagePanelOpen(false);
    setIsLinkPanelOpen((current) => !current);
  };

  const openImagePanel = () => {
    if (!editor || !editable) return;
    setImageInput('');
    setImageAltInput('');
    setSelectedImageFile(null);
    setImageError('');
    setIsLinkPanelOpen(false);
    setIsImagePanelOpen((current) => !current);
  };

  const applyLink = () => {
    if (!editor || !editable) return;
    const href = linkInput.trim();
    if (!isHttpsUrl(href)) {
      setLinkError(copy.linkHttpsOnly);
      return;
    }

    const chain = editor.chain().focus().extendMarkRange('link');
    if (editor.state.selection.empty && !editor.isActive('link')) {
      chain.insertContent({
        type: 'text',
        text: href,
        marks: [{ type: 'link', attrs: { href, target: '_blank', rel: 'noopener noreferrer' } }],
      }).run();
    } else {
      chain.setLink({ href, target: '_blank', rel: 'noopener noreferrer' }).run();
    }
    setIsLinkPanelOpen(false);
    setLinkError('');
  };

  const unsetLink = () => {
    if (!editor || !editable) return;
    editor.chain().focus().extendMarkRange('link').unsetLink().run();
    setIsLinkPanelOpen(false);
    setLinkError('');
  };

  const applyImage = () => {
    if (!editor || !editable) return;
    const src = imageInput.trim();
    if (!isHttpsUrl(src)) {
      setImageError(copy.imageHttpsOnly);
      return;
    }
    editor.chain().focus().setImage({ src, alt: imageAltInput.trim() || undefined }).run();
    setIsImagePanelOpen(false);
    setImageError('');
    setImageInput('');
    setImageAltInput('');
  };

  const uploadImage = async () => {
    if (!editor || !editable || !onUploadImage || !selectedImageFile || isUploadingImage) return;
    const validationError = validateInlineImageFile(selectedImageFile, {
      typeError: copy.inlineImageTypeError,
      sizeError: copy.inlineImageSizeError,
    });
    if (validationError) {
      setImageError(validationError);
      return;
    }
    setIsUploadingImage(true);
    setImageError('');
    try {
      const uploaded = await onUploadImage(selectedImageFile);
      editor.chain().focus().setImage({ src: uploaded.src, alt: imageAltInput.trim() || uploaded.alt || selectedImageFile.name }).run();
      setIsImagePanelOpen(false);
      setSelectedImageFile(null);
      setImageInput('');
      setImageAltInput('');
    } catch (error) {
      setImageError(error instanceof Error ? error.message : copy.imageUploadFailed);
    } finally {
      setIsUploadingImage(false);
    }
  };

  return (
    <div
      style={{
        overflow: 'visible',
        borderRadius: '14px',
        border: `1px solid ${themeTokens.inputBorder}`,
        background: themeTokens.inputBackground,
      }}
    >
      <div
        style={{
          position: editable ? 'sticky' : 'relative',
          top: editable ? toolbarStickyTop : 0,
          zIndex: 24,
          overflow: 'visible',
          borderTopLeftRadius: '14px',
          borderTopRightRadius: '14px',
          borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
          background: themeTokens.placeholderBackground,
          boxShadow: editable ? '0 10px 22px rgba(15, 23, 42, 0.10)' : 'none',
        }}
      >
        <TiptapToolbar
          editor={editor}
          editable={editable}
          themeTokens={themeTokens}
          isLinkPanelOpen={isLinkPanelOpen}
          isImagePanelOpen={isImagePanelOpen}
          copy={copy}
          onOpenLinkPanel={openLinkPanel}
          onOpenImagePanel={openImagePanel}
          runCommand={runCommand}
        />
        {isLinkPanelOpen ? (
          <TiptapLinkPanel
            linkInput={linkInput}
            setLinkInput={setLinkInput}
            linkError={linkError}
            setLinkError={setLinkError}
            isLinkActive={Boolean(editor?.isActive('link'))}
            editable={editable}
            themeTokens={themeTokens}
            copy={copy}
            onApplyLink={applyLink}
            onUnsetLink={unsetLink}
            onClose={() => {
              setIsLinkPanelOpen(false);
              setLinkError('');
            }}
          />
        ) : null}
        {isImagePanelOpen ? (
          <TiptapImagePanel
            imageInput={imageInput}
            setImageInput={setImageInput}
            imageAltInput={imageAltInput}
            setImageAltInput={setImageAltInput}
            selectedImageFile={selectedImageFile}
            setSelectedImageFile={setSelectedImageFile}
            imageError={imageError}
            setImageError={setImageError}
            isUploadingImage={isUploadingImage}
            editable={editable}
            canUploadImage={Boolean(onUploadImage)}
            themeTokens={themeTokens}
            copy={copy}
            onApplyImage={applyImage}
            onUploadImage={uploadImage}
            onClose={() => {
              setIsImagePanelOpen(false);
              setImageError('');
            }}
          />
        ) : null}
      </div>
      <EditorContent editor={editor} />
      <TiptapEditorGlobalStyles themeTokens={themeTokens} />
    </div>
  );
}
