package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	ics "github.com/arran4/golang-ical"
)

func main() {
	http.HandleFunc("/calendar.ics", handleCalendar)

	fmt.Println("Servidor escuchando en :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCalendar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	calendarURL := os.Getenv("CALENDAR_URL")

	if calendarURL == "" {
		http.Error(w, "CALENDAR_URL not configured", http.StatusInternalServerError)
		return
	}

	cal, err := ics.ParseCalendarFromUrl(calendarURL)
	if err != nil {
		http.Error(w, "Error fetching source calendar", http.StatusInternalServerError)
		log.Println("Error fetching calendar:", err)
		return
	}

	filtered := ics.NewCalendar()

	for _, event := range cal.Events() {
		summary := event.GetProperty(ics.ComponentPropertySummary)
		if summary == nil {
			continue
		}

		if !strings.Contains(summary.Value, "Inteligencia Artificial") &&
			!strings.Contains(summary.Value, "Programación Distribuida") {
			continue
		}

		description := event.GetProperty(ics.ComponentPropertyDescription)

		if description != nil && strings.Contains(description.Value, "GR2") {
			continue
		}

		uid := event.GetProperty(ics.ComponentPropertyUniqueId)
		if uid == nil {
			continue
		}

		newEvent := filtered.AddEvent(uid.Value)

		newEvent.SetSummary(summary.Value)

		if startDate, err := event.GetStartAt(); err == nil {
			newEvent.SetStartAt(startDate, ics.WithTZID("Europe/Madrid"))
		}

		if endDate, err := event.GetEndAt(); err == nil {
			newEvent.SetEndAt(endDate, ics.WithTZID("Europe/Madrid"))
		}

		if description != nil {
			newEvent.SetProperty(
				ics.ComponentPropertyDescription,
				description.Value,
			)
		}
	}

	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("Content-Disposition", `inline; filename="calendar.ics"`)

	if _, err := w.Write([]byte(filtered.Serialize())); err != nil {
		log.Println("Error writing response:", err)
	}
}
