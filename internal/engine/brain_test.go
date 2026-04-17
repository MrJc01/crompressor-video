package engine

import (
	"os"
	"testing"
)

func TestBrainLearn_HashConsistency(t *testing.T) {
	b := &AgnosticBrain{}
	tensorA := []float64{0.1, 0.5, 0.9}
	tensorB := []float64{0.1, 0.5, 0.9}
	tensorC := []float64{0.2, 0.5, 0.9}

	hashA := b.Learn(tensorA)
	hashB := b.Learn(tensorB)
	hashC := b.Learn(tensorC)

	if hashA != hashB {
		t.Errorf("Inconsistência SHA256: Tensores idênticos geraram UUIDs desiguais (%d vs %d)", hashA, hashB)
	}
	if hashA == hashC {
		t.Errorf("Colisão CROM fatal: Tensores diferentes geraram hash idêntico (%d)", hashA)
	}

	// Validação de retenção
	if len(b.Memory) != 2 {
		t.Errorf("Vazamento de memória do Cérebro. Esperado: 2. Obtido: %d", len(b.Memory))
	}
}

func TestBrainMatchForced_Accuracy(t *testing.T) {
	b := &AgnosticBrain{}
	id1 := b.Learn([]float64{1.0, 1.0, 1.0})
	id2 := b.Learn([]float64{0.0, 0.0, 0.0})

	// Tensor similar ao ID1 com ruído 0.1
	unseenData := []float64{0.9, 0.9, 0.9}
	matched := b.MatchForced(unseenData)

	if matched != id1 {
		t.Errorf("Qualimetria MSE Falhou. Esperava dar match no %d, mas atracou no %d", id1, matched)
	}

	// Tensor similar ao ID2 
	unseenDark := []float64{0.2, 0.1, 0.0}
	matchedDark := b.MatchForced(unseenDark)
	
	if matchedDark != id2 {
		t.Errorf("Qualimetria MSE Falhou no teste escuro. Atracou errôneamente em %d", matchedDark)
	}
}

func TestBrain_SaveLoad(t *testing.T) {
	b1 := &AgnosticBrain{}
	id := b1.Learn([]float64{3.14, 1.61, 2.71})

	tmpFile := "test_cerebro_tmp.gob"
	defer os.Remove(tmpFile)

	if err := b1.Save(tmpFile); err != nil {
		t.Fatalf("Erro ao serializar Cérebro O(1): %v", err)
	}

	b2 := &AgnosticBrain{}
	if err := b2.Load(tmpFile); err != nil {
		t.Fatalf("Erro ao explodir disco em R.A.M (.Load): %v", err)
	}

	val, ok := b2.Memory[id]
	if !ok {
		t.Errorf("Cérebro carregado teve amnesia. UUID %d despareceu", id)
	}
	if val[0] != 3.14 {
		t.Errorf("Corrupção de Byte Endian ao ler Cérebro persistido. Ligeiro erro MSE: %v", val[0])
	}
}
