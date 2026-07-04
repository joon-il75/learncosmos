import { Suspense } from 'react';

import ResearchMaterialPageClient from '../../../../../../_shared/ResearchMaterialPageClient';

export default function LearningPlanetPointResearchMaterialEntryPage() {
  return (
    <Suspense fallback={null}>
      <ResearchMaterialPageClient />
    </Suspense>
  );
}
