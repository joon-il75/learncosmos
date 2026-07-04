CREATE TABLE IF NOT EXISTS ui_engine_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    value JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO ui_engine_settings (key, value)
VALUES (
    'dashboard',
    jsonb_build_object(
        'compact_breakpoint', 1180,
        'phone_breakpoint', 520,
        'short_viewport_height', 460,
        'cta_min_launch_duration_ms', 920,
        'selection_cta_launch_delay_ms', 720,
        'lumi_quick_action_limit', 2,
        'lumi_desktop_panel_enabled', true,
        'lumi_mobile_sheet_enabled', true
    )
)
ON CONFLICT (key) DO NOTHING;
