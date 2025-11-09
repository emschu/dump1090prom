// dump1090prom
// Copyright (C) 2025 emschu[aet]mailbox.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program.
// If not, see <https://www.gnu.org/licenses/>.

package main

// Dump1090AircraftJson aircraft.json struct
type Dump1090AircraftJson struct {
	Now      float64               `json:"now"`
	Messages uint64                `json:"messages"`
	Aircraft []AircraftStateRecord `json:"aircraft"`
}

// AircraftStateRecord aircraft.json struct
type AircraftStateRecord struct {
	Hex            string    `json:"hex"`
	AltBaro        *float64  `json:"alt_baro,omitempty"`
	AltGeom        *float64  `json:"alt_geom,omitempty"`
	BaroRate       *float64  `json:"baro_rate,omitempty"`
	Category       string    `json:"category,omitempty"`
	Emergency      string    `json:"emergency,omitempty"`
	Flight         string    `json:"flight,omitempty"`
	GS             *float64  `json:"gs,omitempty"`
	GVA            *int      `json:"gva,omitempty"`
	GeomRate       *float64  `json:"geom_rate,omitempty"`
	IAS            *float64  `json:"ias,omitempty"`
	Lat            *float64  `json:"lat,omitempty"`
	Lon            *float64  `json:"lon,omitempty"`
	MLAT           []string  `json:"mlat"`
	Mach           *float64  `json:"mach,omitempty"`
	MagHeading     *float64  `json:"mag_heading,omitempty"`
	Messages       int       `json:"messages"`
	ModeA          bool      `json:"modea,omitempty"`
	ModeC          bool      `json:"modec,omitempty"`
	NACP           *int      `json:"nac_p,omitempty"`
	NACV           *int      `json:"nac_v,omitempty"`
	NIC            *int      `json:"nic,omitempty"`
	NICBaro        *int      `json:"nic_baro,omitempty"`
	NavAltitudeMCP *float64  `json:"nav_altitude_mcp,omitempty"`
	NavHeading     *float64  `json:"nav_heading,omitempty"`
	NavModes       *[]string `json:"nav_modes,omitempty"`
	NavQNH         *float64  `json:"nav_qnh,omitempty"`
	OAT            *int      `json:"oat,omitempty"`
	RC             *int      `json:"rc,omitempty"`
	RSSI           float64   `json:"rssi"`
	Roll           *float64  `json:"roll,omitempty"`
	SDA            *int      `json:"sda,omitempty"`
	SIL            *int      `json:"sil,omitempty"`
	SILType        string    `json:"sil_type,omitempty"`
	Seen           float64   `json:"seen"`
	SeenPos        *float64  `json:"seen_pos,omitempty"`
	SPI            *int      `json:"spi,omitempty"`
	Squawk         string    `json:"squawk,omitempty"`
	TAS            *float64  `json:"tas,omitempty"`
	TAT            *int      `json:"tat,omitempty"`
	TISB           []string  `json:"tisb,omitempty"`
	Track          *float64  `json:"track,omitempty"`
	TrackRate      *float64  `json:"track_rate,omitempty"`
	TrueHeading    *float64  `json:"true_heading,omitempty"`
	Version        int       `json:"version"`
}

// Dump1090ReceiverJson receiver.json struct
type Dump1090ReceiverJson struct {
	Version string   `json:"version"`
	Refresh int      `json:"refresh"`
	History int      `json:"history"`
	Lat     *float64 `json:"lat"`
	Lon     *float64 `json:"lon"`
}

// Dump1090StatsJson stats.json struct
type Dump1090StatsJson struct {
	Latest struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Local struct {
			SamplesProcessed int     `json:"samples_processed"`
			SamplesDropped   int     `json:"samples_dropped"`
			Modeac           int     `json:"modeac"`
			Modes            int     `json:"modes"`
			Bad              int     `json:"bad"`
			UnknownIcao      int     `json:"unknown_icao"`
			Accepted         []int   `json:"accepted"`
			StrongSignals    int     `json:"strong_signals"`
			GainDb           float64 `json:"gain_db"`
		} `json:"local"`
		Remote struct {
			Modeac      int   `json:"modeac"`
			Modes       int   `json:"modes"`
			Bad         int   `json:"bad"`
			UnknownIcao int   `json:"unknown_icao"`
			Accepted    []int `json:"accepted"`
		} `json:"remote"`
		Cpr struct {
			Surface               int `json:"surface"`
			Airborne              int `json:"airborne"`
			GlobalOk              int `json:"global_ok"`
			GlobalBad             int `json:"global_bad"`
			GlobalRange           int `json:"global_range"`
			GlobalSpeed           int `json:"global_speed"`
			GlobalSkipped         int `json:"global_skipped"`
			LocalOk               int `json:"local_ok"`
			LocalAircraftRelative int `json:"local_aircraft_relative"`
			LocalReceiverRelative int `json:"local_receiver_relative"`
			LocalSkipped          int `json:"local_skipped"`
			LocalRange            int `json:"local_range"`
			LocalSpeed            int `json:"local_speed"`
			Filtered              int `json:"filtered"`
		} `json:"cpr"`
		AltitudeSuppressed int `json:"altitude_suppressed"`
		Cpu                struct {
			Demod      int `json:"demod"`
			Reader     int `json:"reader"`
			Background int `json:"background"`
		} `json:"cpu"`
		Tracks struct {
			All           int `json:"all"`
			SingleMessage int `json:"single_message"`
			Unreliable    int `json:"unreliable"`
		} `json:"tracks"`
		Messages     int   `json:"messages"`
		MessagesByDf []int `json:"messages_by_df"`
	} `json:"latest"`
	Last1Min struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Local struct {
			SamplesProcessed int     `json:"samples_processed"`
			SamplesDropped   int     `json:"samples_dropped"`
			Modeac           int     `json:"modeac"`
			Modes            int     `json:"modes"`
			Bad              int     `json:"bad"`
			UnknownIcao      int     `json:"unknown_icao"`
			Accepted         []int   `json:"accepted"`
			Signal           float64 `json:"signal"`
			Noise            float64 `json:"noise"`
			PeakSignal       float64 `json:"peak_signal"`
			StrongSignals    int     `json:"strong_signals"`
			GainDb           float64 `json:"gain_db"`
		} `json:"local"`
		Remote struct {
			Modeac      int   `json:"modeac"`
			Modes       int   `json:"modes"`
			Bad         int   `json:"bad"`
			UnknownIcao int   `json:"unknown_icao"`
			Accepted    []int `json:"accepted"`
		} `json:"remote"`
		Cpr struct {
			Surface               int `json:"surface"`
			Airborne              int `json:"airborne"`
			GlobalOk              int `json:"global_ok"`
			GlobalBad             int `json:"global_bad"`
			GlobalRange           int `json:"global_range"`
			GlobalSpeed           int `json:"global_speed"`
			GlobalSkipped         int `json:"global_skipped"`
			LocalOk               int `json:"local_ok"`
			LocalAircraftRelative int `json:"local_aircraft_relative"`
			LocalReceiverRelative int `json:"local_receiver_relative"`
			LocalSkipped          int `json:"local_skipped"`
			LocalRange            int `json:"local_range"`
			LocalSpeed            int `json:"local_speed"`
			Filtered              int `json:"filtered"`
		} `json:"cpr"`
		AltitudeSuppressed int `json:"altitude_suppressed"`
		Cpu                struct {
			Demod      int `json:"demod"`
			Reader     int `json:"reader"`
			Background int `json:"background"`
		} `json:"cpu"`
		Tracks struct {
			All           int `json:"all"`
			SingleMessage int `json:"single_message"`
			Unreliable    int `json:"unreliable"`
		} `json:"tracks"`
		Messages     int   `json:"messages"`
		MessagesByDf []int `json:"messages_by_df"`
		Adaptive     struct {
			GainDb              float64     `json:"gain_db"`
			DynamicRangeLimitDb float64     `json:"dynamic_range_limit_db"`
			GainChanges         int         `json:"gain_changes"`
			LoudUndecoded       int         `json:"loud_undecoded"`
			LoudDecoded         int         `json:"loud_decoded"`
			NoiseDbfs           float64     `json:"noise_dbfs"`
			GainSeconds         [][]float64 `json:"gain_seconds"`
		} `json:"adaptive"`
	} `json:"last1min"`
	Last5Min struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Local struct {
			SamplesProcessed int     `json:"samples_processed"`
			SamplesDropped   int     `json:"samples_dropped"`
			Modeac           int     `json:"modeac"`
			Modes            int     `json:"modes"`
			Bad              int     `json:"bad"`
			UnknownIcao      int     `json:"unknown_icao"`
			Accepted         []int   `json:"accepted"`
			Signal           float64 `json:"signal"`
			Noise            float64 `json:"noise"`
			PeakSignal       float64 `json:"peak_signal"`
			StrongSignals    int     `json:"strong_signals"`
			GainDb           float64 `json:"gain_db"`
		} `json:"local"`
		Remote struct {
			Modeac      int   `json:"modeac"`
			Modes       int   `json:"modes"`
			Bad         int   `json:"bad"`
			UnknownIcao int   `json:"unknown_icao"`
			Accepted    []int `json:"accepted"`
		} `json:"remote"`
		Cpr struct {
			Surface               int `json:"surface"`
			Airborne              int `json:"airborne"`
			GlobalOk              int `json:"global_ok"`
			GlobalBad             int `json:"global_bad"`
			GlobalRange           int `json:"global_range"`
			GlobalSpeed           int `json:"global_speed"`
			GlobalSkipped         int `json:"global_skipped"`
			LocalOk               int `json:"local_ok"`
			LocalAircraftRelative int `json:"local_aircraft_relative"`
			LocalReceiverRelative int `json:"local_receiver_relative"`
			LocalSkipped          int `json:"local_skipped"`
			LocalRange            int `json:"local_range"`
			LocalSpeed            int `json:"local_speed"`
			Filtered              int `json:"filtered"`
		} `json:"cpr"`
		AltitudeSuppressed int `json:"altitude_suppressed"`
		Cpu                struct {
			Demod      int `json:"demod"`
			Reader     int `json:"reader"`
			Background int `json:"background"`
		} `json:"cpu"`
		Tracks struct {
			All           int `json:"all"`
			SingleMessage int `json:"single_message"`
			Unreliable    int `json:"unreliable"`
		} `json:"tracks"`
		Messages     int   `json:"messages"`
		MessagesByDf []int `json:"messages_by_df"`
		Adaptive     struct {
			GainDb              float64     `json:"gain_db"`
			DynamicRangeLimitDb float64     `json:"dynamic_range_limit_db"`
			GainChanges         int         `json:"gain_changes"`
			LoudUndecoded       int         `json:"loud_undecoded"`
			LoudDecoded         int         `json:"loud_decoded"`
			NoiseDbfs           float64     `json:"noise_dbfs"`
			GainSeconds         [][]float64 `json:"gain_seconds"`
		} `json:"adaptive"`
	} `json:"last5min"`
	Last15Min struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Local struct {
			SamplesProcessed int64   `json:"samples_processed"`
			SamplesDropped   int     `json:"samples_dropped"`
			Modeac           int     `json:"modeac"`
			Modes            int     `json:"modes"`
			Bad              int     `json:"bad"`
			UnknownIcao      int     `json:"unknown_icao"`
			Accepted         []int   `json:"accepted"`
			Signal           float64 `json:"signal"`
			Noise            float64 `json:"noise"`
			PeakSignal       float64 `json:"peak_signal"`
			StrongSignals    int     `json:"strong_signals"`
			GainDb           float64 `json:"gain_db"`
		} `json:"local"`
		Remote struct {
			Modeac      int   `json:"modeac"`
			Modes       int   `json:"modes"`
			Bad         int   `json:"bad"`
			UnknownIcao int   `json:"unknown_icao"`
			Accepted    []int `json:"accepted"`
		} `json:"remote"`
		Cpr struct {
			Surface               int `json:"surface"`
			Airborne              int `json:"airborne"`
			GlobalOk              int `json:"global_ok"`
			GlobalBad             int `json:"global_bad"`
			GlobalRange           int `json:"global_range"`
			GlobalSpeed           int `json:"global_speed"`
			GlobalSkipped         int `json:"global_skipped"`
			LocalOk               int `json:"local_ok"`
			LocalAircraftRelative int `json:"local_aircraft_relative"`
			LocalReceiverRelative int `json:"local_receiver_relative"`
			LocalSkipped          int `json:"local_skipped"`
			LocalRange            int `json:"local_range"`
			LocalSpeed            int `json:"local_speed"`
			Filtered              int `json:"filtered"`
		} `json:"cpr"`
		AltitudeSuppressed int `json:"altitude_suppressed"`
		Cpu                struct {
			Demod      int `json:"demod"`
			Reader     int `json:"reader"`
			Background int `json:"background"`
		} `json:"cpu"`
		Tracks struct {
			All           int `json:"all"`
			SingleMessage int `json:"single_message"`
			Unreliable    int `json:"unreliable"`
		} `json:"tracks"`
		Messages     int   `json:"messages"`
		MessagesByDf []int `json:"messages_by_df"`
		Adaptive     struct {
			GainDb              float64     `json:"gain_db"`
			DynamicRangeLimitDb float64     `json:"dynamic_range_limit_db"`
			GainChanges         int         `json:"gain_changes"`
			LoudUndecoded       int         `json:"loud_undecoded"`
			LoudDecoded         int         `json:"loud_decoded"`
			NoiseDbfs           float64     `json:"noise_dbfs"`
			GainSeconds         [][]float64 `json:"gain_seconds"`
		} `json:"adaptive"`
	} `json:"last15min"`
	Total struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Local struct {
			SamplesProcessed int64   `json:"samples_processed"`
			SamplesDropped   int     `json:"samples_dropped"`
			Modeac           int     `json:"modeac"`
			Modes            int     `json:"modes"`
			Bad              int     `json:"bad"`
			UnknownIcao      int     `json:"unknown_icao"`
			Accepted         []int   `json:"accepted"`
			Signal           float64 `json:"signal"`
			Noise            float64 `json:"noise"`
			PeakSignal       float64 `json:"peak_signal"`
			StrongSignals    int     `json:"strong_signals"`
			GainDb           float64 `json:"gain_db"`
		} `json:"local"`
		Remote struct {
			Modeac      int   `json:"modeac"`
			Modes       int   `json:"modes"`
			Bad         int   `json:"bad"`
			UnknownIcao int   `json:"unknown_icao"`
			Accepted    []int `json:"accepted"`
		} `json:"remote"`
		Cpr struct {
			Surface               int `json:"surface"`
			Airborne              int `json:"airborne"`
			GlobalOk              int `json:"global_ok"`
			GlobalBad             int `json:"global_bad"`
			GlobalRange           int `json:"global_range"`
			GlobalSpeed           int `json:"global_speed"`
			GlobalSkipped         int `json:"global_skipped"`
			LocalOk               int `json:"local_ok"`
			LocalAircraftRelative int `json:"local_aircraft_relative"`
			LocalReceiverRelative int `json:"local_receiver_relative"`
			LocalSkipped          int `json:"local_skipped"`
			LocalRange            int `json:"local_range"`
			LocalSpeed            int `json:"local_speed"`
			Filtered              int `json:"filtered"`
		} `json:"cpr"`
		AltitudeSuppressed int `json:"altitude_suppressed"`
		Cpu                struct {
			Demod      int `json:"demod"`
			Reader     int `json:"reader"`
			Background int `json:"background"`
		} `json:"cpu"`
		Tracks struct {
			All           int `json:"all"`
			SingleMessage int `json:"single_message"`
			Unreliable    int `json:"unreliable"`
		} `json:"tracks"`
		Messages     int   `json:"messages"`
		MessagesByDf []int `json:"messages_by_df"`
		Adaptive     struct {
			GainDb              float64     `json:"gain_db"`
			DynamicRangeLimitDb float64     `json:"dynamic_range_limit_db"`
			GainChanges         int         `json:"gain_changes"`
			LoudUndecoded       int         `json:"loud_undecoded"`
			LoudDecoded         int         `json:"loud_decoded"`
			NoiseDbfs           float64     `json:"noise_dbfs"`
			GainSeconds         [][]float64 `json:"gain_seconds"`
		} `json:"adaptive"`
	} `json:"total"`
}
