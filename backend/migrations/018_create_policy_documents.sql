CREATE TABLE IF NOT EXISTS policy_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type TEXT NOT NULL CHECK (type IN ('terms', 'privacy')),
    version INT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT true,
    is_active BOOLEAN NOT NULL DEFAULT false,
    effective_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (type, version)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_policy_documents_active_type
ON policy_documents(type)
WHERE is_active = true;

CREATE TABLE IF NOT EXISTS user_policy_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    policy_document_id UUID NOT NULL REFERENCES policy_documents(id) ON DELETE CASCADE,
    agreed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, policy_document_id)
);

CREATE INDEX IF NOT EXISTS idx_user_policy_consents_user_id
ON user_policy_consents(user_id);

INSERT INTO policy_documents (type, version, title, content, is_required, is_active, effective_at, published_at, updated_at)
SELECT
    seed.type,
    1,
    seed.title,
    seed.content,
    true,
    true,
    NOW(),
    NOW(),
    NOW()
FROM (
    VALUES
        (
            'terms',
            '서비스 이용약관',
            E'1. 목적\n본 약관은 LearnWeaver가 제공하는 학습 설계 및 콘텐츠 관리 서비스의 이용 조건과 운영 기준을 정합니다.\n\n2. 계정 이용\n이용자는 본인 소셜 계정을 통해 가입하며, 타인의 계정이나 인증 수단을 무단 사용해서는 안 됩니다.\n\n3. 금지 행위\n불법 콘텐츠 등록, 서비스 공격, 자동화 남용, 제3자 권리 침해, 운영 정책 우회 행위는 금지됩니다.\n\n4. AI 기능 이용\nBYOK 또는 플랫폼 제공 포인트를 통해 AI 기능을 사용할 수 있으며, 생성 결과의 적법성과 활용 책임은 이용자에게 있습니다.\n\n5. 서비스 변경\n운영상 필요 시 서비스 일부가 변경되거나 중단될 수 있으며, 중요한 변경은 별도 공지 또는 화면 고지로 안내합니다.'
        ),
        (
            'privacy',
            '개인정보처리방침',
            E'1. 수집 항목\n소셜 로그인 제공자가 전달한 식별자, 이메일, 닉네임, 프로필 이미지와 서비스 이용 중 생성되는 학습 설정 정보를 처리할 수 있습니다.\n\n2. 이용 목적\n회원 식별, 로그인 처리, 개인화된 학습 환경 제공, 고객 문의 대응, 보안 및 부정 이용 방지를 위해 개인정보를 이용합니다.\n\n3. 보관 및 보호\n관련 법령 또는 서비스 운영 목적에 필요한 기간 동안 정보를 보관하며, 민감한 설정 정보와 API 키는 별도 암호화 정책에 따라 보호합니다.\n\n4. 제3자 제공\n법령상 요구가 있는 경우를 제외하고, 이용자 동의 없이 개인정보를 외부에 판매하거나 임의 제공하지 않습니다.\n\n5. 이용자 권리\n이용자는 자신의 개인정보에 대한 조회, 수정, 삭제 요청을 할 수 있으며 관련 문의는 운영 채널을 통해 접수할 수 있습니다.'
        )
) AS seed(type, title, content)
WHERE NOT EXISTS (
    SELECT 1
    FROM policy_documents pd
    WHERE pd.type = seed.type
      AND pd.is_active = true
);
