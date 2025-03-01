package order

import (
	"errors"
	"log"

	"github.com/Rhymond/go-money"
	"github.com/bytedance/sonic"
)

type packagingType uint

const (
	noPackaging packagingType = iota
	bagPackaging
	boxPackaging
	wrapPackaging
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
	bagName  = "bag"
	boxName  = "box"
	wrapName = "wrap"
)

var (
	errNotEnoughWeight = errors.New("not enough weight")
	errWrongPackaging  = errors.New("wrong packaging")
)

type Packaging interface {
	String() string
	Pack(*Order) error
}

func GetPackaging(packaging string) Packaging {
	switch packaging {
	case bagName:
		return newBag()
	case boxName:
		return newBox()
	case wrapName:
		return newWrap()
	default:
		return nil
	}
}

type BasePackaging struct {
	Type        packagingType
	Cost        money.Money
	MinWeight   float64
	CheckWeight bool
}

func (b *BasePackaging) String() string {
	return getPackagingName(b.Type)
}

func (b *BasePackaging) Pack(order *Order) error {
	if b.CheckWeight && order.Weight < b.MinWeight {
		return errNotEnoughWeight
	}

	if order.Packaging == wrapPackaging || order.ExtraPackaging != noPackaging {
		return errWrongPackaging
	}

	if order.Packaging == noPackaging {
		order.Packaging = b.Type
	} else {
		if b.Type != wrapPackaging {
			return errWrongPackaging
		}
		order.ExtraPackaging = b.Type
	}

	tmp, err := order.Price.Add(&b.Cost)
	if err != nil {
		return err
	}
	order.Price = *tmp

	return nil
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

func getPackagingName(packaging packagingType) string {
	switch packaging {
	case bagPackaging:
		return bagName
	case boxPackaging:
		return boxName
	case wrapPackaging:
		return wrapName
	default:
		return ""
	}
}

func formPackagingString(packaging packagingType, extraPackaging packagingType) string {
	var result string
	if packaging == noPackaging {
		result = "none"
	} else {
		result = getPackagingName(packaging)
	}

	if extraPackaging != noPackaging {
		result += " in " + getPackagingName(extraPackaging)
	}

	return result
}

type Bag struct {
	BasePackaging
}

func newBag() Packaging {
	return &Bag{
		BasePackaging{
			Type:        bagPackaging,
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
			Type:        boxPackaging,
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
			Type:        wrapPackaging,
			Cost:        *wrapCost,
			MinWeight:   0,
			CheckWeight: false,
		},
	}
}
