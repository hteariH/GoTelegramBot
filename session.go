package main

import "time"

type Session struct {
	SessionKey       int      `json:"session_key"`
	SessionName      string   `json:"session_name"`
	DateStart        jsonTime `json:"date_start"`
	DateEnd          jsonTime `json:"date_end"`
	GmtOffset        string   `json:"gmt_offset"`
	SessionType      string   `json:"session_type"`
	MeetingKey       int      `json:"meeting_key"`
	Location         string   `json:"location"`
	CountryKey       int      `json:"country_key"`
	CountryCode      string   `json:"country_code"`
	CountryName      string   `json:"country_name"`
	CircuitKey       int      `json:"circuit_key"`
	CircuitShortName string   `json:"circuit_short_name"`
	Year             int      `json:"year"`
}

type jsonTime time.Time

const dateLayout = "2006-01-02T15:04:05"

func (sD *jsonTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(dateLayout, s[1:len(s)-1])
	if err != nil {
		return err
	}
	*sD = jsonTime(t)
	return nil
}
