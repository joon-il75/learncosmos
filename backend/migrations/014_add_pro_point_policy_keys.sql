INSERT INTO point_settings (key, value) VALUES
    ('pro_monthly_points', 300),
    ('tutor_cost', 1),
    ('vision_cost', 5),
    ('quiz_cost', 2)
ON CONFLICT (key) DO NOTHING;
