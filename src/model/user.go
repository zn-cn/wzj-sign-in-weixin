package model

import (
	"bytes"
	"constant"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"model/db"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"util"

	jsoniter "github.com/json-iterator/go"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomodule/redigo/redis"

	"gopkg.in/mgo.v2/bson"
)

// redis 中暂存用户状态
type User struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Openid string        `bson:"openid" json:"openid"`
	Email  string        `bson:"email" json:"email"`
	// redis 中存储当前选中坐标
	Coordinates map[string]Coordinate `bson:"coordinates" json:"coordinates"` // 标签 -> 坐标
}

type Coordinate struct {
	Lon float64 `bson:"lon" json:"lon"`
	Lat float64 `bson:"lat" json:"lat"`
}

type WZJCourse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Topic        string `json:"topic"`
	Code         string `json:"code"`
	College      string `json:"college"`
	Department   string `json:"department"`
	DiscussionID int    `json:"discussionId"`
	Selected     bool   `json:"selected"`
}

func SetUserEmail(openid, email string) error {
	query := bson.M{
		"openid": openid,
	}
	update := bson.M{
		"$set": bson.M{
			"email": email,
		},
	}
	return updateUser(query, update)
}

func AddSignInTask(openid, textOpenid string) error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	_, err := conn.Do("SETEX", fmt.Sprintf(constant.RedisUserTask, openid), 3600*2, textOpenid)
	return err
}

func GetRedisUserStatus(openid string) (int, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	return redis.Int(conn.Do("GET", fmt.Sprintf(constant.RedisUserStatus, openid)))
}

func SetRedisUserStatus(openid string, status int) error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	_, err := conn.Do("SETEX", fmt.Sprintf(constant.RedisUserStatus, openid), 60*10, status)
	return err
}

func AddUserCoordinate(openid, tag string, coordinate Coordinate) error {
	query := bson.M{
		"openid": openid,
	}

	// 相同覆盖
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("coordinates.%s", tag): coordinate,
		},
	}
	err := updateUser(query, update)
	if err != nil {
		return err
	}
	selector := bson.M{
		"coordinates": 1,
	}
	user, _ := findUser(query, selector)
	if len(user.Coordinates) == 0 {
		go SetUserCurCoordinateByTag(openid, tag)
	}
	return err
}

func SetUserCurCoordinateByTag(openid, tag string) (Coordinate, error) {
	var coordinate Coordinate
	selector := bson.M{
		"coordinates": 1,
	}
	user, err := getUserByOpenid(openid, selector)
	if err != nil {
		return coordinate, err
	}

	if c, ok := user.Coordinates[tag]; ok {
		setUserCurCoordinate(openid, c)
		coordinate = c
	}
	return coordinate, err
}

func setUserCurCoordinate(openid string, coordinate Coordinate) error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	_, err := conn.Do("SET", fmt.Sprintf(constant.RedisUserCurCoordinate, openid),
		fmt.Sprintf(constant.RedisUserCoordinateFormat, coordinate.Lon, coordinate.Lat))
	return err
}

func SetUserCurCourse(openid, courseName string, courseID int) error {
	textOpenid, err := getTextOpenid(openid)
	if err != nil {
		return err
	}

	client := &http.Client{}
	data := map[string]interface{}{}
	data["courseId"] = courseID
	data["courseName"] = courseName
	dataByte, _ := jsoniter.Marshal(data)

	req, err := http.NewRequest("POST", constant.URLWZJCourseSelect, bytes.NewBuffer(dataByte))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.2.3 NetType/WIFI Language/zh_CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("openId", textOpenid)
	req.Header.Set("Host", "v18.teachermate.cn")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func GetUserCoordinates(openid string) (map[string]Coordinate, error) {
	selector := bson.M{
		"coordinates": 1,
	}
	user, err := getUserByOpenid(openid, selector)
	return user.Coordinates, err
}

func StartStuSignTask() error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	keys, _ := redis.Strings(conn.Do("KEYS", fmt.Sprintf(constant.RedisUserTask, "*")))
	for _, key := range keys {
		textOpenid, _ := redis.String(conn.Do("GET", key))
		if textOpenid == "" {
			continue
		}
		openid := key[len(constant.RedisUserTask)-2:]
		c, _ := redis.String(conn.Do("GET", fmt.Sprintf(constant.RedisUserCurCoordinate, openid)))
		if c == "" {
			continue
		}
		cs := strings.Split(c, ",")
		if len(cs) != 2 {
			continue
		}
		lon, _ := strconv.ParseFloat(cs[0], 64)
		lat, _ := strconv.ParseFloat(cs[1], 64)
		coordinate := Coordinate{
			Lon: lon,
			Lat: lat,
		}
		go func(textOpenid, openid string) {
			ok, _ := userCheckIn(textOpenid, coordinate)
			if ok {
				signSuccessNotice(openid)
			}
		}(textOpenid, openid)
		go checkIsHaveDiscuss(openid)
	}
	return nil
}

// Find the review items
var nameMap = map[string]string{
	"token-hash": "wx_csrf_name",
	"openid":     "openid",
	"sign-id":    "signId",
	"course-id":  "courseId",
}

func userCheckIn(textOpenid string, coordinate Coordinate) (bool, error) {
	client := &http.Client{}

	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(constant.URLWZJSignIn, textOpenid))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false, errors.New("request wrong")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return false, err
	}

	data := url.Values{}
	// is-sign apply-gps is-qr-sign
	applyGps := false
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		if id, idOk := s.Attr("id"); idOk {
			if name, nameOk := nameMap[id]; nameOk {
				value, _ := s.Attr("value")
				data.Set(name, value)
			}
			if id == "apply-gps" {
				if value, _ := s.Attr("value"); value == "1" {
					applyGps = true
				}
			}
		}
	})

	if data.Get("courseId") != "" && data.Get("openid") != "" && data.Get("signId") != "" {
		if applyGps {
			rand.Seed(time.Now().UnixNano())
			// 随机化处理，防止一致
			coordinate.Lon += float64(rand.Intn(40)-20) * 0.000001
			coordinate.Lat += float64(rand.Intn(40)-20) * 0.000001
			data.Set("lon", strconv.FormatFloat(coordinate.Lon, 'f', 5, 64)) // 5 表示截断为5位小数
			data.Set("lat", strconv.FormatFloat(coordinate.Lat, 'f', 5, 64))
		}

		req, err := http.NewRequest("POST", constant.URLWZJStuSignIn, ioutil.NopCloser(strings.NewReader(data.Encode())))
		if err != nil {
			return false, err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.2.3 NetType/WIFI Language/zh_CN")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		res, err = client.Do(req)
		if err != nil {
			return false, err
		}
		defer res.Body.Close()
		if res.StatusCode == 200 {
			return true, nil
		}
	}
	return false, err
}

func signSuccessNotice(openid string) error {
	user, _ := getUserByOpenid(openid, bson.M{"email": 1})
	if user.Email == "" {
		return errors.New("email empty")
	}
	// 垃圾邮件？
	content := `<body style="margin: 0; padding: 0;">

　<table border="1" cellpadding="0" cellspacing="0" width="100%">

　　<tr>
　　　<td> Hello World! </td>
　　</tr>

　</table>

</body>`
	util.SendEmail("阿楠技术", "微助教签到成功提醒", content, []string{user.Email})
	return nil
}

// 检测是否有讨论并邮件提醒
func checkIsHaveDiscuss(openid string) error {
	user, err := getUserByOpenid(openid, bson.M{"email": 1})
	if err != nil {
		return err
	}
	if user.Email == "" {
		return errors.New("email is empty")
	}
	courses, err := ListWZJDiscussCourses(openid)
	if err != nil {
		return err
	}
	for _, course := range courses {
		if course.Topic != "" {
			content := fmt.Sprintf("课程名: %s<br/>讨论话题: %s<br/>课程是否被选中: %v", course.Name, course.Topic, course.Selected)
			if ok, _ := getRedisUserCourseNotice(openid, course.ID, course.DiscussionID); !ok {
				go util.SendEmail("阿楠技术", "微助教讨论提醒", content, []string{user.Email})
				setRedisUserCourseNotice(openid, course.ID, course.DiscussionID)
			}
		}
	}
	return nil
}

func getRedisUserCourseNotice(openid string, courseID, discussionID int) (bool, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	return redis.Bool(conn.Do("GET", fmt.Sprintf(constant.RedisUserDisCourseNotice, openid, courseID, discussionID)))
}

func setRedisUserCourseNotice(openid string, courseID, discussionID int) error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	_, err := conn.Do("SETEX", fmt.Sprintf(constant.RedisUserDisCourseNotice, openid, courseID, discussionID), 3600*24, true)
	return err
}

func ListWZJDiscussCourses(openid string) ([]WZJCourse, error) {
	textOpenid, err := getTextOpenid(openid)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", constant.URLWZHDisCourseList, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.2.3 NetType/WIFI Language/zh_CN")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("openId", textOpenid)
	req.Header.Set("Host", "v18.teachermate.cn")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	dataByte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = jsoniter.Unmarshal(dataByte, &data)
	if err != nil {
		return nil, err
	}
	coursesByte, err := jsoniter.Marshal(data["courses"])
	if err != nil {
		return nil, err
	}
	courses := []WZJCourse{}
	err = jsoniter.Unmarshal(coursesByte, &courses)
	return courses, err
}

func getTextOpenid(openid string) (string, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	conn := cntrl.GetConn()

	return redis.String(conn.Do("GET", fmt.Sprintf(constant.RedisUserTask, openid)))
}

func createUser(openid string) error {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	userTable := cntrl.GetTable(constant.TableUser)
	user := User{}
	query := bson.M{
		"openid": openid,
	}
	err := userTable.Find(query).One(&user)
	if err != nil {
		user.ID = bson.NewObjectId()
		user.Openid = openid
		err = userTable.Insert(user)
	}
	return err
}

func getUserByOpenid(openid string, selector bson.M) (User, error) {
	query := bson.M{
		"openid": openid,
	}
	return findUser(query, selector)
}

func findUser(query, selector bson.M) (User, error) {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	userTable := cntrl.GetTable(constant.TableUser)
	user := User{}
	err := userTable.Find(query).Select(selector).One(&user)
	return user, err
}

func updateUser(query, update bson.M) error {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	userTable := cntrl.GetTable(constant.TableUser)

	return userTable.Update(query, update)
}
