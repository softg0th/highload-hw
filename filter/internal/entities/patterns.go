package entities

var StopWords = map[string]float64{
	"earn":           0.90,
	"free":           0.95,
	"click":          0.92,
	"win":            0.91,
	"discount":       0.88,
	"money":          0.89,
	"limited":        0.84,
	"offer":          0.83,
	"urgent":         0.87,
	"now":            0.75,
	"guaranteed":     0.90,
	"exclusive":      0.85,
	"risk-free":      0.93,
	"trial":          0.79,
	"subscribe":      0.76,
	"bonus":          0.82,
	"gift":           0.80,
	"buy":            0.70,
	"save":           0.65,
	"promo":          0.86,
	"reward":         0.81,
	"investment":     0.72,
	"opportunity":    0.77,
	"credit":         0.83,
	"deal":           0.71,
	"access":         0.68,
	"miracle":        0.90,
	"income":         0.88,
	"secret":         0.86,
	"pills":          0.93,
	"act":            0.63,
	"order":          0.67,
	"extra":          0.69,
	"winner":         0.89,
	"double":         0.74,
	"cheap":          0.78,
	"unlimited":      0.87,
	"billion":        0.91,
	"cash":           0.90,
	"limited-time":   0.85,
	"instant":        0.84,
	"apply":          0.65,
	"guarantee":      0.83,
	"hidden":         0.82,
	"easy":           0.70,
	"increase":       0.73,
	"pre-approved":   0.89,
	"refund":         0.80,
	"no-cost":        0.88,
	"luxury":         0.86,
	"risk":           0.66,
	"exclusive-deal": 0.91,
	"free-trial":     0.87,
	"work-from-home": 0.94,
	"no-obligation":  0.88,
}

var UrlRegexp = "http[s]?://\\S"
