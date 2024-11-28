package ftracker

import (
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)
	return distance / duration
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки(Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки в часах.
// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки (Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
// lengthPool int — длина бассейна в метрах.
// countPool int — количество заплывов.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	var dist, speed, calories float64

	switch trainingType {
	case "Бег":
		dist = distance(action)                                   // расчет дистанции
		speed = meanSpeed(action, duration)                       // расчет средней скорости
		calories = RunningSpentCalories(action, weight, duration) // расчет сожженных калорий
	case "Ходьба":
		dist = distance(action)                                           // расчет дистанции
		speed = meanSpeed(action, duration)                               // расчет средней скорости
		calories = WalkingSpentCalories(action, duration, weight, height) // расчет сожженных калорий
	case "Плавание":
		dist = distance(action)                                                   // расчет дистанции
		speed = swimmingMeanSpeed(lengthPool, countPool, duration)                // расчет средней скорости
		calories = SwimmingSpentCalories(lengthPool, countPool, duration, weight) // расчет сожженных калорий
	default:
		return "неизвестный тип тренировки"
	}

	// Формируем и возвращаем строку с информацией о тренировке
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, dist, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 — вес пользователя.
// duration float64 — длительность тренировки в часах.
func RunningSpentCalories(action int, weight, duration float64) float64 {
	// Проверяем, чтобы длительность не была равна нулю
	if duration <= 0 {
		return 0
	}

	// Расчет средней скорости
	averageSpeed := meanSpeed(action, duration)

	// Расчет калорий с использованием формулы
	runningCalories := ((runningCaloriesMeanSpeedMultiplier*averageSpeed + runningCaloriesMeanSpeedShift) * weight * duration * minInH)

	return runningCalories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	// Проверка на нулевую длительность
	if duration <= 0 {
		return 0
	}

	// Расчет средней скорости в метрах в секунду (предполагаем, что длина шага 0.65 метра)
	averageSpeed := (float64(action) * lenStep) / (duration * 3600) // переводим часы в секунды

	// Расчет калорий по формуле
	calories := ((walkingCaloriesWeightMultiplier * weight) + ((averageSpeed * averageSpeed) / height * walkingSpeedHeightMultiplier * weight)) * duration * minInH

	return calories
}

// Константы для расчета калорий, расходуемых при плавании.
const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых колорий при плавании относительно скорости.
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании.
)

// swimmingMeanSpeed возвращает среднюю скорость при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	// Проверка на нулевую длительность
	if duration <= 0 {
		return 0
	}

	// Расчет средней скорости плавания (в км/ч)
	totalDistance := float64(lengthPool*countPool) / 1000 // переводим метры в километры
	averageSpeed := totalDistance / duration              // средняя скорость в км/ч

	// Формула для расчета сожженных калорий
	swimmingCalories := (averageSpeed + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration

	return swimmingCalories
}
