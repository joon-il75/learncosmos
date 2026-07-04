BEGIN;

-- Goal Profile v2: summarized_context, rebuild_decision, awaiting_rebuild_decision state 추가
ALTER TABLE course_goal_profiles
  ADD COLUMN summarized_context TEXT NOT NULL DEFAULT '',
  ADD COLUMN rebuild_decision VARCHAR(30)
    CHECK (rebuild_decision IN ('keep_structure', 'rebuild_remaining', 'rebuild_all'));

ALTER TABLE course_goal_profiles
  DROP CONSTRAINT course_goal_profiles_interview_state_check;

ALTER TABLE course_goal_profiles
  ADD CONSTRAINT course_goal_profiles_interview_state_check
    CHECK (interview_state IN (
      'listening',
      'clarifying',
      'proposing_goal',
      'confirmed',
      'revising_goal',
      'awaiting_rebuild_decision'
    ));

COMMIT;
