import { useEffect, useMemo, useState } from 'react';

import type {
  PlanetRouteKind,
  PlanetMetaState,
  ResolvedPoint,
  LearningPointMutationResponse,
} from '../pointPageTypes';
import type { PointLearningCopy } from '@/lib/i18n/pages/pointLearning';
import { resolveSafetyInputMessage, type SafetyAPIErrorPayload } from '@/lib/safetyErrors';
import { createJournalNoteSignature } from './pointLearningUtils';

type UsePointJournalParams = {
  routeKind: PlanetRouteKind;
  pointID: string;
  planetID: string;
  planet: PlanetMetaState | null;
  pointDetail: ResolvedPoint | null;
  applyPointMutationPayload: (payload: LearningPointMutationResponse) => boolean;
  setPointEntryMessage: (message: string | null) => void;
  copy: PointLearningCopy['workspace']['runtime'];
};

export function usePointJournal({
  routeKind,
  pointID,
  planetID,
  planet,
  pointDetail,
  applyPointMutationPayload,
  setPointEntryMessage,
  copy,
}: UsePointJournalParams) {
  const [journalObservation, setJournalObservation] = useState('');
  const [journalReflection, setJournalReflection] = useState('');
  const [journalExamples, setJournalExamples] = useState('');
  const [journalNextStep, setJournalNextStep] = useState('');
  const [savedJournalNoteSignature, setSavedJournalNoteSignature] = useState('');
  const [isSavingJournal, setIsSavingJournal] = useState(false);

  const hasSavedJournalNote = useMemo(
    () => savedJournalNoteSignature.trim().length > 0,
    [savedJournalNoteSignature],
  );

  const readSavedJournalDraft = () => ({
    observation: pointDetail?.journal_entry?.core_concept ?? pointDetail?.journal_entry?.observation ?? '',
    reflection: pointDetail?.journal_entry?.my_explanation ?? pointDetail?.journal_entry?.reflection ?? '',
    examples: pointDetail?.journal_entry?.examples ?? '',
    nextStep: pointDetail?.journal_entry?.confused_parts ?? pointDetail?.journal_entry?.next_step ?? '',
  });

  useEffect(() => {
    if (!pointDetail) return;
    const loaded = readSavedJournalDraft();
    setJournalObservation(loaded.observation);
    setJournalReflection(loaded.reflection);
    setJournalExamples(loaded.examples);
    setJournalNextStep(loaded.nextStep);
    setSavedJournalNoteSignature(createJournalNoteSignature(
      loaded.observation,
      loaded.reflection,
      loaded.examples,
      loaded.nextStep,
    ));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pointDetail]);

  const resetJournalDraftToSaved = () => {
    const loaded = readSavedJournalDraft();
    setJournalObservation(loaded.observation);
    setJournalReflection(loaded.reflection);
    setJournalExamples(loaded.examples);
    setJournalNextStep(loaded.nextStep);
  };

  const handleSaveJournal = async () => {
    if (!planet || !planetID || routeKind !== 'learning' || isSavingJournal) return false;
    setIsSavingJournal(true);
    setPointEntryMessage(null);
    try {
      const res = await fetch(`/api/v1/planets/learning/${planetID}/points/${pointID}/journal`, {
        method: 'PATCH',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          observation: journalObservation,
          reflection: journalReflection,
          next_step: journalNextStep,
          core_concept: journalObservation,
          my_explanation: journalReflection,
          examples: journalExamples,
          confused_parts: journalNextStep,
          reference_links: '',
        }),
      });
      const payload = (await res.json().catch(() => ({}))) as LearningPointMutationResponse & SafetyAPIErrorPayload;
      if (!res.ok || !applyPointMutationPayload(payload)) {
        throw new Error(resolveSafetyInputMessage(payload, payload.error ?? copy.journalSaveFailed));
      }
      setSavedJournalNoteSignature(createJournalNoteSignature(
        journalObservation,
        journalReflection,
        journalExamples,
        journalNextStep,
      ));
      setPointEntryMessage(copy.journalSaved);
      return true;
    } catch (e) {
      setPointEntryMessage(e instanceof Error ? e.message : copy.journalSaveFailed);
      return false;
    } finally {
      setIsSavingJournal(false);
    }
  };

  return {
    journalObservation,
    setJournalObservation,
    journalReflection,
    setJournalReflection,
    journalExamples,
    setJournalExamples,
    journalNextStep,
    setJournalNextStep,
    hasSavedJournalNote,
    isSavingJournal,
    handleSaveJournal,
    resetJournalDraftToSaved,
  };
}
