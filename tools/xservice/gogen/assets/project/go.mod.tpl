module {{.Module}}

require (
	github.com/xinpianchang/xservice latest
)

replace {{.Module}}_pb => ./pb
