module github.com/Kirill-Znamenskiy/WorldOfWisdom/client

go 1.22

require (
	github.com/Kirill-Znamenskiy/kzlogger v0.0.3
	github.com/Kirill-Znamenskiy/WorldOfWisdom/server v0.0.1
)

replace github.com/Kirill-Znamenskiy/WorldOfWisdom/server v0.0.1 => ./../server

require github.com/Kirill-Znamenskiy/kzerror v0.0.1 // indirect
