import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type DashboardSettingsProfileCopy = {
  loading: string
  shell: {
    eyebrow: string
    title: string
    copy: string
    tabs: {
      profile: string
      points: string
      ai: string
    }
  }
  profile: {
    title: string
    subtitle: string
    standardPlan: string
    avatarTitle: string
    avatarAlt: string
    avatarDescription: string
    avatarUploading: string
    avatarSelect: string
    fallbackInitial: string
    labels: {
      email: string
      provider: string
      createdAt: string
      points: string
      editableNickname: string
    }
    unset: string
    defaultProvider: string
    pointTotal: (total: number) => string
    pointBreakdown: (free: number, paid: number) => string
    editHint: string
    saving: string
    saveNickname: string
  }
  language: {
    title: string
    subtitle: string
    optionTitle: string
    optionHelp: string
    saving: string
    save: string
  }
  withdrawal: {
    title: string
    subtitle: string
    expand: string
    collapse: string
    warning: string
    confirmLabel: string
    confirmPlaceholder: string
    confirmValue: string
    processing: string
    action: string
  }
}

const ko: DashboardSettingsProfileCopy = {
  loading: '학습자 정보를 불러오는 중...',
  shell: {
    eyebrow: '학습자 설정',
    title: '학습자 정보 관리',
    copy: '닉네임, Login Method, 요금제, 가입일과 같은 기본 계정 정보를 확인하고 수정합니다.',
    tabs: { profile: '학습자정보', points: '포인트', ai: 'AI 설정' },
  },
  profile: {
    title: '기본 정보',
    subtitle: '계정 식별 정보와 학습자 표시 이름을 관리합니다.',
    standardPlan: 'Standard',
    avatarTitle: '프로필 이미지',
    avatarAlt: '프로필 이미지',
    avatarDescription: 'Social Login 프로필 사진이 있으면 자동으로 표시됩니다. LearnCosmos에서 사용할 전용 이미지는 JPG/PNG/WebP, 1MB 이하로 바꿀 수 있습니다.',
    avatarUploading: '업로드 중...',
    avatarSelect: '이미지 선택',
    fallbackInitial: '학',
    labels: {
      email: '이메일',
      provider: 'Login Method',
      createdAt: '가입일',
      points: '보유 포인트',
      editableNickname: '편집 가능한 닉네임',
    },
    unset: '미설정',
    defaultProvider: '기본',
    pointTotal: (total) => `총 ${total}pt`,
    pointBreakdown: (free, paid) => `무료 ${free}pt · 유료 ${paid}pt`,
    editHint: '상단 공통 헤더와 학습 화면에 표시되는 이름입니다.',
    saving: '저장 중...',
    saveNickname: '닉네임 저장',
  },
  language: {
    title: 'Language Settings',
    subtitle: '언어를 하나만 선택하면 화면 문구와 AI 생성 결과에 함께 적용됩니다.',
    optionTitle: 'Language',
    optionHelp: '메뉴, 안내문, 동의 화면, 목표 채팅, 리슨 생성, 추천 문구에 함께 적용합니다.',
    saving: '저장 중...',
    save: 'Save Language',
  },
  withdrawal: {
    title: '계정 탈퇴',
    subtitle: '계정 연결을 끊고 LearnCosmos 이용을 중단합니다. 탈퇴 후에는 현재 계정으로 Login할 수 없습니다.',
    expand: '펼치기',
    collapse: '접기',
    warning: '탈퇴하면 계정은 비활성화되고 등록한 BYOK/API 키는 삭제됩니다. 학습 기록은 개인정보처리방침의 보관 기준에 따라 계정과 분리된 비활성 상태로 보관되며, 같은 소셜 계정으로 다시 가입하면 보관 기간 내 이전 학습 자료를 다시 확인할 수 있습니다.',
    confirmLabel: '탈퇴 확인 문구',
    confirmPlaceholder: '탈퇴합니다',
    confirmValue: '탈퇴합니다',
    processing: '탈퇴 처리 중...',
    action: '계정 탈퇴',
  },
}

const en: DashboardSettingsProfileCopy = {
  loading: 'Loading learner profile...',
  shell: {
    eyebrow: 'Learner Settings',
    title: 'Profile',
    copy: 'Review and update your nickname, login method, plan, join date, and basic account details.',
    tabs: { profile: 'Profile', points: 'Points', ai: 'AI Settings' },
  },
  profile: {
    title: 'Basic Info',
    subtitle: 'Manage account identifiers and your learner display name.',
    standardPlan: 'Standard',
    avatarTitle: 'Profile Image',
    avatarAlt: 'Profile image',
    avatarDescription: 'Your social profile image is shown automatically when available. You can upload a JPG, PNG, or WebP image up to 1 MB for LearnCosmos.',
    avatarUploading: 'Uploading...',
    avatarSelect: 'Choose Image',
    fallbackInitial: 'L',
    labels: {
      email: 'Email',
      provider: 'Login Method',
      createdAt: 'Joined',
      points: 'Points',
      editableNickname: 'Nickname',
    },
    unset: 'Not set',
    defaultProvider: 'Default',
    pointTotal: (total) => `Total ${total}pt`,
    pointBreakdown: (free, paid) => `Free ${free}pt · Paid ${paid}pt`,
    editHint: 'This name appears in the shared header and learning screens.',
    saving: 'Saving...',
    saveNickname: 'Save Nickname',
  },
  language: {
    title: 'Language Settings',
    subtitle: 'Choose one language for both the interface and learning generation.',
    optionTitle: 'Language',
    optionHelp: 'Applies to menus, guidance, consent screens, goal chat, lesson generation, and recommendation text.',
    saving: 'Saving...',
    save: 'Save Language',
  },
  withdrawal: {
    title: 'Delete Account',
    subtitle: 'Disconnect your account and stop using LearnCosmos. You cannot sign in with this account after deletion.',
    expand: 'Expand',
    collapse: 'Collapse',
    warning: 'Deleting your account deactivates it and removes any saved BYOK/API keys. Learning records are kept in a deactivated state according to the Privacy Policy retention rules, and you may be able to view retained learning data again if you sign up with the same social account during the retention period.',
    confirmLabel: 'Confirmation Text',
    confirmPlaceholder: 'DELETE',
    confirmValue: 'DELETE',
    processing: 'Deleting...',
    action: 'Delete Account',
  },
}

export function getDashboardSettingsProfileCopy(locale: Locale | string | null | undefined): DashboardSettingsProfileCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
