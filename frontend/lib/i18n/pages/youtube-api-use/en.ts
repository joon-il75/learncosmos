import type { YoutubeApiUseCopy } from './types'

export const youtubeApiUseEn: YoutubeApiUseCopy = {
  nav: {
    home: 'Home',
    login: 'Log in',
  },
  hero: {
    eyebrow: 'YouTube Data API use',
    title: 'LearnCosmos uses YouTube search to recommend learning materials',
    description:
      'LearnCosmos helps learners build self-directed learning paths. YouTube search is used to find public educational videos that may fit a learner’s goal, lesson, or learning point.',
  },
  facts: [
    { label: 'Service', value: 'LearnCosmos' },
    { label: 'Primary use', value: 'Educational video search' },
    { label: 'Default language', value: 'Korean, with English support in progress' },
    { label: 'Operator contact', value: 'learnweavr@gmail.com' },
  ],
  sections: [
    {
      title: 'Purpose of YouTube API use',
      body: [
        'LearnCosmos uses the YouTube Data API to search for public videos that can support a learner’s course, lesson, or learning point.',
        'Search results are treated as candidate learning resources. Learners can review suggested items before using them in their learning path.',
      ],
    },
    {
      title: 'How search queries are created',
      body: [
        'Queries are generated from learner-facing context such as the course title, confirmed learning goal, lesson title, point title, and user-edited recommendation criteria.',
        'Internal IDs, account credentials, API keys, and private learner records are not used as YouTube search terms.',
      ],
    },
    {
      title: 'How results are used',
      body: [
        'Returned YouTube results are used to show possible learning materials inside LearnCosmos.',
        'The service stores resource metadata needed for recommendation review, such as title, URL, provider, and thumbnail metadata. LearnCosmos does not claim ownership of YouTube videos.',
      ],
    },
    {
      title: 'User data and privacy',
      body: [
        'LearnCosmos does not send raw API keys, private notes, uploaded files, or authentication tokens to YouTube search.',
        'Learners may report unsuitable material or replace a recommendation. Those actions are used to improve learning-resource quality and safety.',
      ],
    },
    {
      title: 'Why additional quota is needed',
      body: [
        'LearnCosmos recommends learning resources across course planning, lesson exploration, and point-level study workflows.',
        'As more alpha and beta learners create courses and refresh recommendations, the number of educational search requests can exceed the default quota even when requests are user initiated and scoped to learning use cases.',
      ],
    },
  ],
  policyFallback: {
    title: 'Policies',
    body: 'English policy pages are being prepared. Until the English versions are published, the Korean versions are the governing text.',
    terms: 'Terms in Korean',
    privacy: 'Privacy Policy in Korean',
  },
  contact: {
    title: 'Contact',
    body: 'For questions about LearnCosmos or its use of YouTube Data API, contact the operator by email.',
    email: 'learnweavr@gmail.com',
  },
}
