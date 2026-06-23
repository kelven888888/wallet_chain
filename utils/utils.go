package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-gomail/gomail"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"wallet_chain.com/global"
	amtime "wallet_chain.com/utils/time"
)

const (
	TrafficKey = "X-Request-Id"
	LoggerKey  = "_go-admin-logger-request"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}
func IsElementInSlice(slice []int, target int) bool {
	sort.Ints(slice)
	index := sort.SearchInts(slice, target)
	return index < len(slice) && slice[index] == target
}

// 对密码进行加密
func EncryptPassworld(password string) string {
	if len(password) == 0 {
		return ""
	}

	pass, err := NewPass(password, "087DC8269428CD3BC4BB689C5EE9A148")
	if err != nil {
		global.SHOP_LOG.Error("密码加密错误")
		return ""
	}

	return pass
}

// 通过scrypt生成密码
func NewPass(passwd, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(dk), nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(s string) string {

	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// 返回参数的类型
func Type(v interface{}) string {
	t := reflect.TypeOf(v)
	k := t.Kind()
	return k.String()
}

func InArray(in interface{}, list interface{}) bool {
	ret := false
	if in == nil {
		in = ""
	}

	// 判断list是否slice
	l := reflect.TypeOf(list).String()
	t := Type(in)
	if "[]"+t != l {
		return false
	}

	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		for _, l := range list.([]string) {
			v := reflect.ValueOf(l)
			if tv == v.String() {
				ret = true
				break
			}
		}

	case "int":
		tv := reflect.ValueOf(in).Int()
		for _, l := range list.([]int) {
			v := reflect.ValueOf(l)
			if tv == v.Int() {
				ret = true
				break
			}
		}
	}

	return ret
}
func Us_datatime() string {
	var localTime amtime.LocalTime

	return localTime.FormatDateString(localTime.String())

}

func Us_datatimecon() time.Time {

	times := amtime.LocalTime{}
	return times.Local()
}
func Http_Get(url string) ([]byte, *http.Response, error) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Printf("recover:%v", e)
			}
		}()
	}()
	//var resp *http.Response
	//var err error
	//ch := make(chan struct{}, 1)
	//
	//// 创建一个协程来调用数据库访问函数
	//
	//	client := resty.New()
	//	client.SetRetryCount(3)
	//	client.SetTimeout(20 * time.Second)
	//	//client.SetHeader("User-Agent", fmt.Sprintf("Polygon.io GoClient/%v", clientVersion))
	//	client.SetHeader("Accept-Encoding", "gzip")
	//
	//	resp, err := client.R().Get(url)
	//	// 向管道写入一个空结构体
	//	ch <- struct{}{}
	//}()
	//
	//// 使用 select 多路复用来实现一个超时控制
	//select {
	//// 当数据库调用完毕则执行取出
	//case <-ch:
	//// 假如此300毫秒先到了，而readDB()还没有执行完毕则返回超时信息。
	//// 300ms => 此实践中并不会触发超时，这是由于我们模拟的数据库读取还是比较简单的。
	//// 此处使用 100s 来验证超时
	//case <-time.After(100 * time.Millisecond):
	//
	//}
	client := resty.New()
	client.SetRetryCount(3)
	client.SetTimeout(20 * time.Second)
	//client.SetHeader("User-Agent", fmt.Sprintf("Polygon.io GoClient/%v", clientVersion))
	client.SetHeader("Accept-Encoding", "gzip")
	resp, err := client.R().Get(url)
	return resp.Body(), resp.RawResponse, err
	//httpClient := &http.Client{
	//	Transport: &http.Transport{
	//		DialContext: (&net.Dialer{
	//			Timeout: 900 * time.Millisecond, // 连接超时
	//		}).DialContext,
	//		DialTLSContext: (&net.Dialer{
	//			Timeout: 900 * time.Millisecond, // 连接超时
	//		}).DialContext,
	//		DisableKeepAlives:   true,
	//		ForceAttemptHTTP2:   true,
	//		TLSHandshakeTimeout: 900 * time.Millisecond,
	//	},
	//	Timeout: time.Duration(2) * time.Second,
	//}
	//
	//resp, err := httpClient.Get(url)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil, nil, err
	//}
	//fmt.Println(resp.Body)
	//if err != nil {
	//	// 如果是因为超时引起的错误，则不打印错误信息
	//	if resp != nil && resp.StatusCode == http.StatusRequestTimeout {
	//		log.Println("请求超时，不打印错误信息")
	//	} else {
	//		log.Fatal(err) // 打印其他错误信息
	//	}
	//} else {
	//	// 处理响应的代码
	//	defer resp.Body.Close()
	//}
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//if resp.StatusCode == 200 {
	//	return body, resp, err
	//}
	//if err != nil {
	//
	//	logger.Info("Http_Get", err)
	//	return nil, nil, err
	//}

}
func Get_last_tradeday(code string) (string, string, string) {
	dates := Us_datatime()
	return dates[0:4], dates[5:7], dates[8:10]
}
func AnnualToDailyRate(annualRate float64) float64 {
	return annualRate / 252
}
func RandomDistributeOverSevenDays(totalAmount int, div int) []int {
	result := make([]int, div)
	remainingAmount := totalAmount

	for i := 0; i < div-1; i++ {
		randomAmount := rand.Intn(2*totalAmount) - totalAmount
		result[i] = randomAmount
		remainingAmount -= randomAmount
	}
	result[div-1] = remainingAmount
	return result
}

func CalCulateNetWorth(initialWorth float64, dailyProfitsOrLosses []float64) float64 {
	netWorth := initialWorth
	for _, dailyAmount := range dailyProfitsOrLosses {
		netWorth += dailyAmount
	}
	return netWorth
}
func IsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}
func RandFloat(min float64, max float64) float64 {
	rand.Seed(time.Now().UnixNano())

	// 生成 0 到 1 之间的随机浮点数
	randomFloat := rand.Float64()
	fmt.Println(randomFloat)

	randomInRange := min + rand.Float64()*(max-min)
	fmt.Println(randomInRange)
	return randomInRange
}

// 从切片中随机选择指定数量的元素
func SelectRandomElements(slice []string, count int) []string {
	result := make([]string, 0)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		index := rand.Intn(len(slice))
		result = append(result, slice[index])
	}

	return result
}
func GenerateOrderID() string {
	// 获取当前时间的纳秒值
	nano := strconv.Itoa(int(time.Now().UnixNano()))
	// 生成一个随机数
	randNum := rand.Intn(10000)
	return fmt.Sprintf("%s%d", nano[3:], randNum)
}
func Float64ToDecimal(f float64) decimal.Decimal {
	return decimal.NewFromFloat(math.Round(f*10000) / 10000)
}

// 计算回撤
func CalculateDrawdown(prices []float64) float64 {
	maxPrice := prices[0]
	maxDrawdown := 0.0

	for _, price := range prices {
		if price > maxPrice {
			maxPrice = price
		} else {
			drawdown := (maxPrice - price) / maxPrice
			if drawdown > maxDrawdown {
				maxDrawdown = drawdown
			}
		}
	}

	return maxDrawdown
}

// 计算 PnL
func CalculatePnL(initialInvestment float64, currentValue float64) float64 {
	return currentValue - initialInvestment
}

func GetRandomNumberInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
func Pullticketprice(code string) {
	key1 := "sub_ticket_zset"
	key2 := "sub_tickets"

	err := global.SHOP_REDIS.ZAdd(ctx, key1, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: code,
	}).Err()
	if err != nil {
		global.SHOP_LOG.Log(0, err.Error())

		return
	}
	err = global.SHOP_REDIS.ZAdd(ctx, key2, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: code,
	}).Err()
	if err != nil {
		global.SHOP_LOG.Log(0, err.Error())
		return
	}
}
func Get_us_trade_status() int {
	//"""
	//获取美股交易状态
	//0、停盘 1、盘前 2、盘中 3、盘后
	//:return:
	//"""

	et_times := Us_datatime()

	et_time := et_times[11:16]
	marketkey := fmt.Sprintf("market_status_%s", et_times[0:10])
	marketstatus, _ := global.SHOP_REDIS.Get(ctx, marketkey).Result()

	if marketstatus != "1" {
		return 0
	}
	if ("04:00" < et_time) && (et_time < "09:30") {
		return 1
	}

	if ("09:30" < et_time) && (et_time < "16:00") {
		return 2
	}

	if ("16:00" < et_time) && (et_time < "20:00") {
		return 3
	}

	return 0
}
func BatchDeleteKeys(client *redis.Client, pattern string) error {
	keys, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		fmt.Println("获取键列表时出错:", err)
		return err
	}
	for _, key := range keys {

		err := client.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("删除键", key, "时出错:", err)
			return err
		}
	}
	return nil
}
func BatchDeleteKeysbuvalue(client *redis.Client, pattern string, val string) error {
	keys, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		fmt.Println("获取键列表时出错:", err)
		return err
	}
	for _, key := range keys {
		//value, _ := client.Get(ctx, key).Result()

		v, _ := client.HGet(ctx, key, "adminId").Result()
		fmt.Println(key, v)
		if v == val {
			err := client.Del(ctx, key).Err()
			if err != nil {
				fmt.Println("删除键", key, "时出错:", err)
				return err
			}
		}
	}
	return nil
}

type Smsresp struct {
	Result    string `json:"result"`
	Code      string `json:"code"`
	Messageid string `json:"messageid"`
}

func SendSmsNq(apiURL, apiKey, secretKey, phoneNumber, message string, appcode string, retryTimes int) error {
	timestamp := time.Now().UnixMilli()
	data := make(map[string]interface{})

	sign := MD5V(fmt.Sprintf("%s%s%s", apiKey, secretKey, strconv.Itoa(int(timestamp))))

	data["appkey"] = apiKey
	data["phone"] = phoneNumber
	data["msg"] = message
	data["timestamp"] = strconv.Itoa(int(timestamp))
	data["appcode"] = appcode
	data["sign"] = sign
	req, _ := json.Marshal(&data)
	for i := 1; i <= retryTimes; i++ {
		resp, err := http.Post(apiURL, "application/json", bytes.NewReader(req))
		var smsresp Smsresp

		if err != nil {
			fmt.Printf("第 %d 次发送短信时出错: %v\n", i, err)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		print(string(body))
		if err != nil {
			fmt.Println("读取响应体时出错:", err)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New("读取响应体时出错")
		}
		json.Unmarshal(body, &smsresp)
		fmt.Println(smsresp)
		if smsresp.Code == "00000" {
			fmt.Println("短信发送成功")
			return nil
		} else {
			fmt.Printf("第 %d 次短信发送失败，状态码: %d\n", i, resp.StatusCode)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New(smsresp.Result)
		}
	}
	return nil
}
func SendSms(apiURL, apiKey, secretKey, phoneNumber, message string, retryTimes int) error {

	data := url.Values{}
	data.Set("appkey", apiKey)
	data.Set("phone", phoneNumber)
	data.Set("content", message)
	data.Set("secretkey", secretKey)

	for i := 1; i <= retryTimes; i++ {
		resp, err := http.PostForm(apiURL, data)
		var smsresp Smsresp

		if err != nil {
			fmt.Printf("第 %d 次发送短信时出错: %v\n", i, err)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取响应体时出错:", err)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New("读取响应体时出错")
		}
		json.Unmarshal(body, &smsresp)
		fmt.Println(smsresp)
		if smsresp.Code == "0" {
			fmt.Println("短信发送成功")
			return nil
		} else {
			fmt.Printf("第 %d 次短信发送失败，状态码: %d\n", i, resp.StatusCode)
			if i < retryTimes {
				time.Sleep(5 * time.Second) // 等待 5 秒后重试
				continue
			}
			return errors.New(smsresp.Result)
		}
	}
	return nil
}
func SendMailWithRetry(to string, message string, subject string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", global.SHOP_CONFIG.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", global.SHOP_CONFIG.Email.Subject)
	// m.SetBody("text/plain", message)
	m.SetBody("text/html", message)
	emailusername := global.SHOP_CONFIG.Email.EmailUserName
	emailpwd := global.SHOP_CONFIG.Email.EmailPwd

	// 创建发送器
	d := gomail.NewDialer(global.SHOP_CONFIG.Email.Smtp, global.SHOP_CONFIG.Email.Port, emailusername, emailpwd)

	// 重试次数和间隔
	retryTimes := 3
	retryInterval := 10 * time.Second

	//sendMailWithRetry(m, d, retryTimes, retryInterval)
	for i := 1; i <= retryTimes; i++ {
		if err := d.DialAndSend(m); err == nil {
			fmt.Println("邮件发送成功")
			return nil
		} else {
			fmt.Printf("第 %d 次发送邮件失败，错误: %v\n", i, err)
			if i < retryTimes {
				fmt.Printf("等待 %v 后重试\n", retryInterval)
				time.Sleep(retryInterval)
			}
			return errors.New(err.Error())
		}

	}
	fmt.Println("邮件发送多次失败，放弃")
	return nil
}

func ValidateEmail(email string) bool {
	// 定义电子邮件的正则表达式模式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func Get_Code_Key(types int8) string {
	switch types {
	case Reg_code:
		return "reg_code"
	case Change_pwd_code:
		return "reset_pwd_code"
	case Forgot_pwd_code:
		return "forgot_pwd_code"
	case Auth_code:
		return "auth_code"
	case Recharge_Msg:
		return "deposit_notification"
	case Whthdraw_Fail_Msg:
		return "withdrawal_fail"
	case Whthdraw_Success_Msg:
		return "withdrawal_notification_success"
	case Fund_Redemption_Successful:
		return "fund_redemption_successful"
	case Quan_Redemption_Successful:
		return "quan_redem_success"
	case Trade_Account_Pass:
		return "trade_account_apply_pass"
	case Trade_Account_Reject:
		return "trade_account_apply_rejected"
	case BankRecharge_Msg_Wait:
		return "deposit_notification_bank_wait"
	case BankRecharge_Msg_Success:
		return "deposit_notification_bank_success"
	case Review_Passed:
		return "review_passed"
	case Review_Rejected:
		return "review_rejected"
	default:
		return "Invalid_code"
	}
}
func Get_SendCode_Key(types int8) string {
	switch types {
	case Reg_code:
		return "reg_code"
	case Change_pwd_code:
		return "reset_pwd_code"
	case Forgot_pwd_code:
		return "forgot_pwd_code"
	case Auth_code:
		return "auth_code"
	case 5:
		return "cancel_account"
	case 6:
		return "reset_trade_pwd_code"
	case 7:
		return "bind_wallet"
	case 8:
		return "withdraw"

	default:
		return "Invalid_code"
	}
}
func IsEthAddress(address string) bool {
	// 以太坊地址的正则表达式模式
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
func IsBaseAddress(address string) bool {
	// 以太坊地址的正则表达式模式
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}
func IsBTCAddress(address string) bool {
	// 比特币地址的常见模式：以 1、3 或 bc1 开头，后面跟着数字和字母
	re := regexp.MustCompile(`^(1|[3]|[bc1])[a-zA-Z0-9]{25,34}$`)
	return re.MatchString(address)
}
func IsTronAddr(address string) bool {
	// 波场地址通常以 "T" 开头，后面跟着 33 个数字和字母
	re := regexp.MustCompile(`^T[a-zA-Z0-9]{33}$`)
	return re.MatchString(address)
}
func IsSolana(address string) bool {
	// 波场地址通常以 "T" 开头，后面跟着 33 个数字和字母
	re := regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]{32,44}$`)
	return re.MatchString(address)
}

func Get_crypto_current_price(code string) (float64, error) {
	apiKey := global.SHOP_CONFIG.System.POLYGON_API_KEY

	url := fmt.Sprintf("https://api.polygon.io/v1/last/crypto/%s/USD?apiKey=%s", code, apiKey)
	fmt.Println(url)

	result, _, err := Http_Get(url)
	var reslut Crypto
	if err != nil {
		global.SHOP_LOG.Log(0, err.Error())
	}

	json.Unmarshal([]byte(result), &reslut)
	if reslut.Status == "success" {
		return reslut.Last.Price, nil
	} else {
		return 0, errors.New("获取数据失败")
	}

}
func DecimalToFloat(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}
func Get_quan_status() bool {

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("转换时间出错:", err)
		return false
	}
	now := time.Now().In(loc)

	timenow := now.Format("2006-01-02 15:04:05")

	qs_date := timenow[0:10] //"2024-07-26" //
	timeend := fmt.Sprintf("%s 15:50:00", qs_date)
	timebegin := fmt.Sprintf("%s 09:30:00", qs_date)
	layout := "2006-01-02 15:04:05"

	tend, err := time.ParseInLocation(layout, timeend, loc)

	if err != nil {
		fmt.Println("转换时间出错:", err)
		return false

	}
	tbgin, err := time.ParseInLocation(layout, timebegin, loc)
	if err != nil {
		fmt.Println("转换时间出错:", err)
		return false
	}
	fmt.Println(now, tbgin, tend)
	if now.After(tend) {
		global.SHOP_LOG.Log(0, "结束时间到,子程序退出下单")
		return false
	}
	if now.Before(tbgin) {
		global.SHOP_LOG.Log(0, "开始时间未到,子程序退出下单")
		return false
	}
	return true
}
func Languagebycode(code string, content string) string {
	if code == "" {
		code = global.SHOP_CONFIG.System.Language
	}
	var jsonmap = make(map[string]string)
	err := json.Unmarshal([]byte(content), &jsonmap)
	if err != nil {
		fmt.Println("unmarshal failed:", err)
		return content
	}

	val, ok := jsonmap[code]
	if ok {
		return val
	} else {
		return jsonmap[global.SHOP_CONFIG.System.Language]
	}
}
func Languageresponse(code string, language string) string {
	if language == "" {
		language = global.SHOP_CONFIG.System.Language
	}
	content, _ := global.BlackCache.Get(fmt.Sprintf("language_%s", language)) // 输
	contents := content.(string)
	var jsonmap = make(map[string]string)
	err := json.Unmarshal([]byte(contents), &jsonmap)
	if err != nil {
		fmt.Println("unmarshal failed:", err)
		return ""
	}

	val, ok := jsonmap[code]
	if ok {
		return val
	} else {
		return code
	}
}
func ReadFileContent(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := make([]string, 0)
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return content, nil
}
func ReplaceSecondUppercase(s string) string {
	// 使用正则表达式匹配单词的第二个大写字母，但不捕获它
	re := regexp.MustCompile(`\b\p{Lu}\p{Lu}`)
	// 直接替换为横杠
	return re.ReplaceAllString(s, "_")
}

// transformSecondUppercase 将字符串中的第二个大写字母转换为下划线加字母的形式
func TransformSecondUppercase(s string) string {
	var (
		uppercaseCount int             // 用于计数大写字母的数量
		result         strings.Builder // 用于构建结果字符串
		//hasAdded       bool            // 标记是否已经添加了下划线加字母
	)

	for _, r := range s {
		if unicode.IsUpper(r) {
			uppercaseCount++         // 增加大写字母计数
			if uppercaseCount >= 2 { // 检查是否为第二个大写字母，且尚未添加过下划线加字母
				result.WriteRune('_')                // 添加下划线
				result.WriteRune(unicode.ToLower(r)) // 添加小写字母（可选，取决于是否希望保持原字母大小写）
				//hasAdded = true                      // 标记已添加过下划线加字母
			} else if uppercaseCount == 1 { // 检查是否为第二个大写字母，且尚未添加过下划线加字母

				result.WriteRune(unicode.ToLower(r)) // 添加小写字母（可选，取决于是否希望保持原字母大小写）

			} else {
				result.WriteRune(r) // 直接添加大写字母（或小写版本，如果需要的话）
			}
		} else {
			result.WriteRune(r) // 非大写字母直接添加
		}
	}

	return result.String()
}

type Crypto struct {
	Last      Last   `json:"last"`
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Symbol    string `json:"symbol"`
}
type Last struct {
	Conditions []int   `json:"conditions"`
	Exchange   int     `json:"exchange"`
	Price      float64 `json:"price"`
	Size       float64 `json:"size"`
	Timestamp  int64   `json:"timestamp"`
}
type InviteService struct {
	h *hashids.HashID
}

// NewInviteService 初始化服务
// salt: 加盐值，用于增加安全性，不同业务建议使用不同盐值
// minLength: 生成邀请码的最小长度
func NewInviteService(salt string, minLength int) (*InviteService, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	// 可选：自定义字符集，剔除易混淆字符如 0/O, 1/l/I
	hd.Alphabet = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, err
	}

	return &InviteService{h: h}, nil
}

// GenerateCode 根据用户ID生成邀请码
func (s *InviteService) GenerateCode(userID int) (string, error) {

	// Hashids 接受 int 切片，这里只传入一个用户ID
	code, err := s.h.Encode([]int{userID})
	if err != nil {
		return "", err
	}
	return code, nil
}

// DecodeCode 解析邀请码获取用户ID
func (s *InviteService) DecodeCode(code string) (int, error) {

	ids, err := s.h.DecodeWithError(code)
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, fmt.Errorf("invalid invite code")
	}
	id := ids[0]
	return id, nil
}
func IsValidMalaysiaPhone(phone string) bool {
	// 匹配规则：以1开头，第二位为3-9，后接9位数字
	reg := regexp.MustCompile(`^(?:\+?60)1\d{8,9}$`)
	return reg.MatchString(phone)
}
func IsValidPasswd(password string) bool {
	// 1. 检查长度
	if len(password) < 8 || len(password) > 16 {
		return false
	}

	// 2. 检查是否只包含字母和数字 (可选，根据需求决定是否允许特殊字符)
	// 如果允许特殊字符，可跳过此步或修改正则
	validChars := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !validChars.MatchString(password) {
		return false
	}

	// 3. 检查是否包含大写字母
	hasUpper := regexp.MustCompile(`[A-Z]`)
	if !hasUpper.MatchString(password) {
		return false
	}

	// 4. 检查是否包含小写字母
	hasLower := regexp.MustCompile(`[a-z]`)
	if !hasLower.MatchString(password) {
		return false
	}

	// 5. 检查是否包含数字
	hasDigit := regexp.MustCompile(`[0-9]`)
	if !hasDigit.MatchString(password) {
		return false
	}

	return true
}
func IsValidTradePasswd(password string) bool {
	//检查是否包含数字
	pattern := `^\d{6}$`
	hasDigit := regexp.MustCompile(pattern)

	if !hasDigit.MatchString(password) {
		return false
	}

	return true
}
func BuildInviteCode(userID int) string {
	service, err := NewInviteService("ShowGo", 6)
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}
	inviteCode, err := service.GenerateCode(userID)
	return inviteCode
}
func GetuidfromiCode(invitecode string) int {
	service, err := NewInviteService("ShowGo", 6)
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
	}
	uid, err := service.DecodeCode(invitecode)
	if err != nil {
		global.SHOP_LOG.Error(err.Error())
		return 0
	}
	return uid
}
func IsNetworkImage(url string) (bool, error) {
	// Step 1: 检查URL是否有效
	resp, err := http.Head(url) // 使用Head方法只获取响应头，不下载整个body
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码是否为200 OK
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP status code is not OK: %d", resp.StatusCode)
	}

	// Step 2: 检查Content-Type是否为图片类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return false, fmt.Errorf("Content-Type is not an image: %s", contentType)
	}

	// Step 3: 实际下载图片并读取一小部分内容进行进一步验证（可选）
	// 这里我们只验证了Content-Type，如果要进一步验证可以下载一小部分内容进行检查。
	// 但通常Content-Type足够用来确定是否为图片。

	return true, nil // 如果所有检查都通过，则认为是网络图片
}
