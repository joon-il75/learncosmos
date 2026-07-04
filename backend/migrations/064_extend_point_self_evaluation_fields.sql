-- Additive self-evaluation fields for the redesigned point learning screen.
-- Existing columns remain for compatibility with current records and clients.

ALTER TABLE course_point_self_evaluations
  ADD COLUMN IF NOT EXISTS understanding_score INT,
  ADD COLUMN IF NOT EXISTS understanding_reason TEXT,
  ADD COLUMN IF NOT EXISTS application_score INT,
  ADD COLUMN IF NOT EXISTS application_reason TEXT,
  ADD COLUMN IF NOT EXISTS proficiency_score INT,
  ADD COLUMN IF NOT EXISTS proficiency_reason TEXT,
  ADD COLUMN IF NOT EXISTS problem_solving_score INT,
  ADD COLUMN IF NOT EXISTS problem_solving_reason TEXT,
  ADD COLUMN IF NOT EXISTS expression_score INT,
  ADD COLUMN IF NOT EXISTS expression_reason TEXT;
