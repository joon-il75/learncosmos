import type { PointWorkTab } from '../../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';

export type NoteFieldKey = 'observation' | 'reflection' | 'examples' | 'nextStep';
export type LearningWorkGuideDetail = { title: string; message: string };
export type ArtifactKind = '영상자료' | '문서 편집자료' | '첨부자료';

export const artifactTypeOptions: readonly ArtifactKind[] = ['영상자료', '문서 편집자료', '첨부자료'] as const;

export const pointWorkTabs: Array<{ key: PointWorkTab; icon: string; label: string; caption: string }> = [
  { key: 'notes', icon: '📝', label: '1. 내용정리', caption: '핵심 개념과 내 설명' },
  { key: 'questions', icon: '❓', label: '2. 질문관리', caption: '모르는 점과 답변' },
  { key: 'practice', icon: '🔁', label: '3. 연습/활동기록', caption: '시도와 막힘' },
  { key: 'artifacts', icon: '📦', label: '4. 결과물제출', caption: '만든 것과 배운 점' },
  { key: 'attachments', icon: '📎', label: '5. 보조자료', caption: '참고 링크와 파일' },
  { key: 'timeline', icon: '🕘', label: '6. 활동 기록', caption: '자동 활동 타임라인' },
];

export function normalizeArtifactKind(value: string): ArtifactKind {
  if (value === '영상자료' || value === '동영상') return '영상자료';
  if (value === '첨부자료' || value === '파일' || value === '이미지' || value === '링크') return '첨부자료';
  return '문서 편집자료';
}

export function getLearningWorkGuideDetail(
  activeNoteField: NoteFieldKey,
  activeWorkTab: PointWorkTab,
  copy?: PointLearningCopy['workspace']['learningWork'],
): LearningWorkGuideDetail {
  if (copy) {
    if (activeWorkTab === 'questions') return copy.tabs.questions;
    if (activeWorkTab === 'practice') return copy.tabs.practice;
    if (activeWorkTab === 'artifacts') return copy.tabs.artifacts;
    if (activeWorkTab === 'attachments') return copy.tabs.attachments;
    if (activeWorkTab === 'timeline') return copy.tabs.timeline;
    return copy.noteGuides[activeNoteField];
  }
  if (activeWorkTab === 'questions') {
    return {
      title: 'Questions - 질문관리',
      message: '여기서는 이해가 막힌 부분이나 더 확인하고 싶은 점을 질문으로 남기면 됩니다. 질문을 저장한 뒤 나의 답변을 적고, 필요하면 AI 참고 피드백도 받아보세요.',
    };
  }
  if (activeWorkTab === 'practice') {
    return {
      title: 'Practice & Activity Log - 연습/활동기록',
      message: '여기서는 실제로 해 본 연습, 조사, 분석, 적용 시도를 남기면 됩니다. 걸린 시간과 막힌 부분, 다음 계획을 짧게 적어도 충분합니다.',
    };
  }
  if (activeWorkTab === 'artifacts') {
    return {
      title: 'Artifacts - 결과물',
      message: '여기서는 만든 산출물이나 제출할 링크를 남기면 됩니다. 무엇을 만들었는지, 만들면서 배운 점과 어려웠던 점을 함께 정리해 보세요.',
    };
  }
  if (activeWorkTab === 'attachments') {
    return {
      title: 'References - 보조자료',
      message: '여기서는 학습 중 추가로 참고한 링크나 파일을 모아 두면 됩니다. 내용정리와 섞지 않고 자료 출처를 따로 챙길 때 사용하면 좋습니다.',
    };
  }
  if (activeWorkTab === 'timeline') {
    return {
      title: 'Activity Log - 활동 기록',
      message: '여기서는 이 지점에서 저장한 주요 행동을 시간순으로 볼 수 있습니다. 내가 어떤 순서로 학습했는지 천천히 되짚어보세요.',
    };
  }
  if (activeNoteField === 'reflection') {
    return {
      title: 'My Explanation - 내 설명',
      message: '여기서는 학습한 내용을 나의 말로 다시 풀어 쓰면 됩니다. 그대로 베끼기보다 누군가에게 설명하듯 적으면 이해가 더 분명해집니다.',
    };
  }
  if (activeNoteField === 'examples') {
    return {
      title: 'Examples - 예시',
      message: '여기서는 개념이 실제로 쓰이는 상황이나 내가 떠올린 사례를 남겨보세요. 적용 장면을 적어 두면 나중에 다시 이해하기 쉽습니다.',
    };
  }
  if (activeNoteField === 'nextStep') {
    return {
      title: 'Confusing Parts - 헷갈린 부분',
      message: '여기서는 아직 확실하지 않은 내용과 다음에 확인할 질문을 남기면 됩니다. 완벽히 정리되지 않은 상태를 그대로 적어도 괜찮습니다.',
    };
  }
  return {
    title: 'Core Concept - 핵심 개념',
    message: '여기서는 지금 학습한 내용에서 가장 중요한 한두 문장을 남기면 됩니다. 길게 요약하기보다 나중에 다시 봐도 중심이 보이게 적어 보세요.',
  };
}

export function scrollToWorkspaceTarget(targetID: string): void {
  window.requestAnimationFrame(() => {
    window.requestAnimationFrame(() => {
      const target = document.getElementById(targetID);
      if (!target) return;
      const targetTop = target.getBoundingClientRect().top + window.scrollY - 430;
      window.scrollTo({ top: Math.max(0, targetTop), behavior: 'smooth' });
    });
  });
}

export const scrollToLearningRecordSection = () => scrollToWorkspaceTarget('point-learning-record-section');
