package hwclock

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const TIME_LAYOUT = "2006-01-02 15:04:05"

// https://go.googlesource.com/sys/+/master/unix/syscall_linux_test.go
// https://go.googlesource.com/sys/+/master/unix
func main() {
	f, err := os.Open("/dev/rtc0")
	if err != nil {
		fmt.Println("skipping test, %v", err)
	}
	defer f.Close()
	v, err := unix.IoctlGetRTCTime(int(f.Fd()))
	if err != nil {
		fmt.Println("failed to perform ioctl: %v", err)
	}
	ts := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", v.Year+1900, v.Mon+1, v.Mday, v.Hour, v.Min, v.Sec)
	fmt.Println("RTC time: ", ts)
	zone := getTimeZone()
	fmt.Println(zone)
	t2, err := parseWithLocation(zone, ts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t2)
}

func parseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, err
	} else {
		// 转成带时区的时间
		lt, _ := time.ParseInLocation(TIME_LAYOUT, timeStr, l)
		// 直接转成相对时间
		fmt.Println(time.Now().In(l).Format(TIME_LAYOUT))
		return lt, nil
	}
}
func testTime() {
	str := time.Now().Format("2006-01-02 15:04:05")
	// 指定时区
	t1, err := parseWithLocation("America/Cordoba", str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t1)

	t2, err := parseWithLocation("Asia/Shanghai", str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t2)

	t3, err := parseWithLocation("Asia/Chongqing", str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t3)
}

func getTimeZone() (timeZone string) {
	if fi, err := os.Lstat("/etc/localtime"); err == nil {
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			if tzfile, err := os.Readlink("/etc/localtime"); err == nil {
				if strings.Contains(tzfile, "/usr/share/zoneinfo/") {
					l := strings.Split(tzfile, "/usr/share/zoneinfo/")
					timeZone = l[len(l)-1]
					return
				}
			}
		}
	}

	if timezone := echo("/etc/timezone"); timezone != "" {
		timeZone = timezone
		return
	}

	if f, err := os.Open("/etc/sysconfig/clock"); err == nil {
		defer f.Close()
		s := bufio.NewScanner(f)
		for s.Scan() {
			if sl := strings.Split(s.Text(), "="); len(sl) == 2 {
				if sl[0] == "ZONE" {
					timeZone = strings.Trim(sl[1], `"`)
					return
				}
			}
		}
	}
	return
}

func echo(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
