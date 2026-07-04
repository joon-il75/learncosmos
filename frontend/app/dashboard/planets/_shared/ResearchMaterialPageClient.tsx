'use client';

import Link from 'next/link';
import { useCallback, useMemo, useRef, useState } from 'react';
import { useParams, useRouter, useSearchParams } from 'next/navigation';

import { getLearnerHeaderActionsCopy } from '@/components/common/LearnerHeaderActions';
import LumiAvatar from '@/components/lumi/LumiAvatar';
import { getPointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { usePointLearning } from './point/usePointLearning';
import PointWorkspace from './point/PointWorkspace';
import { PointPageHeader } from './point/PointPageHeader';
import { PointLearningToolbar } from './point/PointLearningToolbar';
import { pointToolbarOffsets } from './point/pointGuideCopy';
import { getPointPageThemeTokens, type PointPageTheme } from './pointPageUtils';
import {
  loadingCardStyle,
  loadingPageStyle,
  mainStyle,
  pageBackgroundOverlayStyle,
  pageStyle,
  primaryButtonStyle,
  secondaryButtonStyle,
} from './pointPageStyles';

const POINT_HEADER_FIXED_TOP = 100;
const POINT_TOOLBAR_FLOW_GAP = 12;

export default function ResearchMaterialPageClient() {
  const router = useRouter();
  const params = useParams<{ entry?: string }>();
  const searchParams = useSearchParams();
  const routeEntry = params.entry === 'video' || params.entry === 'content' || params.entry === 'attachments' ? params.entry : null;
  const shouldOpenSelection = searchParams.get('select') === '1';
  const [theme, setTheme] = useState<PointPageTheme>('light');
  const [pointToolbarHeight, setPointToolbarHeight] = useState(() => Number.parseInt(pointToolbarOffsets.collapsed, 10) - POINT_HEADER_FIXED_TOP - POINT_TOOLBAR_FLOW_GAP);
  const [isGuideConfirmed, setIsGuideConfirmed] = useState(() => Boolean(routeEntry) || shouldOpenSelection);
  const workspaceRef = useRef<HTMLDivElement | null>(null);
  const themeTokens = useMemo(() => getPointPageThemeTokens(theme), [theme]);
  const learning = usePointLearning('learning');
  const pointCopy = useMemo(() => getPointLearningCopy(learning.uiLocale), [learning.uiLocale]);
  const learnerHeaderCopy = useMemo(() => getLearnerHeaderActionsCopy(learning.uiLocale), [learning.uiLocale]);
  const pointToolbarOffset = `${POINT_HEADER_FIXED_TOP + pointToolbarHeight + POINT_TOOLBAR_FLOW_GAP}px`;
  const pointDetailHref = learning.planetID && learning.pointID
    ? `/dashboard/planets/learning/${learning.planetID}/points/${learning.pointID}`
    : '/dashboard';
  const planetDetailHref = learning.planetID ? `/dashboard/planets/learning/${learning.planetID}` : '/dashboard';
  const researchMaterialBaseHref = learning.planetID && learning.pointID
    ? '/dashboard/planets/learning/' + learning.planetID + '/points/' + learning.pointID + '/research-material'
    : pointDetailHref;
  const researchMaterialSelectionHref = researchMaterialBaseHref + (researchMaterialBaseHref.includes('?') ? '&' : '?') + 'select=1';
  const [saveSignal, setSaveSignal] = useState(0);
  const [isEntrySaveDisabled, setIsEntrySaveDisabled] = useState(() => routeEntry === 'content');
  const entryPageTitle = routeEntry === 'video'
    ? (learning.uiLocale === 'en' ? 'Research Material Video' : '연구자료 영상')
    : routeEntry === 'content'
      ? (learning.uiLocale === 'en' ? 'Research Material Post' : '연구자료 본문')
      : routeEntry === 'attachments'
        ? (learning.uiLocale === 'en' ? 'Research Material File' : '연구자료 파일')
        : pointCopy.workspace.researchMaterial.pageTitle;
  const entryGuideMessage = routeEntry === 'video'
    ? (learning.uiLocale === 'en' ? 'Upload a video you created for this research material.' : '직접 만든 영상을 연구자료로 올려 주세요.')
    : routeEntry === 'content'
      ? (learning.uiLocale === 'en' ? 'Write the research material body yourself.' : '연구자료 본문을 직접 작성해 주세요.')
      : routeEntry === 'attachments'
        ? (learning.uiLocale === 'en' ? 'Upload a file you created for this research material.' : '직접 만든 파일을 연구자료로 올려 주세요.')
        : pointCopy.workspace.researchMaterial.nextActionNotice;
  const focusedDiaryLessonHref = learning.planetID && learning.pointDetail
    ? `${planetDetailHref}?diaryLessonId=${encodeURIComponent(learning.pointDetail.lessonID)}#planet-route`
    : planetDetailHref;
  const isResearchPoint = learning.pointDetail?.point.point_type === 'research';
  const isPointCompleted = learning.pointDetail?.point.status === 'completed';
  const termsHref = learning.uiLocale === 'en'
    ? '/en/terms#section-6-user-content-and-copyright-responsibility'
    : '/terms#section-7-콘텐츠-및-외부-링크';
  const hasUploadedResearchVideo = learning.attachmentDrafts.some((attachment) => (
    attachment.sourceContext === 'research_material'
    && (attachment.attachmentType === 'video' || attachment.mimeType.toLowerCase().startsWith('video/'))
  ));
  const hasUploadedResearchFile = learning.attachmentDrafts.some((attachment) => (
    attachment.sourceContext === 'research_material'
    && !['video', 'subtitle', 'thumbnail'].includes(attachment.attachmentType)
  ));
  const isBottomConfirmAction = (routeEntry === 'video' && hasUploadedResearchVideo) || (routeEntry === 'attachments' && hasUploadedResearchFile);
  const bottomSecondaryLabel = isBottomConfirmAction
    ? '확인'
    : routeEntry
      ? '취소'
      : '⬅️ 지점 관찰';
  const bottomSecondaryStyle = isBottomConfirmAction
    ? {
        ...secondaryButtonStyle,
        background: '#14B8A6',
        borderColor: '#0F766E',
        color: '#FFFFFF',
      }
    : {
        ...secondaryButtonStyle,
        background: themeTokens.secondaryButtonBackground,
        borderColor: themeTokens.secondaryButtonBorder,
        color: themeTokens.buttonText,
      };
  const shouldShowBottomPrimary = routeEntry === 'content' || !isGuideConfirmed;
  const lumiState = isGuideConfirmed && !routeEntry ? 'planet-point' : 'raise-hand';

  const handleToggleTheme = () => setTheme((current) => (current === 'dark' ? 'light' : 'dark'));
  const handlePointToolbarHeightChange = useCallback((height: number) => {
    setPointToolbarHeight((current) => (current === height ? current : height));
  }, []);
  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined);
    router.replace('/login');
  };
  const handleConfirmGuide = () => {
    setIsGuideConfirmed(true);
    window.requestAnimationFrame(() => {
      workspaceRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    });
  };

  if (learning.isLoading) {
    return (
      <div style={{ ...loadingPageStyle, background: themeTokens.loadingBackground }}>
        <div style={{ ...loadingCardStyle, background: themeTokens.loadingCardBackground, borderColor: themeTokens.loadingCardBorder, color: themeTokens.metaValue }}>
          {pointCopy.loading}
        </div>
      </div>
    );
  }

  if (learning.error || !learning.planet || !learning.pointDetail || !isResearchPoint) {
    return (
      <div style={{ ...loadingPageStyle, background: themeTokens.loadingBackground }}>
        <div style={{ ...loadingCardStyle, display: 'grid', gap: '14px', background: themeTokens.loadingCardBackground, borderColor: themeTokens.loadingCardBorder, color: themeTokens.metaValue }}>
          <div>{learning.error ?? pointCopy.notFound}</div>
          <Link href={routeEntry ? researchMaterialSelectionHref : pointDetailHref} style={{ ...secondaryButtonStyle, background: themeTokens.secondaryButtonBackground, borderColor: themeTokens.secondaryButtonBorder, color: themeTokens.buttonText }}>
            지점 관찰로 돌아가기
          </Link>
        </div>
      </div>
    );
  }

  const pointPathItems = [
    { key: 'system', label: '☀️', href: '/dashboard', title: 'Star System으로 이동' },
    { key: 'diary', label: `📔 ${learning.planet.planet.title}`, href: planetDetailHref, title: '탐험일지로 이동' },
    { key: 'region', label: `🗺️ ${learning.pointDetail.levelTitle}`, href: focusedDiaryLessonHref, title: '탐험일지에서 해당 리슨 보기' },
    ...(learning.pointDetail.lessonTitle !== learning.pointDetail.levelTitle
      ? [{ key: 'subregion', label: `🧭 ${learning.pointDetail.lessonTitle}`, href: focusedDiaryLessonHref, title: '탐험일지에서 해당 리슨 보기' }]
      : []),
  ];
  const currentPointPathItem = { key: 'current', label: `🔬 ${learning.pointDetail.point.title}`, title: '현재 페이지 상단으로 이동' };
  const fullPointPathText = [...pointPathItems.map((item) => item.label), currentPointPathItem.label].join(' > ');

  return (
    <div className={`lw-point-learning-page lw-point-learning-page--${theme}`} style={{ ...pageStyle, background: themeTokens.pageBackground, color: themeTokens.metaValue }}>
      <div style={{ ...pageBackgroundOverlayStyle, background: themeTokens.pageOverlay }} />
      <PointPageHeader
        copy={pointCopy.header}
        theme={theme}
        themeTokens={themeTokens}
        planetDetailHref={planetDetailHref}
        onToggleTheme={handleToggleTheme}
        onLogout={handleLogout}
        learnerHeaderCopy={learnerHeaderCopy}
        locale={learning.uiLocale}
      />
      <PointLearningToolbar
        theme={theme}
        themeTokens={themeTokens}
        pointPathItems={pointPathItems}
        currentPointPathItem={currentPointPathItem}
        fullPointPathText={fullPointPathText}
        onCurrentPathClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        onHeightChange={handlePointToolbarHeightChange}
      />
      <main style={{ ...mainStyle, paddingTop: pointToolbarOffset, paddingBottom: '112px' }}>
        <div style={{ display: 'grid', gap: '16px', maxWidth: '1120px', margin: '0 auto', width: '100%' }}>
          <div style={{ display: 'grid', gap: '16px', minHeight: isGuideConfirmed ? undefined : 'clamp(300px, calc(100dvh - 300px), 520px)', alignContent: 'center' }}>
            <h1 style={{ margin: 0, color: themeTokens.title, fontSize: '24px', fontWeight: 900, lineHeight: 1.2, textAlign: 'center' }}>{entryPageTitle}</h1>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '12px', maxWidth: 'min(560px, 100%)', justifySelf: 'center', width: '100%' }}>
              <LumiAvatar state={lumiState} size={104} />
              <div style={{ position: 'relative', flex: '1 1 0', maxWidth: '360px', minWidth: 0, padding: '12px 14px', borderRadius: '14px', border: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground, color: themeTokens.title, boxShadow: '0 14px 28px rgba(15, 23, 42, 0.12)', textAlign: 'left' }}>
                <span aria-hidden="true" style={{ position: 'absolute', left: '-9px', top: '30px', width: '16px', height: '16px', transform: 'rotate(45deg)', borderLeft: `1px solid ${themeTokens.noticeBorder}`, borderBottom: `1px solid ${themeTokens.noticeBorder}`, background: themeTokens.noticeBackground }} />
                <strong style={{ display: 'block', marginBottom: '5px', color: themeTokens.title, fontSize: '13px', lineHeight: 1.35 }}>Lumi</strong>
                <span style={{ display: 'block', color: themeTokens.mutedText, fontSize: '12.5px', lineHeight: 1.5 }}>{isGuideConfirmed ? entryGuideMessage : pointCopy.workspace.researchMaterial.copyrightNotice}</span>
                {!isGuideConfirmed ? (
                  <a href={termsHref} target="_blank" rel="noopener noreferrer" style={{ display: 'inline-flex', marginTop: '7px', color: themeTokens.buttonText, fontSize: '12.5px', fontWeight: 850, textDecoration: 'underline', textUnderlineOffset: '3px' }}>
                    {pointCopy.workspace.researchMaterial.copyrightAgreementLinkLabel}
                  </a>
                ) : null}
              </div>
            </div>
          </div>
          {isGuideConfirmed ? (
            <div ref={workspaceRef} style={{ scrollMarginTop: pointToolbarOffset }}>
              <PointWorkspace
                copy={pointCopy.workspace}
                learning={learning}
                themeTokens={themeTokens}
                isExplorationPoint={false}
                isSharedRoute={false}
                isLearningContentOpen
                isLearningWorkOpen={false}
                editorToolbarStickyTop={pointToolbarOffset}
                onResearchMaterialConfirmed={() => undefined}
                onLearningWorkGuideChange={() => undefined}
                onSourceContentOpen={() => undefined}
                researchEditModeSignal={isPointCompleted ? 0 : 1}
                showLearningContentTitle={false}
                initialResearchMaterialEntry={routeEntry ?? undefined}
                hideResearchMaterialEntryNav={Boolean(routeEntry)}
                researchMaterialSaveSignal={saveSignal}
                onResearchMaterialEntrySaved={() => router.push(researchMaterialSelectionHref)}
                onResearchMaterialEntrySaveDisabledChange={setIsEntrySaveDisabled}
              />
            </div>
          ) : null}
        </div>
      </main>
      <div style={{ position: "fixed", left: "50%", bottom: "14px", zIndex: 70, width: "min(1120px, calc(100vw - 24px))", transform: "translateX(-50%)", display: "flex", justifyContent: shouldShowBottomPrimary ? "space-between" : "center", gap: "10px", alignItems: "center", padding: "10px", borderRadius: "12px", border: "1px solid " + themeTokens.sectionBorder, background: themeTokens.sectionBackground, boxShadow: "0 16px 36px rgba(15, 23, 42, 0.22)", backdropFilter: "blur(14px)" }}>
        <Link href={routeEntry ? researchMaterialSelectionHref : pointDetailHref} style={{ ...bottomSecondaryStyle, minHeight: "46px", minWidth: shouldShowBottomPrimary ? "min(190px, 48%)" : "min(260px, 100%)", display: "inline-flex", alignItems: "center", justifyContent: "center", gap: "8px", textDecoration: "none" }}>
          {bottomSecondaryLabel}
        </Link>
        {shouldShowBottomPrimary ? (
          <button type="button" disabled={routeEntry === 'content' && isEntrySaveDisabled} onClick={routeEntry ? () => setSaveSignal((current) => current + 1) : handleConfirmGuide} style={{ ...primaryButtonStyle, minHeight: "46px", minWidth: "min(190px, 48%)", display: "inline-flex", alignItems: "center", justifyContent: "center", gap: "8px", background: routeEntry === 'content' && isEntrySaveDisabled ? themeTokens.disabledButtonBackground : "#7F77DD", borderColor: routeEntry === 'content' && isEntrySaveDisabled ? themeTokens.disabledButtonBorder : "rgba(127, 119, 221, 0.76)", color: routeEntry === 'content' && isEntrySaveDisabled ? themeTokens.disabledButtonText : "#FFFDF7", cursor: routeEntry === 'content' && isEntrySaveDisabled ? 'default' : 'pointer', opacity: routeEntry === 'content' && isEntrySaveDisabled ? 0.78 : 1 }}>
            {routeEntry ? '저장' : '확인/다음 ➡️'}
          </button>
        ) : null}
      </div>
    </div>
  );
}
