// Package daysteps формирует сводку о дневной активности по количеству шагов
// и длительности прогулки.
package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах.
	stepLength = 0.65
	// Количество метров в одном километре.
	mInKm = 1000
	// Количество полей в строке данных.
	expectedFields = 2
)

// dayActionFormat задаёт шаблон отчёта о дневной активности.
const dayActionFormat = `Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`

var (
	// errSplitString возвращается, если строку не удалось разбить на две части.
	errSplitString = errors.New("failed to split string")
	// errStepsCount возвращается, если количество шагов не положительно.
	errStepsCount = errors.New("steps count must be positive")
	// errDuration возвращается, если длительность не положительна.
	errDuration = errors.New("duration must be positive")
)

func parsePackage(data string) (int, time.Duration, error) {
	pair := strings.Split(data, ",")
	if len(pair) != expectedFields {
		return 0, 0, errSplitString
	}

	steps, err := strconv.Atoi(pair[0])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing steps count: %w", err)
	}
	if steps <= 0 {
		return 0, 0, errStepsCount
	}

	duration, err := time.ParseDuration(pair[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parsing duration: %w", err)
	}
	if duration <= 0 {
		return 0, 0, errDuration
	}

	return steps, duration, nil
}

// DayActionInfo разбирает строку с количеством шагов и длительностью прогулки
// и возвращает отчёт с дистанцией и потраченными калориями.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distance := stepLength * float64(steps) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(dayActionFormat, steps, distance, calories)
}
