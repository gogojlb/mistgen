package app

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-mistgen/utils"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	pgx "github.com/jackc/pgx/v4"
)

type JobResult struct {
	Executed bool
	Time     time.Time
}

type MeasureResponse struct {
	Success bool        `json:"success"`
	Data    MeasureData `json:"data"`
}

type MeasureData struct {
	Temperature int `json:"temperature"`
	Humidity    int `json:"humidity"`
}

type SwitchResponce struct {
	Success bool   `json:"success"`
	IsOn    bool   `json:"is_on"`
	Error   string `json:"error"`
}

func StartJob(config utils.MistConfig, db *pgx.Conn, lastTime *JobResult, grafanaConfig utils.GrafanaConfig) *JobResult {
	var result = &JobResult{Executed: false}
	m, err := measure(config.MiApiUrl)
	if err != nil {
		log.Error(err)
		return result
	}
	log.Infof("got measurements: %+v", m)
	insertMeasurements(m, db)

	if float64(m.Temperature/100) < float64(config.ActivationTemperature) {
		log.Infof("temperature is lower than threshold (%v), skip", config.ActivationTemperature)
		return result
	}
	if !lastTime.Executed {
		log.Info("first-time execution, for minimal duration")
		err = startMisting(config.MinActiveDuration, config.MiApiUrl)
		if err != nil {
			log.Error(err)
			return result
		}
		result = &JobResult{Executed: true, Time: time.Now()}
		return result
	}
	timeNow := time.Now()
	timeSinceLastActive := timeNow.Sub(lastTime.Time)
	if timeSinceLastActive < config.MinActivationInterval {
		log.Infof("temperature is high enough (>%v), but MinActivationInterval (%v) is not passed yet since last activation (%v)", config.ActivationTemperature, config.MinActivationInterval, lastTime.Time.Format("2006-01-02 15:04:05"))
		return result
	}
	durationPercent := CalculateDurationPercent(config, m.Temperature)
	log.Infof("calculated duration percent: %v\\%", int64(durationPercent))
	duration := time.Duration(int32(durationPercent)*int32(timeSinceLastActive/time.Second)/100) * time.Second
	if duration < config.MinActiveDuration {
		log.Infof("calculated duration (%v) is lower than threshold (%v)", duration, config.MinActiveDuration)
		return result
	}
	log.Infof("start mist machine for %v", duration)
	startTime := time.Now()
	err = startMisting(duration, config.MiApiUrl)
	if err != nil {
		log.Error(err)
		return result
	}
	annotationText := fmt.Sprintf("Start time: %v\nDuration: %v\nCalculated duration percent: %v\n Temperature %v\nHumidity: %v\n", startTime, duration, durationPercent, m.Temperature, m.Humidity)
	grafanaAnnotate(grafanaConfig, annotationText, startTime, (startTime.Add(duration)))
	result = &JobResult{Executed: true, Time: time.Now()}
	return result
}

func measure(MiApiUrl string) (MeasureData, error) {
	var data MeasureData
	var mResp MeasureResponse

	client := http.Client{}
	req := NewSensorReq(MiApiUrl)
	res, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return data, err
		}
		return data, fmt.Errorf("http error %s %s", res.Status, string(body))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &mResp)
	if err != nil {
		return data, err
	}
	if !mResp.Success {
		return data, fmt.Errorf("responce status failed %+v", mResp)
	}
	return mResp.Data, nil
}

func insertMeasurements(m MeasureData, db *pgx.Conn) {
	_, err := db.Exec(context.Background(), "insert into timeseries.measurements(time,temperature,humidity ) values($1,$2,$3)", time.Now(), m.Temperature, m.Humidity)
	if err != nil {
		log.Error(err)
	}
}

func startMisting(duration time.Duration, MiApiUrl string) error {
	errOn := utils.Retry(2, 1*time.Second, func() (err error) {
		err = executeSwitch(MiApiUrl, true)
		return
	})
	if errOn != nil {
		return fmt.Errorf("cannot switch on: %v", errOn)
	}
	log.Infof("Mist machine has been turned on. Waiting %v...", duration)
	time.Sleep(duration)
	errOff := utils.Retry(2, 1*time.Second, func() (err error) {
		err = executeSwitch(MiApiUrl, false)
		return
	})
	if errOn != nil {
		return fmt.Errorf("cannot switch off: %v", errOff)
	}
	log.Info("Mist machine has been turned off")
	return nil
}

func executeSwitch(MiApiUrl string, isOn bool) error {
	client := http.Client{}
	req := NewSwitchReq(MiApiUrl, isOn)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("http error %s %s", res.Status, string(body))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var sResp SwitchResponce
	err = json.Unmarshal(body, &sResp)
	if err != nil {
		return err
	}
	if !sResp.Success {
		return fmt.Errorf("responce status failed %+v", sResp)
	}
	if sResp.IsOn != isOn {
		return fmt.Errorf("the switch status has not changed %+v", sResp)
	}
	return nil
}

func CalculateDurationPercent(config utils.MistConfig, temperature int) float64 {
	// aaa := math.Atan(float64(10000))
	// println(aaa)
	DurationPercent := (float64(2.0) / float64(math.Pi)) * math.Atan(float64(((float64(temperature)/100)-float64(config.ActivationTemperature))/5)) * float64(config.MaxDurationPercent)
	// println(DurationPercent)
	return DurationPercent
}
