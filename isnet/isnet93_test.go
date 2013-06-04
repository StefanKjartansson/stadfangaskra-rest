package isnet

import "testing"

func TestIsnet93ToWgs84(t *testing.T) {

	const in1, in2 = 357337.497727273, 408120.711363636
	const out1, out2 = 64.14614813044987, -21.931836213122278

	x1, x2 := Isnet93ToWgs84(in1, in2)

	if x1 != out1 {
		t.Errorf("Isnet93ToWgs84(%v,%v) was (%v, %v), expected (%v,%v)",
			in1, in2, x1, x2, out1, out2)
	}

}
