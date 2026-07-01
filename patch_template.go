package main

import (
	"io/ioutil"
	"strings"
	"fmt"
)

func main() {
	content, err := ioutil.ReadFile("templates/PO_Template.jrxml")
	if err != nil {
		panic(err)
	}

	str := string(content)

	// Change band height
	str = strings.Replace(str, "<band height=\"168\">", "<band height=\"185\">", 1)

	// Add Expected Delivery
	expectedDeliveryXml := `
			<textField>
				<reportElement x="0" y="55" width="80" height="12" uuid="new-id-exp1"/>
				<textElement><font fontName="Helvetica" size="9"/></textElement>
				<textFieldExpression><![CDATA["Estimasi Kirim"]]></textFieldExpression>
			</textField>
			<textField>
				<reportElement x="80" y="55" width="10" height="12" uuid="new-id-exp2"/>
				<textElement><font fontName="Helvetica" size="9"/></textElement>
				<textFieldExpression><![CDATA[":"]]></textFieldExpression>
			</textField>
			<textField>
				<reportElement x="90" y="55" width="445" height="12" uuid="new-id-exp3"/>
				<textElement><font fontName="Helvetica" size="9"/></textElement>
				<textFieldExpression><![CDATA[$F{expected_delivery} != null && !$F{expected_delivery}.equals("") ? $F{expected_delivery} : "-"]]></textFieldExpression>
			</textField>`
	
	// Inject before the line y=58
	lineXml := `<line>
				<reportElement x="0" y="58" width="535" height="1"`
	newLineXml := expectedDeliveryXml + `
			<line>
				<reportElement x="0" y="73" width="535" height="1"`
	str = strings.Replace(str, lineXml, newLineXml, 1)

	// Shift elements down
	str = strings.Replace(str, `y="65"`, `y="80"`, -1)
	str = strings.Replace(str, `y="80"`, `y="95"`, -1)
	str = strings.Replace(str, `y="95"`, `y="110"`, -1)
	str = strings.Replace(str, `y="110"`, `y="125"`, -1)
	str = strings.Replace(str, `y="125"`, `y="140"`, -1)
	
	// The line at 143 to 158
	str = strings.Replace(str, `y="143"`, `y="158"`, -1)

	// Add Store Address
	storeAddrXml := `
			<textField>
				<reportElement x="270" y="110" width="265" height="30" uuid="new-id-straddr"/>
				<textElement><font fontName="Helvetica" size="8"/></textElement>
				<textFieldExpression><![CDATA[$F{store_address} != null ? $F{store_address} : ""]]></textFieldExpression>
			</textField>`
	
	// Inject after store name at y=80 (now 95)
	storeNameXml := `<textField>
				<reportElement x="270" y="95" width="265" height="12" uuid="1e78a809-1aa5-4fd9-9797-97f03ac541dd"/>
				<textElement>
					<font fontName="Helvetica" size="9"/>
				</textElement>
				<textFieldExpression><![CDATA[$F{store}]]></textFieldExpression>
			</textField>`
	
	newStoreNameXml := storeNameXml + storeAddrXml
	str = strings.Replace(str, storeNameXml, newStoreNameXml, 1)

	// Fix alignment in detail band
	str = strings.Replace(str, `textAlignment="Center" verticalAlignment="Bottom"`, `textAlignment="Left" verticalAlignment="Middle"`, 1) // For namaItem
	
	// Also adjust Qty, Sat, Harga, Total to be Middle vertically
	str = strings.Replace(str, `verticalAlignment="Bottom"`, `verticalAlignment="Middle"`, -1)
	str = strings.Replace(str, `verticalAlignment="Top"`, `verticalAlignment="Middle"`, -1)

	ioutil.WriteFile("templates/PO_Template.jrxml", []byte(str), 0644)
	fmt.Println("Template patched successfully")
}
