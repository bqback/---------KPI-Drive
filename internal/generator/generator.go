package generator

import (
	"math/rand"
	"messagequeue/internal/config"
	"messagequeue/internal/pkg/entities"
	"strconv"
	"time"
)

func GenerateFacts(config *config.AppConfig) []entities.Fact {
	facts := make([]entities.Fact, config.GeneratorPreset.Facts)
	preset := config.GeneratorPreset
	for i := range config.GeneratorPreset.Facts {
		facts[i] = entities.Fact{
			PeriodStart: preset.PeriodStart,
			PeriodEnd:   preset.PeriodEnd,
			PeriodKey:   preset.PeriodKey,
			MoID:        preset.MoID,
			MoFactID:    preset.MoFactID,
			Value:       strconv.Itoa(i + rand.Intn(1997)),
			FactTime:    time.Now().Format(config.DateFormat),
			IsPlan:      preset.IsPlan,
			AuthUserID:  preset.AuthUserID,
			Comment:     "buffer Архаров",
		}
	}
	return facts
}
