type RecommendationQueryTarget = 'new' | 'selected'

export type BuildRecommendationQueryInput = {
  target?: RecommendationQueryTarget
  courseTitle?: string | null
  courseGoal?: string | null
  regionTitle?: string | null
  regionDescription?: string | null
  subRegionTitle?: string | null
  subRegionDescription?: string | null
  nodeTitleInput?: string | null
  selectedNodeTitle?: string | null
  nodeSummary?: string | null
}

const STOP_WORDS = new Set([
  '그리고',
  '또는',
  '및',
  '위한',
  '통한',
  '통해',
  '하는',
  '하고',
  '한다',
  '하며',
  '있다',
  '있는',
  '싶어',
  '싶습니다',
  '이해한다',
  '선택하고',
  '적합한',
  '적합',
  '따라',
  '효과적',
  '효과적으로',
  '체계적',
  '고품질',
  '자료',
  '만들어낼',
  '수',
  '높은',
  '기본',
  '기초',
  '입문',
  '과정',
  '준비',
  '완료',
  'the',
  'and',
  'or',
  'for',
  'with',
  'from',
  'into',
  'about',
  'learn',
  'learning',
  'basic',
  'basics',
  'beginner',
  'course',
  'lesson',
  'step',
  'steps',
])

const LOW_PRIORITY_WORDS = new Set([
  '수익',
  '수익성',
  '수익성이',
  '수익화',
  '수익화에',
  '창출',
  '창출하는',
  '매출',
  '마케팅',
  '판매',
])

const PARTICLE_SUFFIXES = [
  '하고',
  '으로',
  '에서',
  '에게',
  '까지',
  '부터',
  '보다',
  '처럼',
  '하기',
  '하며',
  '하여',
  '해서',
  '하다',
  '한다',
  '할',
  '된',
  '될',
  '한',
  '은',
  '는',
  '이',
  '가',
  '을',
  '를',
  '의',
  '에',
  '와',
  '과',
  '도',
  '로',
]

function compactParts(parts: Array<string | null | undefined>): string[] {
  return parts
    .map((part) => part?.trim() ?? '')
    .filter((part) => part !== '')
}

function normalizeWord(word: string): string {
  let normalized = word.trim().toLowerCase()
  if (normalized === '') return ''

  for (const suffix of PARTICLE_SUFFIXES) {
    if (normalized.endsWith(suffix) && [...normalized].length > [...suffix].length + 1) {
      normalized = normalized.slice(0, -suffix.length)
      break
    }
  }
  return normalized
}

function tokenize(text: string): string[] {
  return text
    .replace(/[^\p{L}\p{N}\s]/gu, ' ')
    .split(/\s+/)
    .map(normalizeWord)
    .filter((word) => [...word].length >= 2)
    .filter((word) => !STOP_WORDS.has(word))
}

function appendUnique(target: string[], seen: Set<string>, words: string[], maxCount: number) {
  for (const word of words) {
    if (target.length >= maxCount) return
    if (seen.has(word)) continue
    seen.add(word)
    target.push(word)
  }
}

function prioritizedWords(text: string, maxCount: number, options?: { allowLowPriority?: boolean }): string[] {
  const words = tokenize(text)
  const highPriority = words.filter((word) => !LOW_PRIORITY_WORDS.has(word))
  const lowPriority = options?.allowLowPriority ? words.filter((word) => LOW_PRIORITY_WORDS.has(word)) : []
  return [...highPriority, ...lowPriority].slice(0, maxCount)
}

function courseDomainWords(courseTitle: string, courseGoal: string): string[] {
  const words = prioritizedWords(compactParts([courseTitle, courseGoal]).join(' '), 5)
  if (words.includes('온라인') && words.includes('강의')) {
    return ['온라인', '강의']
  }
  return words.slice(0, 2)
}

function lessonTitleText(input: BuildRecommendationQueryInput, nodeTitle: string): string {
  return compactParts([
    input.regionTitle,
    input.subRegionTitle,
    nodeTitle,
  ]).join(' ')
}

function lessonObjectiveText(input: BuildRecommendationQueryInput): string {
  return compactParts([
    input.regionDescription,
    input.subRegionDescription,
    input.nodeSummary,
  ]).join(' ')
}

function withSearchIntent(words: string[]): string[] {
  if (words.length === 0) return words
  const hasKoreanIntent = words.some((word) => ['방법', '가이드', '비교', '추천', '기초', '입문'].includes(word))
  const hasEnglishIntent = words.some((word) => ['tutorial', 'guide', 'how', 'practice', 'project'].includes(word))
  if (hasKoreanIntent || hasEnglishIntent) return words
  const hasHangul = words.some((word) => /[가-힣]/.test(word))
  return hasHangul ? [...words, '방법'] : [...words, 'tutorial']
}

export function buildDefaultRecommendationQuery(input: BuildRecommendationQueryInput): string {
  const target = input.target ?? 'new'
  const nodeTitle = target === 'selected'
    ? compactParts([input.nodeTitleInput, input.selectedNodeTitle])[0] ?? ''
    : input.nodeTitleInput ?? ''
  const lessonTitle = lessonTitleText(input, nodeTitle)
  const lessonObjective = lessonObjectiveText(input)

  const result: string[] = []
  const seen = new Set<string>()

  appendUnique(result, seen, courseDomainWords(input.courseTitle ?? '', input.courseGoal ?? ''), 2)
  appendUnique(result, seen, prioritizedWords(lessonTitle, 8), 6)
  appendUnique(result, seen, prioritizedWords(lessonObjective, 8), 7)

  return withSearchIntent(result).slice(0, 8).join(' ')
}
