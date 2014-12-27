package main

import (
    "fmt"
    "net/http"
    "net/url"
    "io/ioutil"
    "os"
    "log"
    "strings"
    "regexp"
    "strconv"
    "errors"
    "bufio"
)

func getPage(zkzh, xm string) []byte {
    client := &http.Client{}
    q_path := requestPath(zkzh, xm)
    //fmt.Println(q_path)

    req, err := http.NewRequest("GET", q_path, nil)
    if err != nil {
        log.Fatal(err);
    }
    req.Header.Add("Referer", `W/"http://www.chsi.com.cn/cet/"`)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err);
    }
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err);
    }

    return bodyBytes
}

func tableNoBlank(bodyBytes []byte) string {
    tableBytes := bodyBytes[5600:7600]
    tableString := string(tableBytes)

    s_tamp := strings.Replace(tableString, "\t", "", -1)
    s_tamp = strings.Replace(s_tamp, " ", "", -1)
    s_tamp = strings.Replace(s_tamp, "\r\n", "", -1)
    return s_tamp;
}

type User struct {
    xx      string
    xm      string
    kslb    string
    zkzh    string
    kssj    string
    zf      string
    tl      string
    yd      string
    xzyfy   string
}

type CachedUser struct {
    kh    string
    xm    string
}



func (self *User) SetUserAll(arr []string){
    self.xx     =  arr[1]
    self.xm     =  arr[0]
    self.kslb   =  arr[2]
    self.zkzh   =  arr[3]
    self.kssj   =  arr[4]
    self.zf     =  arr[5]
    self.tl     =  arr[6]
    self.yd     =  arr[7]
    self.xzyfy  =  arr[8]
}



func parsePage(tbStr string) ([]string) {
    re, err := regexp.Compile(`(<td>.{0,15}<\/td>|>\d{2,3})`)
    if err != nil {
        log.Fatal(err);
    }
    arr := re.FindAllString(tbStr, -1)

    if len(arr) != 9 {
        fmt.Println("parse in tableStringNoBlank error")
        fmt.Println("ERROR: parsePage() error")
        os.Exit(1)
    }

    for i := 0; i < 5; i++ {
        arr[i] = arr[i][4:len(arr[i])-5]
    }
    for i := 5; i < len(arr); i++ {
        arr[i] = arr[i][1:len(arr[i])]
    }
    return arr;
}


func requestPath(zkzh, xm string) string {
    xm = url.QueryEscape(xm)
    s := "http://www.chsi.com.cn/cet/query?zkzh=" + zkzh + "&xm=" + xm;
    //fmt.Println(s)
    return s;
}

func getPageAndParse(zkzh, xm string) ([]string) {
    bodyBytes := getPage(zkzh, xm)
    tableStringNoBlank := tableNoBlank(bodyBytes)
    user_elements := parsePage(tableStringNoBlank)

    return user_elements
}

func mainAgrs() (string, string, error) {
    if len(os.Args) == 3 {
        if len(os.Args[1]) > 12 {
            return os.Args[1], os.Args[2], nil;
        } else {
            return os.Args[2], os.Args[1], nil;
        }
    } else {
        return "", "", errors.New("no main arguments");
    }
    
}



func readFromCacheFile(filename string) ([]CachedUser) {
    f, _ := os.Open(filename)
    defer f.Close()
    reader := bufio.NewReader(f)
    Bytes, _ := ioutil.ReadAll(reader)

    re_kh, _ := regexp.Compile(`\d{15}`)
    re_xm, _ := regexp.Compile(`,.{2,6};`)

    arr_kh := re_kh.FindAllString(string(Bytes), -1)
    arr_xm := re_xm.FindAllString(string(Bytes), -1)

    len_kh := len(arr_kh)
    len_xm := len(arr_xm)

    for i := 0; i < len_xm; i++ {
        length := len(arr_xm[i])
        arr_xm[i] = arr_xm[i][1:length-1]
    }

    var len_min = 0
    if len_xm > len_kh {
        len_min = len_kh
    } else {
        len_min = len_xm
    }

    cachedUsers := make([]CachedUser, len_min)

    for i := 0; i < len_min; i++ {
        cachedUsers[i].xm = arr_xm[i]
        cachedUsers[i].kh = arr_kh[i]
    }

    return cachedUsers
}

func enterFromCommand() (kh, xm string) {
    fmt.Println("WARN: no arguments put int command.\n")

    reader := bufio.NewReader(os.Stdin)

    fmt.Print("输入准考证号：")
    input, _,  _ := reader.ReadLine()
    kh = string(input)

    fmt.Print("    输入姓名：")
    input, _, _ = reader.ReadLine()
    xm = string(input)
    return kh, xm
}

func (self *User) PrintOut() {
    fmt.Println(self.xx)
    fmt.Println(self.xm)
    fmt.Println(self.zkzh)
    fmt.Println(self.kslb)
    fmt.Println(self.kssj)
    fmt.Printf("-----------------\n")
    fmt.Printf("%-9s%6s\n", "总分:", self.zf)
    fmt.Printf("%-9s%6s\n", "听力:", self.tl)
    fmt.Printf("%-9s%6s\n", "阅读:", self.yd)
    fmt.Printf("%-7s%6s\n", "写作翻译:", self.xzyfy)

    zf_int, _ := strconv.Atoi(self.zf)
    if zf_int > 425 {
        fmt.Println("\n通过考试！\n")
    } else {
        fmt.Println("\n考试没有通过\n")
    }
}

func DoOneUser(kh, xm string) (*User) {
    arr := getPageAndParse(kh, xm)
    user := new(User)
    user.SetUserAll(arr)
    user.PrintOut()
    return user
}


func main() {
    kh, xm, err := mainAgrs()
    if err != nil {
        if _, err := os.Stat("user-cache.txt"); err != nil {
            kh, xm = enterFromCommand()
            _ = DoOneUser(kh, xm)
        } else {
            cachedUsers := readFromCacheFile("user-cache.txt")
            for i := 0; i < len(cachedUsers); i++ {
                fmt.Println("===================\n")
                kh = cachedUsers[i].kh
                xm = cachedUsers[i].xm

                _ = DoOneUser(kh, xm)
            }
        }
    } else {
        _ = DoOneUser(kh, xm)
    }


    fmt.Printf("回车键退出...")
    fmt.Scanln()

}