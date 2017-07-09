package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"time"
	"strings"
)

type API_KEY struct {
	Key string
	Secret string
	Ex string
}

type Asset struct{
	Ex string
	Base string
	Pair string
	Start int64
}

type Configuration struct {
	Bucket struct {
		Assets []Asset
	}

	Keys []API_KEY

	HashedKey string

	WebServer struct {
		Port int32
		Address string
	}
}

func Config() *Configuration{
	file, _ := os.OpenFile("config.json",os.O_CREATE|os.O_RDWR,os.FileMode(0644))
	configuration := Configuration{}
	data,_ := ioutil.ReadAll(file)
	if len(data) == 0 {
		jsonfile, _ := json.Marshal(configuration)
		file.Write(jsonfile)
	}else {
		json.Unmarshal(data, &configuration)
	}
	file.Close()
	return &configuration
}

func (conf *Configuration)AddTarget(ex,base,pair string,start time.Time){
	if conf.Bucket.Assets == nil {
		conf.Bucket.Assets = []Asset{}
	}

	for _,target := range conf.Bucket.Assets {
		if target.Ex == ex &&
			target.Base == base &&
			target.Pair == pair {
			target.Start = start.UTC().Unix()
			conf.Save()
			return
		}
	}

	conf.Bucket.Assets = append(conf.Bucket.Assets,
		Asset{
			Ex:ex,
			Base:base,
			Pair:pair,
			Start:start.UTC().Unix(),
		})

	conf.Save()
}

func (conf *Configuration)RemoveTarget(ex,base,pair string,start time.Time){
	for i,target := range conf.Bucket.Assets {
		if target.Ex == ex &&
			target.Base == base &&
			target.Pair == pair {
			conf.Bucket.Assets = append(conf.Bucket.Assets [:i], conf.Bucket.Assets [i+1:]...)
			return
		}
	}
}

func (conf *Configuration)Save()  {
	file, _ := os.OpenFile("config.json",os.O_CREATE|os.O_RDWR|os.O_TRUNC,os.FileMode(0644))
	jsonf, _ := json.Marshal(conf)
	//formatting
	str := strings.Replace(string(jsonf),",",",\n",-1)
	str = strings.Replace(str,"{","\n{\n",-1)
	str = strings.Replace(str,"[","\n[\n",-1)
	str = strings.Replace(str,"]","\n]\n",-1)
	str = strings.Replace(str,"}","\n}\n",-1)
	splits := strings.Split(str,"\n")
	count := 0
	final := ""
	for i,item := range splits {
		if item == "" || item == "," {
			continue
		}

		if strings.Contains(item,"}") || strings.Contains(item,"]") {
			count --
		}
		for i:=0;i<count ;  i++{
			final += "  "
		}
		final += item
		if i != len(splits)-2 {
			if splits[i+1] == "," {
				final += ","
			}
		}
		final += "\n"
		if strings.Contains(item,"{") || strings.Contains(item,"[") {
			count ++
		}
	}
	file.Write([]byte(final))
	file.Sync()
	file.Close()
}