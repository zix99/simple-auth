package db

type ManualStipulation struct {
}

func (s *ManualStipulation) Type() StipulationType {
	return StipulationType("manual")
}

func (s *ManualStipulation) IsSatisfiedBy(spec IStipulation) bool {
	return true
}
