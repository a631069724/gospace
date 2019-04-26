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
			outline := arryTostring(outArry)
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
	enc := mahonia.NewEncoder("gbk")

	bufWriter.WriteString(enc.ConvertString("总笔数	应付总金额	手续费总金额	实付总金额 \t\n"))
	bufWriter.Flush()
	bufWriter.WriteString(fmt.Sprintf("%d\t%.2f\t%.2f\t%.2f", convf.count, convf.yingfu, convf.fee, convf.shifu))
	bufWriter.Flush()
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
	str += "\n"
	return
}
