ALTER TABLE course_drafts
    ADD COLUMN IF NOT EXISTS planet_type_id UUID;

ALTER TABLE courses
    ADD COLUMN IF NOT EXISTS planet_type_id UUID;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_course_drafts_planet_type_id'
    ) THEN
        ALTER TABLE course_drafts
            ADD CONSTRAINT fk_course_drafts_planet_type_id
            FOREIGN KEY (planet_type_id) REFERENCES planet_types(id) ON DELETE SET NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_courses_planet_type_id'
    ) THEN
        ALTER TABLE courses
            ADD CONSTRAINT fk_courses_planet_type_id
            FOREIGN KEY (planet_type_id) REFERENCES planet_types(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_course_drafts_planet_type_id
    ON course_drafts (planet_type_id);

CREATE INDEX IF NOT EXISTS idx_courses_planet_type_id
    ON courses (planet_type_id);
