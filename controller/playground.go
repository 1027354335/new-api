package controller

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseBodyWriter) WriteHeaderNow() {
	// Defer writing headers until we finalized the body size
}

func (w *responseBodyWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.body.Write(b)
}

func (w *responseBodyWriter) WriteString(s string) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.body.WriteString(s)
}

func (w *responseBodyWriter) Status() int {
	if w.statusCode != 0 {
		return w.statusCode
	}
	return w.ResponseWriter.Status()
}

func (w *responseBodyWriter) Size() int {
	return w.body.Len()
}

func (w *responseBodyWriter) Written() bool {
	return w.body.Len() > 0
}

func (w *responseBodyWriter) Flush() {
	// Defer flushing headers and body until response modification is completed
}

func Playground(c *gin.Context) {
	playgroundRelay(c, types.RelayFormatOpenAI)
}

func PlaygroundImage(c *gin.Context) {
	playgroundRelay(c, types.RelayFormatOpenAIImage)
}

func playgroundRelay(c *gin.Context, relayFormat types.RelayFormat) {
	var newAPIError *types.NewAPIError

	defer func() {
		if newAPIError != nil {
			c.JSON(newAPIError.StatusCode, gin.H{
				"error": newAPIError.ToOpenAIError(),
			})
		}
	}()

	useAccessToken := c.GetBool("use_access_token")
	if useAccessToken {
		newAPIError = types.NewError(errors.New("暂不支持使用 access token"), types.ErrorCodeAccessDenied, types.ErrOptionWithSkipRetry())
		return
	}

	relayInfo, err := relaycommon.GenRelayInfo(c, relayFormat, nil, nil)
	if err != nil {
		newAPIError = types.NewError(err, types.ErrorCodeInvalidRequest, types.ErrOptionWithSkipRetry())
		return
	}

	userId := c.GetInt("id")

	// Write user context to ensure acceptUnsetRatio is available
	userCache, err := model.GetUserCache(userId)
	if err != nil {
		newAPIError = types.NewError(err, types.ErrorCodeQueryDataError, types.ErrOptionWithSkipRetry())
		return
	}
	userCache.WriteContext(c)

	tempToken := &model.Token{
		UserId: userId,
		Name:   fmt.Sprintf("playground-%s", relayInfo.UsingGroup),
		Group:  relayInfo.UsingGroup,
	}
	_ = middleware.SetupContextForToken(c, tempToken)

	var originalWriter gin.ResponseWriter
	var bodyWriter *responseBodyWriter
	isImageStorageEnabled := relayFormat == types.RelayFormatOpenAIImage && storage_setting.GetStorageSetting().Enabled

	if isImageStorageEnabled {
		originalWriter = c.Writer
		bodyWriter = &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = bodyWriter
	}

	Relay(c, relayFormat)

	if isImageStorageEnabled {
		// Restore writer
		c.Writer = originalWriter

		statusCode := bodyWriter.Status()
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		if statusCode == http.StatusOK && bodyWriter.body.Len() > 0 {
			var resp dto.ImageResponse
			err := common.Unmarshal(bodyWriter.body.Bytes(), &resp)
			if err == nil {
				modified := false
				for i := range resp.Data {
					item := &resp.Data[i]
					if item.Url != "" && !strings.HasPrefix(item.Url, "/api/") {
						localPath, uploadErr := service.UploadPlaygroundImageFromURL(c.Request.Context(), userId, item.Url)
						if uploadErr == nil {
							item.Url = service.FormatPlaygroundImageURL(localPath)
							modified = true
						} else {
							common.SysError("failed to upload image URL to minio: " + uploadErr.Error())
						}
					} else if item.B64Json != "" {
						localPath, uploadErr := service.UploadPlaygroundImageFromBase64(c.Request.Context(), userId, item.B64Json)
						if uploadErr == nil {
							item.Url = service.FormatPlaygroundImageURL(localPath)
							item.B64Json = "" // Clear base64 payload to reduce size and force client load from url
							modified = true
						} else {
							common.SysError("failed to upload base64 image to minio: " + uploadErr.Error())
						}
					}
				}
				if modified {
					newData, marshalErr := common.Marshal(resp)
					if marshalErr == nil {
						originalWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(newData)))
						originalWriter.WriteHeader(statusCode)
						_, _ = originalWriter.Write(newData)
						return
					}
				}
			} else {
				common.SysError("failed to unmarshal image response for storage: " + err.Error())
			}
			// Fallback: write original response
			originalWriter.Header().Set("Content-Length", fmt.Sprintf("%d", bodyWriter.body.Len()))
			originalWriter.WriteHeader(statusCode)
			_, _ = originalWriter.Write(bodyWriter.body.Bytes())
		} else {
			if bodyWriter.body.Len() > 0 {
				originalWriter.Header().Set("Content-Length", fmt.Sprintf("%d", bodyWriter.body.Len()))
				originalWriter.WriteHeader(statusCode)
				_, _ = originalWriter.Write(bodyWriter.body.Bytes())
			} else if statusCode != 0 {
				originalWriter.WriteHeader(statusCode)
			}
		}
	}
}
