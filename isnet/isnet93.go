package isnet

import "math"

const (
	a    = 6378137.0
	f    = 1 / 298.257222101
	lat1 = 64.25
	lat2 = 65.75
	latc = 65.00
	lonc = 19.00
	eps  = 0.00000000001
	rho  = 57.29577951308232   // 45 / math.Atan2(1.0, 1.0)
	e    = 0.08181919104281579 // math.Sqrt(f * (2 - f))
)

func fx(p float64) float64 {
	return a * math.Cos(p/rho) / math.Sqrt(1-math.Pow(e*math.Sin(p/rho), 2))
}

func f1(p float64) float64 {
	return math.Log((1 - p) / (1 + p))
}

func f2(p float64) float64 {
	return f1(p) - e*f1(e*p)
}

func f3(p float64, pol1 float64, f2sin1 float64, sint float64) float64 {
	return pol1 * math.Exp((f2(math.Sin(p/rho))-f2sin1)*sint/2)
}

func Isnet93ToWgs84(x float64, y float64) (lat float64, lon float64) {
	dum := f2(math.Sin(lat1/rho)) - f2(math.Sin(lat2/rho))
	sint := 2 * (math.Log(fx(lat1)) - math.Log(fx(lat2))) / dum
	f2sin1 := f2(math.Sin(lat1 / rho))
	pol1 := fx(lat1) / sint
	polc := f3(latc, pol1, f2sin1, sint) + 500000.0
	peq := a * math.Cos(latc/rho) / (sint * math.Exp(sint*math.Log((45-latc/2)/rho)))
	pol := math.Sqrt(math.Pow(x-500000, 2) + math.Pow(polc-y, 2))
	lat = 90 - 2*rho*math.Atan(math.Exp(math.Log(pol/peq)/sint))
	lon = 0
	fact := rho * math.Cos(lat/rho) / sint / pol
	fact = rho * math.Cos(lat/rho) / sint / pol
	delta := 1.0

	for math.Abs(delta) > eps {
		delta = (f3(lat, pol1, f2sin1, sint) - pol) * fact
		lat += delta
	}

	lon = -(lonc + rho*math.Atan((500000-x)/(polc-y))/sint)

	return
}
