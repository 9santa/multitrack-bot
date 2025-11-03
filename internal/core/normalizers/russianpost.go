package normalizers

import (
	"multitrack-bot/internal/domain"
	"sort"
	"time"
)

type RussianPostNormalizer struct{}

func NewRussianPostNormalizer() *RussianPostNormalizer {
	return &RussianPostNormalizer{}
}

func (n *RussianPostNormalizer) CanNormalize(courierName string) bool {
	return courierName == "russianpost" || courierName == "Почта России"
}

func (n *RussianPostNormalizer) Normalize(raw *domain.RawTrackingResult) *domain.TrackingResult {
	result := &domain.TrackingResult{
		Courier: raw.Courier,
	}

	if !raw.Successful {
		result.Status = "Не найдено"
		result.Description = "Информация по трек-номеру не найдена"
		return result
	}

	records, ok := raw.RawData.([]domain.HistoryRecord)
	if !ok || len(records) == 0 {
		result.Status = "Не найдено"
		result.Description = "Нет данных по отслеживанию"
		return result
	}

	firstRecord := records[0]
	result.Number = firstRecord.Barcode

	statusMap := map[string]string{
		"Вручение":  "✅ Доставлено",
		"Обработка": "🚚 В обработке",
		"Прием":     "📮 Принято в отделении",
		"Присвоение идентификатора":  "📝 Создана",
		"Покинуло место приема":      "➡️ Покинуло место приема",
		"Прибыло в место вручения":   "🏢 Прибыло в место вручения",
		"Неудачная попытка вручения": "❌ Не удалось вручить",
	}

	lastRecord := records[len(records)-1]
	if humanStatus, exists := statusMap[lastRecord.OperType]; exists {
		result.Status = humanStatus
	} else {
		result.Status = lastRecord.OperType
	}

	result.Description = lastRecord.OperType

	for _, r := range records {
		t, _ := time.Parse(time.RFC3339, r.OperDate)
		result.Checkpoints = append(result.Checkpoints, domain.Checkpoint{
			Date:        t,
			Location:    r.Address,
			Status:      r.OperAttr,
			Description: r.OperType,
		})
	}

	n.sortCheckpoints(result.Checkpoints)
	result.LastUpdated = time.Now()

	return result
}

func (n *RussianPostNormalizer) sortCheckpoints(checkpoints []domain.Checkpoint) {
	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].Date.Before(checkpoints[j].Date)
	})
}
