import { describe, it, expect } from 'vitest';
import { getPlanetLayout, planetSlotPixelSize } from '@/components/dashboard/PlanetMapLayout';

describe('getPlanetLayout — 궤도 기반 행성 배치', () => {
  it('행성 1개: 11시 방향에 배치되고 lg 크기', () => {
    const result = getPlanetLayout({ planetIds: ['p1'], breakpoint: 'desktop' });
    expect(result).toHaveLength(1);
    expect(result[0]!.size).toBe('lg');
    // 시작 각도 -120°: x=34, y≈26.6
    expect(result[0]!.position.x).toBeCloseTo(34, 1);
    expect(result[0]!.position.y).toBeLessThan(50); // 위쪽
  });

  it('행성 2개: 균등 각도(180° 간격) 배치', () => {
    const result = getPlanetLayout({ planetIds: ['p1', 'p2'], breakpoint: 'desktop' });
    expect(result).toHaveLength(2);
    // 두 행성의 x 좌표가 대칭이어야 함 (중심 기준)
    expect(result[0]!.position.x).toBeCloseTo(34, 1);
    expect(result[1]!.position.x).toBeCloseTo(66, 1);
    // y 좌표는 반대
    expect(result[0]!.position.y).toBeLessThan(50);
    expect(result[1]!.position.y).toBeGreaterThan(50);
  });

  it('행성 4개: 90° 간격 배치 (4방향)', () => {
    const result = getPlanetLayout({ planetIds: ['p1', 'p2', 'p3', 'p4'], breakpoint: 'desktop' });
    expect(result).toHaveLength(4);
    result.forEach((r) => {
      expect(r.size).toBe('sm');
      expect(r.position.x).toBeGreaterThanOrEqual(0);
      expect(r.position.x).toBeLessThanOrEqual(100);
      expect(r.position.y).toBeGreaterThanOrEqual(0);
      expect(r.position.y).toBeLessThanOrEqual(100);
    });
  });

  it('행성 5개: 최대 5개까지만 처리', () => {
    const ids = ['p1', 'p2', 'p3', 'p4', 'p5', 'p6'];
    const result = getPlanetLayout({ planetIds: ids, breakpoint: 'desktop' });
    expect(result).toHaveLength(5);
  });

  it('mobile breakpoint는 세로 스택 위치를 사용한다', () => {
    const desktop = getPlanetLayout({ planetIds: ['p1'], breakpoint: 'desktop' });
    const mobile = getPlanetLayout({ planetIds: ['p1'], breakpoint: 'mobile' });
    expect(mobile[0]!.position).toEqual({ x: 50, y: 8 });
    expect(mobile[0]!.position.y).toBeLessThan(desktop[0]!.position.y);
  });

  it('모든 행성 위치가 0~100% 범위 내에 있다', () => {
    const ids = ['p1', 'p2', 'p3', 'p4', 'p5'];
    const result = getPlanetLayout({ planetIds: ids, breakpoint: 'desktop' });
    result.forEach((r) => {
      expect(r.position.x).toBeGreaterThanOrEqual(0);
      expect(r.position.x).toBeLessThanOrEqual(100);
      expect(r.position.y).toBeGreaterThanOrEqual(0);
      expect(r.position.y).toBeLessThanOrEqual(100);
    });
  });
});

describe('planetSlotPixelSize', () => {
  it('desktop: lg=142, md=108, sm=84', () => {
    expect(planetSlotPixelSize('lg', false)).toBe(142);
    expect(planetSlotPixelSize('md', false)).toBe(108);
    expect(planetSlotPixelSize('sm', false)).toBe(84);
  });

  it('mobile: lg=112, md=92, sm=74', () => {
    expect(planetSlotPixelSize('lg', true)).toBe(112);
    expect(planetSlotPixelSize('md', true)).toBe(92);
    expect(planetSlotPixelSize('sm', true)).toBe(74);
  });
});
