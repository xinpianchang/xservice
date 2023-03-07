module {{.Module}}

require (
	github.com/xinpianchang/xservice/v2 v2.0.0
)

replace {{.Module}}_pb => ./pb
