import type { Dispatch, SetStateAction } from 'react';

import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import type { PointPageThemeTokens } from '../pointPageUtils';
import { createToolbarTextButtonStyle } from './tiptapEditorStyles';

type Props = {
  linkInput: string;
  setLinkInput: Dispatch<SetStateAction<string>>;
  linkError: string;
  setLinkError: Dispatch<SetStateAction<string>>;
  isLinkActive: boolean;
  editable: boolean;
  themeTokens: PointPageThemeTokens;
  copy: PointLearningCopy['workspace']['editor'];
  onApplyLink: () => void;
  onUnsetLink: () => void;
  onClose: () => void;
};

export default function TiptapLinkPanel({
  linkInput,
  setLinkInput,
  linkError,
  setLinkError,
  isLinkActive,
  editable,
  themeTokens,
  copy,
  onApplyLink,
  onUnsetLink,
  onClose,
}: Props) {
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
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', alignItems: 'center' }}>
        <input
          value={linkInput}
          onChange={(event) => {
            setLinkInput(event.target.value);
            setLinkError('');
          }}
          onKeyDown={(event) => {
            if (event.key === 'Enter') {
              event.preventDefault();
              onApplyLink();
            }
            if (event.key === 'Escape') {
              event.preventDefault();
              onClose();
            }
          }}
          placeholder="https://..."
          style={{
            flex: '1 1 240px',
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
        <button type="button" onClick={onApplyLink} style={createToolbarTextButtonStyle(themeTokens, editable)}>
          {copy.apply}
        </button>
        <button
          type="button"
          onClick={onUnsetLink}
          disabled={!isLinkActive}
          style={{
            ...createToolbarTextButtonStyle(themeTokens, editable),
            opacity: isLinkActive ? 1 : 0.58,
          }}
        >
          {copy.unset}
        </button>
      </div>
      {linkError ? (
        <span style={{ color: themeTokens.dangerButtonText, fontSize: '13px', fontWeight: 700 }}>
          {linkError}
        </span>
      ) : null}
    </div>
  );
}
