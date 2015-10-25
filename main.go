package main

func main() {
	oraDB := oracleDB{}
	errPrtg := errorXML{}

	ses := oraDB.getOraSession()
	if nil != ses {
		refCursor := oraDB.executeOraSP(ses)
		if refCursor.IsOpen() {
			prtg := prtgXML{}
			if prtg.outputXMLResult(refCursor) < 1 {
				errPrtg.outputXMLError("Fetched 0 result!")
			}
		} else {
			errPrtg.outputXMLError("Cursor not open: " + oraDB.mErr.Error())
		}
	} else {
		errPrtg.outputXMLError("Unable to get Oracle session: " + oraDB.mErr.Error())
	}

	defer oraDB.closeOraSession()
}
