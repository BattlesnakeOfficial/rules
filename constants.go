package rules

type RulesetError string

func (err RulesetError) Error() string { return string(err) }

const (
	MoveUp    = "up"
	MoveDown  = "down"
	MoveRight = "right"
	MoveLeft  = "left"

	BoardSizeSmall   = 7
	BoardSizeMedium  = 11
	BoardSizeLarge   = 19
	BoardSizeXLarge  = 21
	BoardSizeXXLarge = 25

	SnakeMaxHealth = 100
	SnakeStartSize = 3

	// Snake state constants
	NotEliminated                   = ""
	EliminatedByCollision           = "snake-collision"
	EliminatedBySelfCollision       = "snake-self-collision"
	EliminatedByOutOfHealth         = "out-of-health"
	EliminatedByHeadToHeadCollision = "head-collision"
	EliminatedByOutOfBounds         = "wall-collision"
	EliminatedByHazard              = "hazard"

	// Error constants
	ErrorTooManySnakes   = RulesetError("too many snakes for fixed start positions")
	ErrorNoRoomForSnake  = RulesetError("not enough space to place snake")
	ErrorNoRoomForFood   = RulesetError("not enough space to place food")
	ErrorNoMoveFound     = RulesetError("move not provided for snake")
	ErrorZeroLengthSnake = RulesetError("snake is length zero")
	ErrorEmptyRegistry   = RulesetError("empty registry")
	ErrorNoStages        = RulesetError("no stages")
	ErrorStageNotFound   = RulesetError("stage not found")
	ErrorMapNotFound     = RulesetError("map not found")

	// Ruleset / game type names
	GameTypeConstrictor        = "constrictor"
	GameTypeRoyale             = "royale"
	GameTypeSolo               = "solo"
	GameTypeStandard           = "standard"
	GameTypeWrapped            = "wrapped"
	GameTypeWrappedConstrictor = "wrapped_constrictor"

	// Game creation parameter names
	ParamGameType            = "name"
	ParamFoodSpawnChance     = "foodSpawnChance"
	ParamMinimumFood         = "minimumFood"
	ParamHazardDamagePerTurn = "damagePerTurn"
	ParamHazardMap           = "hazardMap"
	ParamHazardMapAuthor     = "hazardMapAuthor"
	ParamShrinkEveryNTurns   = "shrinkEveryNTurns"
	ParamAllowBodyCollisions = "allowBodyCollisions"
	ParamSharedElimination   = "sharedElimination"
	ParamSharedHealth        = "sharedHealth"
	ParamSharedLength        = "sharedLength"
)
