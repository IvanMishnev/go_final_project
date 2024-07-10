package nextdate

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/IvanMishnev/go_final_project/internal/constants"
)

func NextDate(now time.Time, d string, repeat string) (string, error) {
	date, err := time.Parse(constants.DateFormat, d)
	if err != nil {
		return "", errors.New("wrong date format")
	}

	repeatRules := strings.Fields(repeat)
	if len(repeatRules) == 0 {
		return "", errors.New("wrong repeat rule format")
	}

	rule := repeatRules[0]
	switch rule {
	case "d":
		if len(repeatRules) != 2 {
			return "", errors.New("wrong repeat rule format")
		}
		days, err := strconv.Atoi(repeatRules[1])
		if err != nil {
			return "", errors.New("wrong repeat rule format")
		}
		if days > 400 {
			return "", errors.New("number of days should be no more than 400")
		}

		date = date.AddDate(0, 0, days)
		for !date.After(now) {
			date = date.AddDate(0, 0, days)
		}

	case "y":
		if len(repeatRules) != 1 {
			return "", errors.New("wrong repeat rule format")
		}

		date = date.AddDate(1, 0, 0)
		for !date.After(now) {
			date = date.AddDate(1, 0, 0)
		}

	case "w":
		if len(repeatRules) != 2 {
			return "", errors.New("wrong repeat rule format")
		}

		if date.Before(now) {
			date = now
		}

		//В пакете time существует целочисленный тип Weekday, однако
		//в нем нумерация начинается с Sunday = 0, в отличие от ТЗ.
		//Поэтому приведем к стандартному виду
		originalDays := strings.Split(repeatRules[1], ",")
		var days = make([]int, len(originalDays))
		for i, v := range originalDays {
			day, err := strconv.Atoi(v)
			if err != nil {
				return "", errors.New("wrong repeat rule format")
			}
			if day < 1 || day > 7 {
				return "", errors.New("days are not in range 1 - 7")
			}
			if day == 7 {
				day = 0
			}
			days[i] = day
		}

		//Сортировка по возрастанию на случай, если дни указаны не по порядку
		slices.Sort(days)

		//Проверка, будут ли дни с указанными номерами на текущей неделе
		var found bool
		for _, weekDay := range days {
			if int(date.Weekday()) < weekDay {
				date = date.AddDate(0, 0, (weekDay - int(date.Weekday())))
				found = true
				break
			}
		}

		//Если день не найден, выбирается ближайший на следующей неделе
		if !found {
			weekDay := days[0]
			date = date.AddDate(0, 0, (7 - int(date.Weekday())))
			date = date.AddDate(0, 0, weekDay)
		}

	case "m":
		if len(repeatRules) != 2 && len(repeatRules) != 3 {
			return "", errors.New("wrong repeat rule format")
		}

		//Правило повторения по числам месяца
		monthDaysStr := strings.Split(repeatRules[1], ",")
		monthDays := make([]int, len(monthDaysStr))
		for i, v := range monthDaysStr {
			number, err := strconv.Atoi(v)
			if err != nil {
				return "", errors.New("wrong repeat rule format")
			}
			if number < -2 || number > 31 {
				return "", errors.New("days are not in range (-2) - 31")
			}
			monthDays[i] = number
		}

		//Функция вычисляет количество дней в данном месяце
		daysIn := func(m time.Month, year int) int {
			return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
		}

		if date.Before(now) {
			date = now
		}
		year, month, day := date.Date()

		//Если указаны только дни месяца
		if len(repeatRules) == 2 {
			var found = false
			for !found {
				//Определение конкретных чисел для данного месяца и сортировка по возрастанию
				for i := range monthDays {
					if monthDays[i] < 0 {
						monthDays[i] = daysIn(month, year) - (-monthDays[i]) + 1
					}
				}
				slices.Sort(monthDays)

				for _, d := range monthDays {
					if day < d && d <= daysIn(month, year) {
						day = d
						found = true
						break
					}
				}
				if !found {
					month += 1
					day = 1
				}
			}
		}

		if len(repeatRules) == 3 {
			monthesStr := strings.Split(repeatRules[2], ",")
			var monthes = make([]int, len(monthesStr))
			for i, v := range monthesStr {
				m, err := strconv.Atoi(v)
				if err != nil {
					return "", errors.New("wrong repeat rule format")
				}
				if m < 1 || m > 12 {
					return "", errors.New("month is not in range 1 - 12")
				}
				monthes[i] = m
			}
			slices.Sort(monthes)

			var found = false
			for !found {
				for _, m := range monthes {
					if m == int(now.Month()) && year == now.Year() {
						for _, d := range monthDays {
							if day < d && d <= daysIn(month, year) {
								day = d
								found = true
								break
							}
						}
						continue
					} else {
						if int(month) < m {
							month = time.Month(m)
							day = monthDays[0]
							found = true
							break
						}
					}
				}
				if !found {
					month = 1
					year += 1
				}
			}
		}
		date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		return date.Format(constants.DateFormat), nil

	default:
		return "", errors.New("wrong repeat rule")
	}

	return date.Format(constants.DateFormat), nil
}
