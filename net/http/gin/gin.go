package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"io/ioutil"
	"time"
)

type Config struct {
	// 日志文件名，留空不写
	Filename string `json:"filename" toml:"filename"`

	// 日志达到MaxSize(MB)后，进行日志翻滚
	MaxSize int `json:"maxsize" toml:"maxsize"`

	// 多少天翻滚
	MaxAge int `json:"maxage" toml:"maxage"`

	// 日志备份数，超过将删除
	MaxBackups int `json:"maxbackups" toml:"maxbackups"`

	// 是否使用本地时间
	LocalTime bool `json:"localtime" toml:"localtime"`

	// 是否压缩日志存贮
	Compress bool `json:"compress" toml:"compress"`
}

func New(c *Config) (r *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r = gin.New()

	// 注册验证器
	//RegValidation()
	// 初始化日志
	lum := &lumberjack.Logger{
		Filename:   "", //默认是空，不写文件
		MaxSize:    30, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		LocalTime:  false,
		Compress:   false, // disabled by default
	}

	if c.Filename != "" {
		lum.Filename = c.Filename
	}
	if c.MaxSize > 0 {
		lum.MaxSize = c.MaxSize
	}
	if c.MaxBackups > 0 {
		lum.MaxBackups = c.MaxBackups
	}
	if c.MaxAge > 0 {
		lum.MaxAge = c.MaxAge
	}
	if c.LocalTime {
		lum.LocalTime = c.LocalTime
	}
	if c.Compress {
		lum.Compress = c.Compress
	}
	cf := gin.LoggerConfig{
		Output:    lum,
		SkipPaths: []string{"/favicon.ico"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s\" %d %d %s \"%s\" \"%s\" \"%s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.BodySize,
				params.Latency,
				params.Request.UserAgent(),
				params.Request.Referer(),
				params.ErrorMessage,
			)
		},
	}
	if c.Filename != "" {
		r.Use(gin.LoggerWithConfig(cf))
	}

	r.Use(gin.Recovery())
	//r.Use(cors.CORSMiddleware())
	//r.Use(action_log.ActionLog())
	//r.Use(jwt.JWTAuth())
	return
}
