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

func (u *User) SetUserAll(arr []string){
    u.xx     = arr[1]
    u.xm     = arr[0]
    u.kslb   = arr[2]
    u.zkzh   = arr[3]
    u.kssj   = arr[4]
    u.zf     = arr[5]
    u.tl     = arr[6]
    u.yd     = arr[7]
    u.xzyfy  = arr[8]
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

func main() {
    kh, xm, err := mainAgrs()
    if err != nil {
        fmt.Println("WARN: no arguments put int command.\n")

        reader := bufio.NewReader(os.Stdin)
        
        fmt.Print("输入准考证号：")
        input, _,  _ := reader.ReadLine()
        kh = string(input)

        fmt.Print("    输入姓名：")
        input, _, _ = reader.ReadLine()
        xm = string(input)

        //fmt.Printf(": %v --\n", string([]byte(input)))
        //fmt.Printf(": %v --\n", string([]byte(input)))
    }

    arr := getPageAndParse(kh, xm)

    user := new(User)
    user.SetUserAll(arr)

    fmt.Println(user.xx)
    fmt.Println(user.xm)
    fmt.Println(user.zkzh)
    fmt.Println(user.kslb)
    fmt.Println(user.kssj)
    fmt.Printf("-----------------\n")
    fmt.Printf("%-9s%6s\n", "总分:", user.zf)
    fmt.Printf("%-9s%6s\n", "听力:", user.tl)
    fmt.Printf("%-9s%6s\n", "阅读:", user.yd)
    fmt.Printf("%-7s%6s\n", "写作翻译:", user.xzyfy)

    zf_int, _ := strconv.Atoi(user.zf)
    if zf_int > 425 {
        fmt.Println("\n通过考试！\n")
    } else {
        fmt.Println("\n考试没有通过\n")
    }

    fmt.Printf("回车键退出...")
    fmt.Scanln()

}