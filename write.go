package wlgologger

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

func Write(level string, payload any) {
	result, _ := json.Marshal(payload)
	log.Println(level + " " + string(result))
}

type httpLogStructure struct {
	RequestUUID string `json:"request_uuid"`
	Payload     any    `json:"payload"`
}

func WriteHttp(c *gin.Context, level string, payload any) {
	httpLog := httpLogStructure{
		RequestUUID: c.MustGet("request_uuid").(string),
		Payload:     payload,
	}
	result, _ := json.Marshal(httpLog)
	log.Println(level + " " + string(result))
}