package cmd

import (
	"fmt"
	"github.com/rad12000/list-agent/internal/config"
	"github.com/rad12000/list-agent/internal/zillow"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// zillowCmd represents the zillow command
var (
	runData = zillow.RunData{
		FilterState: zillow.FilterState{
			SortSelection: zillow.Filter[string]{
				Value: "globalrelevanceex",
			},
		},
	}

	zillowCmd = &cobra.Command{
		Use:   "zillow",
		Short: "Search zillow for listings based on a set of filters",
		Long:  `This command allows searching for zillow listings within a given geographic area`,
		Run: func(cmd *cobra.Command, args []string) {
			runData.FilterState.IsApartmentOrCondo.Value = runData.FilterState.IsApartment.Value && runData.FilterState.IsCondo.Value
			zillow.Run(runData)
		},
	}
)

func init() {
	flags := zillowCmd.Flags()

	flags.StringVar(&runData.FilePath, "file", filepath.Join(config.Directory(), "zillow-results"), "File to store visited zillow listings in")
	cobra.CheckErr(zillowCmd.MarkFlagFilename("file"))

	flags.Float64VarP(&runData.MapBounds.West, "west", "w", 0.0, "the western most coordinate in which to constrain search results")
	flags.Float64VarP(&runData.MapBounds.East, "east", "e", 0.0, "the eastern most coordinate in which to constrain search results")
	flags.Float64VarP(&runData.MapBounds.South, "south", "s", 0.0, "the southern most coordinate in which to constrain search results")
	flags.Float64VarP(&runData.MapBounds.North, "north", "n", 0.0, "the northern most coordinate in which to constrain search results")

	flags.Var(NewPtrToValue(&runData.FilterState.Price.Min, PtrToIntValuer()), "min-price", "the minimum price of a listing")
	flags.Var(NewPtrToValue(&runData.FilterState.Price.Max, PtrToIntValuer()), "max-price", "the maximum price of a listing")

	flags.Var(NewPtrToValue(&runData.FilterState.Beds.Min, PtrToIntValuer()), "min-beds", "the minimum beds of a listing")
	flags.Var(NewPtrToValue(&runData.FilterState.Beds.Max, PtrToIntValuer()), "max-beds", "the maximum beds of a listing")

	flags.Var(NewPtrToValue(&runData.FilterState.Baths.Min, PtrToIntValuer()), "min-baths", "the minimum baths of a listing")
	flags.Var(NewPtrToValue(&runData.FilterState.Baths.Max, PtrToIntValuer()), "max-baths", "the maximum baths of a listing")

	flags.Var(NewPtrToValue(&runData.FilterState.HOA.Min, PtrToIntValuer()), "min-hoa", "the minimum hoa of a listing")
	flags.Var(NewPtrToValue(&runData.FilterState.HOA.Max, PtrToIntValuer()), "max-hoa", "the maximum hoa of a listing")

	flags.Var(NewPtrToValueDefault(&runData.DurationBetweenRuns, IntToDurationValuer(time.Minute), 60), "run-interval", "The amount of time, in minutes, to wait between each execution.")
	flags.Var(NewPtrToValueDefault(&runData.DurationBetweenPages, IntToDurationValuer(time.Second), 30), "page-interval", "The amount of time, in seconds, to wait between each page during each execution.")

	flags.BoolVar(&runData.FilterState.IsSingleFamily.Value, "single-family", false, "whether to include single family homes in the search results")
	flags.BoolVar(&runData.FilterState.IsApartment.Value, "apartment", false, "whether to include apartments in the search results")
	flags.BoolVar(&runData.FilterState.IsCondo.Value, "condo", false, "whether to include condos in the search results")
	flags.BoolVar(&runData.FilterState.IsTownhouse.Value, "townhouse", false, "whether to include townhouses in the search results")
	flags.BoolVar(&runData.FilterState.IsManufactured.Value, "manufactured", false, "whether to include manufactured homes in the search results")
	flags.BoolVar(&runData.FilterState.IsLotLand.Value, "lot", false, "whether to include land lots in the search results")

	flags.StringVar(&runData.UserAgent, "user-agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36`, "The user agent to use when sending requests to zillow.")
	flags.StringSliceVarP(&runData.SearchTerms, "search", "q", nil, "Search terms to include in the search query. Raw regex patterns are supported. Each term will be joined with the the other with the regex OR (|) operator.")
}

func IntToDurationValuer(unit time.Duration) Valuer[time.Duration] {
	return ValuerFunc[time.Duration](func(v string) (time.Duration, error) {
		if v == "" {
			return 0 * unit, nil
		}

		factor, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}

		if factor < 0 {
			return 0, fmt.Errorf("invalid value provided: %v. Must not be less than 0", factor)
		}

		return time.Duration(factor) * unit, nil
	})
}

func PtrToIntValuer() Valuer[*int] {
	return ValuerFunc[*int](func(v string) (*int, error) {
		if v == "" {
			return new(int), nil
		}

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
func NewPtrToValueDefault[T any](valuePtr *T, valuer Valuer[T], defaultValue T) *PtrToValue[T] {
	*valuePtr = defaultValue
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
	return fmt.Sprintf("%T", *f.ValuePtr)
}
