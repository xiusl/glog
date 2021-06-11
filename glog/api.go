package glog

var defaultLog = NewConsoleLogger("debug")

func Debug(format string, arg ...interface{}) {
	defaultLog.Debug(format, arg...)
}
func Info(format string, arg ...interface{}) {
	defaultLog.Info(format, arg...)

}
func Trace(format string, arg ...interface{}) {
	defaultLog.Trace(format, arg...)

}
func Warning(format string, arg ...interface{}) {
	defaultLog.Warning(format, arg...)

}
func Error(format string, arg ...interface{}) {
	defaultLog.Error(format, arg...)

}
func Fatal(format string, arg ...interface{}) {
	defaultLog.Fatal(format, arg...)

}
