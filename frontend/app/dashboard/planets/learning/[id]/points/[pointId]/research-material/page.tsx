import { Suspense } from 'react';

import ResearchMaterialPageClient from '../../../../../_shared/ResearchMaterialPageClient';

export default function LearningPlanetPointResearchMaterialPage() {
  return (
    <Suspense fallback={null}>
      <ResearchMaterialPageClient />
    </Suspense>
  );
}
