package logger

import (
	"fmt"
	"os"
	"projects_template/internal/entity/global"
	"projects_template/internal/entity/rpc"
	"time"

	"io"
	"log"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UseGinLogger is the logrus logger handler
func UseGinLogger(log *logrus.Logger, errlog *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		// latency := start.Sub(stop)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		// clientUserAgent := c.Request.UserAgent()
		// referer := c.Request.Referer()
		// hostname, err := os.Hostname()
		// if err != nil {
		// hostname = "unknow"
		// }
		// dataLength := c.Writer.Size()
		// if dataLength < 0 {
		// 	dataLength = 0
		// }

		operLogin := c.GetString(global.UserLoginKey)
		token := c.GetString(global.TokenKey)
		methodName := c.GetString(rpc.MethodNameKey)
		// gatewayVersion := c.GetInt(rpc.GatewayVersion)

		logFields := logrus.Fields{
			"operLogin": operLogin,
			"token":     token,
		}

		if methodName != "" {
			logFields["m"] = methodName
			path += "/" + methodName
		}

		entry := logrus.NewEntry(log).WithFields(logFields)

		msg := fmt.Sprintf("| %s | %s | %s | %s", stop, clientIP, c.Request.Method, path)

		rpcError, errorExists := c.Get(rpc.MethodErrorKey)
		rpcErrorMap, errorMapExists := c.Get(rpc.MultiplyMethodErrorKey)
		if errorExists {
			rerr := rpcError.(*rpc.Error)
			msg = fmt.Sprintf("| %s %s, error=%s", bgRed(" ", rerr.Code, " "), msg, rerr.Message)
			entry.Error(msg)
		} else if errorMapExists {
			rerr := rpcErrorMap.(rpc.MethodErrors)
			methods := ""
			methodsErr := ""
			for name, m := range rerr {
				methods += bgRed(" ", fmt.Sprintf("%s=%d", name, m.Code), " ")
				methodsErr += bgRed(" ", fmt.Sprintf("%s=%s", name, m.Message), " ")
			}
			msg = fmt.Sprintf("| %s %s, error=%s", msg, methods, methodsErr)
			entry.Warn(msg)
		} else if statusCode > 499 {
			msg = fmt.Sprintf("| %s %s", bgRed(" ", statusCode, " "), msg)
			entry.Error(msg)
		} else if statusCode > 399 {
			msg = fmt.Sprintf("| %s %s", bgYellow(" ", statusCode, " "), msg)
			entry.Warn(msg)
		} else {
			msg = fmt.Sprintf("| %s %s", bgGreen(" ", statusCode, " "), msg)
			entry.Info(msg)
		}

		if len(c.Errors) > 0 {
			errlog.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}

var (
	bgRed    = color.New(color.BgRed, color.FgWhite).SprintFunc()
	bgYellow = color.New(color.BgYellow, color.FgWhite).SprintFunc()
	bgGreen  = color.New(color.BgGreen, color.FgWhite).SprintFunc()
)

func GetAccessLogWriter(module string) io.Writer {
	defaultLogDir := ""

	logDir := defaultLogDir
	if os.Getenv("LOG_DIR") != "" {
		logDir = os.Getenv("LOG_DIR")
	}

	accessLogFile, err := os.OpenFile(
		fmt.Sprintf("%s/%s_access.log", logDir, module),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatalln("Нет возможности писать в файл-логи. Убедитесь, что существует папка c правом доступа", logDir)
	}

	iw := io.Writer(accessLogFile)
	return iw

}
