export type RotationDirection = 'left' | 'right';

export interface PlanetTextureMapAsset {
  id: string;
  name: string;
  description: string;
  asset_path: string;
  public_url: string;
  width: number;
  height: number;
  columns: number;
  rows: number;
  cell_width: number;
  cell_height: number;
  rotation_duration_seconds: number;
  rotation_direction: RotationDirection;
  is_active: boolean;
  is_builtin: boolean;
  exists: boolean;
  size_bytes: number;
  updated_at?: string;
  version?: number;
  created_at?: string;
}

export interface TextureFrame {
  row: 0 | 1;
  col: 0 | 1 | 2 | 3;
  stage: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7;
  label: string;
}
