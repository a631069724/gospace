package conver

import (
	"bufio"
	"fmt"
	"io"
	"os"
	logger "skpconvfile/log"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
)

const (
	INDXE_TERMINAL = iota
	INDXE_SEP1
	INDXE_ID
	INDXE_DATE
	INDXE_TIME
	INDXE_AMT
	INDXE_FEE
	INDXE_SETAMT
	INDXE_PERCENT
	INDXE_SEP2
	INDXE_TRANTYPE
	INDXE_CODE
	INDXE_CHANNEL
	INDXE_REMAK1
	INDXE_MERCHANTID
	INDXE_MERCHANTNAME
	INDXE_SETDATE
)

var outIndex = []int{
	INDXE_MERCHANTNAME,
	INDXE_SEP1,
	INDXE_MERCHANTID,
	INDXE_SEP1,
	INDXE_SEP1,
	INDXE_SEP1,
	INDXE_TERMINAL,
	INDXE_ID,
	INDXE_DATE,
	INDXE_AMT,
	INDXE_FEE,
	INDXE_SETAMT,
	INDXE_TRANTYPE,
	INDXE_REMAK1,
}

type ConverFile struct {
	infile  string
	outfile string
	count   int
	yingfu  float64
	fee     float64
	shifu   float64
}

func NewConverFile(in string, out string) Conver {
	return &ConverFile{
		infile:  in,
		outfile: out,
	}
}

func (convf *ConverFile) Conver() error {
	infp, err := os.Open(convf.infile)
	if err != nil {
		return err
	}
	defer infp.Close()
	outfp, err := os.Create(convf.outfile)
	if err != nil {
		return err
	}
	defer outfp.Close()
	bufReader := bufio.NewReader(infp)
	bufWriter := bufio.NewWriter(outfp)

	bufWriter.WriteString(mahonia.NewEncoder("gbk").ConvertString("商户名称;商户地址;商户编号;邮政编码;管辖支行号;流水号;终端号;交易卡号;交易日期;交易金额;回扣;净计金额;卡别;系统参考号\n"))
	bufWriter.Flush()
	for {
		inline, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		} else {

			outArry := strings.Split(string(inline), "\t")
			if len(outArry) < 8 {
				logger.Println("conver buf err:", inline)
				bufWriter.WriteString(string(inline))
				bufWriter.Flush()
				continue
			}
			if strings.HasPrefix(outArry[7], "-") && !strings.HasPrefix(outArry[5], "-") {
				logger.Println("conver line:", string(inline))
				outArry[5] = fmt.Sprintf("-%s", outArry[5])
				outArry[6] = fmt.Sprintf("-%s", outArry[6])
			}
			//outline := arryTostring(outArry)
			outline := arryIndexToString(outArry, outIndex)
			bufWriter.WriteString(outline)
			bufWriter.Flush()
			convf.addCount()
			if amt, err := strconv.ParseFloat(outArry[5], 64); err != nil {
				logger.Println("conver yingfu:", outArry[5], "amt err", err)
			} else {
				convf.addYingfu(amt)
			}

			if amt, err := strconv.ParseFloat(outArry[6], 64); err != nil {
				logger.Println("conver fee:", outArry[6], "amt err", err)
			} else {
				convf.addFee(amt)
			}

			if amt, err := strconv.ParseFloat(outArry[7], 64); err != nil {
				logger.Println("conver shifu:", outArry[7], "amt err", err)
			} else {
				convf.addShifu(amt)
			}

		}
	}
	/*
		enc := mahonia.NewEncoder("gbk")

		bufWriter.WriteString(enc.ConvertString("总笔数	应付总金额	手续费总金额	实付总金额 \t\n"))
		bufWriter.Flush()
		bufWriter.WriteString(fmt.Sprintf("%d\t%.2f\t%.2f\t%.2f", convf.count, convf.yingfu, convf.fee, convf.shifu))
		bufWriter.Flush()*/
	return nil
}

func (convf *ConverFile) addCount() {
	convf.count++
}
func (convf *ConverFile) addYingfu(amt float64) {
	convf.yingfu += amt
}
func (convf *ConverFile) addFee(amt float64) {
	convf.fee += amt
}
func (convf *ConverFile) addShifu(amt float64) {
	convf.shifu += amt
}

func arryTostring(arry []string) (str string) {
	for _, v := range arry {
		str += v + string('\t')
	}
	str += "\r\n"
	return
}

func arryIndexToString(arry []string, arryindex []int) (str string) {
	for i, v := range arryindex {
		if v >= len(arry) {
			continue
		}
		if v == INDXE_TRANTYPE {
			if arry[INDXE_TRANTYPE] == "WEXP" {
				str += "01;"
			} else if arry[INDXE_TRANTYPE] == "WEXR" {
				str += "03;"
			} else {
				str += ";"
			}
			continue
		}
		if i == len(arryindex)-1 {
			str += arry[v]
			break
		}
		str += arry[v] + string(';')
	}
	str += "\n"
	return
}
