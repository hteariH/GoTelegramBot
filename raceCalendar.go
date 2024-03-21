package main

type RaceCalendar struct {
	Meta struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Fields      struct {
			RaceID       string `json:"race_id"`
			Name         string `json:"name"`
			Country      string `json:"country"`
			Status       string `json:"status"`
			Season       string `json:"season"`
			StartDate    string `json:"start_date"`
			EndDate      string `json:"end_date"`
			Track        string `json:"track"`
			SessionArray struct {
				ID          string `json:"id"`
				SessionName string `json:"session_name"`
				Date        string `json:"date"`
			} `json:"session_array"`
		} `json:"fields"`
	} `json:"meta"`

	Results []struct {
		RaceID    int    `json:"race_id"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Status    string `json:"status"`
		Season    string `json:"season"`
		Track     string `json:"track"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Sessions  []struct {
			ID          int    `json:"id"`
			SessionName string `json:"session_name"`
			Date        string `json:"date"`
		} `json:"sessions"`
	} `json:"results"`
}
