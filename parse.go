package gop1

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errCOSEMNoMatch     = errors.New("COSEM was no match")
	telegramHeaderRegex = regexp.MustCompile(`^\/(.+)$`)
	cosemOBISRegex      = regexp.MustCompile(`^(\d+-\d+:\d+\.\d+\.\d+)([0-9A-Za-z\(\)\*\-\.\:]+)$`)
	cosemValsRegex      = regexp.MustCompile(`\(([^\)]+)\)`)
	cosemUnitRegex      = regexp.MustCompile(`^([\d\.]+)\*(?i)([a-z0-9]+)$`)
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

		"0-0:96.1.4":  OBISTypeVersionInformation,
		"0-0:96.13.1": OBISTypeConsumerMessageCode,
		"0-0:96.3.10": OBISTypeBreakerState,
		"0-0:17.0.0":  OBISTypeLimiterThreshold,
		"1-0:31.4.0":  OBISTypeFuseThresholdL1,
	}

	// In the specification, there are several OBIS types specified for slave
	// devices as gas meters and such. These have variable OBIS IDs (first 2 digits)
	// and this requires special treatment in the current way of parsing
	addOBISTypes = map[string]OBISType{
		`0-(\d+):96.1.0`: OBISTypeGasEquipmentIdentifier,
		`0-(\d+):24.1.0`: OBISTypeDeviceType,
		`0-(\d+):24.2.1`: OBISTypeGasDelivered,

		`0-(\d+):96.1.1`: OBISTypeGasEquipmentIdentifier,
		`0-(\d+):24.4.0`: OBISTypeGasValveState,
		`0-(\d+):24.2.3`: OBISTypeGasDelivered,
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

	var obj *TelegramObject
	// is this a known COSEM object
	if _, ok := allOBISTypes[matches[1]]; ok {
		obj = &TelegramObject{
			Type: allOBISTypes[matches[1]],
		}
	} else {
		// try to match it to one of the additional types
		for ptr, obisType := range addOBISTypes {
			if regexp.MustCompile(ptr).MatchString(matches[1]) {
				obj = &TelegramObject{Type: obisType}
				break
			}
		}
	}

	if obj == nil {
		return nil, errCOSEMNoMatch
	}

	vmatches := cosemValsRegex.FindAllStringSubmatch(matches[2], -1)
	if len(vmatches) == 0 {
		return nil, errCOSEMNoMatch
	}

	for _, v := range vmatches {
		ov := TelegramValue{}
		// check if the unit of the value is specified as well
		match := cosemUnitRegex.FindStringSubmatch(v[1])
		if len(match) > 1 {
			ov.Value = match[1]
			ov.Unit = match[2]
		} else {
			ov.Value = v[1]
		}

		obj.Values = append(obj.Values, ov)
	}

	return obj, nil
}
