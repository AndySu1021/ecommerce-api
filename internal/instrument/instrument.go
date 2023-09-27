package instrument

import "ecommerce-api/internal/config"

var instance *Instrument

type Instrument struct {
	Request *RequestInstrument
}

func InitInstrument(cfg config.AppConfig) {
	instance = &Instrument{
		Request: NewRequestInstrument(cfg),
	}
}
