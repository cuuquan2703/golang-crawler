package service_test

import (
	"reflect"
	"testing"
	"webcrawler/service"
)

func TestCheckValidTag(t *testing.T) {
	input := "<b>"
	expected := true

	matched, err := service.CheckValidTag(input)

	if !reflect.DeepEqual(matched, expected) {
		t.Errorf("Not working. Expected: %v, Actual: %v", expected, matched)
	}
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestCheckUnValidTag(t *testing.T) {
	input := "/s/s"
	expected := false

	matched, _ := service.CheckValidTag(input)

	if !reflect.DeepEqual(matched, expected) {
		t.Errorf("Not working. Expected: %v, Actual: %v", expected, matched)
	}
}

func TestCheckValidURL(t *testing.T) {
	input := []string{
		"https://vnexpress.net/vff-cham-dut-hop-dong-voi-hlv-troussier-4725844.html",
		"https://vnexpress.net/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html",
		"https://vnexpress.net/tau-hang-dam-sap-cau-o-my-nhieu-xe-roi-xuong-song-4726802.html",
		"https://vnexpress.net/4-nghi-pham-tan-cong-nha-hat-nga-trinh-dien-toa-an-4726109.html",
	}
	// expected := true
	err := service.CheckValidURL(input)
	if err != nil {
		t.Errorf("Error")
	}

}

func TestCheckValidImg(t *testing.T) {
	input := "https://iv1.vnecdn.net/vnexpress/images/web/2024/03/25/4-nghi-pham-tan-cong-nha-hat-nga-trinh-dien-toa-an-1711323350.jpg?w=0\u0026h=0\u0026q=100\u0026dpr=1\u0026fit=crop\u0026s=ry0Ob2YP6G7iCBk3kuRz4A"
	expected := true

	matched, err := service.CheckValidImgURL(input)

	if !reflect.DeepEqual(matched, expected) {
		t.Errorf("Not working. Expected: %v, Actual: %v", expected, matched)
	}
	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestCheckUnValidImg(t *testing.T) {
	input := "type:jpeg/base64, data=asdadad"
	expected := false

	matched, _ := service.CheckValidTag(input)

	if !reflect.DeepEqual(matched, expected) {
		t.Errorf("Not working. Expected: %v, Actual: %v", expected, matched)
	}
}
