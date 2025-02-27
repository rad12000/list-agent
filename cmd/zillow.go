package cmd

import (
	"fmt"
	"github.com/rad12000/list-agent/internal/config"
	"github.com/rad12000/list-agent/internal/zillow"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
	"strings"
)

// zillowCmd represents the zillow command
var (
	zillowStorageFile string
	zillowSearchTerms []string
	zillowUserAgent   string

	zillowBounds zillow.Bounds
	zillowFilter = zillow.FilterState{
		SortSelection: zillow.Filter[string]{
			Value: "globalrelevanceex",
		},
	}

	zillowCmd = &cobra.Command{
		Use:   "zillow",
		Short: "Search zillow for listings based on a set of filters",
		Long:  `This command allows searching for zillow listings within a given geographic area`,
		Run: func(cmd *cobra.Command, args []string) {
			zillowFilter.IsApartmentOrCondo.Value = zillowFilter.IsApartment.Value && zillowFilter.IsCondo.Value

			zillow.Run(zillow.RunData{
				FilePath:    zillowStorageFile,
				MapBounds:   zillowBounds,
				FilterState: zillowFilter,
				SearchTerms: zillowSearchTerms,
				UserAgent:   zillowUserAgent,
			})
		},
	}
)

func init() {
	flags := zillowCmd.Flags()

	flags.StringVar(&zillowStorageFile, "file", filepath.Join(config.Directory(), "zillow-results"), "File to store visited zillow listings in")
	cobra.CheckErr(zillowCmd.MarkFlagFilename("file"))

	flags.Float64VarP(&zillowBounds.West, "west", "w", 0.0, "the western most coordinate in which to constrain search results")
	flags.Float64VarP(&zillowBounds.East, "east", "e", 0.0, "the eastern most coordinate in which to constrain search results")
	flags.Float64VarP(&zillowBounds.South, "south", "s", 0.0, "the southern most coordinate in which to constrain search results")
	flags.Float64VarP(&zillowBounds.North, "north", "n", 0.0, "the northern most coordinate in which to constrain search results")

	flags.Var(NewPtrToValue(&zillowFilter.Price.Min, PtrToIntValuer()), "min-price", "the minimum price of a listing")
	flags.Var(NewPtrToValue(&zillowFilter.Price.Max, PtrToIntValuer()), "max-price", "the maximum price of a listing")

	flags.Var(NewPtrToValue(&zillowFilter.Beds.Min, PtrToIntValuer()), "min-beds", "the minimum beds of a listing")
	flags.Var(NewPtrToValue(&zillowFilter.Beds.Max, PtrToIntValuer()), "max-beds", "the maximum beds of a listing")

	flags.Var(NewPtrToValue(&zillowFilter.Baths.Min, PtrToIntValuer()), "min-baths", "the minimum baths of a listing")
	flags.Var(NewPtrToValue(&zillowFilter.Baths.Max, PtrToIntValuer()), "max-baths", "the maximum baths of a listing")

	flags.Var(NewPtrToValue(&zillowFilter.HOA.Min, PtrToIntValuer()), "min-hoa", "the minimum hoa of a listing")
	flags.Var(NewPtrToValue(&zillowFilter.HOA.Max, PtrToIntValuer()), "max-hoa", "the maximum hoa of a listing")

	flags.BoolVar(&zillowFilter.IsSingleFamily.Value, "single-family", false, "whether to include single family homes in the search results")
	flags.BoolVar(&zillowFilter.IsApartment.Value, "apartment", false, "whether to include apartments in the search results")
	flags.BoolVar(&zillowFilter.IsCondo.Value, "condo", false, "whether to include condos in the search results")
	flags.BoolVar(&zillowFilter.IsTownhouse.Value, "townhouse", false, "whether to include townhouses in the search results")
	flags.BoolVar(&zillowFilter.IsManufactured.Value, "manufactured", false, "whether to include manufactured homes in the search results")
	flags.BoolVar(&zillowFilter.IsLotLand.Value, "lot", false, "whether to include land lots in the search results")

	flags.StringVar(&zillowUserAgent, "user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36`, "The user agent to use when sending requests to zillow.")
	flags.StringSliceVarP(&zillowSearchTerms, "search", "q", nil, "Search terms to include in the search query. Raw regex patterns are supported. Each term will be joined with the the other with the regex OR (|) operator.")
}

func PtrToIntValuer() Valuer[*int] {
	return ValuerFunc[*int](func(v string) (*int, error) {
		result, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		return &result, nil
	})
}

type Valuer[T any] interface {
	Value(v string) (T, error)
}

type ValuerFunc[T any] func(v string) (T, error)

func (v ValuerFunc[T]) Type() string {
	var t T
	return strings.ReplaceAll(fmt.Sprintf("%T", t), "*", "")
}

func (v ValuerFunc[T]) Value(val string) (T, error) {
	return v(val)
}

func NewPtrToValue[T any](valuePtr *T, valuer Valuer[T]) *PtrToValue[T] {
	return &PtrToValue[T]{
		ValuePtr: valuePtr,
		Valuer:   valuer,
	}
}

type PtrToValue[T any] struct {
	ValuePtr *T
	Valuer   Valuer[T]
}

func (f *PtrToValue[T]) String() string {
	return fmt.Sprintf("%+v", *f.ValuePtr)
}

func (f *PtrToValue[T]) Set(s string) error {
	value, err := f.Valuer.Value(s)
	if err != nil {
		return fmt.Errorf("failed to parse value %s: %w", s, err)
	}
	*f.ValuePtr = value
	return nil
}

func (f *PtrToValue[T]) Type() string {
	if v, ok := f.Valuer.(interface{ Type() string }); ok {
		return v.Type()
	}
	return fmt.Sprintf("%T", f)
}

/*
	West:  -112.069650,
	East:  -111.569363,
	South: 40.012749,
	North: 41.127364,
*/

//ADU|basement\s+apartment|\smother\sin(|\s|-)law|\swalk(\s|-)out|\sseparate\sentrance
