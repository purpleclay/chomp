/*
Copyright (c) 2023 - 2024 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/purpleclay/chomp"
)

type GpgPrivateKey struct {
	UserName     string
	UserEmail    string
	SecretKey    GpgKeyDetails
	SecretSubKey GpgKeyDetails
}

type GpgKeyDetails struct {
	CreationDate   int
	ExpirationDate int
	Fingerprint    string
	KeyID          string
	KeyGrip        string
}

func (k GpgPrivateKey) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("fingerprint:     %s\n", k.SecretKey.Fingerprint))
	buf.WriteString(fmt.Sprintf("keygrip:         %s\n", k.SecretKey.KeyGrip))
	buf.WriteString(fmt.Sprintf("key_id:          %s\n", k.SecretKey.KeyID))
	buf.WriteString(fmt.Sprintf("created_on:      %d\n", k.SecretKey.CreationDate))
	if k.SecretKey.ExpirationDate > 0 {
		buf.WriteString(fmt.Sprintf("expires_on:      %d\n", k.SecretKey.ExpirationDate))
	}

	buf.WriteString(fmt.Sprintf("user:            %s <%s>\n", k.UserName, k.UserEmail))

	buf.WriteString(fmt.Sprintf("sub_fingerprint: %s\n", k.SecretSubKey.Fingerprint))
	buf.WriteString(fmt.Sprintf("sub_keygrip:     %s\n", k.SecretSubKey.KeyGrip))
	buf.WriteString(fmt.Sprintf("sub_key_id:      %s\n", k.SecretSubKey.KeyID))
	buf.WriteString(fmt.Sprintf("created_on:      %d\n", k.SecretSubKey.CreationDate))
	if k.SecretSubKey.ExpirationDate > 0 {
		buf.WriteString(fmt.Sprintf("expires_on:      %d\n", k.SecretSubKey.ExpirationDate))
	}

	return buf.String()
}

func Parse(str string) (GpgPrivateKey, error) {
	key := GpgPrivateKey{}

	rem, sec, err := secretKey()(str)
	if err != nil {
		return key, err
	}
	crtDate, _ := strconv.Atoi(sec[1])
	expDate, _ := strconv.Atoi(sec[2])

	key.SecretKey = GpgKeyDetails{
		CreationDate:   crtDate,
		ExpirationDate: expDate,
		KeyID:          sec[0],
	}

	rem, fpr, err := fingerprint()(rem)
	if err != nil {
		return key, err
	}
	key.SecretKey.Fingerprint = fpr

	rem, grp, err := keygrip()(rem)
	if err != nil {
		return key, err
	}
	key.SecretKey.KeyGrip = grp

	rem, usr, err := user()(rem)
	if err != nil {
		return key, err
	}
	key.UserName = usr[0]
	key.UserEmail = usr[1]

	rem, ssb, err := secretSubKey()(rem)
	if err != nil {
		return key, err
	}
	crtDate, _ = strconv.Atoi(ssb[1])
	expDate, _ = strconv.Atoi(ssb[2])

	key.SecretSubKey = GpgKeyDetails{
		CreationDate:   crtDate,
		ExpirationDate: expDate,
		KeyID:          ssb[0],
	}

	rem, fpr, err = fingerprint()(rem)
	if err != nil {
		return key, err
	}
	key.SecretSubKey.Fingerprint = fpr

	rem, grp, err = keygrip()(rem)
	if err != nil {
		return key, err
	}
	key.SecretSubKey.KeyGrip = grp

	return key, nil
}

func secretKey() chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		// sec:-:4096:1:AAC7E54CBD73F690:1664450926:::-:::scESC:::+:::23::0:
		var rem string
		var err error

		if rem, _, err = chomp.Pair(
			chomp.Tag("sec"),
			chomp.Repeat(colon(), 4))(s); err != nil {
			return rem, nil, err
		}

		var ext []string
		if rem, ext, err = chomp.Repeat(colon(), 3)(rem); err != nil {
			return rem, nil, err
		}

		if rem, _, err = eol()(rem); err != nil {
			return rem, nil, err
		}

		return rem, ext, nil
	}
}

func colon() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		// <any>:
		rem, ext, err := chomp.Pair(chomp.Until(":"), chomp.Tag(":"))(s)
		if err != nil {
			return rem, "", err
		}

		return rem, ext[0], nil
	}
}

func eol() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		rem, _, err := chomp.Pair(chomp.Until("\n"), chomp.Crlf())(s)
		if err != nil {
			return rem, "", err
		}

		return rem, "", nil
	}
}

func secretSubKey() chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		// ssb:-:4096:1:17441D4227A0B812:1664450926::::::e:::+:::23:
		var rem string
		var err error

		if rem, _, err = chomp.Pair(
			chomp.Tag("ssb"),
			chomp.Repeat(colon(), 4))(s); err != nil {
			return rem, nil, err
		}

		var ext []string
		if rem, ext, err = chomp.Repeat(colon(), 3)(rem); err != nil {
			return rem, nil, err
		}

		if rem, _, err = eol()(rem); err != nil {
			return rem, nil, err
		}

		return rem, ext, nil
	}
}

func fingerprint() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		// fpr:::::::::28BF65E18407FD2966565284AAC7E54CBD73F690:
		var rem string
		var err error

		if rem, _, err = chomp.Pair(chomp.Tag("fpr"), chomp.Repeat(chomp.Tag(":"), 9))(s); err != nil {
			return rem, "", err
		}

		var fpr string
		if rem, fpr, err = chomp.Until(":")(rem); err != nil {
			return rem, "", err
		}

		if rem, _, err = eol()(rem); err != nil {
			return rem, "", err
		}
		return rem, fpr, nil
	}
}

func keygrip() chomp.Combinator[string] {
	return func(s string) (string, string, error) {
		// grp:::::::::12E86CE47CEB942D2A65B4D02106657BA8D0C92B:
		var rem string
		var err error

		if rem, _, err = chomp.Pair(chomp.Tag("grp"), chomp.Repeat(chomp.Tag(":"), 9))(s); err != nil {
			return rem, "", err
		}

		var grp string
		if rem, grp, err = chomp.Until(":")(rem); err != nil {
			return rem, "", err
		}

		if rem, _, err = eol()(rem); err != nil {
			return rem, "", err
		}
		return rem, grp, nil
	}
}

func user() chomp.Combinator[[]string] {
	return func(s string) (string, []string, error) {
		// uid:-::::1664450926::E6F81442C4BEE48D9ED3E6EE4CAC21231D3C25EB::john.smith <john.smith@testing.com>::::::::::0:
		var rem string
		var err error

		if rem, _, err = chomp.Pair(
			chomp.Tag("uid"),
			chomp.Repeat(colon(), 9))(s); err != nil {
			return rem, nil, err
		}

		var ext []string
		if rem, ext, err = chomp.SepPair(
			chomp.Until(" "),
			chomp.Tag(" "),
			chomp.BracketAngled())(rem); err != nil {
			return rem, nil, err
		}

		if rem, _, err = eol()(rem); err != nil {
			return rem, nil, err
		}
		return rem, ext, nil
	}
}

func main() {
	colonFmt := `sec:-:4096:1:AAC7E54CBD73F690:1664450926:::-:::scESC:::+:::23::0:
fpr:::::::::28BF65E18407FD2966565284AAC7E54CBD73F690:
grp:::::::::12E86CE47CEB942D2A65B4D02106657BA8D0C92B:
uid:-::::1664450926::E6F81442C4BEE48D9ED3E6EE4CAC21231D3C25EB::albert.einstein <albert.einstein@emcsqua.red>::::::::::0:
ssb:-:4096:1:17441D4227A0B812:1664450926::::::e:::+:::23:
fpr:::::::::26965E00791A52ECC33AE88917441D4227A0B812:
grp:::::::::603DAFFC5AAE42C4B8BFCC99DD7CEDD5C443FFA0:
`

	pk, err := Parse(colonFmt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pk)
}
