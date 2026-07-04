import type { Dispatch, SetStateAction } from 'react';

import type {
  PlanetRouteKind, PlanetMetaState, ResolvedPoint,
  PointQuestionType, PointReplacementCandidate, PointWorkTab, PlanetPoint, ObservationNoteType,
  LearningPointSelfEvaluationApplicationQuestion,
} from '../pointPageTypes';
import type {
  ResearchBlockDraft, PointArtifactDraft, PointAttachmentDraft,
  PointPracticeLogDraft, PointQuestionDraft,
} from '../pointPageUtils';
import type { getPointCompletionReadiness } from '../pointPageUtils';

export interface UsePointLearningResult {
  routeKind: PlanetRouteKind;
  planetID: string;
  pointID: string;
  planet: PlanetMetaState | null;
  pointDetail: ResolvedPoint | null;
  error: string | null;
  isLoading: boolean;
  uiLocale: 'ko' | 'en';

  goalMetaItems: string[];
  previousPointID: string | null;
  nextPointID: string | null;
  pointCompletionReadiness: ReturnType<typeof getPointCompletionReadiness>;
  learningGoal: string | null;
  isResearchMaterialConfirmed: boolean;
  canConfirmResearchMaterial: boolean;
  isConfirmingResearchMaterial: boolean;
  handleConfirmResearchMaterial: () => Promise<boolean>;
  handleUnconfirmResearchMaterial: () => Promise<boolean>;

  pointRuntimeMessage: string | null;
  pointEntryMessage: string | null;

  // point goal
  pointGoal: string; setPointGoal: Dispatch<SetStateAction<string>>;
  pointCategory: string; setPointCategory: Dispatch<SetStateAction<string>>;
  isSavingPointGoal: boolean;
  handleSavePointGoal: () => Promise<void>;

  // journal
  journalObservation: string; setJournalObservation: Dispatch<SetStateAction<string>>;
  journalReflection: string; setJournalReflection: Dispatch<SetStateAction<string>>;
  journalExamples: string; setJournalExamples: Dispatch<SetStateAction<string>>;
  journalNextStep: string; setJournalNextStep: Dispatch<SetStateAction<string>>;
  hasSavedJournalNote: boolean;
  isSavingJournal: boolean;
  handleSaveJournal: () => Promise<boolean>;
  resetJournalDraftToSaved: () => void;


  // observation notes
  isSavingObservationNote: boolean;
  handleCreateObservationNote: (noteType: ObservationNoteType, content: string) => Promise<boolean>;
  handleUpdateObservationNote: (noteID: string, noteType: ObservationNoteType, content: string) => Promise<boolean>;
  handleDeleteObservationNote: (noteID: string) => Promise<boolean>;

  // record
  recordStudyMinutes: string; setRecordStudyMinutes: Dispatch<SetStateAction<string>>;
  recordPracticeCount: string; setRecordPracticeCount: Dispatch<SetStateAction<string>>;
  recordConfidenceLevel: string; setRecordConfidenceLevel: Dispatch<SetStateAction<string>>;
  recordApplicationNote: string; setRecordApplicationNote: Dispatch<SetStateAction<string>>;
  isSavingRecord: boolean;
  handleSaveRecord: () => Promise<void>;

  // artifacts
  artifactType: string; setArtifactType: Dispatch<SetStateAction<string>>;
  artifactTitle: string; setArtifactTitle: Dispatch<SetStateAction<string>>;
  artifactURL: string; setArtifactURL: Dispatch<SetStateAction<string>>;
  artifactDescription: string; setArtifactDescription: Dispatch<SetStateAction<string>>;
  artifactPointCategory: string; setArtifactPointCategory: Dispatch<SetStateAction<string>>;
  artifactProductionProcess: string; setArtifactProductionProcess: Dispatch<SetStateAction<string>>;
  artifactLearnedPoints: string; setArtifactLearnedPoints: Dispatch<SetStateAction<string>>;
  artifactDifficultPoints: string; setArtifactDifficultPoints: Dispatch<SetStateAction<string>>;
  artifactVisibility: string; setArtifactVisibility: Dispatch<SetStateAction<string>>;
  artifactDrafts: PointArtifactDraft[];
  isSavingArtifact: boolean;
  handleSaveArtifact: (override?: { artifactType?: string; title?: string }) => Promise<PointArtifactDraft | null>;
  handleChangeArtifactDraft: (id: string, field: keyof Pick<PointArtifactDraft, 'artifactType' | 'title' | 'url' | 'description' | 'pointCategory' | 'productionProcess' | 'learnedPoints' | 'difficultPoints' | 'visibility'>, value: string) => void;
  handleUpdateArtifact: (artifact: PointArtifactDraft) => Promise<boolean>;
  handleDeleteArtifact: (artifact: PointArtifactDraft) => Promise<boolean>;

  // attachments
  attachmentDrafts: PointAttachmentDraft[];
  newAttachment: PointAttachmentDraft;
  newResearchMaterialAttachment: PointAttachmentDraft;
  selectedAttachmentFile: File | null; setSelectedAttachmentFile: Dispatch<SetStateAction<File | null>>;
  selectedResearchMaterialFile: File | null; setSelectedResearchMaterialFile: Dispatch<SetStateAction<File | null>>;
  maxResearchMaterialAttachmentCount: number;
  canAddResearchMaterialAttachment: boolean;
  researchMaterialAttachmentLimitMessage: string | null;
  researchMaterialAttachmentValidationMessage: string | null;
  setResearchMaterialAttachmentValidationMessage: Dispatch<SetStateAction<string | null>>;
  researchMaterialUploadProgress: number | null;
  isSavingAttachment: boolean;
  handleCreateAttachment: () => Promise<boolean>;
  handleCreateResearchMaterialAttachment: () => Promise<void>;
  handleUploadAttachment: (fileOverride?: File | null, titleOverride?: string) => Promise<boolean>;
  handleReplaceAttachmentFile: (attachment: PointAttachmentDraft, file: File, titleOverride?: string) => Promise<boolean>;
  handleUploadResearchMaterialAttachment: (file?: File, forcedAttachmentType?: 'video' | 'subtitle' | 'thumbnail', titleOverride?: string) => Promise<boolean>;
  handleUploadArtifactAttachment: (artifactID: string, file: File, forcedAttachmentType?: 'video' | 'subtitle' | 'thumbnail', titleOverride?: string) => Promise<boolean>;
  handleUploadResearchMaterialInlineImage: (file: File) => Promise<{ src: string; alt?: string }>;
  handleOpenAttachment: (attachment: PointAttachmentDraft) => Promise<void>;
  handleUpdateAttachment: (attachment: PointAttachmentDraft) => Promise<boolean>;
  handleDeleteAttachment: (attachment: PointAttachmentDraft) => Promise<boolean>;
  handleCreateLinkAttachment: (title: string, url: string) => Promise<boolean>;
  handleChangeAttachmentDraft: (id: string, field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => void;
  handleChangeNewAttachment: (field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => void;
  handleChangeNewResearchMaterialAttachment: (field: keyof Omit<PointAttachmentDraft, 'id'>, value: string) => void;

  // material report
  isReportingMaterial: boolean;
  materialReportTarget: 'source' | 'ai_summary' | 'attachment' | 'other';
  setMaterialReportTarget: Dispatch<SetStateAction<'source' | 'ai_summary' | 'attachment' | 'other'>>;
  materialReportType: 'broken_link' | 'wrong_content' | 'unsafe_content' | 'copyright' | 'low_quality' | 'other';
  setMaterialReportType: Dispatch<SetStateAction<'broken_link' | 'wrong_content' | 'unsafe_content' | 'copyright' | 'low_quality' | 'other'>>;
  materialReportMessage: string; setMaterialReportMessage: Dispatch<SetStateAction<string>>;
  handleReportMaterial: () => Promise<void>;
  replacementCandidates: PointReplacementCandidate[];
  replacementMessage: string | null;
  isLoadingReplacementCandidates: boolean;
  handleLoadReplacementCandidates: () => Promise<void>;
  handleUseReplacementCandidate: (candidate: PointReplacementCandidate, reportID: string) => Promise<boolean>;
  handleUseReplacementLink: (title: string, url: string, reportID: string) => Promise<boolean>;
  handleCancelMaterialReport: (reportID: string) => Promise<boolean>;

  // practice logs
  practiceLogDrafts: PointPracticeLogDraft[];
  newPracticeLog: PointPracticeLogDraft;
  isSavingPracticeLog: boolean;
  handleCreatePracticeLog: () => Promise<boolean>;
  handleUpdatePracticeLog: (log: PointPracticeLogDraft) => Promise<boolean>;
  handleDeletePracticeLog: (log: PointPracticeLogDraft) => Promise<boolean>;
  handleChangePracticeLogDraft: (id: string, field: keyof Omit<PointPracticeLogDraft, 'id'>, value: string) => void;
  handleChangeNewPracticeLog: (field: keyof Omit<PointPracticeLogDraft, 'id'>, value: string) => void;

  // AI summary
  isGeneratingAISummary: boolean;
  handleGenerateAISummary: () => Promise<void>;

  // research blocks
  researchBlocks: ResearchBlockDraft[];
  maxResearchBlockCount: number;
  canAddResearchBlock: boolean;
  researchBlockLimitMessage: string | null;
  isSavingResearchBlock: boolean;
  handleAddResearchBlock: (blockType: import('../pointPageTypes').ResearchBlockType) => string | null;
  handleChangeResearchBlock: (id: string, field: keyof Pick<ResearchBlockDraft, 'text' | 'url' | 'title' | 'caption' | 'note'>, value: string) => void;
  handleSaveResearchBlock: (block: ResearchBlockDraft) => Promise<void>;
  handleDeleteResearchBlock: (block: ResearchBlockDraft) => Promise<void>;

  // questions
  pointQuestions: PointQuestionDraft[];
  newQuestionTitle: string; setNewQuestionTitle: Dispatch<SetStateAction<string>>;
  newQuestionText: string; setNewQuestionText: Dispatch<SetStateAction<string>>;
  newQuestionType: PointQuestionType; setNewQuestionType: Dispatch<SetStateAction<PointQuestionType>>;
  isSavingQuestion: boolean;
  isGeneratingFeedbackQuestionID: string | null;
  refreshPointDetail: () => Promise<boolean>;
  handleCreateQuestion: () => Promise<boolean>;
  handleChangeQuestionDraft: (questionID: string, field: keyof Pick<PointQuestionDraft, 'title' | 'question' | 'questionType' | 'answerMethod' | 'answer' | 'aiFeedback'>, value: string) => void;
  handleChangeQuestion: (questionID: string, value: string) => void;
  handleChangeQuestionAnswerMethod: (questionID: string, value: string) => void;
  handleSaveQuestion: (question: PointQuestionDraft) => Promise<boolean>;
  handleDeleteQuestion: (question: PointQuestionDraft) => Promise<void>;
  handleGenerateQuestionFeedback: (question: PointQuestionDraft) => Promise<void>;

  // self-evaluation
  selfEvalUnderstanding: string; setSelfEvalUnderstanding: Dispatch<SetStateAction<string>>;
  selfEvalUnderstandingReason: string; setSelfEvalUnderstandingReason: Dispatch<SetStateAction<string>>;
  selfEvalApplicationScore: string; setSelfEvalApplicationScore: Dispatch<SetStateAction<string>>;
  selfEvalApplicationReason: string; setSelfEvalApplicationReason: Dispatch<SetStateAction<string>>;
  selfEvalProficiency: string; setSelfEvalProficiency: Dispatch<SetStateAction<string>>;
  selfEvalProficiencyReason: string; setSelfEvalProficiencyReason: Dispatch<SetStateAction<string>>;
  selfEvalProblemSolvingScore: string; setSelfEvalProblemSolvingScore: Dispatch<SetStateAction<string>>;
  selfEvalProblemSolvingReason: string; setSelfEvalProblemSolvingReason: Dispatch<SetStateAction<string>>;
  selfEvalExpressionScore: string; setSelfEvalExpressionScore: Dispatch<SetStateAction<string>>;
  selfEvalExpressionReason: string; setSelfEvalExpressionReason: Dispatch<SetStateAction<string>>;
  selfEvalGoalAlignmentNote: string; setSelfEvalGoalAlignmentNote: Dispatch<SetStateAction<string>>;
  selfEvalApplicationQuestions: LearningPointSelfEvaluationApplicationQuestion[];
  selfEvalApplicationAnswers: string[];
  hasSelfEvalDraft: boolean;
  isSelfEvalLocked: boolean;
  finalSelfEvalScore: string; setFinalSelfEvalScore: Dispatch<SetStateAction<string>>;
  isSavingSelfEvaluation: boolean;
  isGeneratingSelfEvalDraft: boolean;
  handleChangeSelfEvalApplicationAnswer: (index: number, value: string) => void;
  handleRestartSelfEvaluation: () => Promise<void>;
  handleSaveSelfEvaluation: () => Promise<boolean>;
  handleGenerateSelfEvaluationDraft: () => Promise<void>;

  // completion/runtime
  isCompletionOpen: boolean; setIsCompletionOpen: Dispatch<SetStateAction<boolean>>;
  completionChecks: boolean[]; setCompletionChecks: Dispatch<SetStateAction<boolean[]>>;
  showCompletionAchievement: boolean;
  isSavingPointRuntime: boolean;
  handleSavePointRuntime: (status: 'in_progress' | 'completed') => Promise<void>;

  // work tab
  activeWorkTab: PointWorkTab; setActiveWorkTab: Dispatch<SetStateAction<PointWorkTab>>;

  answeredQuestionCount: number;
}
