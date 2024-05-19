package generator

import (
	"math/rand"
	"messagequeue/internal/pkg/entities"
	"strconv"
	"time"
)

func GenerateFacts(n int, preset entities.GeneratorPreset) []entities.Fact {
	facts := []entities.Fact{}
	for i := range n {
		facts = append(facts, entities.Fact{
			PeriodStart: preset.PeriodStart,
			PeriodEnd:   preset.PeriodEnd,
			PeriodKey:   preset.PeriodKey,
			MoID:        preset.MoID,
			MoFactID:    strconv.Itoa(i),
			Value:       strconv.Itoa(i + rand.Intn(1997)),
			FactTime:    time.Now().String(),
			IsPlan:      preset.IsPlan,
			AuthUserID:  preset.AuthUserID,
			Comment:     "buffer Архаров",
		})
	}
	return facts
}
