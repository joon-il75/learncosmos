WITH latest_active AS (
    SELECT DISTINCT ON (type)
        id,
        type,
        version
    FROM policy_documents
    WHERE is_active = true
    ORDER BY type, version DESC
),
upsert_terms_draft AS (
    INSERT INTO policy_documents (
        id,
        type,
        version,
        title,
        content,
        is_required,
        is_active,
        is_draft,
        effective_at,
        published_at,
        updated_at
    )
    SELECT
        gen_random_uuid(),
        'terms',
        COALESCE((SELECT version + 1 FROM latest_active WHERE type = 'terms'), 1),
        'LearnWeaver 서비스 이용약관',
        E'제1조 목적\n이 약관은 LearnWeaver 운영주체(이하 "회사")가 제공하는 학습 설계, 콘텐츠 관리, AI 기반 학습 보조, 포인트 및 구독 서비스의 이용과 관련하여 회사와 이용자 간 권리, 의무 및 책임사항을 정하는 것을 목적으로 합니다.\n\n제2조 회원가입 및 계정\n이용자는 회사가 제공하는 소셜 로그인 방식으로 가입할 수 있으며, 회사가 요구하는 필수 약관 및 개인정보처리방침 동의를 완료해야 주요 학습 기능을 이용할 수 있습니다. 계정 식별은 소셜 로그인 제공자의 식별자 기준으로 처리될 수 있습니다.\n\n제3조 서비스 내용\n회사는 학습 목표 설정, 커리큘럼 초안 생성, 콘텐츠 등록 및 관리, AI 학습 보조, 포인트 및 구독 기능 등을 제공합니다. 회사는 운영상 필요에 따라 서비스의 전부 또는 일부를 변경하거나 중단할 수 있습니다.\n\n제4조 AI 기능 이용\n서비스의 AI 기능은 학습 보조 도구이며 결과의 정확성, 완전성, 최신성을 보장하지 않습니다. 이용자는 생성 결과를 스스로 검토하고 자신의 책임으로 활용해야 합니다.\n\n제5조 BYOK 이용\n이용자는 자신이 적법하게 사용할 권한이 있는 외부 AI 서비스 API 키만 등록할 수 있습니다. BYOK 키 사용에 따른 외부 서비스 이용약관 준수, 과금 및 법적 책임은 이용자에게 있습니다.\n\n제6조 포인트 및 구독\n회사는 운영 정책에 따라 웰컴 포인트, 유료 포인트, LearnWeaver Pro 구독 혜택을 제공할 수 있습니다. 포인트의 부여, 차감, 사용 단가, 유효성, 구독 혜택은 회사 정책 및 별도 고지에 따릅니다.\n\n제7조 금지행위\n이용자는 타인 계정 도용, 불법 콘텐츠 등록, 권리 침해, 자동화 남용, 보안 우회, 서비스 운영 방해, 등록한 API 키를 이용한 불법 행위 등을 해서는 안 됩니다.\n\n제8조 이용제한 및 해지\n회사는 이용자가 관련 법령, 본 약관 또는 운영 정책을 위반한 경우 경고, 기능 제한, 이용 정지, 계약 해지 등의 조치를 할 수 있습니다. 이용자는 회사가 정한 절차에 따라 탈퇴를 요청할 수 있습니다.\n\n제9조 책임 제한\n회사는 천재지변, 통신장애, 외부 서비스 장애, 이용자 귀책사유 등 회사의 합리적 통제를 벗어난 사유로 발생한 손해에 대해 책임을 지지 않습니다. 무료 서비스에 대해서는 특별한 사정이 없는 한 손해배상 책임을 제한합니다.\n\n제10조 준거법 및 관할\n본 약관은 대한민국 법령에 따라 해석되며, 회사와 이용자 간 분쟁은 관련 법령에 따른 관할 법원을 따릅니다.',
        true,
        false,
        true,
        NOW(),
        NOW(),
        NOW()
    WHERE NOT EXISTS (
        SELECT 1
        FROM policy_documents
        WHERE type = 'terms' AND is_draft = true
    )
    ON CONFLICT DO NOTHING
),
upsert_privacy_draft AS (
    INSERT INTO policy_documents (
        id,
        type,
        version,
        title,
        content,
        is_required,
        is_active,
        is_draft,
        effective_at,
        published_at,
        updated_at
    )
    SELECT
        gen_random_uuid(),
        'privacy',
        COALESCE((SELECT version + 1 FROM latest_active WHERE type = 'privacy'), 1),
        'LearnWeaver 개인정보처리방침',
        E'1. 처리하는 개인정보 항목\n회사는 소셜 로그인 식별자, 이메일, 닉네임, 프로필 이미지, 서비스 이용 기록, 포인트 및 구독 정보, 정책 문서 동의 이력, 접속 로그, 보안 정보, BYOK 설정 정보와 API 키 암호화 저장값 등을 처리할 수 있습니다.\n\n2. 처리 목적\n회사는 회원 식별, 로그인 및 인증, 서비스 제공, 학습 기능 운영, AI 기능 제공, 정책 문서 동의 이력 관리, 고객 문의 대응, 보안 및 부정 이용 방지를 위해 개인정보를 처리합니다.\n\n3. 보유 및 이용기간\n회사는 개인정보 처리 목적 달성 시까지 개인정보를 보유·이용합니다. 다만 관계 법령상 보관 의무가 있거나 분쟁 대응, 보안 점검, 정책 동의 증빙이 필요한 경우 일정 기간 보관할 수 있습니다.\n\n4. 제3자 제공 및 처리위탁\n회사는 원칙적으로 이용자의 개인정보를 외부에 제공하지 않습니다. 다만 법령상 근거가 있거나 이용자 동의가 있는 경우 예외로 합니다. 서비스 운영에 필요한 범위에서 처리위탁이 발생할 경우 관련 법령에 따라 공개하고 관리합니다.\n\n5. 국외 이전 가능성\nAI 기능, 클라우드, 외부 기술 서비스 연동 과정에서 개인정보 또는 관련 정보가 국외에서 처리될 가능성이 있습니다. 실제 국외 이전이 발생하는 경우 이전받는 자, 국가, 항목, 목적, 보유기간 등을 별도로 고지합니다.\n\n6. 파기절차 및 방법\n회사는 개인정보 보유기간 경과 또는 처리 목적 달성 시 지체 없이 개인정보를 파기합니다. 전자적 파일은 복구 불가능한 방식으로 삭제하고, 출력물은 분쇄 또는 소각합니다.\n\n7. 정보주체 권리\n이용자는 자신의 개인정보에 대한 열람, 정정, 삭제, 처리정지, 동의철회 등을 요청할 수 있습니다. 회사는 관련 법령에 따라 지체 없이 조치합니다.\n\n8. 안전성 확보조치\n회사는 접근권한 관리, 암호화 저장, 접속기록 관리, 보안 점검, 관리자 인증 강화 등 개인정보 보호를 위한 기술적·관리적 조치를 시행합니다.\n\n9. BYOK 및 AI 처리\n이용자가 등록한 외부 AI 서비스 API 키는 서비스 제공 목적 범위에서 암호화 저장될 수 있습니다. 이용자 입력 정보는 AI 기능 제공에 필요한 범위에서 외부 AI 서비스 제공자에게 전달될 수 있으며, 외부 서비스의 이용약관과 과금 책임은 이용자에게 있습니다.\n\n10. 정책 문서 동의 이력\n회사는 서비스 이용약관 및 개인정보처리방침에 대한 동의 여부, 동의 시각, 동의 대상 문서 버전, 접속 정보 등을 법적 의무 이행, 서비스 이용 자격 확인 및 분쟁 대응 목적으로 보관할 수 있습니다.\n\n11. 문의처\n개인정보 보호책임자, 담당부서, 문의 이메일 등 구체 정보는 운영 주체가 확정 후 반영합니다.',
        true,
        false,
        true,
        NOW(),
        NOW(),
        NOW()
    WHERE NOT EXISTS (
        SELECT 1
        FROM policy_documents
        WHERE type = 'privacy' AND is_draft = true
    )
    ON CONFLICT DO NOTHING
)
SELECT 1;
