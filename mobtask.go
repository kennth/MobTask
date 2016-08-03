package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	cmd := flag.Int("cmd", 0, "create runscript for console")
	id := flag.Int("id", 0, "create runscript for console")
	flag.Parse()
	//fmt.Println(*cmd)
	//fmt.Println(*id)
	switch {
	case *cmd == -1:

		fmt.Printf("%s", time.Now())
	case *cmd == 0:
		genRunScript()
	case *cmd == 1:
		updatehook(*id)
	case *cmd == 2:
		getHOOKActivity()
	case *cmd == 8:
		keepActivityClean(*id)
	case *cmd == 9:
		keepActivityAlive(*id)
	}
}

func keepActivityClean(id int) {
	db, err := sql.Open("mysql", "root:funmix@tcp(192.168.99.10:3306)/helper?charset=utf8")
	CheckErr(err)
	sql := "select worker,activity,hook from tcmcctask where status>-9 and id=" + strconv.Itoa(id)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)
	var activity string
	var worker string
	var hook string
	var startphone int64
	var endphone int64
	var phoneid string
	var nextphone string
	var packname string
	if rows.Next() {
		err = rows.Scan(&worker, &activity, &hook)
		CheckErr(err)
		phoneid = worker[0:strings.Index(worker, " ")]
		startphone, err = strconv.ParseInt(phoneid, 10, 0)
		phoneid = worker[strings.Index(worker, " ")+1:]
		endphone, err = strconv.ParseInt(phoneid, 10, 0)
		packname = activity[0:strings.Index(activity, "/")]
	} else {
		fmt.Printf("%d is not exsits!", id)
		return
	}
	db.Close()
	for {
		for i := startphone; i <= endphone; i++ {
			phoneid = "E3CD20" + strconv.Itoa(int(i))
			if i+1 > endphone {
				nextphone = "E3CD20" + strconv.Itoa(int(startphone))
			} else {
				phoneid = "E3CD20" + strconv.Itoa(int(i+1))
			}
			f, err := exec.Command("/bin/sh", "-c", "adb -s "+nextphone+" shell am start -n "+activity).Output()
			if err == nil {
				fmt.Println(string(f))
			} else {
				fmt.Println(err.Error())
				f, err = exec.Command("/bin/sh", "-c", "adb -s "+nextphone+" shell am start -n "+activity).Output()
				if err == nil {
					fmt.Println(string(f))
				} else {
					fmt.Println(err.Error())
				}
			}

			f, err = exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell am force-stop "+packname).Output()
			if err == nil {
				fmt.Println(string(f))
			} else {
				fmt.Println(err.Error())
			}
			f, err = exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell am force-stop "+packname).Output()
			if err == nil {
				fmt.Println(string(f))
			} else {
				fmt.Println(err.Error())
			}
			time.Sleep(time.Second * 1)
			f, err = exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell rm -f /data/data/com.sg.hlw.vivo/files/c_data_store.dat").Output()
			if err == nil {
				fmt.Println(string(f))
			} else {
				fmt.Println(err.Error())
			}
			time.Sleep(time.Second * 5)
		}
		time.Sleep(time.Second * 1)
	}

}

func keepActivityAlive(id int) {
	db, err := sql.Open("mysql", "root:funmix@tcp(192.168.99.10:3306)/helper?charset=utf8")
	CheckErr(err)
	sql := "select worker,activity,hook from tcmcctask where status>-9 and id=" + strconv.Itoa(id)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)
	var activity string
	var hookactivity string
	var worker string
	var hook string
	var curactivity string
	var startphone int64
	var endphone int64
	var phoneid string
	var packname string
	var wait float64
	if rows.Next() {
		err = rows.Scan(&worker, &activity, &hook)
		CheckErr(err)
		phoneid = worker[0:strings.Index(worker, " ")]
		startphone, err = strconv.ParseInt(phoneid, 10, 0)
		phoneid = worker[strings.Index(worker, " ")+1:]
		endphone, err = strconv.ParseInt(phoneid, 10, 0)
		packname = activity[0:strings.Index(activity, "/")]
		if len(hook) == 0 {
			hookactivity = activity
		} else {
			hookactivity = packname + "/" + hook
		}
	} else {
		fmt.Printf("%d is not exsits!", id)
		return
	}
	db.Close()
	t := time.Now()
	wait = 0
	for {
		if time.Now().Sub(t).Seconds() >= wait {
			for i := startphone; i <= endphone; i++ {
				phoneid = "E3CD20" + strconv.Itoa(int(i))
				f, geterr := exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell dumpsys activity | grep \"mFocusedActivity\"").Output()
				//CheckErr(err)
				if geterr == nil {
					curactivity = string(f)
					curactivity = curactivity[strings.Index(curactivity, packname) : len(curactivity)-3]
					//fmt.Println("adb -s " + phoneid + " shell am start -n " + activity)
					if !strings.EqualFold(curactivity, activity) && !strings.EqualFold(curactivity, hookactivity) {
						fmt.Println(phoneid + ":" + curactivity + " || " + activity + " || " + hookactivity)
						f, err := exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell am force-stop "+packname).Output()
						if err == nil {
							fmt.Println(string(f))
						} else {
							fmt.Println(err.Error())
						}
						f, err = exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell am start -n "+activity).Output()
						if err == nil {
							fmt.Println(string(f))
						} else {
							fmt.Println(err.Error())
						}
						wait = 1
					} else {
						//fmt.Println(phoneid + ":" + mainact + " is on top!")
						wait = wait + 1
					}
				} else {
					wait = wait + 1
					fmt.Println(geterr.Error())
				}
			}
			t = time.Now()
			if wait > 5 {
				wait = 60
			} else {
				wait = 1
			}
			fmt.Printf("%d : %ds || %s\n", id, int(wait), t)
		}
		time.Sleep(time.Second * 1)
	}

}

func getHOOKActivity() {
	id := GetHostID() - 5
	db, err := sql.Open("mysql", "root:funmix@tcp(192.168.99.10:3306)/helper?charset=utf8")
	CheckErr(err)
	sql := "select id,worker,activity,hook from tcmcctask where status>-9 and id>" + strconv.Itoa(id*12) + " and id<=" + strconv.Itoa((id+1)*12) + " order by id"
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)
	var taskid string
	var activity string
	var worker string
	var hook string
	var phoneid string
	var curactivity string
	var curhook string
	stmt, err := db.Prepare("UPDATE tcmcctask SET hook = ? WHERE id=?")
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&taskid, &worker, &activity, &hook)
		CheckErr(err)
		phoneid = "E3CD20" + worker[0:strings.Index(worker, " ")]
		f, err := exec.Command("/bin/sh", "-c", "adb -s "+phoneid+" shell dumpsys activity | grep \"mFocusedActivity\"").Output()
		//CheckErr(err)
		if err == nil {
			curactivity = string(f)
			curhook = curactivity[strings.Index(curactivity, "/")+1 : len(curactivity)-3]
			fmt.Println(taskid + ":" + curhook)
			_, err := stmt.Exec(curhook, taskid)
			CheckErr(err)
		}
		//db.Exec("update tcmcctask set hook=")
	}
	db.Close()

	fmt.Println("done!")
}

func updatehook(id int) {
	db, err := sql.Open("mysql", "root:funmix@tcp(192.168.99.10:3306)/helper?charset=utf8")
	CheckErr(err)
	sql := "select worker,activity,hook from tcmcctask where status>=0 and id=" + strconv.Itoa(id)
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)

	if rows.Next() {
		var activity string
		var worker string
		var mainact string
		var hook string
		err = rows.Scan(&worker, &activity, &hook)
		CheckErr(err)
		if len(hook) == 0 {
			mainact = activity[strings.Index(activity, "/")+1:]
		} else if strings.Index(hook, ".") == 1 {
			mainact = activity[0:strings.Index(activity, "/")] + hook
		} else {
			mainact = hook
		}
		fmt.Println(mainact)
		err := ioutil.WriteFile("/home/funmix/ADBWorker/hook.conf", []byte(mainact), 0666)
		CheckErr(err)
	}
	db.Close()

	fmt.Println("done!")
}

func genRunScript() {
	id := GetHostID() - 5
	db, err := sql.Open("mysql", "root:funmix@tcp(192.168.99.10:3306)/helper?charset=utf8")
	CheckErr(err)
	sql := "select worker,activity from tcmcctask where status>-9 and id>" + strconv.Itoa(id*12) + " and id<=" + strconv.Itoa((id+1)*12) + " order by id"
	fmt.Println(sql)
	rows, err := db.Query(sql)
	CheckErr(err)
	script := "cd /home/funmix/ADBWorker\necho \"$(date)\" >> rerun.log\n"
	for rows.Next() {
		var activity string
		var worker string
		err = rows.Scan(&worker, &activity)
		CheckErr(err)
		fmt.Println("./startapp.sh " + worker + " " + activity + " >> rerun.log &")
		script = script + "./startapp.sh " + worker + " " + activity + " >> rerun.log &\n"
	}
	db.Close()

	err = ioutil.WriteFile("/home/funmix/ADBWorker/rerun.sh", []byte(script), 0666)
	CheckErr(err)
	fmt.Println("done!")
}

func GetHostID() int {
	f, err := exec.Command("hostname").Output()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		hostname := strings.Replace(string(f), "\n", "", -1)
		id := hostname[len(hostname)-1 : len(hostname)]
		i, err := strconv.Atoi(id)
		CheckErr(err)
		return i
	}
	return 0
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
