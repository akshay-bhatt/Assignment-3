package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

type UberPrices struct {
	Prices []struct {
		Product_ID       string  `json:"product_id"`
		Currency_Code    string  `json:"currency_code"`
		Display_Name     string  `json:"display_name"`
		Estimate         string  `json:"estimate"`
		Low_Estimate     int     `json:"low_estimate"`
		High_Estimate    int     `json:"high_estimate"`
		Surge_Multiplier int     `json:"surge_multiplier"`
		Duration         int     `json:"duration"`
		Distance         float64 `json:"distance"`
	} `json:"prices"`
}

type Res1 struct {
	Trip_Id          string   `json:"id" bson:"_id"`
	Status           string   `json:"status" bson:"status"`
	StartPt          string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	Bestlocation_ids []string `json:"Bestlocation_ids" bson:"Bestlocation_ids"`
	Costs            int      `json:"total_uber_costs" bson:"total_uber_costs"`
	Duration         int      `json:"total_uber_duration" bson:"total_uber_duration"`
	Distance         float64  `json: "total_distance" bson:"total_distance"`
}

type Timer_Est struct {
	Times []struct {
		Localized_Display_Name string `json:"localized_display_name"`
		Estimate               int    `json:"estimate"`
		Display_Name           string `json:"display_name"`
		Product_ID             string `json:"product_id"`
	} `json:"times"`
}

type Request struct {
	StartPt      string   `json:"StartPt"`
	Location_ids []string `json:"Location_ids"`
}

type Res2 struct {
	Trip_Id          string   `json:"id" bson:"_id"`
	Status           string   `json:"status" bson:"status"`
	StartPt          string   `json:"starting_from_location_id" bson:"starting_from_location_id"`
	NextPt           string   `json:"next_destination_location_id" bson:"next_destination_location_id"`
	Bestlocation_ids []string `json:"Bestlocation_ids" bson:"Bestlocation_ids"`
	Costs            int      `json:"total_uber_costs" bson:"total_uber_costs"`
	Duration         int      `json:"total_uber_duration" bson:"total_uber_duration"`
	Distance         float64  `json: "total_distance" bson:"total_distance"`
	wait_time_eta    int      `json: "uber_wait_time_eta" bson:"uber_wait_time_eta"`
}

type ResMongLab struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Name       string        `json:"name" bson:"name"`
	Address    string        `json:"address" bson:"address" `
	City       string        `json:"city"  bson:"city"`
	State      string        `json:"state"  bson:"state"`
	ZipCode    string        `json:"zip"  bson:"zip" `
	Coordinate struct {
		Lat float64 `json:"lat"   bson:"lat"`
		Lng float64 `json:"lng"   bson:"lng"`
	} `json:"coordinate" bson:"coordinate"`
}

func getSession() *mgo.Session {

	s, err := mgo.Dial("mongodb://user123:pass12345@ds041934.mongolab.com:41934/users")
	if err != nil {
		panic(err)
	}
	return s
}

type UserC struct {
	session *mgo.Session
}

func NewUser(s *mgo.Session) *UserC {
	return &UserC{s}
}

type ShortPath struct {
	GeoWay []struct {
		Geocoder_Status string   `json:"geocoder_status"`
		Place_ID        string   `json:"place_id"`
		Types           []string `json:"types"`
	} `json:"geocoded_waypoints"`
	Routes []struct {
		Wp_Order []int `json:"waypoint_order"`
	} `json:"route"`
	Status string `json:"status"`
}

var Index_Put int

func C_Loc(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	nuser := NewUser(getSession())
	var Ubp UberPrices
	var res Res1
	var req Request

	temp1 := ResMongLab{}

	json.NewDecoder(r.Body).Decode(&req)
	json.NewDecoder(r.Body).Decode(&res)

	route := append(res.Bestlocation_ids, req.StartPt)

	for _, aa := range req.Location_ids {

		route = append(route, aa)
	}

	//	fmt.Println("fk fk fk fk fk")

	var Costs int
	var Duration int
	var Distance float64

	Adj_mat := make([][]int, len(route))
	Dist_mat := make([][]float64, len(route))
	Dur_mat := make([][]int, len(route))
	for i := 0; i < len(route); i++ {

		/*for i := 0; i < len(); i++ {

			Adj_mat[i] = make([]int, len(route))
			Dist_mat[i] = make([]float64, len(route))
			Dur_mat[i] = make([]int, len(route))
			url1 := "start_latitude=" + strconv.FormatFloat(temp1.Coordinate.Lat, 'f', -1, 64) + "&start_longitude=" + strconv.FormatFloat(temp1.Coordinate.Lng, 'f', -1, 64)
			objid = bson.ObjectIdHex(route[j])

			if err := nuser.session.DB("users").C("hello").FindId(objid).One(&temp1); err != nil {
				rw.WriteHeader(404)
				return
			}

		}
		*/
		Adj_mat[i] = make([]int, len(route))
		Dist_mat[i] = make([]float64, len(route))
		Dur_mat[i] = make([]int, len(route))

		for j := 0; j < len(route); j++ {

			objid := bson.ObjectIdHex(route[i])
			if err := nuser.session.DB("users").C("hello").FindId(objid).One(&temp1); err != nil {
				rw.WriteHeader(404)
				return
			}
			url1 := "start_latitude=" + strconv.FormatFloat(temp1.Coordinate.Lat, 'f', -1, 64) + "&start_longitude=" + strconv.FormatFloat(temp1.Coordinate.Lng, 'f', -1, 64)
			objid = bson.ObjectIdHex(route[j])

			if err := nuser.session.DB("users").C("hello").FindId(objid).One(&temp1); err != nil {
				rw.WriteHeader(404)
				return
			}

			url1 = url1 + "&end_latitude=" + strconv.FormatFloat(temp1.Coordinate.Lat, 'f', -1, 64) + "&end_longitude=" + strconv.FormatFloat(temp1.Coordinate.Lng, 'f', -1, 64)
			url1 = url1 + "&access_token=<NO ACCESS.>"

			Url := "https://sandbox-api.uber.com/v1/estimates/price?" + url1

			res, _ := http.Get(Url)
			data, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			_ = json.Unmarshal(data, &Ubp)
			for _, aa := range Ubp.Prices {
				if aa.Display_Name == "uberX" {
					Costs = aa.High_Estimate
					//	fmt.Println("costs")
					Duration = aa.Duration
					Distance = aa.Distance
				}

			}

			//	fmt.Println("Cost of travelling from location  ", i, " to location  ", j, "is ", Costs)

			Adj_mat[i][j] = Costs
			Dist_mat[i][j] = Distance
			Dur_mat[i][j] = Duration

		}
	}

	for i := 0; i < len(Adj_mat); i++ {
		for j := 0; j < len(Adj_mat); j++ {
			Adj_mat[i][i] = 0
			if Adj_mat[i][j] != Adj_mat[j][i] {
				Adj_mat[j][i] = Adj_mat[i][j]
			}
		}
	}

	wayPoints := make([]int, 0)
	wayPoints = append(wayPoints, 0)
	visit := make([]int, len(Adj_mat))
	visit[0] = 1
	Dst := 0
	flag1 := false // flag for minimum
	val := math.MaxInt64
	i := 0
	j := 0
	count := 0
	//	fmt.Println(wayPoints)
	for count < len(Adj_mat)-1 {
		val = math.MaxInt64
		for j < len(Adj_mat) {

			if Adj_mat[i][j] > 1 && visit[j] == 0 {

				if val > Adj_mat[i][j] {

					val = Adj_mat[i][j]
					Dst = j
					flag1 = true

				}

			}

			j++

		}

		if flag1 == true {
			i = Dst
			//	fmt.Println(i)
			visit[Dst] = 1
			flag1 = false
		}
		//	fmt.Println(visit[])
		wayPoints = append(wayPoints, i)
		count++
		//	fmt.Println(count)
		j = 0
	}

	wayPoints = append(wayPoints, 0)
	//	fmt.Println("Shortest route from source to destination is  ", wayPoints)

	Costs = 0
	for a := 0; a < len(wayPoints)-1; a++ {
		//	fmt.Println(Adj_mat[i][j])
		//	fmt.Println("fk fk fk fk fk")
		i := wayPoints[a]
		j := wayPoints[a+1]

		Costs += Adj_mat[i][j]
		//	fmt.Println("Cost for journey")
		//	fmt.Println(Costs)

		//	fmt.Println("fk fk fk fk fk")
		//	fmt.Println("Distance for journey")

		//	fmt.Println(Dist_mat[i][j])

		Distance += Dist_mat[i][j]

		//		fmt.Println("time taken for journey")
		//		fmt.Println(Dur_mat[i][j])
		Duration += Dur_mat[i][j]

	}

	//	fmt.Println("Best suited path way : ", len(req.Location_ids), len(wayPoints))

	for i := 1; i <= len(wayPoints)-2; i++ {

		//	fmt.Println("fk fk fk fk fk")
		j := wayPoints[i]
		res.Bestlocation_ids = append(res.Bestlocation_ids, route[j])
	}

	fmt.Println("final costs  is : ", Costs)
	fmt.Println("final distance is : ", Distance)
	fmt.Println("Total time taken is : ", Duration)

	x := nuser.session.DB("users").C("hello")
	ans_temp := []ResMongLab{}
	_ = x.Find(nil).All(&ans_temp)

	if len(ans_temp) == 0 {
		res.Trip_Id = strconv.Itoa(12345)
	} else {
		res.Trip_Id = strconv.Itoa(12345 + len(ans_temp))
	}

	res.Status = "Planning"
	res.StartPt = req.StartPt
	res.Costs = Costs
	res.Duration = Duration
	res.Distance = Distance

	nuser.session.DB("users").C("hello").Insert(res)
	UJ, _ := json.Marshal(res)
	fmt.Fprintf(rw, "%s", UJ)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)

	//fmt.Println("End")
}

/*
func getLoc(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	uc := NewUserController(getSession())
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)
	v := Response{}
	//users  //Collection name
	if err := uc.session.DB("users").C("hello").FindId(oid).One(&v); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(v)
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)

}
*/
func put_Location(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	nuser := NewUser(getSession())
	id := p.ByName("id")

	var req Res1
	var mmm_res2 Res2
	var wtf Timer_Est
	//var temp1 Res1
	var RML ResMongLab
	var wait_time_eta int

	if err := nuser.session.DB("users").C("hello").FindId(id).One(&req); err != nil {
		rw.WriteHeader(404)
		return
	}

	if Index_Put < len(req.Bestlocation_ids) {

		bestid := req.Bestlocation_ids[Index_Put]

		if !bson.IsObjectIdHex(bestid) {
			rw.WriteHeader(404)
			return
		}

		objid := bson.ObjectIdHex(bestid)

		//	fmt.Println("REquesting", bestid)

		if err := nuser.session.DB("users").C("hello").FindId(objid).One(&RML); err != nil {
			rw.WriteHeader(404)
			return
		}
		//	fmt.Println("ResMongLab is ", RML)

		//	buffer.WriteString(strconv.FormatFloat(startLat, 'g', -1, 64))
		//	buffer.WriteString("&start_longitude=")
		//	buffer.WriteString(strconv.FormatFloat(endlat, 'g', -1, 64))
		//	buffer.WriteString("&end_longitude=")
		//	buffer.WriteString(strconv.FormatFloat(endLng, 'g', -1, 64))

		Url := "https://sandbox-api.uber.com/v1/estimates/time?start_latitude="
		Url = Url + strconv.FormatFloat(RML.Coordinate.Lat, 'f', -1, 64)
		//	buffer.WriteString(strconv.FormatFloat(startLng, 'g', -1, 64))
		//	buffer.WriteString("&end_latitude=")
		Url = Url + "&start_longitude="
		Url = Url + strconv.FormatFloat(RML.Coordinate.Lng, 'f', -1, 64)
		Url = Url + "&access_token=<NOT DISCLOSED FOR SECURITY REASON.>"

		//	buffer.WriteString(strconv.FormatFloat(endLng, 'g', -1, 64))
		//	buffer.WriteString("&end_latitude=")
		//	fmt.Println(Url)
		res, _ := http.Get(Url)
		data, _ := ioutil.ReadAll(res.Body)
		//	fmt.Println(res.Body)
		res.Body.Close()
		_ = json.Unmarshal(data, &wtf)
		//	fmt.Println(wtf)
		for _, aa := range wtf.Times {
			if aa.Localized_Display_Name == "uberX" {
				wait_time_eta = aa.Estimate
			}

		}
		var temp1 Res1

		json.NewDecoder(r.Body).Decode(&temp1)

		//	fmt.Println("Res1 Id found", temp1)

		mmm_res2.Trip_Id = req.Trip_Id
		mmm_res2.StartPt = req.StartPt
		mmm_res2.Bestlocation_ids = req.Bestlocation_ids
		mmm_res2.Costs = req.Costs
		mmm_res2.Duration = req.Duration
		mmm_res2.Distance = req.Distance

		mmm_res2.Status = temp1.Status
		mmm_res2.wait_time_eta = wait_time_eta
		mmm_res2.NextPt = req.Bestlocation_ids[Index_Put]
		//	fmt.Println("Updated Res1 is ", mmm_res2)
		Index_Put++
		if err := nuser.session.DB("users").C("hello").Update(req, mmm_res2); err != nil {
			rw.WriteHeader(404)
			return

		}
		uj, _ := json.Marshal(req)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, "%s", uj)

	} else if Index_Put >= len(req.Bestlocation_ids) {
		Msg := "You have reached the destination"
		Index_Put = 0
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, "%s", Msg)
	}
}

func get_Location(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	nuser := NewUser(getSession())
	id := p.ByName("id")
	var req Res1

	if err := nuser.session.DB("users").C("hello").FindId(id).One(&req); err != nil {
		rw.WriteHeader(404)
		return
	}
	//fmt.Println(req)
	Uj, _ := json.Marshal(req)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	fmt.Fprintf(rw, "%s", Uj)

}
func main() {

	hi := httprouter.New()

	hi.GET("/locations/:id", get_Location)

	hi.POST("/locations", C_Loc)

	hi.PUT("/locations/:id", put_Location)

	http.ListenAndServe("localhost:6666", hi)

}
