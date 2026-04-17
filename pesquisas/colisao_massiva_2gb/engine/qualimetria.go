package engine

import (
	"math"
)

// CalcMSE calcula o Mean Squared Error de duas fatias float64.
func CalcMSE(original, decoded []float64) float64 {
	limit := len(original)
	if len(decoded) < limit {
		limit = len(decoded)
	}
	if limit == 0 {
		return math.MaxFloat64
	}

	var sum float64
	for i := 0; i < limit; i++ {
		diff := original[i] - decoded[i]
		sum += diff * diff
	}
	return sum / float64(limit)
}

// CalcPSNR calcula o Peak Signal-to-Noise Ratio (dB) assumindo que as amostras estão entre 0.0 e 1.0.
func CalcPSNR(mse float64) float64 {
	if mse == 0 {
		return 100.0 // Valor perfeito limitador empírico
	}
	return 10 * math.Log10(1.0/mse)
}
