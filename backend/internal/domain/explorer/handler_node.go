package explorer

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	contentdomain "github.com/learnweaver/backend/internal/domain/content"
	"github.com/learnweaver/backend/internal/pkg/embedder"
	"github.com/learnweaver/backend/internal/pkg/logsafe"
	"github.com/learnweaver/backend/internal/pkg/metaparser"
)

func (h *Handler) CreateExplorationNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateExplorationNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID, err := h.resolveParentCourseID(c, req.ParentKind, req.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	if err := h.service.ValidateExplorationNodeSource(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node := Node{
		ParentKind: req.ParentKind,
		ParentID:   req.ParentID,
		NodeType:   NodeTypeExploration,
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
		Status:     ItemStatusActive,
		BlockCount: 0,
		SourceType: &req.SourceType,
		SourceURL:  req.SourceURL,
		ContentID:  req.ContentID,
	}

	created, err := h.repo.InsertNodeWithLinkedPoint(c.Request.Context(), courseID, node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exploration node"})
		return
	}

	// URL 기반 탐험지점: 비동기로 메타데이터 파싱 → contents 행 upsert → 임베딩 생성 → content_id 연결
	if (req.SourceType == SourceTypeYoutube || req.SourceType == SourceTypeWeb) &&
		req.SourceURL != nil && strings.TrimSpace(*req.SourceURL) != "" {
		go h.indexExternalContent(created.ID, userID, *req.SourceURL, req.SourceType)
	}

	c.JSON(http.StatusCreated, gin.H{"node": created})
}

// indexExternalContent — URL 메타데이터를 contents 테이블에 upsert하고 임베딩을 생성한 뒤
// explorer_nodes.content_id를 업데이트한다. 저작권 안전: title/description/thumbnail URL만 저장.
func (h *Handler) indexExternalContent(nodeID, userID uuid.UUID, rawURL string, sourceType SourceType) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := metaparser.Parse(ctx, rawURL, h.youtubeAPIKey)
	if err != nil {
		log.Printf("[explorer.index] meta parse 실패 node=%s url=%s err=%s", nodeID, logsafe.URL(rawURL), logsafe.Error(err))
		return
	}

	title := strings.TrimSpace(info.Title)
	if title == "" {
		title = rawURL
	}

	var desc *string
	if d := strings.TrimSpace(info.Description); d != "" {
		desc = &d
	}
	var thumb *string
	if t := strings.TrimSpace(info.ThumbnailURL); t != "" {
		thumb = &t
	}
	var author *string
	if a := strings.TrimSpace(info.Author); a != "" {
		author = &a
	}
	urlVal := rawURL

	contentType := contentdomain.ContentType(info.ContentType)
	req := contentdomain.CreateContentRequest{
		ContentType:  contentType,
		URL:          &urlVal,
		Title:        title,
		Description:  desc,
		ThumbnailURL: thumb,
		Author:       author,
		Language:     "ko",
	}

	content, err := h.contentRepo.UpsertExternalContent(ctx, userID, req)
	if err != nil {
		log.Printf("[explorer.index] content upsert 실패 node=%s err=%v", nodeID, err)
		return
	}

	// content_id를 explorer_nodes에 반영
	node, err := h.repo.GetNodeByID(ctx, nodeID)
	if err != nil {
		log.Printf("[explorer.index] node 조회 실패 node=%s err=%v", nodeID, err)
		return
	}

	courseID, err := h.resolveParentCourseIDFromBackground(ctx, node.ParentKind, node.ParentID)
	if err != nil {
		log.Printf("[explorer.index] node parent course 조회 실패 node=%s err=%v", nodeID, err)
		return
	}

	if err := h.repo.SyncNodeContentWithLinkedPoint(ctx, courseID, nodeID, content.ID); err != nil {
		log.Printf("[explorer.index] node content_id 업데이트 실패 node=%s err=%v", nodeID, err)
		return
	}

	descStr := ""
	if desc != nil {
		descStr = *desc
	}
	provider := "embedding_gemma"
	model := embeddingGemmaModel()
	text := embedder.BuildText(title, descStr)
	client, err := embedder.NewProviderClient(provider, strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_ENDPOINT")))
	if err != nil {
		log.Printf("[explorer.index] 임베딩 client 생성 실패 content=%s err=%v", content.ID, err)
		_ = h.contentRepo.SaveEmbeddingFailure(ctx, content.ID, provider, model, embeddingGemmaDimension(), text, err)
		return
	}
	vector, err := client.Embed(ctx, text)
	if err != nil {
		log.Printf("[explorer.index] 임베딩 실패 content=%s err=%v", content.ID, err)
		_ = h.contentRepo.SaveEmbeddingFailure(ctx, content.ID, provider, model, embeddingGemmaDimension(), text, err)
		return
	}
	if err := h.contentRepo.SaveEmbeddingReady(ctx, content.ID, provider, model, text, vector); err != nil {
		log.Printf("[explorer.index] embedding 저장 실패 content=%s err=%v", content.ID, err)
		return
	}
	log.Printf("[explorer.index] 완료 provider=%s model=%s node=%s content=%s dim=%d", provider, model, nodeID, content.ID, len(vector))
}

func embeddingGemmaModel() string {
	model := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_MODEL"))
	if model == "" {
		return "embedding-gemma"
	}
	return model
}

func embeddingGemmaDimension() int {
	raw := strings.TrimSpace(os.Getenv("EMBEDDING_GEMMA_DIMENSION"))
	if raw == "" {
		return 768
	}
	dimension, err := strconv.Atoi(raw)
	if err != nil || dimension <= 0 {
		return 768
	}
	return dimension
}

func (h *Handler) CreateResearchNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req CreateResearchNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseID, err := h.resolveParentCourseID(c, req.ParentKind, req.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	node := Node{
		ParentKind: req.ParentKind,
		ParentID:   req.ParentID,
		NodeType:   NodeTypeResearch,
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
		Status:     ItemStatusActive,
		BlockCount: 0,
	}

	created, err := h.repo.InsertNodeWithLinkedPoint(c.Request.Context(), courseID, node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create research node"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"node": created})
}

func (h *Handler) GetNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	node, err := h.repo.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	courseID, err := h.resolveParentCourseID(c, node.ParentKind, node.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": node})
}

func (h *Handler) UpdateNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	node, err := h.repo.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	courseID, err := h.resolveParentCourseID(c, node.ParentKind, node.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	var req UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parentKind := node.ParentKind
	parentID := node.ParentID
	if req.ParentKind != nil {
		parentKind = *req.ParentKind
	}
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	targetCourseID, err := h.resolveParentCourseID(c, parentKind, parentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}
	if targetCourseID != courseID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "node parent must stay within the same course"})
		return
	}

	if node.NodeType == NodeTypeExploration {
		sourceType := node.SourceType
		sourceURL := node.SourceURL
		contentID := node.ContentID
		if req.SourceType != nil {
			sourceType = req.SourceType
		}
		if req.SourceURL != nil {
			sourceURL = req.SourceURL
		}
		if req.ContentID != nil {
			contentID = req.ContentID
		}
		if sourceType != nil {
			validateReq := CreateExplorationNodeRequest{
				ParentKind: parentKind,
				ParentID:   parentID,
				Title:      node.Title,
				SourceType: *sourceType,
				SourceURL:  sourceURL,
				ContentID:  contentID,
			}
			if req.Title != nil {
				validateReq.Title = *req.Title
			}
			if err := h.service.ValidateExplorationNodeSource(validateReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	updated, err := h.repo.UpdateNodeWithLinkedPoint(c.Request.Context(), courseID, nodeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": updated})
}

func (h *Handler) DeleteNode(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	node, err := h.repo.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	courseID, err := h.resolveParentCourseID(c, node.ParentKind, node.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	preLearning, ok := h.isPreLearningCourse(c, courseID)
	if !ok {
		return
	}

	if preLearning {
		err = h.repo.DeleteNodeAndLinkedDraftPoint(c.Request.Context(), courseID, nodeID)
	} else {
		err = h.repo.UpdateNodeStatus(c.Request.Context(), nodeID, ItemStatusInactive)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) UpdateNodeStatus(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	node, err := h.repo.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	courseID, err := h.resolveParentCourseID(c, node.ParentKind, node.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == ItemStatusActive {
		switch node.ParentKind {
		case ParentKindRegion:
			region, err := h.repo.GetRegionByID(c.Request.Context(), node.ParentID)
			if err != nil || region.Status == ItemStatusInactive {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrParentInactive.Error()})
				return
			}
		case ParentKindSubRegion:
			subregion, err := h.repo.GetSubRegionByID(c.Request.Context(), node.ParentID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrParentInactive.Error()})
				return
			}
			region, err := h.repo.GetRegionByID(c.Request.Context(), subregion.RegionID)
			if err != nil || subregion.Status == ItemStatusInactive || region.Status == ItemStatusInactive {
				c.JSON(http.StatusBadRequest, gin.H{"error": ErrParentInactive.Error()})
				return
			}
		}
	}

	if err := h.repo.UpdateNodeStatus(c.Request.Context(), nodeID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update node status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) SetResearchNodeType(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	nodeID, err := uuid.Parse(c.Param("nodeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid node id"})
		return
	}

	node, err := h.repo.GetNodeByID(c.Request.Context(), nodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	if node.NodeType != NodeTypeResearch {
		c.JSON(http.StatusBadRequest, gin.H{"error": "node is not a research node"})
		return
	}

	courseID, err := h.resolveParentCourseID(c, node.ParentKind, node.ParentID)
	if err != nil {
		h.writeResolveParentError(c, err)
		return
	}

	if !h.verifyCourseOwner(c, courseID, userID) {
		return
	}

	var req SetResearchNodeTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := UpdateNodeRequest{
		ResearchType: &req.ResearchType,
		LayoutType:   &req.LayoutType,
	}

	updated, err := h.repo.UpdateNodeWithLinkedPoint(c.Request.Context(), courseID, nodeID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set research node type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"node": updated})
}
