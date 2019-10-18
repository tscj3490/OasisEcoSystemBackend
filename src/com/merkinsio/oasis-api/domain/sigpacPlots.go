package domain

import (
	geojson "github.com/paulmach/go.geojson"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type SigpacPlot struct {
	ID          bson.ObjectId    `bson:"_id" json:"id" valid:"required"`
	GID         int64            `json:"gid" bson:"gid"`
	SigpacID    int64            `json:"sigpacId" bson:"sigpacId"`
	ProvinceRec int64            `json:"provinceRec" bson:"provinceRec"`
	TownRec     int64            `json:"townRec" bson:"townRec"`
	PolygonNum  int64            `json:"polygonNum" bson:"polygonNum"`
	PlotNum     int64            `json:"plotNum" bson:"plotNum"`
	Enclosure   int64            `json:"enclosure" bson:"enclosure"`
	UsageCode   string           `json:"usageCode" bson:"usageCode"`
	Area        float64          `json:"area" bson:"area"`
	Region      string           `json:"region" bson:"region"`
	GC          string           `json:"gc" bson:"gc"`
	Version     string           `json:"version" bson:"version"`
	GeoJSON     geojson.Geometry `json:"geoJSON" bson:"geoJSON"`
}

func GetSigpacPlotByPlotCode(plotCode PlotCodification) (*SigpacPlot, error) {
	var sigpacPlot SigpacPlot
	if err := DB.C("sigpacPlots").Find(bson.M{"provinceRec": plotCode.ProvinceRec, "townRec": plotCode.TownRec, "polygonNum": plotCode.PolygonNum, "plotNum": plotCode.PlotNum, "enclosure": plotCode.Enclosure}).One(&sigpacPlot); err != nil {
		log.Errorf("Error in domain.GetSigpacPlotByPlotCode -> error: %v", err.Error())
		return nil, err
	}

	return &sigpacPlot, nil
}
