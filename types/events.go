package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type EventRegistrationEmailDTO struct {
	Event *queries.Event
	User  *queries.User
}

func (d *EventRegistrationEmailDTO) EmailRecipient() string {
	return fmt.Sprintf("%s %s <%s>", d.User.GivenName.String(), d.User.FamilyName.String(), d.User.Email.String())
}

func (d *EventRegistrationEmailDTO) lang() string {
	return string(d.User.PreferredLocale)
}

func (d *EventRegistrationEmailDTO) Title() string {
	if d.lang() == "pl" {
		return d.Event.TitlePl
	}
	return d.Event.TitleEn
}

func (d *EventRegistrationEmailDTO) FormattedStartsAt() string {
	t := d.Event.StartsAt.UTC()
	if d.lang() == "pl" {
		return t.Format("2.01.2006, 15:04") + " UTC"
	}
	return t.Format("January 2, 2006 at 15:04") + " UTC"
}

func (d *EventRegistrationEmailDTO) IsVirtual() bool {
	return d.Event.IsVirtual
}

func (d *EventRegistrationEmailDTO) VenueNameLocalized() string {
	if d.lang() == "pl" && d.Event.VenueNamePl != nil {
		return *d.Event.VenueNamePl
	}
	if d.Event.VenueNameEn != nil {
		return *d.Event.VenueNameEn
	}
	return ""
}

func (d *EventRegistrationEmailDTO) VenueAddress() string {
	if d.Event.IsVirtual {
		return ""
	}
	var parts []string
	if name := d.VenueNameLocalized(); name != "" {
		parts = append(parts, name)
	}
	if d.Event.VenueStreet != nil {
		parts = append(parts, *d.Event.VenueStreet)
	}
	var cityLine string
	if d.lang() == "pl" && d.Event.VenueCityPl != nil {
		cityLine = *d.Event.VenueCityPl
	} else if d.Event.VenueCityEn != nil {
		cityLine = *d.Event.VenueCityEn
	}
	if d.Event.VenuePostalCode != nil && cityLine != "" {
		cityLine = *d.Event.VenuePostalCode + " " + cityLine
	}
	if cityLine != "" {
		parts = append(parts, cityLine)
	}
	return strings.Join(parts, ", ")
}

func (d *EventRegistrationEmailDTO) EventURL() string {
	return config.PublicUrl + "/events/" + d.Event.Slug
}

// ICS generates an iCalendar (RFC 5545) payload for the event.
func (d *EventRegistrationEmailDTO) ICS() []byte {
	fmtTime := func(t time.Time) string {
		return t.UTC().Format("20060102T150405Z")
	}
	uid := fmt.Sprintf("%s@homeosapiens.eu", d.Event.ID.String())
	title := escapeICSText(d.Title())
	location := escapeICSText(d.VenueAddress())

	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//Homeo Sapiens//Go//EN\r\n")
	sb.WriteString("METHOD:PUBLISH\r\n")
	sb.WriteString("BEGIN:VEVENT\r\n")
	fmt.Fprintf(&sb, "UID:%s\r\n", uid)
	fmt.Fprintf(&sb, "DTSTAMP:%s\r\n", fmtTime(time.Now()))
	fmt.Fprintf(&sb, "DTSTART:%s\r\n", fmtTime(d.Event.StartsAt))
	fmt.Fprintf(&sb, "DTEND:%s\r\n", fmtTime(d.Event.EndsAt))
	fmt.Fprintf(&sb, "SUMMARY:%s\r\n", title)
	if location != "" {
		fmt.Fprintf(&sb, "LOCATION:%s\r\n", location)
	}
	sb.WriteString("END:VEVENT\r\n")
	sb.WriteString("END:VCALENDAR\r\n")

	return []byte(sb.String())
}

func escapeICSText(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, ";", `\;`)
	s = strings.ReplaceAll(s, ",", `\,`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
