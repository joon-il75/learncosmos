import type { LumiPageContext, LumiQuickActionTemplate } from './lumiTypes';

export const lumiQuickActionTemplates: Record<LumiPageContext, LumiQuickActionTemplate[]> = {
  dashboard: [
    { id: 'new-course', label: '새 행성 만들기' },
    { id: 'resume-course', label: '진행 중인 탐험 보기' },
    { id: 'status-summary', label: '지금 상태 요약' },
  ],
  'planet-map': [
    { id: 'resume-course', label: '진행 중인 탐험 보기' },
    { id: 'status-summary', label: '지금 상태 요약' },
  ],
  'planet-scene': [
    { id: 'open-next-zone', label: '다음 구역 보기' },
    { id: 'scene-guide', label: '탐험 길잡이 보기' },
  ],
  diary: [
    { id: 'refine-goal', label: '목표 더 구체화하기' },
    { id: 'search-again', label: '자료 다시 찾기' },
    { id: 'add-practice', label: '실습 리슨 추가하기' },
    { id: 'review-structure', label: '구조 다시 보기' },
  ],
  player: [
    { id: 'summarize-lesson', label: '핵심만 요약해줘' },
    { id: 'next-step', label: '다음 단계 보기' },
    { id: 'note-help', label: '메모 도와줘' },
  ],
  completion: [
    { id: 'open-base', label: '탐험기지 보기' },
    { id: 'start-review', label: '복습 시작하기' },
    { id: 'open-next-planet', label: '다음 행성 열기' },
  ],
  mobile: [
    { id: 'open-lumi', label: '루미 열기' },
    { id: 'status-summary', label: '지금 상태 요약' },
  ],
};
