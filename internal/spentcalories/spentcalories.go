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
	// Количество полей в строке данных.
	trainingFields = 3
)

// trainingInfoFormat задаёт шаблон отчёта о тренировке.
const trainingInfoFormat = `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`

var (
	// errSplitString возвращается, если строку не удалось разбить на три части.
	errSplitString = errors.New("failed to split string")
	// errStepsCount возвращается, если количество шагов не положительно.
	errStepsCount = errors.New("steps count must be positive")
	// errDuration возвращается, если длительность не положительна.
	errDuration = errors.New("duration must be positive")
	// errHeight возвращается, если рост не положителен.
	errHeight = errors.New("height must be positive")
	// errWeight возвращается, если вес не положителен.
	errWeight = errors.New("weight must be positive")
	// errUnknownTraining возвращается для неизвестного типа тренировки.
	errUnknownTraining = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	arr := strings.Split(data, ",")
	if len(arr) != trainingFields {
		return 0, "", 0, errSplitString
	}

	steps, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("parsing steps count: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errStepsCount
	}

	duration, err := time.ParseDuration(arr[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("parsing duration: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errDuration
	}

	return steps, arr[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return stepLengthCoefficient * height * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
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
		log.Println(errUnknownTraining)
		return "", errUnknownTraining
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
		return 0, errStepsCount
	}
	if weight <= 0 {
		return 0, errWeight
	}
	if height <= 0 {
		return 0, errHeight
	}
	if duration <= 0 {
		return 0, errDuration
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / float64(minInH), nil
}

// WalkingSpentCalories рассчитывает калории, потраченные при ходьбе, по количеству
// шагов, весу, росту и длительности нагрузки.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningCalories, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0, err
	}

	return walkingCaloriesCoefficient * runningCalories, nil
}
