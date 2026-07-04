import { normalizeLocale, type Locale } from '@/lib/i18n/locales';

export type LearnerMainSection = 'home' | 'personal' | 'community' | 'creator' | 'platform';

export type LearnerSubSection =
  | 'homeCreate'
  | 'homeGuide'
  | 'homeOperations'
  | 'homeNews'
  | 'personalToday'
  | 'personalPlan'
  | 'personalMap'
  | 'personalDiary'
  | 'personalArtifacts'
  | 'communityHome'
  | 'communityDiscover'
  | 'communityMine'
  | 'communityCreate'
  | 'creatorIntro'
  | 'creatorContentNew'
  | 'creatorGroupCourses'
  | 'creatorContent'
  | 'creatorFeedback'
  | 'platformNews'
  | 'platformGuide'
  | 'platformOperations'
  | 'platformPolicies';

export interface LearnerMainNavItem {
  key: LearnerMainSection;
  label: string;
  shortLabel: string;
  iconLabel: string;
  href: string;
  publicHref: string;
}

export interface LearnerSubNavItem {
  key: LearnerSubSection;
  main: LearnerMainSection;
  label: string;
  href: string;
}

export interface LearnerNavigationModel {
  mainItems: LearnerMainNavItem[];
  subItemsByMain: Record<LearnerMainSection, LearnerSubNavItem[]>;
}

const mainSectionOrder: LearnerMainSection[] = ['home', 'personal', 'community', 'creator', 'platform'];

const mainHrefBySection: Record<LearnerMainSection, string> = {
  home: '/',
  personal: '/dashboard',
  community: '/dashboard/community',
  creator: '/dashboard/creator',
  platform: '/platform',
};

const publicHrefBySection: Record<LearnerMainSection, string> = {
  home: '/',
  personal: '/dashboard',
  community: '/dashboard/community',
  creator: '/dashboard/creator',
  platform: '/platform',
};

const mainLabels: Record<Locale, Record<LearnerMainSection, { label: string; shortLabel: string; iconLabel: string }>> = {
  ko: {
    home: { label: '홈', shortLabel: '홈', iconLabel: '홈' },
    personal: { label: '개인학습', shortLabel: '학습', iconLabel: '학습' },
    community: { label: '커뮤니티', shortLabel: '커뮤', iconLabel: '커뮤' },
    creator: { label: '크리에이터', shortLabel: '제작', iconLabel: '제작' },
    platform: { label: '플랫폼', shortLabel: '운영', iconLabel: '운영' },
  },
  en: {
    home: { label: 'Home', shortLabel: 'Home', iconLabel: 'Home' },
    personal: { label: 'Personal Learning', shortLabel: 'Learn', iconLabel: 'Learn' },
    community: { label: 'Community', shortLabel: 'Commu', iconLabel: 'Commu' },
    creator: { label: 'Creator', shortLabel: 'Create', iconLabel: 'Create' },
    platform: { label: 'Platform', shortLabel: 'Ops', iconLabel: 'Ops' },
  },
};

const subNavConfig: Record<LearnerMainSection, Array<{ key: LearnerSubSection; href: string }>> = {
  home: [
    { key: 'homeCreate', href: '/#course-create-ready' },
    { key: 'homeGuide', href: '/#intro-video-ready' },
    { key: 'homeOperations', href: '/#public-documents-ready' },
    { key: 'homeNews', href: '/#platform-news-ready' },
  ],
  personal: [
    { key: 'personalToday', href: '/dashboard' },
    { key: 'personalPlan', href: '/dashboard/goal' },
    { key: 'personalMap', href: '/dashboard' },
    { key: 'personalDiary', href: '/dashboard' },
    { key: 'personalArtifacts', href: '/dashboard' },
  ],
  community: [
    { key: 'communityHome', href: '/dashboard/community' },
    { key: 'communityDiscover', href: '/dashboard/community/discover' },
    { key: 'communityMine', href: '/dashboard/community/mine' },
    { key: 'communityCreate', href: '/dashboard/community?intent=create' },
  ],
  creator: [
    { key: 'creatorIntro', href: '/dashboard/creator' },
    { key: 'creatorContentNew', href: '/dashboard/creator/content' },
    { key: 'creatorGroupCourses', href: '/dashboard/creator/group-courses' },
    { key: 'creatorContent', href: '/dashboard/creator/content' },
    { key: 'creatorFeedback', href: '/dashboard/creator/feedback' },
  ],
  platform: [
    { key: 'platformNews', href: '/platform/news' },
    { key: 'platformGuide', href: '/platform/guide' },
    { key: 'platformOperations', href: '/platform/operations' },
    { key: 'platformPolicies', href: '/platform/policies' },
  ],
};

const subLabels: Record<Locale, Record<LearnerSubSection, string>> = {
  ko: {
    homeCreate: '코스 생성',
    homeGuide: '학습 안내',
    homeOperations: '운영 기준',
    homeNews: '새소식(공지사항)',
    personalToday: '오늘 할 일',
    personalPlan: '학습계획',
    personalMap: '내 학습지도',
    personalDiary: '탐험일지',
    personalArtifacts: '결과물',
    communityHome: '커뮤니티 홈',
    communityDiscover: '커뮤니티 찾기',
    communityMine: '내 커뮤니티',
    communityCreate: '만들기',
    creatorIntro: '크리에이터 소개',
    creatorContentNew: '콘텐츠 만들기',
    creatorGroupCourses: '그룹 코스',
    creatorContent: '내 콘텐츠',
    creatorFeedback: '피드백',
    platformNews: '새소식(공지사항)',
    platformGuide: '이용 안내',
    platformOperations: '운영 기준',
    platformPolicies: '정책 문서',
  },
  en: {
    homeCreate: 'Create Course',
    homeGuide: 'Learning Guide',
    homeOperations: 'Operations',
    homeNews: 'News',
    personalToday: 'Today Task',
    personalPlan: 'Learning Plan',
    personalMap: 'My Learning Map',
    personalDiary: 'Explorer Diary',
    personalArtifacts: 'Artifacts',
    communityHome: 'Community Home',
    communityDiscover: 'Discover',
    communityMine: 'My Communities',
    communityCreate: 'Create',
    creatorIntro: 'Creator Intro',
    creatorContentNew: 'Create Content',
    creatorGroupCourses: 'Group Courses',
    creatorContent: 'My Content',
    creatorFeedback: 'Feedback',
    platformNews: 'News & Notices',
    platformGuide: 'Guide',
    platformOperations: 'Operations',
    platformPolicies: 'Policies',
  },
};

export function getLearnerNavigationModel(locale?: Locale | string | null): LearnerNavigationModel {
  const resolvedLocale = normalizeLocale(locale);
  const mainItems = mainSectionOrder.map((key) => ({
    key,
    label: mainLabels[resolvedLocale][key].label,
    shortLabel: mainLabels[resolvedLocale][key].shortLabel,
    iconLabel: mainLabels[resolvedLocale][key].iconLabel,
    href: mainHrefBySection[key],
    publicHref: publicHrefBySection[key],
  }));

  const subItemsByMain = mainSectionOrder.reduce((acc, main) => {
    acc[main] = subNavConfig[main].map((item) => ({
      key: item.key,
      main,
      label: subLabels[resolvedLocale][item.key],
      href: item.href,
    }));
    return acc;
  }, {} as Record<LearnerMainSection, LearnerSubNavItem[]>);

  return { mainItems, subItemsByMain };
}

export function getLearnerMainSectionFromPath(pathname: string | null | undefined): LearnerMainSection {
  if (!pathname) return 'home';
  if (pathname.startsWith('/platform')) return 'platform';
  if (pathname.startsWith('/dashboard/community')) return 'community';
  if (pathname.startsWith('/dashboard/creator')) return 'creator';
  if (pathname.startsWith('/dashboard')) return 'personal';
  return 'home';
}

export function isLearnerMainNavActive(pathname: string | null | undefined, section: LearnerMainSection): boolean {
  return getLearnerMainSectionFromPath(pathname) === section;
}

export function isLearnerSubNavActive(pathname: string | null | undefined, item: LearnerSubNavItem): boolean {
  if (!pathname) return item.href === '/';
  if (item.href === '/') return pathname === '/' || pathname === '/en';
  return pathname === item.href || pathname.startsWith(`${item.href}/`);
}
