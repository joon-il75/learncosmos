'use client';

import { useMemo, useState, type CSSProperties, type FormEvent } from 'react';
import type { CourseItem, TodayTask } from './_types';

type Props = {
  courses: CourseItem[];
  todayTask: TodayTask | null;
  userName: string;
  isSubmitting: boolean;
  creationError: string | null;
  onCreateStar: (query: string) => Promise<void>;
  onOpenCourse: (course: CourseItem) => void;
  onOpenTask: (href: string) => void;
};

const starSlots = [
  { left: 50, top: 42, size: 34 },
  { left: 24, top: 28, size: 18 },
  { left: 72, top: 24, size: 16 },
  { left: 30, top: 62, size: 14 },
  { left: 77, top: 62, size: 18 },
  { left: 47, top: 72, size: 12 },
  { left: 14, top: 48, size: 10 },
  { left: 88, top: 44, size: 10 },
];

const statusTone: Record<CourseItem['status'], { core: string; glow: string; ring: string; label: string }> = {
  draft: { core: '#B6C7D8', glow: 'rgba(182, 199, 216, 0.34)', ring: 'rgba(182, 199, 216, 0.48)', label: '초안' },
  ready: { core: '#75D7F0', glow: 'rgba(117, 215, 240, 0.38)', ring: 'rgba(117, 215, 240, 0.52)', label: '준비' },
  learning: { core: '#8B7CFC', glow: 'rgba(139, 124, 252, 0.42)', ring: 'rgba(139, 124, 252, 0.60)', label: '탐험중' },
  completed: { core: '#F8C85E', glow: 'rgba(248, 200, 94, 0.45)', ring: 'rgba(248, 200, 94, 0.68)', label: '완료' },
};

function getRecommendedCourse(courses: CourseItem[]) {
  const active = courses
    .filter((course) => !course.isInactive && course.status === 'learning')
    .sort((a, b) => new Date(b.lastAccessedAt ?? b.updatedAt).getTime() - new Date(a.lastAccessedAt ?? a.updatedAt).getTime());
  const ready = courses
    .filter((course) => !course.isInactive && course.status === 'ready')
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
  const draft = courses
    .filter((course) => !course.isInactive && course.status === 'draft')
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
  const completed = courses
    .filter((course) => !course.isInactive && course.status === 'completed')
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());

  return active[0] ?? ready[0] ?? draft[0] ?? completed[0] ?? courses[0] ?? null;
}

function getProgress(course: CourseItem | null) {
  if (!course) return 0;
  if (typeof course.progress === 'number') return Math.round(course.progress);
  if (course.lessonCount > 0) return Math.round((course.completedLessonCount / course.lessonCount) * 100);
  return course.status === 'completed' ? 100 : 0;
}

function clampTitle(title: string) {
  return title.length > 24 ? `${title.slice(0, 23)}...` : title;
}

const rootStyle: CSSProperties = {
  minHeight: 'calc(100dvh - 60px)',
  padding: '10px 10px 14px',
  background:
    'linear-gradient(180deg, #FBFFFE 0%, #F4FBFF 38%, #F8F4FF 100%)',
  color: '#12384E',
};

const shellStyle: CSSProperties = {
  width: '100%',
  maxWidth: '360px',
  minHeight: 'calc(100dvh - 84px)',
  margin: '0 auto',
  display: 'grid',
  gridTemplateRows: 'auto minmax(248px, 1fr) auto auto',
  gap: '10px',
};

const routeStyle: CSSProperties = {
  border: '1px solid rgba(151, 190, 210, 0.34)',
  borderRadius: '18px',
  padding: '10px 12px',
  background: 'rgba(255, 255, 255, 0.78)',
  boxShadow: '0 12px 36px rgba(35, 79, 112, 0.08)',
};

const routeEyebrowStyle: CSSProperties = {
  margin: 0,
  color: '#6C7ABF',
  fontSize: '10px',
  lineHeight: 1.2,
  fontWeight: 900,
  letterSpacing: '0.12em',
  textTransform: 'uppercase',
};

const routeTextStyle: CSSProperties = {
  margin: '5px 0 0',
  color: '#12384E',
  fontSize: '13px',
  lineHeight: 1.36,
  fontWeight: 850,
};

const mapStyle: CSSProperties = {
  position: 'relative',
  overflow: 'hidden',
  minHeight: '252px',
  borderRadius: '24px',
  border: '1px solid rgba(174, 198, 230, 0.40)',
  background:
    'radial-gradient(ellipse at 46% 38%, rgba(198, 184, 255, 0.50), transparent 38%), radial-gradient(ellipse at 68% 22%, rgba(128, 220, 244, 0.24), transparent 32%), linear-gradient(160deg, rgba(255,255,255,0.94), rgba(243,247,255,0.90) 54%, rgba(250,245,255,0.94))',
  boxShadow: 'inset 0 1px 0 rgba(255,255,255,0.88), 0 18px 48px rgba(83, 96, 143, 0.14)',
};

const galaxyVeilStyle: CSSProperties = {
  position: 'absolute',
  inset: '10% -18%',
  transform: 'rotate(-14deg)',
  background:
    'linear-gradient(90deg, transparent, rgba(163, 139, 255, 0.18), rgba(124, 211, 252, 0.15), transparent)',
  filter: 'blur(10px)',
};

const dustStyle: CSSProperties = {
  position: 'absolute',
  inset: 0,
  opacity: 0.34,
  backgroundImage:
    'radial-gradient(circle, rgba(42, 91, 124, 0.36) 0 1px, transparent 1.6px), radial-gradient(circle, rgba(146, 126, 250, 0.30) 0 1px, transparent 1.7px)',
  backgroundPosition: '4px 10px, 21px 28px',
  backgroundSize: '38px 42px, 54px 58px',
};

const pathStyle: CSSProperties = {
  position: 'absolute',
  left: '29%',
  top: '45%',
  width: '44%',
  height: '2px',
  transform: 'rotate(24deg)',
  transformOrigin: 'left center',
  borderRadius: '999px',
  background: 'linear-gradient(90deg, rgba(139,124,252,0.06), rgba(139,124,252,0.58), rgba(117,215,240,0.16))',
  boxShadow: '0 0 14px rgba(139,124,252,0.22)',
};

const panelStyle: CSSProperties = {
  border: '1px solid rgba(151, 190, 210, 0.34)',
  borderRadius: '20px',
  padding: '12px',
  background: 'rgba(255, 255, 255, 0.86)',
  boxShadow: '0 14px 36px rgba(35, 79, 112, 0.10)',
};

const panelTitleStyle: CSSProperties = {
  margin: 0,
  color: '#113A50',
  fontSize: '16px',
  lineHeight: 1.18,
  fontWeight: 900,
  letterSpacing: 0,
  overflowWrap: 'anywhere',
};

const panelMetaStyle: CSSProperties = {
  margin: '6px 0 0',
  color: '#557084',
  fontSize: '12px',
  lineHeight: 1.35,
  fontWeight: 750,
};

const progressTrackStyle: CSSProperties = {
  height: '8px',
  marginTop: '10px',
  overflow: 'hidden',
  borderRadius: '999px',
  background: '#E5F1F4',
};

const primaryButtonStyle: CSSProperties = {
  width: '100%',
  minHeight: '42px',
  marginTop: '10px',
  border: 0,
  borderRadius: '14px',
  background: 'linear-gradient(135deg, #41C7B3 0%, #6E8CFB 100%)',
  color: '#FFFFFF',
  fontFamily: 'inherit',
  fontSize: '13px',
  fontWeight: 900,
  cursor: 'pointer',
  boxShadow: '0 12px 28px rgba(65, 199, 179, 0.24)',
};

const createFormStyle: CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) auto',
  gap: '8px',
  alignItems: 'center',
};

const inputStyle: CSSProperties = {
  minWidth: 0,
  height: '40px',
  border: '1px solid rgba(151, 190, 210, 0.48)',
  borderRadius: '14px',
  padding: '0 11px',
  background: 'rgba(255,255,255,0.92)',
  color: '#12384E',
  fontFamily: 'inherit',
  fontSize: '12px',
  fontWeight: 750,
  outline: 'none',
};

const createButtonStyle: CSSProperties = {
  height: '40px',
  minWidth: '72px',
  border: 0,
  borderRadius: '14px',
  background: '#12384E',
  color: '#FFFFFF',
  fontFamily: 'inherit',
  fontSize: '12px',
  fontWeight: 900,
  cursor: 'pointer',
};

export default function CosmosMobileFirstView({
  courses,
  todayTask,
  userName,
  isSubmitting,
  creationError,
  onCreateStar,
  onOpenCourse,
  onOpenTask,
}: Props) {
  const visibleCourses = useMemo(() => courses.filter((course) => !course.isInactive), [courses]);
  const recommendedCourse = useMemo(() => getRecommendedCourse(visibleCourses), [visibleCourses]);
  const [selectedCourseId, setSelectedCourseId] = useState<string | null>(recommendedCourse?.id ?? null);
  const [newTopic, setNewTopic] = useState('');

  const starCourses = useMemo(() => {
    const seed = recommendedCourse ? [recommendedCourse] : [];
    visibleCourses.forEach((course) => {
      if (!seed.some((item) => item.id === course.id)) seed.push(course);
    });
    return seed.slice(0, starSlots.length);
  }, [recommendedCourse, visibleCourses]);

  const selectedCourse = starCourses.find((course) => course.id === selectedCourseId) ?? recommendedCourse ?? starCourses[0] ?? null;
  const selectedTone = statusTone[selectedCourse?.status ?? 'ready'];
  const progress = getProgress(selectedCourse);
  const routeTitle = todayTask?.title ?? selectedCourse?.title ?? '첫 번째 학습별을 만들어볼까요?';
  const routeDescription = todayTask?.description ?? (selectedCourse
    ? `${selectedCourse.title}에서 다음 지점을 이어갈 수 있어요.`
    : `${userName || '학습자'}님의 우주에 새 별을 띄울 준비가 되었어요.`);

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void onCreateStar(newTopic);
  };

  return (
    <div style={rootStyle}>
      <div style={shellStyle}>
        <section style={routeStyle} aria-label="오늘의 추천 항로">
          <p style={routeEyebrowStyle}>Lumi Route</p>
          <p style={routeTextStyle}>{routeDescription}</p>
        </section>

        <section style={mapStyle} aria-label="내 학습별 지도">
          <div aria-hidden="true" style={galaxyVeilStyle} />
          <div aria-hidden="true" style={dustStyle} />
          <div aria-hidden="true" style={pathStyle} />
          {starCourses.length > 0 ? starCourses.map((course, index) => {
            const slot = starSlots[index];
            const tone = statusTone[course.status];
            const selected = selectedCourse?.id === course.id;
            const size = selected ? slot.size + 8 : slot.size;
            return (
              <button
                key={course.id}
                type="button"
                aria-label={`${course.title} ${tone.label}`}
                title={course.title}
                onClick={() => setSelectedCourseId(course.id)}
                style={{
                  position: 'absolute',
                  left: `${slot.left}%`,
                  top: `${slot.top}%`,
                  width: `${size}px`,
                  height: `${size}px`,
                  transform: 'translate(-50%, -50%)',
                  borderRadius: '999px',
                  border: selected ? `2px solid ${tone.ring}` : '1px solid rgba(255,255,255,0.72)',
                  background: `radial-gradient(circle at 35% 30%, #FFFFFF 0 8%, ${tone.core} 36%, rgba(255,255,255,0.34) 100%)`,
                  boxShadow: `0 0 ${selected ? 34 : 18}px ${tone.glow}, 0 0 0 ${selected ? 9 : 0}px ${selected ? tone.ring : 'transparent'}`,
                  cursor: 'pointer',
                  transition: 'transform 160ms ease, box-shadow 160ms ease, width 160ms ease, height 160ms ease',
                }}
              >
                <span style={{
                  position: 'absolute',
                  left: '50%',
                  top: 'calc(100% + 7px)',
                  maxWidth: '84px',
                  transform: 'translateX(-50%)',
                  color: selected ? '#12384E' : 'rgba(18,56,78,0.68)',
                  fontSize: selected ? '10px' : '0',
                  fontWeight: 900,
                  lineHeight: 1.1,
                  textAlign: 'center',
                  overflow: 'hidden',
                  display: '-webkit-box',
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: 'vertical',
                }}>
                  {clampTitle(course.title)}
                </span>
              </button>
            );
          }) : (
            <div style={{
              position: 'absolute',
              left: '50%',
              top: '48%',
              width: '42px',
              height: '42px',
              transform: 'translate(-50%, -50%)',
              borderRadius: '999px',
              border: '2px dashed rgba(139, 124, 252, 0.48)',
              boxShadow: '0 0 36px rgba(139, 124, 252, 0.22)',
            }} />
          )}
        </section>

        <section style={panelStyle} aria-label="선택한 학습별">
          <p style={{ margin: 0, color: selectedTone.core, fontSize: '11px', fontWeight: 900 }}>{selectedTone.label}</p>
          <h1 style={panelTitleStyle}>{selectedCourse?.title ?? routeTitle}</h1>
          <p style={panelMetaStyle}>
            {selectedCourse
              ? `진행률 ${progress}% · 지역 ${selectedCourse.lessonCount || selectedCourse.plannedLessonCount || 0}개`
              : '아직 학습별이 없어요. 새 별을 만들면 이곳에 항로가 열립니다.'}
          </p>
          <div style={progressTrackStyle} aria-hidden="true">
            <div style={{
              width: `${Math.min(100, Math.max(0, progress))}%`,
              height: '100%',
              borderRadius: '999px',
              background: `linear-gradient(90deg, ${selectedTone.core}, #41C7B3)`,
            }} />
          </div>
          <button
            type="button"
            style={primaryButtonStyle}
            onClick={() => {
              if (todayTask?.href) onOpenTask(todayTask.href);
              else if (selectedCourse) onOpenCourse(selectedCourse);
            }}
            disabled={!selectedCourse && !todayTask}
          >
            {todayTask?.cta_label ?? (selectedCourse ? '이어서 탐험하기' : '첫 별 만들기')}
          </button>
        </section>

        <form style={createFormStyle} onSubmit={handleSubmit} aria-label="새 학습별 만들기">
          <input
            value={newTopic}
            onChange={(event) => setNewTopic(event.target.value)}
            style={inputStyle}
            placeholder="무엇을 배우고 싶나요?"
            aria-label="새 학습별 주제"
          />
          <button type="submit" style={createButtonStyle} disabled={isSubmitting}>
            {isSubmitting ? '생성중' : '새 별'}
          </button>
          {creationError ? (
            <p style={{ gridColumn: '1 / -1', margin: 0, color: '#B94747', fontSize: '11px', fontWeight: 800 }}>
              {creationError}
            </p>
          ) : null}
        </form>
      </div>
    </div>
  );
}
