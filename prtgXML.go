package main

import(
	"fmt"
	"os"
	"encoding/xml"
	"gopkg.in/rana/ora.v3"
)

type prtgXML struct {
    XMLName xml.Name `xml:"prtg"`
    Res     []result `xml:"result"`
}

type result struct {
	XMLName xml.Name 	`xml:"result"`
    EChannel string 		`xml:"channel"`
    EValue   int 	    `xml:"value"`
	EFloat	int	    	`xml:"float"`
	EUnit	string 		`xml:"unit"`	
}

func outputXMLResult(refCursor *ora.Rset) {

	prtg := &prtgXML{}
	
	for refCursor.Next() {
		res := &result{}
		
		res.EChannel, _ = refCursor.Row[0].(string)
		res.EValue, _ = refCursor.Row[1].(int)
		res.EFloat, _ = refCursor.Row[2].(int)
		res.EUnit, _ = refCursor.Row[3].(string)

		prtg.Res = append(prtg.Res, *res)
	}

	output, err := xml.MarshalIndent(prtg, "  ", "    ")
    if err == nil {
        os.Stdout.Write(output)	
    }else{
		prtg.outputXMLError(true, "Building xml result!", nil)
	}

}

func (prtgXML *prtgXML) outputXMLError(isError bool, errMSG string, refCursor *ora.Rset) {
	if isError {
		fmt.Println(errMSG)
	}else{
		//enc := xml.NewEncoder(os.Stdout)
		//enc.Indent("  ", "    ")
		//if err := enc.Encode(refCursor); err != nil {
		//    fmt.Printf("error: %v\n", err)
		//}
		
	}
}
