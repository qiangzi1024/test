package awesome

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func doNodeCommand() string {
	cmd := exec.Command("node", "JD_DailyBonus.js")
	stdout, err := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err != nil {
		log.Panicln(err)
	}

	if err := cmd.Start(); err != nil {
		log.Panicln("Do Command Error:", err)
	}
	stdoutBody, err := ioutil.ReadAll(stdout)
	stderrBody, _ := ioutil.ReadAll(stderr)
	if err != nil {
		log.Panicln(err)
	}
	resultStr := string(stdoutBody)
	if len(stderrBody) != 0 {
		ioutil.WriteFile(jdResultFile, stderrBody, 0644)
	} else {
		ioutil.WriteFile(jdResultFile, stdoutBody, 0644)
	}
	log.Println("Result:", resultStr)
	log.Println("执行完毕~")
	return resultStr
}

var JD_COOKIE string
var JD_COOKIE_2 string
var PUSH_KEY string
var jdJSFile string
var jdResultFile string

func JustDoIt() {
	jdJSFile = "JD_DailyBonus.js"
	jdResultFile = "result.md"
	JD_COOKIE = os.Getenv("JD_COOKIE")
	PUSH_KEY = os.Getenv("PUSH_KEY")
	JD_COOKIE_2 = os.Getenv("JD_COOKIE_2")
	if len(JD_COOKIE) == 0 {
		log.Fatalln("请先填写KEY!!!")
	}

	downloadJDJSFile()
	replaceWithKEY()
	result := doNodeCommand()
	if len(PUSH_KEY) != 0 {
		sendNotify(fmt.Sprintf("京东签到结果_%s", time.Now().Format("2006-01-02 15:04:05")), result)
	} else {
		log.Println("未设置PUSH_KEY")
	}
}

func downloadJDJSFile() {
	// url := "https://cdn.jsdelivr.net/gh/NobyDa/Script@master/JD-DailyBonus/JD_DailyBonus.js"
	url := "https://raw.githubusercontent.com/NobyDa/Script/master/JD-DailyBonus/JD_DailyBonus.js"
	if err := downloadFile(jdJSFile, url); err != nil {
		log.Fatalln(err)
	}
	log.Println("JD_DailyBonus.js下载完毕~")
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func replaceWithKEY() {
	// let content = await fs.readFileSync('./JD_DailyBonus.js', 'utf8')
	// content = content.replace(/var Key = ''/, `var Key = '${KEY}'`);
	// if (DualKey) {
	//  content = content.replace(/var DualKey = ''/, `var DualKey = '${DualKey}'`);
	// }
	// await fs.writeFileSync( './JD_DailyBonus.js', content, 'utf8')
	bytes, _ := ioutil.ReadFile(jdJSFile)
	jsContent := string(bytes)
	jsContent = strings.ReplaceAll(jsContent, "var Key = ''", fmt.Sprintf("var Key = '%s'", JD_COOKIE))
	if len(JD_COOKIE_2) != 0 {
		// 替换第二个账号
		jsContent = strings.ReplaceAll(jsContent, "var DualKey = ''", fmt.Sprintf("var DualKey = '%s'", JD_COOKIE))
	}
	ioutil.WriteFile(jdJSFile, []byte(jsContent), 0644)
	log.Println("替换变量完毕~")
}

func sendNotify(text string, desp string) {
	resp, err := http.PostForm(fmt.Sprintf("https://sctapi.ftqq.com/%s.send", PUSH_KEY), url.Values{"title": {text}, "desp": {desp}})
	if err != nil {
		log.Panicln(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
