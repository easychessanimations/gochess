module github.com/easychessanimations/gochess

go 1.14

replace bitbucket.org/zurichess/board v1.0.0 => c:/gomodules/modules/gochess/zurichessboard

replace bitbucket.org/zurichess/zurichess v0.0.0-20181124230012-ee20f164ad4a => c:/gomodules/modules/gochess/zurichess

require bitbucket.org/zurichess/zurichess v0.0.0-20181124230012-ee20f164ad4a // indirect
