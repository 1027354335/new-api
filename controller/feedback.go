package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
)

type createFeedbackRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Priority string `json:"priority"`
}

type adminReplyFeedbackRequest struct {
	Status string `json:"status"`
	Reply  string `json:"reply"`
}

type adminUpdateFeedbackStatusRequest struct {
	Status string `json:"status"`
}

func normalizeFeedbackRequest(req *createFeedbackRequest) {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.Category = strings.TrimSpace(req.Category)
	req.Priority = strings.TrimSpace(req.Priority)
	if req.Category == "" {
		req.Category = model.FeedbackCategoryOther
	}
	if req.Priority == "" {
		req.Priority = model.FeedbackPriorityNormal
	}
}

func validateFeedbackRequest(req createFeedbackRequest) string {
	if req.Title == "" {
		return "反馈标题不能为空"
	}
	if len([]rune(req.Title)) > 255 {
		return "反馈标题不能超过 255 个字符"
	}
	if req.Content == "" {
		return "反馈内容不能为空"
	}
	if len([]rune(req.Content)) > 5000 {
		return "反馈内容不能超过 5000 个字符"
	}
	if !model.IsValidFeedbackCategory(req.Category) {
		return "反馈分类无效"
	}
	if !model.IsValidFeedbackPriority(req.Priority) {
		return "反馈优先级无效"
	}
	return ""
}

func getFeedbackId(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.ApiErrorMsg(c, "id 参数无效")
		return 0, false
	}
	return id, true
}

func CreateFeedback(c *gin.Context) {
	var req createFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}
	normalizeFeedbackRequest(&req)
	if msg := validateFeedbackRequest(req); msg != "" {
		common.ApiErrorMsg(c, msg)
		return
	}

	userId := c.GetInt("id")
	user, err := model.GetUserById(userId, false)
	if err != nil {
		common.ApiErrorMsg(c, "用户不存在")
		return
	}

	now := common.GetTimestamp()
	feedback := &model.Feedback{
		UserId:     userId,
		Username:   user.Username,
		Email:      user.Email,
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		Priority:   req.Priority,
		Status:     model.FeedbackStatusOpen,
		CreateTime: now,
		UpdateTime: now,
	}
	if err := feedback.Insert(); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, feedback)
}

func GetMyFeedbacks(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	if status != "" && !model.IsValidFeedbackStatus(status) {
		common.ApiErrorMsg(c, "反馈状态无效")
		return
	}

	pageInfo := common.GetPageQuery(c)
	feedbacks, total, err := model.GetUserFeedbacks(c.GetInt("id"), status, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(feedbacks)
	common.ApiSuccess(c, pageInfo)
}

func GetMyFeedback(c *gin.Context) {
	id, ok := getFeedbackId(c)
	if !ok {
		return
	}
	feedback, err := model.GetUserFeedbackById(id, c.GetInt("id"))
	if err != nil {
		if model.IsFeedbackNotFound(err) {
			common.ApiErrorMsg(c, "反馈不存在")
			return
		}
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, feedback)
}

func CloseMyFeedback(c *gin.Context) {
	id, ok := getFeedbackId(c)
	if !ok {
		return
	}
	feedback, err := model.GetUserFeedbackById(id, c.GetInt("id"))
	if err != nil {
		if model.IsFeedbackNotFound(err) {
			common.ApiErrorMsg(c, "反馈不存在")
			return
		}
		common.ApiError(c, err)
		return
	}
	if feedback.Status == model.FeedbackStatusClosed {
		common.ApiSuccess(c, feedback)
		return
	}
	feedback.Status = model.FeedbackStatusClosed
	feedback.UpdateTime = common.GetTimestamp()
	if err := feedback.Update(); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, feedback)
}

func AdminListFeedbacks(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	category := strings.TrimSpace(c.Query("category"))
	priority := strings.TrimSpace(c.Query("priority"))
	keyword := strings.TrimSpace(c.Query("keyword"))

	if status != "" && !model.IsValidFeedbackStatus(status) {
		common.ApiErrorMsg(c, "反馈状态无效")
		return
	}
	if category != "" && !model.IsValidFeedbackCategory(category) {
		common.ApiErrorMsg(c, "反馈分类无效")
		return
	}
	if priority != "" && !model.IsValidFeedbackPriority(priority) {
		common.ApiErrorMsg(c, "反馈优先级无效")
		return
	}

	pageInfo := common.GetPageQuery(c)
	feedbacks, total, err := model.GetAllFeedbacks(status, category, priority, keyword, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(feedbacks)
	common.ApiSuccess(c, pageInfo)
}

func AdminGetFeedback(c *gin.Context) {
	id, ok := getFeedbackId(c)
	if !ok {
		return
	}
	feedback, err := model.GetFeedbackById(id)
	if err != nil {
		if model.IsFeedbackNotFound(err) {
			common.ApiErrorMsg(c, "反馈不存在")
			return
		}
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, feedback)
}

func AdminReplyFeedback(c *gin.Context) {
	id, ok := getFeedbackId(c)
	if !ok {
		return
	}
	var req adminReplyFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}
	req.Status = strings.TrimSpace(req.Status)
	req.Reply = strings.TrimSpace(req.Reply)
	if !model.IsValidFeedbackReplyStatus(req.Status) {
		common.ApiErrorMsg(c, "回复处理状态必须为 resolved 或 rejected")
		return
	}
	if req.Reply == "" {
		common.ApiErrorMsg(c, "回复内容不能为空")
		return
	}
	if len([]rune(req.Reply)) > 5000 {
		common.ApiErrorMsg(c, "回复内容不能超过 5000 个字符")
		return
	}

	feedback, err := model.GetFeedbackById(id)
	if err != nil {
		if model.IsFeedbackNotFound(err) {
			common.ApiErrorMsg(c, "反馈不存在")
			return
		}
		common.ApiError(c, err)
		return
	}
	if feedback.Status != model.FeedbackStatusOpen && feedback.Status != model.FeedbackStatusInProgress {
		common.ApiErrorMsg(c, "只有待处理或处理中的反馈可以回复")
		return
	}

	now := common.GetTimestamp()
	feedback.Status = req.Status
	feedback.Reply = req.Reply
	feedback.AdminId = c.GetInt("id")
	feedback.AdminName = c.GetString("username")
	feedback.UpdateTime = now
	feedback.ResolvedTime = now
	if err := feedback.Update(); err != nil {
		common.ApiError(c, err)
		return
	}

	go sendFeedbackReplyEmail(feedback)
	common.ApiSuccess(c, feedback)
}

func AdminUpdateFeedbackStatus(c *gin.Context) {
	id, ok := getFeedbackId(c)
	if !ok {
		return
	}
	var req adminUpdateFeedbackStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}
	req.Status = strings.TrimSpace(req.Status)
	if !model.IsValidFeedbackStatus(req.Status) {
		common.ApiErrorMsg(c, "反馈状态无效")
		return
	}
	if req.Status == model.FeedbackStatusResolved || req.Status == model.FeedbackStatusRejected {
		common.ApiErrorMsg(c, "请通过回复处理反馈")
		return
	}

	feedback, err := model.GetFeedbackById(id)
	if err != nil {
		if model.IsFeedbackNotFound(err) {
			common.ApiErrorMsg(c, "反馈不存在")
			return
		}
		common.ApiError(c, err)
		return
	}
	feedback.Status = req.Status
	feedback.AdminId = c.GetInt("id")
	feedback.AdminName = c.GetString("username")
	feedback.UpdateTime = common.GetTimestamp()
	if req.Status == model.FeedbackStatusClosed {
		feedback.ResolvedTime = feedback.UpdateTime
	}
	if err := feedback.Update(); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, feedback)
}

func sendFeedbackReplyEmail(feedback *model.Feedback) {
	email := strings.TrimSpace(feedback.Email)
	if email == "" {
		userEmail, err := model.GetUserEmail(feedback.UserId)
		if err != nil {
			common.SysLog(fmt.Sprintf("failed to get feedback user email: feedback_id=%d, user_id=%d, error=%v", feedback.Id, feedback.UserId, err))
			return
		}
		email = strings.TrimSpace(userEmail)
	}
	if email == "" {
		common.SysLog(fmt.Sprintf("feedback user has no email, skip notification: feedback_id=%d, user_id=%d", feedback.Id, feedback.UserId))
		return
	}

	subject := "您的问题反馈已处理"
	content := fmt.Sprintf("您好 %s，\n\n您提交的问题反馈已处理。\n\n反馈标题：%s\n处理状态：%s\n处理结果：\n%s\n\n感谢您的反馈。", feedback.Username, feedback.Title, feedback.Status, feedback.Reply)
	content = common.RenderEmailNoticeTemplate(subject, content, "前往控制台", "")
	if err := common.SendEmail(subject, email, content); err != nil {
		common.SysLog(fmt.Sprintf("failed to send feedback reply email: feedback_id=%d, user_id=%d, error=%v", feedback.Id, feedback.UserId, err))
	}
}
