import type { Dispatch, SetStateAction } from 'react';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import {
  inlineImageAccept,
  getInlineImageUploadHint,
  validateInlineImageFile,
} from './tiptapEditorUtils';
import { createToolbarTextButtonStyle } from './tiptapEditorStyles';

type Props = {
  imageInput: string;
  setImageInput: Dispatch<SetStateAction<string>>;
  imageAltInput: string;
  setImageAltInput: Dispatch<SetStateAction<string>>;
  selectedImageFile: File | null;
  setSelectedImageFile: Dispatch<SetStateAction<File | null>>;
  imageError: string;
  setImageError: Dispatch<SetStateAction<string>>;
  isUploadingImage: boolean;
  editable: boolean;
  canUploadImage: boolean;
  themeTokens: PointPageThemeTokens;
  copy: PointLearningCopy['workspace']['editor'];
  onApplyImage: () => void;
  onUploadImage: () => void;
  onClose: () => void;
};

export default function TiptapImagePanel({
  imageInput,
  setImageInput,
  imageAltInput,
  setImageAltInput,
  selectedImageFile,
  setSelectedImageFile,
  imageError,
  setImageError,
  isUploadingImage,
  editable,
  canUploadImage,
  themeTokens,
  copy,
  onApplyImage,
  onUploadImage,
  onClose,
}: Props) {
  const textButtonStyle = createToolbarTextButtonStyle(themeTokens, editable);

  return (
    <div
      style={{
        display: 'grid',
        gap: '8px',
        padding: '10px',
        borderBottom: `1px solid ${themeTokens.surfaceBorder}`,
        background: themeTokens.surfaceBackground,
      }}
    >
      <div style={{ display: 'grid', gap: '8px' }}>
        <input
          value={imageInput}
          onChange={(event) => {
            setImageInput(event.target.value);
            setImageError('');
          }}
          onKeyDown={(event) => {
            if (event.key === 'Enter') {
              event.preventDefault();
              onApplyImage();
            }
            if (event.key === 'Escape') {
              event.preventDefault();
              onClose();
            }
          }}
          placeholder={copy.imageUrlPlaceholder}
          style={{
            width: '100%',
            minWidth: 0,
            height: '38px',
            borderRadius: '8px',
            border: `1px solid ${themeTokens.inputBorder}`,
            background: themeTokens.inputBackground,
            color: themeTokens.inputText,
            padding: '0 12px',
            fontSize: '14px',
            outline: 'none',
          }}
        />
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', alignItems: 'center' }}>
          <input
            value={imageAltInput}
            onChange={(event) => setImageAltInput(event.target.value)}
            placeholder={copy.imageAltPlaceholder}
            style={{
              flex: '1 1 220px',
              minWidth: 0,
              height: '38px',
              borderRadius: '8px',
              border: `1px solid ${themeTokens.inputBorder}`,
              background: themeTokens.inputBackground,
              color: themeTokens.inputText,
              padding: '0 12px',
              fontSize: '14px',
              outline: 'none',
            }}
          />
          <button type="button" onClick={onApplyImage} style={textButtonStyle}>
            {copy.insertImage}
          </button>
        </div>
        {canUploadImage ? (
          <div style={{ display: 'grid', gap: '8px' }}>
            <span style={{ color: themeTokens.description, fontSize: '13px', fontWeight: 700 }}>
              {copy.uploadHint || getInlineImageUploadHint()}
            </span>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', alignItems: 'center' }}>
              <label className="lw-point-button-face" style={{ ...textButtonStyle, cursor: editable ? 'pointer' : 'default' }}>
                {copy.chooseFile}
                <input
                  type="file"
                  accept={inlineImageAccept}
                  onChange={(event) => {
                    const file = event.target.files?.[0] ?? null;
                    if (!file) {
                      setSelectedImageFile(null);
                      setImageError('');
                      return;
                    }
                    const validationError = validateInlineImageFile(file, {
                      typeError: copy.inlineImageTypeError,
                      sizeError: copy.inlineImageSizeError,
                    });
                    if (validationError) {
                      event.target.value = '';
                      setSelectedImageFile(null);
                      setImageError(validationError);
                      return;
                    }
                    setSelectedImageFile(file);
                    setImageError('');
                  }}
                  style={{ display: 'none' }}
                />
              </label>
              {selectedImageFile ? (
                <span style={{ flex: '1 1 180px', minWidth: 0, color: themeTokens.inputText, fontSize: '14px', fontWeight: 700, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {selectedImageFile.name}
                </span>
              ) : null}
              <button
                type="button"
                onClick={onUploadImage}
                disabled={!selectedImageFile || isUploadingImage}
                style={{
                  ...textButtonStyle,
                  opacity: selectedImageFile && !isUploadingImage ? 1 : 0.58,
                }}
              >
                {isUploadingImage ? copy.uploading : copy.uploadInsert}
              </button>
            </div>
          </div>
        ) : null}
      </div>
      {imageError ? (
        <span style={{ color: themeTokens.dangerButtonText, fontSize: '13px', fontWeight: 700 }}>
          {imageError}
        </span>
      ) : null}
    </div>
  );
}
