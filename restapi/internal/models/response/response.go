package response

import (
	"restapi/internal/models/state"

	"github.com/gin-gonic/gin"
)

type Response struct {
	State state.State    `json:"state"`
	Data  map[string]any `json:"data,omitempty"`
}

func Ok(c *gin.Context, code int, data map[string]any) {
	c.JSON(code, Response{
		State: state.OK(),
		Data:  data,
	})
}

func Error(c *gin.Context, code int, errorMsg string) {
	c.JSON(code, Response{
		State: state.Error(errorMsg),
	})
}
