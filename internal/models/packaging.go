package models

import (
	"log"

	"github.com/Rhymond/go-money"
	"github.com/bytedance/sonic"
)

// PackagingType is a type for different packagings
type PackagingType uint

const (
	// NoPackaging is for specifying that no packagings were used
	NoPackaging PackagingType = iota

	// BagPackaging is for bag packaging
	BagPackaging

	// BoxPackaging is for box packaging
	BoxPackaging

	// WrapPackaging is for wrap packaging
	WrapPackaging
)

var (
	bagCost  = money.NewFromFloat(5.00, money.RUB)
	boxCost  = money.NewFromFloat(20.00, money.RUB)
	wrapCost = money.NewFromFloat(1.00, money.RUB)
)

const (
	bagMinWeight float64 = 10.00
	boxMinWeight float64 = 30.00
)

const (
	// NoPackagingName is name for no packaging
	NoPackagingName = "none"

	// BagName is name for bag packaging
	BagName = "bag"

	// BoxName is name for box packaging
	BoxName = "box"

	// WrapName is name for wrap packaging
	WrapName = "wrap"
)

// Packaging is an interface that all packagings must implement
type Packaging interface {
	String() string
	GetType() PackagingType
	GetCost() *money.Money
	GetMinWeight() float64
	GetCheckWeight() bool
}

// GetPackaging is a factory for making packagings
func GetPackaging(packaging string) Packaging {
	switch packaging {
	case BagName:
		return newBag()
	case BoxName:
		return newBox()
	case WrapName:
		return newWrap()
	case NoPackagingName:
		return newNone()
	default:
		return nil
	}
}

// BasePackaging is a struct for basic packaging
type BasePackaging struct {
	// Type is a type of this packaging
	Type PackagingType

	// Cost is a cost of this packaging
	Cost money.Money

	// MinWeight is a minimum weight required for this packaging
	MinWeight float64

	// CheckWeight is a flag whether to check weight or not
	CheckWeight bool
}

func (b *BasePackaging) String() string {
	return GetPackagingName(b.Type)
}

// GetType gets Type of the Packaging
func (b *BasePackaging) GetType() PackagingType {
	return b.Type
}

// GetCost gets Cost of the Packaging
func (b *BasePackaging) GetCost() *money.Money {
	return &b.Cost
}

// GetMinWeight gets MinWeight of the Packaging
func (b *BasePackaging) GetMinWeight() float64 {
	return b.MinWeight
}

// GetCheckWeight gets CheckWeight of the Packaging
func (b *BasePackaging) GetCheckWeight() bool {
	return b.CheckWeight
}

// MarshalJSON is used to marshall packaging to json
func (b *BasePackaging) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(b.String())
}

// UnmarshalJSON is used to unmarshall packaging to json
func (b *BasePackaging) UnmarshalJSON(data []byte) error {
	var name string
	if err := sonic.Unmarshal(data, &name); err != nil {
		return err
	}

	packaging := GetPackaging(name)
	if packaging == nil {
		log.Fatal("unknown packaging type: ", name)
	}

	basePackaging, ok := packaging.(*BasePackaging)
	if !ok {
		log.Fatal("invalid packaging type for: ", name)
	}

	*b = *basePackaging

	return nil
}

// GetPackagingName gets Packaging name from PackagingType
func GetPackagingName(packaging PackagingType) string {
	switch packaging {
	case BagPackaging:
		return BagName
	case BoxPackaging:
		return BoxName
	case WrapPackaging:
		return WrapName
	case NoPackaging:
		return NoPackagingName
	default:
		return ""
	}
}

// Bag is a struct of bag packaging
type Bag struct {
	BasePackaging
}

func newBag() Packaging {
	return &Bag{
		BasePackaging{
			Type:        BagPackaging,
			Cost:        *bagCost,
			MinWeight:   bagMinWeight,
			CheckWeight: true,
		},
	}
}

// Box is a struct of box packaging
type Box struct {
	BasePackaging
}

func newBox() Packaging {
	return &Box{
		BasePackaging{
			Type:        BoxPackaging,
			Cost:        *boxCost,
			MinWeight:   boxMinWeight,
			CheckWeight: true,
		},
	}
}

// Wrap is a struct of wrap packaging
type Wrap struct {
	BasePackaging
}

func newWrap() Packaging {
	return &Wrap{
		BasePackaging{
			Type:        WrapPackaging,
			Cost:        *wrapCost,
			MinWeight:   0,
			CheckWeight: false,
		},
	}
}

// None is a struct of no packaging
type None struct {
	BasePackaging
}

func newNone() Packaging {
	return &None{
		BasePackaging{
			Type: NoPackaging,
			Cost: *money.NewFromFloat(0, money.RUB),
		},
	}
}
