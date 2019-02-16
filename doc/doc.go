package doc

import (
	"fmt"
	"errors"
	"strings"
)


var SEPARATORS = []string{":", ";", ",", "\t"}


type Document struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Domain string `json:"domain"`
	DomainNoTLD string `json:"domain_notld"`
	Tld string `json:"tld"`
	Password string `json:"password"`
	PasswordLen int `json:"password_length"`
	Source int `json:"source"`
}


type MetaDoc struct {
	Id string
	Type string
	Prefix string
	Routing string
	Doc *Document
}


func ParseLine(line string, type_ string, source int, partitionLevels int) (MetaDoc, error) {
	line = strings.TrimSpace(line)
	if len(line) > 255 {
		return MetaDoc{}, errors.New("line is too long")
	}
	var sline []string
	for _, sep := range SEPARATORS {
		sline = strings.SplitN(line, sep, 2)
		if len(sline) == 2 {
			break
		}
	}
	if len(sline) != 2 {
		return MetaDoc{}, errors.New("no separator")
	}
	email := sline[0]
	username := email
	domain := ""
	domain_no_tld := ""
	tld := ""
	semail := strings.SplitN(email, "@", 2)
	if len(semail) == 2 {
		username = semail[0]
		domain = semail[1]
		dotIdx := strings.LastIndex(domain, ".")
		if (dotIdx != -1) {
			domain_no_tld = domain[:dotIdx]
			tld = domain[dotIdx+1:]
		}
	}
	password := sline[1]
	password_len := len(password)
	doc := Document{
		Email: email,
		Username: username,
		Domain: domain,
		DomainNoTLD: domain_no_tld,
		Tld: tld,
		Password: password,
		PasswordLen: password_len,
		Source: source,
	}
	var prefix string
	id := fmt.Sprintf("%s%s%d", email, password, source)
	if len(email) >= partitionLevels {
		prefix = string(email[0:partitionLevels])
	}
	return MetaDoc{
		Doc: &doc,
		Id: id,
		Type: type_,
		Routing: strings.ToLower(doc.Email),
		Prefix: prefix,
	}, nil
}
