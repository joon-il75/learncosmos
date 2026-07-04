import { redirect } from 'next/navigation';

type Props = {
  params: Promise<{ id: string; pointId: string }>;
};

export default async function LearningPlanetPointPage({ params }: Props) {
  const { id, pointId } = await params;
  redirect(`/dashboard/planets/learning/${id}/points/${pointId}/observe`);
}
