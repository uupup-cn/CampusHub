package handler

import (
	"net/url"

	"github.com/campushub/chb-backend/internal/idp"
	"github.com/campushub/chb-backend/pkg/errcode"
	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type OAuthHandler struct {
	svc *idp.IdpService
}

func NewOAuthHandler(svc *idp.IdpService) *OAuthHandler {
	return &OAuthHandler{svc: svc}
}

// GET /oauth/authorize
func (h *OAuthHandler) Authorize(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")
	responseType := c.Query("response_type")
	scope := c.Query("scope")
	state := c.Query("state")

	if clientID == "" || redirectURI == "" || responseType != "code" || state == "" {
		redirectError(c, redirectURI, "invalid_request", state)
		return
	}

	// 开发/测试环境通过 X-User-ID 模拟登录态；生产接入论坛 Session 后由 auth 中间件注入
	userID := getUserID(c)
	if userID == 0 {
		// Redirect to login page
		redirectError(c, redirectURI, "login_required", state)
		return
	}

	trustLevel := getTrustLevel(c)

	result := h.svc.Authorize(&idp.AuthRequest{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		State:       state,
		UserID:      userID,
		TrustLevel:  trustLevel,
	})

	if result.Error != "" {
		redirectError(c, result.RedirectURI, result.Error, state)
		return
	}

	c.Redirect(302, result.RedirectURI)
}

// POST /api/oauth/authorize/confirm
// 授权页用户点击“同意”后调用，校验并发放授权码，返回应跳转的 redirect_uri。
func (h *OAuthHandler) Confirm(c *gin.Context) {
	var req struct {
		ClientID     string `json:"client_id"`
		RedirectURI  string `json:"redirect_uri"`
		ResponseType string `json:"response_type"`
		Scope        string `json:"scope"`
		State        string `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParamInvalid)
		return
	}
	if req.ClientID == "" || req.RedirectURI == "" || req.ResponseType != "code" || req.State == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}

	result := h.svc.Authorize(&idp.AuthRequest{
		ClientID:    req.ClientID,
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
		State:       req.State,
		UserID:      userID,
		TrustLevel:  getTrustLevel(c),
	})
	if result.Error != "" {
		response.ErrorWithMessage(c, errcode.ErrPermissionDenied, result.Error)
		return
	}
	response.Success(c, gin.H{
		"redirect_uri": result.RedirectURI,
		"code":         result.Code,
		"state":        result.State,
	})
}

// POST /oauth/token
func (h *OAuthHandler) Token(c *gin.Context) {
	var req idp.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try form data
		req.GrantType = c.PostForm("grant_type")
		req.Code = c.PostForm("code")
		req.RedirectURI = c.PostForm("redirect_uri")
		req.ClientID = c.PostForm("client_id")
		req.ClientSecret = c.PostForm("client_secret")
		req.RefreshToken = c.PostForm("refresh_token")
	}

	if req.GrantType == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	result, err := h.svc.Token(&req)
	if err != nil {
		if ec, ok := err.(errcode.ErrorCode); ok {
			response.Error(c, ec)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}

	c.JSON(200, result)
}

// POST /oauth/introspect
func (h *OAuthHandler) Introspect(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Token = c.PostForm("token")
	}
	if req.Token == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	result, err := h.svc.Introspect(req.Token)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}

	c.JSON(200, result)
}

// GET /api/oauth/app-info
func (h *OAuthHandler) AppInfo(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		response.Error(c, errcode.ErrParamMissing)
		return
	}

	info, err := h.svc.GetAppInfo(clientID)
	if err != nil {
		response.Error(c, errcode.ErrNotFound)
		return
	}
	response.Success(c, info)
}

func redirectError(c *gin.Context, redirectURI, err, state string) {
	params := url.Values{}
	params.Set("error", err)
	if state != "" {
		params.Set("state", state)
	}
	c.Redirect(302, redirectURI+"?"+params.Encode())
}
