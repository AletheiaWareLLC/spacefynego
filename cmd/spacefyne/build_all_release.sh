#!/bin/bash
#
# Copyright 2020-2021 Aletheia Ware LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -x

(cd $GOPATH/src/aletheiaware.com/spacefynego/ui/data/ && ./gen.sh)
go fmt $GOPATH/src/aletheiaware.com/spacefynego/...
go vet $GOPATH/src/aletheiaware.com/spacefynego/...
go test $GOPATH/src/aletheiaware.com/spacefynego/...
fyne-cross android -app-id com.aletheiaware.space -app-version 1.1.3 -icon ./ui/data/logo.png -keystore=./private/SPACE.keystore -output SPACE -release ./cmd/spacefyne/
fyne-cross darwin -app-id com.aletheiaware.space -app-version 1.1.3 -icon ./ui/data/logo.png -output SPACE -release ./cmd/spacefyne/
fyne-cross ios -app-id com.aletheiaware.space -app-version 1.1.3 -icon ./ui/data/logo.png -output SPACE -release ./cmd/spacefyne/
fyne-cross linux -app-id com.aletheiaware.space -app-version 1.1.3 -icon ./ui/data/logo.png -output space -release ./cmd/spacefyne/
fyne-cross windows -app-id com.aletheiaware.space -app-version 1.1.3 -icon ./ui/data/logo.png -output space -release ./cmd/spacefyne/