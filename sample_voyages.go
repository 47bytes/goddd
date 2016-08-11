package goddd

// A set of sample voyages.
var (
	V100 = NewVoyage("V100", Schedule{
		[]CarrierMovement{
			{DepartureLocation: Hongkong, ArrivalLocation: Tokyo},
			{DepartureLocation: Tokyo, ArrivalLocation: NewYork},
		},
	})

	V300 = NewVoyage("V300", Schedule{
		[]CarrierMovement{
			{DepartureLocation: Tokyo, ArrivalLocation: Rotterdam},
			{DepartureLocation: Rotterdam, ArrivalLocation: Hamburg},
			{DepartureLocation: Hamburg, ArrivalLocation: Melbourne},
			{DepartureLocation: Melbourne, ArrivalLocation: Tokyo},
		},
	})

	V400 = NewVoyage("V400", Schedule{
		[]CarrierMovement{
			{DepartureLocation: Hamburg, ArrivalLocation: Stockholm},
			{DepartureLocation: Stockholm, ArrivalLocation: Helsinki},
			{DepartureLocation: Helsinki, ArrivalLocation: Hamburg},
		},
	})
)

// These voyages are hard-coded into the current pathfinder. Make sure
// they exist.
var (
	V0100S = NewVoyage("0100S", Schedule{[]CarrierMovement{}})
	V0200T = NewVoyage("0200T", Schedule{[]CarrierMovement{}})
	V0300A = NewVoyage("0300A", Schedule{[]CarrierMovement{}})
	V0301S = NewVoyage("0301S", Schedule{[]CarrierMovement{}})
	V0400S = NewVoyage("0400S", Schedule{[]CarrierMovement{}})
)
