package gop1

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTelegram(t *testing.T) {
	fixture, err := ioutil.ReadFile("testdata/parser/output")
	require.NoError(t, err)

	lines := strings.Split(string(fixture), "\n")
	tgram := parseTelegram(lines)

	assert.Equal(t, `ISk5\2MT382-1000`, tgram.Device)
	assert.Equal(t, 32, len(tgram.Objects))
}

func TestParseTelegramLine(t *testing.T) {
	tests := []struct {
		line        string
		result      *TelegramObject
		expectError bool
	}{
		{"foo", nil, true},            // bogus
		{"1-3:0.2.8(50", nil, true},   // missing )
		{"1-100:0.2.8(0)", nil, true}, // unknown OBIS ID
		{
			line: "1-3:0.2.8(50)",
			result: &TelegramObject{
				Type: OBISTypeVersionInformation,
				Values: []TelegramValue{
					{Value: "50"},
				},
			},
		},
		{
			line: "0-0:1.0.0(101209113020W)",
			result: &TelegramObject{
				Type: OBISTypeDateTimestamp,
				Values: []TelegramValue{
					{Value: "101209113020W"},
				},
			},
		},
		{
			line: "0-0:96.1.1(4B384547303034303436333935353037)",
			result: &TelegramObject{
				Type: OBISTypeEquipmentIdentifier,
				Values: []TelegramValue{
					{Value: "4B384547303034303436333935353037"},
				},
			},
		},
		{
			line: "1-0:1.8.1(123456.789*kWh)",
			result: &TelegramObject{
				Type: OBISTypeElectricityDeliveredTariff1,
				Values: []TelegramValue{
					{"123456.789", "kWh"},
				},
			},
		},
		{
			line: "1-0:1.8.2(123456.789*kWh)",
			result: &TelegramObject{
				Type: OBISTypeElectricityDeliveredTariff2,
				Values: []TelegramValue{
					{"123456.789", "kWh"},
				},
			},
		},
		{
			line: "1-0:2.8.1(123456.789*kWh)",
			result: &TelegramObject{
				Type: OBISTypeElectricityGeneratedTariff1,
				Values: []TelegramValue{
					{"123456.789", "kWh"},
				},
			},
		},
		{
			line: "1-0:2.8.2(123456.789*kWh)",
			result: &TelegramObject{
				Type: OBISTypeElectricityGeneratedTariff2,
				Values: []TelegramValue{
					{"123456.789", "kWh"},
				},
			},
		},
		{
			line: "0-0:96.14.0(0002)",
			result: &TelegramObject{
				Type: OBISTypeElectricityTariffIndicator,
				Values: []TelegramValue{
					{Value: "0002"},
				},
			},
		},
		{
			line: "1-0:1.7.0(01.193*kW)",
			result: &TelegramObject{
				Type: OBISTypeElectricityDelivered,
				Values: []TelegramValue{
					{"01.193", "kW"},
				},
			},
		},
		{
			line: "1-0:2.7.0(00.000*kW)",
			result: &TelegramObject{
				Type: OBISTypeElectricityGenerated,
				Values: []TelegramValue{
					{"00.000", "kW"},
				},
			},
		},
		{
			line: "0-0:96.7.21(00004)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfPowerFailures,
				Values: []TelegramValue{
					{Value: "00004"},
				},
			},
		},
		{
			line: "0-0:96.7.9(00002)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfLongPowerFailures,
				Values: []TelegramValue{
					{Value: "00002"},
				},
			},
		},
		/*
			{
				line: "1-0:99.97.0(2)(0-0:96.7.19)(101208152415W)(0000000240*s)(101208151004W)(0000000301*s)",
				result: &TelegramObject{
					Type: OBISTypePowerFailureEventLog,
					Values: []TelegramValue{
						{Value: "2"},
						{Value: "0-0:96.7.19"},
						{Value: "101208152415W"},
						{"0000000240", "s"},
						{Value: "101208151004W"},
						{"0000000301", "s"},
					},
				},
			},
		*/

		{
			line: "1-0:32.32.0(00002)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSagsL1,
				Values: []TelegramValue{
					{Value: "00002"},
				},
			},
		},
		{
			line: "1-0:52.32.0(00001)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSagsL2,
				Values: []TelegramValue{
					{Value: "00001"},
				},
			},
		},
		{
			line: "1-0:72.32.0(00000)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSagsL3,
				Values: []TelegramValue{
					{Value: "00000"},
				},
			},
		},
		{
			line: "1-0:32.36.0(00000)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSwellsL1,
				Values: []TelegramValue{
					{Value: "00000"},
				},
			},
		},
		{
			line: "1-0:52.36.0(00003)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSwellsL2,
				Values: []TelegramValue{
					{Value: "00003"},
				},
			},
		},
		{
			line: "1-0:72.36.0(00000)",
			result: &TelegramObject{
				Type: OBISTypeNumberOfVoltageSwellsL3,
				Values: []TelegramValue{
					{Value: "00000"},
				},
			},
		},
		{
			line: "0-0:96.13.0(303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F)",
			result: &TelegramObject{
				Type: OBISTypeTextMessage,
				Values: []TelegramValue{
					{Value: "303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F303132333435363738393A3B3C3D3E3F"},
				},
			},
		},
		{
			line: "1-0:32.7.0(220.1*V)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousVoltageL1,
				Values: []TelegramValue{
					{"220.1", "V"},
				},
			},
		},
		{
			line: "1-0:52.7.0(220.2*V)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousVoltageL2,
				Values: []TelegramValue{
					{"220.2", "V"},
				},
			},
		},
		{
			line: "1-0:72.7.0(220.3*V)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousVoltageL3,
				Values: []TelegramValue{
					{"220.3", "V"},
				},
			},
		},
		{
			line: "1-0:31.7.0(001*A)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousCurrentL1,
				Values: []TelegramValue{
					{"001", "A"},
				},
			},
		},
		{
			line: "1-0:51.7.0(002*A)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousCurrentL2,
				Values: []TelegramValue{
					{"002", "A"},
				},
			},
		},
		{
			line: "1-0:71.7.0(003*A)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousCurrentL3,
				Values: []TelegramValue{
					{"003", "A"},
				},
			},
		},
		{
			line: "1-0:21.7.0(01.111*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerDeliveredL1,
				Values: []TelegramValue{
					{"01.111", "kW"},
				},
			},
		},
		{
			line: "1-0:41.7.0(02.222*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerDeliveredL2,
				Values: []TelegramValue{
					{"02.222", "kW"},
				},
			},
		},
		{
			line: "1-0:61.7.0(03.333*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerDeliveredL3,
				Values: []TelegramValue{
					{"03.333", "kW"},
				},
			},
		},
		{
			line: "1-0:22.7.0(04.444*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerGeneratedL1,
				Values: []TelegramValue{
					{"04.444", "kW"},
				},
			},
		},
		{
			line: "1-0:42.7.0(05.555*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerGeneratedL2,
				Values: []TelegramValue{
					{"05.555", "kW"},
				},
			},
		},
		{
			line: "1-0:62.7.0(06.666*kW)",
			result: &TelegramObject{
				Type: OBISTypeInstantaneousPowerGeneratedL3,
				Values: []TelegramValue{
					{"06.666", "kW"},
				},
			},
		},
	}

	for _, test := range tests {
		obj, err := parseTelegramLine(test.line)
		if test.expectError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		assert.Equal(t, test.result, obj)
	}
}
