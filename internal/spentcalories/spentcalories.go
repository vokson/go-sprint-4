// Package spentcalories рассчитывает дистанцию, скорость и потраченные калории
// по количеству шагов, весу и росту.
package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// trainingInfoFormat задаёт шаблон отчёта о тренировке.
const trainingInfoFormat = `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`

var (
	// ErrSplitString возвращается, если строку не удалось разбить на три части.
	ErrSplitString = errors.New("ошибка разделения строки на части")
	// ErrStepsCount возвращается, если количество шагов не положительно.
	ErrStepsCount = errors.New("количество шагов <= 0")
	// ErrDuration возвращается, если длительность не положительна.
	ErrDuration = errors.New("промежуток времени <= 0")
	// ErrHeight возвращается, если рост не положителен.
	ErrHeight = errors.New("рост <= 0")
	// ErrWeight возвращается, если вес не положителен.
	ErrWeight = errors.New("вес <= 0")
	// ErrUnknownTraining возвращается для неизвестного типа тренировки.
	ErrUnknownTraining = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	arr := strings.Split(data, ",")
	if len(arr) != 3 {
		return 0, "", 0, ErrSplitString
	}

	steps, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrStepsCount
	}

	duration, err := time.ParseDuration(arr[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, ErrDuration
	}

	return steps, arr[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return stepLengthCoefficient * height * float64(steps) / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo разбирает строку с данными о тренировке (шаги, тип, длительность)
// и возвращает отчёт с дистанцией, скоростью и потраченными калориями.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var caloriesFunc func(int, float64, float64, time.Duration) (float64, error)

	switch trainingType {
	case "Бег":
		caloriesFunc = RunningSpentCalories
	case "Ходьба":
		caloriesFunc = WalkingSpentCalories
	default:
		log.Println(ErrUnknownTraining)
		return "", ErrUnknownTraining
	}

	calories, err := caloriesFunc(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return fmt.Sprintf(trainingInfoFormat, trainingType, duration.Hours(), distance, speed, calories), nil
}

// RunningSpentCalories рассчитывает калории, потраченные при беге, по количеству
// шагов, весу, росту и длительности нагрузки.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, ErrStepsCount
	}
	if weight <= 0 {
		return 0.0, ErrWeight
	}
	if height <= 0 {
		return 0.0, ErrHeight
	}
	if duration <= 0 {
		return 0.0, ErrDuration
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / float64(minInH), nil
}

// WalkingSpentCalories рассчитывает калории, потраченные при ходьбе, по количеству
// шагов, весу, росту и длительности нагрузки.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningCalories, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0.0, err
	}

	return walkingCaloriesCoefficient * runningCalories, nil
}
