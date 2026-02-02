Chaos crashers is a top down, wave based hack and slash inspired by soul night for andriod and ios. it uses the ebitengine game engine/framework.

![player asset attacking zombie hoard](assets/images/exampleOne.png)
![player asset attacking zombie hoard two](assets/images/exampleTwo.png)
![player picking between two seperate rooms](assets/images/exampleThree.png)

Instructions are for debian/ubuntu based distrobutions and windows. Instructions for other distrobutions can be found on ebitengine website.

debian/ubuntu:

```
sudo apt install gcc

sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

go mod init github.com/name/gamename

go mod tidy

go run .
```

windows 10/11:

```
install GO from go.dev

if using windows subsystem for linux (probably just ignore if you do not know what this means) set the environment table to:

GOOS=windows go run github.com/hajimehoshi/ebiten/v2/examples/rotate@latest

go mod init github.com/yourname/yourgame

go mod tidy

go run .
```

wasd for movement & arrow keys for attack direction or standard controller dpad/joystick for movement & abxy for attack direction.

[^1]:
    shoutout to my bff @crungulus for his help with assets <3
[^1]

