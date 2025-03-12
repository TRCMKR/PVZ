package models

import (
	"log"

	"github.com/Rhymond/go-money"
	"github.com/bytedance/sonic"
)

type PackagingType uint

const (
	NoPackaging PackagingType = iota
	BagPackaging
	BoxPackaging
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
	BagName  = "bag"
	BoxName  = "box"
	WrapName = "wrap"
)

type Packaging interface {
	String() string
	GetType() PackagingType
	GetCost() *money.Money
	GetMinWeight() float64
	GetCheckWeight() bool
}

func GetPackaging(packaging string) Packaging {
	switch packaging {
	case BagName:
		return newBag()
	case BoxName:
		return newBox()
	case WrapName:
		return newWrap()
	default:
		return nil
	}
}

type BasePackaging struct {
	Type        PackagingType
	Cost        money.Money
	MinWeight   float64
	CheckWeight bool
}

func (b *BasePackaging) String() string {
	return GetPackagingName(b.Type)
}

func (b *BasePackaging) GetType() PackagingType {
	return b.Type
}

func (b *BasePackaging) GetCost() *money.Money {
	return &b.Cost
}

func (b *BasePackaging) GetMinWeight() float64 {
	return b.MinWeight
}

func (b *BasePackaging) GetCheckWeight() bool {
	return b.CheckWeight
}

func (b *BasePackaging) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(b.String())
}

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

func GetPackagingName(packaging PackagingType) string {
	switch packaging {
	case BagPackaging:
		return BagName
	case BoxPackaging:
		return BoxName
	case WrapPackaging:
		return WrapName
	default:
		return ""
	}
}

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
