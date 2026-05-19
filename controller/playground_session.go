package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

type playgroundSessionRequest struct {
	Id            int             `json:"id"`
	Title         string          `json:"title"`
	Messages      model.JSONValue `json:"messages"`
	SelectedImage model.JSONValue `json:"selected_image"`
}

func normalizePlaygroundSessionTitle(title string, messages model.JSONValue) string {
	title = strings.TrimSpace(title)
	if title != "" {
		if len([]rune(title)) > 128 {
			return string([]rune(title)[:128])
		}
		return title
	}

	var parsed []struct {
		From     string `json:"from"`
		Versions []struct {
			Content string `json:"content"`
		} `json:"versions"`
	}
	if len(messages) > 0 {
		_ = common.Unmarshal(messages, &parsed)
	}
	for _, message := range parsed {
		if message.From != "user" || len(message.Versions) == 0 {
			continue
		}
		content := strings.TrimSpace(message.Versions[0].Content)
		if content == "" {
			continue
		}
		runes := []rune(content)
		if len(runes) > 36 {
			return string(runes[:36]) + "..."
		}
		return content
	}
	return "Untitled conversation"
}

func normalizePlaygroundJSON(value model.JSONValue, fallback string) model.JSONValue {
	if len(value) == 0 {
		return model.JSONValue([]byte(fallback))
	}
	return value
}

func ListPlaygroundSessions(c *gin.Context) {
	userId := c.GetInt("id")
	pageInfo := common.GetPageQuery(c)
	sessions, err := model.GetUserPlaygroundSessions(userId, pageInfo.GetStartIdx(), pageInfo.GetPageSize())
	if err != nil {
		common.ApiError(c, err)
		return
	}
	total, err := model.CountUserPlaygroundSessions(userId)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// Format all sessions' relative paths to full URLs for the frontend
	for i := range sessions {
		sessions[i].Messages = service.ProcessPlaygroundJSONUrls(sessions[i].Messages, service.FormatPlaygroundImageURL)
		sessions[i].SelectedImage = service.ProcessPlaygroundJSONUrls(sessions[i].SelectedImage, service.FormatPlaygroundImageURL)
	}

	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(sessions)
	common.ApiSuccess(c, pageInfo)
}

func GetPlaygroundSession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	session, err := model.GetUserPlaygroundSessionByID(id, c.GetInt("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}

	// Format relative paths to full URLs for the frontend
	session.Messages = service.ProcessPlaygroundJSONUrls(session.Messages, service.FormatPlaygroundImageURL)
	session.SelectedImage = service.ProcessPlaygroundJSONUrls(session.SelectedImage, service.FormatPlaygroundImageURL)

	common.ApiSuccess(c, session)
}

func CreatePlaygroundSession(c *gin.Context) {
	userId := c.GetInt("id")
	var req playgroundSessionRequest
	if err := common.UnmarshalBodyReusable(c, &req); err != nil {
		common.ApiError(c, err)
		return
	}
	req.Messages = normalizePlaygroundJSON(req.Messages, "[]")
	req.SelectedImage = normalizePlaygroundJSON(req.SelectedImage, "null")

	// Strip bucket and prefix from the messages and selected_image before storing in DB
	strippedMessages := service.ProcessPlaygroundJSONUrls(req.Messages, service.StripPlaygroundImageURL)
	strippedSelectedImage := service.ProcessPlaygroundJSONUrls(req.SelectedImage, service.StripPlaygroundImageURL)

	session := &model.PlaygroundSession{
		UserId:        userId,
		Title:         normalizePlaygroundSessionTitle(req.Title, req.Messages),
		Messages:      strippedMessages,
		SelectedImage: strippedSelectedImage,
	}
	if err := session.Insert(); err != nil {
		common.ApiError(c, err)
		return
	}

	// Format bucket and prefix back for the response to the frontend
	session.Messages = service.ProcessPlaygroundJSONUrls(session.Messages, service.FormatPlaygroundImageURL)
	session.SelectedImage = service.ProcessPlaygroundJSONUrls(session.SelectedImage, service.FormatPlaygroundImageURL)

	common.ApiSuccess(c, session)
}

func UpdatePlaygroundSession(c *gin.Context) {
	userId := c.GetInt("id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	var req playgroundSessionRequest
	if err := common.UnmarshalBodyReusable(c, &req); err != nil {
		common.ApiError(c, err)
		return
	}
	session, err := model.GetUserPlaygroundSessionByID(id, userId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	req.Messages = normalizePlaygroundJSON(req.Messages, "[]")
	req.SelectedImage = normalizePlaygroundJSON(req.SelectedImage, "null")

	// Strip bucket and prefix before saving to DB
	strippedMessages := service.ProcessPlaygroundJSONUrls(req.Messages, service.StripPlaygroundImageURL)
	strippedSelectedImage := service.ProcessPlaygroundJSONUrls(req.SelectedImage, service.StripPlaygroundImageURL)

	session.Title = normalizePlaygroundSessionTitle(req.Title, req.Messages)
	session.Messages = strippedMessages
	session.SelectedImage = strippedSelectedImage
	if err := session.Update(); err != nil {
		common.ApiError(c, err)
		return
	}

	// Format back for the response
	session.Messages = service.ProcessPlaygroundJSONUrls(session.Messages, service.FormatPlaygroundImageURL)
	session.SelectedImage = service.ProcessPlaygroundJSONUrls(session.SelectedImage, service.FormatPlaygroundImageURL)

	common.ApiSuccess(c, session)
}

func DeletePlaygroundSession(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if err := model.DeleteUserPlaygroundSessionByID(id, c.GetInt("id")); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, gin.H{
		"message": fmt.Sprintf("deleted playground session %d", id),
	})
}
