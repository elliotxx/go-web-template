package handler

import (
	"fmt"
	"time"
)

type Response struct {
	Success   bool      `json:"success" yaml:"success"`
	Code      string    `json:"code" yaml:"code"`
	Message   string    `json:"message" yaml:"message"`
	Data      any       `json:"data,omitempty" yaml:"data,omitempty"`
	TraceID   string    `json:"traceID,omitempty" yaml:"traceID,omitempty"`
	StartTime time.Time `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty" yaml:"endTime,omitempty"`
	CostTime  Duration  `json:"costTime,omitempty" yaml:"costTime,omitempty"`
}

type Duration time.Duration

func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, (time.Duration(d)).String())), nil
}
