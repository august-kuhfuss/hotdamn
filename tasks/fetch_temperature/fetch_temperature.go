package fetchtemperature

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/august-kuhfuss/hotdamn/domain"
	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/august-kuhfuss/hotdamn/tasks"
	"github.com/go-resty/resty/v2"
)

type task struct {
	HttpClient *resty.Client
	Store      store.Store
	URLs       []string
	Interval   time.Duration
}

func NewTask(themometerIPs []string, interval time.Duration) tasks.Task {
	return &task{
		URLs: func() []string {
			urls := make([]string, len(themometerIPs))
			for i, ip := range themometerIPs {
				urls[i] = fmt.Sprintf("http://%s/values.xml", ip)
			}
			return urls
		}(),
		Interval:   interval,
		HttpClient: resty.New(),
	}
}

type RequestXML struct {
	XMLName xml.Name `xml:"Root"`
	Text    string   `xml:",chardata"`
	Val     string   `xml:"val,attr"`
	Agent   struct {
		Text        string `xml:",chardata"`
		Version     string `xml:"Version"`
		XmlVer      string `xml:"XmlVer"`
		DeviceName  string `xml:"DeviceName"`
		Model       string `xml:"Model"`
		VendorID    string `xml:"vendor_id"`
		MAC         string `xml:"MAC"`
		IP          string `xml:"IP"`
		MASK        string `xml:"MASK"`
		SysName     string `xml:"sys_name"`
		SysLocation string `xml:"sys_location"`
		SysContact  string `xml:"sys_contact"`
	} `xml:"Agent"`
	SenSet struct {
		Text  string `xml:",chardata"`
		Entry struct {
			Text     string `xml:",chardata"`
			ID       string `xml:"ID"`
			Name     string `xml:"Name"`
			Units    string `xml:"Units"`
			Value    string `xml:"Value"`
			Min      string `xml:"Min"`
			Max      string `xml:"Max"`
			Hyst     string `xml:"Hyst"`
			EmailSMS string `xml:"EmailSMS"`
			State    string `xml:"State"`
		} `xml:"Entry"`
	} `xml:"SenSet"`
}

func (r *RequestXML) Map(ts time.Time) *domain.MeasurementEntry {
	v, err := strconv.ParseFloat(r.SenSet.Entry.Value, 32)
	if err != nil {
		slog.Error("unable to parse temperature value", slog.String("value", r.SenSet.Entry.Value), slog.String("msg", err.Error()))
		return nil
	}

	return &domain.MeasurementEntry{
		Timestamp: ts,
		Device: domain.MeasurementDevice{
			ID:   r.Agent.MAC,
			Name: r.Agent.DeviceName,
		},
		Value: float32(v),
		Unit:  r.SenSet.Entry.Units,
	}
}

func (t *task) fetchMeasurement(url string) (*domain.MeasurementEntry, error) {
	ts := time.Now()
	xml := &RequestXML{}
	resp, err := t.HttpClient.R().SetResult(xml).Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	return xml.Map(ts), nil
}

func (t *task) Start(ctx context.Context) error {
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, url := range t.URLs {
				go func() {
					m, err := t.fetchMeasurement(url)
					if err != nil {
						slog.Error("unable to fetch temperature", slog.String("url", url), slog.String("msg", err.Error()))
						return
					}

					t.Store.CreateEntry(*m)
				}()
			}
		case <-ctx.Done():
			fmt.Println("fetch temperature task stopped")
			return nil
		}
	}
}
