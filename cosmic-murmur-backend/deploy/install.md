
## Pi Install

- Install Go
- Install libraries
    - libx11-dev, xvfb, libgl1-mesa-dev, cmake, xorg-dev
- Add go path to sudo
    - visudo, add secure_path="...:/usr/local/go/bin"
- Install github.com/go-gl/gl
    - sudo go get -u github.com/go-gl/gl/v3.1/gles2
- Install github.com/go-gl/glfw
    - sudo go get -u -tags=gles2 github.com/go-gl/glfw/v3.3/glfw
- Build
    - sudo go build ./cmd/runApplication/main.go
- Make Service
    - https://superuser.com/questions/544399/how-do-you-make-a-systemd-service-as-the-last-service-on-boot