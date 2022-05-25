package cmd

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/zu1k/nali/pkg/entity"
)

//查询到的IP信息
type Ipinfo struct {
	Code int    `json:"code" xml:"code"`
	Ip   string `json:"ip" xml:"ip"`
	Addr string `json:"addr" xml:"addr"`
}

//错误状态信息
type Errinfo struct {
	Code int    `json:"code" xml:"code"`
	Msg  string `json:"msg" xml:"msg"`
}

// Web server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start web api server",
	Long:    `Start web api server`,
	Example: "nali server --port 8080",
	Run: func(cmd *cobra.Command, args []string) {
		/*

			ip := res[0].Text
			addr := strings.Split(strings.Replace(res[0].Info,"\t"," ",-1)," ")
			country := addr[0]
			area := addr[1]
			fmt.Println(ip,country,area)
		*/
		port, _ := cmd.Flags().GetString("port")
		port = ":" + port

		//启动echo web服务
		ec := echo.New()
		//路径方式查询
		ec.GET("/:ip", func(c echo.Context) error {
			ip := c.Param("ip")
			if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Help: ?ip=223.5.5.5 or /223.5.5.5 (Default: your ip)",
				}
				return c.JSON(http.StatusOK, reinfo)
			} else if ip == "" {
				ip = c.RealIP()
			}
			//fmt.Println(ip)
			args = nil
			args = append(args, ip)
			res := entity.ParseLine(strings.Join(args, " "))
			//addr := strings.Split(strings.Replace(res[0].Info, "\t", " ", -1), " ")
			//判断是否有结果返回
			if res[0].Info == "" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "No record, please try again",
				}
				return c.JSON(http.StatusOK, reinfo)
			} else {
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
				}
				return c.JSON(http.StatusOK, reinfo)
			}

		})
		//参数方式查询
		ec.GET("/", func(c echo.Context) error {
			ip := c.QueryParam("ip")
			if ip == "help" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "Param: ip=223.5.5.5 or /223.5.5.5 (Default: visitor ip)",
				}
				return c.JSON(http.StatusOK, reinfo)
			} else if ip == "" {
				ip = c.RealIP()
			}
			//fmt.Println(ip)
			args = nil
			args = append(args, ip)
			res := entity.ParseLine(strings.Join(args, " "))
			//判断是否有结果返回
			if res[0].Info == "" {
				reinfo := &Errinfo{
					Code: -1,
					Msg:  "No record, please try again",
				}
				return c.JSON(http.StatusOK, reinfo)
			} else {
				reinfo := &Ipinfo{
					Code: 1,
					Ip:   res[0].Text,
					Addr: strings.Replace(res[0].Info, "\t", " ", -1),
				}
				return c.JSON(http.StatusOK, reinfo)
			}

		})
		ec.Logger.Fatal(ec.Start(port))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().String("port", "8080", "Set web service listen port")
}
