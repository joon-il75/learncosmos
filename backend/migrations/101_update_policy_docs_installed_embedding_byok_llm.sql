BEGIN;

CREATE TEMP TABLE _lw_policy_new_sets (
  document_type TEXT PRIMARY KEY,
  id UUID NOT NULL,
  version INT NOT NULL
) ON COMMIT DROP;

WITH next_versions AS (
  SELECT
    t.document_type,
    COALESCE(MAX(pds.version), 0) + 1 AS next_version
  FROM (VALUES ('terms'), ('privacy')) AS t(document_type)
  LEFT JOIN policy_document_sets pds ON pds.document_type = t.document_type
  GROUP BY t.document_type
),
inserted_sets AS (
  INSERT INTO policy_document_sets (
    id,
    document_type,
    version,
    required,
    active,
    effective_at,
    created_at,
    updated_at
  )
  SELECT
    gen_random_uuid(),
    document_type,
    next_version,
    true,
    false,
    NOW(),
    NOW(),
    NOW()
  FROM next_versions
  ON CONFLICT (document_type, version) DO UPDATE
  SET required = EXCLUDED.required,
      effective_at = EXCLUDED.effective_at,
      updated_at = NOW()
  RETURNING id, document_type, version
)
INSERT INTO _lw_policy_new_sets(document_type, id, version)
SELECT document_type, id, version
FROM inserted_sets;

WITH source_docs AS (
  SELECT
    pd.type,
    pd.locale,
    pd.title,
    pd.content,
    pd.translation_status,
    pns.id AS new_set_id,
    pns.version AS new_version
  FROM policy_document_sets pds
  JOIN policy_documents pd
    ON pd.set_id = pds.id
   AND pd.is_active = true
  JOIN _lw_policy_new_sets pns
    ON pns.document_type = pds.document_type
  WHERE pds.active = true
    AND pds.required = true
    AND pds.document_type IN ('terms', 'privacy')
    AND pd.locale IN ('ko', 'en')
),
transformed_docs AS (
  SELECT
    type,
    locale,
    title,
    new_set_id,
    new_version,
    translation_status,
    CASE
      WHEN locale = 'ko' THEN
        replace(
          replace(
            replace(
              replace(
                content,
                '베타 기간에는 OpenAI BYOK만 먼저 접수합니다. Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE 등은 향후 연결 예정 제공자이며, 실제 지원 제공자, 모델, 기능 범위, 포인트 차감 여부는 서비스 화면과 운영 정책에 따라 단계적으로 제공됩니다.',
                'OpenAI를 시작으로 Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE 등 다양한 LLM 제공자는 단계적으로 연결될 수 있습니다. 실제 지원 제공자, 모델, 기능 범위, 포인트 차감 여부는 서비스 화면과 운영 정책에 따라 제공됩니다.'
              ),
              '베타 기간에는 OpenAI BYOK만 먼저 접수합니다. Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE 등은 향후 연결 예정 제공자이며, 실제 지원 제공자와 모델 범위는 서비스 화면 및 운영 정책에 따라 단계적으로 제공됩니다.',
              'OpenAI를 시작으로 Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE 등 다양한 LLM 제공자는 단계적으로 연결될 수 있습니다. 실제 지원 제공자와 모델 범위는 서비스 화면 및 운영 정책에 따라 제공됩니다.'
            ),
            '현재 OpenAI BYOK 사용은 LearnWeaver가 정한 범위 내에서 포인트를 차감하지 않는 방식으로 운영될 수 있습니다.',
            '현재 BYOK LLM 사용은 LearnWeaver가 정한 범위 내에서 포인트를 차감하지 않는 방식으로 운영될 수 있습니다.'
          ),
          '추천 검색은 BYOK 키를 쓰지 않고 EmbeddingGemma primary를 사용하지만 OpenAI BYOK 활성 사용자는 제품 정책상 0pt이며 BYOK 사용량 기록은 남기지 않음',
          '콘텐츠 검색과 추천은 Google EmbeddingGemma 기반 설치형 임베딩 경로를 사용하며, 사용자가 등록한 BYOK LLM 키와 분리됩니다.'
        ) ||
        E'\n\n콘텐츠 검색과 추천은 Google EmbeddingGemma 기반 설치형 임베딩 모델을 통해 처리되며, 목표 채팅, 코스 생성, 학습 코칭 등 생성형 AI 호출에 사용하는 BYOK LLM 키와 사용 범위가 분리됩니다.'
      ELSE
        replace(
          replace(
            replace(
              replace(
                content,
                'During beta operation, LearnWeaver accepts OpenAI BYOK first. Users may connect their own OpenAI API key for supported AI features.',
                'LearnWeaver starts with OpenAI and may connect Claude, Gemini, Grok, HyperCLOVA X, Solar, Llama, EXAONE, and other LLM providers in stages. Supported providers, models, feature scope, and point policies are provided through service screens and operating policy.'
              ),
              'During beta operation, OpenAI BYOK is accepted first. Other AI providers may be connected gradually after operational review.',
              'OpenAI starts first, and other LLM providers may be connected gradually after operational review.'
            ),
            'OpenAI BYOK users may use supported goal, course generation, point AI, and recommendation search features without point deduction according to the current beta policy. Recommendation search may still use LearnWeaver''s embedding/search infrastructure rather than the user''s BYOK key.',
            'BYOK LLM users may use supported goal, course generation, and point AI features without point deduction according to the current operating policy. Content recommendations and similar-resource search use LearnWeaver''s installed embedding/search infrastructure rather than the user''s BYOK LLM key.'
          ),
          'OpenAI BYOK',
          'BYOK LLM'
        ) ||
        E'\n\nContent search and recommendations are handled through a Google EmbeddingGemma-based installed embedding model on the LearnWeaver server, separate from the scope of user-registered BYOK LLM keys used for generative AI calls such as goal chat, course creation, and learning coaching.'
    END AS content
  FROM source_docs
)
INSERT INTO policy_documents (
  id,
  set_id,
  type,
  version,
  title,
  content,
  is_required,
  is_active,
  is_draft,
  locale,
  translation_status,
  effective_at,
  published_at,
  updated_at
)
SELECT
  gen_random_uuid(),
  new_set_id,
  type,
  new_version,
  title,
  btrim(content, E' \n\r\t'),
  true,
  true,
  false,
  locale,
  translation_status,
  NOW(),
  NOW(),
  NOW()
FROM transformed_docs
ON CONFLICT DO NOTHING;

UPDATE policy_documents pd
SET is_active = false,
    updated_at = NOW()
WHERE pd.type IN ('terms', 'privacy')
  AND pd.set_id NOT IN (SELECT id FROM _lw_policy_new_sets)
  AND pd.is_active = true;

UPDATE policy_document_sets
SET active = false,
    updated_at = NOW()
WHERE document_type IN ('terms', 'privacy')
  AND active = true;

UPDATE policy_document_sets pds
SET active = true,
    required = true,
    effective_at = NOW(),
    updated_at = NOW()
FROM _lw_policy_new_sets pns
WHERE pds.id = pns.id;

COMMIT;
