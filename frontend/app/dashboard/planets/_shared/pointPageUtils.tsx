import type { ReactNode } from 'react';
import type {
  PlanetMetaState, PlanetRouteKind, PlanetStatus, ResolvedPoint, ResearchBlockType,
  PlanetPointDetailResponse, PlanetPoint, PointQuestionType, PointQuestionStatus
} from './pointPageTypes';
import {
  sectionStyle, sectionTitleStyle, sectionSubtitleStyle, pointShellCardStyle,
} from './pointPageStyles';

export function getPlanetDetailHref(routeKind: PlanetRouteKind, planetID: string): string {
  return `/dashboard/planets/${routeKind}/${planetID}`;
}

export function getPointHref(routeKind: PlanetRouteKind, planetID: string, pointID: string): string {
  return `/dashboard/planets/${routeKind}/${planetID}/points/${pointID}`;
}

export function mapPointDetailResponseToResolvedPoint(point: PlanetPointDetailResponse): ResolvedPoint {
  return {
    levelTitle: point.level_title,
    lessonTitle: point.lesson_title,
    lessonID: point.lesson_id,
    previousPointID: point.previous_point_id ?? null,
    nextPointID: point.next_point_id ?? null,
    point: point.point.point,
    research_material_confirmed: Boolean(point.point.research_material_confirmed),
    research_material_confirmed_at: point.point.research_material_confirmed_at ?? null,
    journal_entry: point.point.journal_entry ?? null,
    record_entry: point.point.record_entry ?? null,
    artifact_entry: point.point.artifact_entry ?? null,
    artifacts: point.point.artifacts ?? [],
    attachments: point.point.attachments ?? [],
    material_reports: point.point.material_reports ?? [],
    practice_logs: point.point.practice_logs ?? [],
    events: point.point.events ?? [],
    observation_notes: point.point.observation_notes ?? [],
    ai_summary_entry: point.point.ai_summary_entry ?? null,
    blocks: point.point.blocks ?? [],
    questions: point.point.questions ?? [],
    self_evaluation: point.point.self_evaluation ?? null,
    template_type: point.point.point.template_type,
  };
}

export function mapQuestionToDraft(question: NonNullable<PlanetPoint['questions']>[number]): PointQuestionDraft {
  return {
    id: question.id,
    title: question.title ?? '',
    question: question.question,
    questionType: question.question_type,
    answerMethod: question.answer_method ?? 'self',
    answer: question.answer ?? '',
    aiFeedback: question.ai_feedback ?? '',
    status: question.status,
    createdBy: question.created_by,
    createdAt: question.created_at,
    updatedAt: question.updated_at,
  };
}

export function getPointStatusLabel(status: PlanetPoint['point']['status']): string {
  if (status === 'completed') return '완료';
  if (status === 'learning') return '진행 중';
  if (status === 'ready') return '준비됨';
  return '초안';
}

export function getPlanetStatusLabel(status: PlanetStatus): string {
  if (status === 'ready') return '학습준비중';
  if (status === 'learning') return '학습중';
  return '공유중';
}

export function compactGoalMeta(goalContext: PlanetMetaState['goal_context']): string[] {
  if (!goalContext) return [];
  const items: string[] = [];
  if (goalContext.usage_context?.trim()) {
    items.push(`활용 맥락: ${goalContext.usage_context.trim()}`);
  }
  if (goalContext.motivation?.trim()) {
    items.push(`동기: ${goalContext.motivation.trim()}`);
  }
  if (goalContext.goal_profile_version != null) {
    items.push(`goal version ${goalContext.goal_profile_version}`);
  }
  return items;
}

export function goalAwarePlaceholder(base: string, goal: string | null): string {
  if (!goal) return base;
  const shortGoal = goal.length > 30 ? `${goal.slice(0, 30)}...` : goal;
  return `'${shortGoal}'을 위해 ${base}`;
}

export interface ResearchBlockDraft {
  id: string;
  blockType: ResearchBlockType;
  orderIndex: number;
  text: string;
  url: string;
  title: string;
  caption: string;
  note: string;
  isNew?: boolean;
  isDirty?: boolean;
  updatedAt?: string;
}

export interface PointArtifactDraft {
  id: string;
  artifactType: string;
  title: string;
  url: string;
  description: string;
  pointCategory: string;
  productionProcess: string;
  learnedPoints: string;
  difficultPoints: string;
  visibility: string;
  orderIndex: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface PointAttachmentDraft {
  id: string;
  provider: 'creator' | 'platform' | 'learner';
  attachmentType: 'image' | 'file' | 'link' | 'code' | 'other' | 'video' | 'subtitle' | 'thumbnail';
  sourceContext: 'research_material' | 'work_attachment' | 'artifact';
  artifactID?: string | null;
  title: string;
  url: string;
  filePath: string;
  fileSize: string;
  mimeType: string;
  isDirty?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface PointPracticeLogDraft {
  id: string;
  title: string;
  activityName: string;
  attemptCount: string;
  successCount: string;
  failureCount: string;
  durationMinutes: string;
  blockedPart: string;
  changedMethod: string;
  achievementNote: string;
  achievement: string;
  nextPlan: string;
  nextPractice: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface PointQuestionDraft {
  id: string;
  title: string;
  question: string;
  questionType: PointQuestionType;
  answerMethod: string;
  answer: string;
  aiFeedback: string;
  status: PointQuestionStatus;
  createdBy: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface PointCompletionReadiness {
  answeredQuestionCount: number;
  hasSelfEvaluation: boolean;
  hasSelfEvaluationQuality: boolean;
  hasGoalConnection: boolean;
  hasResearchMaterialConfirmed: boolean;
  isReady: boolean;
}

function asTextValue(value: unknown): string {
  return typeof value === 'string' ? value : '';
}

export function mapBlockToDraft(
  block: NonNullable<PlanetPoint['blocks']>[number],
): ResearchBlockDraft {
  const content = block.content ?? {};
  return {
    id: block.id,
    blockType: block.block_type,
    orderIndex: block.order_index,
    text: asTextValue(content.text),
    url: asTextValue(content.url),
    title: asTextValue(content.title),
    caption: asTextValue(content.caption),
    note: asTextValue(content.note),
    isNew: false,
    isDirty: false,
    updatedAt: block.updated_at,
  };
}

export function createEmptyBlockDraft(blockType: ResearchBlockType, orderIndex: number): ResearchBlockDraft {
  return {
    id: `draft-${blockType}-${Date.now()}-${orderIndex}`,
    blockType,
    orderIndex,
    text: '',
    url: '',
    title: '',
    caption: '',
    note: '',
    isNew: true,
    isDirty: true,
  };
}

export function getResearchBlockPlainText(value: string): string {
  return value
    .replace(/<script[\s\S]*?<\/script>/gi, ' ')
    .replace(/<style[\s\S]*?<\/style>/gi, ' ')
    .replace(/<[^>]+>/g, ' ')
    .replace(/&nbsp;|&#160;/gi, ' ')
    .replace(/\s+/g, ' ')
    .trim();
}

export function hasResearchBlockBodyContent(block: Pick<ResearchBlockDraft, 'blockType' | 'text' | 'url' | 'caption' | 'note'>): boolean {
  if (block.blockType === 'text') return getResearchBlockPlainText(block.text).length > 0;
  if (block.blockType === 'image') return block.url.trim().length > 0 || block.caption.trim().length > 0;
  return block.url.trim().length > 0 || block.note.trim().length > 0;
}

export function mapArtifactToDraft(
  artifact: NonNullable<PlanetPoint['artifacts']>[number],
): PointArtifactDraft {
  const artifactType = (() => {
    if (artifact.artifact_type === 'note' || artifact.artifact_type === '텍스트') return '문서 편집자료';
    if (artifact.artifact_type === '동영상') return '영상자료';
    if (['링크', '파일', '이미지'].includes(artifact.artifact_type)) return '첨부자료';
    return artifact.artifact_type;
  })();

  return {
    id: artifact.id,
    artifactType,
    title: artifact.title,
    url: artifact.url,
    description: artifact.description,
    pointCategory: artifact.point_category ?? '',
    productionProcess: artifact.production_process ?? '',
    learnedPoints: artifact.learned_points ?? '',
    difficultPoints: artifact.difficult_points ?? '',
    visibility: artifact.visibility ?? 'private',
    orderIndex: artifact.order_index,
    createdAt: artifact.created_at,
    updatedAt: artifact.updated_at,
  };
}

export function mapAttachmentToDraft(
  attachment: NonNullable<PlanetPoint['attachments']>[number],
): PointAttachmentDraft {
  return {
    id: attachment.id,
    provider: attachment.provider,
    attachmentType: attachment.attachment_type,
    sourceContext: attachment.source_context ?? 'work_attachment',
    artifactID: attachment.artifact_id ?? null,
    title: attachment.title,
    url: attachment.url,
    filePath: attachment.file_path,
    fileSize: optionalNumberInput(attachment.file_size),
    mimeType: attachment.mime_type,
    isDirty: false,
    createdAt: attachment.created_at,
    updatedAt: attachment.updated_at,
  };
}

export function optionalNumberInput(value?: number | null): string {
  return value == null ? '' : String(value);
}

export function nullableNumberFromInput(value: string): number | null {
  const trimmed = value.trim();
  if (!trimmed) return null;
  const parsed = Number(trimmed);
  return Number.isFinite(parsed) ? Math.max(0, Math.floor(parsed)) : null;
}

export function mapPracticeLogToDraft(
  log: NonNullable<PlanetPoint['practice_logs']>[number],
): PointPracticeLogDraft {
  return {
    id: log.id,
    title: log.title,
    activityName: log.activity_name ?? log.title,
    attemptCount: optionalNumberInput(log.attempt_count),
    successCount: optionalNumberInput(log.success_count),
    failureCount: optionalNumberInput(log.failure_count),
    durationMinutes: optionalNumberInput(log.duration_minutes),
    blockedPart: log.blocked_part,
    changedMethod: log.changed_method,
    achievementNote: log.achievement_note,
    achievement: log.achievement ?? log.achievement_note,
    nextPlan: log.next_plan,
    nextPractice: log.next_practice ?? log.next_plan,
    createdAt: log.created_at,
    updatedAt: log.updated_at,
  };
}

export function createEmptyPracticeLogDraft(): PointPracticeLogDraft {
  return {
    id: '',
    title: '',
    activityName: '',
    attemptCount: '',
    successCount: '',
    failureCount: '',
    durationMinutes: '',
    blockedPart: '',
    changedMethod: '',
    achievementNote: '',
    achievement: '',
    nextPlan: '',
    nextPractice: '',
  };
}

export function buildAttachmentPayload(attachment: PointAttachmentDraft) {
  return {
    provider: attachment.provider,
    attachment_type: attachment.attachmentType,
    source_context: attachment.sourceContext,
    artifact_id: attachment.artifactID ?? null,
    title: attachment.title,
    url: attachment.url,
    file_path: attachment.filePath,
    file_size: nullableNumberFromInput(attachment.fileSize),
    mime_type: attachment.mimeType,
  };
}

export function buildPracticeLogPayload(log: PointPracticeLogDraft) {
  return {
    title: log.title,
    activity_name: log.activityName || log.title,
    attempt_count: nullableNumberFromInput(log.attemptCount),
    success_count: nullableNumberFromInput(log.successCount),
    failure_count: nullableNumberFromInput(log.failureCount),
    duration_minutes: nullableNumberFromInput(log.durationMinutes),
    blocked_part: log.blockedPart,
    changed_method: log.changedMethod,
    achievement_note: log.achievementNote,
    achievement: log.achievement || log.achievementNote,
    next_plan: log.nextPlan,
    next_practice: log.nextPractice || log.nextPlan,
  };
}

export function buildResearchBlockContent(block: ResearchBlockDraft): Record<string, string> {
  if (block.blockType === 'text') {
    return {
      title: block.title,
      text: block.text,
    };
  }
  if (block.blockType === 'image') {
    return {
      url: block.url,
      caption: block.caption,
    };
  }
  return {
    url: block.url,
    title: block.title,
    note: block.note,
  };
}

export function getPointCompletionReadiness(
  questions: PointQuestionDraft[],
  selfEvaluation: PlanetPoint['self_evaluation'],
  point?: ResolvedPoint | null,
): PointCompletionReadiness {
  const answeredQuestionCount = questions.filter((question) => question.answer.trim().length > 0).length;
  const hasSelfEvaluation = Boolean(selfEvaluation);
  const hasGoalConnection = Boolean(selfEvaluation && selfEvaluation.goal_alignment_note.trim().length > 0);
  const hasSelfEvaluationQuality = isSelfEvaluationQualityReady(selfEvaluation);
  const hasResearchMaterialConfirmed = point?.point.point_type === 'research'
    ? Boolean(point.research_material_confirmed)
    : true;

  return {
    answeredQuestionCount,
    hasSelfEvaluation,
    hasSelfEvaluationQuality,
    hasGoalConnection,
    hasResearchMaterialConfirmed,
    isReady: answeredQuestionCount > 0 && hasSelfEvaluationQuality && hasGoalConnection && hasResearchMaterialConfirmed,
  };
}

function isSelfEvaluationQualityReady(selfEvaluation: PlanetPoint['self_evaluation']): boolean {
  if (!selfEvaluation) return false;
  const rubricItems = [
    { score: selfEvaluation.understanding_score, reason: selfEvaluation.understanding_reason },
    { score: selfEvaluation.application_score, reason: selfEvaluation.application_reason },
    { score: selfEvaluation.proficiency_score, reason: selfEvaluation.proficiency_reason },
    { score: selfEvaluation.problem_solving_score, reason: selfEvaluation.problem_solving_reason },
    { score: selfEvaluation.expression_score, reason: selfEvaluation.expression_reason },
  ];
  return rubricItems.every((item) => isSelfEvaluationScoreReady(item.score) && (item.reason ?? '').trim().length > 0)
    && selfEvaluation.goal_alignment_note.trim().length > 0;
}

function isSelfEvaluationScoreReady(score: number | null | undefined): boolean {
  return typeof score === 'number' && score >= 1 && score <= 5;
}

export function getQuestionAuthorLabel(createdBy: string): string {
  if (createdBy === 'ai') return 'AI 질문';
  return '학습자 질문';
}

function getPointEventPayloadTitle(payload?: Record<string, unknown> | null): string {
  const value = payload?.title ?? payload?.question ?? payload?.material_title ?? payload?.attachment_title ?? payload?.artifact_title ?? payload?.practice_log_title;
  return typeof value === 'string' ? value.trim() : '';
}

function withPointEventTitle(label: string, payload?: Record<string, unknown> | null): string {
  const title = getPointEventPayloadTitle(payload);
  return title ? `${title} ${label}` : label;
}

function getPointEventAction(payload?: Record<string, unknown> | null): 'created' | 'updated' | 'deleted' {
  const action = typeof payload?.action === 'string' ? payload.action.trim().toLowerCase() : '';
  if (action === 'updated') return 'updated';
  if (action === 'deleted') return 'deleted';
  return 'created';
}

export function getPointEventLabel(eventType: string, payload?: Record<string, unknown> | null): string {
  const action = getPointEventAction(payload);
  switch (eventType) {
    case 'point_started':
      return '지점 학습을 시작했습니다.';
    case 'point_goal_saved':
      return '탐험목표를 저장했습니다.';
    case 'material_viewed':
      return '학습자료를 확인했습니다.';
    case 'journal_saved':
      return '탐험일지를 저장했습니다.';
    case 'record_saved':
      return '학습기록을 저장했습니다.';
    case 'question_added':
      if (action === 'deleted') return withPointEventTitle('질문을 삭제했습니다.', payload);
      return withPointEventTitle('질문을 추가했습니다.', payload);
    case 'question_answered':
      return withPointEventTitle('질문 답변을 저장했습니다.', payload);
    case 'ai_hint_used':
      return 'AI 참고 도움을 사용했습니다.';
    case 'practice_log_added':
      if (action === 'updated') return withPointEventTitle('연습/활동기록을 수정했습니다.', payload);
      if (action === 'deleted') return withPointEventTitle('연습/활동기록을 삭제했습니다.', payload);
      return withPointEventTitle('연습/활동기록을 남겼습니다.', payload);
    case 'artifact_submitted':
      if (action === 'updated') return withPointEventTitle('탐험결과물을 수정했습니다.', payload);
      if (action === 'deleted') return withPointEventTitle('탐험결과물을 삭제했습니다.', payload);
      return withPointEventTitle('탐험결과물을 저장했습니다.', payload);
    case 'attachment_added':
      if (payload?.source_context === 'research_material') {
        if (action === 'updated') return withPointEventTitle('학습대상자료를 수정했습니다.', payload);
        if (action === 'deleted') return withPointEventTitle('학습대상자료를 삭제했습니다.', payload);
        return withPointEventTitle('학습대상자료를 저장했습니다.', payload);
      }
      if (action === 'updated') return withPointEventTitle('보조자료를 수정했습니다.', payload);
      if (action === 'deleted') return withPointEventTitle('보조자료를 삭제했습니다.', payload);
      return withPointEventTitle('보조자료를 저장했습니다.', payload);
    case 'research_material_saved':
      if (action === 'updated') return withPointEventTitle('학습대상자료를 수정했습니다.', payload);
      if (action === 'deleted') return withPointEventTitle('학습대상자료를 삭제했습니다.', payload);
      return withPointEventTitle('학습대상자료를 저장했습니다.', payload);
    case 'research_material_confirmed':
      return '학습대상자료를 확정했습니다.';
    case 'material_reported':
      return '자료 오류를 신고했습니다.';
    case 'self_evaluation_saved':
      return '자기평가를 저장했습니다.';
    case 'point_completed':
      return '지점을 완료했습니다.';
    case 'point_reopened':
      return '지점을 다시 진행 중으로 열었습니다.';
    default:
      return '학습 활동이 기록되었습니다.';
  }
}

export function formatPointEventTime(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString('ko-KR', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function getResearchTemplateLabel(templateType?: ResolvedPoint['template_type']): string {
  switch (templateType) {
    case 'concept_summary':
      return '개념 정리';
    case 'practice_strategy':
      return '연습 전략';
    case 'problem_solving':
      return '문제 해결';
    case 'free_research':
      return '자유 연구';
    default:
      return '연구 기록';
  }
}

export type PointPageTheme = 'dark' | 'light';

export interface PointPageThemeTokens {
  pageBackground: string;
  pageOverlay: string;
  loadingBackground: string;
  loadingCardBackground: string;
  loadingCardBorder: string;
  navBackground: string;
  navBorder: string;
  headerMetaSlash: string;
  headerMetaText: string;
  heroBackground: string;
  heroBorder: string;
  heroAccentBorder: string;
  sectionBackground: string;
  sectionBorder: string;
  surfaceBackground: string;
  surfaceBorder: string;
  placeholderBackground: string;
  placeholderBorder: string;
  emptyBackground: string;
  emptyBorder: string;
  noticeBackground: string;
  noticeBorder: string;
  title: string;
  subtitle: string;
  description: string;
  mutedText: string;
  metaLabel: string;
  metaValue: string;
  inputBackground: string;
  inputBorder: string;
  inputText: string;
  buttonText: string;
  secondaryButtonBackground: string;
  secondaryButtonBorder: string;
  primaryButtonBackground: string;
  primaryButtonBorder: string;
  disabledButtonBackground: string;
  disabledButtonBorder: string;
  disabledButtonText: string;
  dangerButtonBackground: string;
  dangerButtonBorder: string;
  dangerButtonText: string;
}

export function NodeSection({
  id,
  title,
  subtitle,
  tokens,
  children,
  scrollMarginTop,
  showHeader = true,
}: {
  id?: string;
  title: string;
  subtitle?: string;
  tokens: PointPageThemeTokens;
  children: ReactNode;
  scrollMarginTop?: string;
  showHeader?: boolean;
}) {
  return (
    <section id={id} style={{ ...sectionStyle, scrollMarginTop, background: tokens.sectionBackground, borderColor: tokens.sectionBorder }}>
      {showHeader ? (
        <div style={{ display: 'flex', flexWrap: 'wrap', alignItems: 'baseline', columnGap: '10px', rowGap: '4px' }}>
          <h2 style={{ ...sectionTitleStyle, color: tokens.title }}>{title}</h2>
          {subtitle ? <span style={{ ...sectionSubtitleStyle, margin: 0, color: tokens.mutedText, fontSize: '13px', lineHeight: 1.35 }}>{subtitle}</span> : null}
        </div>
      ) : null}
      {children}
    </section>
  );
}

export function NodeSectionCard({
  tokens,
  children,
}: {
  tokens: PointPageThemeTokens;
  children: ReactNode;
}) {
  return (
    <div style={{ ...pointShellCardStyle, background: tokens.surfaceBackground, borderColor: tokens.surfaceBorder }}>
      {children}
    </div>
  );
}

export function getPointPageThemeTokens(theme: PointPageTheme): PointPageThemeTokens {
  if (theme === 'light') {
    return {
      pageBackground: '#F8FAFC',
      pageOverlay: 'linear-gradient(rgba(248, 250, 252, 0.96), rgba(248, 250, 252, 0.96))',
      loadingBackground: '#F8FAFC',
      loadingCardBackground: '#FFFFFF',
      loadingCardBorder: '#CBD5E1',
      navBackground: 'rgba(255, 255, 255, 0.96)',
      navBorder: '#E2E8F0',
      headerMetaSlash: '#64748B',
      headerMetaText: '#0F172A',
      heroBackground: '#FFFFFF',
      heroBorder: '#D8E0EA',
      heroAccentBorder: '#111827',
      sectionBackground: '#FFFFFF',
      sectionBorder: '#D8E0EA',
      surfaceBackground: '#F8FAFC',
      surfaceBorder: '#D8E0EA',
      placeholderBackground: '#FFFFFF',
      placeholderBorder: '#CBD5E1',
      emptyBackground: '#F8FAFC',
      emptyBorder: '#D8E0EA',
      noticeBackground: '#F1F5F9',
      noticeBorder: '#CBD5E1',
      title: '#020617',
      subtitle: '#111827',
      description: '#111827',
      mutedText: '#334155',
      metaLabel: '#475569',
      metaValue: '#0F172A',
      inputBackground: '#FFFFFF',
      inputBorder: '#CBD5E1',
      inputText: '#020617',
      buttonText: '#020617',
      secondaryButtonBackground: '#FFFFFF',
      secondaryButtonBorder: '#CBD5E1',
      primaryButtonBackground: '#E5E7EB',
      primaryButtonBorder: '#111827',
      disabledButtonBackground: '#E5E7EB',
      disabledButtonBorder: '#CBD5E1',
      disabledButtonText: '#64748B',
      dangerButtonBackground: '#FEF2F2',
      dangerButtonBorder: '#FCA5A5',
      dangerButtonText: '#7F1D1D',
    };
  }

  return {
    pageBackground: '#0B0F14',
    pageOverlay: 'linear-gradient(rgba(11, 15, 20, 0.96), rgba(11, 15, 20, 0.96))',
    loadingBackground: '#0B0F14',
    loadingCardBackground: '#111827',
    loadingCardBorder: '#263244',
    navBackground: 'rgba(17, 24, 39, 0.96)',
    navBorder: '#263244',
    headerMetaSlash: '#CBD5E1',
    headerMetaText: '#F9FAFB',
    heroBackground: '#111827',
    heroBorder: '#263244',
    heroAccentBorder: '#CBD5E1',
    sectionBackground: '#111827',
    sectionBorder: '#263244',
    surfaceBackground: '#0F172A',
    surfaceBorder: '#263244',
    placeholderBackground: '#0F172A',
    placeholderBorder: '#263244',
    emptyBackground: '#0F172A',
    emptyBorder: '#263244',
    noticeBackground: '#111C2D',
    noticeBorder: '#263244',
    title: '#F9FAFB',
    subtitle: '#E5E7EB',
    description: '#E5E7EB',
    mutedText: '#CBD5E1',
    metaLabel: '#CBD5E1',
    metaValue: '#F9FAFB',
    inputBackground: '#0B1220',
    inputBorder: '#263244',
    inputText: '#F9FAFB',
    buttonText: '#F9FAFB',
    secondaryButtonBackground: '#111827',
    secondaryButtonBorder: '#334155',
    primaryButtonBackground: '#1F2937',
    primaryButtonBorder: '#CBD5E1',
    disabledButtonBackground: '#1F2937',
    disabledButtonBorder: '#263244',
    disabledButtonText: '#94A3B8',
    dangerButtonBackground: '#3F121A',
    dangerButtonBorder: '#7F1D1D',
    dangerButtonText: '#FECACA',
  };
}
