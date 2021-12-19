package enum

type AttackOutcome int

const (
	AttackOutcomeInvalid AttackOutcome = iota
	AttackOutcomeHit
	AttackOutcomeMiss
	AttackOutcomeAlreadyHit
	AttackOutcomeHitAndWin
)
