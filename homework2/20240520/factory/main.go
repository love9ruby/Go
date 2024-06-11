package main

type IGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
	getSound() string
}

type Gun struct {
	name  string
	power int
}

func (g *Gun) setName(name string) {
	g.name = name
}

func (g *Gun) getName() string {
	return g.name
}

func (g *Gun) setPower(power int) {
	g.power = power
}

func (g *Gun) getPower() int {
	return g.power
}

type Ak47 struct {
	Gun          // embed Gun, Ak47 is a Gun, so it has all the Gun's methods
	sound string // Ak47 has an additional field
}

// rewrite old method for Ak47
func (g *Ak47) getName() string {
	return g.name + "!"
}

// write less method for interface
func (g *Ak47) getSound() string {
	return g.sound
}

// M16 is a concrete product
type M16 struct {
	Gun // embed Gun
}

// write less method for interface
func (g *M16) getSound() string {
	return "Tak Tak Tak..."
}

type GunFactory struct{}

// MakeGun is a factory method
func (g *GunFactory) MakeGun(gunType string) IGun {
	switch gunType {
	case "ak47":
		return &Ak47{Gun{"AK47", 4}, "Bang Bang Bang..."}
	case "m16":
		return &M16{Gun{"M16", 5}}
	default:
		return nil
	}
}

func main() {
	gunFactory := &GunFactory{}

	ak47 := gunFactory.MakeGun("ak47")
	m16 := gunFactory.MakeGun("m16")
	//gun := Gun{}
	//gun.setName("A-47")
	//gun.setPower(100)
	//fmt.Println("Hollo, World!")
	println(ak47.getName())
	println(ak47.getPower())
	println(m16.getName())
	println(m16.getPower())
}
