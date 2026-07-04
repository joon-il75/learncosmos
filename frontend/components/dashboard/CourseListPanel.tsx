'use client';

import type { CSSProperties } from 'react';
import {
  statusColor,
  type PlanetStatus,
} from '@/components/dashboard/planetMapTokens';
import type { DashboardCourseRecord, DashboardStatusFilter } from '@/lib/world-ui-engine/types';
import type { PageSystemChunk } from '@/components/dashboard/GalaxyMap';
import type { DashboardMainCopy } from '@/lib/i18n/pages/dashboardMain';
import type { Locale } from '@/lib/i18n/locales';
import { buildSystemStatusLine, formatSystemUpdatedAt } from '@/components/dashboard/GalaxyMap';

// ─── Types ────────────────────────────────────────────────────────────────────

type PlanetFilter = DashboardStatusFilter;

const STATUS_FILTERS: Array<{ value: PlanetFilter; color: string }> = [
  { value: 'all',       color: '#8EA4C8' },
  { value: 'learning',  color: '#22C55E' },
  { value: 'completed', color: '#F59E0B' },
];

// ─── Props ────────────────────────────────────────────────────────────────────

interface CourseListPanelProps {
  courses: DashboardCourseRecord[];
  mode: 'galaxy' | 'star-system';
  isCompactLayout?: boolean;
  isListOnly?: boolean;
  page: number;
  totalPages: number;
  hoveredCourseId: string | null;
  selectedCourseId: string | null;
  selectedSystemChunk: PageSystemChunk | null;
  searchQuery: string;
  activeStatusFilter: PlanetFilter;
  listModeLabel: string;
  listPageLabel: string;
  showInactiveCourses: boolean;
  inactiveCourseCount: number;
  navEnabled: boolean;
  prevPageDisabled: boolean;
  nextPageDisabled: boolean;
  copy: DashboardMainCopy['list'];
  galaxyCopy: DashboardMainCopy['galaxy'];
  locale: Locale;
  onSearchChange: (query: string) => void;
  onStatusFilterChange: (status: PlanetFilter) => void;
  onCourseHover: (courseId: string | null) => void;
  onCourseClick: (courseId: string) => void;
  onPagePrev: () => void;
  onPageNext: () => void;
  onToggleInactiveCourses: () => void;
  onModeToggle: () => void;
  onPrevSystem?: () => void;
  onNextSystem?: () => void;
  prevSystemDisabled?: boolean;
  nextSystemDisabled?: boolean;
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

function formatUpdatedAt(iso: string): string {
  return new Date(iso).toLocaleDateString('ko-KR', { month: 'short', day: 'numeric' });
}

function buildProgressMessage(course: DashboardCourseRecord, copy: DashboardMainCopy['list']['progress']): string {
  const { status, levelCount, lessonCount, plannedLevelCount, completedLessonCount, progress } = course;
  if (course.isInactive) {
    return copy.inactive;
  }
  if (status === 'learning') {
    if (lessonCount > 0) {
      return copy.learningWithProgress(Math.round((progress ?? 0) * 100), completedLessonCount, lessonCount);
    }
    return copy.learning;
  }
  if (status === 'draft') {
    if (levelCount > 0 && progress != null) {
      return copy.draftWithProgress(Math.round(progress * 100), plannedLevelCount, levelCount);
    }
    return copy.draft;
  }
  if (status === 'ready') {
    if (levelCount > 0) {
      return copy.readyWithRegions(levelCount);
    }
    if (lessonCount > 0) {
      return copy.readyWithLessons(lessonCount);
    }
    return copy.ready;
  }
  if (status === 'completed') {
    if (lessonCount > 0) {
      return copy.completedWithLessons(lessonCount);
    }
    return copy.completed;
  }
  return '';
}

function getDisplayStatusLabel(
  course: DashboardCourseRecord,
  statusLabels: DashboardMainCopy['galaxy']['statusLabels'],
): string {
  if (course.isInactive) return statusLabels.learning;
  return statusLabels[course.status as PlanetStatus];
}

// ─── Main component ───────────────────────────────────────────────────────────

export default function CourseListPanel({
  courses,
  mode,
  isCompactLayout = false,
  isListOnly = false,
  page,
  totalPages,
  hoveredCourseId,
  selectedCourseId,
  selectedSystemChunk,
  searchQuery,
  activeStatusFilter,
  listModeLabel,
  listPageLabel,
  showInactiveCourses,
  inactiveCourseCount,
  navEnabled,
  prevPageDisabled,
  nextPageDisabled,
  copy,
  galaxyCopy,
  locale,
  onSearchChange,
  onStatusFilterChange,
  onCourseHover,
  onCourseClick,
  onPagePrev,
  onPageNext,
  onToggleInactiveCourses,
  onModeToggle,
  onPrevSystem,
  onNextSystem,
  prevSystemDisabled = true,
  nextSystemDisabled = true,
}: CourseListPanelProps) {
  const compactMode = isCompactLayout || isListOnly;
  return (
    <section style={panelStyle(compactMode, isListOnly)}>
      {/* 헤더 */}
      <div style={panelHeaderStyle(compactMode)}>
        <div style={{ display: 'grid', gap: '4px' }}>
          <div style={eyebrowStyle}>{listModeLabel}</div>
          <div style={metaStyle}>
            {listPageLabel}
            {mode === 'galaxy' ? ` · ${copy.totalPages(Math.max(totalPages, 1))}` : ''}
          </div>
        </div>
        <button
          type="button"
          style={modeToggleButtonStyle(mode === 'star-system', compactMode)}
          onClick={onModeToggle}
        >
          {mode === 'star-system' ? '← Galaxy' : 'Star System'}
        </button>
      </div>

      {/* 상태 필터 */}
      <div style={filterRowStyle(compactMode)}>
          {STATUS_FILTERS.map((f) => {
            const isActive = activeStatusFilter === f.value;
            const label = f.value === 'all' ? copy.statusFilters.all : f.value === 'learning' ? copy.statusFilters.learning : copy.statusFilters.completed;
            return (
            <button
              key={f.value}
              type="button"
              aria-label={label}
              onClick={() => onStatusFilterChange(f.value)}
              style={filterButtonStyle(f.color, isActive, compactMode)}
            >
              <span style={filterDotStyle(f.color, isActive)} />
              <span style={filterLabelStyle(f.color, isActive)}>{label}</span>
            </button>
          );
        })}
        <button
          type="button"
          onClick={onToggleInactiveCourses}
          style={inactiveToggleButtonStyle(showInactiveCourses, compactMode)}
          title={showInactiveCourses ? copy.inactive.hideTitle : copy.inactive.showTitle}
        >
          {showInactiveCourses ? copy.inactive.hide : copy.inactive.show}
          {inactiveCourseCount > 0 ? ` ${inactiveCourseCount}` : ''}
        </button>
      </div>

      {/* 검색 */}
      <div style={searchRowStyle(compactMode)}>
        <input
          type="search"
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder={copy.searchPlaceholder}
          style={searchInputStyle(compactMode)}
          aria-label={copy.searchAria}
        />
      </div>

      {/* Star System 항성계 요약 */}
      {mode === 'star-system' && selectedSystemChunk ? (
        <div style={systemSummaryStyle(compactMode)}>
          {buildSystemStatusLine(selectedSystemChunk.statusCount, galaxyCopy)} ·{' '}
          {formatSystemUpdatedAt(selectedSystemChunk.updatedAt, locale, galaxyCopy)}
        </div>
      ) : null}

      {/* 페이지 네비게이션 */}
      {navEnabled ? (
        <div style={pageNavStyle(compactMode)}>
          <button
            type="button"
            onClick={onPagePrev}
            disabled={prevPageDisabled}
            style={pageNavButtonStyle(prevPageDisabled, compactMode)}
          >
            {copy.pagination.previous}
          </button>
          <div style={pageNavLabelStyle}>{`${page} / ${Math.max(totalPages, 1)}`}</div>
          <button
            type="button"
            onClick={onPageNext}
            disabled={nextPageDisabled}
            style={pageNavButtonStyle(nextPageDisabled, compactMode)}
          >
            {copy.pagination.next}
          </button>
        </div>
      ) : null}

      {/* 항성계 이동 내비 (Star System 모드 전용) */}
      {mode === 'star-system' ? (
        <div style={systemNavStyle(compactMode)}>
          <button
            type="button"
            onClick={onPrevSystem}
            disabled={prevSystemDisabled}
            style={systemNavButtonStyle(prevSystemDisabled, compactMode)}
            aria-label={copy.systemNav.previousAria}
          >
            {copy.systemNav.previous}
          </button>
          <span style={systemNavLabelStyle(compactMode)}>{copy.systemNav.label}</span>
          <button
            type="button"
            onClick={onNextSystem}
            disabled={nextSystemDisabled}
            style={systemNavButtonStyle(nextSystemDisabled, compactMode)}
            aria-label={copy.systemNav.nextAria}
          >
            {copy.systemNav.next}
          </button>
        </div>
      ) : null}

      {/* 코스 목록 */}
      {courses.length > 0 ? (
        <div style={listBodyStyle(compactMode)}>
          {courses.map((course) => {
            const isHovered  = hoveredCourseId  === course.id;
            const isSelected = selectedCourseId === course.id;
            const status = course.status as PlanetStatus;
            return (
              <button
                key={course.id}
                type="button"
                onClick={() => onCourseClick(course.id)}
                onMouseEnter={() => onCourseHover(course.id)}
                onMouseLeave={() => onCourseHover(null)}
                style={courseItemStyle(isHovered, isSelected, compactMode, Boolean(course.isInactive))}
              >
                <div style={courseTitleStyle(compactMode)}>{course.title}</div>
                <div style={courseStatusRowStyle}>
                  <span style={statusDotStyle(status)} />
                  <span style={courseStatusLabelStyle}>{getDisplayStatusLabel(course, galaxyCopy.statusLabels)}</span>
                  {course.isInactive ? <span style={inactiveBadgeStyle}>{copy.inactive.badge}</span> : null}
                </div>
                <div style={compactProgressStyle(status)}>{buildProgressMessage(course, copy.progress)}</div>
              </button>
            );
          })}
        </div>
      ) : (
        <div style={emptyStyle(compactMode)}>
          {searchQuery
            ? copy.emptySearch(searchQuery)
            : copy.emptyFilter}
        </div>
      )}
    </section>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

function panelStyle(compactMode: boolean, isListOnly: boolean): CSSProperties {
  return {
    display: 'flex',
    flexDirection: 'column',
    gap: compactMode ? '12px' : '10px',
    padding: compactMode ? '16px' : '14px',
    flex: 1,
    minHeight: 0,
    boxSizing: 'border-box',
    borderRadius: compactMode && isListOnly ? '24px' : '20px',
    border: '1px solid rgba(214, 226, 255, 0.16)',
    background: compactMode
      ? 'rgba(18, 32, 62, 0.78)'
      : 'rgba(18, 32, 62, 0.68)',
    backdropFilter: 'blur(14px)',
  };
}

function panelHeaderStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'flex',
    alignItems: compactMode ? 'stretch' : 'flex-start',
    justifyContent: 'space-between',
    gap: compactMode ? '10px' : '8px',
    flexWrap: 'wrap',
  };
}

const eyebrowStyle = {
  fontSize: '11px',
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
  color: '#8CC7FF',
  fontWeight: 700,
} satisfies CSSProperties;

const metaStyle = {
  fontSize: '12px',
  color: '#8EA4C8',
} satisfies CSSProperties;

function modeToggleButtonStyle(isActive: boolean, compactMode: boolean): CSSProperties {
  return {
    minHeight: compactMode ? '40px' : '34px',
    padding: compactMode ? '0 16px' : '0 14px',
    borderRadius: '14px',
    border: isActive
      ? '1px solid rgba(248, 214, 70, 0.8)'
      : '1px solid rgba(255, 255, 255, 0.2)',
    background: isActive ? 'rgba(248, 214, 70, 0.18)' : 'rgba(255, 255, 255, 0.06)',
    color: isActive ? '#FDE68A' : '#E8EEFF',
    fontSize: compactMode ? '13px' : '12px',
    fontWeight: 600,
    cursor: 'pointer',
    fontFamily: 'inherit',
    whiteSpace: 'nowrap',
  };
}

// 필터 행
function filterRowStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'flex',
    flexWrap: compactMode ? 'nowrap' : 'wrap',
    gap: compactMode ? '6px' : '5px',
    alignItems: 'center',
    overflowX: compactMode ? 'auto' : 'visible',
    paddingBottom: compactMode ? '2px' : 0,
  };
}

function filterButtonStyle(color: string, isActive: boolean, compactMode: boolean): CSSProperties {
  return {
    display: 'flex',
    alignItems: 'center',
    gap: compactMode ? '6px' : '5px',
    padding: compactMode ? '7px 10px 7px 8px' : '3px 8px 3px 6px',
    borderRadius: '999px',
    border: isActive ? `1px solid ${color}88` : '1px solid rgba(255,255,255,0.1)',
    background: isActive ? `${color}22` : 'rgba(255,255,255,0.04)',
    cursor: 'pointer',
    fontFamily: 'inherit',
    opacity: isActive ? 1 : 0.55,
    transition: 'opacity 140ms ease, background 140ms ease, border-color 140ms ease',
  };
}

function filterDotStyle(color: string, isActive: boolean): CSSProperties {
  return {
    width: '7px',
    height: '7px',
    borderRadius: '999px',
    background: color,
    flexShrink: 0,
    boxShadow: isActive ? `0 0 6px ${color}99` : 'none',
  };
}

function filterLabelStyle(color: string, isActive: boolean): CSSProperties {
  return {
    fontSize: '11px',
    fontWeight: isActive ? 700 : 500,
    color: isActive ? color : '#8EA4C8',
    whiteSpace: 'nowrap',
    lineHeight: 1,
  };
}

function inactiveToggleButtonStyle(isActive: boolean, compactMode: boolean): CSSProperties {
  return {
    minHeight: compactMode ? '32px' : '26px',
    padding: compactMode ? '0 12px' : '0 10px',
    borderRadius: '999px',
    border: isActive ? '1px solid rgba(203, 213, 225, 0.62)' : '1px solid rgba(255,255,255,0.12)',
    background: isActive ? 'rgba(148, 163, 184, 0.18)' : 'rgba(255,255,255,0.04)',
    color: isActive ? '#E2E8F0' : '#8EA4C8',
    fontSize: compactMode ? '12px' : '11px',
    fontWeight: 800,
    cursor: 'pointer',
    fontFamily: 'inherit',
    whiteSpace: 'nowrap',
  };
}

function searchRowStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'grid',
    position: compactMode ? 'sticky' : 'static',
    top: compactMode ? 0 : 'auto',
    zIndex: compactMode ? 1 : 'auto',
  };
}

function searchInputStyle(compactMode: boolean): CSSProperties {
  return {
  height: compactMode ? '42px' : '36px',
  padding: compactMode ? '0 14px' : '0 12px',
  borderRadius: '12px',
  border: '1px solid rgba(180, 205, 255, 0.2)',
  background: 'rgba(10, 20, 40, 0.6)',
  color: '#E8F0FF',
  fontSize: compactMode ? '14px' : '13px',
  fontFamily: 'inherit',
  outline: 'none',
  width: '100%',
  boxSizing: 'border-box',
  };
}

function systemSummaryStyle(compactMode: boolean): CSSProperties {
  return {
    fontSize: compactMode ? '13px' : '12px',
    color: '#E0E7FF',
    lineHeight: 1.5,
    padding: compactMode ? '6px 0' : '4px 0',
  };
}

function pageNavStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'grid',
    gridTemplateColumns: compactMode ? '60px 1fr 60px' : '36px 1fr 36px',
    gap: compactMode ? '8px' : '6px',
    alignItems: 'center',
  };
}

const pageNavLabelStyle = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: '#FDE68A',
  fontWeight: 700,
  fontSize: '12px',
} satisfies CSSProperties;

function pageNavButtonStyle(isDisabled: boolean, compactMode: boolean): CSSProperties {
  return {
    minHeight: compactMode ? '38px' : '32px',
    borderRadius: '999px',
    border: '1px solid rgba(255, 255, 255, 0.2)',
    background: isDisabled ? 'rgba(255, 255, 255, 0.08)' : 'rgba(255, 255, 255, 0.12)',
    color: '#F4F8FF',
    fontSize: compactMode ? '13px' : '12px',
    fontWeight: 600,
    cursor: isDisabled ? 'not-allowed' : 'pointer',
    fontFamily: 'inherit',
  };
}

function listBodyStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'grid',
    alignContent: 'start',
    gap: compactMode ? '10px' : '6px',
    flex: 1,
    minHeight: 0,
    overflowY: 'auto',
    paddingRight: '4px',
  };
}

// 코스 카드
function courseItemStyle(isHovered: boolean, isSelected: boolean, compactMode: boolean, isInactive: boolean): CSSProperties {
  return {
    display: 'grid',
    width: '100%',
    gap: compactMode ? '6px' : '4px',
    padding: compactMode ? '14px 14px' : '10px 12px',
    borderRadius: compactMode ? '16px' : '12px',
    boxSizing: 'border-box',
    background: isInactive
      ? 'rgba(51, 65, 85, 0.48)'
      : isSelected
      ? 'rgba(251, 191, 36, 0.12)'
      : isHovered
        ? 'rgba(255, 255, 255, 0.10)'
        : 'rgba(255, 255, 255, 0.07)',
    border: isInactive
      ? '1px solid rgba(148, 163, 184, 0.18)'
      : isSelected
      ? '1px solid rgba(251, 191, 36, 0.5)'
      : isHovered
        ? '1px solid rgba(255, 255, 255, 0.2)'
        : '1px solid rgba(255, 255, 255, 0.07)',
    cursor: 'pointer',
    textAlign: 'left',
    fontFamily: 'inherit',
    transition: 'background 120ms ease, border-color 120ms ease',
    boxShadow: isSelected ? '0 0 16px rgba(251, 191, 36, 0.2)' : 'none',
    opacity: isInactive ? 0.74 : 1,
  };
}

// 제목: 1줄 말줄임
function courseTitleStyle(compactMode: boolean): CSSProperties {
  return {
    fontSize: compactMode ? '14px' : '13px',
    fontWeight: 600,
    color: '#F4F8FF',
    lineHeight: 1.45,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: compactMode ? 'normal' : 'nowrap',
    display: compactMode ? '-webkit-box' : 'block',
    WebkitLineClamp: compactMode ? 2 : undefined,
    WebkitBoxOrient: compactMode ? 'vertical' : undefined,
  };
}

// 상태 행
const courseStatusRowStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '5px',
} satisfies CSSProperties;

function statusDotStyle(status: PlanetStatus): CSSProperties {
  return {
    width: '7px',
    height: '7px',
    borderRadius: '999px',
    background: statusColor[status],
    flexShrink: 0,
  };
}

const courseStatusLabelStyle = {
  fontSize: '11px',
  color: '#8EA4C8',
  lineHeight: 1,
} satisfies CSSProperties;

const inactiveBadgeStyle = {
  marginLeft: 4,
  padding: '2px 6px',
  borderRadius: 999,
  border: '1px solid rgba(148, 163, 184, 0.28)',
  background: 'rgba(148, 163, 184, 0.12)',
  color: '#CBD5E1',
  fontSize: '10px',
  fontWeight: 800,
  lineHeight: 1,
} satisfies CSSProperties;

function emptyStyle(compactMode: boolean): CSSProperties {
  return {
  padding: compactMode ? '22px' : '18px',
  borderRadius: '18px',
  border: '1px dashed rgba(194, 210, 245, 0.18)',
  background: 'rgba(10, 20, 36, 0.42)',
  color: '#C8D1E8',
  fontSize: compactMode ? '15px' : '14px',
  lineHeight: 1.6,
  };
}

function systemNavStyle(compactMode: boolean): CSSProperties {
  return {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    gap: compactMode ? '10px' : '8px',
    padding: compactMode ? '8px 0' : '6px 2px',
    borderTop: '1px solid rgba(180, 205, 255, 0.1)',
  };
}

function systemNavButtonStyle(disabled?: boolean, compactMode = false): CSSProperties {
  return {
    minHeight: compactMode ? '36px' : '30px',
    padding: compactMode ? '0 12px' : '0 10px',
    borderRadius: '10px',
    border: '1px solid rgba(194, 210, 245, 0.18)',
    background: disabled ? 'transparent' : 'rgba(255, 255, 255, 0.06)',
    color: disabled ? 'rgba(180, 200, 240, 0.3)' : 'rgba(200, 218, 255, 0.82)',
    fontSize: compactMode ? '13px' : '12px',
    fontFamily: 'inherit',
    cursor: disabled ? 'default' : 'pointer',
  };
}

function systemNavLabelStyle(compactMode: boolean): CSSProperties {
  return {
    fontSize: compactMode ? '12px' : '11px',
    color: 'rgba(160, 185, 230, 0.55)',
    letterSpacing: '0.03em',
  };
}

function compactProgressStyle(status: PlanetStatus): CSSProperties {
  return {
    fontSize: '12px',
    lineHeight: 1.45,
    color:
      status === 'learning'
        ? '#7BE0A5'
        : status === 'ready'
          ? '#8BC2FF'
          : status === 'completed'
            ? '#FFD27A'
            : '#B9C3D8',
  };
}
