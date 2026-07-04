'use client';

import { type FormEvent, type RefObject } from 'react';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import {
  creatorSectionStyle,
  creatorSectionInnerStyle,
  creatorHeroSloganStyle,
  creatorLabelStyle,
  creatorDescriptionStyle,
  creatorFormStyle,
  creatorInputStyle,
  creatorButtonStyle,
  creatorErrorStyle,
  queryInputTooltipStyle,
  queryInputTooltipTextStyle,
  chipRowStyle,
  chipStyle,
} from './_styles';

interface CTAViewModel {
  title: string;
  placeholder: string;
  submitLabel: string;
}

interface Props {
  ctaQuery: string;
  isQueryInputFocused: boolean;
  isCreatingCourse: boolean;
  creationError: string | null;
  ctaViewModel: CTAViewModel;
  copy: DashboardMainCopy['cta'];
  isCompletedOnly: boolean;
  queryInputRef: RefObject<HTMLInputElement | null>;
  onQueryChange: (value: string) => void;
  onFocusChange: (focused: boolean) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onChipClick: (chip: string) => void;
}

export default function DashboardCTASection({
  ctaQuery,
  isQueryInputFocused,
  isCreatingCourse,
  creationError,
  ctaViewModel,
  copy,
  isCompletedOnly,
  queryInputRef,
  onQueryChange,
  onFocusChange,
  onSubmit,
  onChipClick,
}: Props) {
  const sectionStyle = isCompletedOnly
    ? {
        ...creatorSectionStyle,
        border: '1px solid rgba(255, 205, 96, 0.42)',
        borderLeft: '2px solid rgba(255, 196, 69, 0.92)',
        background: 'linear-gradient(135deg, rgba(30, 16, 4, 0.92), rgba(11, 18, 33, 0.88))',
        boxShadow: '0 26px 90px rgba(0, 0, 0, 0.42), 0 0 0 1px rgba(255, 196, 69, 0.12) inset',
      }
    : creatorSectionStyle;

  const slogan = isCompletedOnly ? copy.slogan.completedOnly : copy.slogan.default;
  const description = isCompletedOnly ? copy.description.completedOnly : copy.description.default;

  return (
    <section style={sectionStyle}>
      <div style={creatorSectionInnerStyle}>
        <div style={{ display: 'grid', gap: '10px' }}>
          <div style={creatorHeroSloganStyle}>{slogan}</div>
          <div style={creatorLabelStyle}>{ctaViewModel.title}</div>
          <div style={creatorDescriptionStyle}>{description}</div>
          <div style={{ position: 'relative' }}>
            <form onSubmit={onSubmit} style={creatorFormStyle}>
              <input
                ref={queryInputRef}
                type="text"
                value={ctaQuery}
                onChange={(e) => onQueryChange(e.target.value)}
                onFocus={() => onFocusChange(true)}
                onBlur={() => onFocusChange(false)}
                placeholder={ctaViewModel.placeholder}
                style={creatorInputStyle}
              />
              <button
                type="submit"
                style={creatorButtonStyle(isCreatingCourse)}
                disabled={isCreatingCourse}
              >
                {isCreatingCourse ? copy.submittingLabel : ctaViewModel.submitLabel}
              </button>
            </form>
            {isQueryInputFocused ? (
              <div style={queryInputTooltipStyle}>
                <LumiAvatar state="raise-hand" size={28} />
                <span style={queryInputTooltipTextStyle}>{copy.inputHelp}</span>
              </div>
            ) : null}
          </div>
          <div
            style={{
              minHeight: isQueryInputFocused ? '56px' : creationError ? '22px' : 0,
              paddingTop: isQueryInputFocused ? '38px' : 0,
              transition: 'min-height 0.18s ease, padding-top 0.18s ease',
            }}
          >
            {creationError ? <div style={creatorErrorStyle}>{creationError}</div> : null}
          </div>
          <div style={chipRowStyle}>
            {copy.topicChips.map((chip) => (
              <button key={chip} type="button" style={chipStyle} onClick={() => onChipClick(chip)}>
                {chip}
              </button>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
