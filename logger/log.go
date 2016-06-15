package logger

import (
	. "github.com/Dataman-Cloud/omega-billing/config"
	log "github.com/cihub/seelog"
)

func LogInit() {
	logger, err := log.LoggerFromConfigAsString(logConfig())
	if err == nil {
		log.ReplaceLogger(logger)
	} else {
		log.Error(err)
	}
}

func logConfig() string {
	logconfig := `<seelog type="asynctimer" asyncinterval="5000000" minlevel="debug">
			            <outputs formatid="main">`
	if GetConfig().Log.Console {
		logconfig += `<console/>`
	}
	if GetConfig().Log.AppendFile {
		logconfig += `<buffered size="10000" flushperiod="1000">`
		logconfig += `<rollingfile type="size" filename="` + GetConfig().Log.File + `" maxsize="5000000" maxrolls="30" />`
		logconfig += `</buffered>`
	}
	logconfig += `</outputs>
	      <formats>
	   <format id="main" format="%Date(2006-01-02 15:04:05Z07:00) [%LEVEL] %Msg%n"/>
		</formats>
		</seelog>`
	return logconfig
}
