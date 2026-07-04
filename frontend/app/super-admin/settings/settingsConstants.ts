import type { PointPolicy, UIEngineSettings, PolicyDocument } from './settingsTypes'

export const EMPTY_POINT_POLICY: PointPolicy = {
  welcome_points: 0,
  course_gen_cost: 0,
  lesson_rec_cost: 0,
  admin_max_grant: 0,
  pro_monthly_points: 0,
  tutor_cost: 0,
  vision_cost: 0,
  quiz_cost: 0,
}

export const DEFAULT_UI_ENGINE_SETTINGS: UIEngineSettings = {
  compact_breakpoint: 1180,
  phone_breakpoint: 680,
  short_viewport_height: 460,
  cta_min_launch_duration_ms: 920,
  selection_cta_launch_delay_ms: 720,
  lumi_quick_action_limit: 2,
  lumi_desktop_panel_enabled: true,
  lumi_mobile_sheet_enabled: true,
}

export const PROVIDERS = ['openai','anthropic','google','grok','solar','hyperclova','llama','exaone']

export const PROVIDER_LABELS: Record<string,string> = {
  openai:'OpenAI', anthropic:'Anthropic', google:'Google Gemini',
  grok:'xAI Grok',
  solar:'Solar (업스테이지)', hyperclova:'HyperCLOVA X',
  llama:'Llama (Ollama)', exaone:'EXAONE (Ollama)',
}

export const MODELS: Record<string,string[]> = {
  openai:['gpt-4o-mini','gpt-4o','gpt-4-turbo'],
  anthropic:['claude-3-5-sonnet','claude-haiku-4-5','claude-sonnet-4-6','claude-opus-4-6'],
  google:['gemini-2.0-flash','gemini-1.5-pro'],
  grok:['grok-beta','grok-2'],
  solar:['solar-pro2'],
  hyperclova:['HyperCLOVA X'],
  llama:['llama3.2','llama3.1'],
  exaone:['exaone-3.5-7.8b'],
}

export const DEFAULT_POLICY_DOCS: PolicyDocument[] = [
  { type:'terms', locale:'ko', title:'서비스 이용약관', version:1, content:'', required:true, active:true, effective_at:new Date().toISOString() },
  { type:'terms', locale:'en', title:'Terms of Service', version:1, content:'', required:true, active:false, effective_at:new Date().toISOString(), translation_status:'translated' },
  { type:'privacy', locale:'ko', title:'개인정보처리방침', version:1, content:'', required:true, active:true, effective_at:new Date().toISOString() },
  { type:'privacy', locale:'en', title:'Privacy Policy', version:1, content:'', required:true, active:false, effective_at:new Date().toISOString(), translation_status:'translated' },
]

export const DEFAULT_FEATURE = 'default'

export const PRO_FEATURES = ['pro_curriculum','pro_tutor','pro_vision','pro_quiz']

export const FEATURE_LABELS: Record<string,string> = {
  default:'기본값',
  pro_curriculum:'Pro 커리큘럼',
  pro_tutor:'Pro 튜터',
  pro_vision:'Pro 비전 코칭',
  pro_quiz:'Pro 퀴즈',
}

export const SLOT_TYPES = ['category_sponsor','step_complete','dashboard_banner']

export const SLOT_LABELS: Record<string,string> = {
  category_sponsor:'카테고리 스폰서',
  step_complete:'스텝완료 카드',
  dashboard_banner:'대시보드 배너',
}

export const CATEGORIES = ['드로잉·회화','악기','외국어','요리·베이킹','사진·영상','프로그래밍','공예·DIY','글쓰기','댄스','피트니스','보컬·작곡','가드닝','디자인']
