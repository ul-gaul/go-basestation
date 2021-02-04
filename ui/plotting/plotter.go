package plotting

import (
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg/draw"
    "sync"
    "time"
    
    "github.com/ul-gaul/go-basestation/constants"
    "github.com/ul-gaul/go-basestation/utils"
)

var _ plot.Plotter = (*Plotter)(nil)

// Plotter est une structure qui permet de gérer une liste de données pour un graphique
type Plotter struct {
    // Name correspond à la légende
    Name                   string
    // LineStyle est le style des lignes
    LineStyle              draw.LineStyle
    // PointStyleFunc est une fonction retournant le style du point spécifié
    PointStyleFunc         func(int) draw.GlyphStyle
    // DataLimit est la limite de données pouvant être affichées
    DataLimit              int
    xys                    plotter.XYs
    padRatioX, padRatioY   float64
    chChange               chan time.Time
    mut                    sync.Mutex
}

// Data retourne la liste des données
func (p *Plotter) Data() plotter.XYs { return p.xys }

// Prepend ajoute les données spécifiées au début des données existantes
func (p *Plotter) Prepend(xys ...plotter.XY)       { p.PrependAll(xys) }
func (p *Plotter) PrependAll(xys plotter.XYs)      { p.InsertAll(0, xys) }
// Append ajoute les données spécifiées à la fin des données existantes
func (p *Plotter) Append(xys ...plotter.XY)        { p.AppendAll(xys) }
func (p *Plotter) AppendAll(xys plotter.XYs)       { p.InsertAll(len(p.xys), xys) }
// Insert ajoute les données spécifiées à la position spécifiée
func (p *Plotter) Insert(i int, xys ...plotter.XY) { p.InsertAll(i, xys) }
func (p *Plotter) InsertAll(i int, xys plotter.XYs) {
    if len(xys) == 0 {
        return
    }
    p.setXYs(append(p.xys[:i], append(xys, p.xys[i:]...)...))
}

// ReplaceAll remplace toutes les données par celles spécifiées
func (p *Plotter) ReplaceAll(xys plotter.XYs) {
    p.setXYs(xys)
}

func (p *Plotter) setXYs(xys plotter.XYs) {
    p.mut.Lock()
    p.xys = xys
    defer p.mut.Unlock()
    
    select {
    // notifie qu'il y a un changement de données
    case p.chChange <- time.Now():
    default:
    }
}

// PaddingX returns the padding ratio for the X axis.
func (p *Plotter) PaddingX() float64 { return p.padRatioX }

// PaddingY returns the padding ratio for the Y axis.
func (p *Plotter) PaddingY() float64 { return p.padRatioY }

// SetPaddingX sets the padding ratio for the X axis.
// Must be between -1 and 1;
func (p *Plotter) SetPaddingX(padding float64) error { return p.setPadding(padding, p.padRatioY) }

// SetPaddingX sets the data padding for the Y axis.
// Must be between -1 and 1;
func (p *Plotter) SetPaddingY(padding float64) error { return p.setPadding(p.padRatioX, padding) }

func (p *Plotter) setPadding(x, y float64) error {
    if x < -1 || x > 1 || y < -1 || y > 1 {
        return constants.ErrPaddingOutOfRange
    }
    p.padRatioX, p.padRatioY = x, y
    return nil
}

// Plot implements the plot.Plotter interface
func (p *Plotter) Plot(c draw.Canvas, plt *plot.Plot) {
    var xys plotter.XYs
    
    // Crée un copie des données à afficher
    p.mut.Lock()
    if p.DataLimit <= 0 || p.xys.Len() < p.DataLimit {
        xys = p.xys[:]
    } else {
        xys = p.xys[p.xys.Len()-p.DataLimit:]
    }
    p.mut.Unlock()
    
    // Crée une nouvelle ligne avec points représentant les points à afficher
    line, points, err := plotter.NewLinePoints(xys)
    utils.CheckErr(err)
    
    // Applique les styles
    line.LineStyle = p.LineStyle
    points.GlyphStyleFunc = p.PointStyleFunc
    
    // Recalcule les axes
    xmin, xmax, ymin, ymax := utils.FindMinMax(xys...)
    xpad, ypad := p.padRatioX*(xmax-xmin), p.padRatioY*(ymax-ymin)
    plt.X.Min, plt.X.Max = xmin-xpad, xmax+xpad
    plt.Y.Min, plt.Y.Max = ymin-ypad, ymax+ypad
    
    // Dessine le graphique sur le canvas c
    line.Plot(c, plt)
    points.Plot(c, plt)
}
