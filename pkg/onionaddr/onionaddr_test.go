package onionaddr

import "testing"

var (
	privateKeys = map[string]string{
		"bmy7nlgozpyn26tv.onion": `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`,
		"fm2dtevxaby3fovv.onion": `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDe4gUXo20v/+WCWSvcJYl+g6bXhxpiICnfhSLrATNw35ohGFij
LL8AQeSVKXvQjo+vRFg8ytxwWyuCGKUvlsErUnx/4UrvXaGJzZwZ/AJDiocsueSr
xi4NRMt11re2hnWFCKR64+PRl6x/tWnwTB2mLvQWkefPAhPeOQr9VsXJlQIDAQAB
AoGAA1mQTAenx3XoJCpk710dEOq3ojukmN704iglGzUcadDihybPhjxQ7pcO8XL7
KmpKHI9BaECASawFHmJycSGp646VclHg6Sp5SdbS2iXbHuKUakErENObeErCmcDR
uCZcqRTjaWbTxeUgtesNeJXPSs7J5Uc4KuI9Bz0JRXYe10ECQQD3samgeMQAkA3D
xQMZcizqpU7SDkemdbdXRVmj0Tpwv1MtxIMlg6nrvf+W9l1Zla3Yc4JrYxFf5Tdu
ZUY5nJpRAkEA5ltdk2ijmLoTyj7L2agQ/BF4p0ylcYLSnorqMHA9KxKac5EkpNwy
335450B32FAw+YkjUxrvbbSinc/zDHLmBQJATe+j7O8y1O5+tkOuNvp68ZX0GBoQ
J2tQtfAHRYlW9xTsSjFUBqDH/Xo6CrkYJFD5c2rc9Ycld2P7LgxFrWj+EQJBAKi5
9WGafHHKoeI6es/TVZV8KpbIAkqRWzo7X+vY2kTpXG4XuvJyQ2UGWwJvaxjrK+Qq
+c/OY99ENvRGs6zDoA0CQGSyguVtCJQMdloOTrpskVaKUhY5rAIuIKoVFKJuHvT/
/ZbgCz7zaR6lCnZiY088EIkBEkeZiQ4w9Of5uPKgJbc=
-----END RSA PRIVATE KEY-----`,
		"ahf462kszuq5cffl.onion": `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC/LyZvdzgxd+RQhQT1snaaTOVVx/RUZZTrVoXFJjLtPTT1LOxZ
KsuAgrk5HkQJjEBlDAlG+ueMrYEaIPmgG1emYM0wTVPw9pBlfElWX0BgG1q0VtlM
GX7tCJzey4NeWBL8GCc47AOpNf646NdkIu+EJ1RxFJryf1baBQP7J0d4SQIDAQAB
AoGBALnx5OMSxC+w2PnLdnh1O891LLSSxrtzFTUTMQX/0hZVqnUvXSyYZ9c0zWuV
WT0kENl2rGtBywVTFzbPjZpAHaz+gdTMJ1y3n/OLd7mY9beY0xKbkrKVKrv+SaBj
B1ru9y7R0TvmG3UbOQzGmso4RGRKhz3A4iTtquXNXXcLV591AkEA91S3N+lzbWHE
/Oa/D1eQYyB8r6+HK97gBxnWqwyQ0Xlm25CTV3Kssp2Ekh5sCk61jnf3IqvWY8uL
7IfewK2oowJBAMXioiVZCk9C1QAZ1yPD+NsfN+lUZ7k8ZRYVqZXpk1EilhSAhLS4
E2IOZLEbJascSdszFvjWn5rAbBFhBQj1jiMCQGinhRt4geoPy73DmabRQ3xeW8Qv
PsAWf68hhM898u1gNGDFzULceCzgMB9wFgFKitJs+rrGAWKa12tPlrbrBIcCQQCH
WLigpOMZTVPUitgMnWijrxmV3OZI2xck4NIqOCVLtEVEZpbd6J1RTxjtzeyYuXOG
ms4Wiu2FciE4Tcyc0R9TAkBYXNeM4ASl2KOKQdY2z6Z0AKBZdtTdsDx9sWGLtqS+
xxfnYUK4xAIOxqEiRYw1uoUeyPvjkVn9Ychds8r5gwR0
-----END RSA PRIVATE KEY-----`,
	}
)

func TestLoad(t *testing.T) {
	for _, key := range privateKeys {
		p, err := LoadPrivateKey([]byte(key))
		if err != nil {
			t.Errorf("Failed to load test key: %v", err)
		}
		err = p.Validate()
		if err != nil {
			t.Errorf("Failed to validate test key: %v", err)
		}
	}
}

func TestGetAddress(t *testing.T) {
	for testAddr, key := range privateKeys {
		addr, err := GetAddress([]byte(key))
		if err != nil {
			t.Errorf("Failed to get address: %v", err)
		}

		if addr != testAddr {
			t.Errorf("Address did not match: %s != %s", addr, testAddr)
		}
	}
}
