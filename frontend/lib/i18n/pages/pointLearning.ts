import type { Locale } from '@/lib/i18n/locales'
import { normalizeLocale } from '@/lib/i18n/locales'

export type PointLearningCopy = {
  loading: string
  notFound: string
  backToPlanet: string
  header: {
    diaryMeta: string
    diaryLabel: string
    diaryTitle: string
    lightMode: string
    darkMode: string
  }
  toolbar: {
    ariaLabel: string
    expand: string
    collapse: string
    content: string
    researchContent: string
    work: string
    evaluation: string
    evaluationDisabled: string
    evaluationDisabledTitle: string
    lumiOpen: string
    lumiClose: string
    lumiLabel: string
  }
  skillFlow: {
    progressLabel: string
    back: string
    next: string
    completePoint: string
    stages: {
      point: { label: string; title: string; subtitle: string }
      goal: { label: string; title: string; subtitle: string }
      skill_discovery: { label: string; title: string; subtitle: string }
      goal_iteration: { label: string; title: string; subtitle: string }
      final_result: { label: string; title: string; subtitle: string }
      complete: { label: string; title: string; subtitle: string }
    }
    observe: {
      externalTitle: string
      externalDescription: string
      internalTitle: string
      internalDescription: string
      researchTitle: string
      researchDescription: string
      researchEditHint: string
      externalSpeech: string
      internalSpeech: string
      researchSpeech: string
      notePanelTitle: string
      notePanelDescription: string
      noteTypeLabel: string
      noteContentLabel: string
      noteTypes: Record<'core_summary' | 'revisit_part' | 'reference_material', { label: string; placeholder: string }>
      noteAdd: string
      noteSave: string
      noteCancel: string
      noteEdit: string
      noteDelete: string
      noteEmpty: string
      noteSaved: string
      noteDeleted: string
      noteSaveFailed: string
    }
    goal: {
      contextLabel: string
      titleLabel: string
      titlePlaceholder: string
      typeLabel: string
      types: Record<'creative' | 'conceptual' | 'physical' | 'coding' | 'other', { label: string; description: string }>
    }
    skills: {
      prompt: string
      presets: Record<'creative' | 'conceptual' | 'physical' | 'coding' | 'other', string[]>
      emptyGoalHint: string
    }
    repeat: {
      requiredHint: string
      lensLabel: string
      selectedLensLabel: string
      connectedRecordLabel: string
      selectedSkillsLabel: string
      noSelectedSkills: string
      goalObjectLabel: string
      emptyGoalObject: string
      savedCountLabel: (count: number) => string
      lenses: Record<'attempt' | 'evidence' | 'known' | 'question' | 'answer' | 'works_well' | 'needs_practice' | 'resource', { label: string; description: string }>
    }
    finalResult: {
      artifactHint: string
      openArtifacts: string
    }
    complete: {
      readinessHint: string
    }
  }
  workspace: {
    materialTitle: string
    explorationSubtitle: string
    researchContentActions: {
      edit: string
      saveNew: string
      saveChanges: string
      saving: string
      saved: string
    }
    runtime: {
      invalidPointPath: string
      inactiveDiaryOnlyPlanning: string
      pointLoadFailed: string
      completionNeedsReadiness: string
      researchContentRequired: string
      learningPointNotFound: string
      pointStatusSaveFailed: string
      pointCompleted: string
      pointCompletionCancelled: string
      pointInProgress: string
      pointGoalSaveFailed: string
      pointGoalSaved: string
      pointRecordSaveFailed: string
      pointRecordSaved: string
      journalSaveFailed: string
      journalSaved: string
      practiceCreateFailed: string
      practiceCreated: string
      practiceUpdateFailed: string
      practiceUpdated: string
      practiceDeleteFailed: string
      practiceDeleted: string
      artifactFallbackTitles: {
        video: string
        attachment: string
        document: string
      }
      artifactSaveFailed: string
      artifactCreated: string
      artifactUpdateFailed: string
      artifactUpdated: string
      artifactDeleteFailed: string
      artifactDeleted: string
      artifactSubtitleDuplicate: string
      artifactThumbnailDuplicate: string
      artifactVideoDuplicate: string
      artifactAttachmentLimit: string
      artifactUploadFailedStatus: (status: number) => string
      artifactUploadFailed: string
      artifactUploadSuccess: string
      attachmentCreateFailed: string
      attachmentCreated: string
      attachmentLinkInvalid: string
      replacementAttachmentCreateFailed: string
      attachmentUploadFailed: string
      attachmentUploadSuccess: string
      attachmentReplaceFailed: string
      attachmentReplaced: string
      attachmentOpenFailed: string
      attachmentUpdateFailed: string
      attachmentUpdated: string
      attachmentDeleteLocked: string
      attachmentDeleteFailed: string
      questionSaveFailed: string
      questionCreated: string
      questionAnswerSaved: string
      questionAnswerSaveFailed: string
      questionDeleteFailed: string
      questionDeleted: string
      aiFeedbackUnavailable: string
      aiFeedbackFailed: string
      aiFeedbackCreated: string
      selfEvalDraftRequired: string
      selfEvalSaveFailed: string
      selfEvalSaved: string
      selfEvalDraftFailed: string
      selfEvalSavedNotice: string
      selfEvalQuestionsReady: string
    }
    explorationContent: {
      dateLocale: string
      reportTypeLabels: Record<string, string>
      reportStatusLabels: Record<string, string>
      aiSummaryTitle: string
      aiSummaryGenerating: string
      aiSummaryRegenerate: string
      aiSummaryGenerate: string
      aiSummaryEmpty: string
      providedContentTitle: string
      openCurrentContent: string
      youtubeSourceLabel: string
      openYoutubePlayer: string
      openOriginalContent: string
      sourceLayoutHint: string
      emptyExternalLink: string
      replacementCompletedNotice: string
      issueSummary: string
      issueHelp: string
      reportTitle: string
      reportHelp: string
      reportPlaceholder: string
      reportSubmitting: string
      reportAction: string
      reportGateHelp: string
      replacementGateNotice: string
      cancelReport: string
      replacementLoading: string
      loadReplacement: string
      replacementCandidatesAria: string
      candidateColumn: string
      replaceColumn: string
      untitledCandidate: string
      replace: string
      directReplacementTitle: string
      directReplacementHelp: string
      directReplacementTitlePlaceholder: string
      directReplacementAction: string
      myReportsTitle: string
      reportCount: (count: number) => string
      reportHistoryAria: string
      reportColumn: string
      statusColumn: string
      receivedAtColumn: string
      replacedPrefix: (value: string) => string
      emptyReports: string
      cancelReportSuccessNotice: string
      replacementQueryMissing: string
      insufficientPoints: string
      replacementLoadFailed: string
      replacementFound: string
      replacementEmpty: string
      replacementUrlInvalid: string
      replacementFallbackTitle: string
      replacementFailed: string
      replacementSuccess: string
      replacementEntrySuccess: string
      candidateUrlMissing: string
      cancelReportFailed: string
      cancelReportEntrySuccess: string
      reportFailed: string
      reportSuccess: string
      aiUnavailable: string
      aiSummaryFailed: string
      aiSummarySuccess: string
    }
    editor: {
      linkHttpsOnly: string
      imageHttpsOnly: string
      imageUploadFailed: string
      paragraphStyle: string
      paragraph: string
      heading1: string
      heading2: string
      heading3: string
      undo: string
      redo: string
      quote: string
      codeBlock: string
      bold: string
      italic: string
      underline: string
      strike: string
      highlight: string
      inlineCode: string
      superscript: string
      subscript: string
      bulletList: string
      orderedList: string
      alignLeft: string
      alignCenter: string
      alignRight: string
      link: string
      image: string
      insertTable: string
      addColumn: string
      addRow: string
      deleteColumn: string
      deleteRow: string
      deleteTable: string
      apply: string
      unset: string
      imageUrlPlaceholder: string
      imageAltPlaceholder: string
      insertImage: string
      chooseFile: string
      uploadHint: string
      uploading: string
      uploadInsert: string
      inlineImageTypeError: string
      inlineImageSizeError: string
    }
    researchMaterial: {
      pageTitle: string
      nextActionNotice: string
      blockLimitMessage: string
      blockAddUnavailable: string
      blockSaveFailed: string
      blockAddSuccess: string
      blockSaveSuccess: string
      deleteLockedMessage: string
      blockDeleteFailed: string
      confirmUnavailable: string
      confirmFailed: string
      confirmSuccess: string
      attachmentLimitMessage: (count: number) => string
      attachmentDirtyMessage: string
      attachmentAddUnavailable: string
      attachmentAddFailed: string
      attachmentAddSuccess: string
      subtitleDuplicate: string
      thumbnailDuplicate: string
      videoDuplicate: string
      referenceTypeError: string
      imageSizeError: string
      videoTypeError: string
      videoSizeError: string
      videoDurationError: string
      videoMetadataError: string
      subtitleTypeError: string
      subtitleSizeError: string
      thumbnailTypeError: string
      thumbnailSizeError: string
      documentSizeError: string
      uploadTooLargeServer: string
      uploadFailedStatus: (status: number) => string
      fileUploadFailed: string
      fileUploadSuccess: string
      inlineImageUnavailable: string
      inlineImageTypeError: string
      inlineImageSizeError: string
      inlineImageUploadFailed: string
      inlineImageInserted: string
      statusConfirmed: string
      statusDraft: string
      confirmSaving: string
      confirmDone: string
      confirmAction: string
      confirmModalTitle: string
      confirmModalMessage: string
      unconfirmAction: string
      unconfirmModalTitle: string
      unconfirmModalMessage: string
      unconfirmSaving: string
      unconfirmFailed: string
      unconfirmSuccess: string
      confirmedHelp: string
      confirmReadyHelp: string
      confirmBlockedHelp: string
      attachmentFallback: string
      videoUrlFailed: string
      contentFallback: string
      attachmentsTitle: string
      uploadedFileListTitle: string
      openFileLink: string
      attachmentLimitTitle: string
      attachmentAddTitle: string
      uploadingAttachment: string
      addAttachment: string
      attachmentDeleteTitle: string
      delete: string
      attachmentImageAlt: string
      emptyAttachments: string
      attachmentHelp: string
      videoSectionTitle: string
      uploadedVideoListTitle: string
      openVideoLink: string
      videoTitleLabel: string
      videoTitlePlaceholder: string
      saving: string
      saveTitle: string
      videoPlayerTitle: string
      uploadSubtitle: string
      uploadThumbnail: string
      deleteSubtitle: string
      deleteThumbnail: string
      deleteVideo: string
      uploadingVideo: string
      uploadVideo: string
      videoHelp: string
      noVideo: string
      contentSectionTitle: string
      addContent: string
      contentTitleLabel: string
      contentTitlePlaceholder: string
      contentPlaceholder: string
      contentDisplayTitle: (title: string) => string
      emptyContent: string
      deleteContentTitle: string
      deleteContentMessage: (title: string) => string
      cancel: string
      deleting: string
      deleteAttachmentModalTitle: string
      deleteAttachmentMessage: (title: string) => string
      validationTitle: string
      validationConfirm: string
      uploadProgressTitle: string
      uploadProgressMessage: string
      uploadProgressAria: (progress: number) => string
      copyrightNotice: string
      copyrightAgreementLinkLabel: string
      statusLabel: string
      observeDraftTitle: string
      observeDraftDescription: string
      openEditMode: string
      navAria: string
      navVideo: string
      navContent: string
      navAttachments: string
    }
    learningWork: {
      title: string
      subtitle: string
      lockedNotice: string
      miniNavLabel: string
      notesLabel: string
      recordsLabel: string
      recordsTitle: string
      optionalLabel: string
      recordsDescription: string
      recordsLockedNotice: string
      moveModalTitle: string
      moveModalEyebrow: string
      moveModalMessage: string
      keepWriting: string
      move: string
      sourceContentGuide: {
        title: string
        message: string
      }
      notes: {
        title: string
        edit: string
        helper: string
        cancelModalTitle: string
        cancelModalEyebrow: string
        cancelModalMessage: string
        back: string
        cancel: string
        cancelButton: string
        saving: string
        save: string
        progressTitle: string
        addEmpty: string
        progressIncomplete: string
        progressComplete: string
        empty: string
        fields: Record<'observation' | 'reflection' | 'examples' | 'nextStep', { icon: string; label: string; placeholder: string }>
      }
      recordPanels: {
        common: {
          number: string
          title: string
          status: string
          type: string
          previous: string
          next: string
          countRange: (total: number, start: number, end: number) => string
        }
        questions: {
          title: string
          subtitle: string
          add: string
          listAria: string
          answered: string
          pending: string
          empty: string
          fallbackTitle: string
          typeLabels: Record<'reflection' | 'application' | 'goal_alignment' | 'none', string>
          detail: {
            editTitle: string
            viewTitle: string
            addTitle: string
            addSubtitle: string
            type: string
            title: string
            content: string
            answer: string
            answerMeta: string
            answerEmptyMeta: string
            aiCoach: string
            aiCoachMeta: string
            titlePlaceholder: string
            contentPlaceholder: string
            answerPlaceholder: string
            noContent: string
            noAnswer: string
            communityDisabled: string
            addAnswer: string
            editAnswer: string
            createAiCoach: string
            createAiCoachStarter: string
            createAiCoachRevision: string
            generatingAiCoach: string
            saveDetail: string
            savingDetail: string
            saveAnswer: string
            savingAnswer: string
            register: string
            saving: string
            edit: string
            list: string
            delete: string
            deleting: string
            cancel: string
            back: string
            ok: string
            answerInstruction: string
            aiCoachNotice: string
            aiCoachStarterHint: string
            aiCoachRevisionHint: string
            aiCoachSummary: string
            aiCoachNextAction: string
            saveSuccessTitle: string
            saveSuccessMessage: string
            saveFailTitle: string
            saveFailMessage: string
            answerSaveSuccessTitle: string
            answerSaveSuccessMessage: string
            answerSaveFailTitle: string
            answerSaveFailMessage: string
            deleteModalTitle: string
            deleteModalMessage: (title: string) => string
            listReturnTitle: string
            listReturnMessage: string
            addCancelTitle: string
            addCancelMessage: string
          }
        }
        practice: {
          title: string
          subtitle: string
          add: string
          listAria: string
          hasRecord: string
          waiting: string
          empty: string
          fallbackTitle: (index: number) => string
          detail: {
            editTitle: string
            viewTitle: string
            addTitle: string
            title: string
            duration: string
            achievement: string
            reflection: string
            titlePlaceholder: string
            durationPlaceholder: string
            achievementPlaceholder: string
            reflectionPlaceholder: string
            addReflection: string
            noDetail: string
            list: string
            delete: string
            deleting: string
            save: string
            saving: string
            edit: string
            cancel: string
            back: string
            deleteModalTitle: string
            deleteModalMessage: (title: string) => string
            listReturnTitle: string
            listReturnMessage: string
            addCancelTitle: string
            addCancelMessage: string
          }
        }
        artifacts: {
          title: string
          subtitle: string
          add: string
          empty: string
          listAria: string
          fallbackTitle: (index: number) => string
          kindLabels: Record<'영상자료' | '문서 편집자료' | '첨부자료', string>
          detail: {
            editTitle: string
            viewTitle: string
            addTitle: string
            kind: string
            title: string
            editTitleLabel: string
            videoSection: string
            attachmentSection: string
            documentSection: string
            videoPlayerTitle: string
            titlePlaceholder: string
            documentTitlePlaceholder: string
            documentPlaceholder: string
            attachmentFallback: string
            noVideo: string
            noFiles: string
            uploadVideo: string
            uploadingVideo: string
            uploadSubtitle: string
            uploadThumbnail: string
            deleteSubtitle: string
            deleteThumbnail: string
            deleteVideo: string
            uploadFile: string
            uploadingFile: string
            list: string
            delete: string
            deleting: string
            save: string
            saving: string
            edit: string
            cancel: string
            back: string
            videoHelp: string
            deleteModalTitle: string
            deleteModalMessage: (title: string) => string
            deleteAttachmentTitle: string
            deleteAttachmentMessage: (title: string) => string
            listReturnTitle: string
            listReturnMessage: string
            addCancelTitle: string
            addCancelMessage: string
          }
        }
        attachments: {
          title: string
          subtitle: (researchCount: number, workCount: number) => string
          add: string
          listAria: string
          learningContent: string
          reference: string
          empty: string
          fallbackTitle: (index: number) => string
          detail: {
            editTitle: string
            viewTitle: string
            addTitle: string
            sectionTitle: string
            title: string
            editTitleLabel: string
            file: string
            titlePlaceholder: string
            attachmentFallback: string
            replacementPending: string
            selectedFile: (name: string) => string
            researchMaterialNotice: string
            uploadFile: string
            uploadingFile: string
            uploadSelected: string
            list: string
            delete: string
            deleting: string
            save: string
            saving: string
            edit: string
            cancel: string
            back: string
            deleteModalTitle: string
            deleteModalMessage: (title: string) => string
            listReturnTitle: string
            listReturnMessage: string
            addCancelTitle: string
            addCancelMessage: string
          }
        }
        timeline: {
          title: string
          subtitle: string
          listAria: string
          time: string
          activity: string
          empty: string
        }
      }
      tabs: Record<'questions' | 'practice' | 'artifacts' | 'attachments' | 'timeline', { label: string; title: string; message: string }>
      noteGuides: Record<'observation' | 'reflection' | 'examples' | 'nextStep', { title: string; message: string }>
    }
  }
  completion: {
    encouragement: string
    finalAction: string
    completedLocked: string
    readyInstruction: string
    cancelProcessing: string
    cancelCompletion: string
    saving: string
    completePoint: string
    completeModalTitle: string
    completeModalEyebrow: string
    completeModalMessage: string
    close: string
    cancel: string
    complete: string
    cancelModalTitle: string
    cancelModalMessage: string
    cancelDone: string
  }
  selfEvaluation: {
    title: string
    subtitle: string
    savedTitle: string
    applicationTitle: string
    savedDescription: string
    applicationDescription: string
    generatingTitle: string
    generatingDescription: string
    questionLabel: string
    answerPlaceholder: string
    aiReferenceTitle: string
    aiReferenceDescription: string
    generatingButton: string
    restartButton: string
    aiReferenceButton: string
    createQuestionsButton: string
    finalTitle: string
    finalDescription: string
    scoreSuffix: string
    lockedUntilDraft: string
    options: Array<{ value: string; label: string; description: string }>
  }
  completionStatus: {
    readOnlyTitle: string
    readOnlyDescription: string
    achievementTitle: string
    achievementDescription: string
    contentChecked: string
    noteCount: (count: number) => string
    questionCount: (total: number, answered: number) => string
    practiceCount: (count: number) => string
    artifactCount: (count: number) => string
    selfEvaluationState: (completed: boolean) => string
    backToDiary: string
    sharedStatusTitle: string
    readyStatusTitle: string
    researchMaterialConfirmed: (completed: boolean) => string
    answeredQuestionReady: (completed: boolean, count: number) => string
    selfEvaluationSaved: (completed: boolean) => string
    selfEvaluationQuality: (completed: boolean) => string
    goalConnection: (completed: boolean) => string
    sharedNotReady: string
    researchNotReady: string
    learningNotReady: string
    sharedReady: string
    learningReady: string
  }
  lumiGuide: {
    goalLabel: string
    close: string
    closeButton: string
    featureGuideLabel: string
    currentTaskLabel: string
    fallbackGoal: string
    pointTitle: (isExplorationPoint: boolean) => string
    pointMessage: (isExplorationPoint: boolean) => string
    learningContentTitle: string
    learningContentMessage: (isExplorationPoint: boolean) => string
    selfEvaluationTitle: string
    selfEvaluationMessage: string
    completedFlow: string
    researchMaterialFlow: string
    openExternalFlow: string
    firstNoteFlow: string
    selfEvaluationReadyFlow: string
    noQuestionFlow: string
    noAnswerFlow: string
    completionReadyFlow: string
    fillMissingFlow: string
    restTitle: string
    restMessage: string
    restConfirm: string
  }
}

const ko: PointLearningCopy = {
  loading: '탐험지점을 준비하는 중...',
  notFound: '지점을 찾지 못했습니다.',
  backToPlanet: '행성 상세로 돌아가기',
  header: {
    diaryMeta: 'Explorer Diary',
    diaryLabel: '탐험일지',
    diaryTitle: '탐험일지로 돌아가기',
    lightMode: '라이트 모드로 전환',
    darkMode: '다크 모드로 전환',
  },
  toolbar: {
    ariaLabel: '학습 진행 메뉴',
    expand: '학습화면 툴바 펼치기',
    collapse: '학습화면 툴바 접기',
    content: '콘텐츠 보기',
    researchContent: '학습 콘텐츠 작성',
    work: '학습하기',
    evaluation: '자기평가',
    evaluationDisabled: '자기평가 비활성',
    evaluationDisabledTitle: '학습하기 기록을 남기면 자기평가를 열 수 있습니다.',
    lumiOpen: '루미 안내 열기',
    lumiClose: '루미 안내 닫기',
    lumiLabel: '루미안내',
  },
  skillFlow: {
    progressLabel: '탐험 단계',
    back: '이전 단계',
    next: '다음 단계',
    completePoint: '완료 단계로 이동',
    stages: {
      point: { label: '지점 관찰', title: '지점 관찰', subtitle: '먼저 이 지점의 자료를 직접 확인하고 무엇을 발견할지 준비합니다.' },
      goal: { label: '목표 정하기', title: '이번 지점 목표 정하기', subtitle: '코스/리슨/포인트 제목은 단서입니다. 내가 최종적으로 만들거나 설명하거나 수행할 목표물을 직접 정합니다.' },
      skill_discovery: { label: '필요 기술 고르기', title: '필요한 기술 고르기', subtitle: '목표물을 얻기 위해 반복해서 써야 할 기술을 골라 봅니다.' },
      goal_iteration: { label: '기록하며 반복하기', title: '오늘 학습 기록하기', subtitle: '오늘 해 본 시도, 질문, 연습, 결과물을 지점 학습 기록 안에 남깁니다.' },
      final_result: { label: '결과 정리', title: '이번 지점 결과 정리', subtitle: '이번 지점에서 남길 대표 발견물과 배운 점을 정리합니다.' },
      complete: { label: '완료 확인', title: '지점 완료 확인', subtitle: '자기평가를 저장하고 지점을 완료합니다.' },
    },
    observe: {
      externalTitle: '외부 자료 관찰',
      externalDescription: '외부 링크 자료는 현재 브라우저의 새 탭/창으로 열어 확인합니다. LearnCosmos에는 요약과 기록 흐름만 유지합니다.',
      internalTitle: '플랫폼 콘텐츠 관찰',
      internalDescription: '크리에이터가 플랫폼 안에 준비한 자료는 이 화면에서 직접 확인합니다.',
      researchTitle: '연구지점 자료 관찰',
      researchDescription: '직접 작성한 연구 자료는 확정된 내용을 관찰합니다.',
      researchEditHint: '자료를 새로 작성하거나 수정해야 하면 편집 모드에서 작성한 뒤 관찰 단계로 돌아옵니다.',
      externalSpeech: '먼저 자료를 열어 어떤 내용인지 살펴보세요.',
      internalSpeech: '먼저 자료를 열어 어떤 내용인지 살펴보세요.',
      researchSpeech: '먼저 연구자료를 만들어 주세요.',
      notePanelTitle: '관찰 기록',
      notePanelDescription: '자료를 보며 발견한 내용을 구분별로 계속 남겨 보세요. 하나만 적어도 다음 단계로 갈 수 있습니다.',
      noteTypeLabel: '게시글 구분',
      noteContentLabel: '내용',
      noteTypes: {
        core_summary: { label: '핵심부분 요약', placeholder: '자료에서 중요하다고 느낀 부분을 내 말로 짧게 적어보세요.' },
        revisit_part: { label: '다시 보고 싶은 부분', placeholder: '다시 확인하고 싶은 장면, 시간, 문장, 위치를 적어보세요.' },
        reference_material: { label: '관련근거자료', placeholder: '근거가 되는 링크, 자료명, 장면, 설명을 적어보세요.' },
      },
      noteAdd: '추가',
      noteSave: '저장',
      noteCancel: '취소',
      noteEdit: '수정',
      noteDelete: '삭제',
      noteEmpty: '아직 관찰 기록이 없습니다. 핵심부분 요약부터 짧게 남겨 보세요.',
      noteSaved: '관찰 기록을 저장했습니다.',
      noteDeleted: '관찰 기록을 삭제했습니다.',
      noteSaveFailed: '관찰 기록 저장에 실패했습니다.',
    },
    goal: {
      contextLabel: '지점 단서',
      titleLabel: '내가 발견할 목표물',
      titlePlaceholder: '예: 짧은 연주 영상, 핵심 개념 설명, 동작 성공 기록, 작동하는 코드',
      typeLabel: '목표물 유형',
      types: {
        creative: { label: '작품/창작형', description: '그림, 영상, 글, 음악처럼 만들어 낼 결과가 있습니다.' },
        conceptual: { label: '개념/학습형', description: '개념을 이해하고 내 말로 설명하는 결과가 있습니다.' },
        physical: { label: '운동/수행형', description: '몸으로 반복해 성공해야 하는 동작이나 루틴이 있습니다.' },
        coding: { label: '코딩/구현형', description: '작동하는 코드, 기능, 분석 결과처럼 실행 가능한 결과가 있습니다.' },
        other: { label: '기타/복합형', description: '여러 유형이 섞였거나 아직 분류하기 어렵습니다.' },
      },
    },
    skills: {
      prompt: '목표물을 얻기 위해 반복해서 써야 할 기술을 선택하세요.',
      presets: {
        creative: ['관찰하기', '구성하기', '초안 만들기', '수정하기', '표현 다듬기'],
        conceptual: ['핵심 찾기', '내 말로 설명하기', '예시 만들기', '질문 만들기', '적용해 보기'],
        physical: ['자세 확인', '반복 시도', '속도 조절', '오류 고치기', '기록 비교'],
        coding: ['요구사항 읽기', '작게 구현하기', '테스트하기', '디버깅하기', '리팩터링하기'],
        other: ['관찰하기', '시도하기', '기록하기', '비교하기', '다시 시도하기'],
      },
      emptyGoalHint: '목표물 유형을 먼저 고르면 추천 기술이 바뀝니다.',
    },
    repeat: {
      requiredHint: 'Phase 1에서는 시도 기록과 탐험 증거를 권장 표시만 하며 완료 조건으로 강제하지 않습니다.',
      lensLabel: '기록 렌즈',
      selectedLensLabel: '현재 렌즈',
      connectedRecordLabel: '연결 기록 영역',
      selectedSkillsLabel: '오늘 사용한 기술',
      noSelectedSkills: '아직 오늘 사용한 기술을 고르지 않았습니다.',
      goalObjectLabel: '목표물',
      emptyGoalObject: '아직 목표물을 적지 않았습니다.',
      savedCountLabel: (count) => `저장 ${count}`,
      lenses: {
        attempt: { label: '오늘 시도', description: '실제로 해 본 과정, 반복 횟수, 걸린 시간을 남깁니다.' },
        evidence: { label: '탐험 증거', description: '스크린샷, 글, 코드, 영상, 사진, 체크리스트 같은 증거를 남깁니다.' },
        known: { label: '아는 것', description: '확실히 이해했거나 설명할 수 있는 내용을 정리합니다.' },
        question: { label: '궁금한 점', description: '아직 모르는 것과 확인해야 할 질문을 남깁니다.' },
        answer: { label: '답을 찾은 것', description: '질문에 대해 찾은 답과 근거를 기록합니다.' },
        works_well: { label: '잘된 점', description: '오늘 해 보니 안정적으로 된 부분을 기록합니다.' },
        needs_practice: { label: '다음에 해볼 점', description: '막힌 부분과 다음 학습에서 다시 시도할 점을 남깁니다.' },
        resource: { label: '보조자료', description: '추가로 참고한 링크와 파일을 모읍니다.' },
      },
    },
    finalResult: {
      artifactHint: '최종 발견물은 대표 결과물을 새로 강제 선택하지 않고, 기존 결과물/첨부 기록과 연결해 확인합니다.',
      openArtifacts: '결과물 기록 열기',
    },
    complete: {
      readinessHint: '실제 완료 조건은 기존 기준을 유지합니다: 답변 있는 질문 1개 이상과 자기평가 6개 항목 저장.',
    },
  },
  workspace: {
    materialTitle: '학습 콘텐츠',
    explorationSubtitle: '제공된 학습 콘텐츠를 확인하세요',
    researchContentActions: {
      edit: '✏️ 편집',
      saveNew: '💾 콘텐츠 저장',
      saveChanges: '💾 변경 저장',
      saving: '💾 저장 중...',
      saved: '✅ 저장됨',
    },
    runtime: {
      invalidPointPath: '유효하지 않은 지점 경로입니다.',
      inactiveDiaryOnlyPlanning: '비활성 탐험 다이어리는 탐험계획 페이지만 사용할 수 있습니다.',
      pointLoadFailed: '지점 데이터를 불러오지 못했습니다.',
      completionNeedsReadiness: '현재 지점을 완료하려면 답변 있는 질문 1개 이상과 자기평가를 먼저 저장해야 합니다.',
      researchContentRequired: '연구지점은 학습 콘텐츠를 확정해야 학습을 시작하거나 완료할 수 있습니다.',
      learningPointNotFound: '현재 learning point를 찾지 못했습니다.',
      pointStatusSaveFailed: '지점 상태 저장에 실패했습니다.',
      pointCompleted: '현재 지점을 완료로 기록했습니다.',
      pointCompletionCancelled: '학습 지점 완료를 취소했습니다.',
      pointInProgress: '현재 지점을 진행 중으로 기록했습니다.',
      pointGoalSaveFailed: '탐험목표 저장에 실패했습니다.',
      pointGoalSaved: '탐험목표를 저장했습니다.',
      pointRecordSaveFailed: '탐험기록 저장에 실패했습니다.',
      pointRecordSaved: '탐험기록을 저장했습니다.',
      journalSaveFailed: '탐험일지 저장에 실패했습니다.',
      journalSaved: '탐험일지를 저장했습니다.',
      practiceCreateFailed: '연습/활동기록 추가에 실패했습니다.',
      practiceCreated: '연습/활동기록을 추가했습니다.',
      practiceUpdateFailed: '연습/활동기록 수정에 실패했습니다.',
      practiceUpdated: '연습/활동기록을 수정했습니다.',
      practiceDeleteFailed: '연습/활동기록 삭제에 실패했습니다.',
      practiceDeleted: '연습/활동기록을 삭제했습니다.',
      artifactFallbackTitles: {
        video: '영상자료',
        attachment: '첨부자료',
        document: '문서 편집자료',
      },
      artifactSaveFailed: '탐험결과물 저장에 실패했습니다.',
      artifactCreated: '탐험결과물을 추가했습니다.',
      artifactUpdateFailed: '탐험결과물 수정에 실패했습니다.',
      artifactUpdated: '탐험결과물을 수정했습니다.',
      artifactDeleteFailed: '탐험결과물 삭제에 실패했습니다.',
      artifactDeleted: '탐험결과물을 삭제했습니다.',
      artifactSubtitleDuplicate: '결과물 자막은 1개만 업로드할 수 있습니다.',
      artifactThumbnailDuplicate: '결과물 썸네일은 1개만 업로드할 수 있습니다.',
      artifactVideoDuplicate: '결과물 동영상은 1개만 업로드할 수 있습니다.',
      artifactAttachmentLimit: '결과물 첨부파일은 최대 10개까지 업로드할 수 있습니다.',
      artifactUploadFailedStatus: (status) => `결과물 파일 업로드에 실패했습니다. (HTTP ${status})`,
      artifactUploadFailed: '결과물 파일 업로드에 실패했습니다.',
      artifactUploadSuccess: '파일을 결과물에 업로드했습니다.',
      attachmentCreateFailed: '보조자료 추가에 실패했습니다.',
      attachmentCreated: '보조자료를 추가했습니다.',
      attachmentLinkInvalid: '링크는 http:// 또는 https:// 주소로 추가해 주세요.',
      replacementAttachmentCreateFailed: '대체 자료 추가에 실패했습니다.',
      attachmentUploadFailed: '파일 업로드에 실패했습니다.',
      attachmentUploadSuccess: '파일을 보조자료로 업로드했습니다.',
      attachmentReplaceFailed: '보조자료 파일 교체에 실패했습니다.',
      attachmentReplaced: '보조자료 파일을 교체했습니다.',
      attachmentOpenFailed: '보조자료 URL 생성에 실패했습니다.',
      attachmentUpdateFailed: '보조자료 수정에 실패했습니다.',
      attachmentUpdated: '보조자료를 수정했습니다.',
      attachmentDeleteLocked: '확정된 학습 콘텐츠의 첨부파일은 학습하기 흐름을 유지하기 위해 삭제할 수 없습니다.',
      attachmentDeleteFailed: '보조자료 삭제에 실패했습니다.',
      questionSaveFailed: '질문 저장에 실패했습니다.',
      questionCreated: '질문을 등록했습니다. 이제 등록된 질문에 대한 답을 찾아 기록해 주세요.',
      questionAnswerSaved: '질문 답변을 저장했습니다.',
      questionAnswerSaveFailed: '질문 답변 저장에 실패했습니다.',
      questionDeleteFailed: '질문 삭제에 실패했습니다.',
      questionDeleted: '질문을 삭제했습니다.',
      aiFeedbackUnavailable: '현재 사용할 수 있는 LLM이 없습니다. AI 설정을 확인해 주세요.',
      aiFeedbackFailed: 'AI 코칭 생성에 실패했습니다.',
      aiFeedbackCreated: 'AI 코칭을 생성했습니다. 저장해야 이 코칭이 기록됩니다.',
      selfEvalDraftRequired: '적용문제 답변을 바탕으로 AI 평가 참고를 먼저 진행해 주세요.',
      selfEvalSaveFailed: '자기평가 저장에 실패했습니다.',
      selfEvalSaved: '자기평가를 저장했습니다.',
      selfEvalDraftFailed: 'AI 자기평가 초안 생성에 실패했습니다.',
      selfEvalSavedNotice: 'AI 평가 참고를 저장했습니다. 다시 평가하려면 다시 평가 하기를 눌러 주세요.',
      selfEvalQuestionsReady: '적용문제 3개를 준비했습니다. 답변 후 AI 평가 참고를 진행해 주세요.',
    },
    explorationContent: {
      dateLocale: 'ko-KR',
      reportTypeLabels: {
        broken_link: '링크가 열리지 않음',
        wrong_content: '내용이 맞지 않음',
        unsafe_content: '부적절한 내용',
        copyright: '저작권 우려',
        low_quality: '품질이 낮음',
        other: '기타',
      },
      reportStatusLabels: {
        open: '접수됨',
        reviewing: '검토 중',
        resolved: '해결됨',
        dismissed: '반려됨',
        cancelled: '취소됨',
      },
      aiSummaryTitle: 'AI 지점 설명',
      aiSummaryGenerating: '✨ 생성 중...',
      aiSummaryRegenerate: '✨ AI 지점 설명 다시 생성',
      aiSummaryGenerate: '✨ AI 지점 설명 생성',
      aiSummaryEmpty: '아직 생성된 AI 지점 설명이 없습니다. 원문 링크와 지점 정보를 바탕으로 이 지점에서 볼 핵심을 정리할 수 있습니다.',
      providedContentTitle: '제공된 학습 콘텐츠',
      openCurrentContent: '🔗 현재 학습 콘텐츠 열기',
      youtubeSourceLabel: 'YouTube 출처',
      openYoutubePlayer: 'YouTube 공식 플레이어로 보기',
      openOriginalContent: '원본 링크 열기',
      sourceLayoutHint: 'AI 결제 혜택은 영상 시청이 아니라 AI 피드백, 학습 기록 분석, 복습 루틴, 결과물 코칭, 개인 학습 리포트에만 연결됩니다.',
      emptyExternalLink: '이 탐험지점에는 아직 외부 링크가 없습니다.',
      replacementCompletedNotice: '신고한 원문 콘텐츠를 대체 콘텐츠로 교체했습니다. 교체 전 자료는 신고 이력에 보존됩니다.',
      issueSummary: '학습 콘텐츠 오류 신고 / 교체',
      issueHelp: '현재 학습 콘텐츠를 바꾸는 작업입니다. 링크가 열리지 않거나 내용이 맞지 않을 때만 사용하세요.',
      reportTitle: '학습 콘텐츠 오류 신고',
      reportHelp: '신고 후 대체 콘텐츠 추천과 교체가 열립니다',
      reportPlaceholder: '어떤 문제가 있었는지 5자 이상으로 적어 주세요.',
      reportSubmitting: '🚩 접수 중...',
      reportAction: '🚩 오류 신고',
      reportGateHelp: '오류 신고를 접수하면 대체 콘텐츠 추천과 직접 교체가 열립니다.',
      replacementGateNotice: '접수된 신고가 있습니다. 먼저 대체 콘텐츠를 추천받거나 직접 교체를 완료해 주세요.',
      cancelReport: '신고 취소하고 계속 학습',
      replacementLoading: '추천받는 중...',
      loadReplacement: '대체 콘텐츠 추천받기',
      replacementCandidatesAria: '대체 자료 후보',
      candidateColumn: '자료',
      replaceColumn: '교체',
      untitledCandidate: '제목 없는 자료',
      replace: '교체',
      directReplacementTitle: '직접 찾은 콘텐츠로 교체',
      directReplacementHelp: '신고한 원문 콘텐츠 대체',
      directReplacementTitlePlaceholder: '자료 제목',
      directReplacementAction: '직접 교체',
      myReportsTitle: '내 신고 내역',
      reportCount: (count) => `${count}건`,
      reportHistoryAria: '내 학습 콘텐츠 오류 신고 내역',
      reportColumn: '신고',
      statusColumn: '상태',
      receivedAtColumn: '접수일',
      replacedPrefix: (value) => `교체됨: ${value}`,
      emptyReports: '아직 신고한 내용이 없습니다.',
      cancelReportSuccessNotice: '신고를 취소했습니다. 현재 학습 콘텐츠로 계속 진행합니다.',
      replacementQueryMissing: '대체 자료를 찾을 기준이 부족합니다.',
      insufficientPoints: '포인트가 부족해 대체 자료를 추천받지 못했습니다.',
      replacementLoadFailed: '대체 자료 추천에 실패했습니다.',
      replacementFound: '대체 후보를 찾았습니다. 검토 후 이 콘텐츠로 교체할 수 있습니다.',
      replacementEmpty: '조건에 맞는 대체 후보를 찾지 못했습니다. 직접 찾은 링크로 교체할 수 있습니다.',
      replacementUrlInvalid: '링크는 http:// 또는 https:// 주소로 입력해 주세요.',
      replacementFallbackTitle: '대체 학습 콘텐츠',
      replacementFailed: '학습 콘텐츠 교체에 실패했습니다.',
      replacementSuccess: '현재 학습 콘텐츠를 선택한 자료로 교체했습니다.',
      replacementEntrySuccess: '학습 콘텐츠를 교체했습니다. 새 원문 콘텐츠를 열어 학습을 이어가세요.',
      candidateUrlMissing: '선택한 후보에 사용할 URL이 없습니다.',
      cancelReportFailed: '오류 신고 취소에 실패했습니다.',
      cancelReportEntrySuccess: '오류 신고를 취소했습니다. 현재 학습 콘텐츠로 계속 진행할 수 있습니다.',
      reportFailed: '학습 콘텐츠 오류 신고에 실패했습니다.',
      reportSuccess: '학습 콘텐츠 오류 신고를 접수했습니다.',
      aiUnavailable: '현재 사용할 수 있는 LLM이 없습니다. AI 설정을 확인해 주세요.',
      aiSummaryFailed: 'AI 지점 설명 생성에 실패했습니다.',
      aiSummarySuccess: 'AI 지점 설명을 생성했습니다.',
    },
    editor: {
      linkHttpsOnly: 'https://로 시작하는 안전한 링크만 입력할 수 있습니다.',
      imageHttpsOnly: 'https://로 시작하는 이미지 URL만 입력할 수 있습니다.',
      imageUploadFailed: '이미지 업로드에 실패했습니다.',
      paragraphStyle: '문단 스타일',
      paragraph: '본문',
      heading1: '제목 1',
      heading2: '제목 2',
      heading3: '제목 3',
      undo: '되돌리기',
      redo: '다시 실행',
      quote: '인용',
      codeBlock: '코드 블록',
      bold: '굵게',
      italic: '기울임',
      underline: '밑줄',
      strike: '취소선',
      highlight: '하이라이트',
      inlineCode: '인라인 코드',
      superscript: '위첨자',
      subscript: '아래첨자',
      bulletList: '글머리 목록',
      orderedList: '번호 목록',
      alignLeft: '왼쪽 정렬',
      alignCenter: '가운데 정렬',
      alignRight: '오른쪽 정렬',
      link: '링크',
      image: '이미지',
      insertTable: '표 삽입',
      addColumn: '열 추가',
      addRow: '행 추가',
      deleteColumn: '열 삭제',
      deleteRow: '행 삭제',
      deleteTable: '표 삭제',
      apply: '적용',
      unset: '해제',
      imageUrlPlaceholder: 'https://... 이미지 URL',
      imageAltPlaceholder: '이미지 설명',
      insertImage: '삽입',
      chooseFile: '파일 선택',
      uploadHint: '업로드 가능: jpg, jpeg, png, webp, gif / 5MB 이하',
      uploading: '업로드 중...',
      uploadInsert: '업로드 삽입',
      inlineImageTypeError: 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.',
      inlineImageSizeError: '이미지는 5MB 이하만 업로드할 수 있습니다.',
    },
    researchMaterial: {
      pageTitle: '연구 자료 만들기',
      nextActionNotice: '영상, 본문, 파일 중 하나를 선택해 연구자료를 만들어 주세요.',
      blockLimitMessage: '현재 지점에서 이어 쓸 수 있는 게시글을 모두 사용했습니다. 작성한 내용을 저장하고 학습하기로 이어가 주세요.',
      blockAddUnavailable: '게시글을 추가할 수 없습니다.',
      blockSaveFailed: '연구 블록 저장에 실패했습니다.',
      blockAddSuccess: '연구 블록을 추가했습니다.',
      blockSaveSuccess: '연구 블록을 저장했습니다.',
      deleteLockedMessage: '확정된 학습 콘텐츠는 수정하거나 삭제할 수 없습니다. 수정이 필요하면 먼저 미확정으로 변경해 주세요.',
      blockDeleteFailed: '연구 블록 삭제에 실패했습니다.',
      confirmUnavailable: '저장된 학습 콘텐츠가 있어야 확정할 수 있습니다.',
      confirmFailed: '학습 콘텐츠 확정에 실패했습니다.',
      confirmSuccess: '학습 콘텐츠를 확정했습니다. 이제 이 연구지점 학습을 시작할 수 있습니다.',
      attachmentLimitMessage: (count) => `학습 콘텐츠 첨부는 최대 ${count}개까지 추가할 수 있습니다.`,
      attachmentDirtyMessage: '수정 중인 학습 콘텐츠 첨부를 먼저 저장해야 다음 첨부를 추가할 수 있습니다.',
      attachmentAddUnavailable: '학습 콘텐츠 첨부를 추가할 수 없습니다.',
      attachmentAddFailed: '학습 콘텐츠 첨부에 실패했습니다.',
      attachmentAddSuccess: '학습 콘텐츠 첨부를 추가했습니다.',
      subtitleDuplicate: '학습 콘텐츠 자막은 1개만 업로드할 수 있습니다.',
      thumbnailDuplicate: '학습 콘텐츠 썸네일은 1개만 업로드할 수 있습니다.',
      videoDuplicate: '학습 콘텐츠 동영상은 1개만 업로드할 수 있습니다.',
      referenceTypeError: '학습 콘텐츠는 아래 파일만 추가할 수 있습니다.\n이미지: jpg, jpeg, png, webp, gif\n동영상/자막: mp4, webm, mov, vtt, srt\n문서/zip: pdf, doc, docx, hwp, hwpx, txt, xls, xlsx, csv, ppt, pptx, rtf, odt, md, zip',
      imageSizeError: '선택한 이미지는 5MB를 넘었습니다.\n이미지는 파일당 5MB 이하만 업로드할 수 있습니다.',
      videoTypeError: '동영상은 mp4, webm, mov 파일만 업로드할 수 있습니다.',
      videoSizeError: '동영상은 500MB 이하만 업로드할 수 있습니다.',
      videoDurationError: '동영상은 25분 이하만 업로드할 수 있습니다.',
      videoMetadataError: '동영상 정보를 읽지 못했습니다. mp4, webm, mov 파일인지 확인해 주세요.',
      subtitleTypeError: '자막은 vtt 또는 srt 파일만 업로드할 수 있습니다.',
      subtitleSizeError: '자막은 2MB 이하만 업로드할 수 있습니다.',
      thumbnailTypeError: '썸네일은 jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.',
      thumbnailSizeError: '썸네일은 5MB 이하만 업로드할 수 있습니다.',
      documentSizeError: '선택한 문서/zip은 20MB를 넘었습니다.\n문서와 zip은 파일당 20MB 이하만 업로드할 수 있습니다.',
      uploadTooLargeServer: '파일이 앱 기준은 통과했지만 서버 업로드 허용 크기를 넘었습니다. 관리자에게 nginx 업로드 제한 설정 확인을 요청해 주세요.',
      uploadFailedStatus: (status) => `학습 콘텐츠 파일 업로드에 실패했습니다. (HTTP ${status})`,
      fileUploadFailed: '학습 콘텐츠 파일 업로드에 실패했습니다.',
      fileUploadSuccess: '파일을 학습 콘텐츠로 업로드했습니다.',
      inlineImageUnavailable: '이미지를 업로드할 수 없습니다.',
      inlineImageTypeError: 'jpg, jpeg, png, webp, gif 이미지만 업로드할 수 있습니다.',
      inlineImageSizeError: '이미지는 5MB 이하만 업로드할 수 있습니다.',
      inlineImageUploadFailed: '이미지 업로드에 실패했습니다.',
      inlineImageInserted: '이미지를 학습 자료 본문에 삽입했습니다.',
      statusConfirmed: '✅ 확정됨',
      statusDraft: '미확정',
      confirmSaving: '확정 중...',
      confirmDone: '✅ 학습 콘텐츠 확정됨',
      confirmAction: '학습 콘텐츠 확정',
      confirmModalTitle: '학습 콘텐츠를 확정할까요?',
      confirmModalMessage: '확정하면 이 연구자료가 학습하기 흐름에 사용되고, 이후 수정하려면 먼저 미확정으로 변경해야 합니다.',
      unconfirmAction: '미확정으로 변경',
      unconfirmModalTitle: '미확정으로 변경할까요?',
      unconfirmModalMessage: '미확정으로 변경하면 연구자료를 다시 수정할 수 있습니다. 수정 후에는 다시 학습 콘텐츠 확정이 필요합니다.',
      unconfirmSaving: '변경 중...',
      unconfirmFailed: '학습 콘텐츠 미확정 전환에 실패했습니다.',
      unconfirmSuccess: '학습 콘텐츠를 미확정 상태로 변경했습니다.',
      confirmedHelp: '확정된 학습 콘텐츠는 학습하기 흐름에 사용됩니다. 수정이 필요하면 먼저 미확정으로 변경해 주세요.',
      confirmReadyHelp: '저장한 학습 콘텐츠를 확정하면 아래 학습하기 흐름에 사용됩니다.',
      confirmBlockedHelp: '저장된 학습 콘텐츠가 있어야 확정할 수 있습니다.',
      attachmentFallback: '첨부파일',
      videoUrlFailed: '동영상 URL 생성에 실패했습니다.',
      contentFallback: '학습 콘텐츠',
      attachmentsTitle: '첨부파일',
      uploadedFileListTitle: '올라간 파일',
      openFileLink: '파일 열기',
      attachmentLimitTitle: '첨부파일은 최대 10개까지 추가할 수 있습니다.',
      attachmentAddTitle: '첨부파일 추가',
      uploadingAttachment: '📎 업로드 중...',
      addAttachment: '파일',
      attachmentDeleteTitle: '첨부파일 삭제',
      delete: '🗑️ 삭제',
      attachmentImageAlt: '첨부 이미지',
      emptyAttachments: '첨부된 파일이 없습니다.',
      attachmentHelp: '첨부 가능: 이미지 5MB 이하, 문서/zip 20MB 이하, 최대 10개',
      videoSectionTitle: '학습 콘텐츠 영상',
      uploadedVideoListTitle: '올라간 영상',
      openVideoLink: '영상 열기',
      videoTitleLabel: '영상 제목',
      videoTitlePlaceholder: '영상 제목',
      saving: '저장 중...',
      saveTitle: '제목 저장',
      videoPlayerTitle: '학습 영상',
      uploadSubtitle: '자막 올리기',
      uploadThumbnail: '썸네일 올리기',
      deleteSubtitle: '자막 삭제',
      deleteThumbnail: '썸네일 삭제',
      deleteVideo: '🗑️ 동영상 삭제',
      uploadingVideo: '🎬 업로드 중...',
      uploadVideo: '영상',
      videoHelp: 'mp4, webm, mov / 500MB 이하 / 25분 이하',
      noVideo: '등록된 영상이 없습니다.',
      contentSectionTitle: '학습 콘텐츠',
      addContent: '본문',
      contentTitleLabel: '콘텐츠 제목',
      contentTitlePlaceholder: '콘텐츠 제목',
      contentPlaceholder: '학습할 콘텐츠의 원문, 요약, 관찰한 내용을 입력하세요',
      contentDisplayTitle: (title) => `제목: ${title}`,
      emptyContent: '작성된 내용이 없습니다.',
      deleteContentTitle: '학습 콘텐츠 삭제',
      deleteContentMessage: (title) => `${title} 콘텐츠를 삭제하시겠습니까?`,
      cancel: '↩️ 취소',
      deleting: '🗑️ 삭제 중...',
      deleteAttachmentModalTitle: '첨부파일 삭제',
      deleteAttachmentMessage: (title) => `${title} 파일을 삭제하시겠습니까?`,
      validationTitle: '첨부파일을 확인해 주세요',
      validationConfirm: '✅ 확인',
      uploadProgressTitle: '학습 콘텐츠 업로드 중',
      uploadProgressMessage: '선택한 파일을 학습 콘텐츠에 추가하고 있습니다. 동영상은 업로드 용량에 따라 시간이 걸릴 수 있습니다.',
      uploadProgressAria: (progress) => `업로드 진행률 ${progress}%`,
      copyrightNotice: '직접 만든 자료나 사용 권한이 있는 자료만 올려 주세요. 다른 사람의 자료는 허용된 범위에서만 사용하세요.',
      copyrightAgreementLinkLabel: '관련 약관 보기',
      statusLabel: '학습 콘텐츠 상태',
      observeDraftTitle: '작성할 학습 콘텐츠가 필요합니다.',
      observeDraftDescription: '관찰 단계에서는 작성 도구를 바로 열지 않습니다. 연구 자료를 작성하거나 수정하려면 편집으로 들어가 콘텐츠를 저장하고 확정하세요.',
      openEditMode: '✏️ 연구자료 편집',
      navAria: '학습 콘텐츠 하위세션 이동',
      navVideo: '영상',
      navContent: '본문',
      navAttachments: '파일',
    },
    learningWork: {
      title: '학습하기',
      subtitle: '내용정리로 학습을 시작한 뒤, 지점 완료 전에는 질문 답변 1개와 자기평가를 저장해야 합니다. 연습/활동기록과 결과물은 필요할 때 남기면 됩니다.',
      lockedNotice: '연구지점은 학습 콘텐츠 확정 전까지 학습하기를 시작하지 않습니다. 먼저 자료를 저장하고 학습 콘텐츠 확정을 눌러 주세요.',
      miniNavLabel: '학습하기 하위세션 이동',
      notesLabel: '내용정리',
      recordsLabel: '학습기록',
      recordsTitle: '학습 기록',
      optionalLabel: '(선택)',
      recordsDescription: '지점 완료에는 질문 답변 1개가 필요합니다. 연습/활동기록과 결과물은 직접 해본 활동이나 남길 성과가 있을 때 필요한 만큼만 기록하세요.',
      recordsLockedNotice: '내용정리를 하나 저장한 뒤 질문 답변 1개와 자기평가까지 마치면 지점을 완료할 수 있습니다. 연습/활동기록과 결과물은 필요할 때 남기면 됩니다.',
      moveModalTitle: '학습 기록 이동',
      moveModalEyebrow: 'Lumi Confirm',
      moveModalMessage: '저장하지 않은 변경이 있다면 사라질 수 있습니다. 다른 학습 기록으로 이동할까요?',
      keepWriting: '↩️ 계속 작성',
      move: '이동',
      sourceContentGuide: {
        title: 'Learning Notes - 내용정리',
        message: '원문 콘텐츠를 열었어요. 지금 떠오른 핵심을 한 줄만 내용정리에 남겨보세요.',
      },
      notes: {
        title: '내용정리',
        edit: '✏️ 편집',
        helper: '모두 채울 필요는 없습니다. 지금 필요한 항목 하나만 정리해도 다음 기록으로 이어갈 수 있습니다.',
        cancelModalTitle: '내용정리 편집 취소',
        cancelModalEyebrow: 'Lumi Confirm',
        cancelModalMessage: '변경한 내용정리를 저장하지 않고 취소하시겠습니까?',
        back: '↩️ 돌아가기',
        cancel: '🚫 취소',
        cancelButton: '↩️ 취소',
        saving: '💾 저장 중...',
        save: '💾 내용정리 저장',
        progressTitle: '내용정리 진행',
        addEmpty: '비어 있는 항목 추가',
        progressIncomplete: '필요한 항목만 채워도 됩니다. 더 정리하고 싶다면 비어 있는 항목을 추가해 보세요.',
        progressComplete: '내용정리 항목이 모두 채워졌습니다. 필요하면 편집으로 다듬을 수 있습니다.',
        empty: '아직 저장된 내용정리가 없습니다.',
        fields: {
          observation: { icon: '💡', label: '핵심 개념', placeholder: '핵심 개념을 정리하세요.' },
          reflection: { icon: '🗣️', label: '내 설명', placeholder: '내 말로 설명해 보세요.' },
          examples: { icon: '🔎', label: '예시', placeholder: '예시나 적용 상황을 적어 보세요.' },
          nextStep: { icon: '❔', label: '헷갈린 부분', placeholder: '헷갈린 부분이나 다음에 확인할 내용을 적어 보세요.' },
        },
      },
      recordPanels: {
        common: {
          number: '번호',
          title: '제목',
          status: '상태',
          type: '구분',
          previous: '← 이전',
          next: '다음 →',
          countRange: (total, start, end) => `${total}개 중 ${start}-${end}개`,
        },
        questions: {
          title: '질문과 답변',
          subtitle: '',
          add: '➕ 질문',
          listAria: '질문 리스트',
          answered: '답변 기록됨',
          pending: '답변 대기',
          empty: '아직 저장된 질문이 없습니다.',
          fallbackTitle: '질문 없음',
          typeLabels: { reflection: '이해 안 됨', application: '적용 방법', goal_alignment: '더 알고 싶음', none: '구분 없음' },
          detail: {
            editTitle: '질문 수정',
            viewTitle: '질문',
            addTitle: '질문 등록',
            addSubtitle: '궁금한 내용을 한 번에 적어 둡니다',
            type: '구분',
            title: '제목',
            content: '질문 내용',
            answer: '답변',
            answerMeta: '직접 찾은 답과 근거',
            answerEmptyMeta: '한 질문에는 답변을 하나만 등록합니다',
            aiCoach: 'AI 코칭',
            aiCoachMeta: '저장해야 기록됩니다',
            titlePlaceholder: '질문 제목',
            contentPlaceholder: '궁금한 점을 입력하세요.',
            answerPlaceholder: '내가 찾은 답, 확인한 근거, 현재 이해한 내용을 기록하세요.',
            noContent: '질문 내용이 없습니다.',
            noAnswer: '아직 등록된 답변이 없습니다.',
            communityDisabled: '💬 커뮤니티에 묻기 · 정식 오픈 시 이용 가능',
            addAnswer: '✍️ 답변 등록',
            editAnswer: '✏️ 답변 수정',
            createAiCoach: '✨ AI 코칭 생성',
            createAiCoachStarter: '✨ 답변 방향 코칭',
            createAiCoachRevision: '✨ 답변 보완 코칭',
            generatingAiCoach: '✨ 생성 중...',
            saveDetail: '💾 질문 저장',
            savingDetail: '💾 저장 중...',
            saveAnswer: '💾 답변 저장',
            savingAnswer: '💾 저장 중...',
            register: '➕ 질문 등록',
            saving: '💾 저장 중...',
            edit: '✏️ 질문 수정',
            list: '📋 목록',
            delete: '🗑️ 삭제',
            deleting: '🗑️ 삭제 중...',
            cancel: '🚫 취소',
            back: '↩️ 돌아가기',
            ok: '확인',
            answerInstruction: '답변은 질문을 등록한 뒤 리플처럼 하나만 남깁니다.',
            aiCoachNotice: 'AI 코칭은 답변 작성 중에만 사용할 수 있으며, 답변을 저장해야 기록됩니다.',
            aiCoachStarterHint: 'AI가 답변을 시작할 방향을 제안합니다.',
            aiCoachRevisionHint: 'AI가 작성 중인 답변의 빠진 부분과 보완 방향을 제안합니다.',
            aiCoachSummary: '피드백 요약',
            aiCoachNextAction: '다음 행동',
            saveSuccessTitle: '질문 저장',
            saveSuccessMessage: '질문을 저장했습니다.',
            saveFailTitle: '질문 저장 실패',
            saveFailMessage: '질문을 저장하지 못했습니다.',
            answerSaveSuccessTitle: '답변 저장',
            answerSaveSuccessMessage: '답변을 저장했습니다.',
            answerSaveFailTitle: '답변 저장 실패',
            answerSaveFailMessage: '답변을 저장하지 못했습니다.',
            deleteModalTitle: '질문 삭제',
            deleteModalMessage: (title) => `'${title}' 질문을 삭제하시겠습니까?`,
            listReturnTitle: '질문 상세 닫기',
            listReturnMessage: '편집 중인 내용을 저장하지 않고 목록으로 돌아가시겠습니까?',
            addCancelTitle: '질문 등록 취소',
            addCancelMessage: '작성 중인 질문을 저장하지 않고 취소하시겠습니까?',
          },
        },
        practice: {
          title: '연습/활동기록',
          subtitle: '제목을 선택해 상세 기록을 확인합니다',
          add: '➕ 새 연습/활동',
          listAria: '연습/활동기록 리스트',
          hasRecord: '기록 있음',
          waiting: '기록 대기',
          empty: '아직 저장된 연습/활동기록이 없습니다.',
          fallbackTitle: (index) => `연습/활동기록 ${index}`,
          detail: {
            editTitle: '연습/활동기록 편집',
            viewTitle: '연습/활동기록 보기',
            addTitle: '새 연습/활동기록 추가',
            title: '제목',
            duration: '시도 시간(분)',
            achievement: '성과',
            reflection: '막힌 부분/다음 연습기록',
            titlePlaceholder: '연습/시도 제목 예: 통계 개념 적용, 설문조사 시도, 자료 분석 연습',
            durationPlaceholder: '분 단위로 입력해 주세요',
            achievementPlaceholder: '성과',
            reflectionPlaceholder: '막힌 부분과 다음 연습 계획을 함께 기록하세요.',
            addReflection: '➕ 추가 기록(막힌부분/다음연습기록)',
            noDetail: '아직 상세 기록이 없습니다.',
            list: '📋 목록',
            delete: '🗑️ 삭제',
            deleting: '🗑️ 삭제 중...',
            save: '저장',
            saving: '저장 중...',
            edit: '✏️ 편집',
            cancel: '↩️ 취소',
            back: '↩️ 돌아가기',
            deleteModalTitle: '연습/활동기록 삭제',
            deleteModalMessage: (title) => `'${title}' 항목을 삭제하시겠습니까?`,
            listReturnTitle: '목록으로',
            listReturnMessage: '편집 중인 내용을 저장하지 않고 목록으로 돌아가시겠습니까?',
            addCancelTitle: '연습기록 추가 취소',
            addCancelMessage: '작성 중인 내용을 저장하지 않고 취소하시겠습니까?',
          },
        },
        artifacts: {
          title: '결과물제출',
          subtitle: '제목을 선택해 상세 결과물을 확인합니다',
          add: '➕ 결과물',
          empty: '아직 저장된 결과물이 없습니다.',
          listAria: '결과물 리스트',
          fallbackTitle: (index) => `결과물 ${index}`,
          kindLabels: { 영상자료: '영상자료', '문서 편집자료': '문서 편집자료', 첨부자료: '첨부자료' },
          detail: {
            editTitle: '결과물 편집',
            viewTitle: '결과물 보기',
            addTitle: '결과물 작성',
            kind: '자료 구분',
            title: '제목',
            editTitleLabel: '제목 수정',
            videoSection: '결과물 영상',
            attachmentSection: '결과물 첨부',
            documentSection: '문서 편집자료',
            videoPlayerTitle: '결과물 영상',
            titlePlaceholder: '결과물 제목',
            documentTitlePlaceholder: '문서 편집자료 제목',
            documentPlaceholder: '결과물 문서를 입력하세요',
            attachmentFallback: '첨부파일',
            noVideo: '등록된 영상이 없습니다.',
            noFiles: '첨부된 파일이 없습니다.',
            uploadVideo: '🎬 동영상 올리기',
            uploadingVideo: '🎬 업로드 중...',
            uploadSubtitle: '자막 올리기',
            uploadThumbnail: '썸네일 올리기',
            deleteSubtitle: '자막 삭제',
            deleteThumbnail: '썸네일 삭제',
            deleteVideo: '🗑️ 동영상 삭제',
            uploadFile: '➕ 첨부파일',
            uploadingFile: '📎 업로드 중...',
            list: '📋 목록',
            delete: '🗑️ 삭제',
            deleting: '🗑️ 삭제 중...',
            save: '💾 저장',
            saving: '💾 저장 중...',
            edit: '✏️ 편집',
            cancel: '↩️ 취소',
            back: '↩️ 돌아가기',
            videoHelp: 'mp4, webm, mov / 500MB 이하 / 25분 이하',
            deleteModalTitle: '결과물 삭제',
            deleteModalMessage: (title) => `'${title}' 항목을 삭제하시겠습니까?`,
            deleteAttachmentTitle: '첨부 삭제',
            deleteAttachmentMessage: (title) => `${title} 파일을 삭제하시겠습니까?`,
            listReturnTitle: '목록으로',
            listReturnMessage: '편집 중인 내용을 저장하지 않고 목록으로 돌아가시겠습니까?',
            addCancelTitle: '결과물 추가 취소',
            addCancelMessage: '작성 중인 내용을 취소하시겠습니까?',
          },
        },
        attachments: {
          title: '보조자료',
          subtitle: (researchCount, workCount) => `학습 콘텐츠 ${researchCount}개 / 보조자료 ${workCount}개`,
          add: '➕ 보조자료 추가',
          listAria: '보조자료 리스트',
          learningContent: '학습 콘텐츠',
          reference: '보조자료',
          empty: '아직 저장된 보조자료가 없습니다.',
          fallbackTitle: (index) => `보조자료 ${index}`,
          detail: {
            editTitle: '보조자료 편집',
            viewTitle: '보조자료 보기',
            addTitle: '새 보조자료 추가',
            sectionTitle: '보조자료 첨부',
            title: '제목',
            editTitleLabel: '제목 수정',
            file: '파일',
            titlePlaceholder: '자료 제목',
            attachmentFallback: '첨부파일',
            replacementPending: '교체 예정',
            selectedFile: (name) => `선택됨: ${name}`,
            researchMaterialNotice: '학습 콘텐츠 입력 영역에서 수정할 수 있습니다.',
            uploadFile: '📎 파일 업로드',
            uploadingFile: '📎 업로드 중...',
            uploadSelected: '📎 선택한 파일 업로드',
            list: '📋 목록',
            delete: '🗑️ 삭제',
            deleting: '🗑️ 삭제 중...',
            save: '💾 저장',
            saving: '💾 저장 중...',
            edit: '✏️ 편집',
            cancel: '🚫 취소',
            back: '↩️ 돌아가기',
            deleteModalTitle: '보조자료 삭제',
            deleteModalMessage: (title) => `'${title}' 파일을 삭제하시겠습니까?`,
            listReturnTitle: '목록으로',
            listReturnMessage: '편집 중인 내용을 저장하지 않고 목록으로 돌아가시겠습니까?',
            addCancelTitle: '보조자료 추가 취소',
            addCancelMessage: '추가를 취소하시겠습니까?',
          },
        },
        timeline: {
          title: '활동 기록',
          subtitle: '자동 생성 타임라인',
          listAria: '활동 기록 리스트',
          time: '시간',
          activity: '활동 내용',
          empty: '아직 기록된 활동이 없습니다.',
        },
      },
      tabs: {
        questions: { label: '질문관리', title: 'Questions - 질문관리', message: '여기서는 이해가 막힌 부분이나 더 확인하고 싶은 점을 질문으로 남기면 됩니다. 질문을 저장한 뒤 나의 답변을 적고, 필요하면 AI 참고 피드백도 받아보세요.' },
        practice: { label: '연습/활동기록', title: 'Practice & Activity Log - 연습/활동기록', message: '여기서는 실제로 해 본 연습, 조사, 분석, 적용 시도를 남기면 됩니다. 걸린 시간과 막힌 부분, 다음 계획을 짧게 적어도 충분합니다.' },
        artifacts: { label: '결과물제출', title: 'Artifacts - 결과물', message: '여기서는 만든 산출물이나 제출할 링크를 남기면 됩니다. 무엇을 만들었는지, 만들면서 배운 점과 어려웠던 점을 함께 정리해 보세요.' },
        attachments: { label: '보조자료', title: 'References - 보조자료', message: '여기서는 학습 중 추가로 참고한 링크나 파일을 모아 두면 됩니다. 내용정리와 섞지 않고 자료 출처를 따로 챙길 때 사용하면 좋습니다.' },
        timeline: { label: '활동 기록', title: 'Activity Log - 활동 기록', message: '여기서는 이 지점에서 저장한 주요 행동을 시간순으로 볼 수 있습니다. 내가 어떤 순서로 학습했는지 천천히 되짚어보세요.' },
      },
      noteGuides: {
        observation: { title: 'Core Concept - 핵심 개념', message: '여기서는 지금 학습한 내용에서 가장 중요한 한두 문장을 남기면 됩니다. 길게 요약하기보다 나중에 다시 봐도 중심이 보이게 적어 보세요.' },
        reflection: { title: 'My Explanation - 내 설명', message: '여기서는 학습한 내용을 나의 말로 다시 풀어 쓰면 됩니다. 그대로 베끼기보다 누군가에게 설명하듯 적으면 이해가 더 분명해집니다.' },
        examples: { title: 'Examples - 예시', message: '여기서는 개념이 실제로 쓰이는 상황이나 내가 떠올린 사례를 남겨보세요. 적용 장면을 적어 두면 나중에 다시 이해하기 쉽습니다.' },
        nextStep: { title: 'Confusing Parts - 헷갈린 부분', message: '여기서는 아직 확실하지 않은 내용과 다음에 확인할 질문을 남기면 됩니다. 완벽히 정리되지 않은 상태를 그대로 적어도 괜찮습니다.' },
      },
    },
  },
  completion: {
    encouragement: '완벽하지 않아도 괜찮아요. 지금 기준으로 이 탐험지점을 충분히 학습했다고 판단하면 완료 처리하세요.',
    finalAction: '마지막 액션',
    completedLocked: '완료된 지점은 수정할 수 없습니다. 수정이 필요하면 완료를 먼저 취소하세요.',
    readyInstruction: '준비 상태와 체크사항을 확인한 뒤 현재 지점을 완료하세요.',
    cancelProcessing: '취소 처리 중...',
    cancelCompletion: '학습 지점 완료 취소',
    saving: '저장 중...',
    completePoint: '현재 지점 완료하기',
    completeModalTitle: '현재 지점 완료',
    completeModalEyebrow: 'Lumi Confirm',
    completeModalMessage: '이 지점을 완료로 기록할까요? 완료 후에도 기록은 확인할 수 있습니다.',
    close: '닫기',
    cancel: '취소',
    complete: '완료하기',
    cancelModalTitle: '학습 지점 완료 취소',
    cancelModalMessage: '완료 상태를 취소하면 이 지점을 다시 수정할 수 있습니다. 완료를 취소할까요?',
    cancelDone: '완료 취소',
  },
  selfEvaluation: {
    title: '자기평가',
    subtitle: '적용문제 3개로 가볍게 확인합니다',
    savedTitle: '저장된 자기평가',
    applicationTitle: '적용 확인',
    savedDescription: '이 지점의 자기평가가 확정되어 저장되었습니다.',
    applicationDescription: '간단한 적용문제에 답한 뒤 AI 평가 참고를 누르면 자기평가로 저장됩니다.',
    generatingTitle: '자기평가 생성 중',
    generatingDescription: 'AI가 적용문제와 평가 참고 문구를 준비하고 있습니다.',
    questionLabel: '문제',
    answerPlaceholder: '내 답변을 짧게 적어 주세요.',
    aiReferenceTitle: 'AI 평가 참고',
    aiReferenceDescription: '적용문제 답변을 바탕으로 AI가 참고 평가를 정리했습니다. 이 내용은 내부 평가값으로 저장되었습니다.',
    generatingButton: '생성 중...',
    restartButton: '다시 평가 하기',
    aiReferenceButton: 'AI 평가 참고',
    createQuestionsButton: '적용문제 만들기',
    finalTitle: '나의 최종 자기평가',
    finalDescription: 'AI 평가 참고를 확인한 뒤, 지금 내 상태에 가장 가까운 점수를 선택하세요.',
    scoreSuffix: '점',
    lockedUntilDraft: 'AI 평가 참고가 저장된 뒤 선택할 수 있습니다.',
    options: [
      { value: '1', label: '1', description: '아직 어렵다' },
      { value: '2', label: '2', description: '조금 이해했다' },
      { value: '3', label: '3', description: '보통이다' },
      { value: '4', label: '4', description: '잘 이해했다' },
      { value: '5', label: '5', description: '설명할 수 있다' },
    ],
  },
  completionStatus: {
    readOnlyTitle: '읽기 전용 보기',
    readOnlyDescription: 'shared 상태에서는 질문, 자기평가, 기록, 연구 블록을 수정하지 않고 현재 학습 결과만 확인할 수 있습니다.',
    achievementTitle: '탐험지점 완료!',
    achievementDescription: '이번 탐험에서 기록으로 남긴 성과를 확인해 보세요.',
    contentChecked: '✓ 학습 콘텐츠 확인',
    noteCount: (count) => `✓ 내용정리 ${count}개 작성`,
    questionCount: (total, answered) => `✓ 질문 ${total}개 / 답변 ${answered}개`,
    practiceCount: (count) => `✓ 연습/활동기록 ${count}건`,
    artifactCount: (count) => `✓ 결과물 ${count}개`,
    selfEvaluationState: (completed) => `✓ 자기평가 ${completed ? '완료' : '미완료'}`,
    backToDiary: '탐험일지로 돌아가기',
    sharedStatusTitle: '완료 반영 상태',
    readyStatusTitle: '완료 준비 상태',
    researchMaterialConfirmed: (completed) => `• 학습 콘텐츠 확정: ${completed ? '완료' : '미완료'}`,
    answeredQuestionReady: (completed, count) => `• 답변 있는 질문 1개 이상: ${completed ? `완료 (${count}개)` : '미완료'}`,
    selfEvaluationSaved: (completed) => `• 자기평가 저장: ${completed ? '완료' : '미완료'}`,
    selfEvaluationQuality: (completed) => `• 자기평가 내부 기준 반영: ${completed ? '완료' : '미완료'}`,
    goalConnection: (completed) => `• 목표 연결 메모: ${completed ? '완료' : '미완료'}`,
    sharedNotReady: '이 화면은 현재까지 저장된 reflection과 완료 충족 여부를 읽기 전용으로 보여줍니다.',
    researchNotReady: '연구지점은 학습 콘텐츠를 확정한 뒤 질문 답변과 자기평가를 저장해야 완료할 수 있습니다.',
    learningNotReady: '현재 지점 완료를 누르기 전에 질문 답변과 자기평가를 저장해야 합니다.',
    sharedReady: '이 지점은 완료 조건을 충족한 상태로 공유되고 있습니다.',
    learningReady: '현재 지점 완료 조건이 모두 충족되었습니다.',
  },
  lumiGuide: {
    goalLabel: '탐험목표',
    close: '루미 안내 닫기',
    closeButton: '💡 루미안내 ❌',
    featureGuideLabel: '기능 안내',
    currentTaskLabel: '지금 할 일',
    fallbackGoal: '현재 탐험일지에 연결된 학습 목표를 불러오지 못했습니다.',
    pointTitle: (isExplorationPoint) => (isExplorationPoint ? 'Exploration Point - 탐험지점' : 'Research Point - 연구지점'),
    pointMessage: (isExplorationPoint) => (isExplorationPoint
      ? '여기는 준비된 자료를 살펴보고 나의 기록으로 바꾸는 탐험지점입니다. 먼저 흐름을 가볍게 확인해 보세요.'
      : '여기는 내가 준비한 자료를 학습 콘텐츠로 삼아 차근차근 정리하는 연구지점입니다. 자료를 다듬고 확정하면 다음 단계로 이어갈 수 있습니다.'),
    learningContentTitle: 'Learning Content - 학습 콘텐츠',
    learningContentMessage: (isExplorationPoint) => (isExplorationPoint
      ? '여기서는 제공된 학습 콘텐츠를 확인하면 됩니다. 원문과 요약을 살펴보고, 중요한 내용을 나의 기록으로 옮겨보세요.'
      : '여기서는 학습 콘텐츠를 준비하면 됩니다. 게시글과 첨부를 먼저 정리하고, 확정한 뒤 내용정리와 질문으로 이어가세요.'),
    selfEvaluationTitle: 'Self Evaluation - 자기평가',
    selfEvaluationMessage: '마지막에는 지금 이해한 정도와 적용 가능성, 목표와의 연결을 스스로 점검해 보세요. 완벽하지 않아도 괜찮습니다. 지금 기준의 판단을 남기면 됩니다.',
    completedFlow: '이 지점은 완료됐습니다. 남긴 기록과 자기평가를 다시 확인할 수 있습니다.',
    researchMaterialFlow: '지금은 학습 콘텐츠를 게시글이나 첨부로 정리하고 확정해야 합니다.',
    openExternalFlow: '먼저 원문 콘텐츠를 열어 확인해 보세요. 그다음 내용정리에서 필요한 항목 하나만 골라 짧게 저장하면 다음 기록으로 넘어갈 수 있습니다.',
    firstNoteFlow: '먼저 내용정리에서 필요한 항목 하나를 저장해 보세요. 지점을 완료하려면 이후 질문 답변 1개와 자기평가가 필요합니다.',
    selfEvaluationReadyFlow: '학습 기록이 쌓였습니다. 계속해서 필요한 기록을 더 남겨도 좋고, 지금 이해한 만큼 자기평가를 진행해도 좋습니다.',
    noQuestionFlow: '내용정리를 저장했습니다. 지점을 완료하려면 질문을 하나 남기고 답변을 저장해야 합니다. 연습/활동기록이나 결과물은 필요할 때 이어서 남기면 됩니다.',
    noAnswerFlow: '질문을 남겼습니다. 지점을 완료하려면 이 질문에 답변을 저장해야 합니다. 필요하면 연습/활동기록이나 결과물 같은 다른 학습 기록도 이어서 남겨 보세요.',
    completionReadyFlow: '현재 지점을 완료할 준비가 됐습니다. 자기평가 아래의 완료 버튼으로 마무리해 주세요.',
    fillMissingFlow: '하단 완료조건을 보며 아직 비어 있는 항목을 하나씩 채워보세요.',
    restTitle: '🌿 잠깐 쉬어도 괜찮습니다',
    restMessage: '꽤 집중해서 이어왔습니다. 잠깐 쉬었다가 다음에 이어가도 괜찮습니다. 지금 페이스에 맞춰 천천히 진행해도 됩니다.',
    restConfirm: '확인',
  },
}

const en: PointLearningCopy = {
  loading: 'Preparing this point...',
  notFound: 'Point not found.',
  backToPlanet: 'Back to Planet',
  header: {
    diaryMeta: 'Explorer Diary',
    diaryLabel: 'Diary',
    diaryTitle: 'Back to diary',
    lightMode: 'Switch to light mode',
    darkMode: 'Switch to dark mode',
  },
  toolbar: {
    ariaLabel: 'Learning progress menu',
    expand: 'Expand learning toolbar',
    collapse: 'Collapse learning toolbar',
    content: 'Content',
    researchContent: 'Create Learning Content',
    work: 'Learn',
    evaluation: 'Review',
    evaluationDisabled: 'Review unavailable',
    evaluationDisabledTitle: 'Add learning work records before opening review.',
    lumiOpen: 'Open Lumi guide',
    lumiClose: 'Close Lumi guide',
    lumiLabel: 'Lumi Guide',
  },
  skillFlow: {
    progressLabel: 'Point learning flow',
    back: 'Previous step',
    next: 'Next step',
    completePoint: 'Go to completion step',
    stages: {
      point: { label: 'Observe Point', title: 'Observe Point', subtitle: 'Review this point directly before deciding what you want to discover.' },
      goal: { label: 'Set Goal', title: 'Set This Point Goal', subtitle: 'Course, lesson, and point titles are clues. Decide the concrete thing you will make, explain, perform, or solve.' },
      skill_discovery: { label: 'Choose Skills', title: 'Choose Needed Skills', subtitle: 'Choose the skills you will repeat to reach the goal object.' },
      goal_iteration: { label: "Record and Repeat", title: "Record Today Learning", subtitle: "Record today attempts, questions, practice, and artifacts inside this point." },
      final_result: { label: 'Summarize Result', title: 'Summarize This Point Result', subtitle: 'Review the final discovery and what this point produced.' },
      complete: { label: 'Completion Check', title: 'Check Point Completion', subtitle: 'Save self-evaluation and complete this point.' },
    },
    observe: {
      externalTitle: 'Observe external source',
      externalDescription: 'External links open in a new tab or window in the current browser. LearnCosmos keeps the summary and learning record flow here.',
      internalTitle: 'Observe platform content',
      internalDescription: 'Creator content prepared inside the platform can be reviewed directly on this page.',
      researchTitle: 'Observe research material',
      researchDescription: 'Review the confirmed material you wrote for this research point.',
      researchEditHint: 'Create or revise research material in edit mode, then return to observation.',
      externalSpeech: 'Open the source first and see what it contains.',
      internalSpeech: 'Open the source first and see what it contains.',
      researchSpeech: 'Create the research material first.',
      notePanelTitle: 'Observation Notes',
      notePanelDescription: 'Keep adding what you notice while reviewing the source. One saved note unlocks the next step.',
      noteTypeLabel: 'Post type',
      noteContentLabel: 'Content',
      noteTypes: {
        core_summary: { label: 'Key Summary', placeholder: 'Rewrite the important part in your own words.' },
        revisit_part: { label: 'Part to Revisit', placeholder: 'Save a scene, timestamp, sentence, or location you want to review again.' },
        reference_material: { label: 'Supporting Evidence', placeholder: 'Save the link, source name, scene, or explanation that supports your note.' },
      },
      noteAdd: 'Add',
      noteSave: 'Save',
      noteCancel: 'Cancel',
      noteEdit: 'Edit',
      noteDelete: 'Delete',
      noteEmpty: 'No observation notes yet. Start with a short key summary.',
      noteSaved: 'Observation note saved.',
      noteDeleted: 'Observation note deleted.',
      noteSaveFailed: 'Failed to save observation note.',
    },
    goal: {
      contextLabel: 'Point clues',
      titleLabel: 'Goal object I will discover',
      titlePlaceholder: 'Example: short performance video, concept explanation, successful movement record, working code',
      typeLabel: 'Goal object type',
      types: {
        creative: { label: 'Creative artifact', description: 'You will produce a drawing, video, writing, music, or another artifact.' },
        conceptual: { label: 'Concept learning', description: 'You will understand a concept and explain it in your own words.' },
        physical: { label: 'Physical performance', description: 'You will repeat a motion, routine, or performance until it works.' },
        coding: { label: 'Coding implementation', description: 'You will produce working code, a feature, or an executable result.' },
        other: { label: 'Mixed or other', description: 'The goal combines several types or is not clear yet.' },
      },
    },
    skills: {
      prompt: 'Choose the skills you need to repeat for this goal object.',
      presets: {
        creative: ['Observe', 'Compose', 'Draft', 'Revise', 'Polish expression'],
        conceptual: ['Find the core', 'Explain in my words', 'Make examples', 'Ask questions', 'Apply'],
        physical: ['Check posture', 'Repeat attempts', 'Control speed', 'Fix errors', 'Compare records'],
        coding: ['Read requirements', 'Build small', 'Test', 'Debug', 'Refactor'],
        other: ['Observe', 'Try', 'Record', 'Compare', 'Retry'],
      },
      emptyGoalHint: 'Choose a goal object type first to update suggested skills.',
    },
    repeat: {
      requiredHint: 'In Phase 1, attempts and evidence are shown as recommended records, not enforced completion gates.',
      lensLabel: 'Record lens',
      selectedLensLabel: 'Current lens',
      connectedRecordLabel: 'Connected record area',
      selectedSkillsLabel: 'Skills used today',
      noSelectedSkills: 'No skills used today have been selected yet.',
      goalObjectLabel: 'Goal object',
      emptyGoalObject: 'No goal object written yet.',
      savedCountLabel: (count) => `Saved ${count}`,
      lenses: {
        attempt: { label: 'Attempt', description: 'Record what you tried, how many times, and how long it took.' },
        evidence: { label: 'Evidence', description: 'Keep screenshots, writing, code, video, photos, or checklists as evidence.' },
        known: { label: 'Known', description: 'Summarize what you clearly understand or can explain.' },
        question: { label: 'Questions', description: 'Record what you do not know yet and need to check.' },
        answer: { label: 'Answers found', description: 'Record answers and supporting evidence for your questions.' },
        works_well: { label: 'Works well', description: 'Record parts that are becoming stable through repetition.' },
        needs_practice: { label: 'Needs practice', description: 'Record blocks that need more repetition and the next attempt.' },
        resource: { label: 'References', description: 'Collect extra links and files you used.' },
      },
    },
    finalResult: {
      artifactHint: 'Phase 1 does not force a separate representative artifact. Review this through existing artifacts and attachments.',
      openArtifacts: 'Open artifact records',
    },
    complete: {
      readinessHint: 'The actual completion rule stays unchanged: at least one answered question and six saved self-evaluation fields.',
    },
  },
  workspace: {
    materialTitle: 'Content',
    explorationSubtitle: 'Review the provided learning content.',
    researchContentActions: {
      edit: '✏️ Edit',
      saveNew: '💾 Save Content',
      saveChanges: '💾 Save Changes',
      saving: '💾 Saving...',
      saved: '✅ Saved',
    },
    runtime: {
      invalidPointPath: 'Invalid point route.',
      inactiveDiaryOnlyPlanning: 'Inactive explorer diaries can only use the planning page.',
      pointLoadFailed: 'Could not load point data.',
      completionNeedsReadiness: 'Save at least one answered question and your self-evaluation before completing this point.',
      researchContentRequired: 'Research points require confirmed learning content before they can be started or completed.',
      learningPointNotFound: 'Could not find the current learning point.',
      pointStatusSaveFailed: 'Could not save the point status.',
      pointCompleted: 'Point marked as complete.',
      pointCompletionCancelled: 'Point completion cancelled.',
      pointInProgress: 'Point marked as in progress.',
      pointGoalSaveFailed: 'Could not save the point goal.',
      pointGoalSaved: 'Point goal saved.',
      pointRecordSaveFailed: 'Could not save the point record.',
      pointRecordSaved: 'Point record saved.',
      journalSaveFailed: 'Could not save the explorer diary.',
      journalSaved: 'Explorer diary saved.',
      practiceCreateFailed: 'Could not add the practice/activity record.',
      practiceCreated: 'Practice/activity record added.',
      practiceUpdateFailed: 'Could not update the practice/activity record.',
      practiceUpdated: 'Practice/activity record updated.',
      practiceDeleteFailed: 'Could not delete the practice/activity record.',
      practiceDeleted: 'Practice/activity record deleted.',
      artifactFallbackTitles: {
        video: 'Video artifact',
        attachment: 'Attachment artifact',
        document: 'Document artifact',
      },
      artifactSaveFailed: 'Could not save the artifact.',
      artifactCreated: 'Artifact added.',
      artifactUpdateFailed: 'Could not update the artifact.',
      artifactUpdated: 'Artifact updated.',
      artifactDeleteFailed: 'Could not delete the artifact.',
      artifactDeleted: 'Artifact deleted.',
      artifactSubtitleDuplicate: 'Only one subtitle file can be uploaded for an artifact.',
      artifactThumbnailDuplicate: 'Only one thumbnail can be uploaded for an artifact.',
      artifactVideoDuplicate: 'Only one video can be uploaded for an artifact.',
      artifactAttachmentLimit: 'You can upload up to 10 artifact attachment files.',
      artifactUploadFailedStatus: (status) => `Artifact file upload failed. (HTTP ${status})`,
      artifactUploadFailed: 'Artifact file upload failed.',
      artifactUploadSuccess: 'File uploaded to the artifact.',
      attachmentCreateFailed: 'Could not add the reference.',
      attachmentCreated: 'Reference added.',
      attachmentLinkInvalid: 'Add a link that starts with http:// or https://.',
      replacementAttachmentCreateFailed: 'Could not add the replacement reference.',
      attachmentUploadFailed: 'File upload failed.',
      attachmentUploadSuccess: 'File uploaded as a reference.',
      attachmentReplaceFailed: 'Could not replace the reference file.',
      attachmentReplaced: 'Reference file replaced.',
      attachmentOpenFailed: 'Could not create the reference URL.',
      attachmentUpdateFailed: 'Could not update the reference.',
      attachmentUpdated: 'Reference updated.',
      attachmentDeleteLocked: 'Attachments in confirmed learning content cannot be deleted because they are connected to the Learn flow.',
      attachmentDeleteFailed: 'Could not delete the reference.',
      questionSaveFailed: 'Could not save the question.',
      questionCreated: 'Question added. Now find and record your answer.',
      questionAnswerSaved: 'Question answer saved.',
      questionAnswerSaveFailed: 'Could not save the question answer.',
      questionDeleteFailed: 'Could not delete the question.',
      questionDeleted: 'Question deleted.',
      aiFeedbackUnavailable: 'No LLM is currently available. Check your AI settings.',
      aiFeedbackFailed: 'Could not generate AI coaching.',
      aiFeedbackCreated: 'AI coaching generated. Save it to keep this coaching in your record.',
      selfEvalDraftRequired: 'Run AI review from your application answers before saving self-evaluation.',
      selfEvalSaveFailed: 'Could not save the self-evaluation.',
      selfEvalSaved: 'Self-evaluation saved.',
      selfEvalDraftFailed: 'Could not generate the AI self-evaluation draft.',
      selfEvalSavedNotice: 'AI review saved. Use Restart Review if you want to evaluate again.',
      selfEvalQuestionsReady: 'Three application questions are ready. Answer them, then run AI Review.',
    },
    explorationContent: {
      dateLocale: 'en-US',
      reportTypeLabels: {
        broken_link: 'Link does not open',
        wrong_content: 'Content does not match',
        unsafe_content: 'Unsafe content',
        copyright: 'Copyright concern',
        low_quality: 'Low quality',
        other: 'Other',
      },
      reportStatusLabels: {
        open: 'Received',
        reviewing: 'Reviewing',
        resolved: 'Resolved',
        dismissed: 'Dismissed',
        cancelled: 'Cancelled',
      },
      aiSummaryTitle: 'AI Point Guide',
      aiSummaryGenerating: '✨ Generating...',
      aiSummaryRegenerate: '✨ Regenerate AI Guide',
      aiSummaryGenerate: '✨ Generate AI Guide',
      aiSummaryEmpty: 'No AI point guide has been generated yet. You can summarize what to focus on from the source link and point details.',
      providedContentTitle: 'Provided Learning Content',
      openCurrentContent: '🔗 Open Current Learning Content',
      youtubeSourceLabel: 'YouTube Source',
      openYoutubePlayer: 'Open with Official YouTube Player',
      openOriginalContent: 'Open Original Link',
      sourceLayoutHint: 'Payment benefits are only tied to AI feedback, learning-record analysis, review routines, artifact coaching, and personal learning reports, not video viewing.',
      emptyExternalLink: 'This exploration point does not have an external link yet.',
      replacementCompletedNotice: 'The reported source content was replaced. The previous source is preserved in the report history.',
      issueSummary: 'Report or Replace Learning Content',
      issueHelp: 'Use this only when the current learning content needs to be replaced because the link does not open or the content does not match.',
      reportTitle: 'Report Learning Content',
      reportHelp: 'After reporting, replacement recommendations and replacement actions will open.',
      reportPlaceholder: 'Describe the issue in at least 5 characters.',
      reportSubmitting: '🚩 Submitting...',
      reportAction: '🚩 Report Issue',
      reportGateHelp: 'Submit a report to open replacement recommendations and direct replacement.',
      replacementGateNotice: 'A report has been received. Get replacement recommendations or complete a direct replacement.',
      cancelReport: 'Cancel report and keep learning',
      replacementLoading: 'Finding recommendations...',
      loadReplacement: 'Get Replacement Recommendations',
      replacementCandidatesAria: 'Replacement candidates',
      candidateColumn: 'Content',
      replaceColumn: 'Replace',
      untitledCandidate: 'Untitled content',
      replace: 'Replace',
      directReplacementTitle: 'Replace with Content You Found',
      directReplacementHelp: 'Replace the reported source content',
      directReplacementTitlePlaceholder: 'Content title',
      directReplacementAction: 'Replace Directly',
      myReportsTitle: 'My Reports',
      reportCount: (count) => `${count} report${count === 1 ? '' : 's'}`,
      reportHistoryAria: 'My learning content report history',
      reportColumn: 'Report',
      statusColumn: 'Status',
      receivedAtColumn: 'Received',
      replacedPrefix: (value) => `Replaced: ${value}`,
      emptyReports: 'No reports yet.',
      cancelReportSuccessNotice: 'Report cancelled. Continue with the current learning content.',
      replacementQueryMissing: 'There is not enough context to find replacement content.',
      insufficientPoints: 'Not enough points to get replacement recommendations.',
      replacementLoadFailed: 'Could not get replacement recommendations.',
      replacementFound: 'Replacement candidates found. Review them and replace this content.',
      replacementEmpty: 'No matching replacement candidates were found. You can replace it with a link you found.',
      replacementUrlInvalid: 'Enter a link that starts with http:// or https://.',
      replacementFallbackTitle: 'Replacement learning content',
      replacementFailed: 'Could not replace the learning content.',
      replacementSuccess: 'The current learning content was replaced with the selected content.',
      replacementEntrySuccess: 'Learning content replaced. Open the new source content to continue learning.',
      candidateUrlMissing: 'The selected candidate does not have a usable URL.',
      cancelReportFailed: 'Could not cancel the report.',
      cancelReportEntrySuccess: 'Report cancelled. You can continue with the current learning content.',
      reportFailed: 'Could not submit the learning content report.',
      reportSuccess: 'Learning content report submitted.',
      aiUnavailable: 'No LLM is currently available. Check your AI settings.',
      aiSummaryFailed: 'Could not generate the AI point guide.',
      aiSummarySuccess: 'AI point guide generated.',
    },
    editor: {
      linkHttpsOnly: 'Enter a safe link that starts with https://.',
      imageHttpsOnly: 'Enter an image URL that starts with https://.',
      imageUploadFailed: 'Image upload failed.',
      paragraphStyle: 'Paragraph Style',
      paragraph: 'Body',
      heading1: 'Heading 1',
      heading2: 'Heading 2',
      heading3: 'Heading 3',
      undo: 'Undo',
      redo: 'Redo',
      quote: 'Quote',
      codeBlock: 'Code Block',
      bold: 'Bold',
      italic: 'Italic',
      underline: 'Underline',
      strike: 'Strikethrough',
      highlight: 'Highlight',
      inlineCode: 'Inline Code',
      superscript: 'Superscript',
      subscript: 'Subscript',
      bulletList: 'Bulleted List',
      orderedList: 'Numbered List',
      alignLeft: 'Align Left',
      alignCenter: 'Align Center',
      alignRight: 'Align Right',
      link: 'Link',
      image: 'Image',
      insertTable: 'Insert Table',
      addColumn: 'Add Column',
      addRow: 'Add Row',
      deleteColumn: 'Delete Column',
      deleteRow: 'Delete Row',
      deleteTable: 'Delete Table',
      apply: 'Apply',
      unset: 'Remove',
      imageUrlPlaceholder: 'https://... image URL',
      imageAltPlaceholder: 'Image description',
      insertImage: 'Insert',
      chooseFile: 'Choose File',
      uploadHint: 'Allowed: jpg, jpeg, png, webp, gif / up to 5MB',
      uploading: 'Uploading...',
      uploadInsert: 'Upload and Insert',
      inlineImageTypeError: 'Only jpg, jpeg, png, webp, and gif images can be uploaded.',
      inlineImageSizeError: 'Images must be 5MB or smaller.',
    },
    researchMaterial: {
      pageTitle: 'Create Research Material',
      nextActionNotice: 'Choose Video, Post, or File to create your research material.',
      blockLimitMessage: 'You have used all content posts available for this point. Save what you wrote and continue to Learn.',
      blockAddUnavailable: 'Cannot add another post.',
      blockSaveFailed: 'Could not save the research block.',
      blockAddSuccess: 'Research block added.',
      blockSaveSuccess: 'Research block saved.',
      deleteLockedMessage: 'Confirmed learning content cannot be edited or deleted. Mark it unconfirmed first if changes are needed.',
      blockDeleteFailed: 'Could not delete the research block.',
      confirmUnavailable: 'Saved learning content is required before confirmation.',
      confirmFailed: 'Could not confirm the learning content.',
      confirmSuccess: 'Learning content confirmed. You can now start this research point.',
      attachmentLimitMessage: (count) => `You can add up to ${count} learning content attachments.`,
      attachmentDirtyMessage: 'Save the attachment you are editing before adding another one.',
      attachmentAddUnavailable: 'Cannot add a learning content attachment.',
      attachmentAddFailed: 'Could not add the learning content attachment.',
      attachmentAddSuccess: 'Learning content attachment added.',
      subtitleDuplicate: 'Only one subtitle file can be uploaded for learning content.',
      thumbnailDuplicate: 'Only one thumbnail can be uploaded for learning content.',
      videoDuplicate: 'Only one video can be uploaded for learning content.',
      referenceTypeError: 'Only these learning content files can be added.\nImages: jpg, jpeg, png, webp, gif\nVideo/subtitles: mp4, webm, mov, vtt, srt\nDocuments/zip: pdf, doc, docx, hwp, hwpx, txt, xls, xlsx, csv, ppt, pptx, rtf, odt, md, zip',
      imageSizeError: 'The selected image is larger than 5MB.\nImages must be 5MB or smaller per file.',
      videoTypeError: 'Videos must be mp4, webm, or mov files.',
      videoSizeError: 'Videos must be 500MB or smaller.',
      videoDurationError: 'Videos must be 25 minutes or shorter.',
      videoMetadataError: 'Could not read the video metadata. Check that the file is mp4, webm, or mov.',
      subtitleTypeError: 'Subtitles must be vtt or srt files.',
      subtitleSizeError: 'Subtitles must be 2MB or smaller.',
      thumbnailTypeError: 'Thumbnails must be jpg, jpeg, png, webp, or gif images.',
      thumbnailSizeError: 'Thumbnails must be 5MB or smaller.',
      documentSizeError: 'The selected document/zip is larger than 20MB.\nDocuments and zip files must be 20MB or smaller per file.',
      uploadTooLargeServer: 'The file passed app checks, but exceeded the server upload limit. Ask an administrator to check the nginx upload limit.',
      uploadFailedStatus: (status) => `Learning content file upload failed. (HTTP ${status})`,
      fileUploadFailed: 'Learning content file upload failed.',
      fileUploadSuccess: 'File uploaded as learning content.',
      inlineImageUnavailable: 'Cannot upload the image.',
      inlineImageTypeError: 'Only jpg, jpeg, png, webp, and gif images can be uploaded.',
      inlineImageSizeError: 'Images must be 5MB or smaller.',
      inlineImageUploadFailed: 'Image upload failed.',
      inlineImageInserted: 'Image inserted into the learning material body.',
      statusConfirmed: '✅ Confirmed',
      statusDraft: 'Unconfirmed',
      confirmSaving: 'Confirming...',
      confirmDone: '✅ Learning Content Confirmed',
      confirmAction: 'Confirm Learning Content',
      confirmModalTitle: 'Confirm this learning content?',
      confirmModalMessage: 'Once confirmed, this research material is used in the Learn flow. To edit it later, mark it unconfirmed first.',
      unconfirmAction: 'Mark Unconfirmed',
      unconfirmModalTitle: 'Mark this content unconfirmed?',
      unconfirmModalMessage: 'Marking it unconfirmed lets you edit the research material again. Confirm it again after making changes.',
      unconfirmSaving: 'Updating...',
      unconfirmFailed: 'Could not mark the learning content as unconfirmed.',
      unconfirmSuccess: 'Learning content is now unconfirmed.',
      confirmedHelp: 'Confirmed learning content is used in the Learn flow below. Mark it unconfirmed first if changes are needed.',
      confirmReadyHelp: 'Confirm saved learning content to use it in the Learn flow below.',
      confirmBlockedHelp: 'Saved learning content is required before confirmation.',
      attachmentFallback: 'Attachment',
      videoUrlFailed: 'Could not create the video URL.',
      contentFallback: 'Learning Content',
      attachmentsTitle: 'Attachments',
      uploadedFileListTitle: 'Uploaded Files',
      openFileLink: 'Open File',
      attachmentLimitTitle: 'You can add up to 10 attachments.',
      attachmentAddTitle: 'Add Attachment',
      uploadingAttachment: '📎 Uploading...',
      addAttachment: 'File',
      attachmentDeleteTitle: 'Delete Attachment',
      delete: '🗑️ Delete',
      attachmentImageAlt: 'Attached image',
      emptyAttachments: 'No files attached.',
      attachmentHelp: 'Allowed: images up to 5MB, documents/zip up to 20MB, up to 10 files',
      videoSectionTitle: 'Learning Content Video',
      uploadedVideoListTitle: 'Uploaded Video',
      openVideoLink: 'Open Video',
      videoTitleLabel: 'Video Title',
      videoTitlePlaceholder: 'Video title',
      saving: 'Saving...',
      saveTitle: 'Save Title',
      videoPlayerTitle: 'Learning Video',
      uploadSubtitle: 'Upload Subtitles',
      uploadThumbnail: 'Upload Thumbnail',
      deleteSubtitle: 'Delete Subtitles',
      deleteThumbnail: 'Delete Thumbnail',
      deleteVideo: '🗑️ Delete Video',
      uploadingVideo: '🎬 Uploading...',
      uploadVideo: 'Video',
      videoHelp: 'mp4, webm, mov / up to 500MB / up to 25 minutes',
      noVideo: 'No video has been uploaded.',
      contentSectionTitle: 'Learning Content',
      addContent: 'Post',
      contentTitleLabel: 'Content Title',
      contentTitlePlaceholder: 'Content title',
      contentPlaceholder: 'Enter the original text, summary, or observations you want to study.',
      contentDisplayTitle: (title) => `Title: ${title}`,
      emptyContent: 'No content has been written.',
      deleteContentTitle: 'Delete Learning Content',
      deleteContentMessage: (title) => `Delete ${title} content?`,
      cancel: '↩️ Cancel',
      deleting: '🗑️ Deleting...',
      deleteAttachmentModalTitle: 'Delete Attachment',
      deleteAttachmentMessage: (title) => `Delete ${title} file?`,
      validationTitle: 'Check the attachment',
      validationConfirm: '✅ OK',
      uploadProgressTitle: 'Uploading Learning Content',
      uploadProgressMessage: 'Adding the selected file to learning content. Videos may take longer depending on file size.',
      uploadProgressAria: (progress) => `Upload progress ${progress}%`,
      copyrightNotice: 'Upload only material you created or have permission to use. Use other people’s material only within allowed terms.',
      copyrightAgreementLinkLabel: 'View related terms',
      statusLabel: 'Learning Content Status',
      observeDraftTitle: 'Learning content needs to be written.',
      observeDraftDescription: 'Observation does not open authoring tools directly. Enter edit mode to write or revise research material, then save and confirm it.',
      openEditMode: '✏️ Edit Research Material',
      navAria: 'Learning content subsections',
      navVideo: 'Video',
      navContent: 'Post',
      navAttachments: 'File',
    },
    learningWork: {
      title: 'Learn',
      subtitle: 'Start with notes. Before completing this point, save at least one answered question and a review. Practice records and artifacts are optional.',
      lockedNotice: 'Research points stay locked until the learning content is confirmed. Save your material first, then confirm the content.',
      miniNavLabel: 'Learning work subsections',
      notesLabel: 'Notes',
      recordsLabel: 'Records',
      recordsTitle: 'Learning Records',
      optionalLabel: '(Optional)',
      recordsDescription: 'Point completion requires one answered question. Add practice records and artifacts only when you have activities or outcomes worth keeping.',
      recordsLockedNotice: 'After saving one note, complete the point by saving one answered question and a review. Practice records and artifacts are optional.',
      moveModalTitle: 'Move Records',
      moveModalEyebrow: 'Lumi Confirm',
      moveModalMessage: 'Unsaved changes may be lost. Move to another learning record?',
      keepWriting: '↩️ Keep Writing',
      move: 'Move',
      sourceContentGuide: {
        title: 'Learning Notes',
        message: 'You opened the source content. Save one key idea in Notes while it is fresh.',
      },
      notes: {
        title: 'Notes',
        edit: '✏️ Edit',
        helper: 'You do not need to fill everything. Saving one useful item is enough to continue.',
        cancelModalTitle: 'Cancel Note Edit',
        cancelModalEyebrow: 'Lumi Confirm',
        cancelModalMessage: 'Discard unsaved note changes?',
        back: '↩️ Go Back',
        cancel: '🚫 Cancel',
        cancelButton: '↩️ Cancel',
        saving: '💾 Saving...',
        save: '💾 Save Notes',
        progressTitle: 'Notes Progress',
        addEmpty: 'Add Empty Item',
        progressIncomplete: 'Fill only what you need. Add an empty item if you want to organize more.',
        progressComplete: 'All note items are filled. You can refine them with Edit.',
        empty: 'No saved notes yet.',
        fields: {
          observation: { icon: '💡', label: 'Core Concept', placeholder: 'Summarize the core concept.' },
          reflection: { icon: '🗣️', label: 'My Explanation', placeholder: 'Explain it in your own words.' },
          examples: { icon: '🔎', label: 'Examples', placeholder: 'Add examples or application situations.' },
          nextStep: { icon: '❔', label: 'Confusing Parts', placeholder: 'Write what is unclear or what to check next.' },
        },
      },
      recordPanels: {
        common: {
          number: 'No.',
          title: 'Title',
          status: 'Status',
          type: 'Type',
          previous: '← Previous',
          next: 'Next →',
          countRange: (total, start, end) => `${start}-${end} of ${total}`,
        },
        questions: {
          title: 'Questions and Answers',
          subtitle: '',
          add: '➕ Question',
          listAria: 'Question list',
          answered: 'Answered',
          pending: 'Waiting',
          empty: 'No saved questions yet.',
          fallbackTitle: 'No question',
          typeLabels: { reflection: 'Unclear', application: 'How to Apply', goal_alignment: 'Want to Know More', none: 'No Type' },
          detail: {
            editTitle: 'Edit Question',
            viewTitle: 'Question',
            addTitle: 'Add Question',
            addSubtitle: 'Write the question you want to keep.',
            type: 'Type',
            title: 'Title',
            content: 'Question',
            answer: 'My Answer',
            answerMeta: 'Answer and evidence you found',
            answerEmptyMeta: 'One question can have one answer.',
            aiCoach: 'AI Coaching',
            aiCoachMeta: 'Save to keep this record',
            titlePlaceholder: 'Question title',
            contentPlaceholder: 'Enter your question.',
            answerPlaceholder: 'Write the answer you found, evidence you checked, and what you understand now.',
            noContent: 'No question details.',
            noAnswer: 'No answer has been added yet.',
            communityDisabled: '💬 Ask Community · Coming Soon',
            addAnswer: '✍️ Add Answer',
            editAnswer: '✏️ Edit Answer',
            createAiCoach: '✨ Create AI Coaching',
            createAiCoachStarter: '✨ Answer Direction Coaching',
            createAiCoachRevision: '✨ Answer Improvement Coaching',
            generatingAiCoach: '✨ Generating...',
            saveDetail: '💾 Save Question',
            savingDetail: '💾 Saving...',
            saveAnswer: '💾 Save Answer',
            savingAnswer: '💾 Saving...',
            register: '➕ Add Question',
            saving: '💾 Saving...',
            edit: '✏️ Edit Question',
            list: '📋 List',
            delete: '🗑️ Delete',
            deleting: '🗑️ Deleting...',
            cancel: '🚫 Cancel',
            back: '↩️ Go Back',
            ok: 'OK',
            answerInstruction: 'After adding the question, add one answer as a reply.',
            aiCoachNotice: 'AI coaching is available only while writing an answer. Save the answer to keep it.',
            aiCoachStarterHint: 'AI can suggest a direction before you write the answer.',
            aiCoachRevisionHint: 'AI can suggest missing parts and improvements for the answer you are drafting.',
            aiCoachSummary: 'Feedback Summary',
            aiCoachNextAction: 'Next Action',
            saveSuccessTitle: 'Question Saved',
            saveSuccessMessage: 'Question was saved.',
            saveFailTitle: 'Save Failed',
            saveFailMessage: 'Question could not be saved.',
            answerSaveSuccessTitle: 'Answer Saved',
            answerSaveSuccessMessage: 'Answer was saved.',
            answerSaveFailTitle: 'Save Failed',
            answerSaveFailMessage: 'Answer could not be saved.',
            deleteModalTitle: 'Delete Question',
            deleteModalMessage: (title) => `Delete '${title}'?`,
            listReturnTitle: 'Close Question',
            listReturnMessage: 'Return to the list without saving your edits?',
            addCancelTitle: 'Cancel Question',
            addCancelMessage: 'Discard this unsaved question?',
          },
        },
        practice: {
          title: 'Practice Log',
          subtitle: 'Select a title to review details.',
          add: '➕ New Practice',
          listAria: 'Practice log list',
          hasRecord: 'Recorded',
          waiting: 'Waiting',
          empty: 'No saved practice records yet.',
          fallbackTitle: (index) => `Practice ${index}`,
          detail: {
            editTitle: 'Edit Practice',
            viewTitle: 'Practice Details',
            addTitle: 'Add Practice Record',
            title: 'Title',
            duration: 'Time Spent (min)',
            achievement: 'Outcome',
            reflection: 'Blockers / Next Practice',
            titlePlaceholder: 'Practice title, e.g. applying a concept, survey attempt, data analysis',
            durationPlaceholder: 'Enter minutes',
            achievementPlaceholder: 'Outcome',
            reflectionPlaceholder: 'Record blockers and your next practice plan.',
            addReflection: '➕ Add Blockers / Next Practice',
            noDetail: 'No details saved yet.',
            list: '📋 List',
            delete: '🗑️ Delete',
            deleting: '🗑️ Deleting...',
            save: 'Save',
            saving: 'Saving...',
            edit: '✏️ Edit',
            cancel: '↩️ Cancel',
            back: '↩️ Go Back',
            deleteModalTitle: 'Delete Practice',
            deleteModalMessage: (title) => `Delete '${title}'?`,
            listReturnTitle: 'Back to List',
            listReturnMessage: 'Return to the list without saving your edits?',
            addCancelTitle: 'Cancel Practice',
            addCancelMessage: 'Discard this unsaved practice record?',
          },
        },
        artifacts: {
          title: 'Artifacts',
          subtitle: 'Select a title to review artifact details.',
          add: '➕ Artifact',
          empty: 'No saved artifacts yet.',
          listAria: 'Artifact list',
          fallbackTitle: (index) => `Artifact ${index}`,
          kindLabels: { 영상자료: 'Video', '문서 편집자료': 'Document', 첨부자료: 'File' },
          detail: {
            editTitle: 'Edit Artifact',
            viewTitle: 'Artifact Details',
            addTitle: 'Create Artifact',
            kind: 'Artifact Type',
            title: 'Title',
            editTitleLabel: 'Edit Title',
            videoSection: 'Artifact Video',
            attachmentSection: 'Artifact Files',
            documentSection: 'Document',
            videoPlayerTitle: 'Artifact Video',
            titlePlaceholder: 'Artifact title',
            documentTitlePlaceholder: 'Document title',
            documentPlaceholder: 'Write your artifact document.',
            attachmentFallback: 'Attachment',
            noVideo: 'No video has been uploaded.',
            noFiles: 'No files attached.',
            uploadVideo: '🎬 Upload Video',
            uploadingVideo: '🎬 Uploading...',
            uploadSubtitle: 'Upload Subtitles',
            uploadThumbnail: 'Upload Thumbnail',
            deleteSubtitle: 'Delete Subtitles',
            deleteThumbnail: 'Delete Thumbnail',
            deleteVideo: '🗑️ Delete Video',
            uploadFile: '➕ Attach File',
            uploadingFile: '📎 Uploading...',
            list: '📋 List',
            delete: '🗑️ Delete',
            deleting: '🗑️ Deleting...',
            save: '💾 Save',
            saving: '💾 Saving...',
            edit: '✏️ Edit',
            cancel: '↩️ Cancel',
            back: '↩️ Go Back',
            videoHelp: 'mp4, webm, mov / up to 500MB / up to 25 minutes',
            deleteModalTitle: 'Delete Artifact',
            deleteModalMessage: (title) => `Delete '${title}'?`,
            deleteAttachmentTitle: 'Delete Attachment',
            deleteAttachmentMessage: (title) => `Delete ${title}?`,
            listReturnTitle: 'Back to List',
            listReturnMessage: 'Return to the list without saving your edits?',
            addCancelTitle: 'Cancel Artifact',
            addCancelMessage: 'Discard this unsaved artifact?',
          },
        },
        attachments: {
          title: 'References',
          subtitle: (researchCount, workCount) => `Content ${researchCount} / References ${workCount}`,
          add: '➕ Add Reference',
          listAria: 'Reference list',
          learningContent: 'Content',
          reference: 'Reference',
          empty: 'No saved references yet.',
          fallbackTitle: (index) => `Reference ${index}`,
          detail: {
            editTitle: 'Edit Reference',
            viewTitle: 'Reference Details',
            addTitle: 'Add Reference',
            sectionTitle: 'Reference File',
            title: 'Title',
            editTitleLabel: 'Edit Title',
            file: 'File',
            titlePlaceholder: 'Reference title',
            attachmentFallback: 'Attachment',
            replacementPending: 'Replacement selected',
            selectedFile: (name) => `Selected: ${name}`,
            researchMaterialNotice: 'Edit this from the learning content area.',
            uploadFile: '📎 Upload File',
            uploadingFile: '📎 Uploading...',
            uploadSelected: '📎 Upload Selected File',
            list: '📋 List',
            delete: '🗑️ Delete',
            deleting: '🗑️ Deleting...',
            save: '💾 Save',
            saving: '💾 Saving...',
            edit: '✏️ Edit',
            cancel: '🚫 Cancel',
            back: '↩️ Go Back',
            deleteModalTitle: 'Delete Reference',
            deleteModalMessage: (title) => `Delete '${title}'?`,
            listReturnTitle: 'Back to List',
            listReturnMessage: 'Return to the list without saving your edits?',
            addCancelTitle: 'Cancel Reference',
            addCancelMessage: 'Cancel adding this reference?',
          },
        },
        timeline: {
          title: 'Activity Log',
          subtitle: 'Auto-generated timeline',
          listAria: 'Activity log list',
          time: 'Time',
          activity: 'Activity',
          empty: 'No activity has been recorded yet.',
        },
      },
      tabs: {
        questions: { label: 'Questions', title: 'Questions', message: 'Use this area to keep what is unclear or worth checking. Save a question, write your answer, and use AI feedback when helpful.' },
        practice: { label: 'Practice', title: 'Practice Log', message: 'Record practice, research, analysis, or application attempts. A short note about time spent, blockers, and next steps is enough.' },
        artifacts: { label: 'Artifacts', title: 'Artifacts', message: 'Save outputs or submission links here. Include what you made, what you learned, and what was difficult.' },
        attachments: { label: 'References', title: 'References', message: 'Collect extra links or files you used while learning. Keep sources separate from your notes.' },
        timeline: { label: 'Activity', title: 'Activity Log', message: 'Review the key actions saved in this point in chronological order.' },
      },
      noteGuides: {
        observation: { title: 'Core Concept', message: 'Save one or two key sentences from what you learned. Keep it short enough to recognize the main idea later.' },
        reflection: { title: 'My Explanation', message: 'Rewrite the idea in your own words, as if explaining it to someone else.' },
        examples: { title: 'Examples', message: 'Save situations or examples where this concept could be used. Application examples make it easier to revisit later.' },
        nextStep: { title: 'Confusing Parts', message: 'Save what is still unclear and what you want to check next. It is fine if the thought is not fully organized yet.' },
      },
    },
  },
  completion: {
    encouragement: 'It does not need to be perfect. When this point feels sufficiently learned for now, mark it complete.',
    finalAction: 'Final Action',
    completedLocked: 'This point is complete and locked. Cancel completion first if you need to edit it.',
    readyInstruction: 'Check readiness and complete the current point.',
    cancelProcessing: 'Cancelling...',
    cancelCompletion: 'Cancel Completion',
    saving: 'Saving...',
    completePoint: 'Complete Point',
    completeModalTitle: 'Complete Point',
    completeModalEyebrow: 'Lumi Confirm',
    completeModalMessage: 'Mark this point as complete? You can still review records after completion.',
    close: 'Close',
    cancel: 'Cancel',
    complete: 'Complete',
    cancelModalTitle: 'Cancel Completion',
    cancelModalMessage: 'Cancelling completion unlocks this point for editing again. Cancel completion?',
    cancelDone: 'Cancel Completion',
  },
  selfEvaluation: {
    title: 'Review',
    subtitle: 'Check yourself with 3 short application questions',
    savedTitle: 'Saved Review',
    applicationTitle: 'Application Check',
    savedDescription: 'Your review for this point has been saved.',
    applicationDescription: 'Answer the short application questions, then ask AI for a reference review to save it.',
    generatingTitle: 'Preparing Review',
    generatingDescription: 'AI is preparing application questions and reference review notes.',
    questionLabel: 'Question',
    answerPlaceholder: 'Write a short answer.',
    aiReferenceTitle: 'AI Reference Review',
    aiReferenceDescription: 'AI organized reference review notes from your application answers. These notes were saved as internal review values.',
    generatingButton: 'Generating...',
    restartButton: 'Review Again',
    aiReferenceButton: 'AI Reference',
    createQuestionsButton: 'Create Questions',
    finalTitle: 'My Final Review',
    finalDescription: 'After checking the AI reference review, choose the score closest to your current state.',
    scoreSuffix: '',
    lockedUntilDraft: 'You can choose a score after the AI reference review is saved.',
    options: [
      { value: '1', label: '1', description: 'Still difficult' },
      { value: '2', label: '2', description: 'Partly understood' },
      { value: '3', label: '3', description: 'Fair' },
      { value: '4', label: '4', description: 'Well understood' },
      { value: '5', label: '5', description: 'Can explain it' },
    ],
  },
  completionStatus: {
    readOnlyTitle: 'Read-Only View',
    readOnlyDescription: 'In shared mode, you can review current learning results without editing questions, review, records, or research blocks.',
    achievementTitle: 'Point Complete',
    achievementDescription: 'Review what you recorded during this exploration.',
    contentChecked: '✓ Content reviewed',
    noteCount: (count) => `✓ ${count} note${count === 1 ? '' : 's'} written`,
    questionCount: (total, answered) => `✓ ${total} question${total === 1 ? '' : 's'} / ${answered} answered`,
    practiceCount: (count) => `✓ ${count} practice record${count === 1 ? '' : 's'}`,
    artifactCount: (count) => `✓ ${count} artifact${count === 1 ? '' : 's'}`,
    selfEvaluationState: (completed) => `✓ Review ${completed ? 'complete' : 'incomplete'}`,
    backToDiary: 'Back to Diary',
    sharedStatusTitle: 'Completion Status',
    readyStatusTitle: 'Readiness',
    researchMaterialConfirmed: (completed) => `• Content confirmed: ${completed ? 'Done' : 'Not yet'}`,
    answeredQuestionReady: (completed, count) => `• At least 1 answered question: ${completed ? `Done (${count})` : 'Not yet'}`,
    selfEvaluationSaved: (completed) => `• Review saved: ${completed ? 'Done' : 'Not yet'}`,
    selfEvaluationQuality: (completed) => `• Review quality check: ${completed ? 'Done' : 'Not yet'}`,
    goalConnection: (completed) => `• Goal connection note: ${completed ? 'Done' : 'Not yet'}`,
    sharedNotReady: 'This read-only view shows the saved reflection and completion readiness so far.',
    researchNotReady: 'For research points, confirm content, answer a question, and save a review before completion.',
    learningNotReady: 'Save an answered question and review before completing this point.',
    sharedReady: 'This point is shared with completion conditions satisfied.',
    learningReady: 'All completion conditions are satisfied.',
  },
  lumiGuide: {
    goalLabel: 'Goal',
    close: 'Close Lumi guide',
    closeButton: '💡 Lumi Guide ❌',
    featureGuideLabel: 'Guide',
    currentTaskLabel: 'Current Task',
    fallbackGoal: 'Could not load the learning goal connected to this diary.',
    pointTitle: (isExplorationPoint) => (isExplorationPoint ? 'Exploration Point' : 'Research Point'),
    pointMessage: (isExplorationPoint) => (isExplorationPoint
      ? 'This point turns prepared material into your own learning records. Start by checking the flow lightly.'
      : 'This point uses your own material as learning content. Refine and confirm the material before moving on.'),
    learningContentTitle: 'Content',
    learningContentMessage: (isExplorationPoint) => (isExplorationPoint
      ? 'Review the provided content, source, and summary, then move what matters into your own notes.'
      : 'Prepare the learning content here. Organize the post and attachments, confirm them, then continue with notes and questions.'),
    selfEvaluationTitle: 'Review',
    selfEvaluationMessage: 'At the end, check your understanding, ability to apply it, and connection to the goal. It does not need to be perfect. Save your current judgment.',
    completedFlow: 'This point is complete. You can review the records and review notes you saved.',
    researchMaterialFlow: 'Prepare and confirm the learning content with a post or attachments first.',
    openExternalFlow: 'Open the source content first. Then save one useful item in Notes to unlock the next records.',
    firstNoteFlow: 'Save one useful item in Notes first. Completion later requires one answered question and a review.',
    selfEvaluationReadyFlow: 'You have learning records now. Add more records if useful, or start the review based on what you understand.',
    noQuestionFlow: 'Your note is saved. To complete this point, add one question and save an answer. Practice records and artifacts are optional.',
    noAnswerFlow: 'Your question is saved. To complete this point, save an answer. You can also add practice records or artifacts if useful.',
    completionReadyFlow: 'This point is ready to complete. Use the completion button below the review section.',
    fillMissingFlow: 'Check the completion conditions below and fill the missing items one by one.',
    restTitle: '🌿 A Short Break Is Fine',
    restMessage: 'You have been focused for a while. It is fine to rest and continue later. Move at your own pace.',
    restConfirm: 'OK',
  },
}

export function getPointLearningCopy(locale: Locale | string | null | undefined): PointLearningCopy {
  return normalizeLocale(locale) === 'en' ? en : ko
}
