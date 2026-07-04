CREATE TABLE IF NOT EXISTS planet_texture_maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    asset_path TEXT NOT NULL,
    width INTEGER NOT NULL DEFAULT 2048 CHECK (width > 0),
    height INTEGER NOT NULL DEFAULT 768 CHECK (height > 0),
    columns INTEGER NOT NULL DEFAULT 4 CHECK (columns > 0),
    rows INTEGER NOT NULL DEFAULT 3 CHECK (rows > 0),
    cell_width INTEGER NOT NULL DEFAULT 512 CHECK (cell_width > 0),
    cell_height INTEGER NOT NULL DEFAULT 256 CHECK (cell_height > 0),
    rotation_duration_seconds INTEGER NOT NULL DEFAULT 36 CHECK (rotation_duration_seconds BETWEEN 18 AND 90),
    rotation_direction TEXT NOT NULL DEFAULT 'left' CHECK (rotation_direction IN ('left', 'right')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_planet_texture_maps_active_created_at
    ON planet_texture_maps (is_active, created_at DESC);

INSERT INTO planet_texture_maps (
    id,
    name,
    description,
    asset_path,
    width,
    height,
    columns,
    rows,
    cell_width,
    cell_height,
    rotation_duration_seconds,
    rotation_direction,
    is_active
)
VALUES (
    '8f7a8f12-4c7f-4e0f-9f6c-2e9b1f5d0c31',
    '기본 행성 텍스처맵',
    '기본 4x3 자전용 행성 표면 atlas',
    '/textures/planets/maps/learnweaver_planet_atlas_basic_seamfixed_2048x768.webp',
    2048,
    768,
    4,
    3,
    512,
    256,
    36,
    'left',
    true
)
ON CONFLICT (id) DO NOTHING;
