package faker

import (
	"github.com/google/uuid"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type UserData struct {
	ClickId       string    `db:"click_id"`
	BrandId       string    `db:"brand_id"`
	UserId        int64     `db:"user_id"`
	BalanceBefore float32   `db:"balance_before"`
	BalanceAfter  float32   `db:"balance_after"`
	Amount        float32   `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}

func GenerateFakeData(monthsCount int, rowsCount int, brandCount int, userCount int, out []chan<- UserData) {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	begin := time.Date(now.Year(), now.Month()-time.Month(monthsCount), now.Day(), 0, 0, 0, 0, now.Location())

	rowsPerDay := int(math.Round(float64(rowsCount) / (end.Sub(begin).Hours() / 24)))

	users := generateUsers(userCount)
	brands := generateBrands(brandCount)

	usersBalance := make(map[int64]float32, 0)
	for _, u := range users {
		usersBalance[u] = 0.0
	}

	for begin.Unix() <= end.Unix() {
		rowsForDay := generateRandomDays(begin, rowsPerDay)
		for _, date := range rowsForDay {
			clickId, _ := uuid.NewUUID()
			userId := users[rand.Int31n(int32(len(users)-1))]
			bBefore := usersBalance[userId]
			amount := rand.Float32() * 100.0
			bAfter := bBefore + amount
			usersBalance[userId] = bAfter
			ud := UserData{
				ClickId:       clickId.String(),
				BrandId:       brands[rand.Int31n(int32(len(brands)-1))],
				UserId:        userId,
				BalanceBefore: bBefore,
				BalanceAfter:  bAfter,
				Amount:        amount,
				CreatedAt:     date,
			}
			for _, ch := range out {
				ch <- ud
			}
		}
		log.Println("Data for date: " + begin.Format("2006-01-02") + " has been generated")
		begin = begin.Add(24 * time.Hour)
	}
	for _, ch := range out {
		close(ch)
	}
}

func generateUsers(userCount int) []int64 {
	users := make([]int64, 0)
	for uc := 1; uc <= userCount; uc++ {
		users = append(users, int64(uc))
	}

	return users
}

func generateBrands(brandsCount int) []string {
	//brands := make([]string, 0)
	//for bc := 0; bc < brandsCount; bc++ {
	//	uid, _ := uuid.NewUUID()
	//	brands = append(brands, uid.String())
	//}
	//
	//return brands
	brands := make([]string, 0)
	for uc := 1; uc <= brandsCount; uc++ {
		brands = append(brands, strconv.Itoa(uc))
	}

	return brands
}

func generateRandomDays(day time.Time, daysCount int) []time.Time {
	d := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	days := make([]time.Time, 0)
	for c := 0; c < daysCount; c++ {
		days = append(days, generateRandomDay(d))
	}
	sort.SliceStable(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})

	return days
}

func generateRandomDay(day time.Time) time.Time {
	return day.Add(time.Second * time.Duration(rand.Int31n(84000)))
}
