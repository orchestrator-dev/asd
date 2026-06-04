package handlers

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type CertHandler struct{}

func (h *CertHandler) CanHandle(mime, ext string) bool {
	ext = strings.ToLower(ext)
	return ext == ".pem" || ext == ".crt" || ext == ".cer" || ext == ".der"
}

func (h *CertHandler) Render(w io.Writer, r io.Reader, meta FileMeta, opts Options) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if opts.Flat {
		_, err = w.Write(data)
		return err
	}

	var certs []*x509.Certificate

	// Try PEM decoding
	block, rest := pem.Decode(data)
	if block != nil {
		for block != nil {
			if block.Type == "CERTIFICATE" {
				cert, err := x509.ParseCertificate(block.Bytes)
				if err == nil {
					certs = append(certs, cert)
				}
			}
			block, rest = pem.Decode(rest)
		}
	} else {
		// Try DER decoding
		cert, err := x509.ParseCertificate(data)
		if err == nil {
			certs = append(certs, cert)
		}
	}

	if len(certs) == 0 {
		// If we couldn't parse it as a cert, fallback to text
		th := &TextHandler{}
		return th.Render(w, strings.NewReader(string(data)), meta, opts)
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	if opts.NoColor || opts.Plain {
		keyStyle = lipgloss.NewStyle()
		valStyle = lipgloss.NewStyle()
	}

	for i, cert := range certs {
		if i > 0 {
			fmt.Fprintln(w, strings.Repeat("-", 40))
		}
		fmt.Fprintf(w, "%s %s\n", keyStyle.Render("Certificate:"), valStyle.Render(""))
		fmt.Fprintf(w, "  %s\n", keyStyle.Render("Data:"))
		fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("Version:"), valStyle.Render(fmt.Sprintf("%d", cert.Version)))
		fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("Serial Number:"), valStyle.Render(fmt.Sprintf("%x", cert.SerialNumber)))
		fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("Signature Algorithm:"), valStyle.Render(cert.SignatureAlgorithm.String()))
		fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("Issuer:"), valStyle.Render(cert.Issuer.String()))
		fmt.Fprintf(w, "    %s\n", keyStyle.Render("Validity:"))
		fmt.Fprintf(w, "      %s %s\n", keyStyle.Render("Not Before:"), valStyle.Render(cert.NotBefore.String()))
		fmt.Fprintf(w, "      %s %s\n", keyStyle.Render("Not After :"), valStyle.Render(cert.NotAfter.String()))
		fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("Subject:"), valStyle.Render(cert.Subject.String()))
		fmt.Fprintf(w, "    %s\n", keyStyle.Render("Subject Public Key Info:"))
		fmt.Fprintf(w, "      %s %s\n", keyStyle.Render("Public Key Algorithm:"), valStyle.Render(cert.PublicKeyAlgorithm.String()))
		
		if len(cert.DNSNames) > 0 {
			fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("DNS Names:"), valStyle.Render(strings.Join(cert.DNSNames, ", ")))
		}
		if len(cert.IPAddresses) > 0 {
			var ips []string
			for _, ip := range cert.IPAddresses {
				ips = append(ips, ip.String())
			}
			fmt.Fprintf(w, "    %s %s\n", keyStyle.Render("IP Addresses:"), valStyle.Render(strings.Join(ips, ", ")))
		}
	}

	return nil
}
