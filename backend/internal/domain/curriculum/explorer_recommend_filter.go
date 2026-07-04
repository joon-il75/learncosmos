package curriculum

import "github.com/learnweaver/backend/internal/pkg/search/normalizer"

type explorerFilterStage struct {
	requiredTokens []string
	strongTokens   []string
	minStrongMatch int
}

func filterExplorerCandidatesByContext(
	candidates []ContentSearchCandidate,
	queryBundle normalizer.LessonSearchQuery,
	limit int,
) []ContentSearchCandidate {
	return filterExplorerCandidatesByStage(candidates, buildExplorerFilterStages(queryBundle), limit)
}

func filterExplorerCandidatesByStage(
	candidates []ContentSearchCandidate,
	stages []explorerFilterStage,
	limit int,
) []ContentSearchCandidate {
	if len(candidates) == 0 {
		return candidates
	}
	if limit <= 0 {
		limit = len(candidates)
	}

	nonEmptyStages := make([]explorerFilterStage, 0, len(stages))
	for _, stage := range stages {
		if len(stage.requiredTokens) == 0 {
			continue
		}
		nonEmptyStages = append(nonEmptyStages, stage)
	}
	if len(nonEmptyStages) == 0 {
		return []ContentSearchCandidate{}
	}

	filtered := make([]ContentSearchCandidate, 0, minInt(limit, len(candidates)))
	seen := make(map[string]struct{}, len(candidates))

	for _, stage := range nonEmptyStages {
		for _, candidate := range candidates {
			key := explorerCandidateFilterKey(candidate)
			if _, exists := seen[key]; exists {
				continue
			}

			searchText := normalizer.BuildSearchTextKO(
				candidate.Title,
				derefString(candidate.Description),
				derefString(candidate.ExternalURL),
			)
			if !matchesExplorerCandidateContext(searchText, stage.strongTokens, stage.requiredTokens, stage.minStrongMatch) {
				continue
			}

			seen[key] = struct{}{}
			filtered = append(filtered, candidate)
			if len(filtered) >= limit {
				return filtered
			}
		}
	}

	return filtered
}

func filterExplorerCandidatesByAllStages(
	candidates []ContentSearchCandidate,
	stages []explorerFilterStage,
	limit int,
) []ContentSearchCandidate {
	if len(candidates) == 0 {
		return candidates
	}
	if limit <= 0 {
		limit = len(candidates)
	}

	nonEmptyStages := make([]explorerFilterStage, 0, len(stages))
	for _, stage := range stages {
		if len(stage.requiredTokens) == 0 {
			continue
		}
		nonEmptyStages = append(nonEmptyStages, stage)
	}
	if len(nonEmptyStages) == 0 {
		return []ContentSearchCandidate{}
	}

	filtered := make([]ContentSearchCandidate, 0, minInt(limit, len(candidates)))
	seen := make(map[string]struct{}, len(candidates))

	for _, candidate := range candidates {
		key := explorerCandidateFilterKey(candidate)
		if _, exists := seen[key]; exists {
			continue
		}

		searchText := normalizer.BuildSearchTextKO(
			candidate.Title,
			derefString(candidate.Description),
			derefString(candidate.ExternalURL),
		)

		matchesAll := true
		for _, stage := range nonEmptyStages {
			if !matchesExplorerCandidateContext(searchText, stage.strongTokens, stage.requiredTokens, stage.minStrongMatch) {
				matchesAll = false
				break
			}
		}
		if !matchesAll {
			continue
		}

		seen[key] = struct{}{}
		filtered = append(filtered, candidate)
		if len(filtered) >= limit {
			break
		}
	}

	return filtered
}

func buildExplorerSearchStages(nodeBundle, levelBundle, courseBundle normalizer.LessonSearchQuery) []explorerSearchStage {
	inputs := []struct {
		name   string
		bundle normalizer.LessonSearchQuery
	}{
		{name: "node", bundle: nodeBundle},
		{name: "level", bundle: levelBundle},
		{name: "course", bundle: courseBundle},
	}

	stages := make([]explorerSearchStage, 0, len(inputs))
	for _, input := range inputs {
		filter, ok := buildExplorerFilterStage(input.name, input.bundle)
		if !ok {
			continue
		}
		stages = append(stages, explorerSearchStage{
			name:   input.name,
			bundle: input.bundle,
			filter: filter,
		})
	}
	return stages
}
