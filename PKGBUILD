# Maintainer: S0FTS0RR0W <your@email.com>
pkgname=liquid-ui-git
pkgver=r3.gd87383b
pkgrel=1
pkgdesc="A web-based dashboard for monitoring and controlling liquid coolers"
arch=('x86_64')
url="https://github.com/S0FTS0RR0W/liquid-ui"
license=('MIT')
depends=('liquidctl' 'glibc')
makedepends=('go' 'git' 'npm' 'nodejs')
provides=("${pkgname%-git}")
conflicts=("${pkgname%-git}")
source=("git+https://github.com/S0FTS0RR0W/liquid-ui.git")
sha256sums=('SKIP')

pkgver() {
  cd "$srcdir/${pkgname%-git}"
  printf "r%s.g%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

prepare() {
  # Create systemd service file
  cat <<EOF > liquid-ui-backend.service
[Unit]
Description=Liquid UI Backend
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/liquid-ui-backend
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
EOF

  # Create desktop entry
  cat <<EOF > liquid-ui.desktop
[Desktop Entry]
Version=1.0
Type=Application
Name=Liquid UI
Comment=Liquid Cooler Control Dashboard
Exec=xdg-open http://localhost:8765
Icon=utilities-system-monitor
Terminal=false
Categories=System;Monitor;
EOF
}

build() {
  # Build Frontend
  cd "$srcdir/${pkgname%-git}/frontend"
  npm install
  # Install adapter-static and configure it
  npm install -D @sveltejs/adapter-static
  sed -i 's/@sveltejs\/adapter-auto/@sveltejs\/adapter-static/' svelte.config.js
  npm run build

  # Move frontend build to backend for embedding
  rm -rf ../backend/cmd/server/dist
  cp -r build ../backend/cmd/server/dist

  # Build Backend
  cd "$srcdir/${pkgname%-git}/backend/cmd/server"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -o liquid-ui-backend .
}

package() {
  # Install backend binary
  install -Dm755 "$srcdir/${pkgname%-git}/backend/cmd/server/liquid-ui-backend" "$pkgdir/usr/bin/liquid-ui-backend"

  # Install systemd service
  install -Dm644 liquid-ui-backend.service "$pkgdir/usr/lib/systemd/system/liquid-ui-backend.service"

  # Install desktop file
  install -Dm644 liquid-ui.desktop "$pkgdir/usr/share/applications/liquid-ui.desktop"
}
