package model

import (
	"constant"
	"errors"
	"fmt"
	"io/ioutil"
	"model/db"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomodule/redigo/redis"

	"gopkg.in/mgo.v2/bson"
)

// redis 中暂存用户状态
type User struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Openid string        `bson:"openid" json:"openid"`
	// redis 中存储当前选中坐标
	Coordinates map[string]Coordinate `bson:"coordinates" json:"coordinates"` // 标签 -> 坐标
}

type Coordinate struct {
	Lon float64 `bson:"lon" json:"lon"`
	Lat float64 `bson:"lat" json:"lat"`
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
	query := bson.M{
		"openid": openid,
	}
	selector := bson.M{
		"coordinates": 1,
	}
	user, err := findUser(query, selector)
	if err != nil {
		return coordinate, err
	}

	if c, ok := user.Coordinates[tag]; ok {
		cntrl := db.NewRedisDBCntlr()
		defer cntrl.Close()
		conn := cntrl.GetConn()

		_, err = conn.Do("SET", fmt.Sprintf(constant.RedisUserCurCoordinate, openid),
			fmt.Sprintf(constant.RedisUserCoordinateFormat, c.Lon, c.Lat))
		coordinate = c
	}
	return coordinate, err
}

func GetUserCoordinates(openid string) (map[string]Coordinate, error) {
	query := bson.M{
		"openid": openid,
	}
	selector := bson.M{
		"coordinates": 1,
	}
	user, err := findUser(query, selector)
	return user.Coordinates, err
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
		c, _ := redis.String(conn.Do("GET", fmt.Sprintf(constant.RedisUserCurCoordinate, key[len(constant.RedisUserTask)-2:])))
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
		go userCheckIn(textOpenid, coordinate)
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

func userCheckIn(textOpenid string, coordinate Coordinate) error {
	client := &http.Client{}

	// Request the HTML page.
	res, err := http.Get(fmt.Sprintf(constant.URLWZJSignIn, textOpenid))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("request wrong")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	data := url.Values{}
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		if id, idOk := s.Attr("id"); idOk {
			if name, nameOk := nameMap[id]; nameOk {
				value, _ := s.Attr("value")
				data.Set(name, value)
			}
		}
	})

	req, err := http.NewRequest("POST", constant.URLWZJStuSignIn, ioutil.NopCloser(strings.NewReader(data.Encode())))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_4 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Mobile/12H143 MicroMessenger/6.2.3 NetType/WIFI Language/zh_CN")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	_, err = client.Do(req)
	return err
}
