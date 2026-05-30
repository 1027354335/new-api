package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func insertFeedbackForTest(t *testing.T, feedback *Feedback) {
	t.Helper()
	require.NoError(t, DB.Create(feedback).Error)
}

func TestFeedbackValidationHelpers(t *testing.T) {
	assert.True(t, IsValidFeedbackStatus(FeedbackStatusOpen))
	assert.True(t, IsValidFeedbackStatus(FeedbackStatusInProgress))
	assert.True(t, IsValidFeedbackReplyStatus(FeedbackStatusResolved))
	assert.False(t, IsValidFeedbackReplyStatus(FeedbackStatusInProgress))
	assert.True(t, IsValidFeedbackCategory(FeedbackCategoryBilling))
	assert.False(t, IsValidFeedbackCategory("unknown"))
	assert.True(t, IsValidFeedbackPriority(FeedbackPriorityUrgent))
	assert.False(t, IsValidFeedbackPriority("unknown"))
}

func TestGetUserFeedbacksFiltersByOwnerAndStatus(t *testing.T) {
	truncateTables(t)

	insertFeedbackForTest(t, &Feedback{
		UserId:     1,
		Username:   "alice",
		Title:      "open issue",
		Content:    "content",
		Category:   FeedbackCategoryBug,
		Priority:   FeedbackPriorityNormal,
		Status:     FeedbackStatusOpen,
		CreateTime: 100,
		UpdateTime: 100,
	})
	insertFeedbackForTest(t, &Feedback{
		UserId:     1,
		Username:   "alice",
		Title:      "resolved issue",
		Content:    "content",
		Category:   FeedbackCategoryBilling,
		Priority:   FeedbackPriorityHigh,
		Status:     FeedbackStatusResolved,
		CreateTime: 200,
		UpdateTime: 200,
	})
	insertFeedbackForTest(t, &Feedback{
		UserId:     2,
		Username:   "bob",
		Title:      "other user issue",
		Content:    "content",
		Category:   FeedbackCategoryBug,
		Priority:   FeedbackPriorityUrgent,
		Status:     FeedbackStatusOpen,
		CreateTime: 300,
		UpdateTime: 300,
	})

	pageInfo := commonPageInfoForFeedbackTest()
	feedbacks, total, err := GetUserFeedbacks(1, FeedbackStatusOpen, &pageInfo)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, feedbacks, 1)
	assert.Equal(t, "open issue", feedbacks[0].Title)
}

func TestGetAllFeedbacksFiltersAndSearchesKeyword(t *testing.T) {
	truncateTables(t)

	insertFeedbackForTest(t, &Feedback{
		UserId:     1,
		Username:   "alice",
		Email:      "alice@example.com",
		Title:      "billing problem",
		Content:    "charged twice",
		Category:   FeedbackCategoryBilling,
		Priority:   FeedbackPriorityHigh,
		Status:     FeedbackStatusOpen,
		CreateTime: 100,
		UpdateTime: 100,
	})
	insertFeedbackForTest(t, &Feedback{
		UserId:     2,
		Username:   "bob",
		Email:      "bob@example.com",
		Title:      "model question",
		Content:    "latency",
		Category:   FeedbackCategoryModel,
		Priority:   FeedbackPriorityLow,
		Status:     FeedbackStatusClosed,
		CreateTime: 200,
		UpdateTime: 200,
	})

	pageInfo := commonPageInfoForFeedbackTest()
	feedbacks, total, err := GetAllFeedbacks(
		FeedbackStatusOpen,
		FeedbackCategoryBilling,
		FeedbackPriorityHigh,
		"charged",
		&pageInfo,
	)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, feedbacks, 1)
	assert.Equal(t, "billing problem", feedbacks[0].Title)
}

func commonPageInfoForFeedbackTest() common.PageInfo {
	return common.PageInfo{
		Page:     1,
		PageSize: 10,
	}
}
