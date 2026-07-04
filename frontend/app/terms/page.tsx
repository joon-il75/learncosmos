import PolicyDocumentPage from '@/components/policies/PolicyDocumentPage'

export default function TermsPage() {
  return (
    <PolicyDocumentPage
      type="terms"
      fallbackTitle="서비스 이용약관"
      loadingText="약관 문서를 불러오는 중..."
      errorText="약관 문서를 불러오지 못했습니다"
    />
  )
}
