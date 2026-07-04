'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

import DashboardSettingsShell, {
  actionRowStyle,
  fieldLabelStyle,
  helperTextStyle,
  inlineInfoStyle,
  loadingCardStyle,
  loadingPageStyle,
  primaryActionButtonStyle,
  sectionCardStyle,
  sectionHeaderStyle,
  sectionSubtitleStyle,
  sectionTitleStyle,
  settingsInputStyle,
} from '@/components/dashboard/settings/DashboardSettingsShell';
import { formatDateTime, useDashboardSettingsData } from '@/components/dashboard/settings/useDashboardSettingsData';
import { getDashboardSettingsProfileCopy } from '@/lib/i18n/pages/dashboardSettingsProfile';

export default function DashboardSettingsProfilePage() {
  const router = useRouter();
  const [withdrawConfirm, setWithdrawConfirm] = useState('');
  const [isWithdrawalOpen, setIsWithdrawalOpen] = useState(false);
  const [isCompactProfileLayout, setIsCompactProfileLayout] = useState(false);
  const {
    user,
    profileNickname,
    setProfileNickname,
    uiLocale,
    setUILocale,
    setLearningLanguage,
    isSavingProfile,
    isSavingLanguagePreferences,
    profileMessage,
    languagePreferencesMessage,
    avatarMessage,
    isUploadingAvatar,
    isWithdrawingAccount,
    withdrawalMessage,
    isLoading,
    handleLogout,
    handleProfileSave,
    handleLanguagePreferencesSave,
    handleAvatarUpload,
    handleWithdrawAccount,
  } = useDashboardSettingsData(router, '/dashboard/settings/profile');
  const copy = getDashboardSettingsProfileCopy(uiLocale);
  const isSuccessMessage = (message: string) => message.includes('저장') || message.toLowerCase().includes('saved');
  const withdrawalConfirmValue = copy.withdrawal.confirmValue;

  useEffect(() => {
    const syncLayout = () => setIsCompactProfileLayout(window.innerWidth <= 640);
    syncLayout();
    window.addEventListener('resize', syncLayout);
    return () => window.removeEventListener('resize', syncLayout);
  }, []);

  if (isLoading) {
    return (
      <div style={loadingPageStyle}>
        <div style={loadingCardStyle}>{copy.loading}</div>
      </div>
    );
  }

  return (
    <DashboardSettingsShell
      eyebrow={copy.shell.eyebrow}
      title={copy.shell.title}
      copy={copy.shell.copy}
      tabLabels={copy.shell.tabs}
      uiLocale={uiLocale}
      activeTab="profile"
      onLogout={handleLogout}
    >
      <section style={profileSurfaceStyle}>
        <div style={{ ...compactHeaderStyle, ...(isCompactProfileLayout ? compactHeaderMobileStyle : {}) }}>
          <div>
            <h2 style={compactTitleStyle}>{copy.profile.title}</h2>
            <p style={compactSubtitleStyle}>{copy.profile.subtitle}</p>
          </div>
          <div style={planTextStyle}>{user?.premium_access ? 'LearnCosmos Pro' : copy.profile.standardPlan}</div>
        </div>

        <div style={{ ...avatarPanelStyle, ...(isCompactProfileLayout ? avatarPanelMobileStyle : {}) }}>
          <div style={avatarPreviewStyle}>
            {user?.avatar_url ? (
              <img src={user.avatar_url} alt={copy.profile.avatarAlt} style={avatarImageStyle} />
            ) : (
              <span style={avatarInitialStyle}>{(user?.nickname || user?.display_id || user?.email || copy.profile.fallbackInitial).slice(0, 1)}</span>
            )}
          </div>
          <div style={avatarTextStyle}>
            <strong style={avatarTitleStyle}>{copy.profile.avatarTitle}</strong>
            <p style={avatarDescriptionStyle}>{copy.profile.avatarDescription}</p>
            <label style={avatarUploadButtonStyle(isUploadingAvatar)}>
              {isUploadingAvatar ? copy.profile.avatarUploading : copy.profile.avatarSelect}
              <input
                type="file"
                accept="image/jpeg,image/png,image/webp"
                disabled={isUploadingAvatar}
                style={{ display: 'none' }}
                onChange={(event) => {
                  const file = event.target.files?.[0];
                  event.target.value = '';
                  if (file) void handleAvatarUpload(file);
                }}
              />
            </label>
            {avatarMessage ? <div style={inlineInfoStyle(isSuccessMessage(avatarMessage) ? '#48BB78' : '#FFB4A2')}>{avatarMessage}</div> : null}
          </div>
        </div>

        <div style={{ ...profileLayoutStyle, ...(isCompactProfileLayout ? profileLayoutMobileStyle : {}) }}>
          <div style={readOnlyInfoStyle}>
            <div style={{ ...infoRowStyle, ...(isCompactProfileLayout ? infoRowMobileStyle : {}) }}>
              <span style={infoLabelStyle}>{copy.profile.labels.email}</span>
              <span style={infoValueStyle}>{user?.email || copy.profile.unset}</span>
            </div>
            <div style={{ ...infoRowStyle, ...(isCompactProfileLayout ? infoRowMobileStyle : {}) }}>
              <span style={infoLabelStyle}>{copy.profile.labels.provider}</span>
              <span style={infoValueStyle}>{user?.provider || copy.profile.defaultProvider}</span>
            </div>
            <div style={{ ...infoRowStyle, ...(isCompactProfileLayout ? infoRowMobileStyle : {}) }}>
              <span style={infoLabelStyle}>{copy.profile.labels.createdAt}</span>
              <span style={infoValueStyle}>{formatDateTime(user?.created_at, uiLocale)}</span>
            </div>
            <div style={{ ...infoRowStyle, ...(isCompactProfileLayout ? infoRowMobileStyle : {}) }}>
              <span style={infoLabelStyle}>{copy.profile.labels.points}</span>
              <span style={infoValueStyle}>
                {copy.profile.pointTotal(user?.total_points ?? 0)} <span style={{ ...subValueStyle, ...(isCompactProfileLayout ? subValueMobileStyle : {}) }}>{copy.profile.pointBreakdown(user?.free_points ?? 0, user?.paid_points ?? 0)}</span>
              </span>
            </div>
          </div>

          <div style={editPanelStyle}>
            <label style={{ ...fieldLabelStyle, gap: '10px' }}>
              {copy.profile.labels.editableNickname}
              <input
                type="text"
                value={profileNickname}
                onChange={(event) => {
                  setProfileNickname(event.target.value);
                }}
                maxLength={20}
                style={{
                  ...settingsInputStyle,
                  minHeight: '44px',
                  border: '1px solid rgba(106, 210, 193, 0.56)',
                  background: 'rgba(255, 255, 255, 0.86)',
                  boxShadow: 'inset 0 0 0 1px rgba(255, 255, 255, 0.5), 0 0 0 3px rgba(45, 142, 108, 0.08)',
                }}
              />
            </label>
            <p style={editHintStyle}>{copy.profile.editHint}</p>

            {profileMessage ? <div style={inlineInfoStyle(isSuccessMessage(profileMessage) ? '#48BB78' : '#FFB4A2')}>{profileMessage}</div> : null}

            <div style={actionRowStyle}>
              <button type="button" onClick={handleProfileSave} style={primaryActionButtonStyle(isSavingProfile)} disabled={isSavingProfile}>
                {isSavingProfile ? copy.profile.saving : copy.profile.saveNickname}
              </button>
            </div>
          </div>
        </div>
      </section>

      <section style={sectionCardStyle}>
        <div style={sectionHeaderStyle}>
          <div>
            <h2 style={sectionTitleStyle}>{copy.language.title}</h2>
            <p style={sectionSubtitleStyle}>{copy.language.subtitle}</p>
          </div>
        </div>

        <div style={languageGridStyle}>
          <div style={languagePanelStyle}>
            <div>
              <h3 style={languageTitleStyle}>{copy.language.optionTitle}</h3>
              <p style={languageHelpStyle}>{copy.language.optionHelp}</p>
            </div>
            <div style={languageButtonRowStyle}>
              <button
                type="button"
                onClick={() => {
                  setUILocale('ko');
                  setLearningLanguage('ko');
                }}
                style={languageOptionStyle(uiLocale === 'ko')}
              >
                한국어
              </button>
              <button
                type="button"
                onClick={() => {
                  setUILocale('en');
                  setLearningLanguage('en');
                }}
                style={languageOptionStyle(uiLocale === 'en')}
              >
                English
              </button>
            </div>
          </div>
        </div>

        {languagePreferencesMessage ? (
          <div style={inlineInfoStyle(isSuccessMessage(languagePreferencesMessage) ? '#48BB78' : '#FFB4A2')}>
            {languagePreferencesMessage}
          </div>
        ) : null}

        <div style={actionRowStyle}>
          <button
            type="button"
            onClick={handleLanguagePreferencesSave}
            disabled={isSavingLanguagePreferences}
            style={primaryActionButtonStyle(isSavingLanguagePreferences)}
          >
            {isSavingLanguagePreferences ? copy.language.saving : copy.language.save}
          </button>
        </div>
      </section>

      <section
        style={{
          ...sectionCardStyle,
          gap: isWithdrawalOpen ? '16px' : '0',
          borderColor: 'rgba(185, 28, 28, 0.42)',
          borderLeft: '5px solid rgba(185, 28, 28, 0.72)',
          background: 'linear-gradient(180deg, rgba(255, 241, 242, 0.98), rgba(254, 226, 226, 0.94))',
          boxShadow: '0 18px 44px rgba(127, 29, 29, 0.14)',
        }}
      >
        <button
          type="button"
          onClick={() => setIsWithdrawalOpen((open) => !open)}
          aria-expanded={isWithdrawalOpen}
          style={{
            ...sectionHeaderStyle,
            width: '100%',
            padding: 0,
            border: 'none',
            background: 'transparent',
            textAlign: 'left',
            cursor: 'pointer',
            fontFamily: 'inherit',
          }}
        >
          <div>
            <h2 style={{ ...sectionTitleStyle, fontSize: '18px', color: '#7F1D1D' }}>{copy.withdrawal.title}</h2>
            <p style={{ ...sectionSubtitleStyle, color: '#8A3B2F' }}>{copy.withdrawal.subtitle}</p>
          </div>
          <span style={toggleTextStyle}>{isWithdrawalOpen ? copy.withdrawal.collapse : copy.withdrawal.expand}</span>
        </button>

        {isWithdrawalOpen ? (
          <>
            <div style={withdrawalWarningStyle}>
              {copy.withdrawal.warning}
            </div>

            <label style={{ ...fieldLabelStyle, marginTop: '6px' }}>
              {copy.withdrawal.confirmLabel}
              <input
                type="text"
                value={withdrawConfirm}
                onChange={(event) => setWithdrawConfirm(event.target.value)}
                placeholder={copy.withdrawal.confirmPlaceholder}
                style={settingsInputStyle}
              />
            </label>

            {withdrawalMessage ? <div style={inlineInfoStyle('#FFB4A2')}>{withdrawalMessage}</div> : null}

            <div style={actionRowStyle}>
              <button
                type="button"
                onClick={handleWithdrawAccount}
                disabled={isWithdrawingAccount || withdrawConfirm.trim() !== withdrawalConfirmValue}
                style={{
                  ...primaryActionButtonStyle(isWithdrawingAccount || withdrawConfirm.trim() !== '탈퇴합니다'),
                  background:
                    isWithdrawingAccount || withdrawConfirm.trim() !== withdrawalConfirmValue
                      ? 'rgba(146, 111, 62, 0.18)'
                      : 'linear-gradient(135deg, #F97373, #CC5216)',
                  color: isWithdrawingAccount || withdrawConfirm.trim() !== withdrawalConfirmValue ? '#6A5A46' : '#FFF7ED',
                }}
              >
                {isWithdrawingAccount ? copy.withdrawal.processing : copy.withdrawal.action}
              </button>
            </div>
          </>
        ) : null}
      </section>
    </DashboardSettingsShell>
  );
}

const profileSurfaceStyle = {
  display: 'grid',
  gap: '20px',
  borderRadius: '18px',
  border: '1px solid rgba(255, 255, 255, 0.42)',
  background: 'rgba(255, 253, 247, 0.96)',
  padding: '24px',
  boxShadow: '0 24px 60px rgba(3, 8, 20, 0.22)',
} satisfies React.CSSProperties;

const compactHeaderStyle = {
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: '18px',
  paddingBottom: '16px',
  borderBottom: '1px solid rgba(146, 111, 62, 0.14)',
} satisfies React.CSSProperties;

const compactHeaderMobileStyle = {
  flexDirection: 'column',
  gap: '10px',
} satisfies React.CSSProperties;

const compactTitleStyle = {
  margin: 0,
  fontSize: '24px',
  lineHeight: 1.25,
  color: '#160E08',
  fontWeight: 800,
} satisfies React.CSSProperties;

const avatarPanelStyle = {
  display: 'grid',
  gridTemplateColumns: '96px minmax(0, 1fr)',
  gap: '18px',
  alignItems: 'center',
  padding: '18px',
  borderRadius: '18px',
  border: '1px solid rgba(45, 142, 108, 0.18)',
  background: 'rgba(255, 248, 235, 0.54)',
} satisfies React.CSSProperties;

const avatarPanelMobileStyle = {
  gridTemplateColumns: '1fr',
  justifyItems: 'center',
  textAlign: 'center',
} satisfies React.CSSProperties;

const avatarPreviewStyle = {
  width: '96px',
  height: '96px',
  borderRadius: '999px',
  overflow: 'hidden',
  border: '3px solid rgba(45, 142, 108, 0.32)',
  background: 'linear-gradient(135deg, rgba(13, 49, 78, 0.96), rgba(28, 125, 121, 0.92))',
  display: 'grid',
  placeItems: 'center',
} satisfies React.CSSProperties;

const avatarImageStyle = {
  width: '100%',
  height: '100%',
  objectFit: 'cover',
} satisfies React.CSSProperties;

const avatarInitialStyle = {
  color: '#FFFDF7',
  fontSize: '34px',
  fontWeight: 900,
} satisfies React.CSSProperties;

const avatarTextStyle = {
  display: 'grid',
  gap: '8px',
} satisfies React.CSSProperties;

const avatarTitleStyle = {
  color: '#160E08',
  fontSize: '18px',
  lineHeight: 1.35,
} satisfies React.CSSProperties;

const avatarDescriptionStyle = {
  margin: 0,
  color: '#35291E',
  fontSize: '14px',
  lineHeight: 1.65,
} satisfies React.CSSProperties;

const avatarUploadButtonStyle = (disabled: boolean): React.CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  justifySelf: 'start',
  minHeight: '40px',
  padding: '0 16px',
  borderRadius: '999px',
  border: '1px solid rgba(13, 49, 78, 0.24)',
  background: disabled ? 'rgba(146, 111, 62, 0.18)' : 'rgba(255, 253, 247, 0.86)',
  color: '#160E08',
  fontSize: '14px',
  fontWeight: 800,
  cursor: disabled ? 'default' : 'pointer',
});

const compactSubtitleStyle = {
  margin: '6px 0 0',
  fontSize: '15px',
  lineHeight: 1.7,
  color: '#35291E',
  fontWeight: 500,
} satisfies React.CSSProperties;

const planTextStyle = {
  fontSize: '14px',
  fontWeight: 800,
  color: '#0F5C44',
  whiteSpace: 'nowrap',
} satisfies React.CSSProperties;

const profileLayoutStyle = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) minmax(280px, 360px)',
  gap: '24px',
  alignItems: 'start',
} satisfies React.CSSProperties;

const profileLayoutMobileStyle = {
  gridTemplateColumns: '1fr',
  gap: '18px',
} satisfies React.CSSProperties;

const readOnlyInfoStyle = {
  display: 'grid',
  borderTop: '1px solid rgba(146, 111, 62, 0.14)',
} satisfies React.CSSProperties;

const infoRowStyle = {
  display: 'grid',
  gridTemplateColumns: '120px minmax(0, 1fr)',
  gap: '16px',
  padding: '14px 0',
  borderBottom: '1px solid rgba(146, 111, 62, 0.14)',
} satisfies React.CSSProperties;

const infoRowMobileStyle = {
  gridTemplateColumns: '1fr',
  gap: '6px',
  padding: '13px 0',
} satisfies React.CSSProperties;

const infoLabelStyle = {
  fontSize: '13px',
  color: '#4A3724',
  fontWeight: 700,
} satisfies React.CSSProperties;

const infoValueStyle = {
  fontSize: '15px',
  color: '#130D07',
  fontWeight: 800,
  overflowWrap: 'anywhere',
} satisfies React.CSSProperties;

const subValueStyle = {
  marginLeft: '8px',
  fontSize: '13px',
  fontWeight: 700,
  color: '#4A3724',
} satisfies React.CSSProperties;

const subValueMobileStyle = {
  display: 'block',
  marginLeft: 0,
  marginTop: '4px',
} satisfies React.CSSProperties;

const editPanelStyle = {
  display: 'grid',
  gap: '12px',
  borderRadius: '16px',
  border: '1px solid rgba(13, 125, 121, 0.34)',
  background: 'rgba(232, 250, 244, 0.86)',
  padding: '16px',
} satisfies React.CSSProperties;

const editHintStyle = {
  margin: 0,
  fontSize: '13px',
  lineHeight: 1.7,
  color: '#27443A',
  fontWeight: 600,
} satisfies React.CSSProperties;

const languageGridStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(260px, 1fr))',
  gap: '14px',
} satisfies React.CSSProperties;

const languagePanelStyle = {
  display: 'grid',
  gap: '14px',
  padding: '16px',
  borderRadius: '16px',
  border: '1px solid rgba(13, 49, 78, 0.16)',
  background: 'rgba(255, 248, 235, 0.68)',
} satisfies React.CSSProperties;

const languageTitleStyle = {
  margin: 0,
  color: '#160E08',
  fontSize: '17px',
  lineHeight: 1.35,
  fontWeight: 800,
} satisfies React.CSSProperties;

const languageHelpStyle = {
  margin: '6px 0 0',
  color: '#4A3724',
  fontSize: '13px',
  lineHeight: 1.65,
  fontWeight: 600,
} satisfies React.CSSProperties;

const languageButtonRowStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(2, minmax(0, 1fr))',
  gap: '8px',
} satisfies React.CSSProperties;

const languageOptionStyle = (active: boolean): React.CSSProperties => ({
  minHeight: '44px',
  padding: '0 12px',
  borderRadius: '14px',
  border: active ? '2px solid rgba(13, 49, 78, 0.72)' : '1px solid rgba(13, 49, 78, 0.18)',
  background: active ? 'linear-gradient(135deg, rgba(13, 49, 78, 0.96), rgba(28, 125, 121, 0.92))' : 'rgba(255, 255, 255, 0.84)',
  color: active ? '#FFFDF7' : '#1F2933',
  fontFamily: 'inherit',
  fontSize: '14px',
  fontWeight: 900,
  cursor: 'pointer',
});

const toggleTextStyle = {
  minHeight: '32px',
  display: 'inline-flex',
  alignItems: 'center',
  padding: '0 12px',
  borderRadius: '999px',
  border: '1px solid rgba(185, 28, 28, 0.38)',
  background: 'rgba(254, 226, 226, 0.72)',
  color: '#7F1D1D',
  fontSize: '13px',
  fontWeight: 700,
} satisfies React.CSSProperties;

const withdrawalWarningStyle = {
  padding: '14px 16px',
  borderRadius: '14px',
  border: '1px solid rgba(185, 28, 28, 0.32)',
  background: 'rgba(255, 255, 255, 0.48)',
  color: '#7F1D1D',
  fontSize: '13px',
  lineHeight: 1.75,
  fontWeight: 600,
} satisfies React.CSSProperties;
