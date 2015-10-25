package main

import (
	"errors"
	"flag"
	"strings"

	"gopkg.in/rana/ora.v3"
)

type oracleDB struct {
	//mStmt *ora.Stmt
	mSes *ora.Ses
	mSrv *ora.Srv
	mEnv *ora.Env
	mErr error
}

var (
	missingArgs string = ""
	user               = flag.String("user", "", "username")
	pass               = flag.String("pass", "", "password")
	host               = flag.String("host", "localhost", "server hostname")
	port               = flag.String("port", "1521", "port")
	sid                = flag.String("sid", "", "sid")
	sn                 = flag.String("sn", "", "service name")
	sp                 = flag.String("sp", "", "stored procedure name")
	spPar              = flag.String("sp_par", "", "stored procedure parameter")
)

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	//ora.Register(nil)
}

func (oracleDB *oracleDB) closeOraSession() {
	if nil != oracleDB.mSes && oracleDB.mSes.IsOpen() {
		oracleDB.mSes.Close()
	}
	if nil != oracleDB.mSrv && oracleDB.mSrv.IsOpen() {
		oracleDB.mSrv.Close()
	}
	if nil != oracleDB.mEnv && oracleDB.mEnv.IsOpen() {
		oracleDB.mEnv.Close()
	}
}

//returns the new session
func (oracleDB *oracleDB) getOraSession() *ora.Ses {
	var sidOrSn string

	oracleDB.mEnv, oracleDB.mErr = ora.OpenEnv(nil)

	if nil == oracleDB.mErr {

		srvCfg := ora.NewSrvCfg()
		if *sid != "" {
			sidOrSn = "SID=" + *sid
		} else {
			sidOrSn = "SERVICE_NAME=" + *sn
		}

		srvCfg.Dblink = "(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=" +
			*host + ")(PORT=" + *port + "))(CONNECT_DATA=(" + sidOrSn + ")))"

		oracleDB.mSrv, oracleDB.mErr = oracleDB.mEnv.OpenSrv(srvCfg)
		//defer srv.Close()
		if nil == oracleDB.mErr {
			sesCfg := ora.NewSesCfg()
			sesCfg.Username = *user
			sesCfg.Password = *pass
			oracleDB.mSes, oracleDB.mErr = oracleDB.mSrv.OpenSes(sesCfg)

		}
	}

	return oracleDB.mSes
}

//if ses is nill it will use the struc session mSes
//returns a &ora.Rset{}
func (oracleDB *oracleDB) executeOraSP(ses *ora.Ses) *ora.Rset {
	if nil == ses {
		ses = oracleDB.mSes
	}

	refCursor := &ora.Rset{}
	stmt, err := ses.Prep("CALL " + *sp + "(:1, :STATNAME, :ACCESSGROUP, :INTMIN)")

	if nil == err {
		//split param
		var params = strings.Fields(*spPar)

		if len(params) == 3 {
			_, err := stmt.Exe(refCursor, params[0], params[1], params[2])
			if nil != err {
				oracleDB.mErr = err
			}
		} else {
			oracleDB.mErr = errors.New("Incomplete SP parameter!")
		}
	} else {
		oracleDB.mErr = err
	}

	return refCursor
}

func (oracleDB *oracleDB) buildDSNFromArgs() string {
	var lDSN string = ""

	switch {
	case *sp == "":
		missingArgs = "sp "
		fallthrough
	case *user == "":
		missingArgs += "user "
		fallthrough
	case *pass == "":
		missingArgs += "pass "
		fallthrough
	case *host == "":
		missingArgs += "host "
		fallthrough
	case *port == "":
		missingArgs += "port "
		fallthrough
	case (*sid == "") && (*sn == ""):
		missingArgs += "sid sn"
	}

	if !(missingArgs == "") {
		missingArgs = strings.Replace(strings.TrimSpace(missingArgs), " ", ", ", -1)
		//prtgXML.outputXMLError(true, "Missing parameter(s): " + missingArgs, nil)

	} /*else {
		if !(*sid == "") {
			lDSN = *user + "/" + *pass + "@" + *host + ":" + *port + "/" + *sid
		} else if !(*sn == "") {
			lDSN = *user + "/" + *pass + "@" + *host + ":" + *port + "/" + *sn
		}
	}*/

	return lDSN

}
