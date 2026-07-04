export type YoutubeApiUseCopy = {
  nav: {
    home: string
    login: string
  }
  hero: {
    eyebrow: string
    title: string
    description: string
  }
  sections: Array<{
    title: string
    body: string[]
  }>
  facts: Array<{
    label: string
    value: string
  }>
  policyFallback: {
    title: string
    body: string
    terms: string
    privacy: string
  }
  contact: {
    title: string
    body: string
    email: string
  }
}
