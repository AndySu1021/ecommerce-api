package helper

import "time"

func ParseUTC8Datetime(strArr []string) ([]time.Time, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	timeArr := make([]time.Time, 0)

	for i := 0; i < len(strArr); i++ {
		result, parseErr := time.ParseInLocation("2006-01-02 15:04:05", strArr[i], loc)
		if parseErr != nil {
			return nil, parseErr
		}
		timeArr = append(timeArr, result)
	}

	return timeArr, nil
}

func GetCurrentMonth() time.Time {
	now := time.Now().UTC()
	t := now.AddDate(0, 0, -now.Day()+1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetNextMonth() time.Time {
	now := time.Now().UTC()
	t := now.AddDate(0, 1, -now.Day()+1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetLastMonth() time.Time {
	now := time.Now().UTC()
	t := now.AddDate(0, -1, -now.Day()+1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetLastFourMonths() []time.Time {
	result := make([]time.Time, 4)
	result[0] = GetCurrentMonth()
	for i := 1; i <= 3; i++ {
		result[i] = result[0].AddDate(0, -i, 0)
	}
	return result
}

func GetTodayTimeRange() (startTime, endTime time.Time) {
	date := time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02")
	startTime, _ = time.Parse("2006-01-02", date)

	startTime = startTime.Add(-8 * time.Hour)
	endTime = startTime.Add(24 * time.Hour)

	return
}

func GetYesterdayTimeRange() (startTime, endTime time.Time) {
	date := time.Now().UTC().Add(8 * time.Hour).Format("2006-01-02")
	startTime, _ = time.Parse("2006-01-02", date)

	startTime = startTime.AddDate(0, 0, -1).Add(-8 * time.Hour)
	endTime = startTime.AddDate(0, 0, 1)

	return
}
