package zillow

type User struct {
	IsLoggedIn                    bool   `json:"isLoggedIn"`
	HasHousingConnectorPermission bool   `json:"hasHousingConnectorPermission"`
	SavedHomesCount               int    `json:"savedHomesCount"`
	PersonalizedSearchTraceID     string `json:"personalizedSearchTraceID"`
	Guid                          string `json:"guid"`
	Zuid                          string `json:"zuid"`
	IsBot                         bool   `json:"isBot"`
	Email                         string `json:"email"`
	DisplayName                   string `json:"displayName"`
}

type Response struct {
	User     User `json:"user"`
	MapState struct {
		CustomRegionPolygonWkt  interface{} `json:"customRegionPolygonWkt"`
		SchoolPolygonWkt        interface{} `json:"schoolPolygonWkt"`
		IsCurrentLocationSearch bool        `json:"isCurrentLocationSearch"`
		UserPosition            struct {
			Lat interface{} `json:"lat"`
			Lon interface{} `json:"lon"`
		} `json:"userPosition"`
	} `json:"mapState"`
	RegionState struct {
		RegionInfo []struct {
			RegionName    string `json:"regionName"`
			IsPointRegion bool   `json:"isPointRegion"`
		} `json:"regionInfo"`
	} `json:"regionState"`
	SearchPageSeoObject struct {
		BaseUrl         string `json:"baseUrl"`
		WindowTitle     string `json:"windowTitle"`
		MetaDescription string `json:"metaDescription"`
	} `json:"searchPageSeoObject"`
	RequestId int `json:"requestId"`
	Cat1      struct {
		SearchResults struct {
			ListResults []struct {
				Zpid                        string  `json:"zpid"`
				Id                          string  `json:"id"`
				RawHomeStatusCd             string  `json:"rawHomeStatusCd"`
				MarketingStatusSimplifiedCd string  `json:"marketingStatusSimplifiedCd"`
				ImgSrc                      string  `json:"imgSrc"`
				HasImage                    bool    `json:"hasImage"`
				DetailUrl                   string  `json:"detailUrl"`
				StatusType                  string  `json:"statusType"`
				StatusText                  string  `json:"statusText"`
				CountryCurrency             string  `json:"countryCurrency"`
				Price                       string  `json:"price"`
				UnformattedPrice            int     `json:"unformattedPrice"`
				Address                     string  `json:"address"`
				AddressStreet               string  `json:"addressStreet"`
				AddressCity                 string  `json:"addressCity"`
				AddressState                string  `json:"addressState"`
				AddressZipcode              string  `json:"addressZipcode"`
				IsUndisclosedAddress        bool    `json:"isUndisclosedAddress"`
				Beds                        int     `json:"beds"`
				Baths                       float64 `json:"baths"`
				Area                        int     `json:"area"`
				LatLong                     struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"latLong"`
				IsZillowOwned bool `json:"isZillowOwned"`
				VariableData  struct {
					Type string `json:"type"`
					Text string `json:"text"`
					Data struct {
						IsRead  interface{} `json:"isRead"`
						IsFresh bool        `json:"isFresh"`
					} `json:"data,omitempty"`
				} `json:"variableData"`
				HdpData struct {
					HomeInfo struct {
						Zpid            int     `json:"zpid"`
						StreetAddress   string  `json:"streetAddress"`
						Zipcode         string  `json:"zipcode"`
						City            string  `json:"city"`
						State           string  `json:"state"`
						Latitude        float64 `json:"latitude"`
						Longitude       float64 `json:"longitude"`
						Price           float64 `json:"price"`
						Bathrooms       float64 `json:"bathrooms"`
						Bedrooms        float64 `json:"bedrooms"`
						LivingArea      float64 `json:"livingArea"`
						HomeType        string  `json:"homeType"`
						HomeStatus      string  `json:"homeStatus"`
						DaysOnZillow    int     `json:"daysOnZillow"`
						IsFeatured      bool    `json:"isFeatured"`
						ShouldHighlight bool    `json:"shouldHighlight"`
						Zestimate       int     `json:"zestimate,omitempty"`
						RentZestimate   int     `json:"rentZestimate,omitempty"`
						ListingSubType  struct {
							IsFSBA bool `json:"is_FSBA"`
						} `json:"listing_sub_type"`
						IsUnmappable            bool    `json:"isUnmappable"`
						IsPreforeclosureAuction bool    `json:"isPreforeclosureAuction"`
						HomeStatusForHDP        string  `json:"homeStatusForHDP"`
						PriceForHDP             float64 `json:"priceForHDP"`
						TimeOnZillow            int64   `json:"timeOnZillow"`
						IsNonOwnerOccupied      bool    `json:"isNonOwnerOccupied"`
						IsPremierBuilder        bool    `json:"isPremierBuilder"`
						IsZillowOwned           bool    `json:"isZillowOwned"`
						Currency                string  `json:"currency"`
						Country                 string  `json:"country"`
						TaxAssessedValue        float64 `json:"taxAssessedValue,omitempty"`
						LotAreaValue            float64 `json:"lotAreaValue"`
						LotAreaUnit             string  `json:"lotAreaUnit"`
						IsShowcaseListing       bool    `json:"isShowcaseListing"`
						DatePriceChanged        int64   `json:"datePriceChanged,omitempty"`
						PriceReduction          string  `json:"priceReduction,omitempty"`
						PriceChange             int     `json:"priceChange,omitempty"`
					} `json:"homeInfo"`
				} `json:"hdpData"`
				IsSaved                    bool   `json:"isSaved"`
				IsUserClaimingOwner        bool   `json:"isUserClaimingOwner"`
				IsUserConfirmedClaim       bool   `json:"isUserConfirmedClaim"`
				Pgapt                      string `json:"pgapt"`
				Sgapt                      string `json:"sgapt"`
				Zestimate                  int    `json:"zestimate,omitempty"`
				ShouldShowZestimateAsPrice bool   `json:"shouldShowZestimateAsPrice"`
				Has3DModel                 bool   `json:"has3DModel"`
				HasVideo                   bool   `json:"hasVideo"`
				IsHomeRec                  bool   `json:"isHomeRec"`
				HasAdditionalAttributions  bool   `json:"hasAdditionalAttributions"`
				IsFeaturedListing          bool   `json:"isFeaturedListing"`
				IsShowcaseListing          bool   `json:"isShowcaseListing"`
				List                       bool   `json:"list"`
				Relaxed                    bool   `json:"relaxed"`
				Info3String                string `json:"info3String"`
				BrokerName                 string `json:"brokerName"`
				CarouselPhotos             []struct {
					Url string `json:"url"`
				} `json:"carouselPhotos"`
			} `json:"listResults"`
			MapResults         []interface{} `json:"mapResults"`
			RelaxedResults     []interface{} `json:"relaxedResults"`
			ResultsHash        string        `json:"resultsHash"`
			HomeRecCount       int           `json:"homeRecCount"`
			ShowForYouCount    int           `json:"showForYouCount"`
			RelaxedResultsHash string        `json:"relaxedResultsHash"`
		} `json:"searchResults"`
		SearchList struct {
			ExpansionDistance  int         `json:"expansionDistance"`
			StaticBaseUrl      interface{} `json:"staticBaseUrl"`
			ZeroResultsFilters interface{} `json:"zeroResultsFilters"`
			Pagination         struct {
				NextUrl string `json:"nextUrl"`
			} `json:"pagination"`
			AdsConfig struct {
				NavAdSlot     string `json:"navAdSlot"`
				DisplayAdSlot string `json:"displayAdSlot"`
				Targets       struct {
					Guid         string `json:"guid"`
					Vers         string `json:"vers"`
					Premieragent string `json:"premieragent"`
					State        string `json:"state"`
					Dma          string `json:"dma"`
					Cnty         string `json:"cnty"`
					City         string `json:"city"`
					Zip          string `json:"zip"`
					Mlat         string `json:"mlat"`
					Mlong        string `json:"mlong"`
					Bd           string `json:"bd"`
					Ba           string `json:"ba"`
					Prange       string `json:"prange"`
					Listtp       string `json:"listtp"`
					Searchtp     string `json:"searchtp"`
				} `json:"targets"`
				NeedsUpdate bool `json:"needsUpdate"`
			} `json:"adsConfig"`
			TotalResultCount        int         `json:"totalResultCount"`
			ResultsPerPage          int         `json:"resultsPerPage"`
			TotalPages              int         `json:"totalPages"`
			LimitSearchResultsCount interface{} `json:"limitSearchResultsCount"`
			ListResultsTitle        string      `json:"listResultsTitle"`
			ResultContexts          []struct {
				Ssid         int    `json:"ssid"`
				Context      string `json:"context"`
				ContextImage string `json:"contextImage"`
			} `json:"resultContexts"`
			PageRules   string `json:"pageRules"`
			ShareConfig struct {
				CaptchaKey string `json:"captchaKey"`
				CsrfToken  string `json:"csrfToken"`
			} `json:"shareConfig"`
		} `json:"searchList"`
	} `json:"cat1"`
	CategoryTotals struct {
		Cat1 struct {
			TotalResultCount int `json:"totalResultCount"`
		} `json:"cat1"`
		Cat2 struct {
			TotalResultCount int `json:"totalResultCount"`
		} `json:"cat2"`
	} `json:"categoryTotals"`
}
