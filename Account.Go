package velocitylimit

import (
	"log"
	"time"
)

type Account struct {
	CustID      string
	Balance     float64
	DailyLimit  *DailyLimit
	WeeklyLimit *WeeklyLimit
}

type DailyLimit struct {
	Date        time.Time
	DailyBudget float64
	DailyTxn    int
}

type WeeklyLimit struct {
	Date         time.Time
	WeeklyBudget float64
}

func NewDailyLimit(d time.Time) *DailyLimit {
	return &DailyLimit{
		Date:        getStartOfDay(d),
		DailyBudget: 5000,
		DailyTxn:    3,
	}
}

func NewWeeklyLimit(d time.Time) *WeeklyLimit {
	return &WeeklyLimit{
		Date:         getStartOfWeek(d),
		WeeklyBudget: 20000,
	}
}

func NewAccount(custID string) *Account {
	return &Account{
		CustID: custID,
	}
}

func (dl *DailyLimit) Validate(amount float64, t time.Time) bool {
	transDay := getStartOfDay(t)
	if transDay.After(dl.Date) {
		log.Println("Transaction date: ", transDay, " is newer than current daily date: ", dl.Date)
		newdl := NewDailyLimit(t)
		*dl = *newdl
	}
	if dl.DailyBudget-amount < 0 {
		return false
	}
	if dl.DailyTxn-1 < 0 {
		return false
	}
	return true
}

func (dl *DailyLimit) Apply(amount float64) {
	dl.DailyBudget -= amount
	dl.DailyTxn--
}

func (wl *WeeklyLimit) Validate(amount float64, t time.Time) bool {
	transWeek := getStartOfWeek(t)
	if transWeek.After(wl.Date) {
		log.Println("Transaction date: ", transWeek, " is newer than current week start: ", wl.Date)
		newwl := NewWeeklyLimit(t)
		*wl = *newwl
	}
	if wl.WeeklyBudget-amount < 0 {
		return false
	}
	return true
}

func (wl *WeeklyLimit) Apply(amount float64) {
	wl.WeeklyBudget -= amount
}

func (a *Account) LoadFunds(r *Request) bool {
	if a.DailyLimit == nil {
		log.Println("Daily limit is missing. Adding now")
		a.DailyLimit = NewDailyLimit(r.ParsedTime)
	}
	if a.DailyLimit.Validate(r.ParsedAmount, r.ParsedTime) == false {
		log.Println("Daily limit breached. Transaction rejected: ", r.ID)
		return false
	}
	if a.WeeklyLimit == nil {
		log.Println("Weekly limit is missing. Adding now")
		a.WeeklyLimit = NewWeeklyLimit(r.ParsedTime)
	}
	if a.WeeklyLimit.Validate(r.ParsedAmount, r.ParsedTime) == false {
		log.Println("Weekly limit breached. Transaction rejected: ", r.ID)
		return false
	}
	a.Balance += r.ParsedAmount
	a.DailyLimit.Apply(r.ParsedAmount)
	a.WeeklyLimit.Apply(r.ParsedAmount)
	log.Println("Transaction approved: ", r.ID)
	return true
}

func getStartOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}

func getStartOfWeek(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+int(time.Monday-d.Weekday()), 0, 0, 0, 0, time.UTC)
}
