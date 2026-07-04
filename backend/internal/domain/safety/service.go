package safety

import (
	"context"
	"errors"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

var ErrInvalidAIOutputReviewDecision = errors.New("invalid ai output review decision")

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Moderate(ctx context.Context, input ModerateInput) (*ModerationResult, error) {
	normalized := normalizeText(input.Text)
	if normalized == "" {
		return allowResult(), nil
	}
	rules := append([]Rule{}, builtinRules...)
	if s != nil && s.repo != nil {
		dbRules, err := s.repo.ListActiveRules(ctx, input.Locale)
		if err != nil {
			return nil, err
		}
		rules = append(dbRules, rules...)
	}
	matched := matchRule(normalized, rules)
	if matched == nil {
		return allowResult(), nil
	}
	result := &ModerationResult{
		Action:        matched.Action,
		RiskType:      strings.TrimSpace(matched.RiskType),
		MatchedRuleID: matched.ID,
	}
	if result.RiskType == "" {
		result.RiskType = "unknown"
	}
	if result.Action == ActionSoftWarn {
		code := ErrorCodeSoftWarn
		message := userMessage(input.Locale)
		result.ErrorCode = &code
		result.Message = &message
		return result, nil
	}
	result.Action = ActionBlock
	code := ErrorCodeBlocked
	message := userMessage(input.Locale)
	result.ErrorCode = &code
	result.Message = &message
	return result, nil
}

func (s *Service) Enforce(ctx context.Context, input ModerateInput) (*ModerationResult, error) {
	result, err := s.Moderate(ctx, input)
	if err != nil {
		return nil, err
	}
	if s != nil && s.repo != nil {
		hash := textHash(input.Text)
		if err := s.repo.InsertLog(ctx, input, result, hash); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *Service) ObserveAIOutput(ctx context.Context, input ModerateInput) *ModerationResult {
	input.Metadata = aiOutputMetadata(input.Metadata)
	result, err := s.Moderate(ctx, input)
	if err != nil {
		log.Printf("[safety] ai output moderation skipped target_type=%s route=%s err=%v", input.TargetType, input.Route, err)
		return allowResult()
	}
	if s != nil && s.repo != nil {
		hash := textHash(input.Text)
		if err := s.repo.InsertLog(ctx, input, result, hash); err != nil {
			log.Printf("[safety] ai output moderation log skipped target_type=%s route=%s action=%s risk=%s err=%v", input.TargetType, input.Route, result.Action, result.RiskType, err)
		}
	}
	return result
}

func aiOutputMetadata(metadata map[string]any) map[string]any {
	output := make(map[string]any, len(metadata)+1)
	for key, value := range metadata {
		output[key] = value
	}
	output["direction"] = "ai_output"
	return output
}

func (s *Service) ListLogs(ctx context.Context, input ListLogsInput) ([]LogEntry, error) {
	if s == nil || s.repo == nil {
		return []LogEntry{}, nil
	}
	return s.repo.ListLogs(ctx, input)
}

func (s *Service) GetAIOutputSummary(ctx context.Context) (AIOutputSummary, error) {
	generatedAt := time.Now().UTC()
	observations := []AIOutputLogObservation{}
	var err error
	if s != nil && s.repo != nil {
		observations, err = s.repo.ListAIOutputObservations(ctx, 24*7)
		if err != nil {
			return AIOutputSummary{}, err
		}
	}
	return BuildAIOutputSummary(generatedAt, observations), nil
}

func (s *Service) GetAIOutputReviewCandidates(ctx context.Context) (AIOutputReviewCandidates, error) {
	generatedAt := time.Now().UTC()
	observations := []AIOutputLogObservation{}
	var err error
	if s != nil && s.repo != nil {
		observations, err = s.repo.ListAIOutputObservations(ctx, 24*30)
		if err != nil {
			return AIOutputReviewCandidates{}, err
		}
	}
	result := BuildAIOutputReviewCandidates(generatedAt, 24*30, observations)
	if s == nil || s.repo == nil || len(result.Candidates) == 0 {
		return result, nil
	}
	keys := make([]string, 0, len(result.Candidates))
	for _, candidate := range result.Candidates {
		keys = append(keys, candidate.Key)
	}
	decisions, err := s.repo.ListAIOutputReviewDecisions(ctx, keys)
	if err != nil {
		return AIOutputReviewCandidates{}, err
	}
	ApplyAIOutputReviewDecisions(&result, decisions)
	return result, nil
}

func (s *Service) SaveAIOutputReviewDecision(ctx context.Context, input AIOutputReviewDecisionInput) (AIOutputReviewDecision, error) {
	cleaned, err := validateAIOutputReviewDecisionInput(input)
	if err != nil {
		return AIOutputReviewDecision{}, err
	}
	if s == nil || s.repo == nil {
		return AIOutputReviewDecision{
			CandidateKey: cleaned.CandidateKey,
			Status:       AIOutputReviewStatus(cleaned.Status),
			Note:         cleaned.Note,
			ReviewedBy:   cleaned.ReviewedBy,
			Metadata:     cleaned.Metadata,
		}, nil
	}
	return s.repo.UpsertAIOutputReviewDecision(ctx, cleaned)
}

func validateAIOutputReviewDecisionInput(input AIOutputReviewDecisionInput) (AIOutputReviewDecisionInput, error) {
	cleaned := AIOutputReviewDecisionInput{
		CandidateKey: strings.TrimSpace(input.CandidateKey),
		Status:       strings.TrimSpace(input.Status),
		Note:         strings.TrimSpace(input.Note),
		ReviewedBy:   strings.TrimSpace(input.ReviewedBy),
		Metadata:     sanitizeSafetyMetadata(input.Metadata),
	}
	if cleaned.CandidateKey == "" || !isValidAIOutputReviewStatus(cleaned.Status) {
		return AIOutputReviewDecisionInput{}, ErrInvalidAIOutputReviewDecision
	}
	return cleaned, nil
}

func isValidAIOutputReviewStatus(status string) bool {
	switch AIOutputReviewStatus(strings.TrimSpace(status)) {
	case AIOutputReviewStatusUnreviewed, AIOutputReviewStatusFalsePositive, AIOutputReviewStatusNeedsPromptGuard, AIOutputReviewStatusNeedsRuleTuning, AIOutputReviewStatusNeedsMasking, AIOutputReviewStatusNeedsRegenerate, AIOutputReviewStatusResolved:
		return true
	default:
		return false
	}
}

func defaultAIOutputReviewDecision(candidateKey string) AIOutputReviewDecision {
	return AIOutputReviewDecision{CandidateKey: candidateKey, Status: AIOutputReviewStatusUnreviewed, Metadata: map[string]any{}}
}

func ApplyAIOutputReviewDecisions(candidates *AIOutputReviewCandidates, decisions map[string]AIOutputReviewDecision) {
	if candidates == nil {
		return
	}
	for index := range candidates.Candidates {
		key := candidates.Candidates[index].Key
		decision, ok := decisions[key]
		if !ok {
			candidates.Candidates[index].Review = defaultAIOutputReviewDecision(key)
			continue
		}
		decision.Metadata = sanitizeSafetyMetadata(decision.Metadata)
		candidates.Candidates[index].Review = decision
	}
}

func BuildAIOutputReviewCandidates(generatedAt time.Time, windowHours int, observations []AIOutputLogObservation) AIOutputReviewCandidates {
	if windowHours <= 0 || windowHours > 24*30 {
		windowHours = 24 * 30
	}
	since := generatedAt.Add(-time.Duration(windowHours) * time.Hour)
	candidatesByKey := map[string]*AIOutputReviewCandidate{}

	for _, observation := range observations {
		createdAt := observation.CreatedAt
		if createdAt.Location() != time.UTC {
			createdAt = createdAt.UTC()
		}
		if createdAt.Before(since) || createdAt.After(generatedAt) || !isRiskyAIOutputAction(observation.Action) {
			continue
		}

		targetType := normalizedAIOutputDimension(observation.TargetType)
		riskType := normalizedAIOutputRiskType(observation.RiskType)
		sourceFeature := normalizedAIOutputDimension(observation.SourceFeature)
		matchedRuleID := strings.TrimSpace(observation.MatchedRuleID)
		matchedPattern := strings.TrimSpace(observation.MatchedPattern)
		key := strings.Join([]string{targetType, riskType, sourceFeature, matchedRuleID, matchedPattern}, "|")
		candidate := candidatesByKey[key]
		if candidate == nil {
			candidate = &AIOutputReviewCandidate{
				Key:            key,
				TargetType:     targetType,
				RiskType:       riskType,
				SourceFeature:  sourceFeature,
				MatchedRuleID:  matchedRuleID,
				MatchedPattern: matchedPattern,
				LatestSeenAt:   createdAt,
				SampleTargetID: observation.TargetID,
				Review:         defaultAIOutputReviewDecision(key),
			}
			candidatesByKey[key] = candidate
		}
		candidate.Total++
		addAIOutputAction(&candidate.ActionCounts, observation.Action)
		if createdAt.After(candidate.LatestSeenAt) {
			candidate.LatestSeenAt = createdAt
			candidate.SampleTargetID = observation.TargetID
		}
	}

	candidates := make([]AIOutputReviewCandidate, 0, len(candidatesByKey))
	for _, candidate := range candidatesByKey {
		candidates = append(candidates, *candidate)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Total == candidates[j].Total {
			if candidates[i].LatestSeenAt.Equal(candidates[j].LatestSeenAt) {
				return candidates[i].Key < candidates[j].Key
			}
			return candidates[i].LatestSeenAt.After(candidates[j].LatestSeenAt)
		}
		return candidates[i].Total > candidates[j].Total
	})

	return AIOutputReviewCandidates{
		GeneratedAt:     generatedAt,
		WindowHours:     windowHours,
		Since:           since,
		TotalCandidates: len(candidates),
		Candidates:      candidates,
	}
}

func BuildAIOutputSummary(generatedAt time.Time, observations []AIOutputLogObservation) AIOutputSummary {
	windows := []struct {
		label string
		hours int
	}{
		{label: "24h", hours: 24},
		{label: "7d", hours: 24 * 7},
	}
	result := AIOutputSummary{GeneratedAt: generatedAt, Windows: make([]AIOutputSummaryWindow, 0, len(windows))}
	for _, window := range windows {
		since := generatedAt.Add(-time.Duration(window.hours) * time.Hour)
		windowObservations := make([]AIOutputLogObservation, 0)
		for _, observation := range observations {
			createdAt := observation.CreatedAt
			if createdAt.Location() != time.UTC {
				createdAt = createdAt.UTC()
			}
			if !createdAt.Before(since) && !createdAt.After(generatedAt) {
				windowObservations = append(windowObservations, observation)
			}
		}
		result.Windows = append(result.Windows, buildAIOutputSummaryWindow(window.label, window.hours, since, windowObservations))
	}
	return result
}

func buildAIOutputSummaryWindow(label string, hours int, since time.Time, observations []AIOutputLogObservation) AIOutputSummaryWindow {
	window := AIOutputSummaryWindow{Label: label, Hours: hours, Since: since, Targets: []AIOutputTargetSummary{}, RiskTypes: []AIOutputDimensionCount{}, SourceFeatures: []AIOutputDimensionCount{}, MatchedRules: []AIOutputRuleCount{}}
	targets := map[string]*AIOutputTargetSummary{}
	riskTypes := map[string]int{}
	sourceFeatures := map[string]int{}
	matchedRules := map[string]AIOutputRuleCount{}
	maxSourceRiskRepeats := 0

	for _, observation := range observations {
		window.Total++
		addAIOutputAction(&window.ActionCounts, observation.Action)
		risky := isRiskyAIOutputAction(observation.Action)
		if risky {
			window.RiskyTotal++
		}

		targetType := strings.TrimSpace(observation.TargetType)
		if targetType == "" {
			targetType = "unknown"
		}
		target := targets[targetType]
		if target == nil {
			target = &AIOutputTargetSummary{TargetType: targetType}
			targets[targetType] = target
		}
		target.Total++
		addAIOutputAction(&target.ActionCounts, observation.Action)
		if risky {
			target.RiskyTotal++
		}

		if risky {
			riskType := strings.TrimSpace(observation.RiskType)
			if riskType == "" || riskType == "none" {
				riskType = "unknown"
			}
			riskTypes[riskType]++

			sourceFeature := strings.TrimSpace(observation.SourceFeature)
			if sourceFeature == "" {
				sourceFeature = "unknown"
			}
			sourceFeatures[sourceFeature]++
			if sourceFeatures[sourceFeature] > maxSourceRiskRepeats {
				maxSourceRiskRepeats = sourceFeatures[sourceFeature]
			}

			ruleKey := strings.TrimSpace(observation.MatchedRuleID)
			if ruleKey == "" {
				ruleKey = "builtin:" + strings.TrimSpace(observation.MatchedPattern)
			}
			if ruleKey != "builtin:" {
				current := matchedRules[ruleKey]
				current.RuleID = strings.TrimSpace(observation.MatchedRuleID)
				current.Pattern = strings.TrimSpace(observation.MatchedPattern)
				current.Count++
				matchedRules[ruleKey] = current
			}
		}
	}

	for _, target := range targets {
		target.RiskyRate = ratio(target.RiskyTotal, target.Total)
		target.Severity = targetSeverity(*target)
		window.Targets = append(window.Targets, *target)
	}
	sort.Slice(window.Targets, func(i, j int) bool {
		if window.Targets[i].RiskyTotal == window.Targets[j].RiskyTotal {
			return window.Targets[i].TargetType < window.Targets[j].TargetType
		}
		return window.Targets[i].RiskyTotal > window.Targets[j].RiskyTotal
	})
	window.RiskTypes = topDimensionCounts(riskTypes, 6)
	window.SourceFeatures = topDimensionCounts(sourceFeatures, 6)
	window.MatchedRules = topRuleCounts(matchedRules, 6)
	window.RiskyRate = ratio(window.RiskyTotal, window.Total)
	window.Severity = overallSeverity(window, maxSourceRiskRepeats)
	window.Message = aiOutputSummaryMessage(window)
	return window
}

func normalizedAIOutputDimension(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	return value
}

func normalizedAIOutputRiskType(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "none" {
		return "unknown"
	}
	return value
}

func addAIOutputAction(counts *AIOutputActionCounts, action Action) {
	switch action {
	case ActionBlock:
		counts.Block++
	case ActionSoftWarn:
		counts.SoftWarn++
	default:
		counts.Allow++
	}
}

func isRiskyAIOutputAction(action Action) bool {
	return action == ActionBlock || action == ActionSoftWarn
}

func ratio(part int, total int) float64 {
	if total <= 0 || part <= 0 {
		return 0
	}
	return math.Round((float64(part)/float64(total))*10000) / 100
}

func targetSeverity(target AIOutputTargetSummary) string {
	if target.Total == 0 {
		return "none"
	}
	if target.ActionCounts.Block >= 5 || (target.Total >= 10 && target.ActionCounts.Block > 0 && target.RiskyRate >= 10) {
		return "warning"
	}
	if target.RiskyTotal >= 3 {
		return "caution"
	}
	return "low"
}

func overallSeverity(window AIOutputSummaryWindow, maxSourceRiskRepeats int) string {
	if window.Total == 0 {
		return "none"
	}
	maxTargetBlock := 0
	maxTargetRisky := 0
	maxTargetBlockRateCritical := false
	for _, target := range window.Targets {
		if target.ActionCounts.Block > maxTargetBlock {
			maxTargetBlock = target.ActionCounts.Block
		}
		if target.RiskyTotal > maxTargetRisky {
			maxTargetRisky = target.RiskyTotal
		}
		if target.Total >= 10 && target.ActionCounts.Block > 0 && ratio(target.ActionCounts.Block, target.Total) >= 10 {
			maxTargetBlockRateCritical = true
		}
	}
	maxRuleRepeats := 0
	if len(window.MatchedRules) > 0 {
		maxRuleRepeats = window.MatchedRules[0].Count
	}
	if maxRuleRepeats >= 10 || maxTargetBlockRateCritical {
		return "critical"
	}
	if maxTargetBlock >= 5 || (window.Total >= 20 && window.RiskyRate >= 5) || maxSourceRiskRepeats >= 3 {
		return "warning"
	}
	if maxTargetRisky >= 3 || maxRuleRepeats >= 3 || maxSourceRiskRepeats >= 3 {
		return "caution"
	}
	return "low"
}

func aiOutputSummaryMessage(window AIOutputSummaryWindow) string {
	switch window.Severity {
	case "none":
		return "아직 관찰 데이터가 없습니다."
	case "critical":
		return "심각 기준에 해당하는 AI output 반복 신호가 있습니다. 후속 Phase 검토가 필요합니다."
	case "warning":
		return "경고 기준에 해당하는 AI output 반복 신호가 있습니다. target/source_feature별 확인이 필요합니다."
	case "caution":
		return "주의 기준에 해당하는 반복 신호가 있습니다. false-positive 여부를 확인하세요."
	default:
		return "현재 AI output 관찰 상태는 낮음입니다."
	}
}

func topDimensionCounts(values map[string]int, limit int) []AIOutputDimensionCount {
	items := make([]AIOutputDimensionCount, 0, len(values))
	for key, count := range values {
		items = append(items, AIOutputDimensionCount{Key: key, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Key < items[j].Key
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}

func topRuleCounts(values map[string]AIOutputRuleCount, limit int) []AIOutputRuleCount {
	items := make([]AIOutputRuleCount, 0, len(values))
	for _, item := range values {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Pattern < items[j].Pattern
		}
		return items[i].Count > items[j].Count
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}

func allowResult() *ModerationResult {
	return &ModerationResult{Action: ActionAllow, RiskType: "none"}
}

func userMessage(locale string) string {
	if strings.EqualFold(strings.TrimSpace(locale), "en") {
		return "This may not fit a learning purpose. Please rewrite it as a learning goal, task, or study record."
	}
	return "이 내용은 학습 목적과 맞지 않을 수 있어요. 학습 목표나 과제 형태로 바꿔서 다시 작성해 주세요."
}

func IsBlocked(result *ModerationResult) bool {
	if result == nil {
		return false
	}
	return result.Action == ActionBlock || result.Action == ActionSoftWarn
}

func (s *Service) ListRules(ctx context.Context, input ListRulesInput) ([]ModerationRule, error) {
	if s == nil || s.repo == nil {
		return []ModerationRule{}, nil
	}
	return s.repo.ListRules(ctx, input)
}

func (s *Service) CreateRule(ctx context.Context, input RuleMutationInput, adminUserID string) (ModerationRule, error) {
	if s == nil || s.repo == nil {
		return ModerationRule{}, nil
	}
	return s.repo.CreateRule(ctx, input, adminUserID)
}

func (s *Service) UpdateRule(ctx context.Context, id uuid.UUID, input RuleMutationInput, adminUserID string) (ModerationRule, error) {
	if s == nil || s.repo == nil {
		return ModerationRule{}, nil
	}
	return s.repo.UpdateRule(ctx, id, input, adminUserID)
}

func (s *Service) SetRuleActive(ctx context.Context, id uuid.UUID, active bool, adminUserID string) (ModerationRule, error) {
	if s == nil || s.repo == nil {
		return ModerationRule{}, nil
	}
	return s.repo.SetRuleActive(ctx, id, active, adminUserID)
}
