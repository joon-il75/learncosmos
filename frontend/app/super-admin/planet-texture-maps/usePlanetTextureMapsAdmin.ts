'use client';

import { useRouter } from 'next/navigation';
import { useEffect, useMemo, useRef, useState, type ChangeEvent } from 'react';
import {
  RECOMMENDED_HEIGHT,
  RECOMMENDED_WIDTH,
  SAMPLE_ATLAS_URL,
} from './constants';
import type { PlanetTextureMapAsset, RotationDirection } from './types';

export function usePlanetTextureMapsAdmin() {
  const router = useRouter();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);
  const [textureMaps, setTextureMaps] = useState<PlanetTextureMapAsset[]>([]);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [busyTextureMapId, setBusyTextureMapId] = useState<string | null>(null);
  const [name, setName] = useState('기본 행성 텍스처맵');
  const [description, setDescription] = useState('');
  const [previewURL, setPreviewURL] = useState<string | null>(SAMPLE_ATLAS_URL);
  const [fileName, setFileName] = useState<string | null>('learnweaver_planet_atlas_basic_seamfixed_2048x768.webp');
  const [imageSize, setImageSize] = useState<{ width: number; height: number } | null>({
    width: RECOMMENDED_WIDTH,
    height: RECOMMENDED_HEIGHT,
  });
  const [progress, setProgress] = useState(50);
  const [rotationDuration, setRotationDuration] = useState(36);
  const [rotationDirection, setRotationDirection] = useState<RotationDirection>('left');
  const [isPlaying, setIsPlaying] = useState(true);
  const [isActive, setIsActive] = useState(true);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const replaceFileInputRefs = useRef<Record<string, HTMLInputElement | null>>({});

  const applyRegisteredTextureMap = (textureMap: PlanetTextureMapAsset) => {
    if (previewURL?.startsWith('blob:')) URL.revokeObjectURL(previewURL);
    const versionSuffix = textureMap.version ? `?v=${textureMap.version}` : '';
    setPreviewURL(`${textureMap.public_url}${versionSuffix}`);
    setFileName(textureMap.public_url.split('/').pop() ?? textureMap.name);
    setImageSize({ width: textureMap.width, height: textureMap.height });
    setName(textureMap.name);
    setDescription(textureMap.description);
    setRotationDuration(textureMap.rotation_duration_seconds);
    setRotationDirection(textureMap.rotation_direction);
    setIsActive(textureMap.is_active);
  };

  useEffect(() => {
  const loadTextureMaps = async () => {
      const res = await fetch('/api/v1/super-admin/assets/planet-texture-maps', {
        credentials: 'include',
        cache: 'no-store',
      });

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }

      if (!res.ok) {
        setStatusMessage('행성 텍스처맵 목록을 불러오지 못했습니다.');
        setIsAuthorized(true);
        return;
      }

      const payload = await res.json() as { planet_texture_maps?: PlanetTextureMapAsset[] };
      const maps = payload.planet_texture_maps ?? [];
      setTextureMaps(maps);
      const primaryMap = maps.find((item) => item.is_active) ?? maps[0];
      if (primaryMap) {
        applyRegisteredTextureMap(primaryMap);
      }
      setIsAuthorized(true);
    };

    loadTextureMaps().catch(() => {
      setStatusMessage('행성 텍스처맵 목록을 불러오지 못했습니다.');
      setIsAuthorized(true);
    });
  // applyRegisteredTextureMap intentionally reads the current preview URL for object URL cleanup.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [router]);

  useEffect(() => {
    return () => {
      if (previewURL?.startsWith('blob:')) URL.revokeObjectURL(previewURL);
    };
  }, [previewURL]);

  const activeCount = useMemo(() => textureMaps.filter((item) => item.is_active).length, [textureMaps]);
  const atlasIsRecommended = imageSize?.width === RECOMMENDED_WIDTH && imageSize?.height === RECOMMENDED_HEIGHT;

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null;
    if (!file) return;
    setSelectedFile(file);
    if (previewURL?.startsWith('blob:')) URL.revokeObjectURL(previewURL);
    const objectURL = URL.createObjectURL(file);
    setPreviewURL(objectURL);
    setFileName(file.name);
    setImageSize(null);

    const image = new Image();
    image.onload = () => {
      setImageSize({ width: image.naturalWidth, height: image.naturalHeight });
    };
    image.src = objectURL;
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setStatusMessage('등록할 WebP atlas 파일을 먼저 선택해 주세요.');
      return;
    }
    if (!name.trim()) {
      setStatusMessage('텍스처맵 이름을 입력해 주세요.');
      return;
    }


    const formData = new FormData();
    formData.append('file', selectedFile);
    formData.append('name', name.trim());
    formData.append('description', description.trim());
    formData.append('is_active', String(isActive));
    formData.append('rotation_duration_seconds', String(rotationDuration));
    formData.append('rotation_direction', rotationDirection);

    setIsUploading(true);
    setStatusMessage(null);

    try {
      const res = await fetch('/api/v1/super-admin/assets/planet-texture-maps', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }

      const payload = await res.json().catch(() => ({})) as { error?: string; planet_texture_map?: PlanetTextureMapAsset };
      if (!res.ok || !payload.planet_texture_map) {
        setStatusMessage(payload.error ?? '행성 텍스처맵 등록에 실패했습니다.');
        return;
      }

      setTextureMaps((current) => [payload.planet_texture_map!, ...current]);
      setSelectedFile(null);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      applyRegisteredTextureMap(payload.planet_texture_map);
      setStatusMessage('행성 텍스처맵을 등록했습니다.');
    } catch {
      setStatusMessage('행성 텍스처맵 등록에 실패했습니다.');
    } finally {
      setIsUploading(false);
    }
  };

  const handleToggleActive = async (textureMap: PlanetTextureMapAsset) => {
    setBusyTextureMapId(textureMap.id);
    setStatusMessage(null);
    try {
      const res = await fetch(`/api/v1/super-admin/assets/planet-texture-maps/${textureMap.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ is_active: !textureMap.is_active }),
      });

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }

      const payload = await res.json().catch(() => ({})) as { error?: string; planet_texture_map?: PlanetTextureMapAsset };
      if (!res.ok || !payload.planet_texture_map) {
        setStatusMessage(payload.error ?? '행성 텍스처맵 상태를 바꾸지 못했습니다.');
        return;
      }

      setTextureMaps((current) => current.map((item) => item.id === textureMap.id ? payload.planet_texture_map! : item));
      applyRegisteredTextureMap(payload.planet_texture_map);
      setStatusMessage(`'${payload.planet_texture_map.name}' 상태를 ${payload.planet_texture_map.is_active ? '활성' : '비활성'}으로 변경했습니다.`);
    } catch {
      setStatusMessage('행성 텍스처맵 상태를 바꾸지 못했습니다.');
    } finally {
      setBusyTextureMapId(null);
    }
  };

  const handleReplaceFileChange = async (textureMap: PlanetTextureMapAsset, event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] ?? null;
    event.target.value = '';
    if (!file) return;
    const formData = new FormData();
    formData.append('file', file);

    setBusyTextureMapId(textureMap.id);
    setStatusMessage(null);
    try {
      const res = await fetch(`/api/v1/super-admin/assets/planet-texture-maps/${textureMap.id}/file`, {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }

      const payload = await res.json().catch(() => ({})) as { error?: string; planet_texture_map?: PlanetTextureMapAsset };
      if (!res.ok || !payload.planet_texture_map) {
        setStatusMessage(payload.error ?? '행성 텍스처맵 이미지를 교체하지 못했습니다.');
        return;
      }

      setTextureMaps((current) => current.map((item) => item.id === textureMap.id ? payload.planet_texture_map! : item));
      applyRegisteredTextureMap(payload.planet_texture_map);
      setStatusMessage(`'${payload.planet_texture_map.name}' 이미지를 교체했습니다.`);
    } catch {
      setStatusMessage('행성 텍스처맵 이미지를 교체하지 못했습니다.');
    } finally {
      setBusyTextureMapId(null);
    }
  };

  const handleDelete = async (textureMap: PlanetTextureMapAsset) => {
    if (!window.confirm(`'${textureMap.name}' 텍스처맵을 삭제할까요? 업로드된 atlas 파일도 함께 제거됩니다.`)) {
      return;
    }
    setBusyTextureMapId(textureMap.id);
    setStatusMessage(null);
    try {
      const res = await fetch(`/api/v1/super-admin/assets/planet-texture-maps/${textureMap.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }

      const payload = await res.json().catch(() => ({})) as { error?: string };
      if (!res.ok) {
        setStatusMessage(payload.error ?? '행성 텍스처맵 삭제에 실패했습니다.');
        return;
      }

      const nextMaps = textureMaps.filter((item) => item.id !== textureMap.id);
      setTextureMaps(nextMaps);
      const fallbackMap = nextMaps.find((item) => item.is_active) ?? nextMaps[0];
      if (fallbackMap) {
        applyRegisteredTextureMap(fallbackMap);
      } else {
        setPreviewURL(SAMPLE_ATLAS_URL);
        setFileName('learnweaver_planet_atlas_basic_seamfixed_2048x768.webp');
        setImageSize({ width: RECOMMENDED_WIDTH, height: RECOMMENDED_HEIGHT });
      }
      setStatusMessage(`'${textureMap.name}' 텍스처맵을 삭제했습니다.`);
    } catch {
      setStatusMessage('행성 텍스처맵 삭제에 실패했습니다.');
    } finally {
      setBusyTextureMapId(null);
    }
  };

  return {
    activeCount,
    applyRegisteredTextureMap,
    atlasIsRecommended,
    busyTextureMapId,
    description,
    fileInputRef,
    fileName,
    handleDelete,
    handleFileChange,
    handleReplaceFileChange,
    handleToggleActive,
    handleUpload,
    imageSize,
    isActive,
    isAuthorized,
    isPlaying,
    isUploading,
    name,
    previewURL,
    progress,
    replaceFileInputRefs,
    rotationDirection,
    rotationDuration,
    selectedFile,
    setDescription,
    setIsActive,
    setIsPlaying,
    setName,
    setProgress,
    setRotationDirection,
    setRotationDuration,
    statusMessage,
    textureMaps,
  };
}
