package onionaddr

import "testing"

const (
	testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDEV7KcR3ESUI7Gr3TovoRhGSUYfwvCVmC44wKTj/kB/t2bPkmm
uWsV31uRrRBpPPLyZuFu+m8vSztiEbTp5ZpxdctHbglv5ys3EfixrqD9EwxYHkmy
kv6kSnYf5p9p3p8KBFeD02VX5tt3pXuuHlNnvuuWY14tZZKy78KJhkEOEQIDAQAB
AoGAZ0ZZxKovZ5rH/uo7bFEKAKjhQklRPh+BML73k/ae29XbatUQmInfMdoSqEWH
5FMS1z4WRfGkmhPQYH0/0+fZm/bHRDgNADRkMd43XoUiK5q7dn/lheDwJ9tbLvy1
5fN2yBbGQ2n9DsM3CN6DpbGd3N/8rTJrcAfPI51NMR5GegECQQD+S+StqUffN9Hq
l+mh6jTyKKwBAR0+I9GAY1UPBFakRUcY8vRsWM8koSkUHBIgl9llB7dTqEM8l4qj
Vi5EFpCpAkEAxahqihC4+4BUcQm1t8B1zvgZs65evp7A4XVlfRZVC3DJS2mYA+On
eNAM7/sdaFkfvqOM9nTXilxySoQh5surKQJALNTwcfVgKGhM59D0bYk+4FpvSJYL
s8LY0ouwmT8ojzlveWSL1vYpPsny1grE314mA3vCxErr36jP1lABRBu+UQJAU+ti
eIYLE/TzZR7bQU38dshNmUUyQrqCZ/cBBO/jYb0cKeGGQjh41Ul4BLfYT4JvgPBN
nCIVlVAU0mBxSF02qQJBALZsK4cZWWEygXFIcMK6TNlfjP1BGrf/bhjVao0j2sIf
x78TKBDam/6FIZCjH367kkwhyTHfwpeMbMkDrSpug4E=
-----END RSA PRIVATE KEY-----`
	testAddress = "bmy7nlgozpyn26tv.onion"
)

func TestLoad(t *testing.T) {
	p, err := LoadPrivateKey([]byte(testPrivateKey))
	if err != nil {
		t.Errorf("Failed to load test key: %v", err)
	}
	err = p.Validate()
	if err != nil {
		t.Errorf("Failed to validate test key: %v", err)
	}
}

func TestGetAddress(t *testing.T) {
	addr, err := GetAddress([]byte(testPrivateKey))
	if err != nil {
		t.Errorf("Failed to get address: %v", err)
	}

	if addr != testAddress {
		t.Errorf("Address did not match: %s != %s", addr, testAddress)
	}
}
