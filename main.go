package main

import (
	"context"
	"time"

	"github.com/offer365/example/endecrypt/endeaes"
	"github.com/offer365/example/endecrypt/endeaesrsa"
	"github.com/offer365/example/endecrypt/endersa"
	"github.com/offer365/example/winsysinfo"
	"github.com/offer365/odin/config"
	"github.com/offer365/odin/odinX"
	"github.com/offer365/odin/utils"
)

const (
	clusterToken = "G81dau36S8HAI5J4jGs45T61B93d7t08YQ1CjN1Sv0shRrPE46C4sw0WwVCunxSi"
	embedAuthPwd = "u1MMj4Cx4B1Jr2qpU3a9unmc2U97Pz3VLgsQ0hzFR90d16W5j94q3cQD5ShSu5Jt"

	storeLicenseKey = "/ac7RCREONg3OX11K6qYw1Z6CVMdpPb58yJlIz5enBt70mZ6U1nhGHlYUvwRF0faD/73J0xNwDW9aPqYC6cFRwOHT9456q67ltO1c895HqJ6h094Da79p2Fm6c22kZ8o22"
	storeClearLicenseKey = "/xms5xAQzM7q5e290E2Q93LZVP4aQ45I8JX6T48fwhd3D8kFMzC26n16n2Ec09cSm/W12vaU3C4PiyNhEh8uEpaW97ae778QSb72uE79X1Ddhph6sL0nWPnNRi8e3K34b8"
	storeClientConfigKeyPrefix = "/n2qjSrmVKxfLL8yv7h7zDuwfU9z5hy644z6NWFn2qMKWxQt9lmRLvfexxlU1q93e/pwgEA0swhnMw1qXH4VqbnP06SVrsuJl5VW4OXm5t1o7A2Bv7T0XITJ5vW2ogglU4/"
	storeClientKeyPrefix = "/qN5jQ8WL01q175A0nBQxaIJWAI9572o4Q47T668j5lx902vuq4kZ35SF5fA7YWk9/YTv8MV7am8ARUd50dp7700wKH12re5Jzws9u56AGgL46CMvUB0ViOe1f2r1hDRI4/"
	storeTokenKey = "/O50806rQ9tIBWMBW12RQh1250Gw69FLflmL0rn55XBFE0QE7fe8g344b31g1e15T/SqT6RNm49O886ZRdM61j1889YxmPhRD92sq456Qj509Co088YwtFaF0dEXMk4Z1M"
	storeSerialNumKey = "/P900nXZQNq4Ns93Z5S3OV4oTrz3YdP30JbQW16W875ouK4lIc1eR71RlTOTz82h6/xzOx9578qmV0NUu72bU5b195NGS61fw2ddI77MigxC26aOB20DO5rhFkVMu6Jo8a"
	storeHashSalt = "ri5HEbbc6WN4Z88L75F59h5007TqV1w8JOIK7i6d38V2yfnSnt5a1nD4nAXyZ1q2"
	server_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/IzANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZzZXJ2ZXIxEjAQBgNV
BAMMCXNlcnZlci5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMRF
1mgbKYNO2X0iqX89Rrzc+xftqegQ+7V0n9Sa1HE07xQVcgu05faGB4B/29HPQ/gh
JMt1IxkXlISNuQwIDM5XVSahkH1OhQmtQnTKjYXFgboRFHMQUk26lKoIZ3o9AJ8s
QTPCLBw7a9StBpeWhBzEumDymP60hmGhTft4tbY85MrmObfTZ8KbQiHIy22jXNGV
N5ok61q4tlMV8HYK89q4WX7IcQusdK9NNwL1jZNQ4+WICEe2/zs8xY9r4REONKoM
HOME5aS+EvQSVwh5LyvNuPxa8io83EjokT3yRqZllvmXD/hVS/BCM927fgsiDfm0
ezuE5+AGiMR1N0agv1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAHxXVd/v7noVZ
LJ8IsLty3BjMX7ZjSvkyyrchxdQQIfCoMc/UGkDZ5TgvdPkE8eAfdSVwwrcpGf8V
C4ccB9flekd6HrO7Uo9mWrKcjyQn2MjwAZNDhcs5Sxrz8TusJEQk4iQYSq0oc4Nr
qGrR/2kXEirwXi/xQ0saVXalfhkK5W+rO/YWTc8K3leARQ6BDjGbF2BHRtj6HEZL
RnhJEbx+BvplXMlWQ5CBBYt/NQa/MKJDd2stT70Si8E1lIGIGaVQAy43uT7xy8XW
jSrruOAv1SVLovhSxjsMiu/jXwZsVAtaFAuT4ajiWQHzbNqUjVnt7dJIWJPCnL6h
lhaV0MUy7Q==
-----END CERTIFICATE-----

`
	server_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxEXWaBspg07ZfSKpfz1GvNz7F+2p6BD7tXSf1JrUcTTvFBVy
C7Tl9oYHgH/b0c9D+CEky3UjGReUhI25DAgMzldVJqGQfU6FCa1CdMqNhcWBuhEU
cxBSTbqUqghnej0AnyxBM8IsHDtr1K0Gl5aEHMS6YPKY/rSGYaFN+3i1tjzkyuY5
t9NnwptCIcjLbaNc0ZU3miTrWri2UxXwdgrz2rhZfshxC6x0r003AvWNk1Dj5YgI
R7b/OzzFj2vhEQ40qgwc4wTlpL4S9BJXCHkvK824/FryKjzcSOiRPfJGpmWW+ZcP
+FVL8EIz3bt+CyIN+bR7O4Tn4AaIxHU3RqC/VwIDAQABAoIBAAQbHeghoVWw4ZXf
ksIpqwAqc0pF24cSS+G45dsRvh38KIA4DqG2EBV/KksC4bta5aYcM2PaOHi+6Il5
WYSp6nKqmwpq2NX2PYw9RqWg0yMYRaV50/6wObiMja2c7WU+P3QU/ewyRK/2gkP5
tqiXKn5bkzaR/KdfaWxDbpkzJkIArLAELqEBuS0noxikrfypPanGnXk7IDhYo+rZ
cE0UHOhpkeo7gXeVc9tU/cjTRwBK7awKLIDWyknHGrL28nxMqKf06SzxG2oz6Hn3
twOtwAUS7tjophOZ6WCStgCOVFf0Ue6yJmja9xgWy/r2jJsH5/VV0xJZvmWGxr8T
IQh4oskCgYEA92Katy0Cvl1kS1/cf0ExMtOzXIwtCDu35axGl1FR3VMcoboPmH2h
HrRxSpcIgkRXz7wxsj3zttBXu8assjmwtCWzbDIE0YGYQ3v1CwDITihAyhevhW4b
UxN181RhMo1qHIcgULsVR5+P857FAHRSWWewh77ZK7x17fdQJshZujMCgYEAyxuT
R1CthfC7rbIX359tD9jb1XtG+XCgygZYv+6uoknmWMMmUqgDmQ3u8p4kuHudB6gm
/kZkxrluwJM5B8UKC1NRkejHP2ZO8ygpEGQp7t1H3BBFSfUVlu+YmfD5SjHhK9U5
2t+hfyuO8m0r+XdYk6lliEYufVlPMzJffT3rSk0CgYEAhs+jRGMw9ZBrUXAB9w8N
wob/XVW+TJhOlMiXB2r3U8cw+SktyonbvaHTgzRfHK4ltDz4UAvWvi83QEr6XX12
wBUze6ieW5Vl5pCsbryUa5MgC4Fw0yO3nEQkqN+4wBW0V6uDfrsU050ukzJYZPD+
113cI31rV5wyH+YANcJEs2UCgYEAmh0SY8qT4E4KGoJIGyadWqjyJcqk0CDl4GVw
cjJp0DrCzhdFvPI/yKMJ7I6Szmj9fhHZhJdlYGTT5MvROlQIiw9tlYlLpo+62EZg
4k8egmDlZdXyvWt6Nk0XPbfbcLDoapogjDOkFxq2HL054NDuJR0kLYMTQ4nAztgq
HJ4fKwECgYEAinsJM6lw9m3eyRRuPRFE4jNwg5KmZRjVuZ06+nPW/Sb7GXdN+5e6
62y87e63MRTm1r2C4g/esnqAOcS6iRHQtdTFrG8DU/j9F+uaB5TWZTroxqQ6h2F0
OjGZcdCMohluWRztbas01OZKSoDx1pEfP+H4kKFJ8LhWQXLU0lWibEw=
-----END RSA PRIVATE KEY-----

`
	client_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/JDANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZjbGllbnQxEjAQBgNV
BAMMCWNsaWVudC5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAObU
dbPa2MMnp7X0P6TeUM9+gJgRgVdrOm05EPnf4p1xEmFq09bGupZpD+pVoU/yH/oi
wA4gwYtgk5ETtTfTbF8YUma6LYDye2m98zXiyVWpTs9pmxVRUTcpnjmyIS7mXSNE
hShN26OCTk8DtlL9STFnFWQY2Sb9PVjwDWTrXkHalQU3PFEmoQ/QPbTbBN2gydDn
WkK6LxgTaSA9xMw/j5upZh58aoLVwd8IevzKn/YnwQBEC0ptVQGl6B5EUKabhTWh
q6c4gDAhcqhdRFZa4UMcOZnzgwEuR7XzJlTL3AanBXJu5sUjDPTweOhENcdSHBQ/
sX6Cr9NFRm6bQqOrjmsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAD4IIkxuNITUM
bHU2ebLEMPq8Udhcl9s3mBlaWf3ecDi4Yu+MBTy+ggcRhnq7zqaRVaRxdhtyVVIA
3hFwrWZK38jPGKrI9qZLoqQJu3RFq241jOjVol6zAkYhuqvO8n9AKhShjUFHPfA+
TN1BC8qb30lwYZnELaHdKFM16f7uska2lMY6N8uYqySNWFz/B77zIqUACRnvyGfS
gJ8QRDcGAjA0+SEMKtI0tB0qsL4c+de8uPaUjyO5uzWLXJap50gBi5Zx17YE8aMk
wntulWmvYwJSokZIOVi+3PDSc8Zh2ukhddA716NF2U8c7N++BZFBGcGLZK86Hh76
HayLryJQmA==
-----END CERTIFICATE-----

`
	client_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA5tR1s9rYwyentfQ/pN5Qz36AmBGBV2s6bTkQ+d/inXESYWrT
1sa6lmkP6lWhT/If+iLADiDBi2CTkRO1N9NsXxhSZrotgPJ7ab3zNeLJValOz2mb
FVFRNymeObIhLuZdI0SFKE3bo4JOTwO2Uv1JMWcVZBjZJv09WPANZOteQdqVBTc8
USahD9A9tNsE3aDJ0OdaQrovGBNpID3EzD+Pm6lmHnxqgtXB3wh6/Mqf9ifBAEQL
Sm1VAaXoHkRQppuFNaGrpziAMCFyqF1EVlrhQxw5mfODAS5HtfMmVMvcBqcFcm7m
xSMM9PB46EQ1x1IcFD+xfoKv00VGbptCo6uOawIDAQABAoIBAQCTFY5qrGiy8fHL
33cudvrHPLR0MbNZINp5/oLytdaQvBwaNxgFI1yBuzCJAUdoyb/Wg44dcoHhbgiZ
yRUQHYhQkA7xpnCYWeqJ1p/DFl90Vg4B3CkVzFsT61EHMpoyaFewwViX9gSei8ma
T6M9/mdFM4pN3geA8JzGry/ZvqCxFID3Sz4/08zq9UjS54GiZgJb3lyGazdDk3Gf
h2NukbBRtvdh8iILjEM38czgqTBrDqXlFa5q0p9oq+UPn9twaVZcJ9t4IrcWIgaD
l9cYRE/agXj0cRO/IVOi/RB0e/NLiR0XqXSo4Rx7uGcSJys1yuPt96OGMIh2+c99
VGJbzBsBAoGBAP8qagFe5kJNrjweo9yNhs0H/TFx1mhCqQKPNFouDtttaCDcNvXx
d3B5KYKgWpTJPaZ1eGfPeA4OTLhKCLVG7EVQXUUsztyDS1JpuUJkm1texA30g0sw
UWhLfQfFEgWCaIQkbqZv1J5OYrc2xvPqjHfP+NG1CAte1w5QQ7FA541fAoGBAOeV
rO0yF30sDOUXlixfKN35j2FIgVB0DOT6nkpPyh1OYcdcshu3utGqmOiN7twqwyiL
m3Uucix/JbTb2m+HSAX9/s/SHHOoXeUp21wVSGYesknrBEZt3VifINzu/OFCjLk1
Plx4F0am0WrsDnAtQwgpCV29lgQjmFsXQZlUW051AoGBAJdvpbAgkUmCbsixapCn
0fv3JNZmeFgyT7n8IZbvxNOHkAgIifnXEArJbdBfuMKa2KLlDsuVfuvgormw/pAP
goP0mRZH7JFEvrwvkMqNiQJmMLcTiaRjDb13J8InvHVWmw7pzF2s+yPk44NW2CbE
6g7leAeFiDuvUrTk//e/zGzDAoGBAJx4TLaWubghIzVGkni4cuxHydB5JKYvQucT
Tg/3iR/z7ay9vLltkhRHp7i47UJkwieK7CZok0vtPJTOVvAz/z3NN3VDCWY7w/Uq
KsQ0vQ4Cf4Ph/ql3Ya6XFaUw9Dtes6YPi2r+2PsriyMrCzZP3pKM538msU1qn24s
cG4gyPBhAoGBAL+VTkIaLK07qChlT0Y2hwbmfLwAlOrPguJps7D0C0aBUDPXylOO
S5myV8jp+htbP6Mn5MEzZHhvoVSEe9GiCv9E5KMisJjPtQRRRKGNPAnTt9KJVQ1U
BCggzbZzimK/EFR72woV0071B93C4jO07jEmvkCb3gzmyWkgjREZQusj
-----END RSA PRIVATE KEY-----

`
	ca_crt = `-----BEGIN CERTIFICATE-----
MIIDWzCCAkOgAwIBAgIJAKst9d2m1o1TMA0GCSqGSIb3DQEBCwUAMEMxCzAJBgNV
BAYTAkdCMQ4wDAYDVQQHDAVDaGluYTEPMA0GA1UECgwGZ29ib29rMRMwEQYDVQQD
DApnaXRodWIuY29tMCAXDTE5MTAyMDExMDIzOFoYDzIxMTkwOTI2MTEwMjM4WjBD
MQswCQYDVQQGEwJHQjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazET
MBEGA1UEAwwKZ2l0aHViLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBALh9e8OuCRP0zMYjbCqUk5b+d6J3tC9INL3P0VwcmWx5jCpUQLz0SGafnIL8
LworJfQkbDdOKNol9Zt4vzsxdV1k2VaZuAY0qWG5Kg+n1tCml46By9mQH3B6ngKe
cNDdBmRGYYDkuqI9g8UBgRYT4TbIQJ1Ns4wuKQR02/kCUfWypvE+8bEQEXTRKcHo
inILmcO7RvhWkfwWVbTpUv7M7K8wwIGKawDgl3DeW5g+tss0PD/iCdMo0DMRHykx
4KeTsrPYdxpxgf42LwG0aJ+/28GzYCQ4mYJaTADr5pp+vlUZWtYK8m7fFXbpGlrU
5aLTA5aEPdIuyTa2/DZXl4JBxTkCAwEAAaNQME4wHQYDVR0OBBYEFCikHb0Ms/7f
jci0C5Amwvf7cFmYMB8GA1UdIwQYMBaAFCikHb0Ms/7fjci0C5Amwvf7cFmYMAwG
A1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAG+jH3wHkTqx97/9voaftE/b
0tkbV+9b3SxPv5KoW0fm24x6UNrMPE9APt0J00Vlv20LNc/tOWquyKGDIhhe29/x
ehte/l7doGVW0Wg3xQtiIT5aJdMHNy+bSLogzV5D5sbHcPStKNj3M1wwhMj03YZ7
Nj5ua/c5aUU+MBMv0C/FNPnB+GSeRO2MxYHsZP2mBEJaLhPZ+K29kwGPCVWIESCH
IOS/jew/kfpPLavuvyPqoGAfc1xpe6QQXZUEGCtzTDU/rl/hQWMxCJg85E1S5Usx
gahmAgIzeyFCjb2txOo65VtLM0DfzzkIX2PrLz7CyiXP40m8uBMtCDG+IZS0arQ=
-----END CERTIFICATE-----
`
	grpcUser = `axjBj6l4l49k0EXD473HNuOXZbpZfku4`
	grpcPwd = `tGH57Sz2ml5TT6StPIDeX1l911nPT9M9`
	grepSererName = "server.io"

	_pub1Key1024 = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDF3aII08/G7OdUaJVCHdPbhZl+
T8xGX7dTbDvtMaFpkZrjmA4SC3sM/fvlOQ33HYXnZvwU4IZ80BndbhMuNkYi9aWW
VSMId4VSOALwUx6H/1IVKcjXwAJpKrxdr0R6Wn/E+tDL1G1HSD1DHkbWhbgowrI3
TpwtVty7ROZVBINBgQIDAQAB
-----END PUBLIC KEY-----
`

	_pri1Key1024 = `
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAMXdogjTz8bs51Ro
lUId09uFmX5PzEZft1NsO+0xoWmRmuOYDhILewz9++U5Dfcdhedm/BTghnzQGd1u
Ey42RiL1pZZVIwh3hVI4AvBTHof/UhUpyNfAAmkqvF2vRHpaf8T60MvUbUdIPUMe
RtaFuCjCsjdOnC1W3LtE5lUEg0GBAgMBAAECgYEAv8qONmpBe3vE22+oRfctlRqR
5vqocgpzY9yE1eyGnhKyBSwtb1ZLhxNlqBG+tKqcUenkLMRZ9/+rIpSA6QlYvvLn
bBCYxAsWV/xyFJEh8e5MvS0v+usDDdtSt7nxznD8dJffavahITx0HfA6DxWrVQuz
M8/fRsSIuN0RVaUWnPECQQDxS30rdKPzJxTEJdprDc2jbYUubAdBZ06shPwNCfug
yoOQ637kBWztVE7Cj0S9ld139FLu7UMinQePfKUj6jWlAkEA0eyX++byIPDNtDy0
A96hxHhl0j6of2uG41UfNccHS72CrkVpF4L0Ys+rVsW6s9uZOvOazKXK55DUnRuP
cWMtrQJAMx84RMWwmqqUBr6yWO4SvGZOyjgPDXdSvtBqCmUsD7P4TfLm7m6L1nh/
O09ZVAV1Z523GHHiQGoemPLilgpgFQJAGiQNNQgoRKPX6cbZX9X8bPvVKh41W1Cn
hm2WKlszdGIQAOWR1aSwDBHyMycCPd1tsmKddzh6EOX/I+VHsoX4LQJBAKSH8rFb
F11WHWbxbtDG7oTuqS3Q4+X99WP9LxHYMu+sxas4xCRVyvxVWLW1F6J0FGtACDFV
tUL9797heLPgnE0=
-----END PRIVATE KEY-----
`

	_pub1Key2048 = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3+FVZWd+WcnXN7hk1WWB
L3Hu1Rm3cWLkV/pUMq9jW2JlAJwVmO8aOdIazuKRFmy2Bu7rnIx5nW4SNX/EFIhr
2jQcjP/UBQNs4toR8eDvA60RbTiadJAdaiDzReatCLKhHomyt3HDN9CLql9s2GbP
Cji9K7u7f0TEBl+pcDowOrhZI7JsshnW+trgFyIdNWBlSjvtLD6KWrSjnYy/TQb9
lhKv+MtZeSXD5DusEu4L7DXtvRsiCNVxaY1AyT4C9w2Qou02H2KJ1ORPHGKD8Ix/
PXCp98nP9nU4NhWrPqRFpyV8UyjowvLA4shwwj4o+KgXNCTjbyruZ1YdcaJjl/ux
nQIDAQAB
-----END PUBLIC KEY-----
`

	_pri1Key2048 = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDf4VVlZ35Zydc3
uGTVZYEvce7VGbdxYuRX+lQyr2NbYmUAnBWY7xo50hrO4pEWbLYG7uucjHmdbhI1
f8QUiGvaNByM/9QFA2zi2hHx4O8DrRFtOJp0kB1qIPNF5q0IsqEeibK3ccM30Iuq
X2zYZs8KOL0ru7t/RMQGX6lwOjA6uFkjsmyyGdb62uAXIh01YGVKO+0sPopatKOd
jL9NBv2WEq/4y1l5JcPkO6wS7gvsNe29GyII1XFpjUDJPgL3DZCi7TYfYonU5E8c
YoPwjH89cKn3yc/2dTg2Fas+pEWnJXxTKOjC8sDiyHDCPij4qBc0JONvKu5nVh1x
omOX+7GdAgMBAAECggEAIMVCE8LPauMxnpVeJSJjg4dg106ZXH6GQB6DXpvvpjvD
3w/51VYCd746cFgXtrmY93DXiiXB03p+LdiS4hKJ/vmryDPWXBmBQb976vTq55XY
vC0R5sgFljhWg7/dSi2jie3L/DApzCy5lOm86/w4iB2ACzvCmUF+lBRCoAvUbXOy
bsEovqck8gIf1hlA0KTYxNWB8x5ZkGcFdpazBeCXMPixp82UKEq7Bvp9f7Yw8T1Y
RXJlaFWm3XAp1keDr1Cf1C/OZ4quclK5P5Tk5KfpUoRVsTSO/yn1KBbNZtlI2Ngx
wRV6SrBJN8qa5oH40nsprS536nL4EZJUYuNLffZsIQKBgQD4PF1hDr5Yt7kb29tx
jtAX0JyxXh3oycZ25NEti0Yi/Mcnf8fjIemU0s/YewySMU5vuRBUuvCoWEQLZQdS
CGc4G5Q4htkBvR1WtbBRhuxV+KyEz435HNZfeoyXoRbbB6P5iI+F2UymyEefwpx9
iMVqAg5wD1yy6ExNPcNRg8qaKQKBgQDm4fPXXXFVS4Jthbpbb1UaDzeQVILDtEzw
0ljXT+uu7OEJ6fFlldf+4XPR05zSQ16UK7u8ICbvXBYfvLjd9M5YL9OteY1mwwR+
Rv05w2cpBFrjtVfhiCz35pTOL/zQAA4hfIJmTCHEmV+uMmRWDlka+SK5RHSUdA9I
+my13xayVQKBgQCTWZG8aToIA5a3uLv8Hl/boxNAHbP3WL6cGJsqQ7/wSMgW5DzM
0HaMxs5lnDUMGoSKwPm2sfjklPBfKys7QI20unozS6hI5e8iZ1swKbzkE2akt2d4
9esyZdZKs26TuWdWWf+H3kMnxT7u0GCAC83TbUEQt6247TdNqlnkayy6cQKBgQCS
5VYAQ2qVKyq7xiawgCA0KVRf1vUv0OpXGm3959J7BCmV5it5R6Iaf5Tx/mI7gTOO
sFiMtCQxjHRjEu7IATa78woJyFmH9TJJqZ75fnKHLUcqs7lLPBnoS+OHYA7IxBA5
i/9nWK7vZ+nxagxemFhnCfXmzEAkJ8eF1hcOi/bj1QKBgB5p034FRic8mw6k0Q0j
8peNnh8klaM/kTnR/gtNMjS8M3VG10OiSMh8eHq7ebA3Jm2rFTl/BdSLs+gaqgu2
A87Bh481IZgONR1KbcSKfVAl/7ZCGobqdxkllO/S7yry/JHxITLmrl6fqYd0VgKS
iVVMOBJCfKb7b1IAnrtAlH7k
-----END PRIVATE KEY-----
`

	_pubKey4096 = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA1sdVV9n/8tneDlrQDhyE
bLR5KrdeLFmTLsrddMmZjQaBh4hc+eFEMTkoOaClqkajpjLH4RowcJRbcGrogxv4
pW5EsVVtqD4iKsTuaW9C59/rDDAmiszQJQaTtF7Jhp/6Pt1RJAZy+areFqLrcxZw
pcB19NS9JxTjxyAjPB5wPi05eb6t1TV1wmiIP66eYVm+35OLhsPUAjGaYHOmaNs8
cgltNSQ8hZEyTkL+qB7NEkLboOBQKNkXdyMNtcAsi+iWunBp/mY5LI1D1m6uCIfr
2zVvL7qvIGdZMWquPhypE3QD65QHrXC4u2x+1aliL1r1i0yDAiofX87qmLkT3/fP
DYLE5jmABEhiCNO0Pu/HXXDWfCc9k/oLS0r92l9s60h9hj8QfoR5zCY6kXLUrq5E
dWhZGGDKyUnmUhp3iqS0NOBvibMSt3pM6T7kUIRawWkmialWfeKjAqV15dIeKaQ/
cAkyNecDYn8LWYVn7Fy/Qq++H4N0PjTA52YHhsamwpMuYt9CuS6uxLiyD+3QnALy
+rHr16oQXAZxzbHrwnjB8aflFFL9RUnY7lLHyvQFAZaDTXdYYYNozZPOreGruP8q
EvO6BgxBoP+1/f8Ji6O1hAdcOPqKec894gheerE9dqGw4AKyfC+JVh81Quk6Ajza
0f0P0SFQrqLA6hqoHowqg8cCAwEAAQ==
-----END PUBLIC KEY-----
`

	_priKey4096 = `
-----BEGIN PRIVATE KEY-----
MIIJRAIBADANBgkqhkiG9w0BAQEFAASCCS4wggkqAgEAAoICAQDWx1VX2f/y2d4O
WtAOHIRstHkqt14sWZMuyt10yZmNBoGHiFz54UQxOSg5oKWqRqOmMsfhGjBwlFtw
auiDG/ilbkSxVW2oPiIqxO5pb0Ln3+sMMCaKzNAlBpO0XsmGn/o+3VEkBnL5qt4W
outzFnClwHX01L0nFOPHICM8HnA+LTl5vq3VNXXCaIg/rp5hWb7fk4uGw9QCMZpg
c6Zo2zxyCW01JDyFkTJOQv6oHs0SQtug4FAo2Rd3Iw21wCyL6Ja6cGn+ZjksjUPW
bq4Ih+vbNW8vuq8gZ1kxaq4+HKkTdAPrlAetcLi7bH7VqWIvWvWLTIMCKh9fzuqY
uRPf988NgsTmOYAESGII07Q+78ddcNZ8Jz2T+gtLSv3aX2zrSH2GPxB+hHnMJjqR
ctSurkR1aFkYYMrJSeZSGneKpLQ04G+JsxK3ekzpPuRQhFrBaSaJqVZ94qMCpXXl
0h4ppD9wCTI15wNifwtZhWfsXL9Cr74fg3Q+NMDnZgeGxqbCky5i30K5Lq7EuLIP
7dCcAvL6sevXqhBcBnHNsevCeMHxp+UUUv1FSdjuUsfK9AUBloNNd1hhg2jNk86t
4au4/yoS87oGDEGg/7X9/wmLo7WEB1w4+op5zz3iCF56sT12obDgArJ8L4lWHzVC
6ToCPNrR/Q/RIVCuosDqGqgejCqDxwIDAQABAoICADIa9Jz3HY/RJc2hf/Ia0wXt
IGtHte+QwhZje0B4m5rbzrIIrPAajmcRV4ICKUPNEPZ/2EN6cZyB78cNGcskZmBp
lhrsvBVI0X26zYfJTgl8IoCIZyVwXIqWuzST/F2sypuJ1BkcbAw0wXT0cws5S/RP
LvV7/9izNeRJag7nZvYKZOMzCai4vQ0qh8abfRVm83GDIUTCQJ52ZfZkZIkHxFUy
P8jq+DeMxPifBnvAG8VL1aL1UZ4F70R66ALjn0DQdQFvojqYLHRpTE8lKPKSiwJr
t9GhsqNTmOo/YgDZfNQt95Aoy5W5u072I6zCxEYZ6TijE9kYbJNUWURhwPI6BJJL
0y2dGE5MGFFg5iu0ffvzNFCt4/05hmB9eEXXylBR8QIqZ5Ls8gx4rpjsQcm1Zna1
xPI3d+o0bvZPTAaI9JDZdTDXhLPrNN2nf07HhvvkZmlpEw0lY8ej6+xwfDJE2iCY
AiOwIhz8zsPOQ+BJj6MXHZlnIMUOnFtKOG9589vk7uV+ZWLl7e4Aku1qHOrJN74+
nMsenUG3KWYQjrfF3JC+fI7nTfe/dc3Q1QOjdhZFLaqsL3nNjztJuw8c5ETAdcb6
sy+0eRqpZbhVotnq0DyF/8/pQJNB0KQZK1BF4rF656QEHozeI7OrpjtJMMoplBiR
DDEdaolIP53wGEh2b5v5AoIBAQD0p2o2PopsVpeHOfRszloY4PdjQiR4TzKzRVFV
bFtiakm3LxvilUbiJknMDkXQhABIXsHmYuNNCuZGijNq7Qyfr67gxC9uk9frlNU0
zhZpvYYXJyxQ11t+ZoeAFk9k4cqM+DMmr10QMU8qNDK+bx1X02biXCu++/GKc11A
SDNhZlzmSz0c2BsBSW6hQ9aOLfX4bAeV+9tJNmb3Zh1NH6uwMLZvlWyL3OzAKfhZ
TlNt0vXbYerCMP0FgmDLJMHd3BBVQtaACwV1q5DYAPR5zDJDfZlfvDRadhl8H3gH
coOAPMoDSOZ7xp5MS2cV/vtbamoG+tsDmwWk2du3wWUl06P9AoIBAQDgvTtxapN8
7ZjCyDUzum7e7bN6YX0ccJjJrB9qnOIB4awW7eeBLNkTEXTuCukBYVlxbUCnWD+H
djDeHl7CxrVa0bt/jGkA+hIKEpQx93iompim7EnX1lGC/gP6t2sMpDLNik26BN64
2Q8cbxLSzpYNRuwOC3XyZ79yVIXpgV9xoUU6CUMNBOT7ZAraWY2LohUlAyQX5DUd
sVNwrzLMfQqGeePTjRnvtHzU9ZZRWH/bbsqtSlXI/JftlSyFiLIqdHJ+rrxasAQ7
QOY5DAwX14xdioC1+STzDJHH+2g/BDVToiZxw9GZ1M4JkVnEU7w2ccKLU2W0ful+
WRSE/Q0szjgTAoIBAQDu4QBN0qbpvWrauHW2P42tOQuUOSLO7dV9QTN3CwP3hfxQ
BoldpY++hNANk+oK/Lgh8ZO11dxGf1v0iEBIKQjoamuAP05o06ZB+eJrWsZ7nHfu
52rXzE8jjgzDvgTrZaOWHUokfZmKk/rOJIVfd7LY7CtK3eBA7FMdciMc/uJcOcx/
d/tFzKQhj4ebolc+IBZI54JIqc+lHp9O9L+rbD8BG68mKGoB7kakItbArD+9vfwc
pvDHh3mmBXVpJIy+iX7RIR+7igdcq5YTsmsC+aQiTeKRnXwoz6N2lGtoKiHH9pLw
vh99v6MUr5MJ4RugWLkJ86ohTR4npihotUaDtrApAoIBAQCTxxrNSz1MSpfGjQue
xhqdcEQyVuSDzO5KvmmyGxLqFdCpCyrNYAYlabcvx/DLPY4o7aQz5e1wT6F2jRXW
kf8yhvL5vgRV5hnykaDs8kNe6rkyGfG8gWr77bgEJpO3rkjRqv3NMeKaPfCXy9ne
0IUOmfIikhqumNXkgfvEPZPbDiaNMQXsC6nePDx+s6BFjwDEY7paE29x5OZvFGUc
3aQMJR5QP1osqsvi5NJBDyaTzdhr9pNOI/pq+UpbTDWLgSLAdnnYUCGYLOa14Fwb
WVstLyPPhNJtF3jMvV4hAc1m/xq0eATdWHdbBz61wDHkww0fvGkGNOWodT7u6868
BaYNAoIBAQC9cDs5jUq4wAEwO2H8Si9z97WDC1mkPYdJtDWVlI+DQJIgDWXtkOoc
EW0I2SERWOOSDT9kluDwdyNkUxbtggcGYXo8VBBF+m6ZHzk24lYXMppV77SUkcCg
3VmL4j5vjtasTYxdUfLwXceNCJOYriD3JMN9oz5wY/Fwt+CTknSmCdGW6DaNsQzK
+346fWLfutx0r983BmsLMVQaX4A2ZZsunUQWp3EYZx70emBQjl5BGUWAQeosZjNi
rqjiBBMiVg40yXsgp9uW6fmnw6I18VAM0/pZqSSa2O+5/+Hjj00P4YrqOHwMG4/J
naqfw+uOy44bUZt9ghSqSrfB81bOzgzf
-----END PRIVATE KEY-----
`

	_pub2Key1024 = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4lQ3os+rsFyhzJmFiIJslgGpO
WUPjw+AkVA4ywDDHMQLCrRVbfuAAhWqauBfPJbXT7p7P9vAn9lA42qHGanyvLCLy
jp/MYqJb5pBR5V7fquMcsIPQYhODBCy2p1zaoBHMKr5BvcviEmitj/Gl/xkz1BxX
fxyDSqgqJxZAUxfXtwIDAQAB
-----END PUBLIC KEY-----
`

	_pri2Key1024 = `
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALiVDeiz6uwXKHMm
YWIgmyWAak5ZQ+PD4CRUDjLAMMcxAsKtFVt+4ACFapq4F88ltdPuns/28Cf2UDja
ocZqfK8sIvKOn8xiolvmkFHlXt+q4xywg9BiE4MELLanXNqgEcwqvkG9y+ISaK2P
8aX/GTPUHFd/HINKqConFkBTF9e3AgMBAAECgYA+H2RsAkm5pd2mS6+Q4Bp3V63v
qplvydfhQiz5JbgFAljEfo5mmd/4LO8BHZ5dyHpW1sO6iIixWnQLfoHeDq7hj62a
zzwccG3Q/LhEGjIlrSXMJfA8P9lWYDMuOVLStT/F+E888hoxAUYF4HeBUsj+NXGm
SCgnHDTx8WXODQlcMQJBANm/DDW1ggmezoJqC+UPSml4WPeIlWNiOWQ68PnarhGV
4VpL0zrXhmIluxYLvCS/OQ9B0vPbgndKWjOHKnCSMd8CQQDZAnuGXcnWKVxcKVP1
TiI71yzs0n18gnEmSttvCAF/E7CwYFuT5esu7qkYcn3IizOxWFVJDuVQhy/fmVoS
24UpAkAuw5wgsNGztTqOwa26TRVjH2ikCN5kkMTYpNv6HSADQNg8J0q/OWhwDcBn
VK/ciID9qNpgawVTD1Hd/Sp9MLirAkBj/x/ad6dE90Qm96hHdhySRIHgEtJeKGFp
Sr84t5Cw9OrLK2uniB/KPZFPwZoyaeqFAvYxtxp19AVcXHbED8GZAkEAntZ2SIQY
IhCwH3eO/iejoeXqz/OfJ1NZPUa8Wa39Vgy0wstHCjIdEG/Zv0ltzzEB5a3EM570
ETCYO2wPZPwN3Q==
-----END PRIVATE KEY-----
`

	_pub2Key2048 = `
-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArh12GFOrIeteFg+AyXvm
ryUqjwHFpendla03N3us8LGB67MpvjZM56mjBV1U5RVqDhd9wBXMctG68H/EZH5u
i1tFrbRMo2KIpB21pWwTHESefsyN/DFL/hCD896vDy0OVbZLSlEI9AD+mvhz16Qb
cOP6WxVfxXvO88KRQjoGPrho6rx5Q7OwJpSYOKsje7ycV/nuLKIrpDCwvhpY8+gr
ah7UStGDqa+S17T2dPdtOLTNYIcpJx2wcS0rXPlrIUHqPxCv7uui1I5YWAtl+Lit
pzZJgLOL/jq1TA46f4TFRLX4JfEWAhHwfmttELKYJASuBeB9e2cgNBDd/0GFOsVt
0QIDAQAB
-----END RSA PUBLIC KEY-----
`

	_pri2Key2048 = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEArh12GFOrIeteFg+AyXvmryUqjwHFpendla03N3us8LGB67Mp
vjZM56mjBV1U5RVqDhd9wBXMctG68H/EZH5ui1tFrbRMo2KIpB21pWwTHESefsyN
/DFL/hCD896vDy0OVbZLSlEI9AD+mvhz16QbcOP6WxVfxXvO88KRQjoGPrho6rx5
Q7OwJpSYOKsje7ycV/nuLKIrpDCwvhpY8+grah7UStGDqa+S17T2dPdtOLTNYIcp
Jx2wcS0rXPlrIUHqPxCv7uui1I5YWAtl+LitpzZJgLOL/jq1TA46f4TFRLX4JfEW
AhHwfmttELKYJASuBeB9e2cgNBDd/0GFOsVt0QIDAQABAoIBAQCb3lnzKyufQOJI
Y3aKaLW5g08XGKIEhljMfnVY4QmPq4jAJPKwilHMbas3yTaPod0AYn07cQhGnYR5
ehepUxnI/VtiRm75MONb8BDF1vtAqhktMBfHdaYu+j/2GBqlPlN/3aKHFAYs0Zsb
xmGF4S6DoENmOLs0wkIhK8P4ApPGnrqtXHPh7v3HVzTCylFs6IUCIDLzgbpGZu7O
AoXwpxwQMMSN4EjvXqwGKZmg6IBZzHPoE0NpGYQVPGrjndnN93ECWSD0YU+JYhw5
SD8jIH4lFYJeRHK5cV9grKP4y40G0N3tTCAqOGFRfkI6ZSEoc+RFRFqAEhG8HWOL
n3npcGNBAoGBAMSQSk3KKnYokoTdLG5xy97s5VOEWRSaBh65LScX/nnqsjymIJjI
1X9whhQLbsA6uktU6oJNKh342rGM4NnOqLJoYGllmHjc6/SCyqxGp0ay2H5rbHKg
ZOyCq4BUvgBzt9CRRhf2r7cEMCBNYqXZFKzEVPDknQQyo48RWYuDbg+pAoGBAOLD
dfuAS6wjkr7riFLErwokgUNDSkp4pfSC9rdPd0NFYLY1/9yxuC2++keIhK18/Jlw
oBaKlbIgxX8EsbyI6fAby+dWJ0N0pQyNppWPY4ezJ+uFpposDUrveoXRFUtiS+ve
0gWCQ4wKAnQxVMHG9uDp0b/kjFvH+QTnMw0FP+XpAoGBAJDFqURkCyQdu9SJxejO
fZaCKmF5z9ZhnvJP9tadUHthBcevn8CH4t9K9CWdSgPg/UbwkwxHYybSG9i7Zvxk
vlEwmRnnjwYtyMe88SMzoo5quRNbcXN3eP3NPB13zL0ufYrrBJIvybllJ0ETXf3C
xfx9WgZWiuMFnPuJjsc3lP+JAoGAZ6LtSQRZkVKwvpDmvO0nEnucmCEo0uBQ+G7i
UuT+nMAYcy46waJ3inC98fNyr9dvmrDeeW7c+4v+tw5uLLxmLlaF2jSFvU6SICqc
972Qv3QhyoJKoit/57+LP51PHiTOjf5H/jyKonXwqSnikq1cJ261bf4GJ+w84wDH
VCwSCAECgYEAhfbpOp27Osf6scepuZx+jabEETeI2ltUCeXAWbwxU8Yf6GVSzaVB
6pypSMWZlbWycoliA0/q9SsPUYswa6AV4ciGx4lVWfiFWj8FoGl+YGl545sgT4Vp
HH8xelnu80nL0dt17gCL8BDUgNgmcbABVCi4M5MGPix873CIOhxJclI=
-----END RSA PRIVATE KEY-----
`

	// AES密钥，16或32字节 选择AES-128或AES-256。
	_aes1Key32 = `ac88fc9388c4790ad0cf9a98fc6zpr84`
	_aes2Key32 = `f7e8b819l0ad0ccf9a9g8fc5e8c4765q`
)

var hw   =&hardware{}

func main() {
	cfg:=&odinX.Config{
		EmbedCtx:                   context.TODO(),
		EmbedName:                  config.Cfg.Name,
		EmbedDir:                   config.Cfg.Dir,
		EmbedClientAddr:            config.Cfg.LocalClientAddr(),
		EmbedPeerAddr:              config.Cfg.LocalPeerAddr(),
		EmbedClusterToken:          clusterToken,
		EmbedClusterState:          config.Cfg.State,
		EmbedCluster:               config.Cfg.AllPeerAddr(),
		EmbedAuthPwd:               embedAuthPwd,

		EtcdCliCtx:                 context.TODO(),
		EtcdCliAddr:                "127.0.0.1:"+config.Cfg.LocalClientPort(),
		EtcdCliUser:                "root",
		EtcdCliTimeout:             3*time.Second,
		StoreLicenseKey:            storeLicenseKey,
		StoreClearLicenseKey:       storeClearLicenseKey,
		StoreClientConfigKeyPrefix: storeClientConfigKeyPrefix,
		StoreClientKeyPrefix:       storeClientKeyPrefix,
		StoreTokenKey:              storeTokenKey,
		StoreSerialNumKey:          storeSerialNumKey,

		GRpcServerCrt:              server_crt,
		GRpcServerKey:              server_key,
		GRpcClientCrt:              client_crt,
		GRpcClientKey:              client_key,
		GRpcCaCrt:                  ca_crt,
		GRpcUser:                   grpcUser,
		GRpcPwd:                    grpcPwd,
		GRpcServerName:             grepSererName,
		GRpcAllNode:                config.Cfg.AllGRpcAddr(),
		GRpcListen:                 config.Cfg.LocalGRpcAddr(),
		RestfulPwd:                 config.Cfg.Pwd,

		NodeName:                   config.Cfg.LocalName(),
		NodeAddr:                   config.Cfg.LocalGRpcAddr(),
		NodeHardware:               hw,

		// odin & edda
		LicenseEncrypt:             PriEncryptRsa2048Aes256,
		LicenseDecrypt:             PubDecryptRsa2048Aes256,
		SerialEncrypt:              PriEncryptRsa2048Aes256,
		SerialDecrypt:              PubDecryptRsa2048Aes256,
		UntiedEncrypt:              PriEncryptRsa2048Aes256,
		UntiedDecrypt:              PubDecryptRsa2048Aes256,
		TokenHash:                  HashFunc,

		// odin & app
		VerifyDecrypt:              PriDecryptRsa2048,
		CipherEncrypt:              Aes256key1,
		AuthEncrypt:                Aes256key2,
		UuidHash:                   HashFunc,

	}
	odinX.Start(cfg)
}

func PriEncryptRsa2048Aes256(src []byte)([]byte,error)  {
	return endeaesrsa.PriEncrypt(src, []byte(_pri1Key2048), []byte(_aes1Key32))
}

func PubDecryptRsa2048Aes256(src []byte)([]byte,error)  {
	return endeaesrsa.PubDecrypt(src, []byte(_pub1Key2048), []byte(_aes1Key32))
}

func PriDecryptRsa2048(src []byte)([]byte,error)  {
	return endersa.PriDecrypt(src, []byte(_pri2Key2048))
}

func Aes256key1(src []byte)([]byte,error)  {
	return endeaes.AesCbcEncrypt(src, []byte(_aes1Key32))
}

func Aes256key2(src []byte)([]byte,error)  {
	return endeaes.AesCbcEncrypt(src, []byte(_aes2Key32))
}

func HashFunc(src []byte) string {
	return  utils.Sha256sum(src,[]byte(storeHashSalt))
}


type hardware struct {
	// linux
	// sysinfo.SysInfo  // "github.com/zcalusic/sysinfo"
	// windows
	winsysinfo.SysInfo  // "github.com/offer365/example/winsysinfo"
}

func (h *hardware) HostInfo() (machineID, architecture, hypervisor string) {
	return h.Node.MachineID,h.OS.Architecture,h.Node.Hypervisor
}

func (h *hardware) ProductInfo() (name, serial, vendor string) {
	return h.Product.Name,h.Product.Serial,h.Product.Vendor
}

func (h *hardware) BoardInfo() (name, serial, vendor string) {
	return h.Board.Name,h.Board.Serial,h.Board.Vendor
}

func (h *hardware) BiosInfo() (vendor string) {
	return h.BIOS.Vendor
}

func (h *hardware) CpuInfo() (vendor, model string, threads, cache, cores, cpus, speed uint32) {
	return h.CPU.Vendor,h.CPU.Model,uint32(h.CPU.Threads),
		uint32(h.CPU.Cache),uint32(h.CPU.Cores),
		uint32(h.CPU.Cpus),uint32(h.CPU.Speed)
}

func (h *hardware) MemInfo() (speed uint32, tp string) {
	return uint32(h.Memory.Speed),h.Memory.Type
}

func (h *hardware) NetworksInfo() (drivers []*odinX.NetDriver) {
	for _, val := range h.Network {
		nw := new(odinX.NetDriver)
		nw.Speed = uint32(val.Speed)
		nw.Macaddress = val.MACAddress
		nw.Driver = val.Driver
		drivers=append(drivers, nw)
	}
	return
}