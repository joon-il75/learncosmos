package admin

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type lumiRuntimeRuleDraftPayload struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Trigger     string `json:"trigger"`
	Mode        string `json:"mode"`
	Page        string `json:"page"`
	Action      string `json:"action"`
	Scene       string `json:"scene"`
	Context     string `json:"context"`
	DockSlot    string `json:"dockSlot"`
	State       string `json:"state"`
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
	Note        string `json:"note"`
}

type lumiRuntimeMessageDraftPayload struct {
	ID          string `json:"id"`
	Trigger     string `json:"trigger"`
	Label       string `json:"label"`
	Note        string `json:"note"`
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
	Preview     string `json:"preview"`
}

type lumiRuntimeConfigFile struct {
	Rules    []lumiRuntimeRuleDraftPayload    `json:"rules"`
	Messages []lumiRuntimeMessageDraftPayload `json:"messages"`
}

type lumiRuntimeConfigResponse struct {
	Exists    bool                             `json:"exists"`
	Rules     []lumiRuntimeRuleDraftPayload    `json:"rules"`
	Messages  []lumiRuntimeMessageDraftPayload `json:"messages"`
	UpdatedAt string                           `json:"updated_at,omitempty"`
	Version   int64                            `json:"version,omitempty"`
}

type publicLumiRuntimeRuleResponse struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Trigger     string `json:"trigger"`
	Mode        string `json:"mode"`
	Page        string `json:"page"`
	Action      string `json:"action"`
	Scene       string `json:"scene"`
	Context     string `json:"context"`
	DockSlot    string `json:"dockSlot"`
	State       string `json:"state"`
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}

type publicLumiRuntimeMessageResponse struct {
	ID          string `json:"id"`
	Trigger     string `json:"trigger"`
	Label       string `json:"label"`
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
	Preview     string `json:"preview"`
}

type publicLumiRuntimeConfigResponse struct {
	Rules    []publicLumiRuntimeRuleResponse    `json:"rules"`
	Messages []publicLumiRuntimeMessageResponse `json:"messages"`
}

func (h *AdminHandler) getPublicLumiRuntimeConfigHandler(c *gin.Context) {
	response, err := readLumiRuntimeConfigResponse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read lumi runtime config"})
		return
	}

	c.JSON(http.StatusOK, buildPublicLumiRuntimeConfigResponse(response))
}

func (h *AdminHandler) getLumiRuntimeConfigHandler(c *gin.Context) {
	h.writeLumiRuntimeConfigResponse(c)
}

func (h *AdminHandler) updateLumiRuntimeConfigHandler(c *gin.Context) {
	var req lumiRuntimeConfigFile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	payload, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode lumi runtime config"})
		return
	}

	if err := writeAtomicJSONFile(lumiRuntimeConfigFilePath(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist lumi runtime config"})
		return
	}

	response, err := readLumiRuntimeConfigResponse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reload lumi runtime config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "lumi runtime config saved",
		"config":  response,
	})
}

func (h *AdminHandler) writeLumiRuntimeConfigResponse(c *gin.Context) {
	response, err := readLumiRuntimeConfigResponse()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read lumi runtime config"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func buildPublicLumiRuntimeConfigResponse(response lumiRuntimeConfigResponse) publicLumiRuntimeConfigResponse {
	publicRules := make([]publicLumiRuntimeRuleResponse, 0, len(response.Rules))
	for _, rule := range response.Rules {
		publicRules = append(publicRules, publicLumiRuntimeRuleResponse{
			ID:          rule.ID,
			Label:       rule.Label,
			Trigger:     rule.Trigger,
			Mode:        rule.Mode,
			Page:        rule.Page,
			Action:      rule.Action,
			Scene:       rule.Scene,
			Context:     rule.Context,
			DockSlot:    rule.DockSlot,
			State:       rule.State,
			MessageType: rule.MessageType,
			Message:     rule.Message,
		})
	}

	publicMessages := make([]publicLumiRuntimeMessageResponse, 0, len(response.Messages))
	for _, message := range response.Messages {
		publicMessages = append(publicMessages, publicLumiRuntimeMessageResponse{
			ID:          message.ID,
			Trigger:     message.Trigger,
			Label:       message.Label,
			MessageType: message.MessageType,
			Message:     message.Message,
			Preview:     message.Preview,
		})
	}

	return publicLumiRuntimeConfigResponse{Rules: publicRules, Messages: publicMessages}
}

func readLumiRuntimeConfigResponse() (lumiRuntimeConfigResponse, error) {
	info, err := os.Stat(lumiRuntimeConfigFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			return lumiRuntimeConfigResponse{
				Exists:   false,
				Rules:    []lumiRuntimeRuleDraftPayload{},
				Messages: []lumiRuntimeMessageDraftPayload{},
			}, nil
		}
		return lumiRuntimeConfigResponse{}, err
	}

	payload, err := os.ReadFile(lumiRuntimeConfigFilePath())
	if err != nil {
		return lumiRuntimeConfigResponse{}, err
	}

	var config lumiRuntimeConfigFile
	if err := json.Unmarshal(payload, &config); err != nil {
		return lumiRuntimeConfigResponse{}, err
	}

	if config.Rules == nil {
		config.Rules = []lumiRuntimeRuleDraftPayload{}
	}
	if config.Messages == nil {
		config.Messages = []lumiRuntimeMessageDraftPayload{}
	}

	return lumiRuntimeConfigResponse{
		Exists:    true,
		Rules:     config.Rules,
		Messages:  config.Messages,
		UpdatedAt: info.ModTime().Format(time.RFC3339),
		Version:   info.ModTime().Unix(),
	}, nil
}

func writeAtomicJSONFile(target string, payload []byte) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	tmpPath := target + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return err
	}

	return os.Rename(tmpPath, target)
}
