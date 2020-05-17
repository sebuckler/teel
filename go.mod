module github.com/sebuckler/teel

go 1.14

require github.com/sebuckler/teel/pkg/cli v0.0.0-20200517212548-0f0e8b98570f
replace (
	github.com/sebuckler/teel/pkg/cli v0.0.0-20200517212548-0f0e8b98570f => ./pkg/cli
)
