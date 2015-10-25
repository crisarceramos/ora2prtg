package main

import (
	"encoding/xml"
	"os"

	"gopkg.in/rana/ora.v3"
)

type prtgXML struct {
	XMLName xml.Name `xml:"prtg"`
	Res     []result `xml:"result"`
}

type result struct {
	XMLName  xml.Name `xml:"result"`
	EChannel string   `xml:"channel"`
	EValue   int64    `xml:"value"`
	EFloat   int      `xml:"float"`
	EUnit    string   `xml:"unit"`
}

type errorXML struct {
	XMLName xml.Name `xml:"prtg"`
	IsError int      `xml:"error"`
	Msg     string   `xml:"text"`
}

func (errorXml *errorXML) outputXMLError(msg string) {
	errPrtg := &errorXML{IsError: 1, Msg: msg}

	output, err := xml.MarshalIndent(errPrtg, "  ", "    ")
	if nil == err {
		os.Stdout.Write(output)
	}
}

func (prtgXml *prtgXML) outputXMLResult(refCursor *ora.Rset) int {
	var ctr int = 0
	prtg := &prtgXML{}

	for refCursor.Next() {
		res := &result{}

		res.EChannel, _ = refCursor.Row[0].(string)
		res.EValue, _ = refCursor.Row[1].(int64)
		res.EFloat, _ = refCursor.Row[2].(int)
		res.EUnit, _ = refCursor.Row[3].(string)

		prtg.Res = append(prtg.Res, *res)
		ctr += 1
	}

	if ctr > 0 {
		output, err := xml.MarshalIndent(prtg, "  ", "    ")
		if nil == err {
			os.Stdout.Write(output)
		} else {
			errPrtg := errorXML{}
			errPrtg.outputXMLError("Unable to write xml result!")
		}
	}

	return ctr
}
