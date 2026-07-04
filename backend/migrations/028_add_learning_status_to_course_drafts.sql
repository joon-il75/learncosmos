-- 028: course_drafts에 learning 상태 추가
-- 행성탐험중(행성탐험 진행중) 상태를 표현하기 위해 learning 값 추가

ALTER TABLE course_drafts
  DROP CONSTRAINT course_drafts_status_check,
  ADD CONSTRAINT course_drafts_status_check
    CHECK (status IN ('draft', 'learning', 'confirmed', 'archived'));
