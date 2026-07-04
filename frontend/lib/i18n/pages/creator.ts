import { normalizeLocale, type Locale } from '@/lib/i18n/locales';

export type CreatorShellSection = 'overview' | 'content' | 'groupCourses' | 'feedback';

export interface CreatorShellCopy {
  shell: {
    eyebrow: string;
    title: string;
    subtitle: string;
    headerEyebrow: string;
    headerTitle: string;
    galaxyLabel: string;
    galaxyTitle: string;
    statusLabel: string;
  };
  nav: Record<CreatorShellSection, string>;
  sections: Record<CreatorShellSection, {
    eyebrow: string;
    title: string;
    description: string;
    primaryAction: string;
    secondaryAction: string;
    notes: string[];
  }>;
  readiness: {
    title: string;
    subtitle: string;
    items: Array<{
      label: string;
      body: string;
    }>;
  };
}

const ko: CreatorShellCopy = {
  shell: {
    eyebrow: 'Creator Shell',
    title: '가르치는 사람이 학습 여정을 만들 수 있게 준비합니다',
    subtitle: 'Creator 영역은 취미 전문가, 튜터, 운영자가 학습 콘텐츠와 그룹 코스를 구성하는 입구입니다. 지금은 기능을 열기 전 안내 shell로 둡니다.',
    headerEyebrow: '개인학습',
    headerTitle: '크리에이터',
    galaxyLabel: '내 학습지도',
    galaxyTitle: '내 학습지도 Galaxy로 이동',
    statusLabel: 'Creator alpha 준비 중',
  },
  nav: {
    overview: 'Creator 홈',
    content: '콘텐츠',
    groupCourses: '그룹 코스',
    feedback: '피드백',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: '학습자를 위한 작은 우주를 설계하는 자리',
      description: 'Creator shell은 콘텐츠 등록, 그룹 코스, 학습자 피드백을 한 흐름으로 묶기 위한 준비 화면입니다. 실제 운영 기능은 기준이 정리된 뒤 단계적으로 연결합니다.',
      primaryAction: '내 학습지도 보기',
      secondaryAction: 'Creator 신청 준비 중',
      notes: ['처음에는 검증된 취미 주제와 짧은 학습 흐름을 우선합니다.', 'Creator 권한, 검토, 노출 기준은 운영 문서와 함께 확정합니다.', '학습자 기록을 직접 수정하는 기능은 만들지 않습니다.'],
    },
    content: {
      eyebrow: 'Content Studio',
      title: '콘텐츠는 지점 학습에 연결되는 작은 재료가 됩니다',
      description: '영상, 문서, 실습 과제, 참고자료를 학습 지점에 맞게 등록하는 영역이 될 예정입니다. 지금은 등록 API 없이 구조와 방향만 안내합니다.',
      primaryAction: '학습 흐름 확인하기',
      secondaryAction: '콘텐츠 등록 준비 중',
      notes: ['콘텐츠는 course/lesson/point 구조와 연결되는 방향을 유지합니다.', '저작권, YouTube API, 외부 자료 출처 표시 기준이 필요합니다.', '검토 전 콘텐츠는 공개하지 않는 기준을 둡니다.'],
    },
    groupCourses: {
      eyebrow: 'Group Courses',
      title: '같은 목표를 가진 학습자가 같은 지도를 함께 걷는 방식',
      description: '그룹 코스는 Creator가 만든 흐름을 여러 학습자가 함께 따라가도록 확장하는 영역입니다. enrollment와 결제/정산은 이번 단계에서 만들지 않습니다.',
      primaryAction: '커뮤니티 shell 보기',
      secondaryAction: '그룹 코스 준비 중',
      notes: ['개인학습의 course/lesson/point 구조를 재사용하는 방향입니다.', '참여자 진행률, 공개 범위, 운영 권한 설계가 필요합니다.', '수익화와 정산은 별도 Phase로 분리합니다.'],
    },
    feedback: {
      eyebrow: 'Learner Feedback',
      title: '학습자의 막힌 부분을 Creator 개선 신호로 모읍니다',
      description: '학습자가 남긴 질문, 어려웠던 지점, 결과물 반응을 콘텐츠 개선으로 연결하는 영역입니다. 실제 피드백 API는 후속 단계에서 다룹니다.',
      primaryAction: '오늘 할 일로 돌아가기',
      secondaryAction: '피드백 수집 준비 중',
      notes: ['피드백은 개인 학습 기록 보호를 우선합니다.', 'Creator에게 노출되는 정보는 범위와 익명화 기준이 필요합니다.', '신고/안전 검토와 함께 운영해야 합니다.'],
    },
  },
  readiness: {
    title: '열기 전에 필요한 기준',
    subtitle: 'Creator 기능은 학습자 경험과 운영 책임이 함께 움직입니다. 그래서 shell 이후 단계에서 아래 기준을 먼저 닫습니다.',
    items: [
      { label: '권한과 검토', body: '누가 Creator가 될 수 있고 어떤 콘텐츠가 공개되는지 정합니다.' },
      { label: '자료 출처와 안전', body: '외부 자료, 저작권, 신고, 비공개 처리 기준을 둡니다.' },
      { label: '그룹 운영', body: '참여자 관리, 진행률 공개, 수익화 여부를 별도로 설계합니다.' },
    ],
  },
};

const en: CreatorShellCopy = {
  shell: {
    eyebrow: 'Creator Shell',
    title: 'Preparing a space where teachers can design learning journeys',
    subtitle: 'Creator is the future entry point for hobby experts, tutors, and operators to shape content and group courses. For now, it stays as a guided shell before live tools open.',
    headerEyebrow: 'Personal Learning',
    headerTitle: 'Creator',
    galaxyLabel: 'My Learning Map',
    galaxyTitle: 'Go to My Learning Map Galaxy',
    statusLabel: 'Creator alpha preparing',
  },
  nav: {
    overview: 'Creator Home',
    content: 'Content',
    groupCourses: 'Group Courses',
    feedback: 'Feedback',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: 'Design small learning worlds for learners',
      description: 'Creator shell prepares one flow for content, group courses, and learner feedback. Real operation tools will be connected step by step after the standards are settled.',
      primaryAction: 'View My Learning Map',
      secondaryAction: 'Creator applications coming later',
      notes: ['The first scope should focus on verified hobby topics and short learning flows.', 'Creator permission, review, and exposure rules need operation standards.', 'Creator tools must not directly edit learner records.'],
    },
    content: {
      eyebrow: 'Content Studio',
      title: 'Content becomes small material connected to point learning',
      description: 'This area will later register videos, documents, practice tasks, and references by learning point. The current shell has no registration API.',
      primaryAction: 'Review learning flow',
      secondaryAction: 'Content registration coming later',
      notes: ['Content should stay connected to course, lesson, and point structures.', 'Copyright, YouTube API, and external source display rules are required.', 'Unreviewed content should not be public.'],
    },
    groupCourses: {
      eyebrow: 'Group Courses',
      title: 'Let learners with the same goal walk the same map together',
      description: 'Group courses will let multiple learners follow a Creator-made flow. Enrollment, payments, and settlement are not part of this step.',
      primaryAction: 'View Community shell',
      secondaryAction: 'Group courses coming later',
      notes: ['The direction is to reuse the personal course, lesson, and point structure.', 'Participant progress, visibility, and operation authority need design.', 'Monetization and settlement stay in a separate phase.'],
    },
    feedback: {
      eyebrow: 'Learner Feedback',
      title: 'Collect learner blockers as signals for Creator improvement',
      description: 'This area will later connect questions, difficult points, and artifact responses to content improvement. The feedback API stays for a later step.',
      primaryAction: 'Return to Today Task',
      secondaryAction: 'Feedback collection coming later',
      notes: ['Feedback must protect personal learning records first.', 'Creator-visible information needs scope and anonymization standards.', 'Safety review and reporting should operate together.'],
    },
  },
  readiness: {
    title: 'Standards needed before launch',
    subtitle: 'Creator tools carry both learner experience and operation responsibility. After this shell, these standards should close first.',
    items: [
      { label: 'Permission and review', body: 'Define who can become a Creator and what can be published.' },
      { label: 'Sources and safety', body: 'Set rules for external material, copyright, reports, and takedown.' },
      { label: 'Group operation', body: 'Design participant management, progress visibility, and monetization separately.' },
    ],
  },
};

export function getCreatorShellCopy(locale: Locale | string | null | undefined): CreatorShellCopy {
  return normalizeLocale(locale) === 'en' ? en : ko;
}
