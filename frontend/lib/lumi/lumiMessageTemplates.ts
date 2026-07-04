import type { LumiMessageType, LumiTrigger, LumiTriggerPayload } from './lumiTypes';

export interface LumiTriggerMessageTemplate {
  messageType: LumiMessageType;
  message: string;
  preview: string;
}

export interface LumiMessageDraft extends LumiTriggerMessageTemplate {
  id: string;
  trigger: LumiTrigger;
  label: string;
  note: string;
}

export const lumiStatusMessages: Record<NonNullable<LumiTriggerPayload['status']>, { hover: string; select: string }> = {
  draft: {
    hover: '아직 탐험 계획을 차근차근 정리하는 중이에요.',
    select: '먼저 탐험 계획부터 함께 다듬어 보면 좋아요.',
  },
  ready: {
    hover: '탐험 준비가 잘 끝나서 이제 출발할 수 있어요.',
    select: '이제 편하게 첫 탐험을 시작해 보셔도 좋아요.',
  },
  learning: {
    hover: '지금 탐험을 차근차근 이어가고 있는 행성이에요.',
    select: '지금 흐름 그대로 이어서 탐험해 볼까요?',
  },
  completed: {
    hover: '탐험을 마친 행성이에요.',
    select: '여기서 쌓은 흐름을 다음 행성에도 자연스럽게 이어갈 수 있어요.',
  },
};

export const lumiTriggerMessageTemplates: Record<Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>, LumiTriggerMessageTemplate> = {
  cta_focus: {
    messageType: 'question',
    message: '배우고 싶은 주제를 편하게 적어 주세요.',
    preview: '',
  },
  cta_submit: {
    messageType: 'hint',
    message: '좋아요. 탐험 목표를 조금 더 또렷하게 함께 정리해 볼게요.',
    preview: '',
  },
  generation_loading: {
    messageType: 'guide',
    message: '첫 번째 탐험 경로를 차근차근 정리하고 있어요.\n행성도 함께 준비하고 있어요.',
    preview: '',
  },
  generation_success: {
    messageType: 'success',
    message: '새로운 행성이 준비됐어요.\n이제 첫 탐험 경로를 함께 다듬어 볼까요?',
    preview: '',
  },
  idle: {
    messageType: 'encourage',
    message: '천천히 다시 이어가도 괜찮아요.\n조금만 더 가면 다음 단계예요.',
    preview: '',
  },
  search_start: {
    messageType: 'discovery',
    message: '학습 자료를 탐험하고 있어요.\n이 단계에 잘 맞는 자료를 차근차근 찾고 있어요.',
    preview: '',
  },
  search_success: {
    messageType: 'discovery',
    message: '잘 맞는 자료를 찾았어요.\n부담 없이 따라가기 좋은 흐름으로 안내해 드릴게요.',
    preview: '',
  },
  evaluation_start: {
    messageType: 'guide',
    message: '전체 탐험 계획을 차분하게 점검하고 있어요.',
    preview: '',
  },
  evaluation_success: {
    messageType: 'success',
    message: '구조가 많이 좋아졌어요.\n몇 군데만 더 다듬으면 한결 또렷해질 거예요.',
    preview: '',
  },
  note_saved: {
    messageType: 'success',
    message: '좋은 기록이에요.\n나중에 복습할 때 큰 도움이 될 거예요.',
    preview: '',
  },
  lesson_complete: {
    messageType: 'success',
    message: '좋아요, 이 탐험목표를 잘 마쳤어요.\n이제 다음 경로로 편하게 이동하실 수 있어요.',
    preview: '',
  },
  level_complete: {
    messageType: 'success',
    message: '한 구간을 잘 마쳤어요.\n다음 탐험 구역도 차분하게 열어 볼 수 있어요.',
    preview: '',
  },
  course_complete: {
    messageType: 'success',
    message: '이 행성 탐험을 마쳤어요.\n정말 멋진 여정이었어요.',
    preview: '',
  },
  base_built: {
    messageType: 'success',
    message: '이 행성에 탐험기지가 세워졌어요.\n이제 이곳은 당신의 학습 거점이에요.',
    preview: '',
  },
  manual_open: {
    messageType: 'summary',
    message: '지금 필요한 안내를 함께 차근차근 정리해 볼게요.',
    preview: '',
  },
  page_enter: {
    messageType: 'summary',
    message: '지금 바로 이어서 탐험할 수 있는 행성이 있어요.\n원하시는 곳부터 편하게 살펴보세요.',
    preview: '',
  },
};

const lumiTriggerMessageLabels: Record<Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>, string> = {
  cta_focus: 'Hero Focus Prompt',
  cta_submit: 'Hero Goal Submit',
  generation_loading: 'Generation Loading',
  generation_success: 'Generation Success',
  idle: 'Idle Encourage',
  search_start: 'Search Start',
  search_success: 'Search Success',
  evaluation_start: 'Evaluation Start',
  evaluation_success: 'Evaluation Success',
  note_saved: 'Note Saved',
  lesson_complete: 'Lesson Complete',
  level_complete: 'Level Complete',
  course_complete: 'Course Complete',
  base_built: 'Base Built',
  manual_open: 'Manual Open',
  page_enter: 'Page Enter',
};

const lumiTriggerMessageNotes: Partial<Record<Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>, string>> = {
  cta_focus: '랜딩 Hero와 CTA focus 흐름의 질문 메시지입니다.',
  generation_loading: '생성 대기 상태에서 쓰는 대표 안내입니다.',
  generation_success: '초안 생성 완료 직후 또는 생성 성공 상태에 연결됩니다.',
  search_start: '자료 탐색 시작 상태에 연결됩니다.',
  evaluation_success: '평가 완료나 구조 피드백 성공 메시지입니다.',
  course_complete: 'completion 성격의 종료 메시지입니다.',
};

let lumiMessageDraftOverrides: LumiMessageDraft[] | null = null;

function cloneMessageDrafts(messages: LumiMessageDraft[]): LumiMessageDraft[] {
  return messages.map((message) => ({ ...message }));
}

function createDefaultLumiMessageDrafts(): LumiMessageDraft[] {
  return Object.entries(lumiTriggerMessageTemplates).map(([trigger, template]) => ({
    id: `message-${trigger}`,
    trigger: trigger as Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>,
    label: lumiTriggerMessageLabels[trigger as Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>],
    note: lumiTriggerMessageNotes[trigger as Exclude<LumiTrigger, 'planet_hover' | 'planet_select' | 'lesson_enter'>] ?? '',
    ...template,
  }));
}

function currentLumiMessageDrafts(): LumiMessageDraft[] {
  return lumiMessageDraftOverrides ?? createDefaultLumiMessageDrafts();
}

export function getLumiMessageDrafts(): LumiMessageDraft[] {
  return cloneMessageDrafts(currentLumiMessageDrafts());
}

export function setLumiMessageDraftOverrides(next: LumiMessageDraft[] | null) {
  lumiMessageDraftOverrides = next ? cloneMessageDrafts(next) : null;
}

export function resolvePlanetStatusMessage(mode: 'hover' | 'select', status: NonNullable<LumiTriggerPayload['status']>): string {
  return lumiStatusMessages[status][mode];
}

export function resolveLessonEnterMessage(lessonTitle?: string | null): LumiTriggerMessageTemplate {
  return {
    messageType: 'guide',
    message: `${lessonTitle ?? '이번 탐험'}에서는 핵심부터 차근차근 살펴볼게요.`,
    preview: '',
  };
}

export function resolveLumiMessageTemplate(
  trigger: LumiTrigger,
  payload?: LumiTriggerPayload,
): LumiTriggerMessageTemplate | null {
  if (trigger === 'planet_hover') {
    const title = payload?.courseTitle ?? '이 행성';
    const status = payload?.status ?? 'draft';
    return {
      messageType: 'guide',
      message: `${title}
${resolvePlanetStatusMessage('hover', status)}`,
      preview: '',
    };
  }

  if (trigger === 'planet_select') {
    const title = payload?.courseTitle ?? '이 행성';
    const status = payload?.status ?? 'draft';
    return {
      messageType: 'guide',
      message: `${title}
${resolvePlanetStatusMessage('select', status)}`,
      preview: '',
    };
  }

  if (trigger === 'lesson_enter') {
    return resolveLessonEnterMessage(payload?.lessonTitle);
  }

  const draft = currentLumiMessageDrafts().find((entry) => entry.trigger === trigger);
  if (!draft) return null;

  return {
    messageType: draft.messageType,
    message: draft.message,
    preview: draft.preview,
  };
}

export function resolveLumiMessageDraft(trigger: LumiTrigger): LumiMessageDraft | null {
  const draft = currentLumiMessageDrafts().find((entry) => entry.trigger === trigger);
  return draft ? { ...draft } : null;
}

export function createInitialLumiMessageDrafts(): LumiMessageDraft[] {
  return createDefaultLumiMessageDrafts();
}

export function getDefaultLumiMessageDraftId(): string {
  return currentLumiMessageDrafts()[0]?.id ?? '';
}
