BEGIN;

CREATE TEMP TABLE _lw_policy_new_sets (
  document_type TEXT PRIMARY KEY,
  id UUID NOT NULL,
  version INT NOT NULL
) ON COMMIT DROP;

WITH next_versions AS (
  SELECT
    p.document_type,
    COALESCE(MAX(ps.version), 0) + 1 AS next_version
  FROM (VALUES ('terms'), ('privacy')) AS p(document_type)
  LEFT JOIN policy_document_sets ps
    ON ps.document_type = p.document_type
  GROUP BY p.document_type
), created_sets AS (
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
INSERT INTO _lw_policy_new_sets (document_type, id, version)
SELECT document_type, id, version
FROM created_sets;

WITH source_documents AS (
  SELECT
    pd.type,
    pd.locale,
    pd.title,
    CASE
      WHEN pd.locale = 'ko' THEN
        split_part(
          pd.content,
          E'\n## 외부 콘텐츠 접근 및 과금 범위\n',
          1
        )
      ELSE
        split_part(
          pd.content,
          E'\n## External Content Access and Charging Scope\n',
          1
        )
    END AS content_base,
    pd.translation_status,
    ns.id AS new_set_id
  FROM policy_documents pd
  JOIN policy_document_sets s
    ON s.id = pd.set_id
  JOIN _lw_policy_new_sets ns
    ON ns.document_type = s.document_type
  WHERE s.active = true
    AND s.required = true
    AND s.document_type IN ('terms', 'privacy')
    AND pd.is_active = true
    AND pd.locale IN ('ko', 'en')
), transformed AS (
  SELECT
    type,
    locale,
    title,
    new_set_id,
    translation_status,
    CASE
      WHEN locale = 'ko' THEN
        btrim(
          content_base
          || E'\n\n'
          || E'## 외부 콘텐츠 접근 및 과금 범위\n'
          || E'LearnCosmos는 외부 콘텐츠를 판매하지 않습니다.\n\n'
          || E'LearnCosmos는 외부 콘텐츠를 학습 행동으로 바꾸는 자기주도학습 도구를 제공합니다. 자체 플랫폼에서 업로드한 크리에이터 저작물은 별도 운영 기준으로 제공합니다.\n\n'
          || E'추천 영상은 유료/무료에 관계없이 동일한 방식으로 접근할 수 있으며, 학습자 화면에서 유료 CTA를 영상 플레이어 위나 재생 버튼 앞에 배치하지 않습니다.\n\n'
          || E'결제 혜택은 영상이 아니라 AI 피드백, 학습 기록 분석, 복습 루틴, 결과물 코칭, 개인 학습 리포트에만 연결됩니다.\n\n'
          || E'추천 결과에는 YouTube 출처와 원본 링크, 공식 플레이어 링크를 함께 표시해야 합니다.'
        )
      ELSE
        btrim(
          content_base
          || E'\n\n'
          || E'## External Content Access and Charging Scope\n'
          || E'LearnCosmos does not sell external content.\n\n'
          || E'LearnCosmos is a self-directed learning tool that converts external content into learning actions. Creator content uploaded on LearnCosmos itself is handled under separate operational rules.\n\n'
          || E'External videos are always accessible through the same flow, and paid CTAs are not shown above the player or before the play button.\n\n'
          || E'Payment benefits are only tied to AI feedback, learning record analysis, review routines, artifact coaching, and personal learning reports.\n\n'
          || E'Recommendation results must clearly display the YouTube source, original link, and official player link.'
        )
    END AS content
  FROM source_documents
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
  ns.version,
  title,
  content,
  true,
  true,
  false,
  locale,
  translation_status,
  NOW(),
  NOW(),
  NOW()
FROM transformed t
JOIN _lw_policy_new_sets ns
  ON ns.document_type = t.type;

UPDATE policy_documents pd
SET is_active = false,
    updated_at = NOW()
WHERE pd.type IN ('terms', 'privacy')
  AND pd.is_active = true
  AND pd.set_id NOT IN (SELECT id FROM _lw_policy_new_sets);

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
FROM _lw_policy_new_sets ns
WHERE pds.id = ns.id;

COMMIT;
