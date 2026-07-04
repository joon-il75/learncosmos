'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import type { LumiAssetMeta } from './lumiLabTypes';

export type LumiLabAssetsResult = {
  isAuthorized: boolean | null;
  assetMeta: LumiAssetMeta | null;
  assetVersion: number;
  assetStatus: string | null;
  selectedUpload: File | null;
  setSelectedUpload: (f: File | null) => void;
  isUploading: boolean;
  copiedPromptId: string | null;
  handleUpload: () => Promise<void>;
  handleCopyPrompt: (id: string, value: string) => Promise<void>;
};

export function useLumiLabAssets(params: {
  onRuntimeConfigLoad: () => Promise<void>;
  onRuntimeConfigError: (msg: string) => void;
}): LumiLabAssetsResult {
  const router = useRouter();
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);
  const [assetMeta, setAssetMeta] = useState<LumiAssetMeta | null>(null);
  const [assetVersion, setAssetVersion] = useState<number>(() => Date.now());
  const [assetStatus, setAssetStatus] = useState<string | null>(null);
  const [selectedUpload, setSelectedUpload] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [copiedPromptId, setCopiedPromptId] = useState<string | null>(null);

  const onRuntimeConfigLoadRef = useRef(params.onRuntimeConfigLoad);
  onRuntimeConfigLoadRef.current = params.onRuntimeConfigLoad;
  const onRuntimeConfigErrorRef = useRef(params.onRuntimeConfigError);
  onRuntimeConfigErrorRef.current = params.onRuntimeConfigError;

  useEffect(() => {
    const loadMeta = async () => {
      const res = await fetch('/api/v1/super-admin/assets/lumi/meta', {
        credentials: 'include',
        cache: 'no-store',
      });
      if (res.status === 401 || res.status === 403) {
        router.replace('/super-admin/login');
        return;
      }
      if (!res.ok) {
        setAssetStatus('Lumi 자산 정보를 불러오지 못했습니다.');
        setIsAuthorized(true);
        return;
      }
      const data = await res.json() as LumiAssetMeta;
      setAssetMeta(data);
      setAssetVersion(data.version ?? Date.now());
      await onRuntimeConfigLoadRef.current();
      setIsAuthorized(true);
    };

    loadMeta().catch((error) => {
      setAssetStatus('Lumi 자산 정보를 불러오지 못했습니다.');
      onRuntimeConfigErrorRef.current(
        error instanceof Error ? error.message : 'Lumi runtime 설정을 불러오지 못했습니다.',
      );
      setIsAuthorized(true);
    });
  }, [router]);

  const handleUpload = async () => {
    if (!selectedUpload) {
      setAssetStatus('업로드할 WebP 파일을 먼저 선택해 주세요.');
      return;
    }

    const formData = new FormData();
    formData.append('file', selectedUpload);
    setIsUploading(true);
    setAssetStatus(null);
    try {
      const res = await fetch('/api/v1/super-admin/assets/lumi/upload', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      if (res.status === 401 || res.status === 403) { router.replace('/super-admin/login'); return; }
      const data = await res.json().catch(() => ({}));
      if (!res.ok) { setAssetStatus(data.error || 'Lumi 자산 업로드에 실패했습니다.'); return; }
      const nextMeta = data.asset as LumiAssetMeta;
      setAssetMeta(nextMeta);
      setAssetVersion(nextMeta.version ?? Date.now());
      setAssetStatus('Lumi 자산을 교체했습니다.');
      setSelectedUpload(null);
    } catch {
      setAssetStatus('Lumi 자산 업로드에 실패했습니다.');
    } finally {
      setIsUploading(false);
    }
  };

  const handleCopyPrompt = async (id: string, value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      setCopiedPromptId(id);
      window.setTimeout(() => setCopiedPromptId((current) => (current === id ? null : current)), 1600);
    } catch {
      setAssetStatus('클립보드 복사에 실패했습니다.');
    }
  };

  return {
    isAuthorized, assetMeta, assetVersion, assetStatus,
    selectedUpload, setSelectedUpload, isUploading, copiedPromptId,
    handleUpload, handleCopyPrompt,
  };
}
