package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/learnweaver/backend/internal/pkg/db"
	"github.com/learnweaver/backend/internal/pkg/search/normalizer"
)

const evaluationSeedExternalSource = "learnweaver_eval_seed"

type evaluationContentSeed struct {
	ExternalID  string
	ContentType string
	URL         string
	Title       string
	Description string
	Author      string
	Language    string
}

var visualArtEvaluationContentSeeds = []evaluationContentSeed{
	{
		ExternalID:  "visual-art-watercolor-water-control-001",
		ContentType: "youtube",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-watercolor-water-control-001",
		Title:       "수채화 물조절 기초 - 붓의 수분량과 번짐 조절 연습",
		Description: "초보자가 수채화 붓에 머금은 물의 양을 조절하고, 종이 위 번짐을 비교하며 단계별 물조절을 연습합니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
	{
		ExternalID:  "visual-art-watercolor-brush-wash-002",
		ContentType: "youtube",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-watercolor-brush-wash-002",
		Title:       "초보 수채화 붓 사용법과 워시 연습 - 맑은 색 번짐 만들기",
		Description: "수채화 기초 붓 사용, 넓은 워시, 물감 농도, 밝고 어두운 색 번짐을 함께 다룹니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
	{
		ExternalID:  "visual-art-watercolor-gradient-003",
		ContentType: "article",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-watercolor-gradient-003",
		Title:       "수채화 그라데이션과 번짐 통제 체크리스트",
		Description: "물조절, 붓 압력, 농도 변화, 번짐 경계 확인을 순서대로 점검하는 수채화 기초 체크리스트입니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
	{
		ExternalID:  "visual-art-pencil-shading-001",
		ContentType: "youtube",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-pencil-shading-001",
		Title:       "연필 드로잉 명암 기초 - 밝기 단계와 톤 스케일 연습",
		Description: "연필로 밝기 단계를 나누고 명암 스케일을 만든 뒤 간단한 도형에 적용하는 드로잉 기초 연습입니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
	{
		ExternalID:  "visual-art-pencil-sphere-cylinder-002",
		ContentType: "youtube",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-pencil-sphere-cylinder-002",
		Title:       "구와 원기둥 명암 넣기 - 연필 스케치 입체감 기초",
		Description: "빛 방향을 정하고 구, 원기둥, 상자에 연필 명암을 쌓아 입체감을 만드는 드로잉 기초 연습입니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
	{
		ExternalID:  "visual-art-pencil-sketch-basics-003",
		ContentType: "article",
		URL:         "https://dev.learnweavr.com/evaluation-seeds/visual-art-pencil-sketch-basics-003",
		Title:       "연필 스케치 기초와 명암 연습 순서",
		Description: "초보 드로잉에서 선 연습, 형태 잡기, 명암 단계, 스케치 마무리까지 이어지는 연습 순서입니다.",
		Author:      "LearnWeaver Evaluation Seed",
		Language:    "ko",
	},
}

func main() {
	var dryRun bool
	var userIDRaw string
	flag.BoolVar(&dryRun, "dry-run", false, "print seed rows without writing")
	flag.StringVar(&userIDRaw, "user-id", strings.TrimSpace(os.Getenv("EVALUATION_CONTENT_SEED_USER_ID")), "owner user id for inserted seed contents")
	flag.Parse()

	loadEnv()

	ctx := context.Background()
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := db.NewPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()

	userID, err := resolveSeedOwnerUserID(ctx, pool, userIDRaw)
	if err != nil {
		log.Fatalf("failed to resolve seed owner: %v", err)
	}

	if dryRun {
		log.Printf("evaluation content seed rows=%d owner=%s dry_run=true", len(visualArtEvaluationContentSeeds), userID)
		for _, seed := range visualArtEvaluationContentSeeds {
			log.Printf("[seed] external_id=%s type=%s title=%s", seed.ExternalID, seed.ContentType, seed.Title)
		}
		return
	}

	inserted, updated, err := upsertEvaluationContentSeeds(ctx, pool, userID, visualArtEvaluationContentSeeds)
	if err != nil {
		log.Fatalf("failed to upsert evaluation content seeds: %v", err)
	}
	log.Printf("evaluation content seed upsert completed inserted=%d updated=%d owner=%s", inserted, updated, userID)
}

func loadEnv() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../backend/.env")
	_ = godotenv.Load("backend/.env")
}

func resolveSeedOwnerUserID(ctx context.Context, pool *pgxpool.Pool, raw string) (uuid.UUID, error) {
	if strings.TrimSpace(raw) != "" {
		return uuid.Parse(strings.TrimSpace(raw))
	}

	var userID uuid.UUID
	err := pool.QueryRow(ctx, `
		SELECT u.id
		FROM users u
		LEFT JOIN contents c ON c.user_id = u.id
		GROUP BY u.id
		ORDER BY COUNT(c.id) DESC, MIN(u.created_at) NULLS LAST
		LIMIT 1
	`).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("select fallback owner: %w", err)
	}
	return userID, nil
}

func upsertEvaluationContentSeeds(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, seeds []evaluationContentSeed) (int, int, error) {
	inserted := 0
	updated := 0
	for _, seed := range seeds {
		wasInserted, err := upsertEvaluationContentSeed(ctx, pool, userID, seed)
		if err != nil {
			return inserted, updated, err
		}
		if wasInserted {
			inserted++
		} else {
			updated++
		}
	}
	return inserted, updated, nil
}

func upsertEvaluationContentSeed(ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID, seed evaluationContentSeed) (bool, error) {
	if strings.TrimSpace(seed.ExternalID) == "" || strings.TrimSpace(seed.Title) == "" {
		return false, fmt.Errorf("invalid seed: external_id and title are required")
	}
	language := strings.TrimSpace(seed.Language)
	if language == "" {
		language = "ko"
	}
	searchTextKO := normalizer.BuildSearchTextKO(seed.Title, seed.Description, seed.Author, seed.ContentType, language)

	var existingID uuid.UUID
	err := pool.QueryRow(ctx, `
		SELECT id
		FROM contents
		WHERE external_source = $1 AND external_content_id = $2
		LIMIT 1
	`, evaluationSeedExternalSource, seed.ExternalID).Scan(&existingID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("select existing seed %s: %w", seed.ExternalID, err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		_, err = pool.Exec(ctx, `
			INSERT INTO contents (
				user_id, content_type, external_source, external_content_id, canonical_url, url,
				title, description, author, language, is_public, quality_score,
				content_status, health_score, search_text_ko, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $5,
				$6, $7, $8, $9, true, 0.55,
				'active', 1.0, $10, NOW(), NOW()
			)
		`, userID, seed.ContentType, evaluationSeedExternalSource, seed.ExternalID, seed.URL, seed.Title, seed.Description, seed.Author, language, searchTextKO)
		if err != nil {
			return false, fmt.Errorf("insert seed %s: %w", seed.ExternalID, err)
		}
		return true, nil
	}

	_, err = pool.Exec(ctx, `
		UPDATE contents
		SET user_id = $1,
		    content_type = $2,
		    canonical_url = $3,
		    url = $3,
		    title = $4,
		    description = $5,
		    author = $6,
		    language = $7,
		    is_public = true,
		    quality_score = 0.55,
		    content_status = 'active',
		    health_score = 1.0,
		    search_text_ko = $8,
		    updated_at = NOW()
		WHERE id = $9
	`, userID, seed.ContentType, seed.URL, seed.Title, seed.Description, seed.Author, language, searchTextKO, existingID)
	if err != nil {
		return false, fmt.Errorf("update seed %s: %w", seed.ExternalID, err)
	}
	if _, err := pool.Exec(ctx, `
		UPDATE content_embeddings
		SET status = 'stale', updated_at = NOW()
		WHERE content_id = $1
		  AND status = 'ready'
	`, existingID); err != nil {
		return false, fmt.Errorf("mark seed embedding stale %s: %w", seed.ExternalID, err)
	}
	return false, nil
}
