export type LLMSetting = { feature: string; provider: string; model: string; endpoint_url?: string }

export type APIKeyStatus = {
  provider: string
  key_status: string
  key2_status: string
  endpoint_url?: string
  has_key?: boolean
  has_key2?: boolean
  updated_at?: string
}

export type PointPolicy = {
  welcome_points: number
  course_gen_cost: number
  lesson_rec_cost: number
  admin_max_grant: number
  pro_monthly_points: number
  tutor_cost: number
  vision_cost: number
  quiz_cost: number
}

export type Product = {
  id?: string
  name: string
  points: number
  price_krw: number
  is_active: boolean
  sort_order: number
  product_type: 'one_time' | 'subscription'
  billing_period?: 'monthly' | 'yearly'
  monthly_points: number
  ad_free: boolean
  premium_access: boolean
  description?: string
}

export type AdSlot = {
  id?: string
  slot_type: string
  title: string
  link_url: string
  category_slug?: string
  is_active: boolean
  sort_order: number
}

export type AffiliateSettings = {
  coupang: string
  naver: string
  class101: string
  kyobo: string
  adpick: string
}

export type UIEngineSettings = {
  compact_breakpoint: number
  phone_breakpoint: number
  short_viewport_height: number
  cta_min_launch_duration_ms: number
  selection_cta_launch_delay_ms: number
  lumi_quick_action_limit: number
  lumi_desktop_panel_enabled: boolean
  lumi_mobile_sheet_enabled: boolean
}

export type PolicyDocument = {
  id?: string
  set_id?: string
  type: 'terms' | 'privacy'
  locale: 'ko' | 'en'
  title: string
  version: number
  content: string
  required: boolean
  active: boolean
  draft?: boolean
  translation_status?: 'source' | 'translated' | 'fallback'
  effective_at: string
  published_at?: string
  updated_at?: string
}
