import type { ChangeEvent, RefObject } from 'react';
import {
  fieldLabelStyle,
  fieldStyle,
  fileDropStyle,
  fileHelpStyle,
  fileInputStyle,
  inputStyle,
  specGridStyle,
  specItemStyle,
  specLabelStyle,
  specValueStyle,
  textareaStyle,
} from './planetTextureMapStyles';

type TextureMapUploadFieldsProps = {
  atlasIsRecommended: boolean | null | undefined;
  description: string;
  fileInputRef: RefObject<HTMLInputElement | null>;
  fileName: string | null;
  handleFileChange: (event: ChangeEvent<HTMLInputElement>) => void;
  imageSize: { width: number; height: number } | null;
  name: string;
  setDescription: (value: string) => void;
  setName: (value: string) => void;
};

export function TextureMapUploadFields({
  atlasIsRecommended,
  description,
  fileInputRef,
  fileName,
  handleFileChange,
  imageSize,
  name,
  setDescription,
  setName,
}: TextureMapUploadFieldsProps) {
  return (
    <>
      <label style={fieldStyle}>
        <span style={fieldLabelStyle}>이름</span>
        <input value={name} onChange={(event) => setName(event.target.value)} style={inputStyle} />
      </label>
      <label style={fieldStyle}>
        <span style={fieldLabelStyle}>설명</span>
        <textarea value={description} onChange={(event) => setDescription(event.target.value)} rows={4} style={textareaStyle} />
      </label>
      <label style={fileDropStyle}>
        <span style={fieldLabelStyle}>WebP atlas 선택</span>
        <input ref={fileInputRef} type="file" accept=".png,.jpg,.jpeg,.webp,image/png,image/jpeg,image/webp" onChange={handleFileChange} style={fileInputStyle} />
        <span style={fileHelpStyle}>{fileName ?? '2048x768 PNG/JPEG/WebP 파일을 선택해 주세요. 저장은 WebP로 변환됩니다.'}</span>
      </label>

      <div style={specGridStyle}>
        <SpecItem label="권장 크기" value="2048 × 768" active={Boolean(atlasIsRecommended || !imageSize)} />
        <SpecItem label="Grid" value="4열 × 3행" active />
        <SpecItem label="Cell" value="512 × 256" active />
        <SpecItem
          label="선택 파일"
          value={imageSize ? `${imageSize.width} × ${imageSize.height}` : '대기 중'}
          active={Boolean(atlasIsRecommended)}
        />
      </div>
    </>
  );
}

function SpecItem({ label, value, active }: { label: string; value: string; active: boolean }) {
  return (
    <div style={specItemStyle(active)}>
      <span style={specLabelStyle}>{label}</span>
      <strong style={specValueStyle}>{value}</strong>
    </div>
  );
}
