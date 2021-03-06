#!/bin/bash
#
# Copyright 2020 Aletheia Ware LLC
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


fyne bundle -name SpaceIcon -package data space.svg > icon.go
fyne bundle -append -name CameraPhotoIcon -package data camera_photo.svg >> icon.go
fyne bundle -append -name CameraVideoIcon -package data camera_video.svg >> icon.go
fyne bundle -append -name MicrophoneIcon -package data microphone.svg >> icon.go
fyne bundle -append -name StorageIcon -package data storage.svg >> icon.go
#fyne bundle -append -name XYZIcon -package data xyz.svg >> icon.go
