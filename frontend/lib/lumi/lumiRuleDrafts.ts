import type { LumiAction, LumiMode, LumiPage, LumiRuntimeContext, LumiScene } from './lumiEngineTypes';
import type { LumiDockSlot, LumiMessageType, LumiPageContext, LumiState, LumiTrigger, LumiTriggerPayload } from './lumiTypes';

export interface LumiRuleDraft {
  id: string;
  label: string;
  trigger: LumiTrigger;
  mode: LumiMode;
  page: LumiPage;
  action: LumiAction;
  scene: LumiScene;
  context: LumiPageContext;
  dockSlot: LumiDockSlot;
  state: LumiState;
  messageType: LumiMessageType;
  message: string;
  note: string;
}

export const lumiRuleDraftPresets: LumiRuleDraft[] = [
  {
    id: 'hero-focus',
    label: 'Hero Focus',
    trigger: 'cta_focus',
    mode: 'general',
    page: 'hero',
    action: 'open_lumi',
    scene: 'editing',
    context: 'dashboard',
    dockSlot: 'dashboard-cta',
    state: 'raise-hand',
    messageType: 'question',
    message: '배우고 싶은 주제를 편하게 적어 주세요.',
    note: '랜딩 Hero 입력창 포커스 시 기본 질문 상태입니다.',
  },
  {
    id: 'hero-generate',
    label: 'Hero Generate',
    trigger: 'generation_loading',
    mode: 'ai',
    page: 'hero',
    action: 'ai_generate',
    scene: 'waiting',
    context: 'dashboard',
    dockSlot: 'dashboard-cta',
    state: 'thinking',
    messageType: 'guide',
    message: '첫 번째 탐험 경로를 차근차근 정리하고 있어요.\n행성도 함께 준비하고 있어요.',
    note: '랜딩에서 행성 생성 요청 직후 대기 상태입니다.',
  },
  {
    id: 'draft-created',
    label: 'Draft Created',
    trigger: 'generation_success',
    mode: 'ai',
    page: 'course_draft_editor',
    action: 'create_course',
    scene: 'success',
    context: 'diary',
    dockSlot: 'diary-header',
    state: 'planet-hold',
    messageType: 'success',
    message: '새로운 행성이 준비됐어요.\n이제 첫 탐험 경로를 함께 다듬어 볼까요?',
    note: '초안 구조 생성 직후 상단 안내로 쓰는 규칙입니다.',
  },
  {
    id: 'draft-search',
    label: 'Draft Search',
    trigger: 'search_start',
    mode: 'ai',
    page: 'course_draft_editor',
    action: 'search_content',
    scene: 'searching',
    context: 'diary',
    dockSlot: 'diary-sidebar',
    state: 'exploring',
    messageType: 'discovery',
    message: '학습 자료를 탐험하고 있어요.\n이 단계에 잘 맞는 자료를 차근차근 찾고 있어요.',
    note: '자료 재탐색이나 추천 검색이 시작될 때의 규칙입니다.',
  },
  {
    id: 'player-lesson',
    label: 'Player Lesson',
    trigger: 'lesson_enter',
    mode: 'general',
    page: 'course_player',
    action: 'open_lumi',
    scene: 'overview',
    context: 'player',
    dockSlot: 'player-inline',
    state: 'note-read',
    messageType: 'guide',
    message: '이번 탐험에서는 핵심부터 차근차근 살펴볼게요.',
    note: '플레이어 진입면이 준비되면 lesson enter 기본 규칙으로 연결할 항목입니다.',
  },
  {
    id: 'player-complete',
    label: 'Player Complete',
    trigger: 'course_complete',
    mode: 'general',
    page: 'course_player',
    action: 'course_complete',
    scene: 'success',
    context: 'completion',
    dockSlot: 'completion-center',
    state: 'victory',
    messageType: 'success',
    message: '이 행성 탐험을 잘 마쳤어요.\n정말 인상적인 여정이었어요.',
    note: '과정 완료 후 completion UI에 붙일 규칙입니다.',
  },
];

let lumiRuleDraftOverrides: LumiRuleDraft[] | null = null;

function cloneRuleDrafts(rules: LumiRuleDraft[]): LumiRuleDraft[] {
  return rules.map((rule) => ({ ...rule }));
}

function currentLumiRuleDrafts(): LumiRuleDraft[] {
  return lumiRuleDraftOverrides ?? lumiRuleDraftPresets;
}

interface ResolveLumiRuleDraftOptions {
  runtimeContext: LumiRuntimeContext;
  trigger: LumiTrigger;
  payload?: LumiTriggerPayload;
}

function scoreLumiRuleDraft(rule: LumiRuleDraft, options: ResolveLumiRuleDraftOptions): number {
  const { runtimeContext, trigger, payload } = options;

  if (rule.trigger !== trigger) return -1;
  if (rule.page !== runtimeContext.page) return -1;
  if (rule.mode !== runtimeContext.mode) return -1;
  if (rule.action !== runtimeContext.action) return -1;
  if (runtimeContext.scene && rule.scene !== runtimeContext.scene) return -1;

  let score = 10;

  if (payload?.context) {
    if (rule.context !== payload.context) return -1;
    score += 1;
  }

  if (payload?.dockSlot) {
    if (rule.dockSlot !== payload.dockSlot) return -1;
    score += 1;
  }

  return score;
}

export function getLumiRuleDrafts(): LumiRuleDraft[] {
  return cloneRuleDrafts(currentLumiRuleDrafts());
}

export function setLumiRuleDraftOverrides(next: LumiRuleDraft[] | null) {
  lumiRuleDraftOverrides = next ? cloneRuleDrafts(next) : null;
}

export function resolveLumiRuleDraft(options: ResolveLumiRuleDraftOptions): LumiRuleDraft | null {
  let bestMatch: LumiRuleDraft | null = null;
  let bestScore = -1;

  for (const rule of currentLumiRuleDrafts()) {
    const score = scoreLumiRuleDraft(rule, options);
    if (score > bestScore) {
      bestScore = score;
      bestMatch = rule;
    }
  }

  return bestMatch;
}

export function createInitialLumiRuleDrafts(): LumiRuleDraft[] {
  return cloneRuleDrafts(lumiRuleDraftPresets);
}

export function getDefaultLumiRuleDraftId(): string {
  return currentLumiRuleDrafts()[0]?.id ?? '';
}
