package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/campushub/chb-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	discourseURL string
}

func NewAuthHandler() *AuthHandler {
	url := os.Getenv("DISCOURSE_URL")
	if url == "" {
		url = "http://127.0.0.1:9800"
	}
	return &AuthHandler{discourseURL: url}
}

// Me checks if the user is logged into Discourse by forwarding the session cookie.
// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	cookie := c.GetHeader("Cookie")

	req, err := http.NewRequest("GET", h.discourseURL+"/session/current.json", nil)
	if err != nil {
		response.Success(c, gin.H{"logged_in": false})
		return
	}

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		response.Success(c, gin.H{"logged_in": false})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		response.Success(c, gin.H{"logged_in": false})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Success(c, gin.H{"logged_in": false})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		response.Success(c, gin.H{"logged_in": false})
		return
	}

	if user, ok := result["current_user"]; ok {
		if userMap, ok := user.(map[string]interface{}); ok {
			response.Success(c, gin.H{
				"logged_in": true,
				"user_id":   userMap["id"],
				"username":  userMap["username"],
				"trust_level": userMap["trust_level"],
				"avatar_template": userMap["avatar_template"],
			})
			return
		}
	}

	response.Success(c, gin.H{"logged_in": false})
}
