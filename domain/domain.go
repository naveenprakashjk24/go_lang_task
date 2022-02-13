package taskdata

import (
	"encoding/json"
	"log"
	"sort"
)

type InputData struct {
	Ev     string `json:"ev"`
	Et     string `json:"et"`
	Id     string `json:"id"`
	Uid    string `json:"uid"`
	Mid    string `json:"mid"`
	T      string `json:"t"`
	P      string `json:"p"`
	L      string `json:"l"`
	Sc     string `json:"sc"`
	Atrk1  string `json:"atrk1"`
	Atrv1  string `json:"atrv1"`
	Atrt1  string `json:"atrt1"`
	Atrk2  string `json:"atrk2"`
	Atrv2  string `json:"atrv2"`
	Atrt2  string `json:"atrt2"`
	Uatrk1 string `json:"uatrk1"`
	Uatrv1 string `json:"uatrv1"`
	Uatrt1 string `json:"uatrt1"`
	Uatrk2 string `json:"uatrk2"`
	Uatrv2 string `json:"uatrv2"`
	Uatrt2 string `json:"uatrt2"`
	Uatrk3 string `json:"uatrk3"`
	Uatrv3 int    `json:"uatrv3"`
	Uatrt3 string `json:"uatrt3"`
}

type OutputData struct {
	Event            string `json:"even"`
	Event_type       string `json:"event_type"`
	App_id           string `json:"app_id"`
	User_id          string `json:"user_id"`
	Message_id       string `json:"message_id"`
	Page_title       string `json:"page_title"`
	Page_url         string `json:"page_url"`
	Browser_language string `json:"browser_language"`
	Screen_size      string `json:"screen_size"`
	Attributes       struct {
		Form_varient struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"form_varient"`
		Ref struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"ref"`
	} `json:"attributes"`
	Traits struct {
		Name struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"name"`
		Email struct {
			Value string `json:"value"`
			Type  string `json:"type"`
		} `json:"email"`
		Age struct {
			Value int    `json:"value"`
			Type  string `json:"type"`
		} `json:"age"`
	} `json:"traits"`
}

func NewWorker(jobs <-chan InputData, results chan<- OutputData) {
	for input := range jobs {

		r := OutputData{
			Event:            input.Ev,
			Event_type:       input.Et,
			App_id:           input.Id,
			User_id:          input.Uid,
			Message_id:       input.Mid,
			Page_title:       input.T,
			Page_url:         input.P,
			Browser_language: input.L,
			Screen_size:      input.Sc,
			// Attributes {
			// 	Form_varient : InputData.atrk1  {
			// 		Value : inputData.atrv1,
			// 		Type  : InputData.atrt1
			// 	},
			// 	Ref  {
			// 		Value : InputData.atrv2,
			// 		Type  : InputData.atrt2
			// 	}
			// } : "attributes",
			// Traits {
			// 	Name :InputData.uatrk1 {
			// 		Value : InputData.uatrv1,
			// 		Type  : InputData.uatrt1
			// 	},
			// 	Email :InputData.uatrk2 {
			// 		Value : InputData.uatrv2,
			// 		Type  : InputData.uatrt2
			// 	},
			// 	Age : InputData.uatrk3 {
			// 		Value : InputData.uatrv3,
			// 		Type  : InputData.uatrt3
			// 	},
			// }
		}
		results <- r
	}
}
func ProcessJsonInput(inputuserData []byte) []InputData {
	collection := []InputData{}
	var data map[string][]json.RawMessage
	err := json.Unmarshal(inputuserData, &data)
	if err != nil {
		log.Println(err)
		return collection
	}
	for _, thing := range data["userdetails"] {
		collection = addInput(thing, collection)

	}
	return collection
}

func GenerateJsonOutput(output map[int]OutputData) ([]byte, error) {
	//The Output is not sorted by index.
	//We sort it by Index prior to returning the response
	sorted := make([]OutputData, len(output))
	keys := make([]int, len(output))

	for k, _ := range output {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {
		sorted[k] = output[k]
	}
	return json.Marshal(sorted)
}

func addInput(thing json.RawMessage, collection []InputData) []InputData {
	input := InputData{}
	err := json.Unmarshal(thing, &input)

	if err != nil {
		log.Println(err)
	} else {
		if input != *new(InputData) {
			collection = append(collection, input)
		}
	}

	return collection
}
