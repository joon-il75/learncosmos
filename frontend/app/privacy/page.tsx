import PolicyDocumentPage from '@/components/policies/PolicyDocumentPage'

export default function PrivacyPage() {
  return (
    <PolicyDocumentPage
      type="privacy"
      fallbackTitle="개인정보처리방침"
      loadingText="개인정보처리방침을 불러오는 중..."
      errorText="개인정보처리방침을 불러오지 못했습니다"
    />
  )
}
