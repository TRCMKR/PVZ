package models

import (
	"log"

	"github.com/Rhymond/go-money"
	"github.com/bytedance/sonic"
)

// PackagingType ...
type PackagingType uint

const (
	// NoPackaging ...
	NoPackaging PackagingType = iota

	// BagPackaging ...
	BagPackaging

	// BoxPackaging ...
	BoxPackaging

	// WrapPackaging ...
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
	// NoPackagingName ...
	NoPackagingName = "none"

	// BagName ...
	BagName = "bag"

	// BoxName ...
	BoxName = "box"

	// WrapName ...
	WrapName = "wrap"
)

// Packaging ...
type Packaging interface {
	String() string
	GetType() PackagingType
	GetCost() *money.Money
	GetMinWeight() float64
	GetCheckWeight() bool
}

// GetPackaging ...
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

// BasePackaging ...
type BasePackaging struct {
	// Type ...
	Type PackagingType

	// Cost ...
	Cost money.Money

	// MinWeight ...
	MinWeight float64

	// CheckWeight ...
	CheckWeight bool
}

// String ...
func (b *BasePackaging) String() string {
	return GetPackagingName(b.Type)
}

// GetType ...
func (b *BasePackaging) GetType() PackagingType {
	return b.Type
}

// GetCost ...
func (b *BasePackaging) GetCost() *money.Money {
	return &b.Cost
}

// GetMinWeight ...
func (b *BasePackaging) GetMinWeight() float64 {
	return b.MinWeight
}

// GetCheckWeight ...
func (b *BasePackaging) GetCheckWeight() bool {
	return b.CheckWeight
}

// MarshalJSON ...
func (b *BasePackaging) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(b.String())
}

// UnmarshalJSON ...
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

// GetPackagingName ...
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

// Bag ...
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

// Box ...
type Box struct {
	BasePackaging
}

// newBox ...
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

// Wrap ...
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

// None ...
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
