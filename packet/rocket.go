package packet

// Direction représente la direction par rapport aux axes/pôles magnétiques
type Direction rune

const (
    North Direction = 'N'
    South Direction = 'S'
    East  Direction = 'E'
    West  Direction = 'W'
)

// Vector représente une coordonnée ou direction/vecteur tridimensionnelle
type Vector struct {
    X float64 `csv:"x"`
    Y float64 `csv:"y"`
    Z float64 `csv:"z"`
}

// State représente un état (actif/inactif)
type State uint

const (
    On  State = 1
    Off State = 0
)

// SystemStates contient l'état des composantes de la fusée
type SystemStates struct {
    AcquisitionBoard1 State `csv:"acquisition_board_state_1"`
    AcquisitionBoard2 State `csv:"acquisition_board_state_2"`
    AcquisitionBoard3 State `csv:"acquisition_board_state_3"`
    PowerSupply1      State `csv:"power_supply_state_1"`
    PowerSupply2      State `csv:"power_supply_state_2"`
    PayloadBoard1     State `csv:"payload_board_state_1"`
}

// RocketPacket représente un paquet de données reçu de la fusée.
//
// NOTE: Les champs doivent être dans la même ordre que les données reçues
type RocketPacket struct {
    // Time est le temps écoulé en millisecondes (ms) depuis le démarrage ?? FIXME
    Time uint64 `csv:"time_stamp"`
    
    // Latitude GPS en degrés
    Latitude float64 `csv:"latitude"`
    
    // Longitude GPS en degrés
    Longitude float64 `csv:"longitude"`
    
    // IndicatorNS est la direction nord/sud
    IndicatorNS Direction `csv:"ns_indicator"`
    
    // IndicatorEW est la direction est/ouest
    IndicatorEW Direction `csv:"ew_indicator"`
    
    // GpsSatellites est le nombre de satellites auxquels le GPS est connecté
    GpsSatellites uint8 `csv:"gps_satellites"`
    
    // Altitude en mètre (m)
    Altitude float64 `csv:"altitude"`
    
    // Pressure est la pression en Pascal (p)
    Pressure float64 `csv:"pressure"`
    
    // Temperature en degrés Celsius (°C)
    Temperature float64 `csv:"temperature"`
    
    // Acceleration est la force gravitationnelle subit en g
    Acceleration Vector `csv:"acceleration_,inline"`
    
    // Magnetism en millitesla (mT)
    Magnetism Vector `csv:"magnetometer_,inline"`
    
    // AngularSpeed est la vitesse angulaire en radian par seconde
    AngularSpeed Vector `csv:"angular_speed_,inline"`
    
    // States est l'état des systèmes
    States SystemStates
}
