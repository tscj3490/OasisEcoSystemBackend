package domain

import (
	"time"

	validator "github.com/asaskevich/govalidator"
	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// PROVISIONAL NAME = this will change when we know what it is
// Nearly all properties are exported from the Excel TABLA_PARCELAS.xlsx
// It is very likely that some of these properties will be completely useless.

//Plot holds the the plantation's different plots data
type Plot struct {
	Entity            `bson:",inline"`
	Name              string            `bson:"name" json:"name"`
	PlantationID      bson.ObjectId     `bson:"plantationId" json:"plantationId"`
	PlantationName    string            `bson:"plantationName" json:"plantationName"`
	CropsType         string            `bson:"cropsType" json:"cropsType"`
	CropsName         string            `bson:"cropsName" json:"cropsName"`   // CULT_DECL in the Excel
	CropsCode         int               `bson:"cropsCode" json:"cropsCode"`   // CD_CULT_DECL in the Excel
	AmbitoGrab        string            `bson:"ambitoGrab" json:"ambitoGrab"` // PROVISIONAL NAME
	NIF               string            `bson:"nif" json:"nif"`
	Record            string            `bson:"record" json:"record"`                       // EXPE in the Excel
	Surname1          string            `bson:"surname1" json:"surname1"`                   //This is empty on the Excel
	Surname2          string            `bson:"surname2" json:"surname2"`                   //This is empty on the Excel
	DpeNombre         string            `bson:"dpeNombre" json:"dpeNombre"`                 // PROVISIONAL NAME
	UserAdmissionCode string            `bson:"userAdmissionCode" json:"userAdmissionCode"` // CD_USU_ALTA in the Excel
	AdmissionUser     string            `bson:"admissionUser" json:"admissionUser"`         // USU_ALTA in the Excel
	AmbitoGest        string            `bson:"ambitoGest" json:"ambitoGest"`               // PROVISIONAL NAME
	PlotNumber        string            `bson:"plotNumber" json:"plotNumber"`
	PlotCode          PlotCodification  `bson:"plotCode" json:"plotCode"`
	VarietyCode       int               `bson:"varietyCode" json:"varietyCode"`           // CD_VARIEDAD in the Excel
	Variety           string            `bson:"variety" json:"variety"`                   // VARIEDAD in the Excel
	OperatingSystem   string            `bson:"operatingSystem" json:"operatingSystem"`   // SIST_EXPLO in the Excel
	PresAlegSigpac    string            `bson:"presAlegSigpac" json:"presAlegSigpac"`     // PROVISIONAL NAME PRES_ALEG_SIGPAC in the Excel
	RegTenencia       string            `bson:"regTenencia" json:"regTenencia"`           // PROVISIONAL NAME REG_TENENCIA in the Excel
	CapDecl           float64           `bson:"capDecl" json:"capDecl"`                   // PROVISIONAL NAME This is empty on the Excel
	CapSigpac         float64           `bson:"capSigpac" json:"capSigpac"`               // PROVISIONAL NAME This is empty on the Excel
	SupNetaPastos     float64           `bson:"supNetaPastos" json:"supNetaPastos"`       // PROVISIONAL NAME This is empty on the Excel
	SupDecl           float64           `bson:"supDecl" json:"supDecl"`                   // PROVISIONAL NAME
	SupGrafica        float64           `bson:"supGrafica" json:"supGrafica"`             // PROVISIONAL NAME
	AgrarianActivity  string            `bson:"agrarianActivity" json:"agrarianActivity"` // ACTIVIDAD AGRARIA in the Excel
	LessorCIF         string            `bson:"lessorCif" json:"lessorCif"`               // CIF_ARRENDADOR in the Excel
	Shape             []Coordinates     `bson:"shape" json:"shape"`
	SigpacID          string            `bson:"sigpacId" json:"sigpacId"`
	DisplayImage      string            `json:"displayImage" bson:"displayImage"` // TODO: Remove this in favour of issue map with location
	GeoJSON           *geojson.Geometry `json:"geoJSON,omitempty" bson:"geoJSON,omitempty"`
	CropColor         string            `json:"cropColor" bson:"cropColor"`
}

// PlotCodification holds the plot codification
type PlotCodification struct {
	ProvinceRec int `bson:"provinceRec" json:"provinceRec"` // PROVISIONAL NAME PROV_REC in the Excel
	TownRec     int `bson:"townRec" json:"townRec"`         // PROVISIONAL NAME MUN_REC in the Excel
	PolygonNum  int `bson:"polygonNum" json:"polygonNum"`   // POLIGONO in the Excel
	PlotNum     int `bson:"plotNum" json:"plotNum"`         // PARCELA in the Excel
	Enclosure   int `bson:"enclosure" json:"enclosure"`     // PROVISIONAL NAME RECINTO in the Excel
}

//Validate validates the Plot struct
func (plot *Plot) Validate() (bool, error) {
	//Default validation
	_, err := validator.ValidateStruct(plot)
	if err != nil {
		return false, err
	}

	//Custom validations

	return true, err
}

// CreatePlotByPlot creates a new plot with the plot values brought
func CreatePlotByPlot(plot *Plot) (*bson.ObjectId, error) {
	plot.InitializeNewData()

	if _, err := plot.Validate(); err != nil {
		log.Errorf("Error in domain.CreatePlotByPlot.Validate -> error: %s", err.Error())
		return nil, err
	}

	if err := DB.C("plots").Insert(plot); err != nil {
		log.Errorf("Error in domain.CreatePlotByPlot.Insert -> error: %v", err.Error())
		return nil, err
	}

	return &plot.ID, nil
}

//GetPlotByID returns a plot by its given ID
func GetPlotByID(plotID bson.ObjectId) (*Plot, error) {
	var plot Plot
	if err := DB.C("plots").Find(bson.M{"_id": plotID}).One(&plot); err != nil {
		log.Errorf("Error in domain.GetPlotByID -> error: %v", err.Error())
		return nil, err
	}

	return &plot, nil
}

func GetPlotsByPlantationID(plantationID bson.ObjectId) ([]Plot, error) {
	plots := []Plot{}

	if err := DB.C("plots").Find(bson.M{"plantationId": plantationID, "active": true}).All(&plots); err != nil {
		return plots, err
	}

	return plots, nil
}

func DeletePlotByID(plotID bson.ObjectId) error {

	if err := DB.C("plots").UpdateId(plotID, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return err
	}

	return nil
}

func UpdatePlot(plot Plot) error {
	plot.UpdatedAt = time.Now().UTC()
	if err := DB.C("plots").UpdateId(plot.ID, plot); err != nil {
		if err.Error() != "not found" {
			log.Errorf("Error in domain.UpdatePlot -> error: %v", err.Error())
		}
		return err
	}

	return nil
}
