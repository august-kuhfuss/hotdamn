package tasks

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

	"github.com/go-resty/resty/v2"
)

type FetchTemperatureTask struct {
	HttpClient *resty.Client
	Store      store.Store
	URLs       []string
	Interval   time.Duration
}

func NewFetchTemperatureTask(sensorIPs []string, interval time.Duration, s store.Store) Task {
	return &FetchTemperatureTask{
		URLs: func() []string {
			urls := make([]string, len(sensorIPs))
			for i, ip := range sensorIPs {
				urls[i] = fmt.Sprintf("http://%s/values.xml", ip)
			}
			return urls
		}(),
		Interval:   interval,
		HttpClient: resty.New(),
		Store:      s,
	}
}

type SensorXMLResponse struct {
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
	SenSet []struct {
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

func (t *FetchTemperatureTask) Start(ctx context.Context) error {
	ticker := time.NewTicker(t.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, url := range t.URLs {
				go func(url string) {
					ts := time.Now()
					xml := &SensorXMLResponse{}
					resp, err := t.HttpClient.R().SetResult(xml).Get(url)
					if err != nil {
						slog.Error("unable to fetch temperature", slog.String("url", url), slog.String("msg", err.Error()))
						return
					}
					if resp.StatusCode() != http.StatusOK {
						slog.Error("unexpected status code", slog.String("url", url), slog.Int("status", resp.StatusCode()))
						return
					}

					for _, sen := range xml.SenSet {
						cmp := &store.CreateMeasurementParams{
							SensorID:  sen.Entry.ID,
							Timestamp: ts,
							ValueK: func(unit string) float32 {
								f, err := strconv.ParseFloat(sen.Entry.Value, 32)
								if err != nil {
									slog.Error("unable to parse temperature value", slog.String("value", sen.Entry.Value), slog.String("msg", err.Error()))
									return 0
								}

								switch unit {
								case "C":
									return domain.ConvertCelsiusToKelvin(float32(f))
								case "F":
									return domain.ConvertFahrenheitToKelvin(float32(f))
								default:
									slog.Warn("unknown unit", slog.String("unit", unit))
									return 0
								}
							}(sen.Entry.Units),
						}

						if err := t.Store.CreateMeasurement(cmp); err != nil {
							slog.Error("unable to create measurement", slog.Any("params", cmp), slog.String("msg", err.Error()))
						}
						slog.Info("measurement fetched", slog.Any("params", cmp))

						csp := &store.CreateSensorParams{
							ID:   sen.Entry.ID,
							Name: sen.Entry.Name,
						}

						if err := t.Store.CreateOrUpdateSensor(csp); err != nil {
							slog.Error("unable to create or update sensor", slog.Any("params", csp), slog.String("msg", err.Error()))
						}
						slog.Info("sensor fetched", slog.Any("params", csp))

					}

				}(url)
			}
		case <-ctx.Done():
			fmt.Println("fetch temperature task stopped")
			return nil
		}
	}
}
