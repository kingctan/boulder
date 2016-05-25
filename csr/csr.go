package csr

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/letsencrypt/boulder/core"
)

// maxCNLength is the maximum length allowed for the common name as specified in RFC 5280
const maxCNLength = 64

// This map is used to detect algorithms in crypto/x509 that
// are no longer considered sufficiently strong.
// * No MD2, MD5, or SHA-1
// * No DSA
//
// SHA1WithRSA is allowed because there's still a fair bit of it
// out there, but we should try to remove it soon.
var badSignatureAlgorithms = map[x509.SignatureAlgorithm]bool{
	x509.UnknownSignatureAlgorithm: true,
	x509.MD2WithRSA:                true,
	x509.MD5WithRSA:                true,
	x509.DSAWithSHA1:               true,
	x509.DSAWithSHA256:             true,
	x509.ECDSAWithSHA1:             true,
}

// VerifyCSR checks the validity of a x509.CertificateRequest
func VerifyCSR(csr *x509.CertificateRequest, maxNames int, keyPolicy *core.KeyPolicy, pa core.PolicyAuthority, regID int64) error {
	key, ok := csr.PublicKey.(crypto.PublicKey)
	if !ok {
		return errors.New("invalid public key in CSR")
	}
	if err := keyPolicy.GoodKey(key); err != nil {
		return fmt.Errorf("invalid public key in CSR: %s", err)
	}
	if badSignatureAlgorithms[csr.SignatureAlgorithm] {
		// go1.6 provides a stringer for x509.SignatureAlgorithm but 1.5.x
		// does not
		return errors.New("signature algorithm not supported")
	}
	if err := csr.CheckSignature(); err != nil {
		return errors.New("invalid signature on CSR")
	}
	if len(csr.DNSNames) == 0 && csr.Subject.CommonName == "" {
		return errors.New("at least one DNS name is required")
	}
	if len(csr.Subject.CommonName) > maxCNLength {
		return fmt.Errorf("CN was longer than %d bytes", maxCNLength)
	}
	if maxNames > 0 && len(csr.DNSNames) > maxNames {
		return fmt.Errorf("CSR contains more than %d DNS names", maxNames)
	}
	badNames := []string{}
	for _, name := range csr.DNSNames {
		if err := pa.WillingToIssue(core.AcmeIdentifier{
			Type:  core.IdentifierDNS,
			Value: name,
		}, regID); err != nil {
			badNames = append(badNames, name)
		}
	}
	if len(badNames) > 0 {
		return fmt.Errorf("policy forbids issuing for: %s", strings.Join(badNames, ", "))
	}
	return nil
}

// NormalizeCSR deduplicates and lowers the case ofdNSNames and lowers the case of the subject CN. If forceCNFromSAN is true it will
// also hoist a dNSName into the CN if the latter is empty
func NormalizeCSR(csr *x509.CertificateRequest, forceCNFromSAN bool) {
	if forceCNFromSAN && csr.Subject.CommonName == "" {
		if len(csr.DNSNames) > 0 {
			csr.Subject.CommonName = csr.DNSNames[0]
		}
	} else if csr.Subject.CommonName != "" {
		csr.DNSNames = append(csr.DNSNames, csr.Subject.CommonName)
	}
	csr.Subject.CommonName = strings.ToLower(csr.Subject.CommonName)
	csr.DNSNames = core.UniqueLowerNames(csr.DNSNames)
}
