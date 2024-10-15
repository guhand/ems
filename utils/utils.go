package utils

import (
	"ems/infrastructure/config"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(config.Config.TokenDuration).Unix(),
	})
	return token.SignedString([]byte(config.Config.JwtSecretKey))
}

/**
 * @function: GenerateOTP
 * @description: function used to generate random numbers of six digits
 * @param: None
 * @returns: 6 digit string
 */
func GenerateOTP() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer between 100000 and 999999 (inclusive)
	otp := rand.Intn(900000) + 100000

	//converting integer into string
	return strconv.Itoa(otp)
}

/**
 * @function: SendFogotPasswordMail
 * @description: function used to send mail
 * @param: to, otp string, requestedAt time.Time, offset int
 * @returns: error if mail not sent
 */
func SendFogotPasswordMail(to, otp string, requestedAt time.Time) error {

	displayName := config.Config.SmtpDisplayName
	from := config.Config.SmtpUserName
	password := config.Config.SmtpPassword
	smtpHost := config.Config.SmtpHost
	smtpPort := config.Config.SmtpPort

	subject := "EMS OTP"

	body := fmt.Sprintf(`<p>Your EMS OTP is <b>%s</b> requested at: <b>%s</b></p>`, otp, time.Now().Format("15:04:05 2006-01-02"))

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", from, to, subject, body)

	auth := smtp.PlainAuth(displayName, from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

type PaginationResponse struct {
	From       uint        `json:"from"`
	To         uint        `json:"to"`
	TotalCount uint        `json:"totalCount"`
	TotalPages uint        `json:"totalPages"`
	Data       interface{} `json:"data"`
}

func PaginatedResponse(totalCount uint, page uint, data interface{}) *PaginationResponse {

	var (
		limit      uint = 10
		offset     uint = 0
		from       uint = 0
		to         uint = 0
		totalPages uint = 0
	)

	if page > 1 {
		offset = (page - 1) * limit
	}

	if totalCount > 0 && page > 0 {
		totalPages = uint(math.Ceil(float64(totalCount) / float64(limit)))
		from = offset + 1
		if (offset + limit) > totalCount {
			to = totalCount
		} else {
			to = offset + limit
		}
	}

	return &PaginationResponse{
		From:       from,
		To:         to,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Data:       data,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SqlParamValidator(input string) string {
	unsafeChars := []string{"'", "\"", ";", "--"}
	sanitizedInput := input
	for _, char := range unsafeChars {
		if char == "'" {
			sanitizedInput = strings.ReplaceAll(sanitizedInput, char, "''")
		} else {
			sanitizedInput = strings.ReplaceAll(sanitizedInput, char, "")
		}
	}
	return sanitizedInput
}

func ValidateTimeDifference(fromTime, toTime string) error {
	const timeLayout = "15:04"

	from, err := time.Parse(timeLayout, fromTime)
	if err != nil {
		return fmt.Errorf("invalid fromTime format: %v", err)
	}

	to, err := time.Parse(timeLayout, toTime)
	if err != nil {
		return fmt.Errorf("invalid toTime format: %v", err)
	}

	diff := to.Sub(from)

	// Check if the difference is exactly 1 hour
	if diff != time.Hour {
		return errors.New("the difference between fromTime and toTime must be exactly 1 hour")
	}

	return nil
}

func GetDateRangeForMonthAndYear(year int, month int) (startDate, endDate string) {
	if year > 0 && month > 0 {
		startYear := year
		startMonth := month - 1
		if month == 1 {
			startMonth = 12
			startYear = year - 1
		}

		startDateObj := time.Date(startYear, time.Month(startMonth), 27, 0, 0, 0, 0, time.Local)
		endDateObj := time.Date(year, time.Month(month), 26, 0, 0, 0, 0, time.Local)

		startDate = startDateObj.Format("2006-01-02")
		endDate = endDateObj.Format("2006-01-02")
	} else if year > 0 {
		startDateObj := time.Date(year-1, 12, 27, 0, 0, 0, 0, time.Local)
		endDateObj := time.Date(year, 12, 26, 0, 0, 0, 0, time.Local)

		startDate = startDateObj.Format("2006-01-02")
		endDate = endDateObj.Format("2006-01-02")
	} else {
		now := time.Now()
		startDateObj := time.Date(now.Year(), now.Month()-1, 27, 0, 0, 0, 0, time.Local)
		endDateObj := time.Date(now.Year(), now.Month(), 26, 0, 0, 0, 0, time.Local)

		startDate = startDateObj.Format("2006-01-02")
		endDate = endDateObj.Format("2006-01-02")
	}
	return
}

func IsValidDate(dateStr string) (*time.Time, bool) {
	const layout = "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, false
	}
	return &date, true
}
