package curriculum

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/learnweaver/backend/internal/domain/goal"
)

func (h *Handler) ApplyDraftGoalContext(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	draftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	draftAggregate, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}

	activeGoal, err := h.goalRepo.GetActiveGoal(c.Request.Context(), draftID)
	if err == goal.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "goal not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load goal"})
		return
	}
	if activeGoal.ConfirmedGoal == nil || strings.TrimSpace(*activeGoal.ConfirmedGoal) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "confirmed_goal_required"})
		return
	}

	confirmedGoal := strings.TrimSpace(*activeGoal.ConfirmedGoal)
	appliedMode := "goal_context_synced"
	didRebuild := false
	rebuildPending := false
	message := "새 목표 문맥을 초안에 반영했고 기존 리슨 구조는 유지합니다."
	decision := goal.RebuildKeepStructure
	if activeGoal.RebuildDecision != nil {
		decision = goal.NormalizeRebuildDecision(*activeGoal.RebuildDecision)
	}
	if decision == goal.RebuildKeepStructure {
		if err := h.repo.SyncDraftGoalContext(
			c.Request.Context(),
			draftID,
			userID,
			confirmedGoal,
			activeGoal.ID,
			activeGoal.Version,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync draft goal"})
			return
		}
		appliedMode = "keep_structure"
	} else {
		switch decision {
		case goal.RebuildAll:
			cost, billingErr := h.resolveDraftBuildBilling(c.Request.Context(), userID)
			if billingErr != nil {
				if billingErr.Error() == "insufficient_points" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
					return
				}
				switch {
				case strings.Contains(billingErr.Error(), "point policy"):
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point policy"})
				case strings.Contains(billingErr.Error(), "point balance"):
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load point balance"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load ai config"})
				}
				return
			}
			req := buildRebuildDraftRequest(draftAggregate.Draft, activeGoal, confirmedGoal)
			rebuiltDecision, buildErr := h.buildDraft(c.Request.Context(), userID, req, cost)
			if buildErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rebuild draft"})
				return
			}
			normalizeRebuiltDraftAggregate(draftAggregate.Draft, rebuiltDecision.Draft)
			chargeAmount := 0
			if rebuiltDecision.ShouldCharge {
				chargeAmount = rebuiltDecision.EffectiveCost
			}
			if draftAggregate.Draft.ConfirmedCourseID != nil {
				existingCourse, courseErr := h.repo.getCourseAggregateByID(c.Request.Context(), *draftAggregate.Draft.ConfirmedCourseID)
				if courseErr != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load confirmed course"})
					return
				}
				rebuiltCourse := buildFullRebuildCourseAggregate(rebuiltDecision.Draft, existingCourse)
				if err := h.repo.ReplaceDraftAndCoursePlan(c.Request.Context(), rebuiltDecision.Draft, rebuiltCourse, chargeAmount, "rebuild_all"); err != nil {
					if err.Error() == "insufficient_points" {
						c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace draft and course plan"})
					return
				}
			} else {
				rebuiltDecision.Draft.Draft.Status = DraftStatusDraft
				if err := h.repo.ReplaceDraftPlan(c.Request.Context(), rebuiltDecision.Draft, chargeAmount); err != nil {
					if err.Error() == "insufficient_points" {
						c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient_points"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to replace draft plan"})
					return
				}
			}
			appliedMode = "rebuild_all"
			didRebuild = true
			if rebuiltDecision.BillingStatus == "byok_no_charge" {
				message = "BYOK LLM으로 새 목표 기준의 리슨을 생성했고, 기존 지역과 지점은 비활성화했습니다."
			} else {
				message = fmt.Sprintf("코스 생성과 같은 %d포인트를 차감하고 새 목표 기준의 리슨을 생성했으며, 기존 지역과 지점은 비활성화했습니다.", rebuiltDecision.EffectiveCost)
			}
		}
	}

	updatedDraft, err := h.repo.GetDraftByID(c.Request.Context(), draftID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "goal synced but failed to reload draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draft":           updatedDraft,
		"goal_profile_id": activeGoal.ID,
		"goal_version":    activeGoal.Version,
		"decision":        decision,
		"applied_mode":    appliedMode,
		"did_rebuild":     didRebuild,
		"rebuild_pending": rebuildPending,
		"message":         message,
	})
}
