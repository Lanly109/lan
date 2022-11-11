/*
Copyright © 2022 Lanly

*/
package cmd

import (
	"github.com/Lanly109/lan/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var htmlName string

const countDownHtml string = `<!doctype html>
<html>

<head>

	<script type="text/javascript">
		var StartTime = new Date('2022/10/29 08:27:00') // please change it
		var EndTime = new Date('2022/10/29 11:57:00') // please change it
		var password = "password"   // please change it and press F5

		var totleTime = EndTime.getTime() - StartTime.getTime()
	</script>

	<meta charset="utf-8">
	<style type="text/css">
		#CountMsg {
			font-size: 120px;
			text-align: center;
			color: yellow;
			line-height: 1;
		}

		#Title {
			font-size: 150px;
			text-align: center;
			color: yellow;
			line-height: 2;
		}

		#password {
			font-size: 60px;
			text-align: center;
			color: yellow;
			line-height: 2;
			font-family: Consolas, Tahoma;
		}

		#CountMsg {
			font-family: Helvetica, Tahoma, Arial;
		}

		body {
			font-family: "黑体";
			background: black;
		}
	</style>
</head>

<body>
	<div id="password">解压密码：<span id="hiden" ></span></div>
	<div id="Title">
		<span id="info">距离比赛开始还有</span>
		<div id="CountMsg" class="HotData">
			<span id="t_h">00</span>
			<span>:</span>
			<span id="t_m">00</span>
			<span>:</span>
			<span id="t_s">00</span>
			<span>:</span>
			<span id="t_ms">00</span>
		</div>
		<div>
			<meter id="meter" min="0" max=1000 value="100" high="9" low="4" optimum="1">
		</div>
</body>

<script type="text/javascript">
	document.getElementById("meter").max = totleTime;
	function getRTime() {
		var NowTime = new Date();
		var t = EndTime.getTime() - NowTime.getTime();
		var ss = NowTime.getTime() - StartTime.getTime();
		if (ss < 0){
			document.getElementById("info").innerText = "距离比赛开始还有";	
			t = -ss
		}else if (t < 0){
			document.getElementById("info").innerText = "比赛结束!";	
			t = 0
		}else {
			document.getElementById("hiden").innerText = password
			document.getElementById("info").innerText = "距离比赛结束还有";	
		}

		var h = Math.floor(t / 1000 / 60 / 60 % 24);
		if (h < 10) h = "0" + h;
		var m = Math.floor(t / 1000 / 60 % 60);
		if (m < 10) m = "0" + m;
		var s = Math.floor(t / 1000 % 60);
		if (s < 10) s = "0" + s;
		var ms = Math.floor(t % 1000);
		if (ms < 10) ms = "0" + "0" + ms;
		else if (ms < 100) ms = "0" + ms;

		var gone = totleTime - t;

		if (ss < 0)
			gone = 0

		document.getElementById("t_h").innerText = h;
		document.getElementById("t_m").innerText = m;
		document.getElementById("t_s").innerText = s;
		document.getElementById("t_ms").innerText = ms;
		document.getElementById("meter").value = gone;
	}

	setInterval(getRTime, 0);
</script>

</html>
`

// countDownCmd represents the countdown command
var countDownCmd = &cobra.Command{
	Use:   "countdown",
	Short: "Generate html that is the countdown of a contest",
	Long:  `A html source code will be generated that is the countdown of a contest.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("filename: %s", htmlName)
		if utils.FileExist(htmlName) {
			log.Errorf("The File %s Exists!", htmlName)
			return
		}
		err := utils.WriteFile(htmlName, countDownHtml)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Successfully generated sharing script [%s]", htmlName)
		}
	},
}

func init() {
	genCmd.AddCommand(countDownCmd)

	countDownCmd.Flags().StringVarP(&htmlName, "name", "n", "countdown.html", "The countdown html name(default countdown.html)")
}
