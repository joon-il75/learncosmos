BEGIN;

-- Goal-Based Learning System: Goal Profiles & Revision Logs
-- 원칙:
-- - goal_profile은 course_draft와 1:N (버전 이력)
-- - is_active=true인 항목이 현재 활성 목표
-- - interview_state는 서버 주도 상태 머신으로 전이
-- - interview_messages는 JSONB 대화 이력 (role/content 배열)

CREATE TABLE course_goal_profiles (
  id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  course_draft_id     UUID        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
  user_id             UUID        NOT NULL,

  user_intent         TEXT        NOT NULL DEFAULT '',
  motivation          TEXT,
  usage_context       TEXT,
  confirmed_goal      TEXT,
  goal_type           VARCHAR(50),
  output_type         VARCHAR(50),
  difficulty_level    VARCHAR(20),
  time_horizon        VARCHAR(50),

  interview_state     VARCHAR(30) NOT NULL DEFAULT 'listening'
                        CHECK (interview_state IN ('listening','clarifying','proposing_goal','confirmed','revising_goal')),
  interview_messages  JSONB       NOT NULL DEFAULT '[]',

  version             INT         NOT NULL DEFAULT 1,
  is_active           BOOLEAN     NOT NULL DEFAULT true,

  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_goal_profiles_draft     ON course_goal_profiles(course_draft_id);
CREATE INDEX idx_goal_profiles_active    ON course_goal_profiles(course_draft_id, is_active) WHERE is_active = true;

CREATE TABLE course_goal_revision_logs (
  id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  course_draft_id       UUID        NOT NULL REFERENCES course_drafts(id) ON DELETE CASCADE,
  from_goal_profile_id  UUID        REFERENCES course_goal_profiles(id),
  to_goal_profile_id    UUID        NOT NULL REFERENCES course_goal_profiles(id),
  from_goal             TEXT,
  to_goal               TEXT        NOT NULL,
  reason                TEXT,
  rebuild_decision      VARCHAR(30) CHECK (rebuild_decision IN ('keep_structure','rebuild_remaining','rebuild_all')),
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_goal_revision_logs_draft ON course_goal_revision_logs(course_draft_id);

COMMIT;
