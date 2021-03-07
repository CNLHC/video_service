package utils

import (
	"math"
	"time"
)

type Estimator interface {
	EstimatePercentage() float32
}

type ExpotentialEstimator struct {
	Lambda  float32
	StartAt time.Time
}

func (c *ExpotentialEstimator) EstimatePercentage() float32 {
	diff := time.Now().Sub(c.StartAt)
	return 100.0 * float32(1-float32(math.Exp(float64(diff.Seconds())*float64(c.Lambda))))
}

func GetDefaultEstimator() Estimator {
	return &ExpotentialEstimator{
		StartAt: time.Now(),
		Lambda:  -0.03837641,
	}
}
