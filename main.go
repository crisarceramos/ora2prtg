package main

import ()

func main() {
	oraDB := oracleDB{}
	prtg := prtgXML{}
	
	ses := oraDB.getOraSession()
	if ses != nil {
		refCursor := oraDB.executeOraSP(ses)
		if refCursor.IsOpen() {
			outputXMLResult(refCursor)
		} else {
			prtg.outputXMLError(true, "No result!", nil)
		}
	} else {
		prtg.outputXMLError(true, "Unable to get Oracle session!", nil)
	}

	defer oraDB.closeOraSession()
}
