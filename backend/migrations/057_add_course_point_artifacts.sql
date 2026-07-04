CREATE TABLE IF NOT EXISTS course_point_artifacts (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    course_point_id UUID        NOT NULL REFERENCES course_points(id) ON DELETE CASCADE,
    artifact_type   TEXT        NOT NULL DEFAULT '',
    title           TEXT        NOT NULL DEFAULT '',
    url             TEXT        NOT NULL DEFAULT '',
    description     TEXT        NOT NULL DEFAULT '',
    order_index     INTEGER     NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_course_point_artifacts_point_order
    ON course_point_artifacts(course_point_id, order_index, created_at);

INSERT INTO course_point_artifacts (
    course_point_id,
    artifact_type,
    title,
    url,
    description,
    order_index,
    created_at,
    updated_at
)
SELECT
    course_point_id,
    artifact_type,
    title,
    url,
    description,
    0,
    created_at,
    updated_at
FROM course_point_artifact_entries
WHERE (title <> '' OR url <> '' OR description <> '')
  AND NOT EXISTS (
      SELECT 1
      FROM course_point_artifacts existing
      WHERE existing.course_point_id = course_point_artifact_entries.course_point_id
        AND existing.order_index = 0
  )
ON CONFLICT DO NOTHING;
