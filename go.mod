module github.com/easychessanimations/gochess

go 1.14

replace bitbucket.org/zurichess/board v1.0.0 => c:/gomodules/modules/gochess/zurichessboard

replace bitbucket.org/zurichess/zurichess v0.0.0-20181124230012-ee20f164ad4a => c:/gomodules/modules/gochess/zurichess

require (
	bitbucket.org/zurichess/zurichess v0.0.0-20181124230012-ee20f164ad4a // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/motemen/go-quickfix v0.0.0-20200118031250-2a6e54e79a50 // indirect
	github.com/motemen/gore v0.5.0 // indirect
	github.com/peterh/liner v1.2.0 // indirect
	golang.org/x/tools v0.0.0-20200410194907-79a7a3126eef // indirect
)
