package curriculum

import (
	"strings"

	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

func compactExplorerRecommendationParts(parts ...string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
func BuildExplorerRecommendationSearchStages(ctx ExplorerRecommendationContext) []explorerSearchStage {
	contextSourceQuery := strings.Join(compactExplorerRecommendationParts(
		ctx.SourceQuery,
		ctx.LearningGoal,
		ctx.GoalQueryHint,
	), " ")
	if strings.TrimSpace(contextSourceQuery) == "" {
		contextSourceQuery = strings.TrimSpace(ctx.FallbackQuery)
	}
	contextLevelTitle := strings.Join(compactExplorerRecommendationParts(
		ctx.RegionTitle,
		ctx.SubRegionTitle,
	), " ")
	contextLevelObjective := strings.Join(compactExplorerRecommendationParts(
		ctx.RegionObjective,
		ctx.SubRegionSummary,
	), " ")
	contextLessonTitle := strings.TrimSpace(ctx.NodeTitle)
	contextLessonSummary := strings.TrimSpace(ctx.NodeSummary)
	if contextLessonTitle != "" && contextLessonTitle == strings.TrimSpace(ctx.FallbackQuery) {
		contextLessonTitle = ""
		contextLessonSummary = ""
	}
	return buildExplorerSearchStages(
		normalizer.BuildLessonSearchQuery("", "", "", contextLessonTitle, strings.TrimSpace(ctx.FallbackQuery), contextLessonSummary),
		normalizer.BuildLessonSearchQuery("", contextLevelTitle, contextLevelObjective, "", "", ""),
		normalizer.BuildLessonSearchQuery(contextSourceQuery, "", "", "", "", ""),
	)
}

func BuildExplorerRecommendationExternalSearchQuery(ctx ExplorerRecommendationContext) string {
	contextSourceQuery := strings.Join(compactExplorerRecommendationParts(
		ctx.SourceQuery,
		ctx.LearningGoal,
		ctx.GoalQueryHint,
	), " ")
	if strings.TrimSpace(contextSourceQuery) == "" {
		contextSourceQuery = strings.TrimSpace(ctx.FallbackQuery)
	}
	return buildExplorerExternalSearchQuery(
		BuildExplorerRecommendationSearchStages(ctx),
		contextSourceQuery,
		strings.TrimSpace(ctx.FallbackQuery),
	)
}

func BuildExplorerRecommendationExternalSearchQueryWithPrimaryQuery(ctx ExplorerRecommendationContext, primaryQuery string) string {
	contextSourceQuery := strings.Join(compactExplorerRecommendationParts(
		ctx.SourceQuery,
		ctx.LearningGoal,
		ctx.GoalQueryHint,
	), " ")
	if strings.TrimSpace(contextSourceQuery) == "" {
		contextSourceQuery = strings.TrimSpace(ctx.FallbackQuery)
	}
	return buildExplorerExternalSearchQueryWithPrimaryQuery(
		BuildExplorerRecommendationSearchStages(ctx),
		primaryQuery,
		contextSourceQuery,
		strings.TrimSpace(ctx.FallbackQuery),
	)
}

func RerankExplorerRecommendationCandidates(candidates []ContentSearchCandidate, ctx ExplorerRecommendationContext, limit int) []ContentSearchCandidate {
	return rerankExplorerCandidates(candidates, BuildExplorerRecommendationSearchStages(ctx), limit)
}

func RerankExplorerRecommendationCandidatesWithPrimaryQuery(candidates []ContentSearchCandidate, ctx ExplorerRecommendationContext, primaryQuery string, limit int) []ContentSearchCandidate {
	trimmedPrimary := strings.TrimSpace(primaryQuery)
	if trimmedPrimary == "" {
		return RerankExplorerRecommendationCandidates(candidates, ctx, limit)
	}
	return rerankExplorerCandidatesWithPrimaryQuery(candidates, BuildExplorerRecommendationSearchStages(ctx), trimmedPrimary, limit)
}

func buildExplorerExternalSearchQuery(stages []explorerSearchStage, fallbackParts ...string) string {
	parts := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	hasPriorityCourseToken := false

	appendToken := func(token string) bool {
		trimmed := normalizeExplorerFilterToken(token)
		if trimmed == "" || isWeakExplorerFilterToken(trimmed) {
			return false
		}
		if _, exists := seen[trimmed]; exists {
			return false
		}
		seen[trimmed] = struct{}{}
		parts = append(parts, trimmed)
		return true
	}

	stageTokenLimit := func(name string) int {
		switch name {
		case "node":
			return 2
		case "level":
			return 2
		case "course":
			return 1
		default:
			return 2
		}
	}

	appendStageTokens := func(stage explorerSearchStage) bool {
		addedForStage := 0
		for _, token := range stage.filter.strongTokens {
			if stage.name == "course" && !isPriorityExplorerCourseToken(token) {
				continue
			}
			if appendToken(token) {
				addedForStage++
				if stage.name == "course" {
					hasPriorityCourseToken = true
				}
			}
			if len(parts) >= 6 || addedForStage >= stageTokenLimit(stage.name) {
				return len(parts) >= 6
			}
		}
		return len(parts) >= 6
	}

	for _, stageName := range []string{"course", "node", "level"} {
		for _, stage := range stages {
			if stage.name != stageName {
				continue
			}
			if appendStageTokens(stage) {
				break
			}
		}
		if len(parts) >= 6 {
			break
		}
	}

	if !hasPriorityCourseToken && len(parts) < 5 {
		for _, fallback := range fallbackParts {
			for _, token := range normalizer.BuildLessonSearchQuery(fallback, "", "", "", "", "").Tokens {
				appendToken(token)
				if len(parts) >= 6 {
					break
				}
			}
			if len(parts) >= 6 {
				break
			}
		}
	}

	if len(parts) == 0 {
		return strings.TrimSpace(strings.Join(compactExplorerRecommendationParts(fallbackParts...), " "))
	}
	return strings.Join(parts, " ")
}

func buildExplorerExternalSearchQueryWithPrimaryQuery(stages []explorerSearchStage, primaryQuery string, fallbackParts ...string) string {
	trimmedPrimary := strings.TrimSpace(primaryQuery)
	if trimmedPrimary == "" {
		return buildExplorerExternalSearchQuery(stages, fallbackParts...)
	}
	normalized := strings.Join(strings.Fields(trimmedPrimary), " ")
	runes := []rune(normalized)
	if len(runes) > 120 {
		return string(runes[:120])
	}
	return normalized
}
