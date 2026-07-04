ALTER TABLE user_api_keys
    ADD COLUMN IF NOT EXISTS is_enabled BOOLEAN NOT NULL DEFAULT TRUE;

UPDATE policy_documents
SET content = replace(
        content,
        'LearnWeaver는 API 키를 민감 정보로 다루며, 평문 노출을 줄이기 위한 암호화 저장, 접근 제한, 로그 마스킹 등 보안 운영 기준을 적용하거나 강화해 나갑니다.',
        'LearnWeaver는 API 키를 민감 정보로 다룹니다. API 키 원문은 DB에 평문으로 저장하지 않으며, 서버에서 암호화하여 저장합니다. API 키는 AI 제공자 호출 등 서비스 제공에 필요한 순간에만 일시적으로 사용하고, 일반 조회 화면, API 응답, 운영 로그에 원문이 노출되지 않도록 관리합니다.'
    ),
    updated_at = NOW()
WHERE type = 'terms'
  AND content LIKE '%LearnWeaver는 API 키를 민감 정보로 다루며, 평문 노출을 줄이기 위한 암호화 저장, 접근 제한, 로그 마스킹 등 보안 운영 기준을 적용하거나 강화해 나갑니다.%';

UPDATE policy_documents
SET content = replace(
        replace(
            content,
            'LearnWeaver는 API 키를 평문 그대로 저장하지 않는 방향을 지향하며, 암호화 저장, 접근 제한, 로그 마스킹, 사용 후 폐기 등 보안 조치를 적용하거나 강화해 나갑니다.',
            'LearnWeaver는 API 키 원문을 DB에 평문으로 저장하지 않습니다. API 키는 서버에서 암호화하여 저장하며, AI 제공자 호출 등 서비스 제공에 필요한 순간에만 일시적으로 사용합니다. API 키 원문은 일반 조회 화면, API 응답, 운영 로그에 포함하지 않습니다.'
        ),
        'BYOK API 키는 가능한 한 평문 노출을 줄이고, 필요한 AI 호출 순간에만 사용되도록 설계합니다.',
        'BYOK API 키는 암호화 저장되며, 필요한 AI 호출 순간에만 사용되도록 설계합니다. 사용자가 BYOK를 비활성화하면 저장된 키는 보관되지만 AI 호출에는 사용하지 않으며, 해당 요청은 LearnWeaver 시스템 기본 키와 포인트 정책을 따릅니다.'
    ),
    updated_at = NOW()
WHERE type = 'privacy'
  AND (
      content LIKE '%LearnWeaver는 API 키를 평문 그대로 저장하지 않는 방향을 지향하며, 암호화 저장, 접근 제한, 로그 마스킹, 사용 후 폐기 등 보안 조치를 적용하거나 강화해 나갑니다.%'
      OR content LIKE '%BYOK API 키는 가능한 한 평문 노출을 줄이고, 필요한 AI 호출 순간에만 사용되도록 설계합니다.%'
  );
