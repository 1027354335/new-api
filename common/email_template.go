package common

import (
	_ "embed"
	"encoding/base64"
	"html"
	"strconv"
	"strings"
)

//go:embed assets/email-logo.png
var emailLogoPNG []byte

const emailTemplate = `<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge" />
  <title>OSS-Token — {{email_title}}</title>
  <!--[if mso]><noscript><xml><o:OfficeDocumentSettings><o:PixelsPerInch>96</o:PixelsPerInch></o:OfficeDocumentSettings></xml></noscript><![endif]-->
  <style>
    * { box-sizing: border-box; }
    body, table, td, a { -webkit-text-size-adjust: 100%; -ms-text-size-adjust: 100%; }
    table, td { mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
    img { border: 0; display: block; outline: none; text-decoration: none; }
    body {
      margin: 0; padding: 0; background-color: #f5f5f7;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Helvetica Neue',
                   'PingFang SC', 'Hiragino Kaku Gothic ProN', 'Noto Sans CJK SC', sans-serif;
    }
    @media only screen and (max-width: 600px) {
      .email-card  { border-radius: 0 !important; }
      .email-body  { padding: 36px 28px !important; }
      .code-digits { font-size: 38px !important; letter-spacing: 8px !important; }
    }
  </style>
</head>
<body style="margin:0;padding:0;background-color:#f5f5f7;">
  <div style="display:none;max-height:0;overflow:hidden;mso-hide:all;font-size:1px;color:#f5f5f7;line-height:1px;">
    {{preheader}}
    &nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;
  </div>

  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f5f5f7;">
    <tr>
      <td align="center" style="padding: 40px 16px 48px;">
        <table class="email-card" role="presentation" width="480" cellpadding="0" cellspacing="0" border="0"
               style="max-width:480px;width:100%;background:#ffffff;border-radius:12px;
                      box-shadow:0 1px 4px rgba(0,0,0,0.08),0 4px 24px rgba(0,0,0,0.06);">
          <tr>
            <td style="padding: 28px 40px; border-bottom: 1px solid #f0f0f0;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td style="vertical-align:middle;">
                    <table role="presentation" cellpadding="0" cellspacing="0" border="0">
                      <tr>
                        <td style="width:28px;height:28px;text-align:center;vertical-align:middle;">
                          <img src="{{logo_src}}" alt="OSS-Token" width="28" height="28" style="width:28px;height:28px;display:block;border:0;outline:none;text-decoration:none;object-fit:contain;" />
                        </td>
                        <td style="padding-left:8px;vertical-align:middle;">
                          <span style="color:#18181b;font-size:15px;font-weight:600;letter-spacing:-0.2px;">OSS-Token</span>
                        </td>
                      </tr>
                    </table>
                  </td>
                  <td align="right" style="vertical-align:middle;">
                    <span style="font-size:11px;color:#71717a;font-weight:500;letter-spacing:0.3px;">{{badge}}</span>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

          {{body}}

          <tr>
            <td style="padding: 20px 40px 28px; border-top: 1px solid #f0f0f0;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td>
                    <span style="font-size:12px;color:#a1a1aa;">OSS-Energie-Technik</span>
                  </td>
                  <td align="right">
                    <a href="https://www.oss-energietechnik.de/" style="font-size:12px;color:#a1a1aa;text-decoration:none;margin-left:16px;">官网</a>
                    <a href="mailto:info@oss-energietechnik.de" style="font-size:12px;color:#a1a1aa;text-decoration:none;margin-left:16px;">联系我们</a>
                  </td>
                </tr>
              </table>
              <p style="margin:14px 0 0;font-size:11px;color:#d4d4d8;line-height:1.7;">
                © 2026 OSS-Energie-Technik. All rights reserved.
                此邮件由系统自动发送，请勿直接回复。
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`

func RenderEmailVerificationTemplate(code string) string {
	validMinutes := strconv.Itoa(VerificationValidMinutes)
	body := strings.NewReplacer(
		"{{code}}", html.EscapeString(code),
		"{{valid_minutes}}", validMinutes,
	).Replace(emailVerificationBody)
	return renderEmailTemplate("邮箱验证码", "邮箱验证", "您的 OSS-Token 验证码，请在 "+validMinutes+" 分钟内完成验证。", body)
}

func RenderEmailNoticeTemplate(title string, content string, actionText string, actionURL string) string {
	body := strings.NewReplacer(
		"{{title}}", html.EscapeString(title),
		"{{content}}", normalizeEmailContent(content),
		"{{action}}", renderEmailAction(actionText, actionURL),
	).Replace(emailNoticeBody)
	return renderEmailTemplate(title, "账户通知", title, body)
}

func renderEmailTemplate(title string, badge string, preheader string, body string) string {
	return strings.NewReplacer(
		"{{email_title}}", html.EscapeString(title),
		"{{preheader}}", html.EscapeString(preheader),
		"{{badge}}", html.EscapeString(badge),
		"{{logo_src}}", emailLogoDataURI(),
		"{{body}}", body,
	).Replace(emailTemplate)
}

func emailLogoDataURI() string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(emailLogoPNG)
}

func normalizeEmailContent(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}
	if strings.Contains(content, "<") && strings.Contains(content, ">") {
		return content
	}
	content = html.EscapeString(content)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\n", "<br/>")
	return content
}

func renderEmailAction(actionText string, actionURL string) string {
	actionText = strings.TrimSpace(actionText)
	actionURL = strings.TrimSpace(actionURL)
	if actionText == "" || actionURL == "" {
		return ""
	}
	return `<a href="` + html.EscapeString(actionURL) + `"
                 style="display:inline-block;background:#18181b;color:#ffffff;text-decoration:none;
                        font-size:13px;font-weight:500;padding:11px 24px;border-radius:6px;
                        letter-spacing:0.1px;">` + html.EscapeString(actionText) + `</a>`
}

const emailVerificationBody = `<tr>
  <td class="email-body" style="padding: 44px 40px 40px;">
    <p style="margin:0 0 6px;font-size:12px;font-weight:500;color:#a1a1aa;letter-spacing:0.8px;text-transform:uppercase;">Verification Code</p>
    <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#18181b;letter-spacing:-0.4px;line-height:1.3;">邮箱验证码</h1>
    <p style="margin:0 0 32px;font-size:14px;color:#52525b;line-height:1.75;">
      您好，感谢使用 OSS-Token AI 中转服务。<br/>
      请使用以下验证码完成邮箱验证：
    </p>
    <div style="background:#fafafa;border:1px solid #e4e4e7;border-radius:8px;padding:30px 20px;text-align:center;margin-bottom:28px;">
      <div class="code-digits" style="font-family:'SF Mono','Fira Code','Courier New',monospace;font-size:42px;font-weight:700;letter-spacing:12px;color:#18181b;line-height:1;">{{code}}</div>
      <p style="margin:14px 0 0;font-size:12px;color:#a1a1aa;">有效期&ensp;<span style="color:#52525b;font-weight:500;">{{valid_minutes}} 分钟</span></p>
    </div>
    <div style="border-left:2px solid #e4e4e7;padding:10px 16px;margin-bottom:32px;">
      <p style="margin:0;font-size:12px;color:#a1a1aa;line-height:1.75;">
        如果您未发起此请求，请忽略本邮件，账户不会受到任何影响。<br/>
        切勿将验证码告知他人，OSS-Token 不会主动索取验证码。
      </p>
    </div>
    <a href="https://www.oss-energietechnik.de/" style="display:inline-block;background:#18181b;color:#ffffff;text-decoration:none;font-size:13px;font-weight:500;padding:11px 24px;border-radius:6px;letter-spacing:0.1px;">前往控制台</a>
  </td>
</tr>`

const emailNoticeBody = `<tr>
  <td class="email-body" style="padding: 44px 40px 40px;">
    <p style="margin:0 0 6px;font-size:12px;font-weight:500;color:#a1a1aa;letter-spacing:0.8px;text-transform:uppercase;">Account Notice</p>
    <h1 style="margin:0 0 16px;font-size:24px;font-weight:700;color:#18181b;letter-spacing:-0.4px;line-height:1.3;">{{title}}</h1>
    <p style="margin:0 0 28px;font-size:14px;color:#52525b;line-height:1.75;">{{content}}</p>
    {{action}}
  </td>
</tr>`
