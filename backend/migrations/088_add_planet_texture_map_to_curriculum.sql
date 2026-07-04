ALTER TABLE course_drafts
    ADD COLUMN IF NOT EXISTS planet_texture_map_id UUID;

ALTER TABLE courses
    ADD COLUMN IF NOT EXISTS planet_texture_map_id UUID;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_course_drafts_planet_texture_map_id'
    ) THEN
        ALTER TABLE course_drafts
            ADD CONSTRAINT fk_course_drafts_planet_texture_map_id
            FOREIGN KEY (planet_texture_map_id) REFERENCES planet_texture_maps(id) ON DELETE SET NULL;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_courses_planet_texture_map_id'
    ) THEN
        ALTER TABLE courses
            ADD CONSTRAINT fk_courses_planet_texture_map_id
            FOREIGN KEY (planet_texture_map_id) REFERENCES planet_texture_maps(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_course_drafts_planet_texture_map_id
    ON course_drafts (planet_texture_map_id);

CREATE INDEX IF NOT EXISTS idx_courses_planet_texture_map_id
    ON courses (planet_texture_map_id);
