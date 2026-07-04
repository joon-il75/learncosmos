ALTER TABLE recommendation_events
DROP CONSTRAINT IF EXISTS recommendation_events_event_type_check;

ALTER TABLE recommendation_events
ADD CONSTRAINT recommendation_events_event_type_check
CHECK (
    event_type IN (
        'curriculum_generated',
        'curriculum_deleted',
        'lesson_recommended',
        'lesson_selected',
        'lesson_rejected',
        'lesson_added_manually',
        'lesson_generated_by_ai',
        'curriculum_confirmed',
        'curriculum_edited',
        'curriculum_edited_after_start'
    )
);
