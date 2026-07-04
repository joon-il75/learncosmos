import { normalizeLocale, type Locale } from '@/lib/i18n/locales';

export type PlatformShellSection = 'overview' | 'news' | 'guide' | 'operations' | 'policies';

export interface PlatformLinkItem {
  label: string;
  body: string;
  href: string;
}

export interface PlatformNewsBoardCopy {
  badge: string;
  title: string;
  description: string;
  emptyTitle: string;
  emptyDescription: string;
  adminOnlyLabel: string;
}

export interface PlatformShellCopy {
  shell: {
    eyebrow: string;
    title: string;
    subtitle: string;
    homeLabel: string;
    appLabel: string;
    statusLabel: string;
  };
  nav: Record<PlatformShellSection, string>;
  sections: Record<PlatformShellSection, {
    eyebrow: string;
    title: string;
    description: string;
    primaryAction: string;
    secondaryAction: string;
    notes: string[];
  }>;
  links: {
    title: string;
    subtitle: string;
    items: PlatformLinkItem[];
  };
  newsBoard: PlatformNewsBoardCopy;
  standards: {
    title: string;
    items: Array<{
      label: string;
      body: string;
    }>;
  };
}

const ko: PlatformShellCopy = {
  shell: {
    eyebrow: 'Platform Menu',
    title: '서비스 안내와 운영 기준을 한 곳에 모읍니다',
    subtitle: 'Platform menu는 LearnCosmos의 공지, 사용 가이드, 운영 기준, 정책 문서로 이어지는 공개 안내 shell입니다.',
    homeLabel: '홈',
    appLabel: '내 학습지도',
    statusLabel: '공개 안내 shell',
  },
  nav: {
    overview: 'Platform 홈',
    news: '새소식(공지사항)',
    guide: '가이드',
    operations: '운영 기준',
    policies: '정책',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: '학습자가 서비스의 기준을 빠르게 찾는 입구',
      description: '흩어진 약관, 개인정보, 오픈소스, YouTube API 안내와 운영 기준을 한 화면 흐름으로 묶습니다. 이번 단계에서는 새 CMS나 게시판을 만들지 않습니다.',
      primaryAction: '정책 모아보기',
      secondaryAction: '공지 게시판은 준비 중',
      notes: ['기존 공개 문서를 재사용합니다.', '운영 기준은 안내 shell로 먼저 보여줍니다.', '관리자 편집 workflow는 후속 단계로 보류합니다.'],
    },
    news: {
      eyebrow: 'News',
      title: '새소식과 공지사항을 모아두는 게시판',
      description: '서비스 공지, 업데이트 안내, 알파 테스트 소식을 읽기 전용 게시판으로 모아둘 공간입니다. 글쓰기와 수정은 슈퍼관리자 모드에서만 제공합니다.',
      primaryAction: '알파 테스터 신청하기',
      secondaryAction: '글쓰기는 슈퍼관리자 전용',
      notes: ['공개 페이지는 읽기 전용입니다.', '공지 작성/수정/삭제는 슈퍼관리자 모드에서만 제공합니다.', '현재는 실제 공지 DB 없이 빈 게시판 상태로 제공합니다.'],
    },
    guide: {
      eyebrow: 'Guide',
      title: '처음 온 학습자가 길을 잃지 않게 돕는 안내',
      description: '목표 입력, 학습계획, 지점 학습, 기록과 결과물 흐름을 설명하는 공개 가이드 영역입니다.',
      primaryAction: '학습 시작하기',
      secondaryAction: '가이드 문서 준비 중',
      notes: ['현재는 서비스 흐름 안내 shell만 제공합니다.', '상세 튜토리얼과 영상은 후속으로 연결합니다.', '홈 영상은 아직 placeholder 기준을 유지합니다.'],
    },
    operations: {
      eyebrow: 'Operations',
      title: '운영 원칙을 투명하게 정리하는 자리',
      description: '안전, 콘텐츠 검토, 커뮤니티 공개, Creator 권한 같은 운영 기준을 모아둘 공간입니다.',
      primaryAction: '커뮤니티 안내 보기',
      secondaryAction: '운영 문서 준비 중',
      notes: ['커뮤니티 공개/신고/검토 기준은 후속으로 구체화합니다.', 'Creator 권한과 콘텐츠 검토 기준도 별도 단계에서 닫습니다.', '운영 데이터 모델은 이번 단계에서 만들지 않습니다.'],
    },
    policies: {
      eyebrow: 'Policies',
      title: '약관과 정책 문서로 빠르게 이동합니다',
      description: '서비스 이용약관, 개인정보 처리방침, 오픈소스 고지, YouTube API 사용 고지, 동의 화면을 연결합니다.',
      primaryAction: '이용약관 보기',
      secondaryAction: '정책 CMS 준비 중',
      notes: ['기존 정책 페이지를 그대로 연결합니다.', '정책 문서 canonical set은 기존 policy_documents 기준을 유지합니다.', '정책 편집 CMS는 만들지 않습니다.'],
    },
  },
  links: {
    title: '바로가기',
    subtitle: '현재 운영 중인 공개 문서와 화면입니다.',
    items: [
      { label: '서비스 이용약관', body: 'LearnCosmos 이용 기준을 확인합니다.', href: '/terms' },
      { label: '개인정보 처리방침', body: '개인정보 수집과 처리 기준을 확인합니다.', href: '/privacy' },
      { label: '오픈소스 고지', body: '사용한 오픈소스와 라이선스를 확인합니다.', href: '/open-source' },
      { label: 'YouTube API 사용', body: 'YouTube API 관련 고지를 확인합니다.', href: '/youtube-api-use' },
      { label: '동의 화면', body: '필수 정책 동의 흐름을 확인합니다.', href: '/agreements' },
      { label: '언어 설정', body: '학습 언어와 UI 언어 설정 흐름입니다.', href: '/language-setup' },
    ],
  },
  newsBoard: {
    badge: 'Read Only',
    title: '공지 리스트',
    description: '공지사항 목록은 이 영역에 연결됩니다. 지금은 게시판 자리만 먼저 열어둡니다.',
    emptyTitle: '아직 등록된 공지사항이 없습니다',
    emptyDescription: '서비스 공지와 업데이트 소식은 준비되는 대로 이곳에 게시됩니다.',
    adminOnlyLabel: '글쓰기와 수정은 슈퍼관리자 모드에서만 가능합니다',
  },
  standards: {
    title: '이번 단계의 경계',
    items: [
      { label: '재사용 우선', body: '이미 있는 정책/안내 route를 묶고 새 저장 구조는 만들지 않습니다.' },
      { label: '공개 shell', body: '로그인 없이 볼 수 있는 안내 구조로 유지합니다.' },
      { label: '후속 CMS 보류', body: '공지 작성, 정책 편집, 운영 문서 CMS는 별도 단계로 분리합니다.' },
    ],
  },
};

const en: PlatformShellCopy = {
  shell: {
    eyebrow: 'Platform Menu',
    title: 'Service guides and operation standards in one place',
    subtitle: 'Platform menu is a public shell that connects LearnCosmos news, guides, operation standards, and policy documents.',
    homeLabel: 'Home',
    appLabel: 'My Learning Map',
    statusLabel: 'Public guide shell',
  },
  nav: {
    overview: 'Platform Home',
    news: 'News & Notices',
    guide: 'Guide',
    operations: 'Operations',
    policies: 'Policies',
  },
  sections: {
    overview: {
      eyebrow: 'Overview',
      title: 'A quick entry point for service standards',
      description: 'This shell gathers terms, privacy, open source notices, YouTube API notices, and operation standards into one flow. This step does not add a CMS or bulletin board.',
      primaryAction: 'View policies',
      secondaryAction: 'News board coming later',
      notes: ['Existing public documents are reused.', 'Operation standards appear first as a guide shell.', 'Admin editing workflows stay in a later phase.'],
    },
    news: {
      eyebrow: 'News',
      title: 'A board for news and notices',
      description: 'This read-only board will gather service notices, release updates, and alpha testing news. Authoring and editing stay in super-admin mode only.',
      primaryAction: 'Apply for alpha test',
      secondaryAction: 'Writing is super-admin only',
      notes: ['The public page is read-only.', 'Writing, editing, and deleting notices are reserved for super-admin mode.', 'For now, the board is intentionally empty until notice data is connected.'],
    },
    guide: {
      eyebrow: 'Guide',
      title: 'Help new learners understand the path',
      description: 'This area explains goal input, learning plans, point learning, records, and artifacts as a public guide.',
      primaryAction: 'Start learning',
      secondaryAction: 'Detailed guides coming later',
      notes: ['This version only provides a service-flow guide shell.', 'Detailed tutorials and videos connect later.', 'The home intro video remains a placeholder.'],
    },
    operations: {
      eyebrow: 'Operations',
      title: 'Make operation principles visible',
      description: 'This area will collect safety, content review, community publishing, and Creator permission standards.',
      primaryAction: 'View community guide',
      secondaryAction: 'Operation docs coming later',
      notes: ['Community publishing, reports, and review standards need a later pass.', 'Creator permission and content review standards also stay separate.', 'No operation data model is created in this step.'],
    },
    policies: {
      eyebrow: 'Policies',
      title: 'Move quickly to terms and policy documents',
      description: 'This section links terms, privacy, open source notices, YouTube API notices, and the agreement flow.',
      primaryAction: 'View Terms',
      secondaryAction: 'Policy CMS coming later',
      notes: ['Existing policy pages are linked as-is.', 'The canonical policy set continues to use the existing policy_documents standard.', 'No policy editing CMS is created.'],
    },
  },
  links: {
    title: 'Quick links',
    subtitle: 'Current public documents and flows.',
    items: [
      { label: 'Terms of Service', body: 'Review the service usage standard.', href: '/en/terms' },
      { label: 'Privacy Policy', body: 'Review personal data handling standards.', href: '/en/privacy' },
      { label: 'Open Source Notice', body: 'Review open source libraries and licenses.', href: '/en/open-source' },
      { label: 'YouTube API Use', body: 'Review YouTube API related notices.', href: '/en/youtube-api-use' },
      { label: 'Agreement Flow', body: 'Review the required policy agreement flow.', href: '/en/agreements' },
      { label: 'Language Setup', body: 'Review UI and learning language setup.', href: '/en/language-setup' },
    ],
  },
  newsBoard: {
    badge: 'Read Only',
    title: 'Notice list',
    description: 'Notice posts will be connected here later. For now, the board space is ready and intentionally empty.',
    emptyTitle: 'No notices yet',
    emptyDescription: 'Service notices and release updates will appear here when ready.',
    adminOnlyLabel: 'Writing and editing are available in super-admin mode only',
  },
  standards: {
    title: 'Boundary of this step',
    items: [
      { label: 'Reuse first', body: 'Existing policy and guide routes are grouped without a new storage model.' },
      { label: 'Public shell', body: 'The structure stays visible without login.' },
      { label: 'CMS later', body: 'News authoring, policy editing, and operation-doc CMS stay separate.' },
    ],
  },
};

export function getPlatformShellCopy(locale: Locale | string | null | undefined): PlatformShellCopy {
  return normalizeLocale(locale) === 'en' ? en : ko;
}
