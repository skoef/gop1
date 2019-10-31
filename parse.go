package gop1

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errCOSEMNoMatch     = errors.New("COSEM was no match")
	telegramHeaderRegex = regexp.MustCompile(`^\/(.+)$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)(?:\(([^\)]+)\))+$`)
	cosemUnitRegex      = regexp.MustCompile(`^([\d\.]+)\*(?i)([a-z]+)$`)
)

var (
	allOBISTypes = map[string]OBISType{
		"1-3:0.2.8":   OBISTypeVersionInformation,
		"0-0:1.0.0":   OBISTypeDateTimestamp,
		"0-0:96.1.1":  OBISTypeEquipmentIdentifier,
		"1-0:1.8.1":   OBISTypeElectricityDeliveredTariff1,
		"1-0:1.8.2":   OBISTypeElectricityDeliveredTariff2,
		"1-0:2.8.1":   OBISTypeElectricityGeneratedTariff1,
		"1-0:2.8.2":   OBISTypeElectricityGeneratedTariff2,
		"0-0:96.14.0": OBISTypeElectricityTariffIndicator,
		"1-0:1.7.0":   OBISTypeElectricityDelivered,
		"1-0:2.7.0":   OBISTypeElectricityGenerated,
		"0-0:96.7.21": OBISTypeNumberOfPowerFailures,
		"0-0:96.7.9":  OBISTypeNumberOfLongPowerFailures,
		"1-0:99.97.0": OBISTypePowerFailureEventLog,
		"1-0:32.32.0": OBISTypeNumberOfVoltageSagsL1,
		"1-0:52.32.0": OBISTypeNumberOfVoltageSagsL2,
		"1-0:72.32.0": OBISTypeNumberOfVoltageSagsL3,
		"1-0:32.36.0": OBISTypeNumberOfVoltageSwellsL1,
		"1-0:52.36.0": OBISTypeNumberOfVoltageSwellsL2,
		"1-0:72.36.0": OBISTypeNumberOfVoltageSwellsL3,
		"0-0:96.13.0": OBISTypeTextMessage,
		"1-0:32.7.0":  OBISTypeInstantaneousVoltageL1,
		"1-0:52.7.0":  OBISTypeInstantaneousVoltageL2,
		"1-0:72.7.0":  OBISTypeInstantaneousVoltageL3,
		"1-0:31.7.0":  OBISTypeInstantaneousCurrentL1,
		"1-0:51.7.0":  OBISTypeInstantaneousCurrentL2,
		"1-0:71.7.0":  OBISTypeInstantaneousCurrentL3,
		"1-0:21.7.0":  OBISTypeInstantaneousPowerDeliveredL1,
		"1-0:41.7.0":  OBISTypeInstantaneousPowerDeliveredL2,
		"1-0:61.7.0":  OBISTypeInstantaneousPowerDeliveredL3,
		"1-0:22.7.0":  OBISTypeInstantaneousPowerGeneratedL1,
		"1-0:42.7.0":  OBISTypeInstantaneousPowerGeneratedL2,
		"1-0:62.7.0":  OBISTypeInstantaneousPowerGeneratedL3,

		// In the specification, there are several OBIS types specified for slave
		// devices as gas meters and such. These have variable OBIS IDs (first 2 digits)
		// and this requires special treatment in the current way of parsing
	}
)

// parsedTelegram parses lines from P1 data, or telegrams
func parseTelegram(lines []string) *Telegram {
	tgram := &Telegram{}

	for _, l := range lines {
		// try to detect identification header
		match := telegramHeaderRegex.FindStringSubmatch(l)
		if len(match) > 0 {
			tgram.Device = match[1]
			continue
		}

		obj, err := parseTelegramLine(strings.TrimSpace(l))
		if err != nil {
			continue
		}

		tgram.Objects = append(tgram.Objects, obj)
	}

	return tgram
}

func parseTelegramLine(line string) (*TelegramObject, error) {
	matches := cosemOBISRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, errCOSEMNoMatch
	}

	// is this a known COSEM object
	_, ok := allOBISTypes[matches[1]]
	if !ok {
		return nil, errCOSEMNoMatch
	}

	obj := &TelegramObject{
		Type: allOBISTypes[matches[1]],
	}

	for _, v := range matches[2:] {
		ov := TelegramValue{}
		// check if the unit of the value is specified as well
		match := cosemUnitRegex.FindStringSubmatch(v)
		if len(match) > 1 {
			ov.Value = match[1]
			ov.Unit = match[2]
		} else {
			ov.Value = v
		}

		obj.Values = append(obj.Values, ov)
	}

	return obj, nil
}
