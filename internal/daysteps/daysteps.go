// Package daysteps формирует сводку о дневной активности по количеству шагов
// и длительности прогулки.
package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// dayActionFormat задаёт шаблон отчёта о дневной активности.
const dayActionFormat = `Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`

var (
	// ErrSplitString возвращается, если строку не удалось разбить на две части.
	ErrSplitString = errors.New("ошибка разделения строки на части")
	// ErrStepsCount возвращается, если количество шагов не положительно.
	ErrStepsCount = errors.New("количество шагов <= 0")
	// ErrDuration возвращается, если длительность не положительна.
	ErrDuration = errors.New("промежуток времени <= 0")
)

func parsePackage(data string) (int, time.Duration, error) {
	pair := strings.Split(data, ",")
	if len(pair) != 2 {
		return 0, 0, ErrSplitString
	}

	steps, err := strconv.Atoi(pair[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, ErrStepsCount
	}

	duration, err := time.ParseDuration(pair[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, ErrDuration
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

	distance := stepLength * float64(steps) / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf(dayActionFormat, steps, distance, calories)
}
