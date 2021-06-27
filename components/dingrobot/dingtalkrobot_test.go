package dingrobot

import (
	"testing"
)

func TestDingTalkRobot_TextMsg_Send(t *testing.T) {
	url := "https://oapi.dingtalk.com/robot/send?access_token=d3acd76b80dcb5619c42bcb3fc9b26f148605c55ea69527ec6cbc949be97e143"
	secret := "SECd22016a925b16fa14c506bb1d4380eef1773a5b078d642cd7f2ccddcb40a9726"
	robot, _ := NewDingDingTalkRobot(url, secret)

	/*msg := NewTextMsg()

	msg.Text.Content = "这是一个测试aaaa"
	msg.At.AtMobiles = []string{"13683506199"}
	msg.At.AtUserIds = []string{"lx6g9h8"}*/

	msg := NewLinkMsg()

	msg.Link.MessageUrl = "https://www.baidu.com/"
	msg.Link.PicUrl = "https://i1.go2yd.com/image.php?url=0OqXSMXbgX"
	msg.Link.Text = "这是一个寂寞的夜， 下着有点伤心的雨"
	msg.Link.Title = "让累化作伤心雨"

	res, err := robot.Send(msg)
	if err != nil {
		t.Log(err)
	}
	t.Log(res)

}

func TestDingTalkRobot_BulidSign(t *testing.T) {
	url := "https://oapi.dingtalk.com/robot/send?access_token=d3acd76b80dcb5619c42bcb3fc9b26f148605c55ea69527ec6cbc949be97e143"
	secret := "SECd22016a925b16fa14c506bb1d4380eef1773a5b078d642cd7f2ccddcb40a9726"
	robot, _ := NewDingDingTalkRobot(url, secret)

	var millSeconds int64 = 1618228844084

	s := robot.buildSign(millSeconds)
	t.Log(s)
}
