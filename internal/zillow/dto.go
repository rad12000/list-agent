package zillow

type rWants struct {
	Cat1 []string `json:"cat1"`
	Cat2 []string `json:"cat2"`
}

type Filter[T any] struct {
	Value T `json:"value"`
}

type Pricing struct {
	Min *int `json:"min"`
	Max *int `json:"max"`
}

type Bounds struct {
	West  float64 `json:"west"`
	East  float64 `json:"east"`
	South float64 `json:"south"`
	North float64 `json:"north"`
}

type FilterState struct {
	SortSelection      Filter[string] `json:"sortSelection"`
	Price              Pricing        `json:"price"`
	MonthlyPayment     Pricing        `json:"monthlyPayment"`
	Beds               Pricing        `json:"beds"`
	Baths              Pricing        `json:"baths"`
	HOA                Pricing        `json:"hoa"`
	IsAllHomes         Filter[bool]   `json:"isAllHomes"`
	IsTownhouse        Filter[bool]   `json:"isTownhouse"`
	IsCondo            Filter[bool]   `json:"isCondo"`
	IsLotLand          Filter[bool]   `json:"isLotLand"`
	IsApartment        Filter[bool]   `json:"isApartment"`
	IsSingleFamily     Filter[bool]   `json:"isSingleFamily"`
	IsManufactured     Filter[bool]   `json:"isManufactured"`
	IsApartmentOrCondo Filter[bool]   `json:"isApartmentOrCondo"`
}

type pagination struct {
	CurrentPage int `json:"currentPage,omitempty"`
}

type searchState struct {
	Pagination           pagination  `json:"pagination"`
	IsMapVisible         bool        `json:"isMapVisible"`
	MapBounds            Bounds      `json:"mapBounds"`
	FilterState          FilterState `json:"filterState"`
	IsEntirePlaceForRent bool        `json:"isEntirePlaceForRent"`
	IsRoomForRent        bool        `json:"isRoomForRent"`
	IsListVisible        bool        `json:"isListVisible"`
	MapZoom              int         `json:"mapZoom"`
}

type request struct {
	SearchQueryState searchState `json:"searchQueryState"`
	Wants            rWants      `json:"wants"`
	RequestId        int         `json:"requestId"`
	IsDebugRequest   bool        `json:"isDebugRequest"`
}

func Copy[T any](v T, fn func(T) T) T {
	return fn(v)
}
