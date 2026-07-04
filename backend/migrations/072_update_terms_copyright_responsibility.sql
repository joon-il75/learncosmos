WITH terms_v4 AS (
    SELECT
        pd.title,
        replace(
            pd.content,
            E'이용자는 외부 콘텐츠를 사용할 때 관련 법령과 저작권을 준수해야 합니다.',
            E'이용자는 외부 콘텐츠를 사용할 때 관련 법령과 저작권을 준수해야 합니다.\n\n이용자가 연구지점, 학습 콘텐츠, 첨부파일, 학습 기록 등에 직접 작성하거나 업로드하는 자료의 저작권 및 이용 권한 확인 책임은 이용자에게 있습니다. 이용자는 자신이 권리를 보유했거나 적법하게 사용할 수 있는 자료만 작성·첨부해야 하며, 타인의 저작물은 관련 법령, 라이선스, 출처의 이용 조건이 허용하는 범위에서만 사용해야 합니다.'
        ) AS content
    FROM policy_documents pd
    WHERE pd.type = 'terms'
      AND pd.version = 3
    ORDER BY pd.updated_at DESC
    LIMIT 1
), deactivated AS (
    UPDATE policy_documents pd
    SET is_active = false,
        updated_at = NOW()
    WHERE pd.type = 'terms'
      AND pd.version <> 4
      AND pd.is_active = true
    RETURNING pd.id
), updated_v4 AS (
    UPDATE policy_documents pd
    SET title = terms_v4.title,
        content = terms_v4.content,
        is_required = true,
        is_active = true,
        is_draft = false,
        effective_at = TIMESTAMPTZ '2026-04-30 00:00:00+09',
        published_at = COALESCE(pd.published_at, NOW()),
        updated_at = NOW()
    FROM terms_v4
    WHERE pd.type = 'terms'
      AND pd.version = 4
      AND (SELECT COUNT(*) FROM deactivated) >= 0
    RETURNING pd.id
)
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
    4,
    terms_v4.title,
    terms_v4.content,
    true,
    true,
    false,
    TIMESTAMPTZ '2026-04-30 00:00:00+09',
    NOW(),
    NOW()
FROM terms_v4
WHERE NOT EXISTS (SELECT 1 FROM updated_v4)
  AND (SELECT COUNT(*) FROM deactivated) >= 0;
