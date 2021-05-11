package delivery

/*var Value bool

func (th TrackHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error

	defer func() {

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = handlers.WriteJsonToResponse(writer, err.Error())
		}
		if Value == false {
			writer.WriteHeader(http.StatusBadRequest)
			err = handlers.WriteJsonToResponse(writer, "Bad request")
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}

	}()
	vars := mux.Vars(req)

	var trackV = vars["track"]
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks := th.Repo.GetTracksRedis(trackV, artistV)
	if tracks != nil {
		if ElasticGet(tracks) != true {
			err = ElasticAdd(tracks)
			if err != nil {
				log.Println("error elastic add")
			}
		}
	}
	if tracks != nil {
		err = handlers.WriteJsonToResponse(writer, tracks)
		Value = true
		return
	}

	tracks, err = th.Repo.GetTracks(trackV, artistV)
	bytes, err := json.Marshal(tracks)
	if err == nil && tracks != nil {
		th.Repo.Redis.Set(ctx, "Track:"+trackV+"_Artist:"+artistV, bytes, 20*time.Minute)
	}
	re, err := api.TrackSearchReq(trackV, artistV)
	if err != nil {
		fmt.Println(writer, err.Error())
	}
	if re == nil {
		Value = false
		return
	}
	if tracks != nil {
		if tracks[0].Album == "" || tracks[0].Name == "" || tracks[0].Artist == "" {
			Value = false
		}
	}
	if re != nil {
		go func() {
			err = th.Repo.SetTracks(*re)
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}()
	}
	result := structConv(re)
	err = handlers.WriteJsonToResponse(writer, result)
	Value = true
	return
}
*/
