package serviceI2B2_test

import (
	"strconv"
	"github.com/JoaoAndreSa/MedCo/lib"
	"testing"
	"os"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"github.com/JoaoAndreSa/MedCo/services/i2b2"
)

func TestServiceEncGrpAndWhereAttr(t *testing.T) {
	log.LLvl1("***************************************************************************************************")
	os.Remove("pre_compute_multiplications.gob")
	log.SetDebugVisible(2)
	local := onet.NewLocalTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, el, _ := local.GenTree(3, true)
	defer local.CloseAll()

	// Send a request to the service
	client := serviceI2B2.NewMedcoClient(el.List[0], strconv.Itoa(0))
	client1 := serviceI2B2.NewMedcoClient(el.List[1], strconv.Itoa(0))
	client2 := serviceI2B2.NewMedcoClient(el.List[2], strconv.Itoa(0))


	sum := []string{"sum1"}
	count := false
	whereQueryValues := []lib.WhereQueryAttribute{{"w1", *lib.EncryptInt(el.Aggregate, 1)}, {"w2", *lib.EncryptInt(el.Aggregate, 1)}, {"w3", *lib.EncryptInt(el.Aggregate, 1)}} // v1, v3 and v5
	pred := "(v0 == v1 || v2 == v3) && v4 == v5"
	groupBy := []string{}

	nbrDPs := make(map[string]int64)
	//how many data providers for each server
	for _, server := range el.List {
		nbrDPs[server.String()] = 1 // 1 DPs for each server
	}
	data := []lib.ProcessResponse{}
	val := int64(1)

	nbrWhere := 3
	sliceWhere := make(lib.CipherVector, nbrWhere)
	for j := 0 ; j < nbrWhere; j++ {
		sliceWhere[j] = *lib.EncryptInt(el.Aggregate, val)

	}

	sliceWhere1 := make(lib.CipherVector, nbrWhere)
	for j := 0 ; j < nbrWhere; j++ {
		sliceWhere1[j] = *lib.EncryptInt(el.Aggregate, val)

	}

	nbrGrp := 3
	sliceGrp := make(lib.CipherVector, nbrGrp)
	for j := 0 ; j < nbrGrp; j++ {
		sliceGrp[j] = *lib.EncryptInt(el.Aggregate, val)

	}

	nbrAggr := 1
	aggr := make(lib.CipherVector, nbrAggr)
	for j := 0 ; j < nbrAggr; j++ {
		aggr[j] = *lib.EncryptInt(el.Aggregate, val)

	}

	data = append(data, lib.ProcessResponse{WhereEnc:sliceWhere, AggregatingAttributes:aggr}, lib.ProcessResponse{WhereEnc:sliceWhere, AggregatingAttributes:aggr}, lib.ProcessResponse{WhereEnc:sliceWhere1, AggregatingAttributes:aggr})

	wg := lib.StartParallelize(2)
	log.LLvl1("START PARA")
	go func(i int) {
		defer wg.Done()
		log.LLvl1("ICI")
		_, _, _ = client1.SendSurveyDpQuery(el, serviceI2B2.SurveyID("testSurvey"), serviceI2B2.SurveyID(""), nil, nbrDPs, false, false, sum, count, whereQueryValues, pred, groupBy, data, 0)
	}(0)
	go func() {
		defer wg.Done()
		_, _, _ = client2.SendSurveyDpQuery(el, serviceI2B2.SurveyID("testSurvey"), serviceI2B2.SurveyID(""), nil, nbrDPs, false, false, sum, count, whereQueryValues, pred, groupBy, data, 0)
	}()
	_, result, err := client.SendSurveyDpQuery(el, serviceI2B2.SurveyID("testSurvey"), serviceI2B2.SurveyID(""), nil, nbrDPs, false, false, sum, count, whereQueryValues, pred, groupBy, data, 0)

	lib.EndParallelize(wg)

	_= result
	if err != nil {
		t.Fatal("Service did not start.", err)
	}
}