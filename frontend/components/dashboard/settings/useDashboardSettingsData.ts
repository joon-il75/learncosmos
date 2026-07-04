'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import type { AppRouterInstance } from 'next/dist/shared/lib/app-router-context.shared-runtime';

export interface UserInfo {
  id: string;
  email: string;
  display_id: string;
  nickname: string;
  avatar_url?: string;
  role: string;
  premium_access: boolean;
  provider: string;
  created_at: string;
  free_points: number;
  paid_points: number;
  total_points: number;
  terms_agreed: boolean;
  privacy_agreed: boolean;
  required_consent_pending: boolean;
  ui_locale?: 'ko' | 'en';
  learning_language?: 'ko' | 'en';
  language_setup_required?: boolean;
}

export interface UserAISettings {
  mode: string;
  provider: string;
  has_api_key: boolean;
  is_enabled: boolean;
  endpoint_url?: string | null;
  updated_at?: string;
  last_validation_status?: string;
  last_validated_at?: string;
  last_validation_error?: string;
  last_failed_at?: string;
  next_retry_at?: string;
}

export interface UserAIUsageEvent {
  id: string;
  source: string;
  provider: string;
  model?: string;
  feature: string;
  billing_status?: string;
  input_tokens?: number;
  output_tokens?: number;
  estimated_cost_usd?: number;
  success: boolean;
  error_code?: string;
  created_at: string;
}

export interface UserAIUsageSummary {
  total_tokens: number;
  average_daily_tokens: number;
  total_estimated_cost_usd: number;
  average_daily_estimated_cost_usd: number;
}

interface UserAIUsageResponse {
  events?: UserAIUsageEvent[];
  total?: number;
  page?: number;
  limit?: number;
  total_pages?: number;
  start_date?: string;
  end_date?: string;
  summary?: Partial<UserAIUsageSummary>;
  error?: string;
  error_code?: string;
}

export interface UserPointTransaction {
  id: string;
  type: string;
  amount: number;
  feature?: string;
  reference_type?: string;
  reference_id?: string;
  description?: string;
  metadata?: Record<string, unknown>;
  created_at: string;
}

export interface UserPointUsageSummary {
  granted_points: number;
  purchased_points: number;
  used_points: number;
  refunded_points: number;
  net_change: number;
}

interface UserPointUsageResponse {
  transactions?: UserPointTransaction[];
  total?: number;
  page?: number;
  limit?: number;
  total_pages?: number;
  start_date?: string;
  end_date?: string;
  summary?: Partial<UserPointUsageSummary>;
  balance?: {
    free_points?: number;
    paid_points?: number;
    total_points?: number;
  };
  error?: string;
  error_code?: string;
}

type ApiErrorPayload = {
  error?: string;
  message?: string;
  error_code?: string;
}

export const AI_PROVIDERS = [
  { value: 'openai', label: 'OpenAI' },
  { value: 'anthropic', label: 'Anthropic' },
  { value: 'google', label: 'Google' },
  { value: 'grok', label: 'Grok' },
  { value: 'solar', label: 'Solar' },
  { value: 'hyperclova', label: 'HyperCLOVA' },
  { value: 'llama', label: 'Llama' },
  { value: 'exaone', label: 'EXAONE' },
] as const;

export function formatDateTime(iso?: string | null, locale: 'ko' | 'en' = 'ko') {
  if (!iso) return locale === 'en' ? 'Unknown' : '미확인';
  return new Date(iso).toLocaleString(locale === 'en' ? 'en-US' : 'ko-KR', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

const settingsErrorMessages: Record<string, { ko: string; en: string }> = {
  unauthorized: { ko: '로그인이 필요합니다.', en: 'Login is required.' },
  invalid_request: { ko: '잘못된 요청입니다.', en: 'Invalid request.' },
  profile_nickname_required: { ko: '닉네임을 입력해 주세요.', en: 'Enter a nickname.' },
  profile_nickname_too_long: { ko: '닉네임은 20자 이하여야 합니다.', en: 'Nickname must be 20 characters or fewer.' },
  profile_update_failed: { ko: '학습자 정보를 저장하지 못했습니다.', en: 'Could not save profile.' },
  user_not_found: { ko: '사용자를 찾을 수 없습니다.', en: 'Could not find the user.' },
  avatar_file_required: { ko: '프로필 이미지 파일이 필요합니다.', en: 'Choose a profile image file.' },
  avatar_file_too_large: { ko: '프로필 이미지는 1MB 이하여야 합니다.', en: 'Profile image must be 1MB or smaller.' },
  avatar_open_failed: { ko: '프로필 이미지를 열 수 없습니다.', en: 'Could not open the profile image.' },
  avatar_read_failed: { ko: '프로필 이미지를 읽지 못했습니다.', en: 'Could not read the profile image.' },
  avatar_type_unsupported: { ko: '프로필 이미지는 JPG, PNG, WebP만 사용할 수 있습니다.', en: 'Profile image must be JPG, PNG, or WebP.' },
  avatar_dimension_unreadable: { ko: '프로필 이미지 크기를 확인할 수 없습니다.', en: 'Could not read the profile image dimensions.' },
  avatar_dimension_too_large: { ko: '프로필 이미지는 가로/세로 1024px 이하여야 합니다.', en: 'Profile image width and height must be 1024px or smaller.' },
  avatar_storage_failed: { ko: '프로필 이미지 저장소를 준비하지 못했습니다.', en: 'Could not prepare profile image storage.' },
  avatar_filename_failed: { ko: '프로필 이미지 이름을 만들지 못했습니다.', en: 'Could not create a profile image filename.' },
  avatar_write_failed: { ko: '프로필 이미지를 저장하지 못했습니다.', en: 'Could not save the profile image.' },
  avatar_apply_failed: { ko: '프로필 이미지를 반영하지 못했습니다.', en: 'Could not apply the profile image.' },
  avatar_metadata_failed: { ko: '프로필 이미지 정보를 저장하지 못했습니다.', en: 'Could not save profile image metadata.' },
  withdrawal_failed: { ko: '계정 탈퇴 처리에 실패했습니다.', en: 'Could not delete account.' },
  preferences_save_failed: { ko: '언어 설정을 저장하지 못했습니다.', en: 'Could not save language settings.' },
  ai_settings_load_failed: { ko: 'AI 설정을 불러오지 못했습니다.', en: 'Could not load BYOK settings.' },
  ai_usage_invalid_start_date: { ko: '시작일을 확인해 주세요.', en: 'Check the start date.' },
  ai_usage_invalid_end_date: { ko: '종료일을 확인해 주세요.', en: 'Check the end date.' },
  ai_usage_invalid_date_range: { ko: '조회 기간을 확인해 주세요.', en: 'Check the date range.' },
  ai_usage_date_range_too_large: { ko: 'AI 사용량 조회 기간은 최대 3개월입니다.', en: 'BYOK usage range can be up to 3 months.' },
  ai_usage_load_failed: { ko: 'BYOK 사용량을 불러오지 못했습니다.', en: 'Could not load BYOK usage.' },
  point_usage_invalid_start_date: { ko: '시작일을 확인해 주세요.', en: 'Check the start date.' },
  point_usage_invalid_end_date: { ko: '종료일을 확인해 주세요.', en: 'Check the end date.' },
  point_usage_invalid_date_range: { ko: '조회 기간을 확인해 주세요.', en: 'Check the date range.' },
  point_usage_date_range_too_large: { ko: '포인트 내역 조회 기간은 최대 3개월입니다.', en: 'Point history range can be up to 3 months.' },
  point_usage_load_failed: { ko: '포인트 사용 내역을 불러오지 못했습니다.', en: 'Could not load point history.' },
  ai_provider_required: { ko: 'AI 제공자를 선택해 주세요.', en: 'Select an AI provider.' },
  ai_provider_unsupported: { ko: '지원하지 않는 AI 제공자입니다.', en: 'This AI provider is not supported.' },
  ai_api_key_required: { ko: 'API 키를 입력해 주세요.', en: 'Enter an API key.' },
  ai_settings_save_failed: { ko: 'BYOK 설정을 저장하지 못했습니다.', en: 'Could not save BYOK settings.' },
  ai_settings_missing_key: { ko: '저장된 BYOK 키가 없습니다.', en: 'No BYOK key is saved.' },
  ai_settings_enabled_update_failed: { ko: 'BYOK 사용 상태를 변경하지 못했습니다.', en: 'Could not update BYOK status.' },
  ai_settings_delete_failed: { ko: 'BYOK 키를 삭제하지 못했습니다.', en: 'Could not delete the BYOK key.' },
}

function getSettingsErrorMessage(payload: ApiErrorPayload, fallback: string, locale: 'ko' | 'en') {
  if (payload.error_code && settingsErrorMessages[payload.error_code]) {
    return settingsErrorMessages[payload.error_code][locale];
  }
  return payload.error ?? payload.message ?? fallback;
}

export function getDisplayName(user: UserInfo | null) {
  if (!user) return '학습자';
  return user.nickname || user.display_id || user.email || '학습자';
}

function formatDateInput(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, '0');
  const day = `${date.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
}

export function getTodayDateInput() {
  return formatDateInput(new Date());
}

export function useDashboardSettingsData(router: AppRouterInstance, redirectAfter: string) {
  const shouldLoadAIData = redirectAfter === '/dashboard/settings/ai';
  const shouldLoadPointUsage = redirectAfter === '/dashboard/settings/points';
  const [user, setUser] = useState<UserInfo | null>(null);
  const [aiSettings, setAISettings] = useState<UserAISettings | null>(null);
  const [aiUsageEvents, setAIUsageEvents] = useState<UserAIUsageEvent[]>([]);
  const [aiUsageTotal, setAIUsageTotal] = useState(0);
  const [aiUsagePage, setAIUsagePage] = useState(1);
  const [aiUsageTotalPages, setAIUsageTotalPages] = useState(0);
  const [aiUsageSummary, setAIUsageSummary] = useState<UserAIUsageSummary>({
    total_tokens: 0,
    average_daily_tokens: 0,
    total_estimated_cost_usd: 0,
    average_daily_estimated_cost_usd: 0,
  });
  const [aiUsageStartDate, setAIUsageStartDate] = useState(getTodayDateInput);
  const [aiUsageEndDate, setAIUsageEndDate] = useState(getTodayDateInput);
  const [isLoadingAIUsage, setIsLoadingAIUsage] = useState(false);
  const [aiUsageMessage, setAIUsageMessage] = useState<string | null>(null);
  const [pointUsageTransactions, setPointUsageTransactions] = useState<UserPointTransaction[]>([]);
  const [pointUsageTotal, setPointUsageTotal] = useState(0);
  const [pointUsagePage, setPointUsagePage] = useState(1);
  const [pointUsageTotalPages, setPointUsageTotalPages] = useState(0);
  const [pointUsageSummary, setPointUsageSummary] = useState<UserPointUsageSummary>({
    granted_points: 0,
    purchased_points: 0,
    used_points: 0,
    refunded_points: 0,
    net_change: 0,
  });
  const [pointUsageStartDate, setPointUsageStartDate] = useState(getTodayDateInput);
  const [pointUsageEndDate, setPointUsageEndDate] = useState(getTodayDateInput);
  const [isLoadingPointUsage, setIsLoadingPointUsage] = useState(false);
  const [pointUsageMessage, setPointUsageMessage] = useState<string | null>(null);
  const [profileNickname, setProfileNickname] = useState('');
  const [uiLocale, setUILocale] = useState<'ko' | 'en'>('ko');
  const [learningLanguage, setLearningLanguage] = useState<'ko' | 'en'>('ko');
  const [isSavingProfile, setIsSavingProfile] = useState(false);
  const [isSavingLanguagePreferences, setIsSavingLanguagePreferences] = useState(false);
  const [profileMessage, setProfileMessage] = useState<string | null>(null);
  const [languagePreferencesMessage, setLanguagePreferencesMessage] = useState<string | null>(null);
  const [avatarMessage, setAvatarMessage] = useState<string | null>(null);
  const [isUploadingAvatar, setIsUploadingAvatar] = useState(false);
  const [isWithdrawingAccount, setIsWithdrawingAccount] = useState(false);
  const [withdrawalMessage, setWithdrawalMessage] = useState<string | null>(null);
  const [aiProvider, setAIProvider] = useState('openai');
  const [aiAPIKey, setAIAPIKey] = useState('');
  const [aiEndpointURL, setAIEndpointURL] = useState('');
  const [isValidatingAI, setIsValidatingAI] = useState(false);
  const [isSavingAI, setIsSavingAI] = useState(false);
  const [isTogglingAIEnabled, setIsTogglingAIEnabled] = useState(false);
  const [isDeletingAIKey, setIsDeletingAIKey] = useState(false);
  const [aiMessage, setAIMessage] = useState<string | null>(null);
  const [aiKeyMessage, setAIKeyMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [validatedAIKeySignature, setValidatedAIKeySignature] = useState<string | null>(null);
  const currentAIKeySignature = `${aiProvider}::${aiEndpointURL.trim()}::${aiAPIKey.trim()}`;
  const isCurrentAIKeyValidated = Boolean(aiAPIKey.trim() && validatedAIKeySignature === currentAIKeySignature);
  const isEnglishUI = uiLocale === 'en';
  const isEnglishUIRef = useRef(isEnglishUI);

  useEffect(() => {
    isEnglishUIRef.current = isEnglishUI;
  }, [isEnglishUI]);

  const loadAIUsage = useCallback(async (startDate: string, endDate: string, page: number) => {
    setIsLoadingAIUsage(true);
    setAIUsageMessage(null);
    try {
      const errorLocale = isEnglishUIRef.current ? 'en' : 'ko';
      const params = new URLSearchParams({
        start_date: startDate,
        end_date: endDate,
        page: `${page}`,
        limit: '10',
      });
      const aiUsageRes = await fetch(`/api/v1/users/me/ai-usage?${params.toString()}`, {
        credentials: 'include',
        cache: 'no-store',
      });
      const aiUsageData = (await aiUsageRes.json().catch(() => ({}))) as UserAIUsageResponse;
      if (!aiUsageRes.ok) {
        throw new Error(getSettingsErrorMessage(aiUsageData, errorLocale === 'en' ? 'Could not load BYOK usage.' : 'BYOK 사용량을 불러오지 못했습니다.', errorLocale));
      }
      setAIUsageEvents(aiUsageData.events ?? []);
      setAIUsageTotal(aiUsageData.total ?? 0);
      setAIUsagePage(aiUsageData.page ?? page);
      setAIUsageTotalPages(aiUsageData.total_pages ?? 0);
      setAIUsageSummary({
        total_tokens: aiUsageData.summary?.total_tokens ?? 0,
        average_daily_tokens: aiUsageData.summary?.average_daily_tokens ?? 0,
        total_estimated_cost_usd: aiUsageData.summary?.total_estimated_cost_usd ?? 0,
        average_daily_estimated_cost_usd: aiUsageData.summary?.average_daily_estimated_cost_usd ?? 0,
      });
      if (aiUsageData.start_date) setAIUsageStartDate(aiUsageData.start_date);
      if (aiUsageData.end_date) setAIUsageEndDate(aiUsageData.end_date);
    } catch (error) {
      setAIUsageEvents([]);
      setAIUsageTotal(0);
      setAIUsagePage(1);
      setAIUsageTotalPages(0);
      setAIUsageSummary({
        total_tokens: 0,
        average_daily_tokens: 0,
        total_estimated_cost_usd: 0,
        average_daily_estimated_cost_usd: 0,
      });
      const errorLocale = isEnglishUIRef.current ? 'en' : 'ko';
      setAIUsageMessage(error instanceof Error ? error.message : (errorLocale === 'en' ? 'Could not load BYOK usage.' : 'BYOK 사용량을 불러오지 못했습니다.'));
    } finally {
      setIsLoadingAIUsage(false);
    }
  }, []);

  const loadPointUsage = useCallback(async (startDate: string, endDate: string, page: number) => {
    setIsLoadingPointUsage(true);
    setPointUsageMessage(null);
    try {
      const errorLocale = isEnglishUIRef.current ? 'en' : 'ko';
      const params = new URLSearchParams({
        start_date: startDate,
        end_date: endDate,
        page: String(page),
        limit: '10',
      });
      const pointUsageRes = await fetch('/api/v1/users/me/point-usage?' + params.toString(), {
        credentials: 'include',
        cache: 'no-store',
      });
      const pointUsageData = (await pointUsageRes.json().catch(() => ({}))) as UserPointUsageResponse;
      if (!pointUsageRes.ok) {
        throw new Error(getSettingsErrorMessage(pointUsageData, errorLocale === 'en' ? 'Could not load point history.' : '포인트 사용 내역을 불러오지 못했습니다.', errorLocale));
      }
      setPointUsageTransactions(pointUsageData.transactions ?? []);
      setPointUsageTotal(pointUsageData.total ?? 0);
      setPointUsagePage(pointUsageData.page ?? page);
      setPointUsageTotalPages(pointUsageData.total_pages ?? 0);
      setPointUsageSummary({
        granted_points: pointUsageData.summary?.granted_points ?? 0,
        purchased_points: pointUsageData.summary?.purchased_points ?? 0,
        used_points: pointUsageData.summary?.used_points ?? 0,
        refunded_points: pointUsageData.summary?.refunded_points ?? 0,
        net_change: pointUsageData.summary?.net_change ?? 0,
      });
      if (pointUsageData.balance) {
        setUser((current) => current
          ? {
              ...current,
              free_points: pointUsageData.balance?.free_points ?? current.free_points,
              paid_points: pointUsageData.balance?.paid_points ?? current.paid_points,
              total_points: pointUsageData.balance?.total_points ?? current.total_points,
            }
          : current);
      }
      if (pointUsageData.start_date) setPointUsageStartDate(pointUsageData.start_date);
      if (pointUsageData.end_date) setPointUsageEndDate(pointUsageData.end_date);
    } catch (error) {
      setPointUsageTransactions([]);
      setPointUsageTotal(0);
      setPointUsagePage(1);
      setPointUsageTotalPages(0);
      setPointUsageSummary({
        granted_points: 0,
        purchased_points: 0,
        used_points: 0,
        refunded_points: 0,
        net_change: 0,
      });
      const errorLocale = isEnglishUIRef.current ? 'en' : 'ko';
      setPointUsageMessage(error instanceof Error ? error.message : (errorLocale === 'en' ? 'Could not load point history.' : '포인트 사용 내역을 불러오지 못했습니다.'));
    } finally {
      setIsLoadingPointUsage(false);
    }
  }, []);

  useEffect(() => {
    if (!validatedAIKeySignature || currentAIKeySignature === validatedAIKeySignature) return;
    setValidatedAIKeySignature(null);
    setAIKeyMessage(null);
  }, [currentAIKeySignature, validatedAIKeySignature]);

  useEffect(() => {
    const load = async () => {
      const refreshRes = await fetch('/api/v1/auth/refresh', { method: 'POST', credentials: 'include' });
      if (!refreshRes.ok) {
        router.push(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      const meRes = await fetch('/api/v1/auth/me', { credentials: 'include', cache: 'no-store' });
      if (!meRes.ok) {
        router.push(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      const meData = (await meRes.json()) as UserInfo;
      if (meData.language_setup_required) {
        router.replace(`/language-setup?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }
      if (meData.required_consent_pending) {
        const agreementsPath = meData.ui_locale === 'en' ? '/en/agreements' : '/agreements';
        router.replace(`${agreementsPath}?redirect_after=${encodeURIComponent(redirectAfter)}`);
        return;
      }

      setUser(meData);
      setProfileNickname(meData.nickname ?? '');
      setUILocale(meData.ui_locale === 'en' ? 'en' : 'ko');
      setLearningLanguage(meData.learning_language === 'en' ? 'en' : 'ko');

      if (shouldLoadPointUsage) {
        const today = getTodayDateInput();
        setPointUsageStartDate(today);
        setPointUsageEndDate(today);
        await loadPointUsage(today, today, 1);
        setIsLoading(false);
        return;
      }

      if (!shouldLoadAIData) {
        setIsLoading(false);
        return;
      }

      const aiSettingsRes = await fetch('/api/v1/users/me/ai-settings', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (!aiSettingsRes.ok) {
        throw new Error(meData.ui_locale === 'en' ? 'Could not load BYOK settings.' : 'BYOK 설정을 불러오지 못했습니다.');
      }

      const aiSettingsData = (await aiSettingsRes.json()) as UserAISettings;
      setAISettings(aiSettingsData);
      setAIProvider(aiSettingsData.has_api_key ? aiSettingsData.provider || 'openai' : 'openai');
      setAIEndpointURL(aiSettingsData.has_api_key ? aiSettingsData.endpoint_url ?? '' : '');
      setAIAPIKey('');

      const today = getTodayDateInput();
      setAIUsageStartDate(today);
      setAIUsageEndDate(today);
      await loadAIUsage(today, today, 1);
      setIsLoading(false);
    };

    load().catch(() => router.push(`/login?redirect_after=${encodeURIComponent(redirectAfter)}`));
  }, [loadAIUsage, loadPointUsage, redirectAfter, router, shouldLoadAIData, shouldLoadPointUsage]);

  const handleAIUsageSearch = async () => {
    await loadAIUsage(aiUsageStartDate, aiUsageEndDate, 1);
  };

  const handleAIUsagePageChange = async (page: number) => {
    if (page < 1 || (aiUsageTotalPages > 0 && page > aiUsageTotalPages)) return;
    await loadAIUsage(aiUsageStartDate, aiUsageEndDate, page);
  };

  const handlePointUsageSearch = async () => {
    await loadPointUsage(pointUsageStartDate, pointUsageEndDate, 1);
  };

  const handlePointUsagePageChange = async (page: number) => {
    if (page < 1 || (pointUsageTotalPages > 0 && page > pointUsageTotalPages)) return;
    await loadPointUsage(pointUsageStartDate, pointUsageEndDate, page);
  };

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined);
    router.push('/');
  };

  const handleProfileSave = async () => {
    const nickname = profileNickname.trim();
    if (!nickname) {
      setProfileMessage(isEnglishUI ? 'Enter a nickname.' : '닉네임을 입력해 주세요.');
      return;
    }

    setIsSavingProfile(true);
    setProfileMessage(null);

    try {
      const response = await fetch('/api/v1/users/me', {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname }),
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload;
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not save profile.' : '학습자 정보를 저장하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      setUser((current) => (current ? { ...current, nickname } : current));
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      setProfileMessage(isEnglishUI ? 'Profile saved.' : (payload.message ?? '학습자 정보를 저장했습니다.'));
    } catch (error) {
      setProfileMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not save profile.' : '학습자 정보를 저장하지 못했습니다.'));
    } finally {
      setIsSavingProfile(false);
    }
  };

  const handleLanguagePreferencesSave = async () => {
    setIsSavingLanguagePreferences(true);
    setLanguagePreferencesMessage(null);

    try {
      const response = await fetch('/api/v1/users/me/preferences', {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ui_locale: uiLocale,
          learning_language: uiLocale,
        }),
      });
      const payload = (await response.json().catch(() => ({}))) as {
        error?: string;
        error_code?: string;
        ui_locale?: 'ko' | 'en';
        learning_language?: 'ko' | 'en';
      };
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not save language settings.' : '언어 설정을 저장하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      const nextUILocale = payload.ui_locale === 'en' ? 'en' : 'ko';
      const nextLearningLanguage = nextUILocale;
      setUILocale(nextUILocale);
      setLearningLanguage(nextLearningLanguage);
      setUser((current) => current
        ? {
            ...current,
            ui_locale: nextUILocale,
            learning_language: nextLearningLanguage,
            language_setup_required: false,
          }
        : current);
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      router.replace('/dashboard');
    } catch (error) {
      setLanguagePreferencesMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not save language settings.' : '언어 설정을 저장하지 못했습니다.'));
    } finally {
      setIsSavingLanguagePreferences(false);
    }
  };

  const handleAvatarUpload = async (file: File) => {
    setIsUploadingAvatar(true);
    setAvatarMessage(null);

    try {
      const formData = new FormData();
      formData.append('file', file);
      const response = await fetch('/api/v1/users/me/avatar', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload & { avatar_url?: string };
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not save profile image.' : '프로필 이미지를 저장하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }
      const avatarURL = payload.avatar_url ?? '';
      setUser((current) => current ? { ...current, avatar_url: avatarURL } : current);
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      setAvatarMessage(isEnglishUI ? 'Profile image saved.' : (payload.message ?? '프로필 이미지가 저장되었습니다.'));
    } catch (error) {
      setAvatarMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not save profile image.' : '프로필 이미지를 저장하지 못했습니다.'));
    } finally {
      setIsUploadingAvatar(false);
    }
  };

  const handleWithdrawAccount = async () => {
    setIsWithdrawingAccount(true);
    setWithdrawalMessage(null);

    try {
      const response = await fetch('/api/v1/users/me', {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload;
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not delete account.' : '계정 탈퇴 처리에 실패했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      document.cookie = 'is_logged_in=; Max-Age=0; path=/';
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      router.push('/?account=withdrawn');
    } catch (error) {
      setWithdrawalMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not delete account.' : '계정 탈퇴 처리에 실패했습니다.'));
    } finally {
      setIsWithdrawingAccount(false);
    }
  };

  const handleValidateAI = async () => {
    const apiKey = aiAPIKey.trim();
    if (!apiKey) {
      setAIKeyMessage(isEnglishUI ? 'Enter an API key to validate.' : '검증할 API 키를 입력해 주세요.');
      setValidatedAIKeySignature(null);
      return;
    }

    setIsValidatingAI(true);
    setAIKeyMessage(null);
    setValidatedAIKeySignature(null);

    try {
      const response = await fetch('/api/v1/users/me/ai-settings/validate', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: aiProvider,
          api_key: apiKey,
          endpoint_url: aiEndpointURL.trim() || null,
          is_enabled: true,
        }),
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload & { valid?: boolean };
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not validate the API key.' : 'API 키 검증에 실패했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      setAISettings((current) => ({
        mode: current?.mode ?? 'byok',
        provider: aiProvider,
        has_api_key: current?.has_api_key ?? false,
        is_enabled: current?.is_enabled ?? false,
        endpoint_url: aiEndpointURL.trim() || null,
        updated_at: current?.updated_at,
        last_validation_status: payload.valid ? 'valid' : 'invalid',
        last_validated_at: new Date().toISOString(),
        last_validation_error: payload.valid ? '' : payload.message,
        last_failed_at: payload.valid ? current?.last_failed_at : new Date().toISOString(),
        next_retry_at: current?.next_retry_at,
      }));
      setValidatedAIKeySignature(payload.valid ? currentAIKeySignature : null);
      setAIKeyMessage(payload.message ?? (isEnglishUI ? 'API key validated.' : 'API 키를 검증했습니다.'));
    } catch (error) {
      setValidatedAIKeySignature(null);
      setAIKeyMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not validate the API key.' : 'API 키 검증에 실패했습니다.'));
    } finally {
      setIsValidatingAI(false);
    }
  };

  const handleSaveAI = async () => {
    const apiKey = aiAPIKey.trim();
    if (!apiKey) {
      setAIKeyMessage(isEnglishUI ? 'Enter an API key to save.' : '저장할 API 키를 입력해 주세요.');
      return;
    }
    if (!isCurrentAIKeyValidated) {
      setAIKeyMessage(isEnglishUI ? 'Validate the API key first.' : 'API 키 검증을 먼저 완료해 주세요.');
      return;
    }

    setIsSavingAI(true);
    setAIKeyMessage(null);

    try {
      const response = await fetch('/api/v1/users/me/ai-settings', {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          provider: aiProvider,
          api_key: apiKey,
          endpoint_url: aiEndpointURL.trim() || null,
          is_enabled: true,
        }),
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload;
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not save BYOK settings.' : 'BYOK 설정을 저장하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      setAISettings({
        mode: 'byok',
        provider: aiProvider,
        has_api_key: true,
        is_enabled: true,
        endpoint_url: aiEndpointURL.trim() || null,
        updated_at: new Date().toISOString(),
        last_validation_status: aiSettings?.last_validation_status,
        last_validated_at: aiSettings?.last_validated_at,
        last_validation_error: aiSettings?.last_validation_error,
        last_failed_at: aiSettings?.last_failed_at,
        next_retry_at: aiSettings?.next_retry_at,
      });
      setAIAPIKey('');
      setValidatedAIKeySignature(null);
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      setAIKeyMessage(payload.message ?? (isEnglishUI ? 'BYOK settings saved.' : 'BYOK 설정을 저장했습니다.'));
    } catch (error) {
      setAIKeyMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not save BYOK settings.' : 'BYOK 설정을 저장하지 못했습니다.'));
    } finally {
      setIsSavingAI(false);
    }
  };

  const handleToggleAIEnabled = async (isEnabled: boolean) => {
    if (!aiSettings?.has_api_key) {
      setAIMessage(isEnglishUI ? 'No saved BYOK key.' : '저장된 BYOK 키가 없습니다.');
      return;
    }

    setIsTogglingAIEnabled(true);
    setAIMessage(null);

    try {
      const response = await fetch('/api/v1/users/me/ai-settings/enabled', {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ is_enabled: isEnabled }),
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload & { is_enabled?: boolean };
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not update BYOK status.' : 'BYOK 사용 상태를 변경하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      setAISettings((current) => current
        ? {
            ...current,
            mode: isEnabled ? 'byok' : 'managed_credit',
            is_enabled: payload.is_enabled ?? isEnabled,
            updated_at: new Date().toISOString(),
          }
        : current);
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      setAIMessage(isEnabled
        ? (isEnglishUI ? 'BYOK usage enabled.' : 'BYOK 사용을 활성화했습니다.')
        : (isEnglishUI ? 'BYOK usage disabled. The saved key is kept but not used for AI calls.' : 'BYOK 사용을 비활성화했습니다. 저장된 키는 보관되지만 AI 호출에는 사용하지 않습니다.'));
    } catch (error) {
      setAIMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not update BYOK status.' : 'BYOK 사용 상태를 변경하지 못했습니다.'));
    } finally {
      setIsTogglingAIEnabled(false);
    }
  };

  const handleDeleteAIKey = async () => {
    if (!aiSettings?.has_api_key) {
      setAIMessage(isEnglishUI ? 'No BYOK key to delete.' : '삭제할 BYOK 키가 없습니다.');
      return;
    }

    setIsDeletingAIKey(true);
    setAIMessage(null);

    try {
      const response = await fetch('/api/v1/users/me/ai-settings', {
        method: 'DELETE',
        credentials: 'include',
      });
      const payload = (await response.json().catch(() => ({}))) as ApiErrorPayload;
      if (!response.ok) {
        throw new Error(getSettingsErrorMessage(payload, isEnglishUI ? 'Could not delete the BYOK key.' : 'BYOK 키를 삭제하지 못했습니다.', isEnglishUI ? 'en' : 'ko'));
      }

      setAISettings({
        mode: 'managed_credit',
        provider: 'openai',
        has_api_key: false,
        is_enabled: false,
      });
      setAIProvider('openai');
      setAIEndpointURL('');
      setAIAPIKey('');
      setValidatedAIKeySignature(null);
      setAIKeyMessage(null);
      window.dispatchEvent(new Event('learnweaver:user-refresh'));
      setAIMessage(payload.message ?? (isEnglishUI ? 'Saved BYOK key deleted.' : '저장된 BYOK 키가 삭제되었습니다.'));
    } catch (error) {
      setAIMessage(error instanceof Error ? error.message : (isEnglishUI ? 'Could not delete the BYOK key.' : 'BYOK 키를 삭제하지 못했습니다.'));
    } finally {
      setIsDeletingAIKey(false);
    }
  };

  return {
    user,
    aiSettings,
    aiUsageEvents,
    aiUsageTotal,
    aiUsagePage,
    aiUsageTotalPages,
    aiUsageSummary,
    aiUsageStartDate,
    setAIUsageStartDate,
    aiUsageEndDate,
    setAIUsageEndDate,
    isLoadingAIUsage,
    aiUsageMessage,
    pointUsageTransactions,
    pointUsageTotal,
    pointUsagePage,
    pointUsageTotalPages,
    pointUsageSummary,
    pointUsageStartDate,
    setPointUsageStartDate,
    pointUsageEndDate,
    setPointUsageEndDate,
    isLoadingPointUsage,
    pointUsageMessage,
    profileNickname,
    setProfileNickname,
    uiLocale,
    setUILocale,
    learningLanguage,
    setLearningLanguage,
    isSavingProfile,
    isSavingLanguagePreferences,
    profileMessage,
    languagePreferencesMessage,
    avatarMessage,
    isUploadingAvatar,
    isWithdrawingAccount,
    withdrawalMessage,
    aiProvider,
    setAIProvider,
    aiAPIKey,
    setAIAPIKey,
    aiEndpointURL,
    setAIEndpointURL,
    isValidatingAI,
    isSavingAI,
    isTogglingAIEnabled,
    isDeletingAIKey,
    isCurrentAIKeyValidated,
    aiMessage,
    aiKeyMessage,
    isLoading,
    handleLogout,
    handleProfileSave,
    handleLanguagePreferencesSave,
    handleAvatarUpload,
    handleWithdrawAccount,
    handleValidateAI,
    handleSaveAI,
    handleToggleAIEnabled,
    handleDeleteAIKey,
    handleAIUsageSearch,
    handleAIUsagePageChange,
    handlePointUsageSearch,
    handlePointUsagePageChange,
  };
}
