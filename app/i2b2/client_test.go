package main

import (
	"gopkg.in/dedis/onet.v1"
	"os"
	"testing"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/dedis/onet.v1/log"
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1/app"
	"io"
	"strings"
	"strconv"
)

var clientSecKey abstract.Scalar
var clientPubKey abstract.Point
var local *onet.LocalTest
var el *onet.Roster

var nbrTerms = 50
var nbrFlags = 50


// SETUP / TEARDOWN FUNCTIONS
// ----------------------------------------------------------
func testRemoteSetup() {
	log.SetDebugVisible(1)

	log.LLvl1("***************************************************************************************************")
	os.Remove("pre_compute_multiplications.gob")

	clientSecKey, clientPubKey = lib.GenKey()

	// generate el with group file
	f, err := os.Open("test/group.toml")
	if err != nil {
		log.Error("Error while opening group file", err)
		os.Exit(1)
	}
	el, err = app.ReadGroupToml(f)
	if err != nil {
		log.Error("Error while reading group file", err)
		os.Exit(1)
	}
	if len(el.List) <= 0 {
		log.Error("Empty or invalid group file", err)
		os.Exit(1)
	}
}

func testLocalSetup() {
	log.SetDebugVisible(1)

	log.LLvl1("***************************************************************************************************")
	os.Remove("pre_compute_multiplications.gob")

	clientSecKey, clientPubKey = lib.GenKey()

	local = onet.NewLocalTest()
	_, el, _ = local.GenTree(3, true)
}

func testLocalTeardown() {
	os.Remove("pre_compute_multiplications.gob")
	local.CloseAll()
}

// UTILITY FUNCTIONS
// ----------------------------------------------------------
func getXMLReaderDDTRequest(t *testing.T,  variant int) io.Reader {

	/*
	<unlynx_ddt_request>
	    <id>request ID</id>
	    <enc_values>
		<enc_value>adfw25e4f85as4fas57f=</enc_value>
		<enc_value>ADA5D4D45ESAFD5FDads=</enc_value>
	    </enc_values>
	</unlynx_ddt_request>
	*/

	// enc query terms (encrypted with client public key)
	encDDTTermsSlice := make([]string, 0)
	encDDTTermsXML := ""

	for i:=0; i< nbrTerms; i++ {
		val := (*lib.EncryptInt(el.Aggregate, int64(i))).Serialize()
		encDDTTermsSlice = append(encDDTTermsSlice,val)
		encDDTTermsXML += "<enc_value>" + val + "</enc_value>"
	}

	queryID := "query_ID_XYZf"+strconv.Itoa(variant)

	xmlReader := strings.NewReader(`<unlynx_ddt_request>
						<id>` + queryID + `</id>
						<enc_values>` +
							encDDTTermsXML +
						`</enc_values>
					</unlynx_ddt_request>`)

	log.LLvl1("Generated DDTRequest XML:", xmlReader)

	return xmlReader
}

func getXMLReaderDDTRequestV2(t *testing.T,  variant int) io.Reader {

	/*
	<unlynx_ddt_request>
	    <id>request ID</id>
	    <enc_values>
		<enc_value>adfw25e4f85as4fas57f=</enc_value>
		<enc_value>ADA5D4D45ESAFD5FDads=</enc_value>
	    </enc_values>
	</unlynx_ddt_request>
	*/

	// enc query terms (encrypted with client public key)
	encDDTTermsSlice := make([]string, 0)
	encDDTTermsXML := ""

	for i:=0; i< nbrTerms; i++ {
		val := (*lib.EncryptInt(el.Aggregate, int64(i))).Serialize()
		encDDTTermsSlice = append(encDDTTermsSlice,val)
		encDDTTermsXML += "<enc_value>" + val + "</enc_value>"
	}

	queryID := "query_ID_XYZf"+strconv.Itoa(variant)

	var stringBuf bytes.Buffer

	stringBuf.WriteString(`<unlynx_ddt_request>
				<id>` + queryID + `</id>
				<enc_values>` + encDDTTermsXML + `</enc_values>
			       </unlynx_ddt_request>`)

	log.LLvl1("Generated DDTRequest XML v2:", stringBuf.String())
	return strings.NewReader(stringBuf.String())
}

func getXMLReaderAggRequest(t *testing.T, variant int) io.Reader {

	/*
	<unlynx_agg_request>
	    <id>request ID</id>
	    <client_public_key>5D4D45ESAFD5FDads==</client_public_key>

	    <enc_dummy_flags>
		<enc_dummy_flag>adfw25e4f85as4fas57f=</enc_dummy_flag>
		<enc_dummy_flag>ADA5D4D45ESAFD5FDads=</enc_dummy_flag>
	    </enc_dummy_flags>
	</unlynx_agg_request>
	*/

	// client public key serialization
	clientPubKeyB64, err := lib.SerializeElement(clientPubKey)
	assert.True(t, err == nil)

	// enc query terms (encrypted with client public key)
	encFlagsSlice := make([]string, 0)
	encFlagsXML := ""

	for i:=0; i< nbrFlags; i++ {
		val := (*lib.EncryptInt(el.Aggregate, int64(1))).Serialize()
		encFlagsSlice = append(encFlagsSlice,val)
		encFlagsXML += "<enc_dummy_flag>" + val + "</enc_dummy_flag>"
	}

	queryID := "query_ID_XYZf"+strconv.Itoa(variant)

	xmlReader := strings.NewReader(`<unlynx_agg_request>
						<id>` + queryID + `</id>
						<client_public_key>` + clientPubKeyB64 + `</client_public_key>
						<enc_dummy_flags>` +
							encFlagsXML +
						`</enc_dummy_flags>
					</unlynx_agg_request>`)

	log.LLvl1("Generated AggRequest XML:", xmlReader)

	return xmlReader
}

func getXMLReaderAggRequestV2(t *testing.T,  variant int) io.Reader {

	/*
	<unlynx_agg_request>
	    <id>request ID</id>
	    <client_public_key>5D4D45ESAFD5FDads==</client_public_key>

	    <enc_dummy_flags>
		<enc_dummy_flag>adfw25e4f85as4fas57f=</enc_dummy_flag>
		<enc_dummy_flag>ADA5D4D45ESAFD5FDads=</enc_dummy_flag>
	    </enc_dummy_flags>
	</unlynx_agg_request>
	*/

	// client public key serialization
	clientPubKeyB64, err := lib.SerializeElement(clientPubKey)
	assert.True(t, err == nil)

	// enc query terms (encrypted with client public key)
	encFlagsSlice := make([]string, 0)
	encFlagsXML := ""

	for i:=0; i< nbrFlags; i++ {
		val := (*lib.EncryptInt(el.Aggregate, int64(1))).Serialize()
		encFlagsSlice = append(encFlagsSlice,val)
		encFlagsXML += "<enc_dummy_flag>" + val + "</enc_dummy_flag>"
	}

	queryID := "query_ID_XYZf"+strconv.Itoa(variant)

	var stringBuf bytes.Buffer

	stringBuf.WriteString(`<unlynx_agg_request>
					<id>` + queryID + `</id>
					<client_public_key>` + clientPubKeyB64 + `</client_public_key>
					<enc_dummy_flags>` + encFlagsXML + `</enc_dummy_flags>
			       </unlynx_agg_request>`)

	log.LLvl1("Generated AggRequest XML v2:", stringBuf.String())
	return strings.NewReader(stringBuf.String())
}

func parseDTTResponse(t *testing.T, xmlString string) lib.XMLMedCoDTTResponse {
	parsed_xml := lib.XMLMedCoDTTResponse{}

	err := xml.Unmarshal([]byte(xmlString), &parsed_xml)
	assert.Equal(t, err, nil)

	return parsed_xml
}

func parseAggResponse(t *testing.T, xmlString string) lib.XMLMedCoAggResponse {
	parsed_xml := lib.XMLMedCoAggResponse{}
	err := xml.Unmarshal([]byte(xmlString), &parsed_xml)
	assert.Equal(t, err, nil)

	return parsed_xml
}

// DDT TEST FUNCTIONS
// ----------------------------------------------------------
func TestMedcoDDTRequest(t *testing.T) {
	testLocalSetup()

	// Start queriers (3 nodes)
	wg := lib.StartParallelize(2)
	var writer, writer1, writer2 bytes.Buffer

	go func() {
		defer wg.Done()
		err1 := unlynxDDTRequest(getXMLReaderDDTRequest(t, 1), &writer1, el, 1, false)
		assert.True(t, err1 == nil)
	}()
	go func() {
		defer wg.Done()
		err2 := unlynxDDTRequest(getXMLReaderDDTRequest(t, 2), &writer2, el, 2, false)
		assert.True(t, err2 == nil)
	}()
	err := unlynxDDTRequest(getXMLReaderDDTRequest(t, 0), &writer, el, 0, false)
	assert.True(t, err == nil)
	lib.EndParallelize(wg)

	// Check results
	finalResponses :=  make([]lib.XMLMedCoDTTResponse,0)

	finalResponses = append(finalResponses, parseDTTResponse(t, writer.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer1.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer2.String()))

	for i, response := range finalResponses {
		assert.True(t, response.Error == "")
		assert.Equal(t, len(response.TaggedValues),  nbrTerms, "(" + string(i) + ") The number of tags is different from the number of initial terms")

		for _, el := range response.TaggedValues {

			for j:=i+1; j<len(finalResponses); j++ {
				assert.NotContains(t, finalResponses[j].TaggedValues, el, "There are tags that are the same among nodes")
			}
		}

	}

	testLocalTeardown()
}

func TestMedCoDDTRequestV2(t *testing.T) {
	testLocalSetup()

	// Start queriers (3 nodes)
	wg := lib.StartParallelize(2)
	var writer, writer1, writer2 bytes.Buffer

	go func() {
		defer wg.Done()
		err1 := unlynxDDTRequest(getXMLReaderDDTRequestV2(t, 1), &writer1, el, 1, false)
		assert.True(t, err1 == nil)
	}()
	go func() {
		defer wg.Done()
		err2 := unlynxDDTRequest(getXMLReaderDDTRequestV2(t, 2), &writer2, el, 2, false)
		assert.True(t, err2 == nil)
	}()
	err := unlynxDDTRequest(getXMLReaderDDTRequestV2(t, 0), &writer, el, 0, false)
	assert.True(t, err == nil)
	lib.EndParallelize(wg)

	// Check results
	finalResponses :=  make([]lib.XMLMedCoDTTResponse,0)

	finalResponses = append(finalResponses, parseDTTResponse(t, writer.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer1.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer2.String()))

	for i, response := range finalResponses {
		assert.True(t, response.Error == "")
		assert.Equal(t, len(response.TaggedValues),  nbrTerms, "(" + string(i) + ") The number of tags is different from the number of initial terms")

		for _, el := range response.TaggedValues {

			for j:=i+1; j<len(finalResponses); j++ {
				assert.NotContains(t, finalResponses[j].TaggedValues, el, "There are tags that are the same among nodes")
			}
		}

	}
	testLocalTeardown()
}

func TestUnlynxQueryRemote(t *testing.T) {
	testRemoteSetup()

	// start queries
	wg := lib.StartParallelize(2)
	var writer, writer1, writer2 bytes.Buffer

	go func() {
		defer wg.Done()
		err1 := unlynxDDTRequest(getXMLReaderDDTRequest(t, 1), &writer1, el, 1, false)
		assert.True(t, err1 == nil)
	}()
	go func() {
		defer wg.Done()
		err2 := unlynxDDTRequest(getXMLReaderDDTRequest(t, 2), &writer2, el, 2, false)
		assert.True(t, err2 == nil)
	}()

	err := unlynxDDTRequest(getXMLReaderDDTRequest(t, 0), &writer, el, 0, false)
	assert.True(t, err == nil)
	lib.EndParallelize(wg)

	// Check results
	finalResponses :=  make([]lib.XMLMedCoDTTResponse,0)

	finalResponses = append(finalResponses, parseDTTResponse(t, writer.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer1.String()))
	finalResponses = append(finalResponses, parseDTTResponse(t, writer2.String()))

	for i, response := range finalResponses {
		assert.True(t, response.Error == "")
		assert.Equal(t, len(response.TaggedValues),  nbrTerms, "(" + string(i) + ") The number of tags is different from the number of initial terms")

		for _, el := range response.TaggedValues {

			for j:=i+1; j<len(finalResponses); j++ {
				assert.Contains(t, finalResponses[j].TaggedValues, el, "There are tags that are the same among nodes")
			}


		}

	}
}

// AGG TEST FUNCTIONS
// ----------------------------------------------------------
