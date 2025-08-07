package models

type BaseModel struct{}

func MaxTeamsPerPlan(plan string) (maxGroups int) {
	if plan == StarterPlan {
		return 5
	} else if plan == BasicPlan {
		return 10
	} else if plan == PremiumPlan {
		return 20
	} else {
		return 0
	}
}
