import type { DashboardSortKey } from '@/lib/world-ui-engine/types';
import type { PlanetFilter } from './_types';

export const PAGE_SIZE = 5;

export const TOPIC_CHIPS = ['기타 입문', '수채화 기초', '라떼아트', 'DSLR 사진촬영'];

export const SORT_OPTIONS: Array<{ key: DashboardSortKey; label: string }> = [
  { key: 'updated_desc', label: '최신순' },
  { key: 'updated_asc', label: '오래된순' },
  { key: 'title_asc', label: '제목 오름차순' },
  { key: 'title_desc', label: '제목 내림차순' },
];

export const planetFilters: Array<{ value: PlanetFilter; label: string }> = [
  { value: 'all', label: '전체' },
  { value: 'learning', label: '학습중' },
  { value: 'completed', label: '완료' },
];
