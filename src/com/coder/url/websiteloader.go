package url

import (
	_ "fmt"
	"os"
	"bufio"
    "encoding/csv"
	"io"
	"fmt"
)

const TotalWebsitesCount = 500000

var websiteList [TotalWebsitesCount]string
var websiteCounter int

func LoadWebsites(fileToLoad string) {
    //csvFile, _ := os.Open("/tmp/url_list.csv")
    defaultFile := "/tmp/top-1m.csv"
    if fileToLoad == ""{
    	fileToLoad = defaultFile
    }
    csvFile, _ := os.Open(fileToLoad)
    reader := csv.NewReader(bufio.NewReader(csvFile))
    
    for {
        line, error := reader.Read()
        if error == io.EOF {
            break
        } else if error != nil {
            fmt.Println("Error : ",error)
        }
       	for i, website := range line {
       		if i == TotalWebsitesCount {
       			break
       		}
       		websiteList[i] = website
		   	//fmt.Printf("%d: %v\n", i, websiteList[i])		
		}
        
    }
    //logging.Logger.Infof("Websites loaded : ",websiteList)
}

func GetWebsite() string{
	toReturn := websiteList[websiteCounter]
	websiteCounter+=1
	return toReturn
}

func GetWebsiteList() [TotalWebsitesCount]string{
	return websiteList
}
