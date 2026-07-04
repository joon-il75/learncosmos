import { redirect } from 'next/navigation';

type Props = {
  params: Promise<{ id: string; pointId: string }>;
};

export default async function SharedPlanetPointPage({ params }: Props) {
  const { id, pointId } = await params;
  redirect(`/dashboard/planets/shared/${id}/points/${pointId}/observe`);
}
