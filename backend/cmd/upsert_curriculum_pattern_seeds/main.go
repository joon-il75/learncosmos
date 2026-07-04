package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/domain/curriculum"
	"github.com/learnweaver/backend/internal/pkg/db"
)

func main() {
	_ = godotenv.Load()

	var (
		seedPath                    = flag.String("seed", defaultSeedPath(), "curriculum pattern seed JSON path")
		provider                    = flag.String("provider", "embedding_gemma", "pattern embedding provider")
		model                       = flag.String("model", defaultString("EMBEDDING_GEMMA_MODEL", "embedding-gemma"), "pattern embedding model")
		dimension                   = flag.Int("dimension", defaultInt("EMBEDDING_GEMMA_DIMENSION", 768), "pattern embedding dimension")
		dryRun                      = flag.Bool("dry-run", false, "validate and summarize seed file without writing")
		skipEmbeddingPlaceholder    = flag.Bool("skip-embedding-placeholder", false, "do not create pending curriculum_pattern_embeddings rows")
		staleOldEmbeddingPlaceholds = flag.Bool("mark-stale-old-embeddings", true, "mark old pending/failed embeddings stale when embedding_text_hash changes")
	)
	flag.Parse()

	ctx := context.Background()
	seeds, err := readSeeds(*seedPath)
	if err != nil {
		log.Fatalf("read seeds: %v", err)
	}
	if err := validateSeeds(seeds); err != nil {
		log.Fatalf("validate seeds: %v", err)
	}
	log.Printf("loaded curriculum pattern seeds=%d path=%s", len(seeds), *seedPath)

	if *dryRun {
		for _, seed := range seeds {
			log.Printf("[seed] key=%s version=%s language=%s hash=%s title=%s", seed.PatternKey, seed.PatternVersion, seed.Language, embeddingTextHash(seed.EmbeddingText), seed.Title)
		}
		return
	}

	if *dimension <= 0 {
		log.Fatal("-dimension must be greater than 0")
	}
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	result, err := upsertSeeds(ctx, pool, seeds, strings.TrimSpace(*provider), strings.TrimSpace(*model), *dimension, !*skipEmbeddingPlaceholder, *staleOldEmbeddingPlaceholds)
	if err != nil {
		log.Fatalf("upsert seeds: %v", err)
	}
	log.Printf(
		"curriculum pattern seed upsert completed patterns=%d language_counts=%v embedding_placeholders=%d stale_embeddings=%d provider=%s model=%s dimension=%d",
		result.Patterns,
		result.LanguageCounts,
		result.EmbeddingPlaceholders,
		result.StaleEmbeddings,
		*provider,
		*model,
		*dimension,
	)
}

type upsertResult struct {
	Patterns              int
	EmbeddingPlaceholders int
	StaleEmbeddings       int64
	LanguageCounts        map[string]int
}

func readSeeds(path string) ([]curriculum.CurriculumPatternSeed, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("seed path is required")
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var seeds []curriculum.CurriculumPatternSeed
	if err := json.Unmarshal(payload, &seeds); err != nil {
		return nil, err
	}
	return seeds, nil
}

func validateSeeds(seeds []curriculum.CurriculumPatternSeed) error {
	if len(seeds) == 0 {
		return fmt.Errorf("no seeds found")
	}
	seen := map[string]struct{}{}
	for _, seed := range seeds {
		key := strings.TrimSpace(seed.PatternKey)
		version := strings.TrimSpace(seed.PatternVersion)
		if key == "" {
			return fmt.Errorf("empty pattern_key")
		}
		if version == "" {
			return fmt.Errorf("empty pattern_version for key=%s", key)
		}
		if strings.TrimSpace(seed.Title) == "" {
			return fmt.Errorf("empty title for key=%s version=%s", key, version)
		}
		if strings.TrimSpace(seed.EmbeddingText) == "" {
			return fmt.Errorf("empty embedding_text for key=%s version=%s", key, version)
		}
		language := normalizeSeedLanguage(seed.Language)
		if language == "" {
			return fmt.Errorf("invalid language for key=%s version=%s: %q", key, version, seed.Language)
		}
		uniqueKey := key + "\x00" + version + "\x00" + language
		if _, ok := seen[uniqueKey]; ok {
			return fmt.Errorf("duplicate pattern key/version/language: %s %s %s", key, version, language)
		}
		seen[uniqueKey] = struct{}{}
	}
	return nil
}

func upsertSeeds(ctx context.Context, pool *pgxpool.Pool, seeds []curriculum.CurriculumPatternSeed, provider, model string, dimension int, createEmbeddingPlaceholder, markStaleOldEmbeddings bool) (upsertResult, error) {
	if provider == "" {
		provider = "embedding_gemma"
	}
	if model == "" {
		model = "embedding-gemma"
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return upsertResult{}, err
	}
	defer tx.Rollback(ctx)

	result := upsertResult{LanguageCounts: map[string]int{}}
	for _, seed := range seeds {
		patternID, err := upsertPattern(ctx, tx, seed)
		if err != nil {
			return result, err
		}
		result.Patterns++
		result.LanguageCounts[normalizeSeedLanguage(defaultIfBlank(seed.Language, "ko"))]++

		if !createEmbeddingPlaceholder {
			continue
		}
		hash := embeddingTextHash(seed.EmbeddingText)
		if markStaleOldEmbeddings {
			tag, err := tx.Exec(ctx, `
				UPDATE curriculum_pattern_embeddings
				SET status = 'stale',
				    updated_at = NOW()
				WHERE pattern_id = $1
				  AND provider = $2
				  AND model = $3
				  AND dimension = $4
				  AND embedding_text_hash <> $5
				  AND status IN ('pending', 'ready', 'failed')
			`, patternID, provider, model, dimension, hash)
			if err != nil {
				return result, err
			}
			result.StaleEmbeddings += tag.RowsAffected()
		}
		if err := upsertEmbeddingPlaceholder(ctx, tx, patternID, provider, model, dimension, hash); err != nil {
			return result, err
		}
		result.EmbeddingPlaceholders++
	}

	if err := tx.Commit(ctx); err != nil {
		return result, err
	}
	return result, nil
}

func upsertPattern(ctx context.Context, tx pgx.Tx, seed curriculum.CurriculumPatternSeed) (string, error) {
	triggerKeywords, err := json.Marshal(seed.TriggerKeywords)
	if err != nil {
		return "", fmt.Errorf("marshal trigger_keywords key=%s: %w", seed.PatternKey, err)
	}
	guidance, err := json.Marshal(seed.Guidance)
	if err != nil {
		return "", fmt.Errorf("marshal guidance key=%s: %w", seed.PatternKey, err)
	}
	stageRules, err := json.Marshal(seed.StageRules)
	if err != nil {
		return "", fmt.Errorf("marshal stage_rules key=%s: %w", seed.PatternKey, err)
	}
	recommendedSequence, err := json.Marshal(seed.RecommendedSequence)
	if err != nil {
		return "", fmt.Errorf("marshal recommended_sequence key=%s: %w", seed.PatternKey, err)
	}
	badPatterns, err := json.Marshal(seed.BadPatterns)
	if err != nil {
		return "", fmt.Errorf("marshal bad_patterns key=%s: %w", seed.PatternKey, err)
	}
	searchSpecTemplate, err := json.Marshal(seed.RecommendationSearchSpecTemplate)
	if err != nil {
		return "", fmt.Errorf("marshal recommendation_search_spec_template key=%s: %w", seed.PatternKey, err)
	}
	source, err := json.Marshal(seed.Source)
	if err != nil {
		return "", fmt.Errorf("marshal source key=%s: %w", seed.PatternKey, err)
	}

	var id string
	err = tx.QueryRow(ctx, `
		INSERT INTO curriculum_patterns (
			pattern_key, pattern_version, domain, goal_type, learner_level, language,
			title, summary, trigger_keywords, guidance, stage_rules, recommended_sequence,
			bad_patterns, recommendation_search_spec_template, embedding_text, source, is_active, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9::jsonb, $10::jsonb, $11::jsonb, $12::jsonb,
			$13::jsonb, $14::jsonb, $15, $16::jsonb, $17, NOW(), NOW()
		)
		ON CONFLICT (pattern_key, pattern_version, language)
		DO UPDATE SET
			domain = EXCLUDED.domain,
			goal_type = EXCLUDED.goal_type,
			learner_level = EXCLUDED.learner_level,
			language = EXCLUDED.language,
			title = EXCLUDED.title,
			summary = EXCLUDED.summary,
			trigger_keywords = EXCLUDED.trigger_keywords,
			guidance = EXCLUDED.guidance,
			stage_rules = EXCLUDED.stage_rules,
			recommended_sequence = EXCLUDED.recommended_sequence,
			bad_patterns = EXCLUDED.bad_patterns,
			recommendation_search_spec_template = EXCLUDED.recommendation_search_spec_template,
			embedding_text = EXCLUDED.embedding_text,
			source = EXCLUDED.source,
			is_active = EXCLUDED.is_active,
			updated_at = NOW()
		RETURNING id::text
	`,
		seed.PatternKey,
		defaultIfBlank(seed.PatternVersion, "v1"),
		seed.Domain,
		seed.GoalType,
		defaultIfBlank(seed.LearnerLevel, "any"),
		normalizeSeedLanguage(defaultIfBlank(seed.Language, "ko")),
		seed.Title,
		seed.Summary,
		string(triggerKeywords),
		string(guidance),
		string(stageRules),
		string(recommendedSequence),
		string(badPatterns),
		string(searchSpecTemplate),
		seed.EmbeddingText,
		string(source),
		seed.IsActive,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("upsert pattern key=%s version=%s language=%s: %w", seed.PatternKey, seed.PatternVersion, seed.Language, err)
	}
	return id, nil
}

func normalizeSeedLanguage(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "ko", "kr", "kor", "korean":
		return "ko"
	case "en", "eng", "english":
		return "en"
	default:
		return ""
	}
}

func upsertEmbeddingPlaceholder(ctx context.Context, tx pgx.Tx, patternID, provider, model string, dimension int, hash string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO curriculum_pattern_embeddings (
			pattern_id, provider, model, dimension, embedding_text_hash,
			status, error_message, created_at, updated_at
		)
		VALUES ($1::uuid, $2, $3, $4, $5, 'pending', '', NOW(), NOW())
		ON CONFLICT (pattern_id, provider, model, dimension, embedding_text_hash)
		DO UPDATE SET
			status = CASE
				WHEN curriculum_pattern_embeddings.status = 'ready' THEN curriculum_pattern_embeddings.status
				ELSE 'pending'
			END,
			error_message = CASE
				WHEN curriculum_pattern_embeddings.status = 'ready' THEN curriculum_pattern_embeddings.error_message
				ELSE ''
			END,
			updated_at = NOW()
	`, patternID, provider, model, dimension, hash)
	return err
}

func embeddingTextHash(text string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(sum[:])
}

func defaultSeedPath() string {
	return filepath.Clean("../docs/data/curriculum-patterns/seed_curriculum_patterns_v1.json")
}

func defaultIfBlank(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func defaultString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func defaultInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
