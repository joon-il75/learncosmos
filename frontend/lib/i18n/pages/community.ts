import { normalizeLocale, type Locale } from '@/lib/i18n/locales';

export type CommunityShellSection = 'overview' | 'discover' | 'mine';

export interface CommunityShellCopy {
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
  nav: Record<CommunityShellSection, string>;
  sections: Record<CommunityShellSection, {
    eyebrow: string;
    title: string;
    description: string;
    primaryAction: string;
    secondaryAction: string;
    notes: string[];
  }>;
  roadmap: {
    title: string;
    items: Array<{
      label: string;
      body: string;
    }>;
  };
}

const ko: CommunityShellCopy = {
  shell: {
    eyebrow: 'Community Shell',
    title: '함께 배우는 공간을 준비하고 있어요',
    subtitle: '지금은 다른 학습자의 여정, 공유 학습지도, 함께 이어갈 그룹 학습을 예고하는 안내 화면입니다.',
    headerEyebrow: '개인학습',
    headerTitle: '커뮤니티',
    galaxyLabel: '내 학습지도',
    galaxyTitle: '내 학습지도 Galaxy로 이동',
    statusLabel: '알파 준비 중',
  },
  nav: {
    overview: '커뮤니티 홈',
    discover: '둘러보기',
    mine: '내 공유',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: '학습 여정이 쌓이면 함께 보는 지도가 됩니다',
      description: '커뮤니티는 완성된 결과물보다 과정, 기록, 다시 시도한 흔적을 중심으로 설계합니다. 지금은 shell만 열어 두고 실제 공개와 승인 흐름은 후속 단계에서 다룹니다.',
      primaryAction: '내 학습지도 보기',
      secondaryAction: '둘러보기 준비 상태',
      notes: ['내가 완료한 행성의 기록을 나중에 선택 공개합니다.', '다른 학습자의 지점 기록은 학습 참고 흐름으로 보여줄 예정입니다.', '승인, 신고, 검색 모델은 이번 단계에서 만들지 않습니다.'],
    },
    discover: {
      eyebrow: 'Discover',
      title: '다른 학습자의 탐험을 발견하는 입구',
      description: '공개된 학습지도, 결과물, 지점별 기록을 찾는 화면이 될 예정입니다. 현재는 검색 데이터 모델 없이 안내 상태만 제공합니다.',
      primaryAction: '내 학습 계속하기',
      secondaryAction: '공개 학습지도 준비 중',
      notes: ['추천/검색은 운영 기준이 정리된 뒤 연결합니다.', '첫 버전은 취미 주제와 학습 단계별 탐색을 우선합니다.', '외부 노출 전 안전 검토와 신고 흐름이 필요합니다.'],
    },
    mine: {
      eyebrow: 'My Shares',
      title: '내가 공유할 학습 흔적을 고르는 자리',
      description: '내 행성의 기록, 결과물, 회고를 골라 공개하는 관리 화면으로 확장할 예정입니다. 지금은 비공개 기본값을 유지합니다.',
      primaryAction: '완료한 행성 확인',
      secondaryAction: '공유 설정 준비 중',
      notes: ['사용자 동의 전에는 어떤 학습 기록도 공개하지 않습니다.', '공개 범위는 행성, 지점, 결과물 단위로 나누는 방향입니다.', '공유 취소와 비공개 전환은 필수 기준으로 둡니다.'],
    },
  },
  roadmap: {
    title: '후속 연결 기준',
    items: [
      { label: '공개 전 승인', body: '학습 기록 공개 전 사용자 선택과 운영 검토 기준을 둡니다.' },
      { label: '검색과 추천', body: '주제, 단계, 결과물 유형으로 찾을 수 있게 확장합니다.' },
      { label: '그룹 학습', body: '같은 목표를 가진 학습자가 같은 지도를 함께 걷는 흐름으로 연결합니다.' },
    ],
  },
};

const en: CommunityShellCopy = {
  shell: {
    eyebrow: 'Community Shell',
    title: 'A shared learning space is on the way',
    subtitle: 'This shell previews learner journeys, shared learning maps, and future group learning without enabling public sharing yet.',
    headerEyebrow: 'Personal Learning',
    headerTitle: 'Community',
    galaxyLabel: 'My Learning Map',
    galaxyTitle: 'Go to My Learning Map Galaxy',
    statusLabel: 'Alpha preparation',
  },
  nav: {
    overview: 'Community Home',
    discover: 'Discover',
    mine: 'My Shares',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: 'Learning journeys become maps we can learn from together',
      description: 'Community is designed around process, records, retries, and artifacts rather than polished outcomes alone. For now, this is a shell; publishing and review flows stay in a later phase.',
      primaryAction: 'View My Learning Map',
      secondaryAction: 'Discover is preparing',
      notes: ['Completed planet records may later be shared by choice.', 'Other learners point records will be framed as learning references.', 'Review, reporting, and search models are not part of this step.'],
    },
    discover: {
      eyebrow: 'Discover',
      title: 'The entry point for finding learner explorations',
      description: 'This will become a place to browse shared maps, artifacts, and point records. The current version keeps it as an empty state without a search data model.',
      primaryAction: 'Continue my learning',
      secondaryAction: 'Public maps coming later',
      notes: ['Search and recommendations will wait for operation rules.', 'The first version should prioritize hobby topics and learning stages.', 'Safety review and reporting flows are required before public exposure.'],
    },
    mine: {
      eyebrow: 'My Shares',
      title: 'Choose which learning traces you may share later',
      description: 'This area will grow into controls for publishing selected records, artifacts, and reflections. The default remains private.',
      primaryAction: 'Check completed planets',
      secondaryAction: 'Sharing settings coming later',
      notes: ['No learning record is public without learner consent.', 'Sharing should be scoped by planet, point, and artifact.', 'Unsharing and private fallback are required standards.'],
    },
  },
  roadmap: {
    title: 'Future connection rules',
    items: [
      { label: 'Review before public sharing', body: 'Learner choice and operation review standards come first.' },
      { label: 'Search and recommendation', body: 'Learners should later browse by topic, stage, and artifact type.' },
      { label: 'Group learning', body: 'Learners with similar goals can eventually walk the same map together.' },
    ],
  },
};

export function getCommunityShellCopy(locale: Locale | string | null | undefined): CommunityShellCopy {
  return normalizeLocale(locale) === 'en' ? en : ko;
}
