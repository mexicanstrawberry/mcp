package web

import (
	"crypto/tls"
	"net/http"

	clog "github.com/morriswinkler/cloudglog"
	"rsc.io/letsencrypt"
)

var Localhost bool

var sslManager letsencrypt.Manager

// load letsencrypt
func init() {
	if err := sslManager.CacheFile("letsencrypt.cache"); err != nil {
		clog.Fatalln("[web][letsencrypt]", err)

	}

	sslManager.SetHosts([]string{"localhost", "app.mexicanstrawberry.com"})

}

func ListenAndServeTLS(addr string, handler http.Handler) error {

	var getCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)

	if Localhost {
		getCertificate = localhostGetCertificate
	} else {
		getCertificate = sslManager.GetCertificate
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: getCertificate,
		},
	}
	return srv.ListenAndServeTLS("", "")
}

func localhostGetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {

	cert, err := tls.X509KeyPair([]byte(localCertFile), []byte(localKeyFile))
	if err != nil {
		return &cert, err

	}

	return &cert, nil
}

var localCertFile string = `-----BEGIN CERTIFICATE-----
MIIEszCCA5ugAwIBAgIJAN7ZEgncocHAMA0GCSqGSIb3DQEBCwUAMIGXMQswCQYD
VQQGEwJERTEPMA0GA1UECBMGQkVSTElOMQ8wDQYDVQQHEwZCRVJMSU4xGjAYBgNV
BAoTEU1leGljYW5zdHJhd2JlcnJ5MQwwCgYDVQQLEwNERVYxEjAQBgNVBAMTCWxv
Y2FsaG9zdDEoMCYGCSqGSIb3DQEJARYZZGV2QG1leGljYW5zdHJhd2JlcnJ5LmNv
bTAeFw0xNjA5MTQxMzI1MDdaFw0yNjA5MTIxMzI1MDdaMIGXMQswCQYDVQQGEwJE
RTEPMA0GA1UECBMGQkVSTElOMQ8wDQYDVQQHEwZCRVJMSU4xGjAYBgNVBAoTEU1l
eGljYW5zdHJhd2JlcnJ5MQwwCgYDVQQLEwNERVYxEjAQBgNVBAMTCWxvY2FsaG9z
dDEoMCYGCSqGSIb3DQEJARYZZGV2QG1leGljYW5zdHJhd2JlcnJ5LmNvbTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJzzwxrHKOk74wc9efrX0wcUvG0
hyGCKQsNhRpTCYHTDHvw3k0S3snYXPL0M0S/glbQyzG20KZNoRnKRmZRos5iTVzr
cCIOWyQ5R4AA++0vpsM5J8x7c5pUl5PyCl2bVNnsJ5Ks4jofK8oWmuzl3n+hBig1
K6HceXGjdgGjUDvx6TsMUzGXAK4DmBUXBGghJUjFfM8muHkQ9MlM8pf8zjaws/jZ
8BCVfheXjfCyxymiFU69LvcE0vuj9XJDs3d22KuslFSxEwKBXvx7zhFWZ+ugC0BG
3HAsb1Ap0Mj70/Yfqg/iWiN57JjOjCFDWCCBcu7fZ7jUF52Or92loV2Ccu0CAwEA
AaOB/zCB/DAdBgNVHQ4EFgQUJDg3Z7O04pJpyoK5vVlC/lScIA4wgcwGA1UdIwSB
xDCBwYAUJDg3Z7O04pJpyoK5vVlC/lScIA6hgZ2kgZowgZcxCzAJBgNVBAYTAkRF
MQ8wDQYDVQQIEwZCRVJMSU4xDzANBgNVBAcTBkJFUkxJTjEaMBgGA1UEChMRTWV4
aWNhbnN0cmF3YmVycnkxDDAKBgNVBAsTA0RFVjESMBAGA1UEAxMJbG9jYWxob3N0
MSgwJgYJKoZIhvcNAQkBFhlkZXZAbWV4aWNhbnN0cmF3YmVycnkuY29tggkA3tkS
CdyhwcAwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAgKj51gFk6NoY
PpGRBlY+m7DkzP0xLHFmJrhxDrDjI/YAa29QOapFY4IfyPgvovjfr9ztNRrd0Fxg
0mHFxvXs3xw3FMIGN5De3AJiBrQjBRNxMtBQbbPGC+MDAUWy0Pn1mbayiLvTSZn8
osFWO3NDqSnWSXEAwCZ9O5RcCboplssAXKpBh5Si2R0E4DJEMK0HBbivRweZtNyZ
+0F+FCpS3gzZwTJVzyDQG0A915lwjLztkYHHh7PMnf77NPK4fnQHVECcOKivqOoq
iklh6/iSqpO9c201sEE+5WJLuzh6FDnoJKvEuWmS6qYWgeRwKqNan5meISAQEpXu
1WfILNMbdQ==
-----END CERTIFICATE-----}`

var localKeyFile string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsnPPDGsco6TvjBz15+tfTBxS8bSHIYIpCw2FGlMJgdMMe/De
TRLeydhc8vQzRL+CVtDLMbbQpk2hGcpGZlGizmJNXOtwIg5bJDlHgAD77S+mwzkn
zHtzmlSXk/IKXZtU2ewnkqziOh8ryhaa7OXef6EGKDUrodx5caN2AaNQO/HpOwxT
MZcArgOYFRcEaCElSMV8zya4eRD0yUzyl/zONrCz+NnwEJV+F5eN8LLHKaIVTr0u
9wTS+6P1ckOzd3bYq6yUVLETAoFe/HvOEVZn66ALQEbccCxvUCnQyPvT9h+qD+Ja
I3nsmM6MIUNYIIFy7t9nuNQXnY6v3aWhXYJy7QIDAQABAoIBAFRNQVKshysHj+Kx
C7o0ByD9gHGOxwedZaZDDM4SzDr4aL1kXKAsefMAs2hS1KV1ky1QFa22n3rw0VpN
pFRR3IeDCOkMkDyGa6gBJzXhQSIbkLxJE/QVndcaf0D05tCxwLPyS/+OjJDIiPc/
FpEzRpkkiLQV6jbc4MI+ZlD/xbeLFZzOkFGtDcr2hv/ttUU55q14lFcEr9SVivOb
uHdG1rh6dEjQkoBnjein32Qs98/DjSWQWD+pqbyReXlRAP4rlLUumNWjXFKshRyi
54UIsjDuhCqAofXkCqxFLUyuP8z7iQGSZ1mSNsuqgNiARQAaRe/0njHod2lOaFQH
zE4SCgECgYEA7E3z9kU93NxF4WxynlA+bXvRunneC0uvy1yzyJc0/hADKjfjeHCF
TpBtdCE3FLKUHd/U1HRvd9QKzuDtz55x11/iEZ8n3i7TQUCwMoByznbfXxbuyR4u
WhjgSqXy2Fn0OeTIM3r2MoLjZCVE2O3eLN4TvmGv7HLPLETElODcuhUCgYEAwVN1
7DOt8s1W5n91QdGkLI9wXtjY9XfzwPbL3gnIW92BpC/KL4OJHx7/K06FY2B4OpDt
7aCZgw7D9z9a9ufdjlI/FUMgzj6jNAQ6tIFGodf4uba7mBhkemw1VdUD2i15Lk6Y
sCamQQyDugS3ePoE/kZKJ8PWKsFjrzxmiVRPQ3kCgYEAjIYK50/j6vx+/gAc5TJ4
/WidnxQrzHHU982ICGiLFe71wtx7hDr9u2u9+0ppVACifmWGTlVzmEHbr40pPsdN
kbOuX6ZS8hjMfkh2v4GNRGSCjyy3EZjGHcQfVaT8FlbgGrGHsL2VvRIDIaHcIFjM
P8hM23GCSc04kG3QrWxPNsUCgYAwGENN781mig8EaNES/sSJEWYzMl9HMgBCESPG
qUhfEkwePIVgLKkARQXWEEK+5lECwOtwInQOVq4J5IkMw8Iqlet7rqeKp6qSVjsE
jOS1frUx/nPM8sSMcD8Ui1nZ/VYYXxU9PWA+7o4WyPWb8xcq6vGn0uCE4neaMLyR
jZfqgQKBgQChjUcHU4NHO4vgYOgtSeopZ90oVPkYxiyiIKpIvr0lh7uwYt6ihDFU
XX3L5SZhBT9E4yO1G4DzDjMsF706SUGYKl8mStbv1BT8vTFhwSwixypMr0NhoL9b
i2gAbYU89sFCUuf0DtWVgsy3SGV+cBbhevB0lTNuH9km66vsjjZiTA==
-----END RSA PRIVATE KEY-----`
