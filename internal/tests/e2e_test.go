package tests

import (
	"bytes"
	"encoding/json"
	"time"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func getToken(t *testing.T, role string) string {
	body, _ := json.Marshal(map[string]string{"role": role})
	res, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer res.Body.Close() //nolint:errcheck

	var data map[string]string
	json.NewDecoder(res.Body).Decode(&data) //nolint:errcheck
	require.NotEmpty(t, data["token"])
	return data["token"]
}

func authReq(t *testing.T, method, path, token string, body any) *http.Response {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, baseURL+path, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return res
}

// E2E: создание переговорки → расписания → брони
func TestE2E_CreateRoomScheduleBooking(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	// создаём комнату
	res := authReq(t, "POST", "/rooms/create", adminToken, map[string]any{
		"name":        "E2E Переговорка",
		"description": "Тестовая",
		"capacity":    5,
	})
	require.Equal(t, 201, res.StatusCode)

	var roomResp map[string]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&roomResp))
	res.Body.Close() //nolint:errcheck

	roomID := roomResp["room"]["id"].(string)

	// создаём расписание
	res = authReq(t, "POST", "/rooms/"+roomID+"/schedule/create", adminToken, map[string]any{
		"days_of_week": []int{1, 2, 3, 4, 5, 6, 7},
		"start_time": "09:00",
		"end_time":   "18:00",
	})
	require.Equal(t, 201, res.StatusCode)
	res.Body.Close() //nolint:errcheck
 
	// берем дату внутри окна генерации слотов
	date := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	res = authReq(
		t,
		"GET",
		"/rooms/"+roomID+"/slots/list?date="+date,
		userToken,
		nil,
	)

	require.Equal(t, 200, res.StatusCode)

	var slotsResp map[string][]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&slotsResp))
	res.Body.Close() //nolint:errcheck

	slots := slotsResp["slots"]
	require.NotEmpty(t, slots)

	slotID := slots[0]["id"].(string)

	// создаём бронь
	res = authReq(t, "POST", "/bookings/create", userToken, map[string]any{
		"slot_id": slotID,
	})

	require.Equal(t, 201, res.StatusCode)

	var bookingResp map[string]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&bookingResp))
	res.Body.Close() //nolint:errcheck

	assert.Equal(t, "active", bookingResp["booking"]["status"])
	assert.Equal(t, slotID, bookingResp["booking"]["slot_id"])
}

// E2E: отмена брони
func TestE2E_CancelBooking(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	// создаём комнату
	res := authReq(t, "POST", "/rooms/create", adminToken, map[string]any{
		"name": "E2E Cancel Room",
	})

	require.Equal(t, 201, res.StatusCode)

	var roomResp map[string]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&roomResp))
	res.Body.Close() //nolint:errcheck

	t.Logf("roomResp: %+v", roomResp)
	roomID := roomResp["room"]["id"].(string)
	t.Logf("roomID: %s", roomID)

	// создаём расписание
	res = authReq(t, "POST", "/rooms/"+roomID+"/schedule/create", adminToken, map[string]any{
		"days_of_week": []int{1, 2, 3, 4, 5, 6, 7},
		"start_time": "09:00",
		"end_time":   "18:00",
	})

	require.Equal(t, 201, res.StatusCode)
	res.Body.Close() //nolint:errcheck

	date := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// получаем слот
	res = authReq(
		t,
		"GET",
		"/rooms/"+roomID+"/slots/list?date="+date,
		userToken,
		nil,
	)

	require.Equal(t, 200, res.StatusCode)

	var slotsResp map[string][]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&slotsResp))
	res.Body.Close() //nolint:errcheck

	require.NotEmpty(t, slotsResp["slots"])

	slotID := slotsResp["slots"][0]["id"].(string)


	// бронируем
	res = authReq(t, "POST", "/bookings/create", userToken, map[string]any{
		"slot_id": slotID,
	})

	require.Equal(t, 201, res.StatusCode)

	var bookingResp map[string]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&bookingResp))
	res.Body.Close() //nolint:errcheck

	bookingID := bookingResp["booking"]["id"].(string)
	// отменяем
	res = authReq(
		t,
		"POST",
		"/bookings/"+bookingID+"/cancel",
		userToken,
		nil,
	)

	require.Equal(t, 200, res.StatusCode)

	var cancelResp map[string]map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&cancelResp))
	res.Body.Close() //nolint:errcheck

	assert.Equal(t, "cancelled", cancelResp["booking"]["status"])

	// идемпотентность
	res = authReq(
		t,
		"POST",
		"/bookings/"+bookingID+"/cancel",
		userToken,
		nil,
	)

	assert.Equal(t, 200, res.StatusCode)
	res.Body.Close() //nolint:errcheck
}