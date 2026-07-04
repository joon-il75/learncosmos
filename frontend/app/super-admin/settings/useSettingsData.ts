'use client'
import { useState, useEffect } from 'react'
import type {
  LLMSetting, APIKeyStatus, PointPolicy, Product,
  AdSlot, AffiliateSettings, UIEngineSettings, PolicyDocument,
} from './settingsTypes'
import {
  EMPTY_POINT_POLICY, DEFAULT_UI_ENGINE_SETTINGS, DEFAULT_POLICY_DOCS,
  DEFAULT_FEATURE, PROVIDER_LABELS,
} from './settingsConstants'

const API_BASE = ''

export function useSettingsData() {
  const [saving, setSaving] = useState(false)
  const [toast, setToast] = useState('')
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null)
  const [settingsLoadStatus, setSettingsLoadStatus] = useState<Record<string, string>>({})

  const [llmSettings, setLlmSettings] = useState<LLMSetting[]>([])
  const [apiKeys, setApiKeys] = useState<Record<string,string>>({})
  const [apiKey2s, setApiKey2s] = useState<Record<string,string>>({})
  const [endpoints, setEndpoints] = useState<Record<string,string>>({})
  const [apiKeyStatuses, setApiKeyStatuses] = useState<Record<string,string>>({})
  const [apiKey2Statuses, setApiKey2Statuses] = useState<Record<string,string>>({})
  const [apiKeyUpdatedAts, setApiKeyUpdatedAts] = useState<Record<string,string>>({})

  const [policy, setPolicy] = useState<PointPolicy | null>(null)
  const [products, setProducts] = useState<Product[]>([])
  const [adSlots, setAdSlots] = useState<AdSlot[]>([])
  const [affiliate, setAffiliate] = useState<AffiliateSettings>({ coupang:'', naver:'', class101:'', kyobo:'', adpick:'' })
  const [uiEngineSettings, setUiEngineSettings] = useState<UIEngineSettings | null>(null)
  const [policyDocs, setPolicyDocs] = useState<PolicyDocument[]>(DEFAULT_POLICY_DOCS)
  const [policyDocsLoaded, setPolicyDocsLoaded] = useState(false)

  const defaultLLMSettings = (): LLMSetting[] => ([
    { feature:DEFAULT_FEATURE, provider:'openai', model:'gpt-4o-mini' },
    { feature:'pro_curriculum', provider:'anthropic', model:'claude-3-5-sonnet' },
    { feature:'pro_tutor', provider:'grok', model:'grok-2' },
    { feature:'pro_vision', provider:'openai', model:'gpt-4o' },
    { feature:'pro_quiz', provider:'openai', model:'gpt-4o-mini' },
  ])

  const mergeLLMSettings = (incoming: LLMSetting[]) => {
    const base = defaultLLMSettings()
    return base.map(def => incoming.find(item => item.feature === def.feature) || def)
  }

  const showToast = (msg: string) => {
    setToast(msg)
    setTimeout(() => setToast(''), 3000)
  }

  const apiCall = async (path: string, method = 'GET', body?: object) => {
    try {
      const res = await fetch(`${API_BASE}${path}`, {
        method,
        headers: {
          ...(body ? { 'Content-Type': 'application/json' } : {}),
        },
        body: body ? JSON.stringify(body) : undefined,
        credentials: 'include',
      })
      if (res.status === 401 || res.status === 403) {
        setIsAuthorized(false)
        setSettingsLoadStatus(prev => ({ ...prev, [path]: `인증 실패 (${res.status})` }))
        return null
      }
      if (!res.ok) {
        setSettingsLoadStatus(prev => ({ ...prev, [path]: `요청 실패 (${res.status})` }))
        return null
      }
      const data = await res.json()
      setSettingsLoadStatus(prev => ({ ...prev, [path]: `정상 (${res.status})` }))
      return data
    } catch {
      setSettingsLoadStatus(prev => ({ ...prev, [path]: '네트워크/파싱 오류' }))
      return null
    }
  }

  useEffect(() => {
    apiCall('/api/v1/super-admin/settings/llm').then(d => {
      if (d?.settings) setLlmSettings(mergeLLMSettings(d.settings))
      else setLlmSettings(defaultLLMSettings())
    })
    apiCall('/api/v1/super-admin/settings/api-keys').then(d => {
      if (d?.keys) {
        const statuses: Record<string,string> = {}
        const statuses2: Record<string,string> = {}
        const eps: Record<string,string> = {}
        const updatedAts: Record<string,string> = {}
        d.keys.forEach((k: APIKeyStatus) => {
          statuses[k.provider] = k.key_status
          statuses2[k.provider] = k.key2_status
          if (k.endpoint_url) eps[k.provider] = k.endpoint_url
          if (k.updated_at) updatedAts[k.provider] = k.updated_at
        })
        setApiKeyStatuses(statuses)
        setApiKey2Statuses(statuses2)
        setEndpoints(eps)
        setApiKeyUpdatedAts(updatedAts)
      }
    })
    apiCall('/api/v1/super-admin/settings/point-policy').then(d => {
      if (d?.policy) setPolicy(d.policy)
      else setPolicy(EMPTY_POINT_POLICY)
    })
    apiCall('/api/v1/super-admin/settings/products').then(d => {
      if (d?.products) setProducts(d.products)
    })
    apiCall('/api/v1/super-admin/settings/ad-slots').then(d => {
      if (d?.slots) setAdSlots(d.slots)
    })
    apiCall('/api/v1/super-admin/settings/affiliate').then(d => {
      if (d?.settings) setAffiliate(prev => ({ ...prev, ...d.settings }))
    })
    apiCall('/api/v1/super-admin/settings/ui-engine').then(d => {
      if (d?.settings) setUiEngineSettings({ ...DEFAULT_UI_ENGINE_SETTINGS, ...d.settings })
      else setUiEngineSettings(DEFAULT_UI_ENGINE_SETTINGS)
    })
    apiCall('/api/v1/super-admin/settings/policies').then(d => {
      if (d?.documents?.length) setPolicyDocs(d.documents)
      setPolicyDocsLoaded(true)
    })
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const saveLLM = async () => {
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/llm', 'PUT', { settings: llmSettings })
    showToast('LLM 설정이 저장되었습니다')
    setSaving(false)
  }

  const saveAPIKey = async (provider: string) => {
    if (!apiKeys[provider] && !apiKey2s[provider] && !endpoints[provider]) return
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/api-keys', 'PUT', {
      provider,
      api_key: apiKeys[provider] || '',
      api_key2: apiKey2s[provider] || '',
      endpoint_url: endpoints[provider] || null,
    })
    showToast(`${PROVIDER_LABELS[provider]} API 키가 저장되었습니다`)
    setApiKeyStatuses(prev => ({ ...prev, [provider]: apiKeys[provider] ? 'SET' : prev[provider] }))
    setApiKey2Statuses(prev => ({ ...prev, [provider]: apiKey2s[provider] ? 'SET' : prev[provider] }))
    setApiKeyUpdatedAts(prev => ({ ...prev, [provider]: new Date().toLocaleString('ko-KR') }))
    setApiKeys(prev => ({ ...prev, [provider]: '' }))
    setApiKey2s(prev => ({ ...prev, [provider]: '' }))
    setSaving(false)
  }

  const savePolicy = async () => {
    if (!policy) return
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/point-policy', 'PUT', policy)
    showToast('포인트 정책이 저장되었습니다')
    setSaving(false)
  }

  const saveProducts = async () => {
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/products', 'PUT', { products })
    showToast('상품 목록이 저장되었습니다')
    setSaving(false)
  }

  const saveAdSlots = async () => {
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/ad-slots', 'PUT', { slots: adSlots })
    showToast('광고 슬롯이 저장되었습니다')
    setSaving(false)
  }

  const saveAffiliate = async () => {
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/affiliate', 'PUT', affiliate)
    showToast('제휴 설정이 저장되었습니다')
    setSaving(false)
  }

  const saveUIEngineSettings = async () => {
    if (!uiEngineSettings) return
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/ui-engine', 'PUT', { settings: uiEngineSettings })
    showToast('UI 엔진 운영값이 저장되었습니다')
    setSaving(false)
  }

  const savePolicyDrafts = async () => {
    setSaving(true)
    const documents = (['terms', 'privacy'] as const).flatMap((policyType) => (['ko', 'en'] as const)
      .map((locale) => policyDocs.find(item => item.type === policyType && item.locale === locale && !!item.draft) || policyDocs.find(item => item.type === policyType && item.locale === locale && item.active))
    )
      .filter(Boolean)
      .map((doc) => ({
        type: doc!.type,
        locale: doc!.locale,
        title: doc!.title,
        content: doc!.content,
        required: doc!.required,
        effective_at: doc!.effective_at,
      }))
    await apiCall('/api/v1/super-admin/settings/policies', 'PUT', { documents })
    showToast('정책 초안이 저장되었습니다')
    const refreshed = await apiCall('/api/v1/super-admin/settings/policies')
    if (refreshed?.documents) setPolicyDocs(refreshed.documents)
    setSaving(false)
  }

  const confirmPolicyDraft = async (policyType: 'terms' | 'privacy', locale: 'ko' | 'en') => {
    setSaving(true)
    await apiCall('/api/v1/super-admin/settings/policies/confirm', 'POST', { types: [policyType], locale })
    showToast('정책 문서가 확정되었습니다')
    const refreshed = await apiCall('/api/v1/super-admin/settings/policies')
    if (refreshed?.documents) setPolicyDocs(refreshed.documents)
    setSaving(false)
  }

  const handleLogout = async () => {
    await fetch('/api/v1/auth/logout', { method: 'POST', credentials: 'include' }).catch(() => undefined)
    window.location.href = '/super-admin/login'
  }

  const updateLLM = (feature: string, field: string, value: string) => {
    setLlmSettings(prev => prev.map(s => s.feature === feature ? { ...s, [field]: value } : s))
  }

  const updatePolicyDraftField = (
    policyType: 'terms' | 'privacy',
    locale: 'ko' | 'en',
    field: keyof Pick<PolicyDocument, 'title' | 'effective_at' | 'required' | 'content'>,
    value: string | boolean,
  ) => {
    setPolicyDocs(prev => {
      const draftIndex = prev.findIndex(item => item.type === policyType && item.locale === locale && !!item.draft)
      if (draftIndex >= 0) {
        return prev.map((item, index) => index === draftIndex ? { ...item, [field]: value } : item)
      }
      const activeDoc = prev.find(item => item.type === policyType && item.locale === locale && item.active)
      const koActiveDoc = prev.find(item => item.type === policyType && item.locale === 'ko' && item.active)
      const baseDoc = activeDoc ?? {
        type: policyType,
        locale,
        title: locale === 'en'
          ? (policyType === 'terms' ? 'Terms of Service' : 'Privacy Policy')
          : (policyType === 'terms' ? '서비스 이용약관' : '개인정보처리방침'),
        version: locale === 'en' ? (koActiveDoc?.version || 0) : 0,
        content: '',
        required: true,
        active: false,
        effective_at: new Date().toISOString(),
        translation_status: locale === 'en' ? 'translated' : 'source',
      }
      const nextVersion = locale === 'en'
        ? (baseDoc.version || koActiveDoc?.version || 1)
        : (baseDoc.version || 0) + 1
      return [
        ...prev,
        {
          ...baseDoc,
          id: undefined,
          version: nextVersion,
          active: false,
          draft: true,
          published_at: '',
          updated_at: '',
          [field]: value,
        },
      ]
    })
  }

  const keyStatusText = (provider: string) => {
    const parts: string[] = []
    if (apiKeyStatuses[provider] === 'SET') parts.push('기본 키 저장됨')
    if (apiKey2Statuses[provider] === 'SET') parts.push('보조 키 저장됨')
    if (endpoints[provider]) parts.push('엔드포인트 저장됨')
    if (parts.length === 0) return '미설정'
    return parts.join(' · ')
  }

  return {
    saving, toast, isAuthorized, settingsLoadStatus,
    llmSettings, apiKeys, setApiKeys, apiKey2s, setApiKey2s,
    endpoints, setEndpoints, apiKeyStatuses, apiKey2Statuses, apiKeyUpdatedAts,
    policy, setPolicy,
    products, setProducts,
    adSlots, setAdSlots,
    affiliate, setAffiliate,
    uiEngineSettings, setUiEngineSettings,
    policyDocs, policyDocsLoaded,
    saveLLM, saveAPIKey, savePolicy, saveProducts, saveAdSlots,
    saveAffiliate, saveUIEngineSettings, savePolicyDrafts, confirmPolicyDraft,
    handleLogout, updateLLM, updatePolicyDraftField, keyStatusText,
  }
}
