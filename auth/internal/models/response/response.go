package response

import (
	"github.com/ayayaakasvin/auth/internal/models/data"
	"github.com/ayayaakasvin/auth/internal/models/state"
	"github.com/gin-gonic/gin"
)

type Response struct {
	State 	state.State 	`json:"state"`
	Data 	data.Data 		`json:"data,omitempty"`
}

func Ok (c *gin.Context, code int, data data.Data)  {
	c.JSON(code, Response{
		State: state.Ok(),
		Data: data,
	})
}

func Error (c *gin.Context, code int, errorMsg string)  {
	c.JSON(code, Response{
		State: state.Error(errorMsg),
	})
}