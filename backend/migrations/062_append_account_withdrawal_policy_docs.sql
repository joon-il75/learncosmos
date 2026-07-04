UPDATE policy_documents
SET content = content || $append$

## 계정 탈퇴와 재가입

이용자는 설정 화면에서 계정 탈퇴를 요청할 수 있습니다. 탈퇴가 처리되면 계정은 비활성화되고 더 이상 해당 계정으로 서비스를 이용할 수 없습니다.

같은 소셜 계정으로 다시 가입하면 보관 기간 내 학습 기록을 다시 확인할 수 있습니다. 다만 개인정보와 보관 항목은 개인정보처리방침과 관련 법령에 따라 처리됩니다.
$append$,
    updated_at = NOW()
WHERE type = 'terms'
  AND content NOT LIKE '%## 계정 탈퇴와 재가입%';

UPDATE policy_documents
SET content = content || $append$

## 계정 탈퇴 시 개인정보 처리

이용자가 계정 탈퇴를 요청하면 LearnWeaver는 계정을 비활성화하고 로그인 및 서비스 이용을 중단합니다. 등록된 BYOK/API 키와 같이 계정 이용에 직접 필요한 민감 설정 정보는 삭제합니다.

학습 기록, 포인트 거래 기록, 동의 이력, 부정 이용 방지 및 분쟁 대응에 필요한 최소 기록은 관련 법령 또는 개인정보처리방침에서 정한 보관 기준에 따라 분리 또는 비활성 상태로 보관될 수 있습니다.

같은 소셜 계정으로 다시 가입하는 경우, 보관 기간 내 남아 있는 학습 기록을 다시 확인할 수 있습니다. 보관 기간이 끝났거나 법령상 보존 필요가 없는 개인정보는 지체 없이 파기합니다.
$append$,
    updated_at = NOW()
WHERE type = 'privacy'
  AND content NOT LIKE '%## 계정 탈퇴 시 개인정보 처리%';
